package tasker

import (
	"arithmometer/orchestr/parsing"
)

type WorkingSpace struct {
	Tasks       *Tasks
	Expressions *Expressions
	Timings     *Timings
	AllNodes    *map[string]*parsing.Node
}

func (ws *WorkingSpace) save() error {
	db := DataBase{
		Expressions: ws.Expressions,
		Tasks:       ws.Tasks,
		Timings:     ws.Timings,
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	err := SafeJSON("db", db)
	if err != nil {
		return err
	}
	return nil
}

// Обновляет рабочее пространство
// TODO
func (ws *WorkingSpace) Update() {
	ws.Expressions.mu.RLock()
	defer ws.Expressions.mu.RUnlock()

	/*
		for _, expr := range ws.Expressions.ListExpr {
			if expr.Calculated() {

			}

		}
	*/
	// Сохранение рабочего пространства
	ws.save()
}

// При получении ответа от вычислителей обновить очереди

// При завершении вычисления выражения сохранить результат
