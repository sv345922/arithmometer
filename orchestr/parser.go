package orchestr

import (
	"fmt"
	//"strconv"
	"strings"
	"text/scanner"
)

// Функция рекурсивно создает дерево выражения и возвращает его корень и ошибку
func GetTree(inputList []string) (Tree, error) {
	tree := Tree {nodes: make([]*Node, 0), relations: make(map[int][]int)}
	for i, symb := range inputList {
		
	}
	
	
	// TODO
	return tree, nil
}

func Parse(input string) (*Tree, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanFloats | scanner.ScanInts // | scanner.ScanIdents
	inputList := make([]string, 0)
	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		text := s.TokenText()
		switch token {
		case scanner.Int, scanner.Float:
			inputList = append(inputList, text)
		default:
			switch text{
			case "+", "-", "*", "/", "(", ")":
				inputList = append(inputList, text)
			default:
				return nil, fmt.Errorf("invalid expression: %s", text)
			}
		}	
	}	
	tree, err := GetTree(inputList)
	if err != nil {
		return nil, err
	}
	return &tree, nil
}


func createChild(n *Node) *Node {
	child := &Node{x: n.x, y: n.y, op: n.op, level: n.level+1, parent: n}
	return child
}

