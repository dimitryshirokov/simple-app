package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dimitryshirokov/simple-app/internal/internal_error"
	"github.com/dimitryshirokov/simple-app/internal/logger"
	"io"
	"net/http"
)

const defaultInternalErrorResponse = "{\"code\":500,\"message\":\"Internal error\"}"
const whoAmI = "simple-app"

type errorResponse struct {
	Error        bool   `json:"error"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	WhoAmI       string `json:"who_am_i"`
	ErrorMessage string `json:"error_message"`
}

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type baseHandler struct{}

func (bh *baseHandler) validationError(w http.ResponseWriter, method string, message string, request interface{}) {
	logger.LogError(message, map[string]interface{}{
		"method":  method,
		"request": request,
	}, nil)
	bh.writeError(w, 400, message, nil)
	return
}

func (bh *baseHandler) writeError(w http.ResponseWriter, code int, message string, e error) {
	er := errorResponse{
		Error:   true,
		Code:    code,
		Message: message,
		WhoAmI:  whoAmI,
	}
	if e != nil {
		er.ErrorMessage = e.Error()
	}
	logger.LogError("Internal error", nil, internal_error.NewErrorFromHandler(message, e, nil))
	errorBody, err := json.Marshal(er)
	if err != nil {
		logger.LogError("error while marshal error response", map[string]interface{}{
			"code":    code,
			"message": message,
		}, err)
		errorBody = []byte(defaultInternalErrorResponse)
	}
	bh.writeBody(w, errorBody)
}

func (bh *baseHandler) writeResult(w http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.LogError("can't create result response body", nil, err)
		bh.writeError(w, 500, "can't create result response body", err)
		return
	}
	bh.writeBody(w, body)
}

func (bh *baseHandler) writeBody(w http.ResponseWriter, body []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := w.Write(body)
	if err != nil {
		logger.LogError("error while writing result response", nil, err)
	}
}

func (bh *baseHandler) getPostRequestBody(r *http.Request) ([]byte, error) {
	if r.Method != "POST" {
		return nil, fmt.Errorf("method %s not allowed", r.Method)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("can't read request body", nil, err)
		return nil, fmt.Errorf("can't read request body: %v", err)
	}
	return body, nil
}
