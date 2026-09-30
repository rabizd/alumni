package main

import (
	_ "embed"
	"net/http"
)

// The OpenAPI document is baked into the binary, like the HTML pages, so the
// running server always serves the spec that was built alongside it.
//
//go:embed openapi.json
var openAPISpec []byte

// GET /api/swagger -> the Swagger UI page, which reads the spec below
func handleSwagger(w http.ResponseWriter, r *http.Request) {
	render(w, "swagger.html")
}

// GET /api/swagger.json -> the OpenAPI document itself
func handleSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(openAPISpec)
}
