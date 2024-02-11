package handler

import (
	"arithmometer/orchestr/parsing"
	"fmt"
)

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

func (t *Timings) String() string {
	return fmt.Sprintf("+: %d s, -: %d s, *: %d s, /: %d s", t.Plus, t.Minus, t.Mult, t.Div)
}

type dataBase struct {
	// TODO
	// список выражений (с таймингами)
	// список задач (ожидающие/выполняющиеся/готовые)

}
type NewExpr struct {
	Expr    string   `json:"expr"`
	Timings *Timings `json:"timings"`
}
