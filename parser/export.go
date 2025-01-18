package parser

import (
	"github.com/Dor1ma/Strawberry/ast"
	"github.com/Dor1ma/Strawberry/lexer"
)

// вспомогательные функции

func ParseExpr(input string) (expr ast.Expression, err error) {
	l := lexer.New(input)
	p := New(l)
	defer func() {
		if r := recover(); r != nil {
			if parseErr, ok := r.(parseError); ok {
				err = &parseErr
				expr = nil
			} else {
				panic(r)
			}
		}
	}()
	return p.parseExpression(), nil
}

func ParseStmts(input string) ([]ast.Statement, error) {
	l := lexer.New(input)
	p := New(l)
	return p.Parse()
}
