package timings

import (
	"fmt"
	"time"
)

const T_const = time.Second

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
		return time.Duration(t.Plus) * T_const
	case "-":
		return time.Duration(t.Minus) * T_const
	case "*":
		return time.Duration(t.Mult) * T_const
	case "/":
		return time.Duration(t.Div) * T_const
	default:
		return 0 * T_const
	}
}
