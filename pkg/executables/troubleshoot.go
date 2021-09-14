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

func (t *Troubleshoot) Version(ctx context.Context) error {
	params := []string{"version"}
	_, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return fmt.Errorf("error executing version: %v", err)
	}
	return nil
}
