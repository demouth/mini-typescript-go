package ast

type Kind int16

const (
	KindUnknown Kind = iota
	KindEndOfFile
	KindNumericLiteral
	KindStringLiteral
	KindCommaToken
	KindSemicolonToken
	KindExpressionStatement
	KindSourceFile
)
