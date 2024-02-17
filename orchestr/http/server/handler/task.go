package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Даёт задачу калькулятору
func GetTask(w http.ResponseWriter, r *http.Request) {
	// Получаем рабочее пространство из контекста
	ws, ok := tasker.GetWs(r.Context())
	if !ok {
		log.Println("ошибка контекста")
		return
	}
	// Проверить метод
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("требуется метод Get"))
		return
	}
	// Читаем id из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("не найден id в запросе вычислителя")
		return
	}
	log.Println("Calculator IdExpression =", id)
	calcId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("ошибка конвертации id  в число", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Получаем задачу из очереди
	task := ws.Tasks.GetTask(calcId)
	if task == nil {
		// Если активных задач нет
		w.WriteHeader(http.StatusNoContent)
	}

	// структура для передачи вычислителю
	type TaskForCalc struct {
		Id       string         `json:"id"`
		TaskN    tasker.Task    `json:"taskN"`
		TimingsN tasker.Timings `json:"timingsN"`
	}
	// Создаем структуру для передачи вычислителю
	container := TaskForCalc{
		Id:       fmt.Sprintf("%d", task.IdTask),
		TaskN:    task.TaskN,
		TimingsN: task.TimingsN,
	}
	// Маршалим её
	data, _ := json.Marshal(container) //ошибку пропускаем
	// и записываем в ответ
	w.Write(data)
}
