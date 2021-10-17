package ast

import (
	"errors"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func (e *Expression) Parse(lex *lexer.PeekingLexer) (err error) {
	expr, err := parseExpr(lex, 0)
	if err != nil {
		return err
	}
	*e = *expr
	return nil
}

func parseExpr(lex *lexer.PeekingLexer, minPrec int) (re *Expression, rerr error) {
	t, feh := lex.Peek(0)
	if feh != nil {
		return nil, feh
	}
	lhs, feh := parseAtom(lex)
	if feh != nil {
		return nil, feh
	}
	defer func() {
		if re != nil {
			re.Pos = t.Pos

			t, _ := lex.RawPeek(0)
			re.EndPos = t.Pos
		}
	}()
	for {
		tok, feh := peek(lex)
		if feh != nil {
			return nil, feh
		}
		if tok.Value == "." {
			tok2, feh := peek2(lex)
			if feh != nil {
				return nil, feh
			}

			if tok2.Type == '(' {
				_, feh := lex.Next()
				if feh != nil {
					panic(feh)
				}
				v := Argumentleiste{}
				err := argLeisteParser.ParseFromLexer(lex, &v, participle.AllowTrailing(true))
				if err != nil {
					return nil, err
				}

				return &Expression{
					FunktionErsteKlasseAufruf: &FunktionErsteKlasseAufruf{
						Funktion:   *lhs,
						Argumenten: v,
					},
				}, nil
			}
		}
		if tok.EOF() || !isOp(tok) || info[tok.Value] == nil || info[tok.Value].Priority < minPrec {
			break
		}
		op := tok.Value
		nextMinPrec := info[op].Priority
		if !info[op].RightAssociative {
			nextMinPrec++
		}
		lex.Next()
		rhs, feh := parseExpr(lex, nextMinPrec)
		if feh != nil {
			return nil, feh
		}
		lhs = parseOp(op, lhs, rhs)

		tok2, feh := peek(lex)
		if feh != nil {
			return nil, feh
		}
		if isOp(tok2) && info[op] != nil && info[tok2.Value] != nil && (info[op].NonAssociative || info[tok2.Value].NonAssociative) {
			return nil, errorf(tok.Pos, "Ich weiß nicht wie man (%s) und (%s) gruppiert. Fügen Sie runde Klammern hinzu!", tok.Value, tok2.Value)
		}
	}
	return lhs, nil
}

func parseAtom(lex *lexer.PeekingLexer) (*Expression, error) {
	tok, feh := peek(lex)
	if feh != nil {
		return nil, feh
	}
	if tok.Type == '(' {
		lex.Next()
		val, feh := parseExpr(lex, 1)
		if feh != nil {
			return nil, feh
		}
		peeked, feh := peek(lex)
		if feh != nil {
			return nil, feh
		}
		if peeked.Value != ")" {
			return nil, errors.New("unmatched (")
		}
		lex.Next()
		return val, nil
	} else if tok.EOF() {
		return nil, errors.New("unexpected EOF")
	} else if isOp(tok) {
		return nil, errors.New("expected a terminal not " + tok.String())
	} else {
		v := Terminal{}
		err := TerminalParser.ParseFromLexer(lex, &v, participle.AllowTrailing(true))
		if err != nil {
			return nil, err
		}
		return &Expression{Terminal: &v}, nil
	}
}
func isOp(t lexer.Token) bool {
	return t.Type == operator
}

func peek(lex *lexer.PeekingLexer) (lexer.Token, error) {
	tok, err := lex.Peek(0)
	if err != nil {
		return lexer.Token{}, err
	}
	return tok, nil
}

func peek2(lex *lexer.PeekingLexer) (lexer.Token, error) {
	tok, err := lex.Peek(1)
	if err != nil {
		return lexer.Token{}, err
	}
	return tok, nil
}

func parseOp(op string, lhs *Expression, rhs *Expression) *Expression {
	k := info[op].Enum
	return &Expression{
		Links:  lhs,
		Op:     &k,
		Rechts: rhs,
	}
}
