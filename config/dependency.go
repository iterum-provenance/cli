package config

import (
	"errors"
	"log"
	"os/exec"
)

// Dep contains all information about a dependency
type Dep struct {
	Cmd  string
	Name string
}

// Enum-like values allowed for dependencies type
const (
	gitCmd       string = "git"
	githubCmd    string = "hub"
	gitlabCmd    string = gitCmd
	bitbucketCmd string = "bitbucket-cli"
	microk8sCmd  string = "microk8s.status"
	dockerCmd    string = "docker"
)

// The dependencies that iterum depends on, used to test whether we can use all functionality
var (
	GitDep       Dep = Dep{Name: "Git", Cmd: gitCmd}
	GithubDep    Dep = Dep{Name: "Remote Github", Cmd: githubCmd}
	GitlabDep    Dep = Dep{Name: "Remote GitLab", Cmd: gitlabCmd}
	BitbucketDep Dep = Dep{Name: "Remote Bitbucket", Cmd: bitbucketCmd}
	Microk8sDep  Dep = Dep{Name: "Microk8s", Cmd: microk8sCmd}
	DockerDep    Dep = Dep{Name: "Docker", Cmd: dockerCmd}
)

// Dependencies holds all the Iterum dependency structs in 1 neat slice
var Dependencies = []Dep{GitDep, GithubDep, GitlabDep, BitbucketDep, Microk8sDep, DockerDep}

// IsValid checks the validity of the RepoType
func (d Dep) IsValid() error {
	switch d {
	case GitDep, GithubDep, GitlabDep, BitbucketDep, Microk8sDep:
		return nil
	}
	return errors.New("Error: Invalid Dep")
}

// IsUsable checks whether the go program can use this dep using os.exec
func (d Dep) IsUsable() bool {
	_, err := exec.LookPath(d.Cmd)
	return err == nil
}

// EnsureDep checks whether a given dependency is met
func EnsureDep(dep Dep) {
	if !dep.IsUsable() {
		log.Fatal(dep.Name, " dependency is not accessible to iterum, run `iterum check` to verify")
	}
}
