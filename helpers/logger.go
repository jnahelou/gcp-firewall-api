package helpers

import (
	"log"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
)

// InitLogger initializes logrus to be compatible with google stackdriver
func InitLogger() {
	logrus.SetFormatter(stackdriver.NewFormatter())
	log.SetOutput(logrus.StandardLogger().Writer())
}