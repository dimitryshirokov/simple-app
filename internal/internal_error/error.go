package internal_error

import (
	"errors"
	"fmt"
	"runtime"
)

type InternalError struct {
	message  string
	file     string
	line     int
	previous []error
}

func (ie InternalError) Error() string {
	return ie.message
}

func (ie InternalError) GetTrace() ([]string, []string) {
	messages := make([]string, 0)
	fullTrace := make([]string, 0)
	for _, oe := range ie.previous {
		var oeIe *InternalError
		if errors.As(oe, &oeIe) {
			fullTrace = append(fullTrace, fmt.Sprintf("%s:%d", oeIe.file, oeIe.line))
			messages = append(messages, oeIe.message)
		} else {
			if oe != nil {
				messages = append(messages, oe.Error())
			}
		}
	}
	fullTrace = append(fullTrace, fmt.Sprintf("%s:%d", ie.file, ie.line))
	fullTrace = ie.reverseSlice(fullTrace)
	messages = append(messages, ie.message)
	messages = ie.reverseSlice(messages)
	return fullTrace, messages
}

func (ie InternalError) reverseSlice(data []string) []string {
	j := len(data) - 1
	i := 0
	var reverse func()
	reverse = func() {
		if i < j {
			data[i], data[j] = data[j], data[i]
			i, j = i+1, j-1
			reverse()
		}
	}
	reverse()
	return data
}

func NewErrorLog(message string, err error, addData map[string]interface{}) *InternalError {
	return newError(message, err, 4)
}

func NewErrorFromHandler(message string, err error, addData map[string]interface{}) *InternalError {
	return newError(message, err, 3)
}

func NewError(message string, err error, addData map[string]interface{}) *InternalError {
	return newError(message, err, 2)
}

func newError(message string, err error, skipCaller int) *InternalError {
	_, file, line, _ := runtime.Caller(skipCaller)
	ie := &InternalError{
		message: message,
		file:    file,
		line:    line,
	}
	var previous []error
	var oldIe *InternalError
	if errors.As(err, &oldIe) {
		previous = oldIe.previous
	} else {
		previous = make([]error, 0)
	}
	previous = append(previous, err)
	ie.previous = previous
	return ie
}
