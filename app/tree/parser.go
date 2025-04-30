package tree

import (
	"fmt"
	"os"
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

func (p *Parser) Parse() Expr {
	// Will be added to later with statements
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
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
			left:     expr,
			operator: operator,
			right:    right,
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
			left:     expr,
			operator: operator,
			right:    right,
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
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{
			operator: operator,
			right:    right,
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(TRUE) {
		return LiteralExpr{
			value: true,
		}
	}
	if p.match(FALSE) {
		return LiteralExpr{
			value: false,
		}
	}
	if p.match(NIL) {
		return LiteralExpr{
			value: nil,
		}
	}

	if p.match(STRING, NUMBER) {
		return LiteralExpr{
			value: p.previous().Literal,
		}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		err := p.consume(RIGHT_PAREN)
		if err != nil {

		}
		return GroupingExpr{
			expression: expr,
		}
	}
	p.error(p.peek(), "Expect Expression")
	return nil
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
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

func (p *Parser) consume(tokenType TokenType) error {
	if p.check(tokenType) {
		return nil
	}
	return fmt.Errorf("expected %v", tokenType)
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
