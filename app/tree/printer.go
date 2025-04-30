package tree

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

// VisitBinaryExpr implements ExprVisitor.
func (a AstPrinter) VisitBinaryExpr(expr BinaryExpr) string {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

// VisitCallExpr implements ExprVisitor.
func (a AstPrinter) VisitCallExpr(expr CallExpr) string {
	panic("unimplemented")
}

// VisitGetExpr implements ExprVisitor.
func (a AstPrinter) VisitGetExpr(expr GetExpr) string {
	panic("unimplemented")
}

// VisitGroupingExpr implements ExprVisitor.
func (a AstPrinter) VisitGroupingExpr(expr GroupingExpr) string {
	return a.parenthesize("group", expr.expression)
}

// VisitLiteralExpr implements ExprVisitor.
func (a AstPrinter) VisitLiteralExpr(expr LiteralExpr) string {
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
func (a AstPrinter) VisitLogicalExpr(expr LogicalExpr) string {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

// VisitSetExpr implements ExprVisitor.
func (a AstPrinter) VisitSetExpr(expr SetExpr) string {
	panic("unimplemented")
}

// VisitSuperExpr implements ExprVisitor.
func (a AstPrinter) VisitSuperExpr(expr SuperExpr) string {
	panic("unimplemented")
}

// VisitThisExpr implements ExprVisitor.
func (a AstPrinter) VisitThisExpr(expr ThisExpr) string {
	return expr.keyword.Lexeme
}

// VisitUnaryExpr implements ExprVisitor.
func (a AstPrinter) VisitUnaryExpr(expr UnaryExpr) string {
	return a.parenthesize(expr.operator.Lexeme, expr.right)
}

// VisitVariableExpr implements ExprVisitor.
func (a AstPrinter) VisitVariableExpr(expr VariableExpr) string {
	return expr.name.Lexeme
}

func (a AstPrinter) Print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitAssignExpr(expr AssignExpr) string {
	return a.parenthesize("= "+expr.name.Lexeme, expr.value)
}
func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var b strings.Builder
	b.WriteString("(" + name)
	for _, e := range exprs {
		fmt.Fprintf(&b, " %s", a.Print(e))
	}
	b.WriteString(")")
	return b.String()
}
