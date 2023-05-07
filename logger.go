package logger

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

var (
	errorLogger    = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags)
	infoLogger     = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	positiveLogger = log.New(os.Stdout, "\u001b[32mINFO: \u001b[0m", log.LstdFlags)
	debugLogger    = log.New(os.Stdout, "\u001b[33mDEBUG: \u001b[0m", log.LstdFlags)
	defaultDepth   = 2
	debugMode      = true
)

type CustomError struct {
	msg string
	err error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// WriteLogsToFile allow to write logs to `logPath`
func WriteLogsToFile(logPath string) {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	errorLogger.SetOutput(logFile)
	infoLogger.SetOutput(logFile)
	positiveLogger.SetOutput(logFile)
	debugLogger.SetOutput(logFile)
}

func (b *CustomError) Error() string {
	j, err := json.MarshalIndent(b, " ", " ")
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func makeStackTrace(depth int, err error) string {
	v, ok := err.(stackTracer)
	if ok {
		trace := v.StackTrace()
		resultLog := "\n"
		for _, frame := range trace[depth:] {
			text, _ := frame.MarshalText()
			resultLog += string(text) + "\n"
		}
		return resultLog
	} else if err != nil {
		return err.Error()
	} else {
		return ""
	}
}

func LogError(err error) {
	if err != nil {
		logError(defaultDepth, err.Error(), makeStackTrace(defaultDepth, err))
	} else {
		logError(defaultDepth, "given error message doesnt implement error interface")
	}
}

func logError(depth int, v ...any) {
	_ = errorLogger.Output(depth+1, fmt.Sprintln(v...))
}

func LogInfo(msg ...any) {
	infoLogger.Println(msg...)
}

func LogPositive(msg ...any) {
	positiveLogger.Println(msg...)
}

func ChangeDebugMode(mode bool) {
	debugMode = mode
}

func LogDebug(msg ...any) {
	if debugMode {
		debugLogger.Println(msg...)
	}
}
