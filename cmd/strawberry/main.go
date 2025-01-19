package main

import (
	"fmt"
	"github.com/Dor1ma/Strawberry/ast"
	bytecodegen "github.com/Dor1ma/Strawberry/bytecode"
	"github.com/Dor1ma/Strawberry/cmd/strawberry/repl"
	"github.com/Dor1ma/Strawberry/interpreter"
	"github.com/Dor1ma/Strawberry/lexer"
	"github.com/Dor1ma/Strawberry/parser"
	virtm "github.com/Dor1ma/Strawberry/vm"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) >= 2 {
		/*name := os.Args[1]*/

		/*name := ".\\tasks\\001-bubble-sort.berry"*/
		/*name := ".\\tasks\\002-factorial.berry"*/
		/*name := ".\\tasks\\003-era.berry"*/
		name := ".\\example\\0-arrays.berry"

		b, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		l := lexer.New(string(b))
		p := parser.New(l)
		if statements, err := p.Parse(); err == nil && len(statements) != 0 {

			interpreter.Interpret(statements)

			for i := 0; i < len(statements); i++ {
				fmt.Println(ast.PrettyPrint(statements[i], 1))
			}

			generator := bytecodegen.CodeGenerator{}

			generator.GenerateProgram(statements)
			generator.PrintBytecode()

			vm := virtm.NewVirtualMachine(generator.GetBytecodes())

			vm.Run()

			vm.PrintBytecode()
		}
		return
	}

	fmt.Fprintln(os.Stdout, "Strawberry.")
	fmt.Fprintln(os.Stdout, "Type \"exit\" to exit.")
	repl.Start(os.Stdin, os.Stdout)
}
