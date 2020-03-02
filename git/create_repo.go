package git

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/util"
)

func createGithubRepo(path string) url.URL {
	var output string
	hub := exec.Command("hub", "create")
	output = util.RunCommand(hub, path)
	rawURL := strings.Split(output, "\n")[1]
	uri, err := url.Parse(rawURL)
	if err != nil {
		log.Print("parsing repo URL failed, most likely due to updated `hub` package, create an issue for iterum-cli")
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

// CreateRepo inits, adds, commits, possibly creates remote, and pushes a repo to target platform
// returns the url that you can find the repo at eventually
func CreateRepo(message string, platform config.GitPlatform, relPath string) url.URL {
	ensureGitDeps(platform)

	path := os.Getenv("PWD") + "/" + relPath
	init := exec.Command("git", "init")
	add := exec.Command("git", "add", ".")
	commit := exec.Command("git", "commit", "-m", message)

	util.RunCommand(init, path)
	util.RunCommand(add, path)
	util.RunCommand(commit, path)

	var uri url.URL
	switch platform {
	case config.Github:
		uri = createGithubRepo(path)
	case config.Gitlab:
		uri = createGitlabRepo(path)
	case config.Bitbucket:
		uri = createBitbucketRepo(path)
	case config.None:
		u, _ := url.Parse("")
		uri = *u
	}
	if platform != config.None {
		push := exec.Command("git", "push", "-u", "origin", "HEAD")
		util.RunCommand(push, path)
	}
	return uri
}
