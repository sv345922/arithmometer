package expressions

import (
	"arithmometer/pkg/parser"
	"arithmometer/pkg/timings"
	"arithmometer/pkg/treeExpression"
	"log"
	"sync"
)

// Expression Выражение
type Expression struct {
	IdExpression uint64           `json:"id"`        // Id запроса клиента
	UserTask     string           `json:"userTask"`  // Задание клиента
	Postfix      []*parser.Symbol `json:"postfix"`   // Постфиксная запись выражения
	Times        timings.Timings  `json:"times"`     // Тайминги
	Result       float64          `json:"result"`    // Результат,
	Status       string           `json:"status"`    // ""/"done"/"деление на ноль"
	RootId       uint64           `json:"rootId"`    // Id корневого узла
	ParsError    error            `json:"parsError"` // Ошибка парсинга
	mu           sync.Mutex
}

// CreateId Создает id выражения
func (e *Expression) CreateId() {
	e.mu.Lock()
	s := ""
	for _, symbol := range e.Postfix {
		s = s + symbol.Val
	}
	e.IdExpression = treeExpression.NewId(s)
}

// Calculated Определяет, вычислено ли выражение
func (e *Expression) Calculated() bool {
	if e.Status == "done" {
		return true
	}
	return false
}

// Expressions Содержит выражения пользователя
// Dict - словарь ссылок на выражения, ключ IdExpression
// ListExpr - список со ссылками на выражения
// mu - мьютекс для блокировки
type Expressions struct {
	Dict     map[uint64]*Expression `json:"dict"` // Ключ - id выражения/запроса клиента
	ListExpr []*Expression          `json:"listExpr"`
	mu       sync.RWMutex           `json:"-"`
}

// NewExpressions Создает новый список выражений
// Dict - словарь со ссылками на выражения
// ListExpr - список со ссылками на выражения, повторяет Dict
func NewExpressions() *Expressions {
	return &Expressions{
		Dict:     make(map[uint64]*Expression),
		ListExpr: make([]*Expression, 0),
		mu:       sync.RWMutex{},
	}
}

// UpdateStatus Обновляет список выражений, при вычисленном корне выражения, либо делении на ноль,
// ставит статус вычислено/деление на ноль и результат вычислений
func (es *Expressions) UpdateStatus(root *treeExpression.Node, status string, result float64) {
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

// Add Добавляет выражение в список выражений
func (es *Expressions) Add(expression *Expression) {
	es.mu.Lock()
	defer es.mu.Unlock()
	lenPrev := len(es.Dict)
	es.Dict[expression.IdExpression] = expression
	// Проверка на наличие идентичного выражения среди имеющихся
	// если уже есть, то длина мапы не изменится, и значит добавлять в список выражение
	// не надо
	if lenPrev+1 == len(es.Dict) {
		es.ListExpr = append(es.ListExpr, expression)
	}
}

// FindExpression возвращает выражение из списка задач
func FindExpression(id uint64, e *Expressions) *Expression {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if task, ok := e.Dict[id]; ok {
		return task
	}
	return nil
}
