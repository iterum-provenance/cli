package git

import (
	"os"
	"os/exec"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/util"
)

// PullRepo pulls a git repo based on a remote uri
func PullRepo(uri string, relativePath string) {
	ensureGitDeps(config.None)

	path := os.Getenv("PWD") + "/" + relativePath
	init := exec.Command("git", "init")

	util.RunCommand(init, path)

}
