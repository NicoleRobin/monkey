package ast

import (
	"bytes"

	"github.com/nicolerobin/monkey/token"
)

// Identifier identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) String() string {
	var out bytes.Buffer

	out.WriteString(i.TokenLiteral() + " ")
	out.WriteString(i.Value)
	return out.String()
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

var _ Node = &Identifier{}
