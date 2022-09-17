package parser

import (
	"fmt"
	"github.com/nicolerobin/monkey/ast"
	"github.com/nicolerobin/monkey/lexer"
	"github.com/nicolerobin/monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parse program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.nextToken()
		return true
	} else {
		p.peekError(tokenType)
		return false
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	identifier := &ast.Identifier{}
	identifier.Token = p.curToken
	return identifier
}

func (p *Parser) parseExpression() ast.Expression {
	if p.curToken.Type == token.INT {
		if p.peekToken.Type == token.PLUS {
			return p.parseOperatorExpression()
		} else if p.peekToken.Type == token.SEMICOLON {
			return p.parseIntegerLiteral()
		}
	} else if p.curToken.Type == token.LPAREN {
		return p.parseGroupedExpression()
	}
	return nil
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	return nil
}

func (p *Parser) parseOperatorExpression() ast.Expression {
	return nil
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	return nil
}
