package main

import (
	"github.com/czhj/ahfs/modules/log"
	"go.uber.org/zap"
)

func main() {
	log.Init()
	log.AddLogger("gogo", "console", `{
		"level": "info",
		"stacktracklevel": "warn"
	}`)
	log.New("{}")

	log.Info("hello", zap.String("atom", "beta"))
}
