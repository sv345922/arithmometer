package orchestr

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	s := "-1 +2 - 3 +4 "
	tree, err := Parse(s)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(tree)
}
