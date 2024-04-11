package parser

import (
	"arithmometer/internal/entities"
	"errors"
	"fmt"
	"testing"
)

type parseTest struct {
	input string
	id    uint64
	err   error
}

var tests = []parseTest{
	{"", 0, nil},
	{"1", 49, nil},
	{"2 + 2 * 2", 245, nil},
	{"1 + 2 * (3 + 4) * 5", 461, nil},
	{"1+ & 2 * 3 + 4 * 5", 0, entities.ErrExpression},
	{"(2+2)*)5-3 * (2 + 4)", 0, entities.ErrExpression},
	{"-5 +3*7/0", 412, nil},
}
var exprs = make([]*Expression, len(tests))

func TestExpression(T *testing.T) {
	for i, test := range tests {
		expr := NewExpression()
		err := expr.Do(test.input)
		//fmt.Println(expr.Id)
		exprs[i] = expr
		if !errors.Is(err, test.err) || expr.Id != test.id {
			T.Errorf("%d: expect %d with err-'%v', got %d with err-'%v'", i, test.id, test.err, expr.Id, err)
		}
	}
}

var symbols = make([][]*Symbol, len(tests))
var postfix = make([][]*Symbol, len(tests))

func TestExpression_Do(t *testing.T) {
	// TODO
}
func TestParse(T *testing.T) {
	for i, test := range tests {
		s, err := Parse(test.input)
		symbols[i] = s
		_ = err
		fmt.Println(i+1, s, err)
	}
}

func TestGetPostfix(t *testing.T) {
	t.Run("parse", TestParse)
	for i, ps := range symbols {
		p, err := GetPostfix(ps)
		postfix[i] = p
		fmt.Println(i+1, p, err)
	}
}
func TestGetTree(t *testing.T) {
	t.Run("get postfix", TestGetPostfix)

	for i, ps := range postfix {
		_, nodes, _ := GetTree(ps)
		fmt.Print(i+1, "\t")
		for _, node := range nodes {
			fmt.Printf("%d ", node.Id)
		}
		fmt.Println()

	}
}
func TestGetNodes(t *testing.T) {
	n5 := Node{
		Op: "+",
	}
	n4 := Node{
		Op:     "*",
		Parent: &n5,
	}
	n1 := Node{
		Val:        2,
		Sheet:      true,
		Calculated: true,
		Parent:     &n4,
	}
	n2 := Node{
		Val:        2,
		Sheet:      true,
		Calculated: true,
		Parent:     &n4,
	}
	n3 := Node{
		Val:        2,
		Sheet:      true,
		Calculated: true,
		Parent:     &n5,
	}
	n5.X = &n4
	n5.Y = &n3
	n4.X = &n1
	n4.Y = &n2

	nodes := make([]*Node, 0)
	GetNodesOfExpression(&n5, &nodes)
	fmt.Println(nodes)
}
