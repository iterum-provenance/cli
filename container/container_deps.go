package container

import "github.com/iterum-provenance/cli/deps"

// Enum-like values allowed for dependencies type
const (
	microk8sCmd string = "microk8s.status"
	dockerCmd   string = "docker"
)

// The container management related dependencies that iterum depends on
var (
	Microk8sDep deps.Dep = deps.Dep{Name: "Microk8s", Cmd: microk8sCmd}
	DockerDep   deps.Dep = deps.Dep{Name: "Docker", Cmd: dockerCmd}
)

// adds these dependencies to the global list of dependencies
func init() {
	deps.Register(Microk8sDep, DockerDep)
}
