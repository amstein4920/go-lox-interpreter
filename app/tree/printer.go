package tree

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

// VisitBinaryExpr implements ExprVisitor.
func (a AstPrinter) VisitBinaryExpr(expr BinaryExpr) any {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

// VisitCallExpr implements ExprVisitor.
func (a AstPrinter) VisitCallExpr(expr CallExpr) any {
	panic("unimplemented")
}

// VisitGetExpr implements ExprVisitor.
func (a AstPrinter) VisitGetExpr(expr GetExpr) any {
	panic("unimplemented")
}

// VisitGroupingExpr implements ExprVisitor.
func (a AstPrinter) VisitGroupingExpr(expr GroupingExpr) any {
	return a.parenthesize("group", expr.expression)
}

// VisitLiteralExpr implements ExprVisitor.
func (a AstPrinter) VisitLiteralExpr(expr LiteralExpr) any {
	if expr.value == nil {
		return "nil"
	}
	if _, ok := expr.value.(float64); ok {
		if expr.value == float64(int64(expr.value.(float64))) {
			return fmt.Sprintf("%.1f", expr.value)
		}
	}
	return fmt.Sprint(expr.value)
}

// VisitLogicalExpr implements ExprVisitor.
func (a AstPrinter) VisitLogicalExpr(expr LogicalExpr) any {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

// VisitSetExpr implements ExprVisitor.
func (a AstPrinter) VisitSetExpr(expr SetExpr) any {
	panic("unimplemented")
}

// VisitSuperExpr implements ExprVisitor.
func (a AstPrinter) VisitSuperExpr(expr SuperExpr) any {
	panic("unimplemented")
}

// VisitThisExpr implements ExprVisitor.
func (a AstPrinter) VisitThisExpr(expr ThisExpr) any {
	return expr.keyword.Lexeme
}

// VisitUnaryExpr implements ExprVisitor.
func (a AstPrinter) VisitUnaryExpr(expr UnaryExpr) any {
	return a.parenthesize(expr.operator.Lexeme, expr.right)
}

// VisitVariableExpr implements ExprVisitor.
func (a AstPrinter) VisitVariableExpr(expr VariableExpr) any {
	return expr.name.Lexeme
}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitAssignExpr(expr AssignExpr) any {
	return a.parenthesize("= "+expr.name.Lexeme, expr.value)
}
func (a AstPrinter) parenthesize(name string, exprs ...Expr) any {
	var b strings.Builder
	b.WriteString("(" + name)
	for _, e := range exprs {
		fmt.Fprintf(&b, " %s", a.Print(e))
	}
	b.WriteString(")")
	return b.String()
}
