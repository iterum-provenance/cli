package git

import "github.com/iterum-provenance/cli/deps"

const (
	bitbucketCmd string = "bitbucket-cli"
	gitCmd       string = "git"
	githubCmd    string = "hub"
	gitlabCmd    string = gitCmd
)

// The dependencies for git related functionality
var (
	GitDep       deps.Dep = deps.Dep{Name: "Git", Cmd: gitCmd}
	GithubDep    deps.Dep = deps.Dep{Name: "Remote Github", Cmd: githubCmd}
	GitlabDep    deps.Dep = deps.Dep{Name: "Remote GitLab", Cmd: gitlabCmd}
	BitbucketDep deps.Dep = deps.Dep{Name: "Remote Bitbucket", Cmd: bitbucketCmd}
)

// adds these dependencies to the global list of dependencies
func init() {
	deps.Register(GitDep, GithubDep, GitlabDep, BitbucketDep)
}

func ensureGitDeps(platform Platform) {
	deps.EnsureDep(GitDep)
	switch platform {
	case Github:
		deps.EnsureDep(GithubDep)
	case Gitlab:
		deps.EnsureDep(GitlabDep)
	case Bitbucket:
		deps.EnsureDep(BitbucketDep)
	case None:
		deps.EnsureDep(GitDep)
	}
}
