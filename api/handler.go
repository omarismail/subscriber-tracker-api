package api

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jonfriesen/subscriber-tracker-api/storage/postgresql"
)

var (
	ErrNotFound               = errors.New("Error: not found")
	ErrRquestTypeNotSupported = errors.New("Error: HTTP Request type not supported")
)

type Handler struct {
	Database *postgresql.PostgreSQL
}

// New creates a new http handler
func New(db *postgresql.PostgreSQL) http.Handler {
	mux := http.NewServeMux()

	h := Handler{
		Database: db,
	}

	mux.Handle("/list/", wrapper(h.list))

	return mux
}

func (h *Handler) SetDatabase(db *postgresql.PostgreSQL) {
	h.Database = db
}

func (h *Handler) list(w io.Writer, r *http.Request) (interface{}, int, error) {
	switch r.Method {
	case "GET":
		rKey := strings.TrimPrefix(r.URL.Path, "/v1/get/")
		log.Printf("Looking up %v", rKey)

		v, err := h.Database.ListSubscribers()
		if err != nil {
			return nil, http.StatusNotFound, ErrNotFound
		}

		return v, http.StatusOK, nil
	}

	return nil, http.StatusBadRequest, ErrRquestTypeNotSupported
}

func wrapper(f func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", " application/json")
		w.WriteHeader(status)

		_, err = io.WriteString(w, data.(string))
		if err != nil {
			panic(err)
		}
	}
}
