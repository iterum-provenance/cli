package git

import (
	"os"
	"os/exec"

	"github.com/iterum-provenance/cli/util"
)

// PullRepo pulls a git repo based on a relative path to the existing repo
func PullRepo(relativePath string) {
	ensureGitDeps(None)

	path := os.Getenv("PWD") + "/" + relativePath
	pull := exec.Command("git", "pull")

	util.RunCommand(pull, path, true)
}
