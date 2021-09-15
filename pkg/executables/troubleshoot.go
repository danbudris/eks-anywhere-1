package executables

import (
	"context"
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

func (t *Troubleshoot) CollectAndAnalyze(ctx context.Context, bundlePath string, kubeconfig string) (archivePath string, analysis string, err error) {
	params := []string{bundlePath, "--kubeconfig", kubeconfig, "--interactive=false"}
	output, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return "", "", err
	}
	analysis, archivePath = parseCollectAndAnalyzeOutputs(output.String())
	return analysis, archivePath, nil
}

func parseCollectAndAnalyzeOutputs(tsLogs string) (analysis string, archivePath string) {
	logStart := "logs["
	logsStartIndex := strings.Index(tsLogs, logStart)

	logEnd := "]"
	logsEndIndex := strings.Index(tsLogs, logEnd) + 1

	analysis = tsLogs[logsStartIndex:logsEndIndex]
	archivePath = tsLogs[logsEndIndex:]

	return archivePath, analysis
}
