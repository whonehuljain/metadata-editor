package logger

import (
	"log"
	"os"
)

type Logger struct {
	verbose bool
	logger  *log.Logger
}

func New(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		logger:  log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string) {
	l.logger.Printf("INFO: %s", msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Printf("ERROR: %s", msg)
}

func (l *Logger) Debug(msg string) {
	if l.verbose {
		l.logger.Printf("DEBUG: %s", msg)
	}
}

func (l *Logger) Fatal(msg string) {
	l.logger.Fatalf("FATAL: %s", msg)
}
