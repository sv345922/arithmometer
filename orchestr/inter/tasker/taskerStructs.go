package tasker

import (
	"arithmometer/orchestr/parsing"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sync"
)

// Зачада для вычислителя
type TaskContainer struct {
	Id       string  `json:"id"`
	TaskN    Task    `json:"taskN"`
	Err      error   `json:"err"`
	TimingsN Timings `json:"timingsN"`
}
type Task struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Op string  `json:"op"`
}

// Содержит задачи для вычислителей
type Tasks struct {
	Dict map[string]*TaskContainer `json:"dict"`
	mu   sync.RWMutex              `json:"-"`
}

func (t *Tasks) Add(task TaskContainer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Dict[task.Id] = &task
}

// Содержит список выражений пользователя
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
type Expression struct {
	Id       string            `json:"id"`
	UserTask string            `json:"userTask"` // задание клиента
	Postfix  []*parsing.Symbol `json:"postfix"`
	Times    Timings           `json:"times"`
	Result   string            `json:"result"`
	Status   string            `json:"status"` // "done" - рассчитано
	Root     *parsing.Node
}

func (e *Expression) CreateId() {
	s := ""
	for _, symbol := range e.Postfix {
		s += symbol.Val
	}
	hasher := sha1.New()
	hasher.Write([]byte(s))
	e.Id = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (e *Expression) Calculated() bool {
	if e.Status == "done" {
		return true
	}
	return false
}

type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

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
