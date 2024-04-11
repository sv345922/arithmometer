package giveAnswer

import (
	"arithmometer/internal/entities"
	"arithmometer/internal/wSpace"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// var ZeroDiv = errors.New("zero division")

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
		var container entities.MessageResult
		err := json.NewDecoder(r.Body).Decode(&container)
		if err != nil {
			log.Println("ошибка json при обработке ответа вычислителя")
			return
		}
		log.Println("Получен ответ от вычислителя", container.Result)
		// парсим id задачи в виде uint64
		id := container.Id
		// Если деление на ноль
		if errors.Is(container.Err, entities.ZeroDiv) {
			// найти корень выражения и внести ошибку деления на ноль
			root, expression := ws.GetRoot(id)
			expression.Status = "zero division"
			// удалить все узлы дерева выражения
			nodesId := make([]uint64, 0)
			ws.GetExpressionNodesID(root.Id, nodesId)

		}
		// Обновляем очередь задач с учетом выполненной задачи и заносим результат вычисления
		rootFlag, err := ws.Queue.AddAnswer(id, container.Result)
		if err != nil {
			log.Println("ошибка обновления задач:", err)
		}
		// Если вычислен корень выражения
		// обновим выражение и запишем результат
		if rootFlag {
			for _, expr := range ws.Expressions {
				if expr.RootId == id {
					expr.Result = container.Result
					expr.Status = "done"
				}
			}
		}
		w.WriteHeader(http.StatusOK)
		ws.Save()
	}
}
