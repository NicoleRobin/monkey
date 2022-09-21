package parser

import (
	"testing"

	"github.com/nicolerobin/monkey/ast"
	"github.com/nicolerobin/monkey/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838 383;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	checkPeekError(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statments. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got=%s", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *astLetStatement, got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s', got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func checkPeekError(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range p.errors {
		t.Errorf("parse error %q", msg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 993 322;	
`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkPeekError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		if returnStmt, ok := stmt.(*ast.ReturnStatement); !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got:%T", stmt)
			continue
		} else {
			if returnStmt != nil && returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral() not 'return', got:%q", returnStmt.TokenLiteral())
			}
		}
	}
}
