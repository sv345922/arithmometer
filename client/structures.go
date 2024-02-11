package main

import (
	"fmt"
)

type NewExp struct {
	Expr    string   `json:"expr"`
	Timings *Timings `json:"timings"`
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
