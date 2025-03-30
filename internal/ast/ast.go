package ast

import (
	"github.com/demouth/mini-typescript-go/internal/core"
	"github.com/demouth/mini-typescript-go/internal/tspath"
)

type nodeData interface {
	AsNode() *Node
}

type NodeFactory struct {
	tokenPool    core.Pool[Token]
	nodeListPool core.Pool[NodeList]
}

func newNode(kind Kind, data nodeData) *Node {
	n := data.AsNode()
	n.Kind = kind
	n.data = data
	return n
}

type NodeList struct {
	Loc   core.TextRange
	Nodes []*Node
}

func (list *NodeList) End() int { return list.Loc.End() }
func (list *NodeList) HasTrailingComma() bool {
	return len(list.Nodes) > 0 && list.Nodes[len(list.Nodes)-1].End() < list.End()
}

type Node struct {
	Kind  Kind
	Flags NodeFlags
	Loc   core.TextRange
	data  nodeData
}

func (n *Node) Pos() int { return n.Loc.Pos() }
func (n *Node) End() int { return n.Loc.End() }
func (n *Node) AsSourceFile() *SourceFile {
	return n.data.(*SourceFile)
}
func (n *Node) AsExpressionStatement() *ExpressionStatement {
	return n.data.(*ExpressionStatement)
}
func (n *Node) AsNumericLiteral() *NumericLiteral {
	return n.data.(*NumericLiteral)
}

type NodeDefault struct {
	Node
}

type (
	Statement       = Node
	Expression      = Node
	LiteralLikeNode = Node
)

func (node *NodeDefault) AsNode() *Node { return &node.Node }

type NodeBase struct {
	NodeDefault
}
type ExpressionBase struct {
	NodeBase
}

type SourceFile struct {
	NodeBase
	Text       string
	fileName   string
	path       tspath.Path
	Statements *NodeList
}

// fixme: this is a hack
func (node *SourceFile) SetFileName(fileName string) {
	node.fileName = fileName
}
func (node *SourceFile) FileName() string {
	return node.fileName
}
func (node *SourceFile) SetPath(p tspath.Path) {
	node.path = p
}

type LiteralLikeBase struct {
	Text string
}
type Token struct {
	NodeBase
}
type NumericLiteral struct {
	ExpressionBase
	LiteralLikeBase
}
type StatementBase struct {
	NodeBase
}
type ExpressionStatement struct {
	StatementBase
	Expression *Expression // Expression
}

func (f *NodeFactory) NewSourceFile(text string, fileName string, statements *NodeList) *Node {
	data := &SourceFile{}
	data.Text = text
	data.fileName = fileName
	data.Statements = statements
	return newNode(KindSourceFile, data)
}
func (f *NodeFactory) NewToken(kind Kind) *Node {
	return newNode(kind, f.tokenPool.New())
}
func (f *NodeFactory) NewNodeList(nodes []*Node) *NodeList {
	list := f.nodeListPool.New()
	list.Loc = core.UndefinedTextRange()
	list.Nodes = nodes
	return list
}
func (f *NodeFactory) NewNumericLiteral(text string) *Node {
	data := &NumericLiteral{}
	data.Text = text
	return newNode(KindNumericLiteral, data)
}
func (f *NodeFactory) NewExpressionStatement(expression *Expression) *Node {
	data := &ExpressionStatement{}
	data.Expression = expression
	return newNode(KindExpressionStatement, data)
}
