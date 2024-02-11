package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"fmt"
	"log"
	"net/http"
)

// Обрабатывает запросы клиента о проверке результата вычисллений
func GetResult(w http.ResponseWriter, r *http.Request) {
	// Получаем список выражений из контекста
	expressions, err := tasker.GetExpressions(r.Context())
	if err != nil {
		log.Println(err)
		return
	}

	// Проверить метод
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("требуется метод Get"))
		return
	}
	// Читаем id параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("не найден id в запросе")
		return
	}
	log.Println("Id =", id)

	// Обновление списка выражений
	expressions.Update()
	// Поиск выражения в списке выражений
	expression := tasker.FindExpression(id, expressions)

	// выражение не найдено
	if expression == nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("id не найден"))
		return
	}
	if expression.Calculated() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expression.Result))
		return
	}
	w.Write([]byte(fmt.Sprintf("выражение %s еще не посчитано", expression.UserTask)))

}
