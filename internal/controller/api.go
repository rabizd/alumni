package controller

import (
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// GET /api/health -> {"status":"ok"}
//
// Reaching this handler at all is the real check: it means the process is up
// and the router is wired correctly.
func Health(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, model.Health{Status: "ok"})
}

// GET /api/swagger -> the Swagger UI page, which reads the spec below
func Swagger(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "swagger.html", nil)
}

// GET /api/swagger.json -> the OpenAPI document itself
func SwaggerSpec(w http.ResponseWriter, r *http.Request) {
	view.OpenAPI(w)
}
