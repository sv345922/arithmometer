package eval

import "fmt"

// вычисление узла, запускается только если посчитаны операнды
func (n *Node) Eval() (float64, error) {
	switch Op {
	case "+":
		return Sum(n.x, n.y), nil
	case "-":
		return Sub(n.x, n.y), nil
	case "*":
		return Mult(n.x, n.y), nil
	case "/":
		return Div(n.x, n.y), nil
	default:
		return 0, fmt.Errorf("invalid operator")
	}
}

// запускает вычисление дочерних узлов
func Count(n *Node) float64 {
	if n.Counted {
		return n.Val
	} else {
		n.x.Eval()
		n.y.Eval()
		n.Val = n.Eval()
		n.Counted = true
	}

}
