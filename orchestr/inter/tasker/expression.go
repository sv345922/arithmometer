package tasker

import (
	"arithmometer/orchestr/parsing"
	"log"
	"sync"
)

// Содержит выражения пользователя
// .Dict - словарь ссылок на выражения, ключ IdExpression
// .ListExpr - список с ссылками на выражения
// .mu - мьютекс для блокировки
type Expressions struct {
	Dict     map[uint64]*Expression `json:"dict"` // ключ - id выражения/запроса клиента
	ListExpr []*Expression          `json:"listExpr"`
	mu       sync.RWMutex           `json:"-"`
}

// Создает новый список выражений
// .Dict - словарь с ссылками на выражения
// .ListExpr - список с ссылками на выраженя, повторяет .Dict
func NewExpressions() *Expressions {
	res := Expressions{}
	res.Dict = make(map[uint64]*Expression)
	res.ListExpr = make([]*Expression, 0)
	return &res
}

// Обновляет список выражений, при вычисленном корне выражения, либо делении на ноль,
// ставит статус вычислено/деление на ноль и результат вычислений
func (es *Expressions) UpdateStatus(root *parsing.Node, status string, result float64) {
	es.mu.Lock()
	defer es.mu.Unlock()
	if expression, ok := es.Dict[root.NodeId]; ok {
		expression.mu.Lock()
		expression.Status = status
		expression.Result = result
		expression.mu.Unlock()
	} else {
		log.Println("Выражение не найдено")
	}
}

// Добавляет выражение в список выражений
func (es *Expressions) Add(expression *Expression) {
	es.mu.Lock()
	defer es.mu.Unlock()
	l_prev := len(es.Dict)
	es.Dict[expression.IdExpression] = expression
	// Проверка на наличие идентичного выражения среди имеющихся
	// если уже есть, то длина мапы не изменится, и значит добавлять в список выражение
	// не надо
	if l_prev+1 == len(es.Dict) {
		es.ListExpr = append(es.ListExpr, expression)
	}
}

// возвращает выражение из списка задач
func FindExpression(id uint64, e *Expressions) *Expression {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if task, ok := e.Dict[id]; ok {
		return task
	}
	return nil
}

// Выражение
type Expression struct {
	IdExpression uint64            `json:"id"`        // id запроса клиента
	UserTask     string            `json:"userTask"`  // задание клиента
	Postfix      []*parsing.Symbol `json:"postfix"`   // постфиксная запись выражения
	Times        Timings           `json:"times"`     // тайминги
	Result       float64           `json:"result"`    // результат,
	Status       string            `json:"status"`    // ""/"done"/"деление на ноль"
	RootId       uint64            `json:"rootId"`    // id кореневого узла
	ParsError    error             `json:"parsError"` // Ошибка парсинга
	mu           sync.Mutex
}

// Создает id выражения
func (e *Expression) CreateId() {
	e.mu.Lock()
	s := ""
	for _, symbol := range e.Postfix {
		s = s + symbol.Val
	}
	e.IdExpression = parsing.GetId(s)
}

// Определяет, вычислено ли выражение
func (e *Expression) Calculated() bool {
	if e.Status == "done" {
		return true
	}
	return false
}
