package stack

// допустимые типы
type additiveStack interface {
	any
	//treeExpression.Node | parser.Symbol
}

// Стэк для реализации алгоритма Дийксты
type Stack[T additiveStack] struct {
	val []*T
	ind int
}

func NewStack[T additiveStack](size int) *Stack[T] {
	return &Stack[T]{
		val: make([]*T, size),
		ind: -1,
	}
}

// Извлечь верхний элемент из стека и удалить его
func (s *Stack[T]) Pop() *T {
	if !s.IsEmpty() {
		res := s.val[s.ind]
		s.val[s.ind] = nil
		s.ind--
		return res
	} else {
		return nil
	}
}

// Добавить элемент в стек
func (s *Stack[T]) Push(l *T) {
	s.ind++
	s.val[s.ind] = l
}

// Проверить стек на пустоту
func (s *Stack[T]) IsEmpty() bool {
	return s.ind < 0
}

// Вернуть значение верхнего элемента стека
func (s *Stack[T]) Top() *T {
	if s.IsEmpty() {
		return nil
	}
	return s.val[s.ind]
}
