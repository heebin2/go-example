package helper

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Sets the location and level of log globally.
func InitLogger(appName, level string, isDev bool) {

	if isDev {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
			DisableColors: false,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// with both terminal and log file
	multiWriter := io.MultiWriter(&lumberjack.Logger{
		Filename:   "/var/log/tms/" + appName + ".log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}, os.Stdout)
	log.SetOutput(multiWriter)

	// callstack
	log.SetReportCaller(true)

	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.InfoLevel)
		log.Warn("log option is invalid. Check the -log option. The log level is set to ", log.InfoLevel.String(), ". See --help")
		return
	}

	log.SetLevel(logLevel)
	log.Debug("log level : ", logLevel.String())
}

func SetLogLevel(level string) {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(log.InfoLevel)
		return
	}

	log.SetLevel(logLevel)
	log.Info("log level : ", logLevel.String())
}
