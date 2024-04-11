package treeExpression

import (
	"fmt"
	"sync"
)

type Nodes struct {
	AllNodes map[uint64]*Node
	mu       sync.RWMutex
}

// Создать новой список узлов
func NewNodes() *Nodes {
	nodes := make(map[uint64]*Node)
	return &Nodes{AllNodes: nodes}
}

// Возвращает длину словаряя
func (ns *Nodes) Len() int {
	return len(ns.AllNodes)
}

// Если узел есть в списке, возвращает true, иниче false
func (ns *Nodes) Contains(key uint64) bool {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	if _, ok := ns.AllNodes[key]; ok {
		return true
	}
	return false
}

// Добавить узел в список узлов, если узел с таким id существует в списке, возвращает false
func (ns *Nodes) Add(n *Node) bool {
	key := n.Id
	if ns.Contains(key) {
		return false
	}
	ns.mu.Lock()
	ns.AllNodes[key] = n
	ns.mu.Unlock()
	return true
}

// возвращает узел из исписка по его id, если id нет в списке узлов, возвращает nil
func (ns *Nodes) Get(key uint64) *Node {
	if ns.Contains(key) {
		ns.mu.RLock()
		defer ns.mu.RUnlock()
		return ns.AllNodes[key]
	}
	return nil
}

// Удалить узел из списка узлов по его id, если узла с таким id нет, возвращает false
func (ns *Nodes) Remove(key uint64) bool {
	if ns.Contains(key) {
		ns.mu.Lock()
		delete(ns.AllNodes, key)
		return true
	}
	return false
}

// Стрингер
func (ns *Nodes) String() string {
	result := ""
	for key, val := range ns.AllNodes {
		result += fmt.Sprintf("key=%d, val=%s\n",
			key,
			val.getVal())
	}
	return result
}

// Возвращает id всех узлов
func (ns *Nodes) GetAllIDs() []uint64 {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	result := make([]uint64, 0, len(ns.AllNodes))
	for key := range ns.AllNodes {
		result = append(result, key)
	}
	return result
}
