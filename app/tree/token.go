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
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	//Single- or Double-character types
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

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
