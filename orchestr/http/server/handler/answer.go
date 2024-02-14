package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Обработчик, отдает клиенту ответ
func GetAnswer(w http.ResponseWriter, r *http.Request) {
	// Получаем рабочее пространство из контекста
	ws, ok := tasker.GetWs(r.Context())
	if !ok {
		log.Println("ошибка контекста")
		return
	}
	// Получаем мапу узлов
	nodes := ws.AllNodes

	// Проверить что это запрос POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("требуется метод POST"))
		return
	}
	// Читаем тело запроса, в котором записан ответ
	var container tasker.AnswerContainer
	err := json.NewDecoder(r.Body).Decode(&container)
	if err != nil {
		log.Println("ошибка json в ответе")
		return
	}
	// Получаем узел из списка узлов
	task := tasker.FindNodes(container.Id, *nodes)
	if task != nil {
		w.WriteHeader(http.StatusOK)
	}
	// Устанавливаем значения результата вычисления
	task.Calculated = true
	// Если ошибка вычислителя это деление на ноль
	if container.AnswerN.Err != nil {
		log.Println("Ошибка выражения, ", err)
		task.Err = fmt.Errorf("ошибка вычисления: %v", container.AnswerN.Err)
	}
	// Установить вычисленное значение
	task.Calculated = true
	task.Val = container.AnswerN.Result
	go ws.Update() // обновление списка задач
}
