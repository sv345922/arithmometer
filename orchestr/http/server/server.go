package server

import (
	"arithmometer/orchestr/http/server/handler"
	"arithmometer/orchestr/inter/tasker"
	"log"
	"net/http"
)

func RunServer(ws *tasker.WorkingSpace) error{
	mux := http.NewServeMux()

	// Дать ответ клиенту о результатах вычисления выражений
	mux.HandleFunc("/getresult", handler.GetResult(ws))
	//mux.Handle("/getresult", stateContext(http.HandlerFunc(handler.GetResult)))

	// Получение нового выражения от клиента
	mux.HandleFunc("/newexpression", handler.NewExpression(ws))
	//mux.Handle("/newexpression", stateContext(http.HandlerFunc(handler.NewExpression)))

	// Дать задачу вычислителю
	mux.HandleFunc("/gettask", handler.GetTask(ws))
	//mux.Handle("/gettask", stateContext(http.HandlerFunc(handler.GetTask)))

	// Получить ответ от вычислителя
	mux.HandleFunc("/giveanswer", handler.GiveAnswer(ws))
	//mux.Handle("/giveanswer", stateContext(http.HandlerFunc(handler.GiveAnswer)))

	log.Println("Server is working")
	defer log.Println("Server stopped")
	err := http.ListenAndServe("localhost:8000", mux)
	log.Fatal(err)
	return err
}
