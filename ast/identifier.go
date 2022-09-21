package ast

import "github.com/nicolerobin/monkey/token"

// Identifier identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
