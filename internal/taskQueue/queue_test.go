// Необходимо запускать все тесты последовательно
package taskQueue

import (
	"arithmometer/internal/entities"
	"fmt"
	"testing"
	"time"
)

var queue = NewQueue()

var task1 = NewTask(&entities.Node{Id: 1})
var task2 = NewTask(&entities.Node{Id: 2})
var task3 = NewTask(&entities.Node{Id: 3})

func TestQueue_AddRemoveTask(t *testing.T) {
	task := task1
	if !queue.AddTask(task) && queue.L != 1 {
		t.Errorf("invalid counter while addTask")
	}
	queue.RemoveTask(1)
	if queue.L != 0 {
		t.Errorf("invalid counter while removeTask from NotReady")
	}
}

func TestQueue_GetTask(t *testing.T) {
	var queue = NewQueue()
	_ = queue.AddTask(task1)
	ok := queue.AddTask(task2)
	if !ok && queue.L != 2 {
		t.Errorf("invalid counter while addTask, L=%d", queue.L)
	}
	result := queue.GetTask()
	if result != nil {
		t.Errorf("invalid GetTask while no ready task")
	}
	task1.XReady = true
	task1.YReady = true
	queue.UpdateReady()
	result = queue.GetTask()
	if result == nil || result.GetID() != 1 {
		t.Errorf("Test error: GetTask=%v, id=%d", result, result.GetID())
	}
	result.SetCalc(12)
	if queue.Working.Len() != 1 && queue.NotReady.Len() != 1 && len(queue.ReadyToCalc) != 0 {
		t.Errorf("len(Working)=%d, wont 1; len(NotReady)=%d, wont 1; len(ReadyToCalc)=%d, wont 0",
			queue.Working.Len(),
			queue.NotReady.Len(),
			len(queue.ReadyToCalc),
		)
	}
	_ = queue.AddTask(task3)
	if queue.L != 3 {
		t.Errorf("invalid counter while addTask")
	}
	fmt.Println(queue.Info())
	result.SetDeadline(1 * time.Second)
	result.SetDuration(2 * time.Second)
	queue.CheckDeadlines()
	time.Sleep(2 * time.Second)
	queue.CheckDeadlines()
	fmt.Println(queue.Info())
	fmt.Print(queue.NotReady.String())
	if queue.AddTask(task2) {
		t.Errorf("invalid add existing Task ")
	}
	queue.RemoveTask(1)
	if queue.L != 2 {
		t.Errorf("invalid counter while removeTask from Ready")
	}
	queue.RemoveTask(3)
	if queue.L != 1 {
		t.Errorf("invalid counter while removeTask from NotReady")
	}
	fmt.Println(queue.Info())
	task2.XReady, task2.YReady = true, true
	queue.GetTask()
	fmt.Println(queue.Info())
	queue.RemoveTask(2)
	fmt.Println(queue.Info())
}
func GetQoueue() *Queue {
	nodes := []*entities.Node{
		&entities.Node{Id: 1, Sheet: true, Calculated: true, Val: 1, Parent: 4},
		&entities.Node{Id: 2, Sheet: true, Calculated: true, Val: 2, Parent: 4},
		&entities.Node{Id: 3, Sheet: true, Calculated: true, Val: 3, Parent: 5},
		&entities.Node{Id: 4, Sheet: false, X: 1, Y: 2, Op: "*", Parent: 5},
		&entities.Node{Id: 5, Sheet: false, X: 4, Y: 3, Op: "+"},
	}
	q := NewQueue()
	q.AddExpressionNodes(nodes)
	return q
}

func TestQueue_CalculatingExpression(t *testing.T) {
	q := GetQoueue()
	fmt.Println(q.Info())
	q.UpdateReady()
	fmt.Println("обновление очереди")
	fmt.Println(q.Info())
	w := q.GetTask()
	fmt.Println("Получена задача из очереди, с id: ", w.Node.Id)
	ok, err := q.AddAnswer(4, 3)
	fmt.Println("добавлен результат вычисления", ok, err)
	fmt.Println(q.Info())
	q.UpdateReady()
	fmt.Println("обновление очереди")
	fmt.Println(q.Info())
	w = q.GetTask()
	fmt.Println("Получена задача из очереди, с id: ", w.Node.Id)
	ok, err = q.AddAnswer(5, 6)
	fmt.Println("добавлен результат вычисления", ok, err)
	fmt.Println(q.Info())
	if q.L != 0 {
		t.Errorf("ошибка счетчика после вычисления выражения")
	}
}
