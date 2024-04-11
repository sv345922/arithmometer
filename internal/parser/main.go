package parser

import (
	"arithmometer/internal/entities"
	"errors"
)

type Expression struct {
	entities.Expression
	Postfix []*Symbol        `json:"postfix"` // Постфиксная запись выражения
	Times   entities.Timings `json:"times"`   // Тайминги
	Root    *entities.Node   `json:"root"`    // Корень дерева выражения
	Nodes   []*entities.Node `json:"nodes"`   // Узлы выражения
}

func NewExpression() *Expression {
	return &Expression{}
}

// Парсит выражения и заполняет поля структуры Expression, возвращает обшибку
func (e *Expression) Do(expr string, t entities.Timings) error {
	e.UserTask = expr
	e.Times = t
	// получаем корректные символы выражения
	symbols, err := Parse(expr)
	if err != nil {
		// ошибка выражения
		//log.Println(err)
		e.ParsError = err
		return err
	}
	// получаем постфиксную запись выражения
	e.Postfix, err = GetPostfix(symbols)
	if err != nil {
		// Ошибка скобок
		//log.Println(err)
		e.ParsError = errors.Join(e.ParsError, err)
		return err
	}
	// записываем поля Root, Nodes
	e.Root, e.Nodes, err = GetTree(e.Postfix)
	// создаем идентификатор выражения
	e.SetID()
	if err != nil {
		e.ParsError = errors.Join(e.ParsError, err)
		return err
	}
	// записываем идентификатор корня выражения
	e.RootId = e.Root.Id
	return nil
}
func (e *Expression) SetID() {
	for i, symbol := range e.Postfix {
		e.Id += uint64(i)
		for _, v := range []byte(symbol.Val) {
			e.Id += uint64(v)
		}
	}
	e.Id = e.Id + uint64(entities.GetDelta(3))
}
