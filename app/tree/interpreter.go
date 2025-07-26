package tree

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Interpreter struct {
	StdErr   io.Writer
	HadError bool
}

func NewInterpreter(stdErr io.Writer) *Interpreter {
	return &Interpreter{
		StdErr: stdErr,
	}
}

// VisitAssignExpr implements ExprVisitor.
func (in *Interpreter) VisitAssignExpr(expr AssignExpr) any {
	panic("unimplemented")
}

// VisitBinaryExpr implements ExprVisitor.
func (in *Interpreter) VisitBinaryExpr(expr BinaryExpr) any {
	left := in.evaluate(expr.left)
	right := in.evaluate(expr.right)

	switch expr.operator.TokenType {
	case MINUS:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) - right.(float64)
	case SLASH:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) * right.(float64)
	case PLUS:
		leftType := reflect.TypeOf(left)
		rightType := reflect.TypeOf(right)

		if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		if leftType.Kind() == reflect.String && rightType.Kind() == reflect.String {
			return left.(string) + right.(string)
		}
		in.error(expr.operator, "Operands must be two numbers or two strings.")
	case GREATER:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		in.checkNumberOperands(expr.operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	}
	return nil // Shouldn't happen
}

// VisitCallExpr implements ExprVisitor.
func (in *Interpreter) VisitCallExpr(expr CallExpr) any {
	panic("unimplemented")
}

// VisitGetExpr implements ExprVisitor.
func (in *Interpreter) VisitGetExpr(expr GetExpr) any {
	panic("unimplemented")
}

// VisitGroupingExpr implements ExprVisitor.
func (in *Interpreter) VisitGroupingExpr(expr GroupingExpr) any {
	return in.evaluate(expr.expression)
}

// VisitLogicalExpr implements ExprVisitor.
func (in *Interpreter) VisitLogicalExpr(expr LogicalExpr) any {
	panic("unimplemented")
}

// VisitSetExpr implements ExprVisitor.
func (in *Interpreter) VisitSetExpr(expr SetExpr) any {
	panic("unimplemented")
}

// VisitSuperExpr implements ExprVisitor.
func (in *Interpreter) VisitSuperExpr(expr SuperExpr) any {
	panic("unimplemented")
}

// VisitThisExpr implements ExprVisitor.
func (in *Interpreter) VisitThisExpr(expr ThisExpr) any {
	panic("unimplemented")
}

func (in *Interpreter) VisitUnaryExpr(expr UnaryExpr) any {
	right := in.evaluate(expr.right)
	switch expr.operator.TokenType {
	case MINUS:
		in.checkNumberOperand(expr.operator, right)
		return -right.(float64)
	case BANG:
		return !isTruthy(right)
	}
	return nil //Shouldn't be reached
}

// VisitVariableExpr implements ExprVisitor.
func (in *Interpreter) VisitVariableExpr(expr VariableExpr) any {
	panic("unimplemented")
}

func (in *Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.value
}

func (in *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(in)
}

func (in *Interpreter) checkNumberOperand(operator Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	in.error(operator, "Operand must be a number.")
}

func (in *Interpreter) checkNumberOperands(operator Token, left any, right any) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	in.error(operator, "Operands must be numbers.")
}

func (in *Interpreter) error(token Token, msg string) {
	in.HadError = true
	_, _ = in.StdErr.Write([]byte(fmt.Sprintf("[line %d] Error: %s\n", token.Line, msg)))
}

func (in *Interpreter) Interpret(expr Expr) {
	value := in.evaluate(expr)
	fmt.Println(stringify(value))
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
		// _, frac := math.Modf(object.(float64))
		// if frac == 0 {
		// 	return fmt.Sprintf("%.1f", object)
		// }
		// text := fmt.Sprintf("%.2f", object)

		// return text
		return strconv.FormatFloat(object.(float64), 'f', -1, 64)
	}
	return fmt.Sprintf("%v", object)
}
