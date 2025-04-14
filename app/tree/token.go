package tree

import (
	"fmt"
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func (t Token) String() string {
	var literal any
	if t.Literal == nil {
		literal = "null"
	} else {
		literal = t.Literal
	}
	return fmt.Sprintf("%s %s %s", t.TokenType.String(), t.Lexeme, literal)
}

type TokenType int

const (
	//Single-character types
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	//Single- or Double-character types
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	//Literals
	Identifier
	String
	Number

	//Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	EOF
)
