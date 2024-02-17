package main

import (
	"arithmometer/orchestr/http/server"
	"arithmometer/orchestr/inter/tasker"
	"context"
	"log"
	"time"
)

// создать задачу (выражение) +
// зафиксировать тайминги операторов +

// сохранить выражение в БД

// сделать список задач для вычисления

// сохранить список задач

// отдать задачу вычислителю

// получить ответ от вычислителя

// обновить список задач (по завершении удалить файл

// повторить до завершения всех задач

// вернуть ответ клиенту при запросе

func main() {

	// сделать список задач для вычисления
	ws, err := tasker.RunTasker()
	if err != nil {
		log.Print("main: %v", err)
	}
	ctx := context.WithValue(context.Background(), "ws", ws)
	time.Sleep(time.Second)
	server.RunServer(ctx)
}
