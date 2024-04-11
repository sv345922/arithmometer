package calculator

import (
	"arithmometer/internal/entities"
	"log"
	"time"
)

var zeroDiv = entities.ZeroDiv
var OpEr = entities.OpEr

// Выбирает оператор вычисления, производит вычисления с учетом тайминга
// и возвращает результат с ошибкой
func Do(c *entities.MessageTask) (float64, error) {
	timings := c.Timings
	Op := c.Op
	x := c.X
	y := c.Y
	log.Println("задача получена", x, Op, y)
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
			return 0, zeroDiv
		}
		t := time.Duration(timings.Div)
		time.Sleep(time.Second * t)
		return x / y, nil
	default:
		return 0, OpEr
	}
}
