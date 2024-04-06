// Необходимо запускать все тесты последовательно
package taskQueue

import (
	"fmt"
	"testing"
	"time"
)

type Task struct {
	Id       uint64
	Deadline time.Time
	calcId   uint64
}

func (t *Task) IsReadyToCalc() bool {
	//TODO возможно другие варианты нужно использовать на тесте
	if (t.Id+t.calcId)%2 == 0 {
		return true
	}
	return false
}

func (t *Task) IsTimeout() bool {
	if time.Now().After(t.Deadline) {
		return true
	}
	return false
}
func (t *Task) SetDeadline(duration time.Duration) {
	t.Deadline.Add(duration)
}
func (t *Task) GetID() uint64 {
	return t.Id
}
func (t *Task) SetCalc(id uint64) {
	t.calcId = id
}
func (t *Task) String() string {
	return fmt.Sprintf("Id=%d, deadline=%v, calcId=%d",
		t.Id,
		t.Deadline,
		t.calcId,
	)
}

type Test struct {
	task   *Task
	result bool
}

// var task Element = Taks{}
var cases = []Test{
	{&Task{1, time.Now(), 0}, true},
	{&Task{2, time.Now(), 0}, false},
}
var tasks = NewTasks()

//task := cases[0].task

func TestTasks_AddTask(t *testing.T) {
	var task Element = &Task{uint64(1), time.Now(), uint64(0)}
	_ = tasks.AddTask(&task)
	if tasks.L != 1 {
		t.Errorf("invalid counter while addTask")
	}
}
func TestTasks_RemoveTask(t *testing.T) {
	tasks.RemoveTask(uint64(1))
	if tasks.L != 0 {
		t.Errorf("invalid counter while removeTask from NotReady")
	}
}
func TestTasks_GetTask(t *testing.T) {
	var task1 Element = &Task{uint64(1), time.Now(), uint64(0)}
	_ = tasks.AddTask(&task1)
	var task2 Element = &Task{uint64(1), time.Now(), uint64(0)}
	_ = tasks.AddTask(&task2)
	if tasks.L != 2 {
		t.Errorf("invalid counter while addTask")
	}
	result := tasks.GetTask()
	if (*result).GetID() != 1 {
		t.Errorf("GetTask error")
	}
	if len(tasks.Working) != 1 && len(tasks.NotReady) != 1 && len(tasks.Waiting) != 0 {
		t.Errorf("len(Working)=%d, wont 1; len(NotReady)=%d, wont 1; len(Waiting)=%d, wont 0",
			len(tasks.Working),
			len(tasks.NotReady),
			len(tasks.Waiting),
		)
	}
	var task3 Element = &Task{uint64(3), time.Now(), uint64(0)}
	_ = tasks.AddTask(&task3)
	if tasks.L != 3 {
		t.Errorf("invalid counter while addTask")
	}
	tasks.String()
}
