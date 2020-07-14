package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jonfriesen/subscriber-tracker-api/model"
	"github.com/jonfriesen/subscriber-tracker-api/storage/postgresql"
	"github.com/rs/cors"
)

var (
	ErrNotFound               = errors.New("Error: not found")
	ErrRquestTypeNotSupported = errors.New("Error: HTTP Request type not supported")
)

type MissingDB struct {
	ErrorMessage string `json:"error_message,omitempty"`
}

type handler struct {
	Database *postgresql.PostgreSQL
}

// New creates a new http handler
func New() http.Handler {
	mux := http.NewServeMux()

	h := handler{}

	mux.Handle("/subscribers/", wrapper(h.list))

	corsHandler := cors.Default().Handler(mux)
	return corsHandler
}

func (h *handler) list(w io.Writer, r *http.Request) (interface{}, int, error) {
	switch r.Method {
	case "GET":

		if h.checkDBConnection() != nil {
			return MissingDB{ErrorMessage: "Database appears to be missing. These changes will not be saved."}, http.StatusOK, nil
		}

		v, err := h.Database.ListSubscribers()
		if err != nil {
			return nil, http.StatusNotFound, ErrNotFound
		}

		return v, http.StatusOK, nil
	case "POST":

		if h.checkDBConnection() != nil {
			return MissingDB{ErrorMessage: "Database appears to be missing. These changes will not be saved."}, http.StatusOK, nil
		}

		var sub *model.Subscriber
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			return nil, http.StatusBadRequest, err
		}

		fmt.Printf("%+v", sub)

		v, err := h.Database.AddSubscriber(sub)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return v, http.StatusCreated, nil
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

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *handler) checkDBConnection() error {
	if h.Database != nil && h.Database.IsUp() == nil {
		return nil
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("no database URL set")
	}

	adb, err := postgresql.NewConnection(dbURL)
	if err != nil {
		return err
	}

	h.Database = adb

	return nil
}
