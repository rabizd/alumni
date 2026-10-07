package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode"

	"github.com/rabizd/alumni/internal/view"
)

// GET / -> the landing page
func Root(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "main.html", nil)
}

// GET /main -> the page the temporary redirect points at
func Main(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "main.html", nil)
}

// GET /about -> a temporary about page
func About(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "about.html", nil)
}

// GET /hello -> Hello, World!
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

// GET /hello/{name} -> Hello, Emre!
func HelloName(w http.ResponseWriter, r *http.Request) {
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

// GET /sum/{number1}/{number2} -> the sum of the two numbers
func Sum(w http.ResponseWriter, r *http.Request) {
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

// GET /temporary -> temporary redirect to the main page
func Temporary(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)
}
