package ast

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"github.com/alecthomas/participle/v2/lexer"
)

// Error represents an error while lexing.
//
// It complies with the participle.Error interface.
type lexerError struct {
	Msg string
	Pos lexer.Position
}

// Creates a new Error at the given position.
func errorf(pos lexer.Position, format string, args ...interface{}) *lexerError {
	return &lexerError{Msg: fmt.Sprintf(format, args...), Pos: pos}
}

func (e *lexerError) Message() string          { return e.Msg } // nolint: golint
func (e *lexerError) Position() lexer.Position { return e.Pos } // nolint: golint

// Error formats the error with FormatError.
func (e *lexerError) Error() string { return formatError(e.Pos, e.Msg) }

// An error in the form "[<filename>:][<line>:<pos>:] <message>"
func formatError(pos lexer.Position, message string) string {
	msg := ""
	if pos.Filename != "" {
		msg += pos.Filename + ":"
	}
	if pos.Line != 0 || pos.Column != 0 {
		msg += fmt.Sprintf("%d:%d:", pos.Line, pos.Column)
	}
	if msg != "" {
		msg += " " + message
	} else {
		msg = message
	}
	return msg
}

type lexFac struct{}
type lexScan struct {
	file    string
	scanner *scanner.Scanner
	err     error
}

var _ lexer.Definition = &lexFac{}
var _ lexer.Lexer = &lexScan{}

func (l *lexFac) Lex(f string, r io.Reader) (lexer.Lexer, error) {
	s := &scanner.Scanner{}
	s.Init(r)
	lex := &lexScan{
		file:    f,
		scanner: s,
	}
	lex.scanner.Error = func(s *scanner.Scanner, msg string) {
		// This is to support single quoted strings. Hacky.
		if !strings.HasSuffix(msg, "char literal") {
			lex.err = errorf(lexer.Position(lex.scanner.Pos()), msg)
		}
	}

	return lex, nil
}

const (
	eof      = -1
	operator = scanner.Comment - 1
)

func (l *lexFac) Symbols() map[string]lexer.TokenType {
	return map[string]lexer.TokenType{
		"EOF":       eof,
		"Char":      scanner.Char,
		"Ident":     scanner.Ident,
		"Int":       scanner.Int,
		"Float":     scanner.Float,
		"String":    scanner.String,
		"RawString": scanner.RawString,
		"Comment":   scanner.Comment,
		"Operator":  operator,
	}
}

var puncts = map[rune]struct{}{
	'+': {},
	'-': {},
	'*': {},
	'/': {},
	'^': {},
	'%': {},
	'=': {},
	'!': {},
	'<': {},
	'>': {},
	'&': {},
	'|': {},
	':': {},
	'.': {},
}

func isPunct(s rune) bool {
	_, k := puncts[s]
	return k
}

func (t *lexScan) Next() (lexer.Token, error) {
	typ := t.scanner.Scan()
	text := t.scanner.TokenText()
	pos := lexer.Position(t.scanner.Position)
	pos.Filename = t.file
	if t.err != nil {
		return lexer.Token{}, t.err
	}

	if !isPunct(typ) {
		return lexer.Token{
			Type:  lexer.TokenType(typ),
			Value: text,
			Pos:   pos,
		}, nil
	}

	typ = operator

	for isPunct(t.scanner.Peek()) {
		t.scanner.Scan()
		text += t.scanner.TokenText()

		if t.err != nil {
			return lexer.Token{}, t.err
		}
	}

	return lexer.Token{
		Type:  lexer.TokenType(operator),
		Value: text,
		Pos:   pos,
	}, nil
}
