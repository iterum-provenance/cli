package git

import (
	"github.com/Mantsje/iterum-cli/config"
)

func ensureGitDeps(platform config.GitPlatform) {
	config.EnsureDep(config.GitDep)
	switch platform {
	case config.Github:
		config.EnsureDep(config.GithubDep)
	case config.Gitlab:
		config.EnsureDep(config.GitlabDep)
	case config.Bitbucket:
		config.EnsureDep(config.BitbucketDep)
	case config.None:
		config.EnsureDep(config.GitDep)
	}
}
