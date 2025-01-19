package virtm

import (
	"fmt"
	"strings"
)

type VirtualMachine struct {
	stack          []interface{}          // Стек для вычислений
	bytecode       []string               // Байткод
	programCounter int                    // Указатель текущей инструкции
	variables      map[string]interface{} // Переменные
	labels         map[string]int         // Метки переходов
}

func NewVirtualMachine(bytecode []string) *VirtualMachine {
	return &VirtualMachine{
		stack:          make([]interface{}, 0),
		bytecode:       bytecode,
		programCounter: 0,
		variables:      make(map[string]interface{}),
		labels:         make(map[string]int),
	}
}

func (virtualMachine *VirtualMachine) Run() {
	virtualMachine.prepareLabels()
	virtualMachine.cleanBytecode()

	for virtualMachine.programCounter < len(virtualMachine.bytecode) {

		command := virtualMachine.bytecode[virtualMachine.programCounter]
		virtualMachine.programCounter++

		virtualMachine.execute(command)
	}
}

// Выполнение команды
func (virtualMachine *VirtualMachine) execute(command string) {
	instructions := strings.Fields(command)

	instruction := instructions[0]

	var argument string
	if len(instructions) > 1 {
		argument = instructions[1]
	}

	switch instruction {
	case "PUSH_CONST":
		constant := virtualMachine.parseArgument(argument)
		virtualMachine.stack = append(virtualMachine.stack, constant)

	case "PUSH_VAR":
		value, exists := virtualMachine.variables[argument]
		if !exists {
			panic(fmt.Sprintf("Variable %s not defined", argument))
		}
		virtualMachine.stack = append(virtualMachine.stack, value)

	case "STORE_VAR":
		value := virtualMachine.pop()
		virtualMachine.variables[argument] = value

	case "ADD":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int)+b.(int))

	case "SUB":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int)-b.(int))

	case "MUL":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int)*b.(int))

	case "DIV":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		if b.(float64) == 0 {
			panic("Division by zero")
		}
		virtualMachine.stack = append(virtualMachine.stack, a.(int)/b.(int))

	case "NEG":
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, -a.(int))

	case "NOT":
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, !a.(bool))

	case "AND":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(bool) && b.(bool))

	case "OR":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(bool) || b.(bool))

	case "LESS_THAN":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int) < b.(int))

	case "GREATER_THAN":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int) > b.(int))

	case "LESS_EQUAL_THAN":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int) <= b.(int))

	case "GREATER_EQUAL_THAN":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a.(int) >= b.(int))

	case "EQUAL":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a == b)

	case "NOT_EQUAL":
		b := virtualMachine.pop()
		a := virtualMachine.pop()
		virtualMachine.stack = append(virtualMachine.stack, a != b)

	case "JUMP":
		label := argument
		if index, exists := virtualMachine.labels[label]; exists {
			virtualMachine.programCounter = index
		} else {
			panic(fmt.Sprintf("Label not found: %s", label))
		}

	case "JUMP_IF_FALSE":
		condition := virtualMachine.pop()
		label := argument
		if !condition.(bool) {
			if index, exists := virtualMachine.labels[label]; exists {
				virtualMachine.programCounter = index
			} else {
				panic(fmt.Sprintf("Label not found: %s", label))
			}
		}

	case "NEW_ARRAY":
		size := virtualMachine.parseArgument(argument)
		newArray := make([]interface{}, size)
		virtualMachine.stack = append(virtualMachine.stack, newArray)

	case "ARRAY_GET":
		index := virtualMachine.pop().(int)
		array := virtualMachine.pop().([]interface{})
		virtualMachine.stack = append(virtualMachine.stack, array[index])

	case "ARRAY_SET":
		value := virtualMachine.pop()
		index := virtualMachine.pop().(int)
		array := virtualMachine.pop().([]interface{})
		array[index] = value

	case "LABEL":

	case "FALSE_LABEL", "LOOP_START_LABEL", "LOOP_END_LABEL", "END_LABEL":

	case "SCOPE_START", "SCOPE_END":
		// для GC

	case "CALL_FUNCTION":
		fmt.Println("CALL_FUNCTION not implemented")

	case "RETURN":
		fmt.Println("RETURN not implemented")

	case "PRINT":
		value := virtualMachine.pop()
		fmt.Println(value)

	default:
		fmt.Printf("Unknown command: %s\n", instruction)
	}
}

func (virtualMachine *VirtualMachine) parseArgument(arg string) int {
	var index int
	_, err := fmt.Sscanf(arg, "%d", &index)
	if err != nil {
		panic(fmt.Sprintf("Invalid argument: %s", arg))
	}
	return index
}

func (virtualMachine *VirtualMachine) pop() interface{} {
	if len(virtualMachine.stack) == 0 {
		panic("Stack underflow")
	}
	value := virtualMachine.stack[len(virtualMachine.stack)-1]
	virtualMachine.stack = virtualMachine.stack[:len(virtualMachine.stack)-1]
	return value
}

func (virtualMachine *VirtualMachine) prepareLabels() {
	for i, command := range virtualMachine.bytecode {
		if strings.HasPrefix(command, "LABEL") {
			parts := strings.Fields(command)
			if len(parts) > 1 {
				virtualMachine.labels[parts[1]] = i
			}
		}
	}
}

func (virtualMachine *VirtualMachine) cleanBytecode() {
	newBytecode := []string{}
	for _, command := range virtualMachine.bytecode {
		if !strings.HasPrefix(command, "LABEL") {
			newBytecode = append(newBytecode, command)
		}
	}
	virtualMachine.bytecode = newBytecode
}

func (virtualMachine *VirtualMachine) PrintBytecode() {
	fmt.Printf("\n\n")
	fmt.Printf("PRINT BYTECODE FUNCTION:\n\n")

	for _, bc := range virtualMachine.bytecode {
		fmt.Printf("%s", bc)
	}
}
