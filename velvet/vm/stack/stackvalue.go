package stack

import (
	"fmt"
	"strings"
)

type ValueKind int16

const (
	Any                = -1
	Number   ValueKind = 0b0
	String             = 0b1
	Bool               = 0b10
	List               = 0b100
	Function           = 0b1000
)

func (vk ValueKind) Name() string {
	return map[ValueKind]string{Any: "Any", Number: "Number", String: "String", Bool: "Bool", List: "List", Function: "Function"}[vk]
}

type StackValue struct {
	numVal    float32
	stringVal string
	listVal   []StackValue
	boolVal   bool
	funcVal   func(st *Stack) error
	kind      ValueKind
}

func NewNumberValue(value float32) StackValue {
	return StackValue{kind: Number, numVal: value}
}

func NewBoolValue(value bool) StackValue {
	return StackValue{kind: Bool, boolVal: value}
}

func NewStringValue(value string) StackValue {
	return StackValue{kind: String, stringVal: value}
}

func NewListValue(values ...StackValue) StackValue {
	return StackValue{kind: List, listVal: values}
}

func AllocListValue(size int) StackValue {
	return StackValue{kind: List, listVal: make([]StackValue, size)}
}

func AllocInitListValue(size int, value StackValue) StackValue {
	l := AllocListValue(size)
	for range size {
		l.listVal = append(l.listVal, value)
	}
	return l
}

func NewFuncValue(value func(st *Stack) error) StackValue {
	return StackValue{kind: Function, funcVal: value}
}

func (sv StackValue) Dump() string {
	return fmt.Sprintf("{%s, '%s', %f, %v}", sv.kind.Name(), sv.stringVal, sv.numVal, sv.boolVal)
}

func (sv StackValue) Is(kind ValueKind) bool {
	return sv.kind&kind == sv.kind || Any&kind == Any
}

func (sv StackValue) GetKind() ValueKind {
	return sv.kind
}

func (sv StackValue) GetNum() float32 {
	return sv.numVal
}

func (sv StackValue) GetString() string {
	return sv.stringVal
}

func (sv StackValue) GetBool() bool {
	return sv.boolVal
}

func (sv StackValue) GetList() []StackValue {
	return sv.listVal
}

func (sv StackValue) GetFunc() func(st *Stack) error {
	return sv.funcVal
}

func (sv StackValue) GetAny() any {
	switch sv.kind {
	case Number:
		return sv.GetNum()
	case String:
		return sv.GetString()
	case Bool:
		return sv.GetBool()
	case List:
		return sv.GetList()
	case Function:
		return sv.GetFunc()
	}
	panic("unreachable")
}

func (sv StackValue) Format() string {
	if sv.kind == Function {
		return "<Function>"
	} else if sv.kind == List {
		fi := []string{}

		for _, item := range sv.GetList() {
			if item.Is(String) {
				fi = append(fi, "\""+item.Format()+"\"")
			} else {
				fi = append(fi, item.Format())
			}
		}

		return fmt.Sprintf("[ %s ]", strings.Join(fi, " "))
	}
	return fmt.Sprintf("%v", sv.GetAny())
}

func (sv StackValue) Equals(other StackValue) bool {
	return sv.GetAny() == sv.GetAny()
}
