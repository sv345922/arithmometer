package treeExpression

import (
	"fmt"
	"sync"
	"testing"
)

func TestNodes(t *testing.T) {
	// Выражение 1 + 2 * 3, постфиксная запись 23*1+
	var TestingNodes = []*Node{
		{
			NodeId:     1,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        2.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     2,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        3.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     3,
			Op:         "*",
			X:          uint64(1),
			Y:          uint64(2),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     4,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        1.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     5,
			Op:         "+",
			X:          uint64(3),
			Y:          uint64(4),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(0),
			Mu:         sync.RWMutex{},
		},
	}
	node := NewNode()
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
	var ids []uint64
	for i, val := range TestingNodes {
		val.CreateId()
		if ok := nodesMap.Add(val); !ok {
			t.Errorf("invalid 'Add' function in %d case", i)
		}
		if !nodesMap.Contains(val.NodeId) {
			t.Errorf("check contains faildin %d case", i)
		}

		ids = append(ids, val.NodeId)
	}
	for i, val := range ids {
		res := nodesMap.Get(val)
		if res.NodeId != val {
			t.Errorf("invalid 'Get' function in %d case", i)
		}
	}
	// Добавляем существующий узел
	if nodesMap.Add(nodesMap.Get(ids[0])) == true {
		t.Errorf("invalid addin contains value")
	}
	nodesMap.Add(node)
	if !nodesMap.Remove(uint64(0)) {
		t.Errorf("invalid removing node")
	}
	fmt.Print(nodesMap.String())
}

func TestNode_GetParent(t *testing.T) {
	// Выражение 1 + 2 * 3, постфиксная запись 23*1+
	var TestingNodes = []*Node{
		{
			NodeId:     1,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        2.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     2,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        3.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     3,
			Op:         "*",
			X:          uint64(1),
			Y:          uint64(2),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     4,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        1.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     5,
			Op:         "+",
			X:          uint64(3),
			Y:          uint64(4),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(0),
			Mu:         sync.RWMutex{},
		},
	}
	nodesMap2 := NewNodes()
	for _, val := range TestingNodes {
		nodesMap2.Add(val)
	}
	node := TestingNodes[0]
	parent := node.GetParent(nodesMap2)
	if parent.NodeId != uint64(3) {
		t.Errorf("invalid check parent")
	}
}
func TestNode_IsReadyToCalc(t *testing.T) {
	// Выражение 1 + 2 * 3, постфиксная запись 23*1+
	var TestingNodes = []*Node{
		{
			NodeId:     1,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        2.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     2,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        3.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     3,
			Op:         "*",
			X:          uint64(1),
			Y:          uint64(2),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     4,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        1.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     5,
			Op:         "+",
			X:          uint64(3),
			Y:          uint64(4),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(0),
			Mu:         sync.RWMutex{},
		},
	}
	nodesMap3 := NewNodes()
	for _, val := range TestingNodes {
		nodesMap3.Add(val)
	}
	result := []bool{
		false,
		false,
		true,
		false,
		false,
	}
	for i, node := range TestingNodes {
		if node.IsReadyToCalc(nodesMap3) != result[i] {
			t.Errorf("IsReadyToCalc error on %d case in node %s",
				i,
				node.String(),
			)
		}

	}
}
func TestNode_GetType(t *testing.T) {
	// Выражение 1 + 2 * 3, постфиксная запись 23*1+
	var TestingNodes = []*Node{
		{
			NodeId:     1,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        2.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     2,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        3.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(3),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     3,
			Op:         "*",
			X:          uint64(1),
			Y:          uint64(2),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     4,
			Op:         "",
			X:          uint64(0),
			Y:          uint64(0),
			Val:        1.,
			Sheet:      true,
			Calculated: true,
			Parent:     uint64(5),
			Mu:         sync.RWMutex{},
		},
		{
			NodeId:     5,
			Op:         "+",
			X:          uint64(3),
			Y:          uint64(4),
			Val:        0.,
			Sheet:      false,
			Calculated: false,
			Parent:     uint64(0),
			Mu:         sync.RWMutex{},
		},
	}
	if TestingNodes[0].GetType() != "num" {
		t.Errorf("error get 'num' type")
	}
	if TestingNodes[2].GetType() != "Op" {
		t.Errorf("error get 'Op' type")
	}
}
