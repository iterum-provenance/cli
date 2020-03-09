package idv

import (
	"os"

	"github.com/prometheus/common/log"
)

func InitIDV() {
	log.Warnln("Still need to pull initial stuff from remote, after that remoge dummy.go")
	os.MkdirAll(localFolder, 744)
	os.MkdirAll(remoteFolder, 744)

}
