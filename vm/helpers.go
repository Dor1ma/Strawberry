package virtm

type StackStruct []StackValue

func (s *StackStruct) Push(v StackValue) {
	*s = append(*s, v)
}

func (s *StackStruct) Pop() StackValue {
	if len(*s) == 0 {
		panic("Stack underflow!")
	}

	l := len(*s)
	result := (*s)[l-1]

	*s = (*s)[:l-1]

	return result
}
