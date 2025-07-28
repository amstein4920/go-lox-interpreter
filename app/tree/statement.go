package tree

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt BlockStmt)
	VisitClassStmt(stmt ClassStmt)
	VisitExpressionStmt(stmt ExpressionStmt)
	VisitFunctionStmt(stmt FunctionStmt)
	VisitIfStmt(stmt IfStmt)
	VisitPrintStmt(stmt PrintStmt)
	VisitReturnStmt(stmt ReturnStmt)
	VisitVariableStmt(stmt VariableStmt)
	VisitWhileStmt(stmt WhileStmt)
}

type BlockStmt struct {
	statements []Stmt
}

func (s BlockStmt) Accept(v StmtVisitor) {
	v.VisitBlockStmt(s)
}

type ClassStmt struct {
	name       Token
	superClass VariableExpr
	methods    []FunctionStmt
}

func (s ClassStmt) Accept(v StmtVisitor) {
	v.VisitClassStmt(s)
}

type ExpressionStmt struct {
	expression Expr
}

func (s ExpressionStmt) Accept(v StmtVisitor) {
	v.VisitExpressionStmt(s)
}

type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func (s FunctionStmt) Accept(v StmtVisitor) {
	v.VisitFunctionStmt(s)
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (s IfStmt) Accept(v StmtVisitor) {
	v.VisitIfStmt(s)
}

type PrintStmt struct {
	expression Expr
}

func (s PrintStmt) Accept(v StmtVisitor) {
	v.VisitPrintStmt(s)
}

type ReturnStmt struct {
	keyword Token
	value   Expr
}

func (s ReturnStmt) Accept(v StmtVisitor) {
	v.VisitReturnStmt(s)
}

type VariableStmt struct {
	name        Token
	initializer Expr
}

func (s VariableStmt) Accept(v StmtVisitor) {
	v.VisitVariableStmt(s)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s WhileStmt) Accept(v StmtVisitor) {
	v.VisitWhileStmt(s)
}
