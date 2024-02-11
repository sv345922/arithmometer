package handler

import (
	"arithmometer/orchestr/parsing"
	"testing"
)

var case1 = &parsing.Node{
	X:  &parsing.Node{X: nil, Y: nil, Sheet: true, Val: 1},
	Y:  &parsing.Node{X: nil, Y: nil, Sheet: true, Val: 2},
	Op: "+"}

func TestSafeJSON(t *testing.T) {
	err := SafeJSON("test", case1)
	if err != nil {
		t.Error(err)
	}
}
