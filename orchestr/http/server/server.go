package server

import (
	"log"
	"net/http"
)

func newExpression(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	expression, err := Parse(r.Form.Get("expr"))
}

func runServer() {
	http.HandleFunc("/expression", newExpression)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
