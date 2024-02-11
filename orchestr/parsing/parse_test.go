package parsing

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	s := "-1+(2*3+4*5)+6"
	symbols, err := Parse(s)
	var postfix string
	for _, val := range symbols {
		postfix += val.Val
	}
	if postfix == "1-23*45*++6+" && err != nil {
		t.Error("parsing error")
	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(postfix)
}
