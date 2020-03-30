package helpers

import (
	"log"
	"os"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
)

// InitLogger initializes logrus to be compatible with google stackdriver
func InitLogger() {
	// Set logger as Stackdriver compliant when runtime is GCP
	// https://cloud.google.com/run/docs/reference/container-contract#env-vars
	if os.Getenv("K_SERVICE") != "" {
		logrus.SetFormatter(stackdriver.NewFormatter())
		log.SetOutput(logrus.StandardLogger().Writer())
	}
}
