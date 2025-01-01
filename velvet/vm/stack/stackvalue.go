package stack

type ValueKind uint8

const (
	Number ValueKind = iota
	String
	Bool
	List
	Any
)

func (vk ValueKind) Type() string {
	return []string{"Number", "String", "Bool", "List", "Any"}[vk]
}

type StackValue struct {
	numVal    float32
	stringVal string
	listVal   []StackValue
	boolVal   bool
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
	return StackValue{kind: String, listVal: values}
}

func (sv StackValue) Is(kind ValueKind) bool {
	return sv.kind == kind
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
	}
	panic("unreachable")
}

func (sv StackValue) Equals(other StackValue) bool {
	return sv.GetAny() == sv.GetAny()
}
