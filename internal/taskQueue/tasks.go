package taskQueue

import (
	"fmt"
	"sync"
)

type Tasks struct {
	AllTask map[uint64]*Task
	mu      sync.RWMutex
}

// Создать новой список узлов
func NewTasks() *Tasks {
	nodes := make(map[uint64]*Task)
	return &Tasks{AllTask: nodes}
}

// Возвращает длину словаряя
func (t *Tasks) Len() int {
	return len(t.AllTask)
}

// Если узел есть в списке, возвращает true, иниче false
func (t *Tasks) Contains(key uint64) bool {
	//t.mu.RLock()
	//defer t.mu.RUnlock()
	if _, ok := t.AllTask[key]; ok {
		return true
	}
	return false
}

// Добавить узел в список узлов, если узел с таким id существует в списке, возвращает false
func (t *Tasks) Add(n *Task) bool {
	key := n.Node.Id
	if t.Contains(key) {
		return false
	}
	//t.mu.Lock()
	t.AllTask[key] = n
	//t.mu.Unlock()
	return true
}

// возвращает узел из исписка по его id, если id нет в списке узлов, возвращает nil
func (t *Tasks) Get(key uint64) *Task {
	if t.Contains(key) {
		//t.mu.RLock()
		result := t.AllTask[key]
		//t.mu.RUnlock()
		return result
	}
	return nil
}

// Удалить узел из списка узлов по его id, если узла с таким id нет, возвращает false
func (t *Tasks) Remove(key uint64) bool {
	if t.Contains(key) {
		t.mu.Lock()
		defer t.mu.Unlock()
		delete(t.AllTask, key)
		return true
	}
	return false
}

// Стрингер
func (t *Tasks) String() string {
	result := ""
	for key, val := range t.AllTask {
		result += fmt.Sprintf("key=%d, val=%.2f%s%.2f\n",
			key,
			val.X,
			val.Node.Op,
			val.Y)
	}
	return result
}

// Возвращает id всех узлов
func (t *Tasks) GetAllIDs() []uint64 {
	result := make([]uint64, len(t.AllTask))
	i := 0
	//t.mu.RLock()
	//defer t.mu.RUnlock()
	for key := range t.AllTask {
		result[i] = key
		i++
	}
	return result
}
