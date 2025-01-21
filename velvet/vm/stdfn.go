package vm

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvet/vm/stack"
)

var stdfn = map[string]func(st *stack.Stack) error{
	"error": func(st *stack.Stack) error {
		return errors.New("")
	},
	"reset": func(st *stack.Stack) error {
		return nil
	},

	// operator functions
	"eq": func(st *stack.Stack) error {
		st.Expect(stack.Any, stack.Any)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.Equals(y)))
		return nil
	},
	"neq": func(st *stack.Stack) error {
		st.Expect(stack.Any, stack.Any)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(!x.Equals(y)))
		return nil
	},
	"not": func(st *stack.Stack) error {
		st.Expect(stack.Bool)
		st.Push(stack.NewBoolValue(!st.Pop().GetBool()))
		return nil
	},
	"lt": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() < y.GetNum()))
		return nil
	},
	"gt": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() > y.GetNum()))
		return nil
	},
	"lte": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() <= y.GetNum()))
		return nil
	},
	"gte": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewBoolValue(x.GetNum() >= y.GetNum()))
		return nil
	},
	"add": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() + y.GetNum()))
		return nil
	},
	"sub": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() - y.GetNum()))
		return nil
	},
	"mul": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() * y.GetNum()))
		return nil
	},
	"div": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(x.GetNum() / y.GetNum()))
		return nil
	},
	"pow": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(math.Pow(float64(x.GetNum()), float64(y.GetNum())))))
		return nil
	},
	"log": func(st *stack.Stack) error {
		st.Expect(stack.Number)
		st.Push(stack.NewNumberValue(float32(math.Log(float64(st.Pop().GetNum())))))
		return nil
	},
	"and": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) & int(y.GetNum()))))
		return nil
	},
	"or": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) | int(y.GetNum()))))
		return nil
	},
	"xor": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Number)
		y, x := st.Pop(), st.Pop()
		st.Push(stack.NewNumberValue(float32(int(x.GetNum()) ^ int(y.GetNum()))))
		return nil
	},
	// end operator functions

	// IO functions
	"print": func(st *stack.Stack) error {
		st.Expect(stack.Any)
		fmt.Print(st.Pop().Format())
		return nil
	},
	"println": func(st *stack.Stack) error {
		st.Expect(stack.Any)
		fmt.Printf("%v\n", st.Pop().Format())
		return nil
	},
	"putc": func(st *stack.Stack) error {
		st.Expect(stack.Number)
		fmt.Print(string(rune(int(st.Pop().GetNum()))))
		return nil
	},
	"putcln": func(st *stack.Stack) error {
		st.Expect(stack.Number)
		fmt.Println(string(rune(int(st.Pop().GetNum()))))
		return nil
	},
	"readn": func(st *stack.Stack) error {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}

		if num, err := strconv.ParseFloat(scanner.Text(), 32); err != nil {
			return err
		} else {
			st.Push(stack.NewNumberValue(float32(num)))
		}

		return nil
	},
	"readt": func(st *stack.Stack) error {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}

		st.Push(stack.NewStringValue(scanner.Text()))

		return nil
	},
	"readb": func(st *stack.Stack) error {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}

		b := []stack.StackValue{}
		for _, byte := range scanner.Bytes() {
			b = append(b, stack.NewNumberValue(float32(byte)))
		}

		st.Push(stack.NewListValue(b...))

		return nil
	},
	"readc": func(st *stack.Stack) error {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}

		st.Push(stack.NewNumberValue(float32(scanner.Bytes()[0])))

		return nil
	},
	// end IO functions

	// string operations
	"strip": func(st *stack.Stack) error {
		st.Expect(stack.String)
		st.Push(stack.NewStringValue(strings.TrimSpace(st.Pop().GetString())))
		return nil
	},
	"split": func(st *stack.Stack) error {
		st.Expect(stack.String, stack.String)

		y, x := st.Pop().GetString(), st.Pop().GetString()
		l := []stack.StackValue{}
		for _, v := range strings.Split(x, y) {
			l = append(l, stack.NewStringValue(v))
		}

		st.Push(stack.NewListValue(l...))
		return nil
	},
	// end string operations

	// seqence operations
	"allocList": func(st *stack.Stack) error {
		st.Expect(stack.Number)
		st.Push(stack.AllocListValue(int(st.Pop().GetNum())))
		return nil
	},
	"allocInitList": func(st *stack.Stack) error {
		st.Expect(stack.Number, stack.Any)
		y, x := st.Pop(), int(st.Pop().GetNum())
		st.Push(stack.AllocInitListValue(x, y))
		return nil
	},
	"len": func(st *stack.Stack) error {
		st.Expect(stack.List | stack.String)

		if seq := st.Pop(); seq.Is(stack.String) {
			st.Push(stack.NewNumberValue(float32(len(seq.GetString()))))
		} else {
			st.Push(stack.NewNumberValue(float32(len(seq.GetList()))))
		}

		return nil
	},
	"index": func(st *stack.Stack) error {
		st.Expect(stack.List|stack.String, stack.Number)

		if i, seq := st.Pop(), st.Pop(); seq.Is(stack.String) {
			st.Push(stack.NewNumberValue(float32(seq.GetString()[int(i.GetNum())])))
		} else {
			st.Push(seq.GetList()[int(i.GetNum())])
		}

		return nil
	},
	// end seqence operations
}
