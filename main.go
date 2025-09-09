package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charseft=utf-8")
	fmt.Fprint(w, "<h1>Welcome</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charseft=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p> mail me at <a href=\"mailto:bolatlabakbay@gmail.com\">bolatlabakbay@gmail.com</a>")
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/", pathHandler)
	fmt.Println("Server starting on :3000...")
	http.ListenAndServe(":3000", nil)
}
