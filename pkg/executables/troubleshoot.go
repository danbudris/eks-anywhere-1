package executables

import (
	"context"
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

func (t *Troubleshoot) CollectAndAnalyze(ctx context.Context, bundlePath string, kubeconfig string) error {
	params := []string{bundlePath, "--kubeconfig", kubeconfig}
	_, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return err
	}
	return nil
}
