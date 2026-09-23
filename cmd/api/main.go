package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"unicode"
)

// The HTML lives in its own files; go:embed bakes them into the binary at
// build time, so the server stays a single executable with no loose files.
//
//go:embed templates/*.html
var templateFiles embed.FS

var templates = template.Must(template.ParseFS(templateFiles, "templates/*.html"))

func main() {
	mux := http.NewServeMux()

	// 1. GET / -> OK
	mux.HandleFunc("GET /{$}", handleRoot)

	// 2. GET /hello -> Hello, World!
	mux.HandleFunc("GET /hello", handleHello)

	// 3. GET /hello/{name} -> Hello, Emre!
	mux.HandleFunc("GET /hello/{name}", handleHelloName)

	// 4. GET /sum/{number1}/{number2} -> the sum of the two numbers
	mux.HandleFunc("GET /sum/{number1}/{number2}", handleSum)

	// 5. GET /temporary -> temporary redirect to the main page
	mux.HandleFunc("GET /temporary", handleTemporary)

	// The main page the temporary redirect points at.
	mux.HandleFunc("GET /main", handleMain)

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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func handleHelloName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	fmt.Fprintf(w, "Hello, %s!", capitalize(name))
}

// capitalize upper-cases the first letter so /hello/emre answers "Hello, Emre!".
func capitalize(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func handleSum(w http.ResponseWriter, r *http.Request) {
	a, err := strconv.Atoi(r.PathValue("number1"))
	if err != nil {
		http.Error(w, "number1 must be a whole number", http.StatusBadRequest)
		return
	}

	b, err := strconv.Atoi(r.PathValue("number2"))
	if err != nil {
		http.Error(w, "number2 must be a whole number", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%d", a+b)
}

func handleTemporary(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "main.html", nil); err != nil {
		log.Printf("rendering main.html: %v", err)
	}
}

