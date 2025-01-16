package repl

import (
	"bufio"
	"fmt"
	"github.com/Dor1ma/Strawberry/interpreter"
	"github.com/Dor1ma/Strawberry/lexer"
	"github.com/Dor1ma/Strawberry/parser"
	"io"
	"strings"
)

const prompt = ">> "

// Start creates a REPL for Lox.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	interpreter.SetEvalEnv("repl")
	for {
		fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if line == "exit" {
			fmt.Fprintln(out, "bye.")
			return
		}
		l := lexer.New(line)
		p := parser.New(l)
		statements, err := p.Parse()
		if err == nil && len(statements) != 0 {
			interpreter.Interpret(statements)
		}
	}
}
