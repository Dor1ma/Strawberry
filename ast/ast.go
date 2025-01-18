package ast

import (
	"fmt"
	"github.com/Dor1ma/Strawberry/token"
	"strings"
)

// Node Нода ! все остальные ноды, должны реализовывать этот интефрейс
type Node interface {
	node()
	String() string
}

// Expression Выражение
type Expression interface {
	Node
	expr()
}

// Statement Оператор
type Statement interface {
	Node
	stmt()
}

// LeftExpr левое выражение для присваивания
type LeftExpr interface {
	Expression
	leftExpr()
}

func (*Identifier) node() {}

func (*Literal) node() {}

func (*AssignExpr) node()   {}
func (*BinaryExpr) node()   {}
func (*CallExpr) node()     {}
func (*GetExpr) node()      {}
func (*GroupingExpr) node() {}
func (*LogicalExpr) node()  {}
func (*SetExpr) node()      {}
func (*SuperExpr) node()    {}
func (*ThisExpr) node()     {}
func (*UnaryExpr) node()    {}
func (*VariableExpr) node() {}

func (*ArrayExpr) node()  {}
func (*ArrayIndex) node() {}

func (*BlockStmt) node()    {}
func (*ClassStmt) node()    {}
func (*ExprStmt) node()     {}
func (*FunctionStmt) node() {}
func (*IfStmt) node()       {}
func (*PrintStmt) node()    {}
func (*ReturnStmt) node()   {}
func (*VarStmt) node()      {}
func (*WhileStmt) node()    {}

type Identifier struct {
	Name string
}

func (ident *Identifier) String() string { return ident.Name }

type Literal struct {
	Token token.Token
	Value string
}

func (*Literal) expr() {}

func (lit *Literal) String() string {
	switch lit.Token {
	case token.Nil:
		return "null"
	case token.True:
		return "true"
	case token.False:
		return "false"
	case token.Number, token.String:
		return lit.Value
	}
	panic("unknown Literal.")
}

type (
	AssignExpr struct {
		Left  LeftExpr
		Value Expression
	}
	BinaryExpr struct {
		Left     Expression
		Operator token.Token
		Right    Expression
	}
	CallExpr struct {
		Callee    Expression
		Arguments []Expression
	}

	// ----

	ArrayExpr struct {
		Elements []Expression
	}

	ArrayIndex struct {
		Array Expression
		Index Expression
	}

	// ----

	GetExpr struct {
		Object Expression
		Name   string
	}
	GroupingExpr struct {
		Expression Expression
	}
	LogicalExpr struct {
		Left     Expression
		Operator token.Token
		Right    Expression
	}
	SetExpr struct {
		Object Expression
		Name   string
		Value  Expression
	}
	SuperExpr struct {
		// Method  Identifier
		Keyword token.Token
		Method  token.Token
	}
	ThisExpr  struct{}
	UnaryExpr struct {
		Operator token.Token
		Right    Expression
	}
	VariableExpr struct {
		Name     string
		Distance int // NOTE!! -1 используем, когда переменная ГЛОБАЛЬНАЯ
	}
)

func (*AssignExpr) expr()   {}
func (*BinaryExpr) expr()   {}
func (*CallExpr) expr()     {}
func (*GetExpr) expr()      {}
func (*GroupingExpr) expr() {}
func (*LogicalExpr) expr()  {}
func (*SetExpr) expr()      {}
func (*SuperExpr) expr()    {}
func (*ThisExpr) expr()     {}
func (*UnaryExpr) expr()    {}
func (*VariableExpr) expr() {}

func (*ArrayExpr) expr()  {}
func (*ArrayIndex) expr() {}

func (e *AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", e.Left, e.Value)
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator, e.Right)
}

func (e *CallExpr) String() string {
	args := make([]string, len(e.Arguments))
	for i, arg := range e.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", e.Callee, strings.Join(args, ", "))
}

func (e *ArrayExpr) String() string {
	var elements []string
	for _, el := range e.Elements {
		elements = append(elements, el.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

func (e *ArrayIndex) String() string {
	return fmt.Sprintf("%s[%s]", e.Array.String(), e.Index.String())
}

func (e *GetExpr) String() string {
	return e.Object.String() + "." + e.Name
}

func (e *GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", e.Expression)
}

func (e *LogicalExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left, e.Operator, e.Right)
}

func (e *SetExpr) String() string {
	return fmt.Sprintf("%s.%s = %s", e.Object, e.Name, e.Value)
}

func (e *SuperExpr) String() string {
	return "TODO"
}

func (e *ThisExpr) String() string {
	return "this"
}

func (e *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right)
}

func (e *VariableExpr) String() string {
	return e.Name
}

type (
	BlockStmt struct {
		Statements []Statement
	}
	ClassStmt struct {
		Name       string
		SuperClass VariableExpr
		Methods    []*FunctionStmt
	}
	ExprStmt struct {
		Expression Expression
	}
	FunctionStmt struct {
		Name          string
		Params        []*Identifier
		Body          []Statement
		IsInitializer bool
	}
	IfStmt struct {
		Condition  Expression
		ThenBranch Statement
		ElseBranch Statement
	}
	PrintStmt struct {
		Expression Expression
	}
	ReturnStmt struct {
		Keyword token.Token
		Value   Expression
	}
	VarStmt struct {
		Name        *Identifier
		Initializer Expression
	}
	WhileStmt struct {
		Condition Expression
		Body      Statement
	}
)

func (*BlockStmt) stmt()    {}
func (*ClassStmt) stmt()    {}
func (*ExprStmt) stmt()     {}
func (*FunctionStmt) stmt() {}
func (*IfStmt) stmt()       {}
func (*PrintStmt) stmt()    {}
func (*ReturnStmt) stmt()   {}
func (*VarStmt) stmt()      {}
func (*WhileStmt) stmt()    {}

func (s *BlockStmt) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, stmt := range s.Statements {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(" }")
	return sb.String()
}

func (s *ClassStmt) String() string {
	return "class " + s.Name
}

func (s *ExprStmt) String() string {
	return s.Expression.String() + ";"
}

func (s *FunctionStmt) String() string {
	var sb strings.Builder
	sb.WriteString("fun ")
	sb.WriteString(s.Name)
	sb.WriteString("(")
	params := make([]string, len(s.Params))
	for i, p := range s.Params {
		params[i] = p.Name
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") { ")
	for _, stmt := range s.Body {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(" }")
	return sb.String()
}

func (s *IfStmt) String() string {
	var sb strings.Builder
	sb.WriteString("if (")
	sb.WriteString(s.Condition.String())
	sb.WriteString(") ")
	sb.WriteString(s.ThenBranch.String())
	if s.ElseBranch != nil {
		sb.WriteString(" else ")
		sb.WriteString(s.ElseBranch.String())
	}
	return sb.String()
}

func (s *PrintStmt) String() string {
	var sb strings.Builder
	sb.WriteString("print ")
	sb.WriteString(s.Expression.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *ReturnStmt) String() string {
	str := "return"
	if s.Value != nil {
		str += " " + s.Value.String()
	}
	return str + ";"
}

func (s *VarStmt) String() string {
	var sb strings.Builder
	sb.WriteString("var ")
	sb.WriteString(s.Name.String())
	sb.WriteString(" = ")
	sb.WriteString(s.Initializer.String())
	sb.WriteRune(';')
	return sb.String()
}

func (s *WhileStmt) String() string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(s.Condition.String())
	sb.WriteString(") ")
	sb.WriteString(s.Body.String())
	return sb.String()
}

// PrettyPrint - печатает AST с отступами для лучшего восприятия.
func PrettyPrint(node Statement, indentLevel int) string {
	indent := strings.Repeat("  ", indentLevel)
	return fmt.Sprintf("%s%s", indent, node.String())
}

func (e *VariableExpr) leftExpr() {}
func (e *ArrayIndex) leftExpr()   {}
