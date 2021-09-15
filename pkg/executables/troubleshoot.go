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

func (t *Troubleshoot) CollectAndAnalyze(ctx context.Context, bundlePath string, kubeconfig string) (archivePath string, analysis string, err error) {
	params := []string{bundlePath, "--kubeconfig", kubeconfig, "--interactive=false"}
	output, err := t.executable.Execute(ctx, params...)
	if err != nil {
		return "", "", fmt.Errorf("error when executing support-bundle: %s", err)
	}
	analysis, archivePath, err = parseCollectAndAnalyzeOutputs(output.String())
	if err != nil {
		return "", "", fmt.Errorf("error when parsing support-bundle output: %v", err)
	}
	return analysis, archivePath, nil
}

func parseCollectAndAnalyzeOutputs(tsLogs string) (analysis string, archivePath string, err error) {
	logStart := "["
	logsStartIndex := strings.Index(tsLogs, logStart)

	logEnd := "]"
	logsEndIndex := strings.Index(tsLogs, logEnd) + 1

	analysis = tsLogs[logsStartIndex:logsEndIndex]
	archivePath = tsLogs[logsEndIndex:]

	var i []interface{}
	err = json.Unmarshal([]byte(analysis), &i)
	if err != nil {
		return "", "", err
	}

	fmt.Println(i)
	return analysis, archivePath, nil
}
