package server

import (
	"log"
	"net/http"

	"arithmometer/orchestr/http/server/handler"
)

func RunServer() {
	//http.HandleFunc("/getresult", handler.GetResult)
	log.Println("Starting Server")
	http.HandleFunc("/newexpression", handler.NewExpression)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
