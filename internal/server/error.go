package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (s *Server) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s err: %s", r.Method, r.URL.Path, err)
	s.Logger.Errorw("internal server error",
		"path", r.URL.Path,
		"method", r.Method,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusInternalServerError, "there was a problem")
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found: %s path: %s err: %s", r.Method, r.URL.Path, err)
	s.Logger.Warnw("not found",
		"path", r.URL.Path,
		"method", r.Method,
		"error", err.Error(),
	)

	writeJSONError(w, http.StatusNotFound, "the requested resource could not be found")
}

func (s *Server) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request: %s path: %s err: %s", r.Method, r.URL.Path, err)
	s.Logger.Warnw("bad request",
		"path", r.URL.Path,
		"method", r.Method,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (s *Server) Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("unauthorized: %s path: %s", r.Method, r.URL.Path)
	s.Logger.Warnw("unauthorized",
		"path", r.URL.Path,
		"method", r.Method,
		"err", err.Error(),
	)

	writeJSONError(w, http.StatusUnauthorized, "you are not authorized to access this resource")
}

func (s *Server) Forbidden(w http.ResponseWriter, r *http.Request) {
	log.Printf("forbidden: %s path: %s", r.Method, r.URL.Path)
	s.Logger.Warnw("forbidden",
		"path", r.URL.Path,
		"method", r.Method,
	)
	writeJSONError(w, http.StatusForbidden, "you do not have permission to access this resource")
}
func (s *Server) Conflict(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict: %s path: %s err: %s", r.Method, r.URL.Path, err)
	s.Logger.Warnw("conflict",
		"path", r.URL.Path,
		"method", r.Method,
		"error", err.Error(),
	)
	writeJSONError(w, http.StatusConflict, err.Error())
}

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "Server.json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxByes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, map[string]string{"error": message})
}

func (s *Server) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &envelope{Data: data})
}
