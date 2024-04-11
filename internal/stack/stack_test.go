package stack

import (
	"reflect"
	"testing"
)

type testCase struct {
	input  []int
	output []int
}

var testCases = []testCase{
	{input: []int{1, 2, 3, 4, 5, 6}, output: []int{6, 5, 4, 3, 2, 1}},
}

func TestStack(t *testing.T) {
	stack := NewStack[int](10)
	if !stack.IsEmpty() {
		t.Errorf("stack.IsEmpty() = false, want true")
	}
	for _, testCase := range testCases {
		for i := 0; i < len(testCase.input); i++ {
			stack.Push(&testCase.input[i])
			if stack.IsEmpty() {
				t.Errorf("stack.IsEmpty() = true, want false")
			}
			if stack.Top() != &testCase.input[i] {
				t.Errorf("stack.Top() = %d, want %d", stack.Top(), &testCase.input[i])
			}
		}
		result := []int{}
		for i := 0; i < len(testCase.output); i++ {
			result = append(result, *stack.Pop())
		}
		if !reflect.DeepEqual(result, testCase.output) {
			t.Errorf("uotput stack = %d, want %d", result, testCase.output)
		}
	}
	if stack.Pop() != nil {
		t.Errorf("stack.Pop() = %v, want nil", stack.Pop())
	}
	if stack.Top() != nil {
		t.Errorf("stack.Top() = %v, want nil", stack.Top())
	}
}
