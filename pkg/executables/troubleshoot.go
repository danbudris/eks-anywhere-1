package executables

import (
	"context"
	"fmt"
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
	params := []string{bundlePath, "--kubeconfig", kubeconfig, "--interactive=false"}
	output, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
