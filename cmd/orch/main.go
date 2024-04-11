package main

import (
	"arithmometer/internal/app"
	"arithmometer/internal/dataBase"
	"arithmometer/internal/wSpace"
	"log"
	"os"
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
	if len(os.Args) > 1 {
		if os.Args[1] == "new" {
			err := dataBase.CreateEmptyDb()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	// сделать список задач для вычисления
	ws, err := RunTasker()
	if err != nil {
		log.Printf("main: %v", err)
	}
	err = app.RunServer(ws)
	log.Println(err)
}

// Создает рабочее пространство из сохраненной базы данных
func RunTasker() (*wSpace.WorkingSpace, error) {
	// Восстанавливаем выражения и задачи из базы данных
	// Загрузка сохраненной БД
	ws, err := wSpace.LoadDB()
	if err != nil {
		log.Println("ошибка загрузки БД", err)
		return ws, err
	}
	return ws, err
}
