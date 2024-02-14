package tasker

import (
	"arithmometer/orchestr/parsing"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sync"
)

// Контейнер задач для отправки вычислителю
type TaskContainer struct {
	// идентификатор
	Id string `json:"id"`
	// задача
	TaskN Task `json:"taskN"`
	// ощибка
	Err error `json:"err"`
	// тайминги
	TimingsN Timings `json:"timingsN"`
}

// Зачада для вычислителя
type Task struct {
	// операнд X
	X float64 `json:"x"`
	// операнд Y
	Y float64 `json:"y"`
	// операция
	Op string `json:"op"`
}

// Содержит задачи для вычислителей
// .Dict - словарь задач, ключ Id
// .mu - мьютекс для блокировки словаря
type Tasks struct {
	Dict map[string]*TaskContainer `json:"dict"`
	mu   sync.RWMutex              `json:"-"`
}

// Потокобезопасно добавляет задачу в список задач
func (t *Tasks) Add(task TaskContainer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Dict[task.Id] = &task
}

// Содержит выражения пользователя
// .Dict - словарь ссылок на выражения, ключ Id
// .ListExpr - список с ссылками на выражения
// .mu - мьютекс для блокировки
type Expressions struct {
	Dict     map[string]*Expression `json:"dict"`
	ListExpr []*Expression          `json:"listExpr"`
	mu       sync.RWMutex           `json:"-"`
}

/*
	func (e *Expressions) Remove(id string) error {
		e.mu.Lock()
		defer e.mu.Unlock()
		delete(e.dict, id)
		e.listExpr = append(e.listExpr[:index], e.listExpr[index+1:]...)
	}
*/
// Выражение
type Expression struct {
	Id       string            `json:"id"`
	UserTask string            `json:"userTask"` // задание клиента
	Postfix  []*parsing.Symbol `json:"postfix"`  // постфиксная запись выражения
	Times    Timings           `json:"times"`    // тайминги
	Result   string            `json:"result"`   // результат, возможно нужно float64
	Status   string            `json:"status"`   // "done" - рассчитано
	Root     *parsing.Node     //корень дерева выражения TODO нужно заводить возможно
}

// Создает id выражения
func (e *Expression) CreateId() {
	s := ""
	for _, symbol := range e.Postfix {
		s += symbol.Val
	}
	hasher := sha1.New()
	hasher.Write([]byte(s))
	e.Id = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// Определяет, вычислено ли выражение
func (e *Expression) Calculated() bool {
	if e.Status == "done" {
		return true
	}
	return false
}

// Тайминги для операторов
type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

// Стрингер
func (t *Timings) String() string {
	return fmt.Sprintf("+: %d s, -: %d s, *: %d s, /: %d s", t.Plus, t.Minus, t.Mult, t.Div)
}

type NewExpr struct {
	Expr    string   `json:"expr"`
	Timings *Timings `json:"timings"`
}
type Answer struct {
	Result float64 `json:"result"`
	Err    error   `json:"err"`
}
type AnswerContainer struct {
	Id      string `json:"id"`
	AnswerN Answer `json:"answerN"`
}
