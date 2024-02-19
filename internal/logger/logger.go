package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dimitryshirokov/simple-app/internal/internal_error"
	"log"
	"os"
	"time"
)

var debug = os.Getenv("LOG_DEBUG")
var logTimeFormat = os.Getenv("LOG_TIME_FORMAT")

type jsonTime struct {
	time.Time
}

func (t jsonTime) MarshalJSON() ([]byte, error) {
	tf := logTimeFormat
	if logTimeFormat == "" {
		tf = "2006-01-02 15:04:05 -07:00"
	}
	stamp := fmt.Sprintf("\"%s\"", t.Format(tf))
	return []byte(stamp), nil
}

type logMsg struct {
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	BaseMessage string                 `json:"base_message"`
	Data        map[string]interface{} `json:"data"`
	Trace       []string               `json:"trace"`
	Messages    []string               `json:"messages"`
	Time        jsonTime               `json:"time"`
}

func LogDebug(message string, data map[string]interface{}) {
	if debug == "1" {
		doLog("DEBUG", message, data, nil)
	}
}

func LogInfo(message string, data map[string]interface{}) {
	doLog("INFO", message, data, nil)
}

func LogWarning(message string, data map[string]interface{}, err error) {
	doLog("WARNING", message, data, err)
}

func LogError(message string, data map[string]interface{}, err error) {
	doLog("ERROR", message, data, err)
}

func doLog(level string, message string, data map[string]interface{}, err error) {
	var ie *internal_error.InternalError
	if errors.As(err, &ie) {
		newErr := internal_error.NewErrorLog(message, ie, data)
		logInternalError(level, newErr, data)
		return
	}
	if err != nil {
		newErr := internal_error.NewErrorLog(message, err, data)
		logInternalError(level, newErr, nil)
		return
	}
	msg := logMsg{
		Level:       level,
		Message:     message,
		BaseMessage: message,
		Data:        data,
		Trace:       make([]string, 0),
		Messages:    make([]string, 0),
		Time:        jsonTime{time.Now()},
	}
	b, e := json.Marshal(msg)
	if e != nil {
		log.Println(e)
	}
	_, _ = fmt.Println(string(b))
}

func logInternalError(level string, err *internal_error.InternalError, data map[string]interface{}) {
	trace, messages := err.GetTrace()
	msg := logMsg{
		Level:       level,
		Message:     err.Error(),
		BaseMessage: messages[len(messages)-1],
		Data:        data,
		Trace:       trace,
		Messages:    messages,
		Time:        jsonTime{time.Now()},
	}
	b, e := json.Marshal(msg)
	if e != nil {
		log.Println(e)
	}
	_, _ = fmt.Println(string(b))
}
