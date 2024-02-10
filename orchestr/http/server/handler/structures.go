package handler

import "arithmometer/orchestr/parsing"

type Expression struct {
	Id      string            `json:"id"`
	Postfix []*parsing.Symbol `json:"postfix"`
	Times   Timings           `json:"times"`
}

func (e *Expression) doId() {
	for _, symbol := range e.Postfix {
		e.Id += symbol.Val
	}
}

type additiveJSON interface {
	Expression | []*parsing.Node | *parsing.Node | []*parsing.Symbol | Timings
}

type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

type dataBase struct {
	// TODO
	// список выражений (с таймингами)
	// список задач (ожидающие/выполняющиеся/готовые)

}
