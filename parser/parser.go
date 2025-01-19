package parser

import (
	"fmt"
	"github.com/Dor1ma/Strawberry/ast"
	"github.com/Dor1ma/Strawberry/lexer"
	"github.com/Dor1ma/Strawberry/token"
	"os"
	"strings"
)

type Parser struct {
	l *lexer.Lexer

	tok token.Token
	lit string

	trace  bool
	indent int
}

func (p *Parser) nextToken() token.Token {
	if p.isAtEnd() {
		return token.EOF
	}
	tok, lit := p.l.NextToken()
	p.tok = tok
	p.lit = lit
	return tok
}

// Parse возвращает все операторы
func (p *Parser) Parse() (statements []ast.Statement, err error) {
	defer func() {
		if r := recover(); r != nil {
			if parseErr, ok := r.(parseError); ok {
				statements = nil
				err = &parseErr
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()
	// Здесь происходит парсинг всей программы
	for !p.isAtEnd() {
		stmt := p.parseDeclaration()
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *Parser) parseDeclaration() ast.Statement {
	if p.match(token.Var) {
		return p.parseVarDeclaration()
	}
	if p.match(token.Fun) {
		return p.parseFunctionDeclaration()
	}
	if p.match(token.Class) {
		return p.parseClassDeclaration()
	}
	return p.parseStatement()
}

func (p *Parser) parseVarDeclaration() *ast.VarStmt {
	name := p.lit
	p.expect(token.Identifier, "Expect variable name.")
	var stmt = &ast.VarStmt{
		Name: &ast.Identifier{
			Name: name,
		},
	}
	var initializer ast.Expression
	if p.match(token.Equal) {
		initializer = p.parseExpression()
	}
	p.expect(token.Semicolon, "Expect ';' after variable declaration.")
	stmt.Initializer = initializer
	return stmt
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionStmt {
	name := p.lit
	p.expect(token.Identifier, "Expect function name.")
	p.expect(token.LeftParen, "Expect '(' after function name.")
	fun := &ast.FunctionStmt{
		Name:   name,
		Params: make([]*ast.Identifier, 0),
		Body:   make([]ast.Statement, 0),
	}
	if !p.match(token.RightParen) {
		for {
			lit := p.lit
			p.expect(token.Identifier, "Expect parameter name.")
			if len(fun.Params) >= 255 {
				p.error("Cannot have more than 255 parameters.")
			}
			ident := &ast.Identifier{Name: lit}
			fun.Params = append(fun.Params, ident)
			if !p.match(token.Comma) {
				break
			}
		}
		p.expect(token.RightParen, "Expect ')' after parameters.")
	}
	p.expect(token.LeftBrace, "Expect '{' before function body.")
	fun.Body = p.parseBlockStatement().Statements
	return fun
}

func (p *Parser) parseClassDeclaration() *ast.ClassStmt {
	name := p.lit
	p.expect(token.Identifier, "Expect class name.")
	p.expect(token.LeftBrace, "Expect '{' after class name.")

	methods := make([]*ast.FunctionStmt, 0)
	for p.check(token.Identifier) {
		method := p.parseFunctionDeclaration()
		method.IsInitializer = method.Name == "init"
		methods = append(methods, method)
	}

	p.expect(token.RightBrace, "Expect '}' after class block.")

	return &ast.ClassStmt{
		Name:    name,
		Methods: methods,
	}
}

func (p *Parser) parseStatement() ast.Statement {
	if p.match(token.Print) {
		return p.parsePrintStatement()
	}
	if p.match(token.If) {
		return p.parseIfStatement()
	}
	if p.match(token.While) {
		return p.parseWhileStatement()
	}
	if p.match(token.For) {
		return p.parseForStatement()
	}
	if p.match(token.LeftBrace) {
		return p.parseBlockStatement()
	}
	if p.match(token.Return) {
		return p.parseReturnStatement()
	}
	return p.parseExprStatement()
}

func (p *Parser) parsePrintStatement() ast.Statement {
	expr := p.parseExpression()
	p.expect(token.Semicolon, "Expect ';' after value.")
	return &ast.PrintStmt{
		Expression: expr,
	}
}

func (p *Parser) parseIfStatement() ast.Statement {
	p.expect(token.LeftParen, "Expect '(' after 'if'.")
	condition := p.parseExpression()
	p.expect(token.RightParen, "Expect ')' after if condition.")
	thenBranch := p.parseStatement()
	var elseBranch ast.Statement
	if p.match(token.Else) {
		elseBranch = p.parseStatement()
	}
	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	p.expect(token.LeftParen, "Expect '(' after 'while'.")
	condition := p.parseExpression()
	p.expect(token.RightParen, "Expect ')' after while condition.")
	body := p.parseStatement()
	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseForStatement() ast.Statement {
	p.expect(token.LeftParen, "Expect '(' after 'for'.")
	var initializer ast.Statement
	if !p.match(token.Semicolon) {
		if p.match(token.Var) {
			initializer = p.parseVarDeclaration()
		} else {
			initializer = p.parseExprStatement()
		}
	}

	var condition ast.Expression
	if !p.match(token.Semicolon) {
		condition = p.parseExpression()
		p.expect(token.Semicolon, "Expect ';' after loop condition.")
	}

	var increment ast.Expression
	if !p.match(token.RightParen) {
		increment = p.parseExpression()
		p.expect(token.RightParen, "Expect ')' after for clause.")
	}

	body := p.parseStatement()

	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Statement{
				body,
				&ast.ExprStmt{
					Expression: increment,
				},
			},
		}
	}
	if condition == nil {
		condition = &ast.Literal{
			Token: token.True,
			Value: "true",
		}
	}
	body = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Statement{
				initializer,
				body,
			},
		}
	}
	return body
}

func (p *Parser) parseBlockStatement() *ast.BlockStmt {
	statements := make([]ast.Statement, 0)
	for !(p.check(token.RightBrace) || p.isAtEnd()) {
		statements = append(statements, p.parseDeclaration())
	}
	p.expect(token.RightBrace, "Expect '}' after block.")
	return &ast.BlockStmt{
		Statements: statements,
	}
}

func (p *Parser) parseExprStatement() ast.Statement {
	expr := p.parseExpression()
	p.expect(token.Semicolon, "Expect ';' after expression.")
	return &ast.ExprStmt{
		Expression: expr,
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStmt{}
	if !p.match(token.Semicolon) {
		stmt.Value = p.parseExpression()
		p.expect(token.Semicolon, "Expect ';' after return value.")
	}
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() ast.Expression {
	expr := p.parseOr()
	if p.match(token.Equal) {
		// ToDo: подумать над рекурсией здесь..
		v := p.parseAssignment()
		switch e := expr.(type) {
		default:
			p.error("Invalid assignment target.")
		case ast.LeftExpr:
			return &ast.AssignExpr{
				Left:  e,
				Value: v,
			}
		case *ast.GetExpr:
			return &ast.SetExpr{
				Object: e.Object,
				Name:   e.Name,
				Value:  v,
			}
		}
	}
	return expr
}

func (p *Parser) parseOr() ast.Expression {
	expr := p.parseAnd()
	if p.match(token.Or) {
		right := p.parseAnd()
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: token.Or,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) parseAnd() ast.Expression {
	expr := p.parseEquality()
	if p.match(token.And) {
		right := p.parseEquality()
		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: token.And,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) parseEquality() ast.Expression {
	expr := p.parseComparison()
	operator := p.tok
	for p.match(token.EqualEqual, token.NotEqual) {
		right := p.parseComparison()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseComparison() ast.Expression {
	expr := p.parseAddition()
	operator := p.tok
	for p.match(token.Greater, token.GreaterThanOrEqual, token.Less, token.LessThanOrEqual) {
		right := p.parseAddition()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseAddition() ast.Expression {
	expr := p.parseMultiplacation()
	operator := p.tok
	for p.match(token.Plus, token.Minus) {
		right := p.parseMultiplacation()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseMultiplacation() ast.Expression {
	expr := p.parseUnary()
	operator := p.tok
	for p.match(token.Slash, token.Star) {
		right := p.parseUnary()
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
		operator = p.tok
	}
	return expr
}

func (p *Parser) parseUnary() ast.Expression {
	operator := p.tok
	if p.match(token.Not, token.Minus) {
		right := p.parseUnary()
		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}
	return p.parseCall()
}

func (p *Parser) parseCall() ast.Expression {
	expr := p.parsePrimary()
	for {
		if p.match(token.LeftParen) {
			expr = p.finishCall(expr)
		} else if p.match(token.Dot) {
			name := p.lit
			p.expect(token.Identifier, "Expect property or method name after '.'.")
			expr = &ast.GetExpr{Object: expr, Name: name}
		} else if p.match(token.LeftBracket) {
			index := p.parseExpression()
			p.expect(token.RightBracket, "Expect ']' after array index.")
			expr = &ast.ArrayIndex{
				Array: expr,
				Index: index,
			}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) parseArrayExpr() *ast.ArrayExpr {
	var elements []ast.Expression

	if p.match(token.RightBracket) {
		return &ast.ArrayExpr{Elements: elements}
	}

	for {
		element := p.parseExpression()
		elements = append(elements, element)

		if p.match(token.Comma) {
			continue
		}

		p.expect(token.RightBracket, "Expect ']' after array elements.")
		break
	}

	return &ast.ArrayExpr{Elements: elements}
}

func (p *Parser) finishCall(expr ast.Expression) ast.Expression {
	call := &ast.CallExpr{
		Callee:    expr,
		Arguments: make([]ast.Expression, 0),
	}
	if p.match(token.RightParen) {
		return call
	}
	for {
		arg := p.parseExpression()
		// ToDo: ?
		if len(call.Arguments) >= 255 {
			p.error("Cannot have more than 255 arguments.")
		}
		call.Arguments = append(call.Arguments, arg)
		if !p.match(token.Comma) {
			break
		}
	}
	p.expect(token.RightParen, "Expect ')' after arguments.")
	return call
}

func (p *Parser) parsePrimary() (expr ast.Expression) {
	tok, lit := p.tok, p.lit
	switch tok {
	default:
		p.error("Expect expression.")
	case token.True, token.False, token.Nil, token.String, token.Number:
		expr = &ast.Literal{
			Token: tok,
			Value: lit,
		}
	case token.Identifier:
		expr = &ast.VariableExpr{
			Name:     lit,
			Distance: -1,
		}
	case token.This:
		expr = &ast.ThisExpr{}
	case token.LeftParen:
		p.nextToken()
		inner := p.parseExpression()
		p.expect(token.RightParen, "Expect ) after expression.")
		expr = &ast.GroupingExpr{
			Expression: inner,
		}
		return

	case token.LeftBracket:
		p.nextToken()
		return p.parseArrayExpr()
	}
	p.nextToken()
	return expr
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		switch p.tok {
		case token.Semicolon:
			p.nextToken()
			return
		case token.Class, token.Fun, token.Var, token.If, token.While, token.Print, token.Return:
			return
		}
		p.nextToken()
	}
}

func (p *Parser) match(tokens ...token.Token) bool {
	for _, tok := range tokens {
		if p.check(tok) {
			p.nextToken()
			return true
		}
	}
	return false
}

func (p *Parser) expect(tok token.Token, msg string) {
	if p.check(tok) {
		p.nextToken()
		return
	}
	p.error(msg)
}

func (p *Parser) error(msg string) {
	s := fmt.Sprintf("%s %s", p.l.Pos(), msg)
	fmt.Fprintln(os.Stderr, s)
	panic(parseError{s})
}

func (p *Parser) check(tok token.Token) bool {
	if p.isAtEnd() {
		return false
	}
	return p.tok == tok
}

func (p *Parser) isAtEnd() bool {
	return p.tok == token.EOF
}

// конструктор
func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		l: l,
	}
	parser.nextToken()
	return parser
}

// trace
func identLevel(count int) string {
	count--
	if count < 0 {
		count = 0
	}
	return strings.Repeat("\t", count)
}

func trace(p *Parser, msg string) (*Parser, string) {
	fmt.Printf("%sBEGIN %s\n", identLevel(p.indent), msg)
	p.indent++
	return p, msg
}

func unTrace(p *Parser, msg string) *Parser {
	count := p.indent
	if count < 0 {
		count = 0
	}
	fmt.Printf("%sEND %s\n", identLevel(p.indent), msg)
	p.indent--
	return p
}
