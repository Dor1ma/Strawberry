package virtm

import (
	"fmt"
	bytecode_gen "github.com/Dor1ma/Strawberry/bytecode"
	"strconv"
	"strings"
)

const (
	INT    = "int"
	BOOL   = "bool"
	STRING = "string"
	ARRAY  = "array"
)

var isTailOptimizationEnabled = false

type ValueType string

type StackValue struct {
	Value     interface{}
	ValueType ValueType
}

func (sv StackValue) String() string {
	switch sv.ValueType {
	case INT:
		return fmt.Sprintf("%d", sv.Value.(int))
	case BOOL:
		return fmt.Sprintf("%t", sv.Value.(bool))
	case STRING:
		return fmt.Sprintf("'%s'", sv.Value.(string))
	case ARRAY:
		return fmt.Sprintf("[%s]", sv.Value)
	default:
		return "UNKNOWN TYPE"
	}
}

type VirtualMachine struct {
	stack           StackStruct
	bytecode        []string
	programCounter  int
	variables       Variables
	labels          map[string]int
	heap            map[string]GCObject
	arrayCounter    int
	callStack       []StackStruct
	returnAddresses StackStruct
}

func (vm *VirtualMachine) newArrayID() string {
	vm.arrayCounter++
	return fmt.Sprintf("array_%d", vm.arrayCounter)
}

func NewVirtualMachine(bytecode []string) *VirtualMachine {
	return &VirtualMachine{
		stack:           make(StackStruct, 0),
		bytecode:        bytecode,
		programCounter:  0,
		variables:       *CreateVariables(),
		labels:          make(map[string]int),
		arrayCounter:    0,
		heap:            make(map[string]GCObject),
		returnAddresses: make(StackStruct, 0),
	}
}

func (virtualMachine *VirtualMachine) Run() {
	virtualMachine.prepareLabels()
	/*virtualMachine.cleanBytecode()*/

	for virtualMachine.programCounter < len(virtualMachine.bytecode) {

		command := virtualMachine.bytecode[virtualMachine.programCounter]
		virtualMachine.programCounter++

		virtualMachine.execute(command)
	}
}

func (virtualMachine *VirtualMachine) execute(command string) {
	instructions := strings.Fields(command)

	instruction := instructions[0]

	var nonParsedArgument string
	if len(instructions) > 1 {
		nonParsedArgument = strings.Join(instructions[1:], " ")
	}

	argument, argType := virtualMachine.parseArgument(nonParsedArgument)

	value := StackValue{Value: argument, ValueType: argType}

	switch instruction {
	case bytecode_gen.PUSH_CONST:
		virtualMachine.stack.Push(value)

	case bytecode_gen.PUSH_VAR:
		varValue := virtualMachine.variables.Get(nonParsedArgument)
		virtualMachine.stack.Push(varValue)

	case bytecode_gen.STORE_VAR:
		poppedValue := virtualMachine.stack.Pop()
		virtualMachine.variables.Set(nonParsedArgument, poppedValue)

	case bytecode_gen.ADD:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		result := StackValue{}
		switch a.ValueType {

		case INT:
			result.Value = a.Value.(int) + b.Value.(int)
			result.ValueType = INT
		case STRING:
			result.Value = a.Value.(string) + b.Value.(string)
			result.ValueType = STRING
		case ARRAY:
			arr := virtualMachine.heap[a.Value.(string)]

			arr.data = append(arr.data, b)

			virtualMachine.heap[a.Value.(string)] = arr
			result.Value = a.Value.(string)
			result.ValueType = ARRAY
		default:
			panic("unsupported operation ADD for this type")
		}

		virtualMachine.stack.Push(result)

	case bytecode_gen.SUB:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		result := StackValue{}
		switch a.ValueType {

		case INT:
			result.Value = a.Value.(int) - b.Value.(int)
			result.ValueType = INT
		default:
			panic("unsupported operation SUB for this type")
		}

		virtualMachine.stack.Push(result)

	case bytecode_gen.MUL:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		result := StackValue{}
		switch a.ValueType {

		case INT:
			result.Value = a.Value.(int) * b.Value.(int)
			result.ValueType = INT
		default:
			panic("unsupported operation MUL for this type")
		}

		virtualMachine.stack.Push(result)

	case bytecode_gen.DIV:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		result := StackValue{}
		if a.ValueType == INT && b.ValueType == INT {
			if b.Value.(int) == 0 {
				panic("Division by zero")
			}
			result.Value = a.Value.(int) / b.Value.(int)
			result.ValueType = INT
		} else {
			panic("unsupported operation DIV for these types")
		}

		virtualMachine.stack.Push(result)

	case bytecode_gen.NEG:
		a := virtualMachine.stack.Pop()

		if a.ValueType != INT {
			panic("unsupported operation NEG for non-integer type")
		}

		result := StackValue{Value: -a.Value.(int), ValueType: INT}
		virtualMachine.stack.Push(result)

	case bytecode_gen.NOT:
		a := virtualMachine.stack.Pop()

		if a.ValueType != BOOL {
			panic("unsupported operation NOT for non-boolean type")
		}

		result := StackValue{Value: !a.Value.(bool), ValueType: BOOL}
		virtualMachine.stack.Push(result)

	case bytecode_gen.AND:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == BOOL && b.ValueType == BOOL {
			result := StackValue{Value: a.Value.(bool) && b.Value.(bool), ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation AND for non-boolean types")
		}

	case bytecode_gen.OR:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == BOOL && b.ValueType == BOOL {
			result := StackValue{Value: a.Value.(bool) || b.Value.(bool), ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation OR for non-boolean types")
		}

	case bytecode_gen.LESS_THAN:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == INT && b.ValueType == INT {
			result := StackValue{Value: a.Value.(int) < b.Value.(int), ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation LESS_THAN for these types")
		}

	case bytecode_gen.GREATER_THAN:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == INT && b.ValueType == INT {
			result := StackValue{Value: a.Value.(int) > b.Value.(int), ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation GREATER_THAN for these types")
		}

	case bytecode_gen.LESS_EQUAL_THAN:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == INT && b.ValueType == INT {
			result := StackValue{Value: a.Value.(int) <= b.Value.(int), ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation LESS_EQUAL_THAN for these types")
		}

	case bytecode_gen.GREATER_EQUAL_THAN:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == INT && b.ValueType == INT {
			result := StackValue{Value: a.Value.(int) >= b.Value.(int), ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation GREATER_EQUAL_THAN for these types")
		}

	case bytecode_gen.EQUAL:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == b.ValueType {
			result := StackValue{Value: a.Value == b.Value, ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation EQUAL for mismatched types")
		}

	case bytecode_gen.NOT_EQUAL:
		b := virtualMachine.stack.Pop()
		a := virtualMachine.stack.Pop()

		if a.ValueType == b.ValueType {
			result := StackValue{Value: a.Value != b.Value, ValueType: BOOL}
			virtualMachine.stack.Push(result)
		} else {
			panic("unsupported operation NOT_EQUAL for mismatched types")
		}

	case bytecode_gen.JUMP:
		label := nonParsedArgument
		if index, exists := virtualMachine.labels[nonParsedArgument]; exists {
			virtualMachine.programCounter = index
		} else {
			panic(fmt.Sprintf("Label not found: %s", label))
		}

	case bytecode_gen.JUMP_IF_FALSE:
		condition := virtualMachine.stack.Pop()
		label := nonParsedArgument

		if condition.ValueType != BOOL {
			panic("JUMP_IF_FALSE requires a boolean condition")
		}

		if !condition.Value.(bool) {
			if index, exists := virtualMachine.labels[label]; exists {
				virtualMachine.programCounter = index
			} else {
				panic(fmt.Sprintf("Label not found: %s", label))
			}
		}

	case bytecode_gen.NEW_ARRAY:
		size, err := strconv.Atoi(nonParsedArgument)
		if err != nil || size < 0 {
			panic("Invalid size for NEW_ARRAY")
		}

		newArray := make([]StackValue, size)

		obj := GCObject{
			data:   newArray,
			marked: false,
		}

		for i := 0; i < size; i++ {
			if len(virtualMachine.stack) == 0 {
				panic("Stack underflow while initializing array")
			}
			newArray[size-i-1] = virtualMachine.stack.Pop() // Инициализация с вершины стека
		}

		arrayID := virtualMachine.newArrayID()
		virtualMachine.heap[arrayID] = obj
		virtualMachine.stack.Push(StackValue{Value: arrayID, ValueType: ARRAY})

	case bytecode_gen.ARRAY_GET:
		index := virtualMachine.stack.Pop()
		arrayRef := virtualMachine.stack.Pop()

		if index.ValueType != INT {
			panic("ARRAY_GET requires an integer index")
		}
		if arrayRef.ValueType != ARRAY {
			panic("ARRAY_GET requires an array reference")
		}

		arrayID := arrayRef.Value.(string)
		arr, exists := virtualMachine.heap[arrayID]
		if !exists {
			panic("Array not found")
		}

		idx := index.Value.(int)
		if idx < 0 || idx >= len(arr.data) {
			panic("Index out of bounds for ARRAY_GET")
		}

		virtualMachine.stack.Push(arr.data[idx])

	case bytecode_gen.ARRAY_SET:
		index := virtualMachine.stack.Pop()
		arrayRef := virtualMachine.stack.Pop()
		pop := virtualMachine.stack.Pop()

		if index.ValueType != INT {
			panic("ARRAY_SET requires an integer index")
		}
		if arrayRef.ValueType != ARRAY {
			panic("ARRAY_SET requires an array reference")
		}

		arrayID := arrayRef.Value.(string)
		arr, exists := virtualMachine.heap[arrayID]
		if !exists {
			panic("Array not found")
		}

		idx := index.Value.(int)
		if idx < 0 || idx >= len(arr.data) {
			panic("Index out of bounds for ARRAY_SET")
		}

		arr.data[idx] = pop
		virtualMachine.heap[arrayID] = arr

	case bytecode_gen.LABEL:

	case bytecode_gen.FALSE_LABEL, bytecode_gen.LOOP_START_LABEL, bytecode_gen.LOOP_END_LABEL, bytecode_gen.SCOPE_START:

	case bytecode_gen.SCOPE_END:
		virtualMachine.Collect()

	case bytecode_gen.CALL_FUNCTION:
		argumentCount := virtualMachine.stack.Pop()

		newStack := StackStruct{}
		for i := 0; i < argumentCount.Value.(int); i++ {
			newStack.Push(virtualMachine.stack.Pop())
		}

		if isTailOptimizationEnabled {
			nextInstruction := virtualMachine.bytecode[virtualMachine.programCounter]
			if strings.TrimSpace(nextInstruction) == bytecode_gen.RETURN {
				virtualMachine.stack = newStack
				virtualMachine.programCounter = virtualMachine.labels[nonParsedArgument]
				virtualMachine.variables.NewScope()
				return
			}
		}

		savedStack := virtualMachine.stack
		virtualMachine.callStack = append(virtualMachine.callStack, savedStack)
		virtualMachine.stack = newStack
		virtualMachine.variables.NewScope()
		virtualMachine.returnAddresses.Push(StackValue{virtualMachine.programCounter, INT})

		virtualMachine.programCounter = virtualMachine.labels[nonParsedArgument]

	case bytecode_gen.RETURN:
		if len(virtualMachine.callStack) == 0 {
			return
		}

		savedStack := virtualMachine.callStack[len(virtualMachine.callStack)-1]
		virtualMachine.callStack = virtualMachine.callStack[:len(virtualMachine.callStack)-1]

		returnedValue := virtualMachine.stack.Pop()
		virtualMachine.stack = savedStack
		virtualMachine.stack.Push(returnedValue)

		virtualMachine.variables.PopScope()
		returnAddress := virtualMachine.returnAddresses.Pop()
		virtualMachine.programCounter = returnAddress.Value.(int)

	case bytecode_gen.PRINT:
		pop := virtualMachine.stack.Pop()

		if pop.ValueType == ARRAY {
			fmt.Println(virtualMachine.heap[pop.Value.(string)].data)
		} else {
			fmt.Println(pop)
		}

	case bytecode_gen.FUNC:
		if len(instructions) < 2 {
			panic("FUNC requires a function name")
		}

		functionName := instructions[1]

		virtualMachine.labels[functionName] = virtualMachine.programCounter

		for virtualMachine.programCounter < len(virtualMachine.bytecode) {
			currentInstructions := strings.Fields(virtualMachine.bytecode[virtualMachine.programCounter])

			if currentInstructions[0] == bytecode_gen.END_FUNC {
				break
			}
			virtualMachine.programCounter++
		}

		if virtualMachine.programCounter >= len(virtualMachine.bytecode) {
			panic(fmt.Sprintf("END_FUNC not found for function %s", functionName))
		}

	case bytecode_gen.END_FUNC:
		if len(instructions) < 2 {
			panic("FUNC_END requires a function name")
		}

		functionName := instructions[1]
		if _, exists := virtualMachine.labels[functionName]; !exists {
			panic(fmt.Sprintf("END_FUNC found for unknown function %s", functionName))
		}

	default:
		fmt.Printf("Unknown command: %s\n", instruction)
	}
}

func (virtualMachine *VirtualMachine) prepareLabels() {
	for i, command := range virtualMachine.bytecode {
		if strings.HasPrefix(command, bytecode_gen.LABEL) {
			parts := strings.Fields(command)
			if len(parts) > 1 {
				virtualMachine.labels[parts[1]] = i
			}
		}
	}
}

/*func (virtualMachine *VirtualMachine) cleanBytecode() {
	newBytecode := []string{}
	for _, command := range virtualMachine.bytecode {
		if !strings.HasPrefix(command, bytecode_gen.LABEL) {
			newBytecode = append(newBytecode, command)
		}
	}
	virtualMachine.bytecode = newBytecode
}*/

func (virtualMachine *VirtualMachine) PrintBytecode() {
	fmt.Printf("\n\n")
	fmt.Printf("PRINT BYTECODE FUNCTION:\n\n")

	for _, bc := range virtualMachine.bytecode {
		fmt.Printf("%s", bc)
	}
}

func (virtualMachine *VirtualMachine) parseArgument(arg string) (interface{}, ValueType) {
	if intValue, err := strconv.Atoi(arg); err == nil {
		return intValue, INT
	}

	if arg == "true" {
		return true, BOOL
	}
	if arg == "false" {
		return false, BOOL
	}

	return arg, STRING
}

func (virtualMachine *VirtualMachine) EnableTailRecursionOptimization() {
	isTailOptimizationEnabled = true
}
