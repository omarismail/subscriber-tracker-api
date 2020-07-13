package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jonfriesen/subscriber-tracker-api/api"
	"github.com/jonfriesen/subscriber-tracker-api/storage/postgresql"
)

const devPGConnStr = "postgresql://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	handler := api.New(nil)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: handler,
	}

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
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)

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

	go func() {
		var db *postgresql.PostgreSQL
		err := errors.New("database is not ready yet")
		for err != nil {
			log.Println("Checking if database is ready yet.")
			db, err = postgresql.NewConnection(devPGConnStr)
			if err != nil {
				log.Println("error")
				time.Sleep(5 * time.Second)
				err = err
			} else {
				log.Println("connected to db")
				handler.SetDatabase(db)
				defer db.Close()
			}
		}
	}()

	wg.Wait()
}
