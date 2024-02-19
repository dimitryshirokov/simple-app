package handler

import (
	"encoding/json"
	"github.com/dimitryshirokov/simple-app/internal/service"
	"net/http"
)

func NewSubtractionHandler(calculatorService *service.CalculatorService) *SubtractionHandler {
	return &SubtractionHandler{calculatorService: calculatorService}
}

type SubtractionHandler struct {
	baseHandler
	calculatorService *service.CalculatorService
}

func (h *SubtractionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.writeError(w, 405, "method not allowed", nil)
		return
	}
	body, err := h.getPostRequestBody(r)
	if err != nil {
		h.writeError(w, 500, "can't get POST request body", err)
		return
	}
	dto := &CalcDto{}
	err = json.Unmarshal(body, dto)
	if err != nil {
		h.writeError(w, 500, "can't unmarshal body to calc dto", err)
		return
	}
	result, err := h.calculatorService.Subtraction(dto.A, dto.B)
	if err != nil {
		h.writeError(w, 500, "execution error", err)
		return
	}
	h.writeResult(w, result)
}
