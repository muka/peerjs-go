package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

// return time in millis
// credits https://stackoverflow.com/questions/24122821/go-golang-time-now-unixnano-convert-to-milliseconds
func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func createLogger(ctx string, opts Options) *logrus.Entry {
	logger := logrus.New()
	level, err := logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		logger.Fatalf("Cannot parse log level %s", opts.LogLevel)
	}
	logger.SetLevel(level)

	return logger.WithFields(logrus.Fields{
		"context": ctx,
	})
}
