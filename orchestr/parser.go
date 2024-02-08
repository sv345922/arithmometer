package orchestr

import (
	"fmt"
	"strconv"

	//"strconv"
	"strings"
	"text/scanner"
)

// парсит строку в дерево
func Parse(input string) (*Node, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanFloats | scanner.ScanInts // | scanner.ScanIdents
	var NodeList []*Node
	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		text := s.TokenText()
		switch token {
		case scanner.Int, scanner.Float:
			num, _ := strconv.ParseFloat(text, 64)
			NodeList = append(NodeList, &Node{val: num, sheet: true})
		default:
			switch text {
			case "+", "-", "*", "/", "(", ")":
				NodeList = append(NodeList, &Node{op: text})
			default:
				return nil, fmt.Errorf("invalid expression: %s", text)
			}
		}
	}
	postFix, err := getPostfix(NodeList)
	if err != nil {
		return nil, err
	}
	root, err := getTree(postFix)
	if err != nil {
		return nil, err
	}
	return root, nil
}
func getPostfix(input []*Node) ([]*Node, error) {
	var opStack stack
	var postFix []*Node

	for _, val := range input {
		switch val.getType() {
		case "num":
			postFix = append(postFix, val)
		case "op":
			switch val.op {
			case "(":
				opStack.push(val)
			case ")":
				for n := opStack.pop(); n.op != "("; n = opStack.pop() {
					if n == nil {
						return nil, fmt.Errorf("invalid paranthesis")
					}
					postFix = append(postFix, n)
				}
			default: // val оператор
				priorCur := val.getPriority()
				for !opStack.isEmpty() && opStack.get().getPriority() >= priorCur {
					postFix = append(postFix, opStack.pop())
				}
				opStack.push(val)
			}
		}
	}
	for !opStack.isEmpty() {
		postFix = append(postFix, opStack.pop())
	}
	return postFix, nil
}

// Строит дерево выражения и возвращает корневой узел
func getTree(postfix []*Node) (*Node, error) {
	if len(postfix) == 0 {
		return nil, fmt.Errorf("expression is empty")
	}
	stack := new(stack)
	for _, node := range postfix {
		if node.getType() == "op" {
			y := stack.pop()
			x := stack.get()
			// при первом отрицательном числе
			if x == nil {
				node = y
				node.val = -node.val
			} else {
				node.x = stack.pop()
				node.y = y
			}
			stack.push(node)
		} else {
			stack.push(node)
		}
	}
	return stack.get(), nil
}
