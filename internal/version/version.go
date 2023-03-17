package version

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

//version.go

var (
	version string = "UNDEFINED"
	commit  string = "UNDEFINED"
	date    string = "UNDEFINED"
	builtBy string = "UNDEFINED"
)

func PrintVersion() {
	fmt.Printf("go-spawn version %s\n", version)
	log.Debugf("{version: %s, commit: %s, date: %s, builtBy: %s}\n", version, commit, date, builtBy)
	os.Exit(0)
}
