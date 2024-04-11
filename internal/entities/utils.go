package entities

import (
	"math"
	"time"
)

// Возвращает идентификатор (число) заданной длины (не более n), зависящий от времени
func GetDelta(n int) int {
	k := int64(math.Pow10(n + 2))
	return int((time.Now().UnixNano() % k) / 100)
}
