package calc

import (
	"fmt"
	"log"
	"time"
)

// Выбирает оператор вычисления, производит вычисления с учетом тайминга
// и возвращает результат с ошибкой
func Calculate(c *TaskContainer) (float64, error) {
	timings := c.TimingsN
	n := c.TaskN
	Op := n.Op
	x := n.X
	y := n.Y
	log.Println("вычислитель получил задачу:", x, Op, y)
	switch Op {
	case "+":
		t := time.Duration(timings.Plus)
		time.Sleep(time.Second * t)
		return x + y, nil
	case "-":
		t := time.Duration(timings.Minus)
		time.Sleep(time.Second * t)
		return x - y, nil
	case "*":
		t := time.Duration(timings.Mult)
		time.Sleep(time.Second * t)
		return x * y, nil
	case "/":
		if y == 0 {
			return 0, fmt.Errorf("division zero")
		}
		t := time.Duration(timings.Div)
		time.Sleep(time.Second * t)
		return x / y, nil
	default:
		return 0, fmt.Errorf("invalid operator")
	}
}
