package framework

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/aws/eks-anywhere/pkg/git"
	"github.com/aws/eks-anywhere/pkg/git/gogithub"
	"github.com/aws/eks-anywhere/pkg/git/providers/github"
	"path/filepath"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/filewriter"
	gitFactory "github.com/aws/eks-anywhere/pkg/git/factory"
)

func (e *ClusterE2ETest) NewGitTools(ctx context.Context, cluster *v1alpha1.Cluster, fluxConfig *v1alpha1.FluxConfig, writer filewriter.FileWriter, repoPath string) (*gitFactory.GitTools, error) {
	if fluxConfig == nil {
		return nil, nil
	}

	var localGitWriterPath string
	var localGitRepoPath string
	if repoPath == "" {
		r := e.gitRepoName()
		localGitWriterPath = filepath.Join("git", r)
		localGitRepoPath = filepath.Join(cluster.Name, "git", r)
	} else {
		localGitWriterPath = repoPath
		localGitRepoPath = repoPath
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

func (e *ClusterE2ETest) TestGithubClient(ctx context.Context, githubToken string, owner string, repository string, personal bool) (git.ProviderClient, error) {
		auth := git.TokenAuth{Token: githubToken, Username: owner}
		gogithubOpts := gogithub.Options{Auth: auth}
		githubProviderClient := gogithub.New(ctx, gogithubOpts)

		config := &v1alpha1.GithubProviderConfig{
			Owner:      owner,
			Repository: repository,
			Personal:   personal,
		}
		provider, err := github.New(githubProviderClient, config, auth)
		if err != nil {
			return nil, fmt.Errorf("creating test git provider: %v", err)
		}

		return provider, nil
}
