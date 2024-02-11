package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	// Получаем список задач
	tasks, err := tasker.GetTasks(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	// Получаем мапу узлов
	nodes, err := tasker.GetNodes(r.Context())
	if err != nil {
		log.Println(err)
		return
	}
	// Проверить что это запрос POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("требуется метод POST"))
		return
	}
	// Читаем тело запроса, в котором записан ответ
	var container tasker.AnswerContainer
	err = json.NewDecoder(r.Body).Decode(&container)
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
	go tasks.Update(nodes) // обновление списка задач
}
