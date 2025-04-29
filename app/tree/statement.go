package tree

type Stmt interface {
	Accept(visitor *StmtVisitor)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *BlockStmt) any
	VisitClassStmt(stmt *ClassStmt) any
	VisitExpressionStmt(stmt *ExpressionStmt) any
	VisitFunctionStmt(stmt *FunctionStmt) any
	VisitIfStmt(stmt *IfStmt) any
	VisitPrintStmt(stmt *PrintStmt) any
	VisitReturnStmt(stmt *ReturnStmt) any
	VisitVariableStmt(stmt *VariableStmt) any
	VisitWhileStmt(stmt *WhileStmt) any
}

type BlockStmt struct {
	statements []Stmt
}

func (s *BlockStmt) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(s)
}

type ClassStmt struct {
	name       Token
	superClass VariableExpr
	methods    []FunctionStmt
}

func (s *ClassStmt) Accept(v StmtVisitor) any {
	return v.VisitClassStmt(s)
}

type ExpressionStmt struct {
	expression Expr
}

func (s *ExpressionStmt) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(s)
}

type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func (s *FunctionStmt) Accept(v StmtVisitor) any {
	return v.VisitFunctionStmt(s)
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (s *IfStmt) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(s)
}

type PrintStmt struct {
	expression Expr
}

func (s *PrintStmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(s)
}

type ReturnStmt struct {
	keyword Token
	value   Expr
}

func (s *ReturnStmt) Accept(v StmtVisitor) any {
	return v.VisitReturnStmt(s)
}

type VariableStmt struct {
	name        Token
	initializer Expr
}

func (s *VariableStmt) Accept(v StmtVisitor) any {
	return v.VisitVariableStmt(s)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s *WhileStmt) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(s)
}
