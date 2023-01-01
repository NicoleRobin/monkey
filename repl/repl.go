package repl

import (
	"bufio"
	"fmt"
	"github.com/nicolerobin/monkey/evaluator"
	"github.com/nicolerobin/monkey/vm"
	"io"

	"github.com/nicolerobin/log"
	"github.com/nicolerobin/monkey/compiler"
	"github.com/nicolerobin/monkey/lexer"
	"github.com/nicolerobin/monkey/object"
	"github.com/nicolerobin/monkey/parser"
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

func StartVM(in io.Reader, out io.Writer) {
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

		comp := compiler.NewCompiler()
		err = comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed, error: %s\n", err)
			continue
		}

		machine := vm.NewVm(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed, error: %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
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

		// 解释器求值
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
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
