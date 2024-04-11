package entities

type MessageResult struct {
	Id     uint64  `json:"id"`
	Result float64 `json:"result"`
	Err    error   `json:"err"`
}
type MessageTask struct {
	Id      uint64   `json:"id"`
	X       float64  `json:"x"`
	Y       float64  `json:"y"`
	Op      string   `json:"op"`
	Timings *Timings `json:"timings"`
}
