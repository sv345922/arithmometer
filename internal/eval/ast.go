package eval

// Интерфейс арифметического выражения
type Expr interface {
	// метод, вычиляющий значение выражения
	Eval(env Env) float64
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
