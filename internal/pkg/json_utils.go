package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func EncodeJSON(w http.ResponseWriter, data Response, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func DecodeValidJSON[T Validator](r *http.Request) (T, map[string]string, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, nil, err
	}

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("valid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}

func DecodeJSON[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func ValidateAndParseUUID(raw_id string) (uuid.UUID, error) {
	if err := uuid.Validate(raw_id); err != nil {
		return uuid.UUID{}, err
	}

	id := uuid.MustParse(raw_id)
	return id, nil
}
