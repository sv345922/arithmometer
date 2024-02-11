package calc

type Task struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Op string  `json:"op"`
}
type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}
type TaskContainer struct {
	Id       string  `json:"id"`
	TaskN    Task    `json:"taskN"`
	TimingsN Timings `json:"timingsN"`
}
type Answer struct {
	Result float64 `json:"result"`
	Err    error   `json:"err"`
}
type AnswerContainer struct {
	Id      string `json:"id"`
	AnswerN Answer `json:"answerN"`
}
