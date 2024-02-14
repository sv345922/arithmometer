package server

import (
	"context"
	"log"
	"net/http"

	"arithmometer/orchestr/http/server/handler"
)

func RunServer(ctx context.Context) {

	http.HandleFunc("/getresult", handler.GetResult)
	http.HandleFunc("/newexpression", handler.NewExpression)
	http.HandleFunc("/gettask", handler.GetTask)
	http.HandleFunc("/getanswer", handler.GetAnswer)
	//http.ServerContextKey
	log.Println("Starting Server")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
