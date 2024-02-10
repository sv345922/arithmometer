package parsing

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	s := "-1+(2*3+4*5)+6"
	_, _, root, err := Parse(s)
	if fmt.Sprint(root) == "((-1.000000+((2.000000*3.000000)+(4.000000*5.000000)))+6.000000)" && err != nil {
		t.Error("parsing error")
	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(root)
}
