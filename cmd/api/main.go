package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"unicode"
)

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
	fmt.Fprint(w, mainPage)
}

const mainPage = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Alumni</title>
<style>
  body { font-family: system-ui, sans-serif; max-width: 40rem; margin: 4rem auto; padding: 0 1rem; line-height: 1.6; }
  h1 { margin-bottom: 0.25rem; }
  p.sub { color: #666; margin-top: 0; }
  li { margin: 0.35rem 0; }
  code { background: #f2f2f2; padding: 0.1rem 0.3rem; border-radius: 3px; }
</style>
</head>
<body>
  <h1>Alumni</h1>
  <p class="sub">Istanbul University alumni network — development server</p>
  <h2>Routes</h2>
  <ul>
    <li><a href="/">/</a> — <code>OK</code></li>
    <li><a href="/hello">/hello</a> — <code>Hello, World!</code></li>
    <li><a href="/hello/emre">/hello/{name}</a> — greets the name in the path</li>
    <li><a href="/sum/7/35">/sum/{number1}/{number2}</a> — adds the two numbers</li>
    <li><a href="/temporary">/temporary</a> — temporary redirect back to this page</li>
  </ul>
</body>
</html>
`
