package newExpression

import (
	"arithmometer/internal/entities"
)

// Для получения выражения от клиента
type NewExpr struct {
	Expression string            `json:"expression"`
	Timings    *entities.Timings `json:"timings"`
}
