package executables

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	troulbeshootPath = "support-bundle"
)

type Troubleshoot struct {
	executable Executable
}

func NewTroubleshoot(executable Executable) *Troubleshoot {
	return &Troubleshoot{
		executable: executable,
	}
}

func (t *Troubleshoot) Collect(ctx context.Context, bundlePath string, kubeconfig string) (archivePath string, err error) {
	params := []string{bundlePath, "--kubeconfig", kubeconfig, "--interactive=false"}
	output, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return "", fmt.Errorf("error when executing support-bundle: %s", err)
	}
	archivePath, err = parseCollectOutput(output.String())
	if err != nil {
		return "", fmt.Errorf("error when parsing support-bundle output: %v", err)
	}
	return archivePath, nil
}

func (t *Troubleshoot) Analyze(ctx context.Context, bundleSpecPath string, archivePath string) ([]*SupportBundleAnalysis, error) {
	params := []string{"analyze", bundleSpecPath, "--bundle", archivePath}
	output, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("error when analyzing support bundle %s with analyzers %s", archivePath, bundleSpecPath)
	}
	var analysisOutput []*SupportBundleAnalysis
	err = json.Unmarshal(output.Bytes(), &analysisOutput)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling support-bundle analyze output: %s", err)
	}
	return analysisOutput, err
}

func parseCollectOutput(tsLogs string) (archivePath string, err error) {
	// output parsing logic to be modified once upstream PR to make output more machine-readable is completed
	// https://github.com/replicatedhq/troubleshoot/pull/419
	logEnd := "]"
	logsEndIndex := strings.Index(tsLogs, logEnd) + 1
	archivePath = tsLogs[logsEndIndex:]
	return archivePath, nil
}

type SupportBundleAnalysis struct {
	IsPass  bool   `json:"isPass"`
	IsFail  bool   `json:"isFail"`
	IsWarn  bool   `json:"isWarn"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Uri     string `json:"URI"`
}
