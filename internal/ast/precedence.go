package ast

type OperatorPrecedence int

const (
	OperatorPrecedenceComma OperatorPrecedence = iota

	OperatorPrecedenceLowest = OperatorPrecedenceComma

	OperatorPrecedenceInvalid OperatorPrecedence = -1

	OperatorPrecedencePrimary
)

func GetBinaryOperatorPrecedence(operatorKind Kind) OperatorPrecedence {
	return OperatorPrecedenceInvalid
}

func GetExpressionPrecedence(expression *Expression) OperatorPrecedence {
	return OperatorPrecedencePrimary
}
