package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ValidationError struct {
	Field   string `json:"field" bidirekt:"field"`
	Message string `json:"message" bidirekt:"field"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func pathID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil
}
