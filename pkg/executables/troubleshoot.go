package executables

import (
	"bytes"
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

func (t *Troubleshoot) Analyze(ctx context.Context, bundlePath string) (bytes.Buffer, error) {
	params := []string{"analyze", "--bundle", bundlePath, "--output", "json"}
	output, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return output, nil
}
