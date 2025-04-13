package main

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/interpreter-starter-go/app/tree"
)

type Scanner struct {
	source   string
	start    int
	current  int
	line     int
	Tokens   []tree.Token
	StdErr   io.Writer
	HadError bool
}

func NewScanner(source string, stdErr io.Writer) *Scanner {
	return &Scanner{
		source: source,
		StdErr: stdErr,
		line:   1,
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
	case ',':
		s.addToken(tree.COMMA)
	case '.':
		s.addToken(tree.DOT)
	case '-':
		s.addToken(tree.MINUS)
	case '+':
		s.addToken(tree.PLUS)
	case ';':
		s.addToken(tree.SEMICOLON)
	case '/':
		s.addToken(tree.SLASH)
	case '*':
		s.addToken(tree.STAR)
	default:
		s.error("Unexpected character: " + string(char))
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

func (s *Scanner) error(message string) {
	s.HadError = true
	_, _ = s.StdErr.Write([]byte(fmt.Sprintf("[line %d] Error: %s\n", s.line, message)))
}
