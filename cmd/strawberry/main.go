package main

import (
	"fmt"
	"github.com/Dor1ma/Strawberry/cmd/strawberry/repl"
	"github.com/Dor1ma/Strawberry/interpreter"
	"github.com/Dor1ma/Strawberry/lexer"
	"github.com/Dor1ma/Strawberry/parser"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) >= 2 {
		name := os.Args[1]
		b, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		l := lexer.New(string(b))
		p := parser.New(l)
		if statements, err := p.Parse(); err == nil && len(statements) != 0 {
			interpreter.Interpret(statements)

			/*for i := 0; i < len(statements); i++ {
				fmt.Println(statements[i])
			}*/
		}
		return
	}

	fmt.Fprintln(os.Stdout, "Lox programing language.")
	fmt.Fprintln(os.Stdout, "Feel free to type commands.")
	fmt.Fprintln(os.Stdout, "Type \"exit\" to exit.")
	repl.Start(os.Stdin, os.Stdout)
}
