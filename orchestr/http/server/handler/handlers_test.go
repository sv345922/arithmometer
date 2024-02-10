package handler

import (
	"arithmometer/orchestr/parsing"
	"fmt"
	"testing"
)

var case1 = &parsing.Node{
	X:  &parsing.Node{X: nil, Y: nil, Sheet: true, Val: 1},
	Y:  &parsing.Node{X: nil, Y: nil, Sheet: true, Val: 2},
	Op: "+"}

func TestGetNodes(t *testing.T) {
	var nodes []*parsing.Node
	nodes = getNodes(case1, nodes)
	if len(nodes) != 3 {
		t.Error("error")
	}
	fmt.Println(nodes)
}
func TestSafeJSON(t *testing.T) {
	err := SafeJSON("test", case1)
	if err != nil {
		t.Error(err)
	}
}
