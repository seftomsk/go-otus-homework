package logger

import (
	"io"
	"log"
	"strings"
	"time"
)

type Logger struct {
	logLevel string
	writer   io.Writer
}

const (
	ERROR = "ERROR"
	WARN  = "WARN"
	INFO  = "INFO"
	DEBUG = "DEBUG"
)

func New(logLevel string, writer io.Writer) *Logger {
	return &Logger{
		logLevel: logLevel,
		writer:   writer,
	}
}

func constructMsg(level string, msg string) string {
	return strings.ToUpper(level) + ": [" + time.Now().UTC().Format("2006-01-02 03:04:05 PM") + "] " + msg + "\n"
}

func (l *Logger) Info(msg string) {
	_, err := l.writer.Write([]byte(constructMsg(INFO, msg)))
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) Error(msg string) {
	_, err := l.writer.Write([]byte(constructMsg(ERROR, msg)))
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) Warn(msg string) {
	_, err := l.writer.Write([]byte(constructMsg(WARN, msg)))
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) Debug(msg string) {
	_, err := l.writer.Write([]byte(constructMsg(DEBUG, msg)))
	if err != nil {
		log.Println(err)
	}
}
