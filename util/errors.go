package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func New400Response(errors map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Message: "request contains errors",
		Errors:  errors,
	}
}

func New500Response() *ErrorResponse {
	return &ErrorResponse{
		Message: "internal server error",
	}
}

// Calls `RespondTo`, and logs a fatal error if it failed
func (errors *ErrorResponse) RespondToFatal(w http.ResponseWriter) {
	err := errors.RespondTo(w)
	if err != nil {
		log.Fatalf("unable to submit response: %v", err)
	}
}

// Marshals the response into Json, and writes it into the response writer.
// Returns an error if writing or marshalling failed
func (errors *ErrorResponse) RespondTo(w http.ResponseWriter) error {
	responseJson, err := json.Marshal(errors)
	if err != nil {
		return fmt.Errorf("error marshalling response: %v", err)
	}
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(responseJson)
	if err != nil {
		return fmt.Errorf("error writing to socket: %v", err)
	}
	return nil
}
