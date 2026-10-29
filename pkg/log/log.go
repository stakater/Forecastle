package log

import (
	"os"
	"sync"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

var (
	sharedLogger *logrus.Logger
	once         sync.Once
)

// New function initialize logrus and return a new logger
// We use an abstraction so that our logs are consistent and if there's anything that needs change
// related to logs, we can just change here
//
// All callers share a single logger. The filename hook keeps its wrapped
// formatter in a package-global variable, so with one logger per package every
// switch between loggers wrapped the formatter once more, and the cost of each
// log line kept growing for the lifetime of the process.
func New() *logrus.Logger {
	once.Do(func() {
		sharedLogger = newLogger()
	})
	return sharedLogger
}

func newLogger() *logrus.Logger {
	// Filename is a hook for logrus that adds file name and line number to the log as well.
	// It's useful for indicating where the log was originated from
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
