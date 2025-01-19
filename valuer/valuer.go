package valuer

import (
	"fmt"
	"github.com/Dor1ma/Strawberry/ast"
	"strconv"
	"strings"
)

var typeMap = map[Type]string{
	NumberType:   "number",
	StringType:   "string",
	BooleanType:  "bool",
	ArrayType:    "array",
	NilType:      "nil",
	FunctionType: "function",
	ReturnType:   "return",
	ClassType:    "class",
}

type Type int

const (
	NumberType   Type = iota + 1 // number
	StringType                   // string
	BooleanType                  // bool
	ArrayType                    // array
	NilType                      // nil
	FunctionType                 // function
	ReturnType                   // return
	ClassType                    // class
	InstanceType                 // instance
)

func (typ Type) String() string {
	if s, ok := typeMap[typ]; ok {
		return s
	}
	return "unknown"
}

type Valuer interface {
	Type() Type
	String() string
}

type Callable interface {
	call()
	Arity() int
}

type Number struct {
	Value float64
}

func (*Number) Type() Type { return NumberType }

func (num *Number) String() string {
	return strconv.FormatFloat(num.Value, 'f', -1, 64)
}

type String struct {
	Value string
}

func (*String) Type() Type { return StringType }

func (s *String) String() string { return s.Value }

type Boolean struct {
	Value bool
}

func (*Boolean) Type() Type { return BooleanType }

func (b *Boolean) String() string { return strconv.FormatBool(b.Value) }

type Nil struct{}

// Type returns its Type.
func (*Nil) Type() Type { return NilType }

func (*Nil) String() string { return "nil" }

type Function struct {
	Name          string
	Params        []*ast.Identifier
	Body          []ast.Statement
	Closure       *Environment
	IsInitializer bool
}

func (*Function) Type() Type { return FunctionType }

func (*Function) call() {}

func (fn *Function) String() string {
	return "<fn " + fn.Name + ">"
}

// что то типа "арность"
func (fn *Function) Arity() int {
	return len(fn.Params)
}

func (fn *Function) Bind(instance *Instance) *Function {
	environment := NewEnclosing(fn.Closure)
	environment.Define("this", instance)
	return &Function{
		Name:    fn.Name,
		Params:  fn.Params,
		Body:    fn.Body,
		Closure: environment,
	}
}

type ReturnValue struct {
	Value Valuer
}

func (*ReturnValue) Type() Type { return ReturnType }

func (rt *ReturnValue) String() string {
	return rt.Value.String()
}

type Array struct {
	Elements []Valuer
}

func (a *Array) Type() Type {
	return ArrayType
}

func (a *Array) String() string {
	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

type ClassValue struct {
	Name    string
	Methods map[string]*Function
}

func (*ClassValue) Type() Type { return ClassType }

func (*ClassValue) call() {}

func (c *ClassValue) Arity() int {
	initializer := c.FindMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *ClassValue) String() string {
	return "class " + c.Name
}

func (c *ClassValue) FindMethod(key string) *Function {
	if method, ok := c.Methods[key]; ok {
		return method
	}
	return nil
}

type Instance struct {
	Klass  *ClassValue
	Fields map[string]Valuer
}

func (*Instance) Type() Type { return ClassType }

func (i *Instance) String() string {
	return i.Klass.Name + " instance"
}

func (i *Instance) Get(key string) (Valuer, bool) {
	if v, ok := i.Fields[key]; ok {
		return v, ok
	}
	if method := i.Klass.FindMethod(key); method != nil {
		return method.Bind(i), true
	}
	return nil, false
}

func (i *Instance) Set(key string, v Valuer) {
	if i.Fields == nil {
		i.Fields = make(map[string]Valuer)
	}
	i.Fields[key] = v
}
