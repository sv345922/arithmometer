package server

import (
	"arithmometer/orchestr/http/server/handler"
	"context"
	"log"
	"net/http"
)

func RunServer(ctx context.Context) {
	mux := http.NewServeMux()

	// Дать ответ клиенту о результатах вычисления выражений
	//mux.HandleFunc("/getresult", handler.GetResult)
	mux.Handle("/getresult", stateContext(http.HandlerFunc(handler.GetResult)))

	// Получение нового выражения от клиента
	//mux.HandleFunc("/newexpression", handler.NewExpression)
	mux.Handle("/newexpression", stateContext(http.HandlerFunc(handler.NewExpression)))

	// Дать задачу вычислителю
	//mux.HandleFunc("/gettask", handler.GetTask)
	mux.Handle("/gettask", stateContext(http.HandlerFunc(handler.GetTask)))

	// Получить ответ от вычислителя
	//mux.HandleFunc("/giveanswer", handler.GiveAnswer)
	mux.Handle("/giveanswer", stateContext(http.HandlerFunc(handler.GiveAnswer)))

	log.Println("Starting Server")
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
