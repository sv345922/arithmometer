package entities

import (
	"time"
)

type Task struct {
	Id       uint64        `json:"id"`        // идентификатор
	ParentId uint64        `json:"parent_id"` // id родительского выражения
	X        float64       `json:"x"`         // операнд X
	Xid      uint64        `json:"xid"`
	XReady   bool          `json:"x_ready"` // готовность операнда X
	Op       string        `json:"op"`      // операция
	Y        float64       `json:"y"`       // операнд Y
	Yid      uint64        `json:"yid"`
	YReady   bool          `json:"y_ready"`  // готовность операнда Y
	Error    error         `json:"error"`    // ошибка
	CalcId   uint64        `json:"calc_id"`  // id вычислителя задачи
	Deadline time.Time     `json:"deadline"` // дедлайн задачи
	Duration time.Duration `json:"duration"` // длительность операции
}
