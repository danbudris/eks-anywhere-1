package framework

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/filewriter"
	"github.com/aws/eks-anywhere/pkg/git"
	gitFactory "github.com/aws/eks-anywhere/pkg/git/factory"
	"github.com/aws/eks-anywhere/pkg/git/gitclient"
)

type GitOptions struct {
	Git    git.Provider
	Writer filewriter.FileWriter
}

func (e *ClusterE2ETest) NewGitOptions(ctx context.Context, cluster *v1alpha1.Cluster, fluxConfig *v1alpha1.FluxConfig, writer filewriter.FileWriter, repoPath string) (*GitOptions, error) {
	if fluxConfig == nil {
		return nil, nil
	}

	var localGitRepoPath string
	if repoPath == "" {
		localGitRepoPath = filepath.Join(cluster.Name, "git", fluxConfig.Spec.Github.Repository)
	} else {
		localGitRepoPath = repoPath
	}

	gogitOptions := gitclient.Options{
		RepositoryDirectory: localGitRepoPath,
	}
	goGit := gitclient.New(gogitOptions)

	gitProviderFactory := gitFactory.New(goGit, writer)
	gitProvider, gitwriter, err := gitProviderFactory.BuildProvider(ctx, &fluxConfig.Spec)
	if err != nil {
		return nil, fmt.Errorf("creating Git provider: %v", err)
	}
	err = gitProvider.Validate(ctx)
	if err != nil {
		return nil, err
	}
	return &GitOptions{
		Git:    gitProvider,
		Writer: gitwriter,
	}, nil
}
