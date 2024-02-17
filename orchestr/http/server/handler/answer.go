package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Обработчик, принимает от вычислителя ответ
func GiveAnswer(w http.ResponseWriter, r *http.Request) {
	// Получаем рабочее пространство из контекста
	ws, ok := tasker.GetWs(r.Context())
	if !ok {
		log.Println("ошибка контекста")
		return
	}
	// Получаем узлы
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
	//id, _ := strconv.Atoi(container.Id)
	id, _ := strconv.ParseUint(container.Id, 10, 64)
	// Обновляем очередь задач с учетом выполненной задачи
	err = ws.UpdateTasks(id, &container.AnswerN)

	// Получаем узел из списка узлов
	task := tasker.FindNodes(id, nodes)
	if task != nil {
		w.WriteHeader(http.StatusOK)
	}
	// Устанавливаем значения результата вычисления
	task.Calculated = true
	// Если ошибка вычислителя, это деление на ноль
	if container.AnswerN.Err != nil {
		log.Println("Ошибка выражения, ", err)
		task.ErrZeroDiv = fmt.Errorf("ошибка вычисления: %v", container.AnswerN.Err)
	}
	// Установить вычисленное значение
	/*
		task.Calculated = true
		task.Val = container.AnswerN.Result

		// обновление очереди задач
		err, ok := ws.UpdateTasks(id)
		if err != nil {
			log.Println("ошибка обновления задач по результатам вычисления", err)
		}
	*/
	if !ok {
		w.Write([]byte("обновление задач не удалось"))
		w.WriteHeader(http.StatusNoContent)
	}

}
