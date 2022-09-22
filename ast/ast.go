package ast

// Node base node
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement statement node
type Statement interface {
	Node
	statementNode()
}

// Expression expression node
type Expression interface {
	Node
	expressionNode()
}
