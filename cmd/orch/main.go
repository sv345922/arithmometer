package main

import (
	"arithmometer/orchestr/http/server"
	"arithmometer/orchestr/inter/tasker"
	"log"
)

// создать задачу (выражение)
// зафиксировать тайминги операторов
// сохранить выражение в БД
// сделать список задач для вычисления
// сохранить список задач
// отдать задачу вычислителю
// получить ответ от вычислителя
// обновить список задач
// повторить до завершения всех задач
// вернуть ответ клиенту при запросе

func main() {
	// создать пустую базу
	//tasker.CreateEmptyDb()
	// сделать список задач для вычисления
	ws, err := tasker.RunTasker()
	if err != nil {
		log.Print("main: %v", err)
	}
	server.RunServer(ws)
}
