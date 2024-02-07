package eval

// Интерфейс арифметического выражения
type Expr interface {
	// метод, вычиляющий значение выражения
	Eval() (float64, error)
}

// Узел
type Node struct {
	Op      rune
	x, y    *Node
	Val     float64
	Counted bool
}
type Tree struct {
	nodes     []*Node
	relations map[int][]int
}

// переменные
type Env map[Var]float64

// Имя переменной
type Var string

// константа выражения
type Literal float64

// Бинарная операция
type binary struct {
	op   rune
	x, y Expr
}
