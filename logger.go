package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

var (
	errorLogger    = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags)
	infoLogger     = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	positiveLogger = log.New(os.Stdout, "\u001b[32mINFO: \u001b[0m", log.LstdFlags)
	debugLogger    = log.New(os.Stdout, "\u001b[33mDEBUG: \u001b[0m", log.LstdFlags)
	defaultDepth   = 2
	debugMode      = true
)

const (
	HeaderContentType          = "Content-Type"
	HeaderContentTypeValueJson = "application/json"
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

// SafeHandler panic handler
func SafeHandler(_ http.Handler, serveHTTP func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL.String())
		b, err := io.ReadAll(req.Body)
		// catching errors when reading the body
		if err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(b))
		} else {
			handleInternalError(err, b, req, w)
			return
		}
		// catching errors during request processing
		defer func() {
			if e := recover(); e != nil {
				handleInternalError(e, b, req, w)
			}
		}()
		serveHTTP(w, req)
	})
}

func castError(err any) error {
	var ie error
	switch er := err.(type) {
	case string:
		ie = errors.New(er)
	case *runtime.Error:
	case error:
		ie = er
	}
	return errors.Wrap(ie, "capturing panic error")
}

func handleInternalError(err any, body []byte, r *http.Request, w http.ResponseWriter) {
	if e := castError(err); e != nil {
		LogError(errors.Wrap(e, r.URL.String()))
	} else {
		LogError(errors.Wrap(errors.New("trying to log panic error, but err object is nil"), r.URL.String()))
	}
	if body != nil {
		if len(body) > 0 {
			LogError(errors.New(string(body)))
		}
	}
	w.Header().Add(HeaderContentType, HeaderContentTypeValueJson)
	w.WriteHeader(http.StatusInternalServerError)
	if ee := json.NewEncoder(w).Encode(err); ee != nil {
		log.Println(ee)
	}
}
