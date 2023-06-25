package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewResponse(success bool, message string, err error) Response {
	response := Response{}
	response.Success = success
	response.Message = message
	if err != nil {
		response.Error = err.Error()
	}
	return response
}

func Write(w http.ResponseWriter, status int, resp interface{}) {
	jsonResponse(w, status, resp)
}
func jsonResponse(w http.ResponseWriter, status int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func SanitizeRequest(r *http.Request, req interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(req); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}
	return nil
}
