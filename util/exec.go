package util

import (
	"fmt"
	"log"
	"os/exec"
)

// RunCommand runs an arbitrary os/exec.Cmd command as if you were in a terminal at the given path
func RunCommand(cmd *exec.Cmd, path string) string {
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(out))
	return string(out)
}
