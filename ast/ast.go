package ast

// Node base node
type Node interface {
	TokenLiteral() string
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

// Program root node
type Program struct {
	Statements []Statement
}

// TokenLiteral return token literal
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
