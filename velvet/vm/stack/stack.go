package stack

import (
	"fmt"
	"os"
	"strings"
)

type Stack []StackValue

func New() Stack {
	return Stack([]StackValue{})
}

func (s Stack) Dump() string {
	dumped := []string{}

	for _, item := range s {
		dumped = append(dumped, item.Dump())
	}

	return "[\n" + strings.Join(dumped, "\n") + "\n]"
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
	if len(kinds) > len(s) {
		return fmt.Errorf("expected '%s' on the stack, but the stack is not large enough", kinds[len(kinds)-1].Name())
	}

	ki, si := 0, len(s)-1

	for ki < len(kinds) {
		if !s[si].Is(kinds[ki]) && kinds[ki] == Any {
			return fmt.Errorf("expected '%s' on the stack, but found '%s' instead", kinds[ki].Name(), s[si].GetKind().Name())
		}
		ki += 1
		si -= 1
	}
	return nil
}

func (s Stack) Expect(kinds ...ValueKind) {
	if err := s.ExpectErr(kinds...); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
