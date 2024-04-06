package server

import (
	"arithmometer/internal/wSpace"
	"arithmometer/orchestr/http/server/handler"
	"log"
	"net/http"
)

func RunServer(ws *wSpace.WorkingSpace) error {
	mux := http.NewServeMux()

	// Дать ответ клиенту о результатах вычисления выражений
	mux.HandleFunc("/getresult", handler.GetResult(ws))

	// Получение нового выражения от клиента
	mux.HandleFunc("/newexpression", handler.NewExpression(ws))

	// Дать задачу вычислителю
	mux.HandleFunc("/gettask", handler.GetTask(ws))

	// Получить ответ от вычислителя
	mux.HandleFunc("/giveanswer", handler.GiveAnswer(ws))

	log.Println("Server is working")
	defer log.Println("Server stopped")
	err := http.ListenAndServe("localhost:8000", mux)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
