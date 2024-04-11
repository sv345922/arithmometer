package app

import (
	"arithmometer/internal/configs"
	"arithmometer/internal/useCases/getResult"
	"arithmometer/internal/useCases/getTask"
	"arithmometer/internal/useCases/giveAnswer"
	"arithmometer/internal/useCases/newExpression"
	"arithmometer/internal/wSpace"
	"log"
	"net/http"
)

func RunServer(ws *wSpace.WorkingSpace) error {
	mux := http.NewServeMux()

	// Дать ответ клиенту о результатах вычисления выражений
	mux.HandleFunc("/getresult", getResult.GetResult(ws))

	// Получение нового выражения от клиента
	mux.HandleFunc("/newexpression", newExpression.NewExpression(ws))

	// Дать задачу вычислителю
	mux.HandleFunc("/gettask", getTask.GetTask(ws))

	// Получить ответ от вычислителя
	mux.HandleFunc("/giveanswer", giveAnswer.GiveAnswer(ws))

	log.Println("Server is working")
	defer log.Println("Server stopped")
	err := http.ListenAndServe("localhost:"+configs.Port, mux)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
