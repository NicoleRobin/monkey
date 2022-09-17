package repl

import (
	"bufio"
	"fmt"
	"github.com/nicolerobin/log"
	"github.com/nicolerobin/monkey/lexer"
	"github.com/nicolerobin/monkey/parser"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		_, err := fmt.Fprintf(out, ">> ")
		if err != nil {
			log.Info("fmt.Fprintf() failed, err:%s", err)
			break
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		/*
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Fprintf(out, "%+v\n", tok)
			}

		*/

		p := parser.NewParser(l)
		program := p.ParseProgram()
		for i, stmt := range program.Statements {
			_, err := fmt.Fprintf(out, "statement %d:%s\n", i, stmt.TokenLiteral())
			if err != nil {
				log.Error("fmt.Fprintf() failed, err:%s", err)
			}
		}
	}
}
