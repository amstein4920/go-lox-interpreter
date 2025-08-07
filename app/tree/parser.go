package tree

import (
	"fmt"
	"os"
	"slices"

	. "github.com/codecrafters-io/interpreter-starter-go/app/definitions"
)

type Parser struct {
	tokens   []Token
	current  int
	HadError bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) ParseExpression() Expr {
	return p.expression()
}

func (p *Parser) ParseStatements() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	stmt := p.statement()
	if p.HadError {
		p.synchronize()
		return nil
	}
	return stmt
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expected variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return VariableStmt{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if varExpr, isVarExpr := expr.(VariableExpr); isVarExpr {
			return AssignExpr{
				Name:  varExpr.Name,
				Value: value,
			}
		}
		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(STAR, SLASH) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TRUE) {
		return LiteralExpr{
			Value: true,
		}
	}
	if p.match(FALSE) {
		return LiteralExpr{
			Value: false,
		}
	}
	if p.match(NIL) {
		return LiteralExpr{
			Value: nil,
		}
	}

	if p.match(STRING, NUMBER) {
		return LiteralExpr{
			Value: p.previous().Literal,
		}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expected ) after (")
		return GroupingExpr{
			Expression: expr,
		}
	}
	if p.match(IDENTIFIER) {
		return VariableExpr{
			Name: p.previous(),
		}
	}

	p.error(p.peek(), "Expect Expression")
	return nil
}

func (p *Parser) statement() Stmt {
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return BlockStmt{
			Statements: p.blockStatement(),
		}
	}
	return p.expressionStmt()
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() Stmt {
	if p.tokens[p.current].TokenType != SEMICOLON {
		value := p.expression()
		p.consume(SEMICOLON, "Expect ';' after expression.")
		return PrintStmt{
			Expression: value,
		}
	}
	return nil
}

func (p *Parser) blockStatement() []Stmt {
	stmts := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expected '}' after block.")
	return stmts
}

func (p *Parser) expressionStmt() Stmt {
	if p.HadError {
		for _, token := range p.tokens {
			if token.TokenType == SEMICOLON {
				break
			}
		}
		return ExpressionStmt{}
	}

	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return ExpressionStmt{
		Expression: value,
	}
}

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return Token{}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) error(token Token, message string) {
	p.HadError = true
	var where string
	if token.TokenType == EOF {
		where = " at end"
	} else {
		where = " at '" + token.Lexeme + "'"
	}
	_, _ = os.Stderr.Write([]byte(fmt.Sprintf("[line %d] Error%s: %s\n", token.Line+1, where, message)))
}
