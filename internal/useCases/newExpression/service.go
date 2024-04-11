package newExpression

import (
	"arithmometer/internal/entities"
	"arithmometer/internal/parser"
)

const t = 5

func TransformParseExpression(node *parser.Expression) *entities.Expression {
	return &entities.Expression{
		Id:        node.Id,
		UserTask:  node.UserTask,
		Result:    node.Result,
		Status:    node.Status,
		RootId:    node.RootId,
		ParsError: node.ParsError,
	}
}

func setTimingsWhileZero(n *NewExpr) {
	times := n.Timings
	if times.Plus == 0 {
		times.Plus = t
	}
	if times.Minus == 0 {
		times.Mult = t
	}
	if times.Mult == 0 {
		times.Mult = t
	}
	if times.Div == 0 {
		times.Mult = t
	}
}
