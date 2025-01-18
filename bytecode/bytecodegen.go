package bytecode_gen

import (
	"fmt"
	"github.com/Dor1ma/Strawberry/ast"
	"github.com/Dor1ma/Strawberry/token"
)

const (
	NEG = "NEG" // Унарный минус
	NOT = "NOT" // Логическое отрицание
	ADD = "ADD" // Сложение
	SUB = "SUB" // Вычитание
	MUL = "MUL" // Умножение
	DIV = "DIV" // Деление
	AND = "AND" // Деление
	OR  = "OR"  // Деление

	LESS_THAN          = "LESS_THAN"          // Меньше
	GREATER_THAN       = "GREATER_THAN"       // Больше
	LESS_EQUAL_THAN    = "LESS_EQUAL_THAN"    // Меньше или равно
	GREATER_EQUAL_THAN = "GREATER_EQUAL_THAN" // Больше или равно
	EQUAL              = "EQUAL"              // Равенство
	NOT_EQUAL          = "NOT_EQUAL"          // Неравенство

	JUMP          = "JUMP"          // Безусловный переход
	JUMP_IF_FALSE = "JUMP_IF_FALSE" // Переход, если условие ложно
	CALL_FUNCTION = "CALL_FUNCTION" // Вызов функции
	RETURN        = "RETURN"        // Возврат из функции

	PUSH_CONST = "PUSH_CONST" // Поместить константу в стек
	PUSH_VAR   = "PUSH_VAR"   // Поместить значение переменной в стек
	STORE_VAR  = "STORE_VAR"  // Сохранить значение в переменной
	NEW_ARRAY  = "NEW_ARRAY"  // Создать новый массив
	ARRAY_GET  = "ARRAY_GET"  // Получить значение из массива
	ARRAY_SET  = "ARRAY_SET"  // Установить значение в массиве

	GET_PROPERTY = "GET_PROPERTY" // Получить свойство объекта
	SET_PROPERTY = "SET_PROPERTY" // Установить свойство объекта

	PRINT = "PRINT" // Вывод значения
	SUPER = "SUPER" // Обращение к суперклассу
	THIS  = "THIS"  // Ссылка на текущий объект
	FUNC  = "FUNC"  // Объявление функции

	LABEL            = "label"        // Метка для перехода
	FALSE_LABEL      = "false_label_" // Метка для перехода при falsr
	LOOP_START_LABEL = "loop_start_"  // Метка старта цикла
	LOOP_END_LABEL   = "loop_end_"    // Метка конца цикла
	END_LABEL        = "end_label_"   // Метка конца

	SCOPE_START = "scope_start"
	SCOPE_END   = "scope_end"

	NULL = "NULL"
)

type Bytecode struct {
	Opcode string
	Arg    string
}

type CodeGenerator struct {
	Bytecodes []Bytecode
}

func (cg *CodeGenerator) emit(opcode, arg string) {
	cg.Bytecodes = append(cg.Bytecodes, Bytecode{Opcode: opcode, Arg: arg})
}

func (cg *CodeGenerator) GenerateProgram(statements []ast.Statement) {
	for _, stmt := range statements {
		cg.GenerateStatement(stmt)
	}
}

func (cg *CodeGenerator) PrintBytecode() {
	for _, bc := range cg.Bytecodes {
		fmt.Printf("%s %s\n", bc.Opcode, bc.Arg)
	}
}

func (cg *CodeGenerator) GenerateExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.Literal:
		cg.GenerateLiteral(e)
	case *ast.VariableExpr:
		cg.GenerateVariable(e)
	case *ast.BinaryExpr:
		cg.GenerateBinaryExpr(e)
	case *ast.CallExpr:
		cg.GenerateCallExpr(e)
	case *ast.AssignExpr:
		cg.GenerateAssignExpr(e)
	case *ast.UnaryExpr:
		cg.GenerateUnaryExpr(e)
	case *ast.LogicalExpr:
		cg.GenerateLogicalExpr(e)
	case *ast.GroupingExpr:
		cg.GenerateGroupingExpr(e)
	case *ast.ArrayExpr:
		cg.GenerateArrayExpr(e)
	case *ast.ArrayIndex:
		cg.GenerateArrayIndex(e)
	case *ast.GetExpr:
		cg.GenerateGetExpr(e)
	case *ast.SetExpr:
		cg.GenerateSetExpr(e)
	case *ast.SuperExpr:
		cg.GenerateSuperExpr(e)
	case *ast.ThisExpr:
		cg.GenerateThisExpr(e)
	}
}

func (cg *CodeGenerator) GenerateLiteral(lit *ast.Literal) {
	var value string
	switch lit.Token {
	case token.Number, token.String:
		value = lit.Value
	case token.Nil:
		value = "null"
	case token.True:
		value = "true"
	case token.False:
		value = "false"
	default:
		panic("unhandled token for literal")
	}
	cg.emit(PUSH_CONST, value)
}

func (cg *CodeGenerator) GenerateVariable(varExpr *ast.VariableExpr) {
	cg.emit(PUSH_VAR, varExpr.Name)
}

func (cg *CodeGenerator) GenerateBinaryExpr(binExpr *ast.BinaryExpr) {
	cg.GenerateExpression(binExpr.Left)
	cg.GenerateExpression(binExpr.Right)

	var opcode string
	switch binExpr.Operator {
	case token.Plus:
		opcode = ADD
	case token.Minus:
		opcode = SUB
	case token.Star:
		opcode = MUL
	case token.Slash:
		opcode = DIV
	case token.Less:
		opcode = LESS_THAN
	case token.Greater:
		opcode = GREATER_THAN
	case token.GreaterThanOrEqual:
		opcode = GREATER_EQUAL_THAN
	case token.LessThanOrEqual:
		opcode = LESS_EQUAL_THAN
	case token.EqualEqual:
		opcode = EQUAL
	case token.NotEqual:
		opcode = NOT_EQUAL
	default:
		panic("unhandled token for binary operation")
	}
	cg.emit(opcode, "")
}

func (cg *CodeGenerator) GenerateCallExpr(call *ast.CallExpr) {
	for _, arg := range call.Arguments {
		cg.GenerateExpression(arg)
	}

	cg.emit(CALL_FUNCTION, call.Callee.String())
}

func (cg *CodeGenerator) GenerateLeftExpr(left ast.LeftExpr) {
	switch l := left.(type) {
	case *ast.VariableExpr:
		cg.emit(STORE_VAR, l.Name)

	case *ast.ArrayIndex:
		cg.GenerateExpression(l.Array)
		cg.GenerateExpression(l.Index)
		cg.emit(ARRAY_SET, "")
	default:
		panic(fmt.Sprintf("unsupported left-side expression: %T", l))
	}
}

func (cg *CodeGenerator) GenerateAssignExpr(assign *ast.AssignExpr) {
	cg.GenerateExpression(assign.Value)
	cg.GenerateLeftExpr(assign.Left)
}

func (cg *CodeGenerator) GenerateIfStmt(ifStmt *ast.IfStmt) {
	cg.GenerateExpression(ifStmt.Condition)

	falseLabel := fmt.Sprintf("%s%d", FALSE_LABEL, len(cg.Bytecodes))
	cg.emit(JUMP_IF_FALSE, falseLabel)

	cg.GenerateStatement(ifStmt.ThenBranch)

	if ifStmt.ElseBranch != nil {
		falseLabelEnd := fmt.Sprintf("%s%d", END_LABEL, len(cg.Bytecodes))
		cg.emit(JUMP, falseLabelEnd)
		cg.emit(LABEL, falseLabel)

		cg.GenerateStatement(ifStmt.ElseBranch)
		cg.emit(LABEL, falseLabelEnd)
	} else {
		cg.emit(LABEL, falseLabel)
	}
}

func (cg *CodeGenerator) GenerateWhileStmt(whileStmt *ast.WhileStmt) {
	loopStartLabel := fmt.Sprintf("%s%d", LOOP_START_LABEL, len(cg.Bytecodes))
	loopEndLabel := fmt.Sprintf("%s%d", LOOP_END_LABEL, len(cg.Bytecodes))

	cg.emit(LABEL, loopStartLabel)
	cg.GenerateExpression(whileStmt.Condition)

	cg.emit(JUMP_IF_FALSE, loopEndLabel)

	cg.GenerateStatement(whileStmt.Body)

	cg.emit(JUMP, loopStartLabel)
	cg.emit(LABEL, loopEndLabel)
}

func (cg *CodeGenerator) GenerateFunctionStmt(funcStmt *ast.FunctionStmt) {
	cg.emit(FUNC, funcStmt.Name)

	for _, stmt := range funcStmt.Body {
		cg.GenerateStatement(stmt)
	}
}

func (cg *CodeGenerator) GeneratePrintStmt(printStmt *ast.PrintStmt) {
	cg.GenerateExpression(printStmt.Expression)

	cg.emit(PRINT, "")
}

func (cg *CodeGenerator) GenerateStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		cg.GenerateExpression(s.Expression)
	case *ast.VarStmt:
		cg.GenerateAssignExpr(&ast.AssignExpr{Left: &ast.VariableExpr{Name: s.Name.Name}, Value: s.Initializer})
	case *ast.IfStmt:
		cg.GenerateIfStmt(s)
	case *ast.WhileStmt:
		cg.GenerateWhileStmt(s)
	case *ast.PrintStmt:
		cg.GeneratePrintStmt(s)
	case *ast.FunctionStmt:
		cg.GenerateFunctionStmt(s)
	case *ast.BlockStmt:
		cg.GenerateBlockStmt(s)
	case *ast.ReturnStmt:
		cg.GenerateReturnStmt(s)
	/*case *ast.ClassStmt:
	cg.GenerateClassStmt(s)*/
	default:
		panic(fmt.Sprintf("unknown ast type to generate bytecode statement: %T", s))
	}
}

func (cg *CodeGenerator) GenerateBlockStmt(stmt *ast.BlockStmt) {
	cg.emit(SCOPE_START, "")

	for _, statement := range stmt.Statements {
		cg.GenerateStatement(statement)
	}

	cg.emit(SCOPE_END, "")
}

func (cg *CodeGenerator) GenerateReturnStmt(stmt *ast.ReturnStmt) {
	if stmt.Value != nil {
		cg.GenerateExpression(stmt.Value)
		cg.emit(RETURN, stmt.Value.String())
	} else {
		cg.emit(PUSH_CONST, NULL)
	}
}

/* ToDO: дописать реализацию байт кодов для классов
func (cg *CodeGenerator) GenerateClassStmt(stmt *ast.ClassStmt) {
	cg.Emit(OpClass, cg.AddConstant(stmt.Name))
	if stmt.SuperClass.Name != "" {
		cg.GenerateExpression(&stmt.SuperClass)
		cg.Emit(OpInherit)
	}
	for _, method := range stmt.Methods {
		cg.GenerateFunctionStmt(method)
		cg.Emit(OpMethod, cg.AddConstant(method.Name))
	}
}*/

func (cg *CodeGenerator) GenerateUnaryExpr(unary *ast.UnaryExpr) {
	cg.GenerateExpression(unary.Right)

	var opcode string

	switch unary.Operator {
	case token.Minus:
		opcode = NEG
	case token.Not:
		opcode = NOT
	default:
		panic("unhandled token for unary operator")
	}
	cg.emit(opcode, "")
}

func (cg *CodeGenerator) GenerateLogicalExpr(logical *ast.LogicalExpr) {
	cg.GenerateExpression(logical.Left)
	cg.GenerateExpression(logical.Right)

	var opcode string
	switch logical.Operator {
	case token.And:
		opcode = AND
	case token.Or:
		opcode = OR
	default:
		panic("unhandled token for logical expression")
	}
	cg.emit(opcode, "")
}

func (cg *CodeGenerator) GenerateGroupingExpr(grouping *ast.GroupingExpr) {
	cg.GenerateExpression(grouping.Expression)
}

func (cg *CodeGenerator) GenerateArrayExpr(array *ast.ArrayExpr) {
	for _, element := range array.Elements {
		cg.GenerateExpression(element)
	}

	cg.emit(NEW_ARRAY, fmt.Sprintf("%d", len(array.Elements)))
}

func (cg *CodeGenerator) GenerateArrayIndex(arrayIndex *ast.ArrayIndex) {
	cg.GenerateExpression(arrayIndex.Array)
	cg.GenerateExpression(arrayIndex.Index)

	cg.emit(ARRAY_GET, "")
}

func (cg *CodeGenerator) GenerateGetExpr(get *ast.GetExpr) {
	cg.GenerateExpression(get.Object)

	cg.emit(GET_PROPERTY, get.Name)
}

func (cg *CodeGenerator) GenerateSetExpr(set *ast.SetExpr) {
	cg.GenerateExpression(set.Object)

	cg.GenerateExpression(set.Value)

	cg.emit(SET_PROPERTY, set.Name)
}

func (cg *CodeGenerator) GenerateSuperExpr(super *ast.SuperExpr) {
	// ToDo: реализовать генерацию байткода
	panic("npt implemented bytecode gen for super")
	cg.emit(SUPER, "")
}

func (cg *CodeGenerator) GenerateThisExpr(this *ast.ThisExpr) {
	// ToDo: реализовать генерацию байткода
	panic("npt implemented bytecode gen for this")
	cg.emit(THIS, "")
}
