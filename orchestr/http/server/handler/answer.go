package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Обработчик, принимает от вычислителя ответ
func GiveAnswer(ws *tasker.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	// Получаем рабочее пространство из контекста
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это метод POST
		if r.Method != http.MethodPost {
			log.Println("метод не POST")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("требуется метод POST"))
			return
		}

		// Читаем тело запроса, в котором записан ответ
		defer r.Body.Close()
		var container tasker.AnswerContainer
		err := json.NewDecoder(r.Body).Decode(&container)
		if err != nil {
			log.Println("ошибка json при обработке ответа вычислителя")
			return
		}
		log.Println("Получен ответ от вычислителя", container.AnswerN.Result)
		// парсим id задачи до uint64
		id, _ := strconv.ParseUint(container.Id, 10, 64)
		// Обновляем очередь задач с учетом выполненной задачи и заносим результат вычисления
		err = ws.UpdateTasks(id, &container.AnswerN)
		w.WriteHeader(http.StatusOK)
	}
}
