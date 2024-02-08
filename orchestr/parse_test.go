package orchestr

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	s := "-1 +(2 + 3) / 4 +5 "
	s = "-1+(2*3+4*5)+6"
	tree, err := Parse(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tree)
}
