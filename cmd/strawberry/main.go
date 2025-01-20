package main

import (
	"fmt"
	bytecodegen "github.com/Dor1ma/Strawberry/bytecode"
	"github.com/Dor1ma/Strawberry/cmd/strawberry/repl"
	"github.com/Dor1ma/Strawberry/lexer"
	"github.com/Dor1ma/Strawberry/parser"
	virtm "github.com/Dor1ma/Strawberry/vm"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) >= 2 {
		name := ".\\tasks\\001-bubble-sort.berry"
		//name := ".\\tasks\\002-factorial.berry"
		//name := ".\\tasks\\003-era.berry"

		b, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		l := lexer.New(string(b))
		p := parser.New(l)
		if statements, err := p.Parse(); err == nil && len(statements) != 0 {
			generator := bytecodegen.CodeGenerator{}

			generator.EnableLoopEnrolling()

			generator.GenerateProgram(statements)

			generator.EliminateDeadCode()

			vm := virtm.NewVirtualMachine(generator.GetBytecodes())

			vm.EnableTailRecursionOptimization()

			vm.Run()
		}
		return
	}

	fmt.Fprintln(os.Stdout, "Strawberry.")
	fmt.Fprintln(os.Stdout, "Type \"exit\" to exit.")
	repl.Start(os.Stdin, os.Stdout)
}
