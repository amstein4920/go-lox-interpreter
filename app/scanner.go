package main

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/tree"
)

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	Tokens  []tree.Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
	}
}

func (s *Scanner) ScanTokens() {
	for s.current < len(s.source) {
		s.start = s.current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, tree.Token{
		TokenType: tree.EOF,
		Lexeme:    "",
		Literal:   nil,
		Line:      s.line,
	})
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(tree.LEFT_PAREN)
	case ')':
		s.addToken(tree.RIGHT_PAREN)
	case '{':
		s.addToken(tree.LEFT_BRACE)
	case '}':
		s.addToken(tree.RIGHT_BRACE)
	default:
		fmt.Println("Do something for lexical errors")
	}
}

func (s *Scanner) advance() rune {
	char := rune(s.source[s.current])
	s.current++
	return char
}

func (s *Scanner) addToken(tokenType tree.TokenType) {
	s.addFullToken(tokenType, nil)
}

func (s *Scanner) addFullToken(tokenType tree.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.Tokens = append(s.Tokens, tree.Token{
		TokenType: tokenType,
		Lexeme:    text,
		Literal:   literal,
		Line:      s.line,
	})
}
