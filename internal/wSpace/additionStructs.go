package wSpace

import "arithmometer/pkg/timings"

// Для получения выражения от клиента
type NewExpr struct {
	Expr    string           `json:"expr"`
	Timings *timings.Timings `json:"timings"`
}

// Для получения ответа на задачу
type Answer struct {
	Result float64 `json:"result"`
	Err    error   `json:"err"`
}
type AnswerContainer struct {
	Id      string `json:"id"`
	AnswerN Answer `json:"answerN"`
}
