package virtm

type GCObject struct {
	marked bool
	data   []StackValue
}

func (gc *VirtualMachine) MarkRoots() {
	for _, item := range gc.stack {
		gc.Mark(item)
	}
	for _, scope := range gc.variables {
		for _, value := range scope {
			gc.Mark(value)
		}
	}
}

func (gc *VirtualMachine) Mark(value StackValue) {
	if obj, ok := value.Value.(*GCObject); ok && !obj.marked {

		obj.marked = true
	}
}

func (gc *VirtualMachine) Sweep() {
	for key, obj := range gc.heap {
		if obj.marked {
			obj.marked = false
			delete(gc.heap, key)
		}
	}
}

func (gc *VirtualMachine) Collect() {
	gc.MarkRoots()
	gc.Sweep()
}
