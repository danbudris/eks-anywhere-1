package framework

import (
	"context"
	_ "embed"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/filewriter"
	gitFactory "github.com/aws/eks-anywhere/pkg/git/factory"
)

func (e *ClusterE2ETest) NewGitTools(ctx context.Context, cluster *v1alpha1.Cluster, fluxConfig *v1alpha1.FluxConfig, writer filewriter.FileWriter, p string) (*gitFactory.GitTools, error) {
	if fluxConfig == nil {
		return nil, nil
	}

	var localGitWriterPath string
	var localGitRepoPath string
	if p == "" {
		localGitWriterPath = filepath.Join("git", repoPath(fluxConfig))
		localGitRepoPath = filepath.Join(cluster.Name, "git", repoPath(fluxConfig))
	} else {
		localGitWriterPath = p
		localGitRepoPath = p
	}

	tools, err := gitFactory.Build(ctx, cluster, fluxConfig, writer, gitFactory.WithRepositoryDirectory(localGitRepoPath))
	if err != nil {
		return nil, fmt.Errorf("creating Git provider: %v", err)
	}
	if tools.Provider != nil {
		err = tools.Provider.Validate(ctx)
		if err != nil {
			return nil, err
		}
	}
	gitwriter, err := writer.WithDir(localGitWriterPath)
	if err != nil {
		return nil, fmt.Errorf("creating file writer: %v", err)
	}
	gitwriter.CleanUpTemp()
	tools.Writer = gitwriter
	return tools, nil
}

func repoPath(config *v1alpha1.FluxConfig) string {
	var p string
	if config.Spec.Github != nil {
		p = config.Spec.Github.Repository
	}
	if config.Spec.Git != nil {
		p = path.Base(strings.Trim(config.Spec.Git.RepositoryUrl, filepath.Ext(config.Spec.Git.RepositoryUrl)))
	}
	return p
}
