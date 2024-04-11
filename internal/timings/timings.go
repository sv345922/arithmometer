package timings

import (
	"arithmometer/internal/entities"
	"fmt"
	"time"
)

// TODO  не используется, можно удалить файл
const T_const = time.Second

// Timings Тайминги для операторов
type Timings struct {
	entities.Timings
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
