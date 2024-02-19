package http_server

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"regexp"
)

const (
	slash        = string(filepath.Separator)
	space        = " "
	openBracket  = "{"
	closeBracket = "}"
)

var bracketRegex = regexp.MustCompile(`\{(.*?)\}`)

type (
	wildcardParamsKey struct{}

	// UserRole is the representation of a user role enum
	UserRole string
)

// Response ...
type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
	Data    any      `json:"data"`
}

// ErrorResponse ...
func ErrorResponse(w http.ResponseWriter, code int, err error) {
	resp := &Response{
		Code:    code,
		Message: err.Error(),
		Details: []string{},
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization")
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	})
}
