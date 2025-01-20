package virtm

import "fmt"

type Scope map[string]StackValue

func (scope *Scope) get(name string) StackValue {
	if value, ok := (*scope)[name]; ok {
		return value
	}

	panic(fmt.Sprintf("No variable in scope named: %s", name))
}

func (scope *Scope) set(name string, value StackValue) {
	(*scope)[name] = value
}

type Variables []Scope

func (variables Variables) Get(name string) StackValue {
	for index, _ := range variables {
		scope := variables[len(variables)-1-index]

		if value, ok := scope[name]; ok {
			return value
		}
	}

	panic(fmt.Sprintf("Variable %s is not defined", name))
}

func (variables *Variables) Set(name string, value StackValue) {
	scope := (*variables)[len(*variables)-1]

	scope.set(name, value)
}

func (variables *Variables) PopScope() {
	*variables = (*variables)[:len(*variables)-1]
}

func (variables *Variables) NewScope() {
	*variables = append(*variables, make(Scope))
}

func CreateVariables() *Variables {
	vars := &Variables{}

	vars.NewScope()

	return vars
}
