package peer

import (
	"github.com/sirupsen/logrus"
)

func createLogger(source string, debugLevel int8) *logrus.Entry {

	log := logrus.New()

	// 0 Prints no logs.
	// 1 Prints only errors.
	// 2 Prints errors and warnings.
	// 3 Prints all logs.
	switch debugLevel {
	case 0:
		log.SetLevel(logrus.PanicLevel)
		break
	case 1:
		log.SetLevel(logrus.ErrorLevel)
		break
	case 2:
		log.SetLevel(logrus.WarnLevel)
		break
	default:
		log.SetLevel(logrus.DebugLevel)
		break
	}

	//TODO configure logger format
	// log.SetFormatter(&logrus.JSONFormatter{})
	log.SetFormatter(&logrus.TextFormatter{})

	// log to stderr by default
	// log.SetOutput(os.Stderr)
	// log.SetOutput(os.Stdout)

	return log.WithFields(logrus.Fields{
		"module": "peer",
		"source": source,
	})
}
