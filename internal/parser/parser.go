package parser

import (
	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/core"
	"github.com/demouth/mini-typescript-go/internal/scanner"
)

type ParsingContext int

const (
	PCSourceElements ParsingContext = iota // Elements in source file
)

type Parser struct {
	scanner *scanner.Scanner
	factory ast.NodeFactory

	fileName   string
	sourceText string

	token ast.Kind

	nodeSlicePool core.Pool[*ast.Node]
}

func ParseSourceFile(fileName string, sourceText string) *ast.SourceFile {
	var p Parser
	p.initializeState(fileName, sourceText)
	p.nextToken()
	return p.parseSourceFileWorker()
}

func (p *Parser) initializeState(fileName string, sourceText string) {
	p.scanner = scanner.NewScanner()
	p.fileName = fileName
	p.sourceText = sourceText
	p.scanner.SetText(p.sourceText)
}

func (p *Parser) nextToken() ast.Kind {
	p.token = p.scanner.Scan()
	return p.token
}

func (p *Parser) nodePos() int {
	return p.scanner.TokenFullStart()
}

func (p *Parser) parseSourceFileWorker() *ast.SourceFile {

	pos := p.nodePos()
	statements := p.parseListIndex(PCSourceElements, (*Parser).parseToplevelStatement)
	eof := p.parseTokenNode()
	if eof.Kind != ast.KindEndOfFile {
		panic("Expected end of file token from scanner.")
	}
	node := p.factory.NewSourceFile(p.sourceText, p.fileName, statements)
	p.finishNode(node, pos)
	result := node.AsSourceFile()
	return result
}
func (p *Parser) parseTokenNode() *ast.Node {
	pos := p.nodePos()
	kind := p.token
	p.nextToken()
	result := p.factory.NewToken(kind)
	p.finishNode(result, pos)
	return result
}
func (p *Parser) parseListIndex(kind ParsingContext, parseElement func(p *Parser, index int) *ast.Node) *ast.NodeList {
	pos := p.nodePos()

	list := make([]*ast.Node, 0, 16)
	for i := 0; !p.isListTerminator(kind); i++ {
		if p.isListElement(kind) {
			list = append(list, parseElement(p, i))
			continue
		}
	}
	slice := p.nodeSlicePool.NewSlice(len(list))
	copy(slice, list)
	return p.newNodeList(core.NewTextRange(pos, p.nodePos()), slice)
}
func (p *Parser) newNodeList(loc core.TextRange, nodes []*ast.Node) *ast.NodeList {
	list := p.factory.NewNodeList(nodes)
	list.Loc = loc
	return list
}
func (p *Parser) parseToplevelStatement(i int) *ast.Node {
	statement := p.parseStatement()
	return statement
}
func (p *Parser) isListTerminator(kind ParsingContext) bool {
	if p.token == ast.KindEndOfFile {
		return true
	}
	return false
}
func (p *Parser) isListElement(parsingContext ParsingContext) bool {
	switch parsingContext {
	case PCSourceElements:
		return true
	}
	panic("Unhandled case in isListElement")
}

func (p *Parser) parseStatement() *ast.Statement {
	switch p.token {
	}
	return p.parseExpressionOrLabeledStatement()
}

func (p *Parser) parseExpressionOrLabeledStatement() *ast.Statement {
	pos := p.nodePos()
	expression := p.parseExpression()
	p.parseSemicolon()
	result := p.factory.NewExpressionStatement(expression)
	p.finishNode(result, pos)
	return result
}

func (p *Parser) parseSemicolon() bool {
	return p.tryParseSemicolon()
}
func (p *Parser) canParseSemicolon() bool {
	return p.token == ast.KindEndOfFile
}
func (p *Parser) tryParseSemicolon() bool {
	if !p.canParseSemicolon() {
		return false
	}
	if p.token == ast.KindSemicolonToken {
		p.nextToken()
	}
	return true
}
func (p *Parser) parseExpression() *ast.Expression {
	expr := p.parseAssignmentExpressionOrHigher()
	return expr
}

func (p *Parser) parseAssignmentExpressionOrHigher() *ast.Expression {
	return p.parseAssignmentExpressionOrHigherWorker(true)
}

func (p *Parser) parseAssignmentExpressionOrHigherWorker(allowReturnTypeInArrowFunction bool) *ast.Node {
	pos := p.nodePos()
	expr := p.parseBinaryExpressionOrHigher(ast.OperatorPrecedenceLowest)
	return p.parseConditionalExpressionRest(expr, pos, allowReturnTypeInArrowFunction)
}

func (p *Parser) parseBinaryExpressionOrHigher(precedence ast.OperatorPrecedence) *ast.Expression {
	pos := p.nodePos()
	leftOperand := p.parseUnaryExpressionOrHigher()
	return p.parseBinaryExpressionRest(precedence, leftOperand, pos)
}

func (p *Parser) parseUnaryExpressionOrHigher() *ast.Expression {
	if p.isUpdateExpression() {
		updateExpression := p.parseUpdateExpression()
		return updateExpression
	}
	return nil // fixme
}
func (p *Parser) parseBinaryExpressionRest(precedence ast.OperatorPrecedence, leftOperand *ast.Expression, pos int) *ast.Expression {
	for {
		newPrecedence := ast.GetBinaryOperatorPrecedence(p.token)

		var consumeCurrentOperator bool
		consumeCurrentOperator = newPrecedence > precedence
		if !consumeCurrentOperator {
			break
		}
	}
	return leftOperand
}
func (p *Parser) isUpdateExpression() bool {
	return true
}

func (p *Parser) parseUpdateExpression() *ast.Expression {
	expression := p.parseLeftHandSideExpressionOrHigher()
	return expression
}

func (p *Parser) parseLeftHandSideExpressionOrHigher() *ast.Expression {
	pos := p.nodePos()
	expression := p.parseMemberExpressionOrHigher()
	return p.parseCallExpressionRest(pos, expression)
}

func (p *Parser) parseMemberExpressionOrHigher() *ast.Expression {
	pos := p.nodePos()
	expression := p.parsePrimaryExpression()
	return p.parseMemberExpressionRest(pos, expression, true)
}

func (p *Parser) parsePrimaryExpression() *ast.Expression {
	// switch p.token {
	// case ast.KindNumericLiteral:
	return p.parseLiteralExpression()
	// }
}

func (p *Parser) parseLiteralExpression() *ast.Expression {
	pos := p.nodePos()
	text := p.scanner.TokenValue()
	var result *ast.Node
	switch p.token {
	case ast.KindNumericLiteral:
		result = p.factory.NewNumericLiteral(text)
	default:
		panic("Unhandled case in parseLiteralExpression")
	}
	p.nextToken()
	p.finishNode(result, pos)
	return result
}
func (p *Parser) parseMemberExpressionRest(pos int, expression *ast.Expression, allowOptionalChain bool) *ast.Expression {
	return expression
}
func (p *Parser) parseCallExpressionRest(pos int, expression *ast.Expression) *ast.Expression {
	return expression
}
func (p *Parser) parseConditionalExpressionRest(leftOperand *ast.Expression, pos int, allowReturnTypeInArrowFunction bool) *ast.Expression {
	return leftOperand
}
func (p *Parser) finishNode(node *ast.Node, pos int) {
	p.finishNodeWithEnd(node, pos, p.nodePos())
}
func (p *Parser) finishNodeWithEnd(node *ast.Node, pos int, end int) {
	node.Loc = core.NewTextRange(pos, end)
	// node.Flags |= p.contextFlags
}
