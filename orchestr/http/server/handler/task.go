package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Даёт задачу калькулятору
func GetTask(ws *tasker.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить метод
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод Get"))
			return
		}
		// Читаем id вычислителя из параметров запроса
		id := r.URL.Query().Get("id")
		if id == "" {
			log.Println("не найден id в запросе вычислителя")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("Запрос задачи от вычислителя с id =", id)
		calcId, err := strconv.Atoi(id)
		if err != nil {
			log.Println(id, "id вычислителя не число", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Обновляем очередь задач, чтобы убрать дедлайны у просроченных
		ws.Tasks.Queue.UpdateWithTimeouts()
		// Получаем задачу из очереди
		task := ws.Tasks.GetTask(calcId)
		if task == nil {
			// Если активных задач нет
			log.Println("нет задач для вычислителя")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// записываем дедлайн для узла с учетом времени выполнения операции
		// достаем оператор из задачи
		// словарь таймигнгов операций
		d := map[string]int{
			"+": task.TimingsN.Plus,
			"-": task.TimingsN.Minus,
			"*": task.TimingsN.Mult,
			"/": task.TimingsN.Div,
		}
		// дедлайн равен времени выполнения операции + 50%, статус задач обновляется при
		// обновлении очереди
		timeout := d[task.TaskN.Op] * 150 / 100
		task.Deadline = time.Now().Add(time.Second * time.Duration(timeout))
		// Сохраняем БД
		ws.Save()
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
		data, _ := json.Marshal(&container) //ошибку пропускаем
		// и записываем в ответ вычислителю
		w.Write(data)
		log.Println("вычислителю дана задача", container.TaskN)
	}
}
