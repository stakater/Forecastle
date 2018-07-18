package log

import (
	"os"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

// New function initialize logrus and return a new logger
func New() *logrus.Logger {
	filenameHook := filename.NewHook()

	log := &logrus.Logger{
		Hooks:     make(logrus.LevelHooks),
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{},
		Level:     logrus.InfoLevel,
	}

	log.Hooks.Add(filenameHook)

	return log
}
