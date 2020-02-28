package cmd

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
)

func createGithubRepo(path string) url.URL {
	var output string
	hub := exec.Command("hub", "create")
	output = runCommand(hub, path)
	rawURL := strings.Split(output, "\n")[1]
	uri, err := url.Parse(rawURL)
	if err != nil {
		log.Print("parsing repo URL failed, most likely fault of updated `hub` package, create an issue for us")
		log.Fatal(err)
	}
	return *uri
}

func createGitlabRepo(path string) url.URL {
	log.Fatal("Gitlab repo creation not implemented yet")
	uri, _ := url.Parse("www.test.com")
	return *uri
}

func createBitbucketRepo(path string) url.URL {
	log.Fatal("Bitbucket repo creation not implemented yet")
	uri, _ := url.Parse("www.test.com")
	return *uri
}

// Inits, adds, commits, creates remote and pushes a repo to target platform
// returns the url that you can find the repo at eventually
func createRepo(message string, platform git.GitPlatform, relPath string) url.URL {
	ensureDeps(platform)

	path := os.Getenv("PWD") + "/" + relPath
	init := exec.Command("git", "init")
	add := exec.Command("git", "add", ".")
	commit := exec.Command("git", "commit", "-m", message)

	runCommand(init, path)
	runCommand(add, path)
	runCommand(commit, path)

	var uri url.URL
	switch platform {
	case git.Github:
		uri = createGithubRepo(path)
	case git.Gitlab:
		uri = createGitlabRepo(path)
	case git.Bitbucket:
		uri = createBitbucketRepo(path)
	case git.None:
		u, _ := url.Parse("")
		uri = *u
	}
	if platform != git.None {
		push := exec.Command("git", "push", "-u", "origin", "HEAD")
		runCommand(push, path)
	}
	return uri
}

func ensureDeps(platform git.GitPlatform) {
	config.EnsureDep(config.GitDep)
	switch platform {
	case git.Github:
		config.EnsureDep(config.GithubDep)
	case git.Gitlab:
		config.EnsureDep(config.GitlabDep)
	case git.Bitbucket:
		config.EnsureDep(config.BitbucketDep)
	case git.None:
		config.EnsureDep(config.GitDep)
	}
}
