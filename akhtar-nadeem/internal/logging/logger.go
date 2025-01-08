package logging

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func Init() {
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
