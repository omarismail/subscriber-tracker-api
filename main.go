package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jonfriesen/subscriber-tracker-api/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		if err := os.Setenv("DATABASE_URL", "postgresql://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"); err != nil {
			panic("failed to set db connection string")
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	handler := api.New()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: handler,
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Unable to shut down server: %v", err)
		} else {
			log.Println("Server stopped")
		}
		wg.Done()
	}()

	go func() {
		log.Printf("Magic is happening on port %s", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("%v", err)
			wg.Done()
		} else {
			log.Println("Server closed")
			wg.Done()
		}
	}()

	wg.Wait()
}
