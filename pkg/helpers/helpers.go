package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func PrintError(w http.ResponseWriter, err error) {
	slog.Error("internal error", "error", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func PrintJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		PrintError(w, err)
	}
}
