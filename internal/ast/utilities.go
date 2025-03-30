package ast

func isLeftHandSideExpressionKind(kind Kind) bool {
	switch kind {
	case KindNumericLiteral:
		return true
	}
	return false
}

func IsLeftHandSideExpression(node *Node) bool {
	return isLeftHandSideExpressionKind(node.Kind)
}
func IsPrologueDirective(node *Node) bool {
	return node.Kind == KindExpressionStatement &&
		node.AsExpressionStatement().Expression.Kind == KindStringLiteral
}

func PositionIsSynthesized(pos int) bool {
	return pos < 0
}
