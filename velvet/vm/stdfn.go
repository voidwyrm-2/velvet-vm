package vm

import (
	"fmt"
	"math"

	"github.com/voidwyrm-2/velvet-vm/velvet/vm/stack"
)

var stdfn = map[string]func(st stack.Stack) bool{
	"error": func(st stack.Stack) bool {
		return false
	},
	"reset": func(st stack.Stack) bool {
		return false
	},

	"eq": func(st stack.Stack) bool { // operator functions
		st.Expect(stack.Any, stack.Any)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.Equals(y)))
		return false
	},
	"neq": func(st stack.Stack) bool {
		st.Expect(stack.Any, stack.Any)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(!x.Equals(y)))
		return false
	},
	"not": func(st stack.Stack) bool {
		st.Expect(stack.Bool)
		st.Push(stack.NewBoolValue(!st.Pop().GetBool()))
		return false
	},
	"lt": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() < y.GetNum()))
		return false
	},
	"gt": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() > y.GetNum()))
		return false
	},
	"lte": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() <= y.GetNum()))
		return false
	},
	"gte": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() >= y.GetNum()))
		return false
	},
	"add": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() + y.GetNum()))
		return false
	},
	"sub": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() - y.GetNum()))
		return false
	},
	"mul": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() * y.GetNum()))
		return false
	},
	"div": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() / y.GetNum()))
		return false
	},
	"pow": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(math.Pow(float64(x.GetNum()), float64(y.GetNum())))))
		return false
	},
	"log": func(st stack.Stack) bool {
		st.Expect(stack.Number)
		st.Push(stack.NewNumberValue(float32(math.Log(float64(st.Pop().GetNum())))))
		return false
	},
	"and": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) & int(y.GetNum()))))
		return false
	},
	"or": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) | int(y.GetNum()))))
		return false
	},
	"xor": func(st stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) ^ int(y.GetNum()))))
		return false
	}, // end operator functions

	"print": func(st stack.Stack) bool {
		st.Expect(stack.Any)
		fmt.Print(st.Pop().GetAny())
		return false
	},
	"println": func(st stack.Stack) bool {
		st.Expect(stack.Any)
		fmt.Printf("%v\n", st.Pop().GetAny())
		return false
	},
}
