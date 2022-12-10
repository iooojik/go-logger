package logger

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"
	"runtime"
)

var (
	errorLogger    = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags)
	infoLogger     = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	positiveLogger = log.New(os.Stdout, "\u001b[32mINFO: \u001b[0m", log.LstdFlags)
	debugLogger    = log.New(os.Stdout, "\u001b[33mDEBUG: \u001b[0m", log.LstdFlags)
)

type CustomError struct {
	msg string
	err error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func WriteLogsToFile(write bool, logPath string) {
	if write {
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		errorLogger.SetOutput(logFile)
		infoLogger.SetOutput(logFile)
		positiveLogger.SetOutput(logFile)
		debugLogger.SetOutput(logFile)
	}
}

func (b *CustomError) Error() string {
	j, err := json.MarshalIndent(b, " ", " ")
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func makeStackTrace(trace errors.StackTrace) string {
	resultLog := "\n"
	for _, frame := range trace {
		text, _ := frame.MarshalText()
		resultLog += string(text) + "\n"
	}
	return resultLog
}

func MakeError(msg *string, err error) *CustomError {
	return &CustomError{
		msg: err.Error(),
		err: err,
	}
}

func LogError(err any) {
	var internalError error
	switch er := err.(type) {
	case string:
		internalError = errors.New(er)
	case *CustomError:
		internalError = er.err
	case *runtime.Error:
	case error:
		internalError = er
	}
	errorLogger.Println(internalError.Error())
	errorLogger.Println(makeStackTrace(errors.WithStack(internalError).(stackTracer).StackTrace()))
}

func LogInfo(msg any) {
	infoLogger.Println(msg)
}

func LogPositive(msg any) {
	positiveLogger.Println(msg)
}

func LogDebug(msg any) {
	debugLogger.Println(msg)
}
