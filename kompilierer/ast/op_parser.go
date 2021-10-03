package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func (e *Expression) Parse(lex *lexer.PeekingLexer) (err error) {
	defer func() {
		switch k := recover().(type) {
		case error:
			err = k
		case nil:
			return
		default:
			panic(k)
		}
	}()
	*e = *parseExpr(lex, 0)
	return nil
}

func parseExpr(lex *lexer.PeekingLexer, minPrec int) *Expression {
	lhs := parseAtom(lex)
	for {
		tok := peek(lex)
		if tok.EOF() || !isOp(tok) || info[tok.Value].Priority < minPrec {
			break
		}
		op := tok.Value
		nextMinPrec := info[op].Priority
		if !info[op].RightAssociative {
			nextMinPrec++
		}
		lex.Next()
		rhs := parseExpr(lex, nextMinPrec)
		lhs = parseOp(op, lhs, rhs)

		tok2 := peek(lex)
		if isOp(tok2) && (info[op].NonAssociative || info[tok2.Value].NonAssociative) {
			panic(errorf(tok.Pos, "Ich weiß nicht wie man (%s) und (%s) gruppiert. Fügen Sie runde Klammern hinzu!", tok.Value, tok2.Value))
		}
	}
	return lhs
}

func parseAtom(lex *lexer.PeekingLexer) *Expression {
	tok := peek(lex)
	if tok.Type == '(' {
		lex.Next()
		val := parseExpr(lex, 1)
		if peek(lex).Value != ")" {
			panic("unmatched (")
		}
		lex.Next()
		return val
	} else if tok.EOF() {
		panic("unexpected EOF")
	} else if isOp(tok) {
		panic("expected a terminal not " + tok.String())
	} else {
		v := Terminal{}
		err := terminal.ParseFromLexer(lex, &v, participle.AllowTrailing(true))
		if err != nil {
			panic(err)
		}
		return &Expression{Terminal: &v}
	}
}
func isOp(t lexer.Token) bool {
	return t.Type == operator
}

func peek(lex *lexer.PeekingLexer) lexer.Token {
	tok, err := lex.Peek(0)
	if err != nil {
		panic("??")
	}
	return tok
}

func parseOp(op string, lhs *Expression, rhs *Expression) *Expression {
	k := info[op].Enum
	return &Expression{
		Links:  lhs,
		Op:     &k,
		Rechts: rhs,
	}
}
