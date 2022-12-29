package repl

import (
	"bufio"
	"fmt"
	"github.com/nicolerobin/log"
	"github.com/nicolerobin/monkey/lexer"
	"github.com/nicolerobin/monkey/parser"
	"io"
)

const (
	PROMPT      = ">> "
	MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		_, err := fmt.Fprintf(out, PROMPT)
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
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParseErrors(out, p.Errors())
			continue
		}
		fmt.Fprintf(out, "program:%s\n", program.String())
	}
}

func printParseErrors(out io.Writer, errors []string) {
	fmt.Fprintf(out, MONKEY_FACE)
	fmt.Fprintf(out, "Woops! We ran into some monkey business here!\n")
	fmt.Fprintf(out, " parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintf(out, "\t"+msg+"\n")
	}
}
