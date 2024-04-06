package handler

import (
	"arithmometer/internal/wSpace"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Обработчик, принимает от вычислителя ответ
func GiveAnswer(ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
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
		var container wSpace.AnswerContainer
		err := json.NewDecoder(r.Body).Decode(&container)
		if err != nil {
			log.Println("ошибка json при обработке ответа вычислителя")
			return
		}
		log.Println("Получен ответ от вычислителя", container.AnswerN.Result)
		// парсим id задачи в виде uint64
		id, _ := strconv.ParseUint(container.Id, 10, 64)
		// Обновляем очередь задач с учетом выполненной задачи и заносим результат вычисления
		err = ws.UpdateTasks(id, &container.AnswerN)
		if err != nil {
			log.Println("ошибка обновления задач:", err)
		}
		w.WriteHeader(http.StatusOK)
	}
}
