package vm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvet/vm/stack"
)

var stdfn = map[string]func(st *stack.Stack) bool{
	"error": func(st *stack.Stack) bool {
		return false
	},
	"reset": func(st *stack.Stack) bool {
		return false
	},

	// operator functions
	"eq": func(st *stack.Stack) bool {
		st.Expect(stack.Any, stack.Any)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.Equals(y)))
		return false
	},
	"neq": func(st *stack.Stack) bool {
		st.Expect(stack.Any, stack.Any)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(!x.Equals(y)))
		return false
	},
	"not": func(st *stack.Stack) bool {
		st.Expect(stack.Bool)
		st.Push(stack.NewBoolValue(!st.Pop().GetBool()))
		return false
	},
	"lt": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() < y.GetNum()))
		return false
	},
	"gt": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() > y.GetNum()))
		return false
	},
	"lte": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() <= y.GetNum()))
		return false
	},
	"gte": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() >= y.GetNum()))
		return false
	},
	"add": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() + y.GetNum()))
		return false
	},
	"sub": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() - y.GetNum()))
		return false
	},
	"mul": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() * y.GetNum()))
		return false
	},
	"div": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() / y.GetNum()))
		return false
	},
	"pow": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(math.Pow(float64(x.GetNum()), float64(y.GetNum())))))
		return false
	},
	"log": func(st *stack.Stack) bool {
		st.Expect(stack.Number)
		st.Push(stack.NewNumberValue(float32(math.Log(float64(st.Pop().GetNum())))))
		return false
	},
	"and": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) & int(y.GetNum()))))
		return false
	},
	"or": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) | int(y.GetNum()))))
		return false
	},
	"xor": func(st *stack.Stack) bool {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) ^ int(y.GetNum()))))
		return false
	},
	// end operator functions

	// IO functions
	"print": func(st *stack.Stack) bool {
		st.Expect(stack.Any)
		fmt.Print(st.Pop().Format())
		return false
	},
	"println": func(st *stack.Stack) bool {
		st.Expect(stack.Any)
		fmt.Printf("%v\n", st.Pop().Format())
		return false
	},
	"putc": func(st *stack.Stack) bool {
		st.Expect(stack.Number)

		fmt.Print(string(rune(int(st.Pop().GetNum()))))

		return false
	},
	"putcln": func(st *stack.Stack) bool {
		st.Expect(stack.Number)

		fmt.Println(string(rune(int(st.Pop().GetNum()))))

		return false
	},
	"readNumber": func(st *stack.Stack) bool {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return true
		}

		if num, err := strconv.ParseFloat(scanner.Text(), 32); err != nil {
			return true
		} else {
			st.Push(stack.NewNumberValue(float32(num)))
		}

		return false
	},

	// end IO functions

	// string operations
	"strip": func(st *stack.Stack) bool {
		st.Expect(stack.String)
		st.Push(stack.NewStringValue(strings.TrimSpace(st.Pop().GetString())))
		return false
	},
	"split": func(st *stack.Stack) bool {
		st.Expect(stack.String, stack.String)

		y, x := st.Pop().GetString(), st.Pop().GetString()
		l := []stack.StackValue{}
		for _, v := range strings.Split(x, y) {
			l = append(l, stack.NewStringValue(v))
		}

		st.Push(stack.NewListValue(l...))
		return false
	},
	// end string operations

	// seqence operations
	"allocList": func(st *stack.Stack) bool {
		st.Expect(stack.Number)
		st.Push(stack.AllocListValue(int(st.Pop().GetNum())))
		return false
	},
	"len": func(st *stack.Stack) bool {
		st.Expect(stack.List | stack.String)

		if seq := st.Pop(); seq.Is(stack.String) {
			st.Push(stack.NewNumberValue(float32(len(seq.GetString()))))
		} else {
			st.Push(stack.NewNumberValue(float32(len(seq.GetList()))))
		}

		return false
	},
	"index": func(st *stack.Stack) bool {
		st.Expect(stack.List|stack.String, stack.Number)

		if i, seq := st.Pop(), st.Pop(); seq.Is(stack.String) {
			st.Push(stack.NewNumberValue(float32(seq.GetString()[int(i.GetNum())])))
		} else {
			st.Push(seq.GetList()[int(i.GetNum())])
		}

		return false
	},
	// end seqence operations
}
