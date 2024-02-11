package main

import (
	"arithmometer/orchestr/http/server"
	"arithmometer/orchestr/inter/tasker"
	"context"
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
	tasks, expressions, timings, nodes := tasker.RunTasker()
	ctx := context.WithValue(context.Background(), "tasks", tasks)
	ctx = context.WithValue(ctx, "expressions", expressions)
	ctx = context.WithValue(ctx, "timings", timings)
	ctx = context.WithValue(ctx, "nodes", nodes)
	server.RunServer(ctx)
}
