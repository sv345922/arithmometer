package treeExpression

import (
	"arithmometer/internal/entities"
	"fmt"
	"testing"
)

// Выражение 1 + 2 * 3, постфиксная запись 23*1+
var TestingNodes = []entities.Node{
	{
		Id:         1,
		Op:         "",
		X:          uint64(0),
		Y:          uint64(0),
		Val:        2.,
		Sheet:      true,
		Calculated: true,
		Parent:     uint64(3),
	},
	{
		Id:         2,
		Op:         "",
		X:          uint64(0),
		Y:          uint64(0),
		Val:        3.,
		Sheet:      true,
		Calculated: true,
		Parent:     uint64(3),
	},
	{
		Id:         3,
		Op:         "*",
		X:          uint64(1),
		Y:          uint64(2),
		Val:        0.,
		Sheet:      false,
		Calculated: false,
		Parent:     uint64(5),
	},
	{
		Id:         4,
		Op:         "",
		X:          uint64(0),
		Y:          uint64(0),
		Val:        1.,
		Sheet:      true,
		Calculated: true,
		Parent:     uint64(5),
	},
	{
		Id:         5,
		Op:         "+",
		X:          uint64(3),
		Y:          uint64(4),
		Val:        0.,
		Sheet:      false,
		Calculated: false,
		Parent:     uint64(0),
	},
}

func TestNodes(t *testing.T) {

	node := NewNode(&entities.Node{})
	if fmt.Sprintf("%T", node) != "*treeExpression.Node" {
		t.Errorf("invalid creating Node")
	}
	nodesMap := NewNodes()
	if fmt.Sprintf("%T", nodesMap) != "*treeExpression.Nodes" {
		t.Errorf("invalid creating Nodes")
	}
	if nodesMap.Get(uint64(100)) != nil {
		t.Errorf("invalid 'Get' in empty Nodes")
	}
	if ok := nodesMap.Remove(uint64(100)); ok {
		t.Errorf("invalid 'Remove' in empty Nodes")
	}
	var testingNodes = make([]entities.Node, len(TestingNodes))
	copy(testingNodes, TestingNodes)
	for i, val := range testingNodes {
		node = NewNode(&val)
		id := node.CreateId()
		if ok := nodesMap.Add(node); !ok {
			t.Errorf("invalid 'Add' function in %d case", i)
		}
		if !nodesMap.Contains(id) {
			t.Errorf("check contains faild in %d case", i)
		}
	}
	ids := nodesMap.GetAllIDs()
	for i, val := range ids {
		res := nodesMap.Get(val)
		if res.Id != val {
			t.Errorf("invalid 'Get' function in %d case", i)
		}
	}
	// Добавляем существующий узел
	if nodesMap.Add(nodesMap.Get(ids[0])) == true {
		t.Errorf("invalid addin contains value")
	}
	nodesMap.Add(node)
	if nodesMap.Remove(ids[0]) && nodesMap.Len() != 4 {
		t.Errorf("invalid removing node, len=%d (wont 4)", nodesMap.Len())
	}
	fmt.Print(nodesMap.String())
}

func TestNode_GetParent(t *testing.T) {
	var testingNodes = make([]entities.Node, len(TestingNodes))
	copy(testingNodes, TestingNodes)

	nodesMap2 := NewNodes()
	for _, val := range testingNodes {
		nodesMap2.Add(NewNode(&val))
	}
	node := NewNode(&testingNodes[0])
	parent := node.GetParent(nodesMap2)
	if parent.Id != uint64(3) {
		t.Errorf("invalid check parent")
	}
}
func TestNode_IsReadyToCalc(t *testing.T) {
	var testingNodes = make([]entities.Node, len(TestingNodes))
	copy(testingNodes, TestingNodes)

	nodesMap3 := NewNodes()
	for _, val := range testingNodes {
		nodesMap3.Add(NewNode(&val))
	}
	result := []bool{
		false,
		false,
		true,
		false,
		false,
	}
	for i, node := range testingNodes {
		if NewNode(&node).IsReadyToCalc(nodesMap3) != result[i] {
			t.Errorf("IsReadyToCalc error on %d case in node %s",
				i,
				NewNode(&node).String(),
			)
		}
	}
}
func TestNode_GetType(t *testing.T) {
	var testingNodes = make([]entities.Node, len(TestingNodes))
	copy(testingNodes, TestingNodes)

	if NewNode(&testingNodes[0]).GetType() != "num" {
		t.Errorf("error get 'num' type")
	}
	if NewNode(&testingNodes[2]).GetType() != "Op" {
		t.Errorf("error get 'Op' type")
	}
}
