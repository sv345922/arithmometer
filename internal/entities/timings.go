package entities

import (
	"arithmometer/internal/configs"
	"fmt"
	"time"
)

var tConst = configs.TConst

// Timings Тайминги для операторов
type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

// Стрингер
func (t *Timings) String() string {
	return fmt.Sprintf("+: %ds, -: %ds, *: %ds, /: %ds", t.Plus, t.Minus, t.Mult, t.Div)
}

// GetDuration Возвращает время выполнения конкретной операции
// Если оператор неизвестен, возвращает нулевую длительность
func (t *Timings) GetDuration(op string) time.Duration {
	switch op {
	case "+":
		return time.Duration(t.Plus) * tConst
	case "-":
		return time.Duration(t.Minus) * tConst
	case "*":
		return time.Duration(t.Mult) * tConst
	case "/":
		return time.Duration(t.Div) * tConst
	default:
		return 0 * tConst
	}
}
