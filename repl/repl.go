package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nicolerobin/log"
	"github.com/nicolerobin/monkey/compiler"
	"github.com/nicolerobin/monkey/lexer"
	"github.com/nicolerobin/monkey/parser"
	"github.com/nicolerobin/monkey/vm"
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
			_, err := fmt.Fprintf(out, "Woops! Compilation failed, error: %s\n", err)
			if err != nil {
				log.Error("fmt.Fprintf failed, error:%s", err)
			}
			continue
		}

		machine := vm.NewVm(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			_, err := fmt.Fprintf(out, "Woops! Executing bytecode failed, error: %s\n", err)
			if err != nil {
				log.Error("fmt.Fprintf failed, error:%s", err)
			}
			continue
		}

		// stackTop := machine.StackTop()
		stackTop := machine.LastPoppedStackElem()
		_, err = io.WriteString(out, stackTop.Inspect())
		if err != nil {
			log.Error("io.WriteString failed, error:%s", err)
		}
		_, err = io.WriteString(out, "\n")
		if err != nil {
			log.Error("io.WriteString failed, error:%s", err)
		}
	}
}

/*
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
			_, err := io.WriteString(out, evaluated.Inspect())
			if err != nil {
				log.Error("io.WriteString failed, error:%s", err)
			}
			_, err = io.WriteString(out, "\n")
			if err != nil {
				log.Error("io.WriteString failed, error:%s", err)
			}
		}
	}
}
*/

func printParseErrors(out io.Writer, errors []string) {
	_, err := fmt.Fprintf(out, MONKEY_FACE)
	if err != nil {
		log.Error("fmt.Fprintf failed, error:%s", err)
	}
	_, err = fmt.Fprintf(out, "Woops! We ran into some monkey business here!\n")
	if err != nil {
		log.Error("fmt.Fprintf failed, error:%s", err)
	}
	_, err = fmt.Fprintf(out, " parser errors:\n")
	if err != nil {
		log.Error("fmt.Fprintf failed, error:%s", err)
	}
	for _, msg := range errors {
		_, err = fmt.Fprintf(out, "\t"+msg+"\n")
		if err != nil {
			log.Error("fmt.Fprintf failed, error:%s", err)
		}
	}
}
