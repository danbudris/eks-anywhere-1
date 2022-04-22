package gitfactory

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/filewriter"
	"github.com/aws/eks-anywhere/pkg/git"
	"github.com/aws/eks-anywhere/pkg/git/gitclient"
	"github.com/aws/eks-anywhere/pkg/git/gogithub"
	"github.com/aws/eks-anywhere/pkg/git/providers/github"
)

type gitProviderFactory struct {
	GitClient github.GitClient
	writer    filewriter.FileWriter
}

type Options struct {
	GithubGitClient github.GitClient
	Writer          filewriter.FileWriter
}

func New(client github.GitClient, writer filewriter.FileWriter) *gitProviderFactory {
	return &gitProviderFactory{
		GitClient: client,
		writer:    writer,
	}
}

// BuildProvider will configure and return the proper Github based on the given GitOps configuration.
func (g *gitProviderFactory) BuildProvider(ctx context.Context, fluxConfig *v1alpha1.FluxConfigSpec) (git.Provider, filewriter.FileWriter, error) {
	switch {
	case fluxConfig.Github != nil:
		return g.buildGitHubProvider(ctx, fluxConfig.Github.Repository, fluxConfig.Github.Owner, fluxConfig.Github.Personal)
	case fluxConfig.Git != nil:
		return g.buildGitProvider()
	default:
		return nil, nil, fmt.Errorf("no valid GitOps provider found")
	}
}

func (g *gitProviderFactory) buildGitHubProvider(ctx context.Context, repository string, owner string, personal bool) (git.Provider, filewriter.FileWriter, error) {
	token, err := github.GetGithubAccessTokenFromEnv()
	if err != nil {
		return nil, nil, err
	}

	auth := git.TokenAuth{
		Token:    token,
		Username: owner,
	}

	githubProviderClient := gogithub.New(ctx, auth)
	githubProviderOpts := github.Options{
		Repository: repository,
		Owner:      owner,
		Personal:   personal,
	}
	provider, err := github.New(g.GitClient, githubProviderClient, githubProviderOpts, auth)
	if err != nil {
		return nil, nil, err
	}

	writer, err := g.newWriter(repository)
	if err != nil {
		return nil, nil, err
	}
	writer.CleanUpTemp()

	return provider, writer, nil
}

func (g *gitProviderFactory) buildGitProvider() (git.Provider, filewriter.FileWriter, error) {
	return nil, nil, nil
}

func (g *gitProviderFactory) newWriter(repository string) (filewriter.FileWriter, error) {
	localGitWriterPath := filepath.Join("git", repository)
	return g.writer.WithDir(localGitWriterPath)
}

func (g *gitProviderFactory) NewGitClient(clusterName string, repository string) *gitclient.GitClient {
	localGitRepoPath := filepath.Join(clusterName, "git", repository)
	gitClientOptions := gitclient.Options{
		RepositoryDirectory: localGitRepoPath,
	}
	return gitclient.New(gitClientOptions)
}
