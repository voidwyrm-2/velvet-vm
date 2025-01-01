package stack

import (
	"fmt"
	"os"
)

type Stack []StackValue

func New() Stack {
	return Stack([]StackValue{})
}

func (s *Stack) Push(v StackValue) {
	*s = append(*s, v)
}

func (s Stack) Empty() bool {
	return len(s) == 0
}

func (s *Stack) Pop() StackValue {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *Stack) TryPop() (StackValue, bool) {
	if s.Empty() {
		return StackValue{}, false
	}
	return s.Pop(), true
}

func (s Stack) ExpectErr(k ...ValueKind) error {
	for i, v := range s {
		if i >= len(k) {
			return fmt.Errorf("expected '%s' on the stack, but the stack is not large enough", k[i].Type())
		} else if !v.Is(k[i]) && k[i] != Any {
			return fmt.Errorf("expected '%s' on the stack, but found '%s' instead", k[i].Type(), v.GetKind().Type())
		}
	}
	return nil
}

func (s Stack) Expect(k ...ValueKind) {
	if err := s.ExpectErr(k...); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
