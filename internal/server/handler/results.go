package handler

import (
	"encoding/json"
	"github.com/dimitryshirokov/simple-app/internal/service"
	"net/http"
)

type resultsDto struct {
	Type   string `json:"type"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ResultsHandler struct {
	baseHandler
	calculatorService *service.CalculatorService
}

func NewResultsHandler(calculatorService *service.CalculatorService) *ResultsHandler {
	return &ResultsHandler{calculatorService: calculatorService}
}

func (h *ResultsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.writeError(w, 405, "method not allowed", nil)
		return
	}
	body, err := h.getPostRequestBody(r)
	if err != nil {
		h.writeError(w, 500, "can't get POST request body", err)
		return
	}
	dto := &resultsDto{}
	err = json.Unmarshal(body, dto)
	if err != nil {
		h.writeError(w, 500, "can't unmarshal body to calc dto", err)
		return
	}
	result, count, err := h.calculatorService.Results(dto.Type, dto.Limit, dto.Offset)
	if err != nil {
		h.writeError(w, 500, "can't list results", err)
		return
	}
	h.writeResult(w, map[string]interface{}{
		"data":  result,
		"count": count,
	})
}
