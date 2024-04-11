package entities

// Expression Выражение
type Expression struct {
	Id        uint64  `json:"id"`        // Id запроса клиента
	UserTask  string  `json:"userTask"`  // Задание клиента
	Result    float64 `json:"result"`    // Результат,
	Status    string  `json:"status"`    // ""/"done"/"zero division"
	RootId    uint64  `json:"rootId"`    // Id корневого узла
	ParsError error   `json:"parsError"` // Ошибка парсинга
}
