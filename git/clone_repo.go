package git

import (
	"os"
	"os/exec"

	"github.com/Mantsje/iterum-cli/util"
)

// CloneRepo clones a previously unexistent git repo based on a config
func CloneRepo(protocol Protocol, platform Platform, repoPath string, relativePath string) {
	ensureGitDeps(None)

	cloneURI := ""
	switch protocol {
	case SSH:
		cloneURI += "git@"
		cloneURI += platform.String() + ".com"
		cloneURI += ":"
	case HTTPS:
		cloneURI += "https://"
		cloneURI += platform.String() + ".com"
		cloneURI += "/"
	}
	cloneURI += repoPath

	path := os.Getenv("PWD")
	clone := exec.Command("git", "clone", cloneURI, relativePath)
	util.RunCommand(clone, path, false)
}
