package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var globalLogger *FileLogger

type LogLevel int

const (
	LogInfo    LogLevel = 1 << iota
	LogWarning LogLevel = 2
	LogError   LogLevel = 4
	LogDebug   LogLevel = 8
)

type FileLogger struct {
	logger   *log.Logger
	logLevel LogLevel
}

func NewFileLogger() *FileLogger {
	return &FileLogger{
		logger: nil,
	}
}

func InitLog(logfilename string, loglevel string) *FileLogger {

	newLogger := FileLogger{logLevel: ParseLogLevel(loglevel)}

	if len(logfilename) == 0 {
		logfilename = fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()) + ".log"
	}
	f, err := os.OpenFile(logfilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	newLogger.logger = log.New(f, "", log.Ldate|log.Ltime)
	globalLogger = &newLogger
	return &newLogger
}

func (level LogLevel) SetLogLevel(newLevel LogLevel) LogLevel {
	return level | newLevel
}

func (level LogLevel) UnsetLogLevel(removeLevel LogLevel) LogLevel {
	return level &^ removeLevel
}

func (level LogLevel) HasLogLevel(compareLevel LogLevel) bool {
	return level&compareLevel != 0
}

func (level LogLevel) ToString() string {
	var levelString string
	if level.HasLogLevel(LogInfo) {
		levelString = levelString + "Info"
	}
	if level.HasLogLevel(LogWarning) {
		levelString = levelString + "Warning"
	}
	if level.HasLogLevel(LogError) {
		levelString = levelString + "Error"
	}
	if level.HasLogLevel(LogDebug) {
		levelString = levelString + "Debug"
	}
	return levelString
}

// Function takes level as string (see enum LogLevel) and returns the LogLevel
func LogLevelToInt(level string) LogLevel {
	var return_level LogLevel
	switch strings.ToLower(level) {
	case "info":
		return_level = LogInfo
	case "warning":
		return_level = LogWarning
	case "error":
		return_level = LogError
	case "debug":
		return_level = LogDebug
	}
	return return_level
}

// Func takes in list of comma-separated Log-Levels, e.g. "error,warning,info"
func ParseLogLevel(input string) LogLevel {
	var level LogLevel
	levels := strings.Split(input, ",")
	for _, l := range levels {
		level.SetLogLevel(LogLevelToInt(l))
	}
	return level
}

func (logger *FileLogger) Info(msg string) {
	logger.logger.Print(msg)
}

func (logger *FileLogger) Warning(msg string) {
	logger.logger.Print(msg)
}

func (logger *FileLogger) Error(msg string) {
	logger.logger.Print(msg)
}

func (logger *FileLogger) Debug(msg string) {
	logger.logger.Print(msg)
	//	logger.logger.Print(msg, LogDebug)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}

func Println(v ...interface{}) {
	log.Println(v)
}

func Fatal(v ...interface{}) {
	log.Fatal(v)
}
