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

func (s Stack) ExpectErr(kinds ...ValueKind) error {
	for i, k := range kinds {
		if i >= len(s) {
			return fmt.Errorf("expected '%s' on the stack, but the stack is not large enough", k.Type())
		} else if !s[i].Is(k) && k != Any {
			return fmt.Errorf("expected '%s' on the stack, but found '%s' instead", k.Type(), s[i].GetKind().Type())
		}
	}
	return nil
}

func (s Stack) Expect(kinds ...ValueKind) {
	if err := s.ExpectErr(kinds...); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
