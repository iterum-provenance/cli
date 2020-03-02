package git

import (
	"net/url"
	"os"
	"os/exec"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/util"
)

// CloneRepo clones a previously unexistent git repo based on a config
func CloneRepo(uri url.URL, relativePath string) {
	ensureGitDeps(config.None)

	path := os.Getenv("PWD")
	clone := exec.Command("git", "clone", uri.String(), relativePath)

	util.RunCommand(clone, path)

}
