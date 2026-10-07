package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rabizd/alumni/internal/controller"
)

// main only wires things together: every route points at a controller, and
// the controllers do the rest.
func main() {
	mux := http.NewServeMux()

	// 1. GET / -> the landing page
	mux.HandleFunc("GET /{$}", controller.Root)

	// 2. GET /hello -> Hello, World!
	mux.HandleFunc("GET /hello", controller.Hello)

	// 3. GET /hello/{name} -> Hello, Emre!
	mux.HandleFunc("GET /hello/{name}", controller.HelloName)

	// 4. GET /sum/{number1}/{number2} -> the sum of the two numbers
	mux.HandleFunc("GET /sum/{number1}/{number2}", controller.Sum)

	// 5. GET /temporary -> temporary redirect to the main page
	mux.HandleFunc("GET /temporary", controller.Temporary)

	// The main page the temporary redirect points at.
	mux.HandleFunc("GET /main", controller.Main)

	// 6. GET /about -> a temporary about page
	mux.HandleFunc("GET /about", controller.About)

	// GET /api/health -> a JSON health check
	mux.HandleFunc("GET /api/health", controller.Health)

	// GET /api/swagger -> the API documentation, generated from openapi.json
	mux.HandleFunc("GET /api/swagger", controller.Swagger)
	mux.HandleFunc("GET /api/swagger.json", controller.SwaggerSpec)

	// The users resource, kept in memory for now, served two ways: JSON for
	// programs under /api/users, HTML pages for people under /users.
	api := controller.ApiUserController{}
	mux.HandleFunc("GET /api/users", api.Index)
	mux.HandleFunc("POST /api/users", api.Store)
	mux.HandleFunc("GET /api/users/{id}", api.Show)
	mux.HandleFunc("PUT /api/users/{id}", api.Update)
	mux.HandleFunc("PATCH /api/users/{id}", api.Patch)
	mux.HandleFunc("DELETE /api/users/{id}", api.Destroy)

	web := controller.UserController{}
	mux.HandleFunc("GET /users", web.Index)
	mux.HandleFunc("GET /users/new", web.Create)
	mux.HandleFunc("POST /users", web.Store)
	mux.HandleFunc("GET /users/{id}", web.Show)
	mux.HandleFunc("GET /users/{id}/edit", web.Edit)
	mux.HandleFunc("POST /users/{id}", web.Update)
	mux.HandleFunc("POST /users/{id}/delete", web.Destroy)

	// The alumni resource itself.
	mux.HandleFunc("GET /alumni", controller.ListAlumni)
	mux.HandleFunc("POST /alumni", controller.CreateAlumni)

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
