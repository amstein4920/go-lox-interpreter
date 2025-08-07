package tree

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/app/definitions"
	. "github.com/codecrafters-io/interpreter-starter-go/app/definitions"
)

type Interpreter struct {
	StdErr   io.Writer
	HadError bool
	Env      *Environment
}

// VisitBlockStmt implements StmtVisitor.
func (in *Interpreter) VisitBlockStmt(stmt definitions.BlockStmt) {
	in.executeBlock(stmt.Statements, &Environment{
		Values: make(map[string]any),
		Parent: in.Env,
	})
}

// VisitClassStmt implements StmtVisitor.
func (in *Interpreter) VisitClassStmt(stmt definitions.ClassStmt) {
	panic("unimplemented")
}

// VisitFunctionStmt implements StmtVisitor.
func (in *Interpreter) VisitFunctionStmt(stmt definitions.FunctionStmt) {
	panic("unimplemented")
}

// VisitIfStmt implements StmtVisitor.
func (in *Interpreter) VisitIfStmt(stmt definitions.IfStmt) {
	if isTruthy(in.evaluate(stmt.Condition)) {
		in.execute(stmt.ThenBranch)
		return
	}
	if stmt.ElseBranch != nil {
		in.execute(stmt.ElseBranch)
	}
}

// VisitReturnStmt implements StmtVisitor.
func (in *Interpreter) VisitReturnStmt(stmt definitions.ReturnStmt) {
	panic("unimplemented")
}

// VisitVariableStmt implements StmtVisitor.
func (in *Interpreter) VisitVariableStmt(stmt definitions.VariableStmt) {
	var value any = nil
	if stmt.Initializer != nil {
		value = in.evaluate(stmt.Initializer)
	}
	in.Env.Define(stmt.Name.Lexeme, value)
}

// VisitWhileStmt implements StmtVisitor.
func (in *Interpreter) VisitWhileStmt(stmt definitions.WhileStmt) {
	panic("unimplemented")
}

func NewInterpreter(stdErr io.Writer) *Interpreter {

	return &Interpreter{
		StdErr: stdErr,
		Env:    definitions.NewEnvironment(),
	}
}

// VisitAssignExpr implements ExprVisitor.
func (in *Interpreter) VisitAssignExpr(expr definitions.AssignExpr) any {
	value := in.evaluate(expr.Value)
	in.Env.Assign(expr.Name, value)
	return value
}

// VisitBinaryExpr implements ExprVisitor.
func (in *Interpreter) VisitBinaryExpr(expr definitions.BinaryExpr) any {
	left := in.evaluate(expr.Left)
	right := in.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case definitions.MINUS:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}
		return left.(float64) - right.(float64)
	case SLASH:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}
		return left.(float64) / right.(float64)
	case STAR:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}

		return left.(float64) * right.(float64)
	case PLUS:
		leftType := reflect.TypeOf(left)
		rightType := reflect.TypeOf(right)

		if (leftType != nil && leftType.Kind() == reflect.Float64) && (rightType != nil && rightType.Kind() == reflect.Float64) {
			return left.(float64) + right.(float64)
		}
		if (leftType != nil && leftType.Kind() == reflect.String) && (rightType != nil && rightType.Kind() == reflect.String) {
			return left.(string) + right.(string)
		}
		in.error(expr.Operator, "Operands must be two numbers or two strings.")
		return ""
	case GREATER:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}
		return left.(float64) >= right.(float64)
	case LESS:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		err := in.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return ""
		}
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}
	return nil // Shouldn't happen
}

// VisitCallExpr implements ExprVisitor.
func (in *Interpreter) VisitCallExpr(expr definitions.CallExpr) any {
	panic("unimplemented")
}

// VisitGetExpr implements ExprVisitor.
func (in *Interpreter) VisitGetExpr(expr definitions.GetExpr) any {
	panic("unimplemented")
}

// VisitGroupingExpr implements ExprVisitor.
func (in *Interpreter) VisitGroupingExpr(expr definitions.GroupingExpr) any {
	return in.evaluate(expr.Expression)
}

// VisitLogicalExpr implements ExprVisitor.
func (in *Interpreter) VisitLogicalExpr(expr definitions.LogicalExpr) any {
	left := in.evaluate(expr.Left)

	if expr.Operator.TokenType == OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return in.evaluate(expr.Right)
}

// VisitSetExpr implements ExprVisitor.
func (in *Interpreter) VisitSetExpr(expr definitions.SetExpr) any {
	panic("unimplemented")
}

// VisitSuperExpr implements ExprVisitor.
func (in *Interpreter) VisitSuperExpr(expr definitions.SuperExpr) any {
	panic("unimplemented")
}

// VisitThisExpr implements ExprVisitor.
func (in *Interpreter) VisitThisExpr(expr definitions.ThisExpr) any {
	panic("unimplemented")
}

func (in *Interpreter) VisitUnaryExpr(expr definitions.UnaryExpr) any {
	right := in.evaluate(expr.Right)
	switch expr.Operator.TokenType {
	case MINUS:
		err := in.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return ""
		}
		return -right.(float64)
	case BANG:
		return !isTruthy(right)
	}
	return nil //Shouldn't be reached
}

// VisitVariableExpr implements ExprVisitor.
func (in *Interpreter) VisitVariableExpr(expr definitions.VariableExpr) any {
	value, err := in.Env.Get(expr.Name)
	if err != nil {
		in.error(expr.Name, err.Error())
	}
	return value
}

func (in *Interpreter) VisitLiteralExpr(expr definitions.LiteralExpr) any {
	return expr.Value
}

func (in *Interpreter) VisitExpressionStmt(stmt definitions.ExpressionStmt) {
	val := in.evaluate(stmt.Expression)
	println(val)
}

func (in *Interpreter) VisitPrintStmt(stmt definitions.PrintStmt) {
	if stmt.Expression != nil {
		value := in.evaluate(stmt.Expression)
		if !in.HadError {
			fmt.Println(stringify(value))
		}
	}
}

func (in *Interpreter) evaluate(expr definitions.Expr) any {
	return expr.Accept(in)
}

func (in *Interpreter) executeBlock(stmts []Stmt, env *Environment) {
	previous := in.Env
	in.Env = env
	for _, stmt := range stmts {
		in.execute(stmt)
	}
	in.Env = previous
}

func (in *Interpreter) checkNumberOperand(operator definitions.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return in.error(operator, "Operand must be a number.")
}

func (in *Interpreter) checkNumberOperands(operator definitions.Token, left any, right any) error {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}
	return in.error(operator, "Operands must be numbers.")
}

func (in *Interpreter) error(token definitions.Token, msg string) error {
	in.HadError = true
	_, _ = in.StdErr.Write([]byte(fmt.Sprintf("[line %d] Error: %s\n", token.Line, msg)))
	return fmt.Errorf("[line %d] Error: %s\n", token.Line, msg)
}

func (in *Interpreter) InterpretExpression(expr definitions.Expr) {
	value := in.evaluate(expr)
	fmt.Println(stringify(value))
}

func (in *Interpreter) Interpret(stmts []definitions.Stmt) {
	for _, stmt := range stmts {
		if stmt != nil && !in.HadError {
			in.execute(stmt)
		}
	}
}

func (in *Interpreter) execute(stmt definitions.Stmt) {
	stmt.Accept(in)
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if boolean, ok := obj.(bool); ok {
		return boolean
	}
	return true
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	if _, ok := object.(float64); ok {
		return strconv.FormatFloat(object.(float64), 'f', -1, 64)
	}
	return fmt.Sprintf("%v", object)
}
