package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	_ "time/tzdata" // the alpine image has no time zone files; this makes TZ work

	"github.com/rabizd/alumni/internal/middleware"
	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/routes"
)

// main only wires things together: the routes package maps every URL to its
// controller, and the controllers do the rest.
func main() {
	// Errors go to the terminal and to logs/error.log; the file survives a restart.
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Fatalf("logs folder: %v", err)
	}
	errorLog, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("error log: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, errorLog))
	middleware.ErrorLog = errorLog

	if err := model.Connect(context.Background()); err != nil {
		log.Fatalf("database: %v", err)
	}

	mux := http.NewServeMux()
	routes.Web(mux)
	routes.API(mux)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	fmt.Println("listening on http://localhost" + addr) // not an error, so not log
	if err := http.ListenAndServe(addr, middleware.Logger(middleware.ServerErrorGuardian(middleware.NotFoundGuardian(mux)))); err != nil {
		log.Fatal(err)
	}
}
