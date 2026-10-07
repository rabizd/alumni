package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rabizd/alumni/internal/routes"
)

// main only wires things together: the routes package maps every URL to its
// controller, and the controllers do the rest.
func main() {
	mux := http.NewServeMux()
	routes.Web(mux)
	routes.API(mux)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
