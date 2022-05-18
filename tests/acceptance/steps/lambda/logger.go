package lambda

import (
	"log"
	"os"
	"path"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/sirupsen/logrus"
)

var logger *Logger

// GetLogger returns the logger for lambda operations.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	if logger != nil {
		return logger
	}
	dir := golium.GetConfig().Log.Directory
	logsPath := path.Join(dir, "lambda.log")
	logger, err := NewLogger(logsPath)
	if err != nil {
		logrus.Fatalf("Error creating lambda logger with file: '%s'. %s", logsPath, err)
	}
	return logger
}

// Logger logs the lambda operations in a configurable file.
type Logger struct {
	log *log.Logger
}

// NewLogger creates an instance of the logger.
// It configures the file path where the lambda operations are logged.
func NewLogger(logsPath string) (*Logger, error) {
	file, err := os.OpenFile(logsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	os.Chmod(file.Name(), 0766)
	return &Logger{
		log: log.New(file, "", log.Ldate|log.Lmicroseconds|log.LUTC),
	}, nil
}

// Log a lambda message
func (l Logger) LogMessage(message string) {
	l.log.Printf("%s", message)
}
