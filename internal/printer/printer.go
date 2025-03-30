package printer

import (
	"fmt"

	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/core"
	"github.com/demouth/mini-typescript-go/internal/scanner"
)

type ListFormat int

const (
	LFMultiLine ListFormat = 1 << 0 // Prints the list on multiple lines.

	LFIndented ListFormat = 1 << 7 // The list should be indented.
)

type Printer struct {
	currentSourceFile  *ast.SourceFile
	nextListElementPos int
	writer             EmitTextWriter
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Write(node *ast.Node, sourceFile *ast.SourceFile, writer EmitTextWriter) {
	savedCurrentSourceFile := p.currentSourceFile
	savedWriter := p.writer

	p.currentSourceFile = sourceFile
	p.writer = writer
	p.writer.Clear()

	switch node.Kind {
	case ast.KindSourceFile:
		p.emitSourceFile(node.AsSourceFile())
	}

	p.writer = savedWriter
	p.currentSourceFile = savedCurrentSourceFile
}

func (p *Printer) emitSourceFile(node *ast.SourceFile) {
	savedCurrentSourceFile := p.currentSourceFile
	p.currentSourceFile = node

	p.writeLine()

	_ = core.FindIndex(node.Statements.Nodes, isNotPrologueDirective)

	p.emitListRange(
		(*Printer).emitStatement,
		node.AsNode(),
		node.Statements,
		1,
		0,
		-1,
	)

	p.currentSourceFile = savedCurrentSourceFile
}

func (p *Printer) writeLine() {
	p.writer.WriteLine()
}
func (p *Printer) emitListRange(
	emit func(p *Printer, node *ast.Node),
	parentNode *ast.Node,
	children *ast.NodeList,
	format ListFormat,
	start int,
	count int,
) {
	length := len(children.Nodes)

	if start < 0 {
		start = 0
	}

	if count < 0 {
		count = length - start
	}

	end := start + count
	if end > length {
		end = length
	}

	p.emitListItems(
		emit,
		parentNode,
		children.Nodes[start:end],
		format,
		children.HasTrailingComma(),
		children.Loc,
	)
}
func (p *Printer) emitStatement(node *ast.Statement) {
	switch node.Kind {
	case ast.KindExpressionStatement:
		p.emitExpressionStatement(node.AsExpressionStatement())
	default:
		panic(fmt.Sprintf("unhandled statement: %v", node.Kind))
	}
}

func (p *Printer) emitListItems(
	emit func(p *Printer, node *ast.Node),
	parentNode *ast.Node,
	children []*ast.Node,
	format ListFormat,
	hasTrailingComma bool,
	childrenTextRange core.TextRange,
) {
	leadingLineTerminatorCount := 0
	if len(children) > 0 {
		leadingLineTerminatorCount = p.getLeadingLineTerminatorCount(parentNode, children[0], format)
	}
	if leadingLineTerminatorCount > 0 {
		for range leadingLineTerminatorCount {
			p.writeLine()
		}
	}

	_ = greatestEnd(-1, parentNode)

	for _, child := range children {
		p.nextListElementPos = child.Pos()
		emit(p, child)
	}

	closingLineTerminatorCount := 1
	if closingLineTerminatorCount > 0 {
		for range closingLineTerminatorCount {
			p.writeLine()
		}
	}
}

func (p *Printer) getLeadingLineTerminatorCount(parentNode *ast.Node, firstChild *ast.Node, format ListFormat) int {
	return core.IfElse(format&LFMultiLine != 0, 1, 0)
}

func (p *Printer) emitExpressionStatement(node *ast.ExpressionStatement) {
	p.enterNode(node.AsNode())
	p.emitExpression(node.Expression, ast.OperatorPrecedenceComma)

	// if p.currentSourceFile == nil {
	p.writeTrailingSemicolon()
	// }
	p.exitNode(node.AsNode())
}

func (p *Printer) enterNode(node *ast.Node) {
	// !!!
}
func (p *Printer) exitNode(node *ast.Node) {
	// !!!
}

func (p *Printer) emitExpression(node *ast.Expression, precedence ast.OperatorPrecedence) {
	switch node.Kind {
	case ast.KindNumericLiteral:
		p.emitNumericLiteral(node.AsNumericLiteral())

	}
}
func (p *Printer) emitNumericLiteral(node *ast.NumericLiteral) {
	p.enterNode(node.AsNode())
	p.emitLiteral(node.AsNode(), getLiteralTextFlagsNone)
	p.exitNode(node.AsNode())
}

func (p *Printer) emitLiteral(node *ast.LiteralLikeNode, flags getLiteralTextFlags) {
	text := p.getLiteralTextOfNode(node, nil, flags)
	p.writer.WriteStringLiteral(text)
}

func (p *Printer) getLiteralTextOfNode(node *ast.LiteralLikeNode, sourceFile *ast.SourceFile, flags getLiteralTextFlags) string {
	return getLiteralText(node, p.currentSourceFile, flags)
}
func getLiteralText(node *ast.LiteralLikeNode, sourceFile *ast.SourceFile, flags getLiteralTextFlags) string {
	return scanner.GetSourceTextOfNodeFromSourceFile(sourceFile, node, false)
}
func (p *Printer) writeTrailingSemicolon() {
	p.writer.WriteTrailingSemicolon(";")
}
