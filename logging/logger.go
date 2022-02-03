package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	LogNone = iota
	LogInfo
	LogWarning
	LogError
	LogVerbose
	LogDebug
)

type FileLogger struct {
	logger   *log.Logger
	logLevel int
}

func NewFileLogger() *FileLogger {
	return &FileLogger{
		logger:   nil,
		logLevel: LogInfo,
	}
}

func (logger *FileLogger) StartLog() error {
	t := time.Now()
	logfilename := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()) + ".log"
	f, err := os.OpenFile(logfilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger.logger = log.New(f, "", log.Lmicroseconds|log.Llongfile)

	//    logger.log.Print("Log startet")
	return nil
}

func (logger *FileLogger) Print(msg string, level int) error {
	if logger.logger == nil {
		return errors.New("Logger nicht korrekt initialisiert")
	} else {
		fmt.Printf("Adresse des Loggers: %p ; Message: %v", logger, msg)
	}
	logger.logger.Print(msg)
	return nil
}
