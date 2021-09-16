package supportbundle

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"
	"github.com/replicatedhq/troubleshoot/pkg/k8sutil"
	"github.com/replicatedhq/troubleshoot/pkg/supportbundle"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/cluster"
	"github.com/aws/eks-anywhere/pkg/filewriter"
	"github.com/aws/eks-anywhere/pkg/logger"
)

type EksaDiagnosticBundleOpts struct {
	AnalyzerFactory  AnalyzerFactory
	BundlePath       string
	CollectorFactory CollectorFactory
	Client           BundleClient
	Kubeconfig       string
	Writer           filewriter.FileWriter
}

type EksaDiagnosticBundle struct {
	bundle           *v1beta2.SupportBundle
	bundlePath       string
	Spec             *cluster.Spec
	analyzerFactory  AnalyzerFactory
	collectorFactory CollectorFactory
	client           BundleClient
	kubeconfig       string
	Writer           filewriter.FileWriter
}

func NewDiagnosticBundle(clusterSpec *cluster.Spec, opts EksaDiagnosticBundleOpts) (*EksaDiagnosticBundle, error) {
	if opts.BundlePath == "" {
		// user did not provide any bundle-config to the support-bundle command, generate one using the default collectors & analyzers
		return NewBundleFromSpec(clusterSpec, opts), nil
	}
	return NewCustomBundleConfig(opts), nil
}

func NewCustomBundleConfig(opts EksaDiagnosticBundleOpts) *EksaDiagnosticBundle {
	return &EksaDiagnosticBundle{
		bundlePath:       opts.BundlePath,
		analyzerFactory:  opts.AnalyzerFactory,
		collectorFactory: opts.CollectorFactory,
		client:           opts.Client,
		kubeconfig:       opts.Kubeconfig,
	}
}

func NewDefaultBundleConfig(af AnalyzerFactory, cf CollectorFactory) *EksaDiagnosticBundle {
	b := &EksaDiagnosticBundle{
		bundle: &v1beta2.SupportBundle{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SupportBundle",
				APIVersion: "troubleshoot.sh/v1beta2",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "defaultBundle",
			},
			Spec: v1beta2.SupportBundleSpec{},
		},
		analyzerFactory:  af,
		collectorFactory: cf,
	}
	return b.WithDefaultAnalyzers().WithDefaultCollectors()
}

func NewBundleFromSpec(spec *cluster.Spec, opts EksaDiagnosticBundleOpts) *EksaDiagnosticBundle {
	b := &EksaDiagnosticBundle{
		bundle: &v1beta2.SupportBundle{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SupportBundle",
				APIVersion: "troubleshoot.sh/v1beta2",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%sBundle", spec.Name),
			},
			Spec: v1beta2.SupportBundleSpec{},
		},
		analyzerFactory:  opts.AnalyzerFactory,
		collectorFactory: opts.CollectorFactory,
		client:           opts.Client,
		kubeconfig:       opts.Kubeconfig,
	}
	return b.
		WithGitOpsConfig(spec.GitOpsConfig).
		WithOidcConfig(spec.OIDCConfig).
		WithExternalEtcd(spec.Spec.ExternalEtcdConfiguration).
		WithDatacenterConfig(spec.Spec.DatacenterRef).
		WithDefaultAnalyzers().
		WithDefaultCollectors()
}

func (e *EksaDiagnosticBundle) CollectAndAnalyze(ctx context.Context) error {
	archivePath, err := e.client.Collect(ctx, e.bundlePath, e.kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to collect and analyze support bundle: %v", err)
	}
	fmt.Println(archivePath)
	analysis, err := e.client.Analyze(ctx, e.bundlePath, archivePath)
	if err != nil {
		return err
	}
	fmt.Println(analysis[0])
	analysisJson, err := json.Marshal(analysis)
	if err != nil {
		return err
	}
	fmt.Println(string(analysisJson))
	return nil
}

func (e *EksaDiagnosticBundle) CollectBundleFromSpec(sinceTimeValue *time.Time) (string, error) {
	k8sConfig, err := k8sutil.GetRESTConfig()
	if err != nil {
		return "", fmt.Errorf("failed to convert kube flags to rest config: %v", err)
	}

	progressChan := make(chan interface{})
	go func() {
		var lastMsg string
		for {
			msg := <-progressChan
			switch msg := msg.(type) {
			case error:
				logger.Info(fmt.Sprintf("\r * %v", msg))
			case string:
				if lastMsg != msg {
					logger.Info(fmt.Sprintf("\r \033[36mCollecting support bundle\033[m %v", msg))
					lastMsg = msg
				}
			}
		}
	}()

	collectorCB := func(c chan interface{}, msg string) {
		c <- msg
	}
	additionalRedactors := &v1beta2.Redactor{}
	createOpts := supportbundle.SupportBundleCreateOpts{
		CollectorProgressCallback: collectorCB,
		KubernetesRestConfig:      k8sConfig,
		ProgressChan:              progressChan,
		SinceTime:                 sinceTimeValue,
	}

	archivePath, err := supportbundle.CollectSupportBundleFromSpec(&e.bundle.Spec, additionalRedactors, createOpts)
	if err != nil {
		return "", err
	}
	return archivePath, nil
}

func (e *EksaDiagnosticBundle) PrintBundleConfig() error {
	bundleYaml, err := yaml.Marshal(e.bundle)
	if err != nil {
		return fmt.Errorf("error outputting yaml: %v", err)
	}
	fmt.Println(string(bundleYaml))
	return nil
}

func (e *EksaDiagnosticBundle) WithDefaultCollectors() *EksaDiagnosticBundle {
	e.bundle.Spec.Collectors = append(e.bundle.Spec.Collectors, e.collectorFactory.DefaultCollectors()...)
	return e
}

func (e *EksaDiagnosticBundle) WithDefaultAnalyzers() *EksaDiagnosticBundle {
	e.bundle.Spec.Analyzers = append(e.bundle.Spec.Analyzers, e.analyzerFactory.DefaultAnalyzers()...)
	return e
}

func (e *EksaDiagnosticBundle) WithDatacenterConfig(config v1alpha1.Ref) *EksaDiagnosticBundle {
	e.bundle.Spec.Analyzers = append(e.bundle.Spec.Analyzers, e.analyzerFactory.DataCenterConfigAnalyzers(config)...)
	return e
}

func (e *EksaDiagnosticBundle) WithOidcConfig(config *v1alpha1.OIDCConfig) *EksaDiagnosticBundle {
	if config != nil {
		e.bundle.Spec.Analyzers = append(e.bundle.Spec.Analyzers, e.analyzerFactory.EksaOidcAnalyzers()...)
	}
	return e
}

func (e *EksaDiagnosticBundle) WithExternalEtcd(config *v1alpha1.ExternalEtcdConfiguration) *EksaDiagnosticBundle {
	if config != nil {
		e.bundle.Spec.Analyzers = append(e.bundle.Spec.Analyzers, e.analyzerFactory.EksaExternalEtcdAnalyzers()...)
	}
	return e
}

func (e *EksaDiagnosticBundle) WithGitOpsConfig(config *v1alpha1.GitOpsConfig) *EksaDiagnosticBundle {
	if config != nil {
		e.bundle.Spec.Analyzers = append(e.bundle.Spec.Analyzers, e.analyzerFactory.EksaGitopsAnalyzers()...)
	}
	return e
}

func ParseTimeOptions(since string, sinceTime string) (*time.Time, error) {
	var sinceTimeValue time.Time
	var err error
	if sinceTime == "" && since == "" {
		return &sinceTimeValue, nil
	} else if sinceTime != "" && since != "" {
		return nil, fmt.Errorf("at most one of `sinceTime` or `since` could be specified")
	} else if sinceTime != "" {
		sinceTimeValue, err = time.Parse(time.RFC3339, sinceTime)
		if err != nil {
			return nil, fmt.Errorf("unable to parse --since-time option: %v", err)
		}
	} else if since != "" {
		duration, err := time.ParseDuration(since)
		if err != nil {
			return nil, fmt.Errorf("unable to parse --since option: %v", err)
		}
		now := time.Now()
		sinceTimeValue = now.Add(0 - duration)
	}
	return &sinceTimeValue, nil
}
