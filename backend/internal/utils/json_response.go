package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJson(w http.ResponseWriter, r *http.Request, dst any) error {
	const maxBytes = 1_048_576

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decodedBody := json.NewDecoder(r.Body)
	decodedBody.DisallowUnknownFields()

	err := decodedBody.Decode(&dst)
	if err != nil {
		return err
	}

	err = decodedBody.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain one JSON object")
	}

	return nil
}

func SuccessJson(w http.ResponseWriter, status int, message string, data any) error {
	response := JsonResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	// json.NewEncoder(w).Encode(&resp)

	json, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(json)

	return err
}

func ErrorJson(w http.ResponseWriter, status int, message string) {
	response := JsonResponse{
		Success: false,
		Message: message,
	}

	json, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = w.Write([]byte(`{
			"success": false,
			"message": "internal server error"
		}`))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(json)
}
