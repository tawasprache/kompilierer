package ast

type BinaryOperator int

const (
	BinOpCus BinaryOperator = iota
	BinOpAdd
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpPow
	BinOpMod

	BinOpGleich
	BinOpNichtGleich
	BinOpWeniger
	BinOpWenigerGleich
	BinOpGrößer
	BinOpGrößerGleich
)

type opInfo struct {
	Enum             BinaryOperator
	NonAssociative   bool
	RightAssociative bool
	Priority         int
}

var info = map[string]opInfo{
	"^": {Enum: BinOpPow, Priority: 8, RightAssociative: true},

	"*": {Enum: BinOpMul, Priority: 7},
	"/": {Enum: BinOpDiv, Priority: 7},
	"%": {Enum: BinOpMod, Priority: 7},

	"+": {Enum: BinOpAdd, Priority: 6},
	"-": {Enum: BinOpSub, Priority: 6},

	"==": {Enum: BinOpGleich, Priority: 4, NonAssociative: true},
	"!=": {Enum: BinOpNichtGleich, Priority: 4, NonAssociative: true},
	"<":  {Enum: BinOpWeniger, Priority: 4, NonAssociative: true},
	">":  {Enum: BinOpWenigerGleich, Priority: 4, NonAssociative: true},
	"<=": {Enum: BinOpGrößer, Priority: 4, NonAssociative: true},
	">=": {Enum: BinOpGrößerGleich, Priority: 4, NonAssociative: true},
}
