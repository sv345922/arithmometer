package parser

import (
	"fmt"
	"strings"
	"text/scanner"

	"arithmometer/internal/eval"
)

// тип ошибки оработки лексем
type lexPanic string

// лексер
type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}
func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

// парсер
func Parce(input string) (_ eval.Expr, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil: // нет паники
		case lexPanic:
			err = fmt.Errorf("%s".x)
		default:
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	lex.next()
	e := parceExpr(lex)
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex.describe())
	}
}
