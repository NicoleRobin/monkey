package compiler

import (
	"fmt"
	"github.com/nicolerobin/monkey/ast"
	"github.com/nicolerobin/monkey/code"
	"github.com/nicolerobin/monkey/object"
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Compiler struct {
	instructions code.Instructions // 指令序列
	constants    []object.Object   // 常量池

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction

	symbolTable *SymbolTable
}

// NewCompiler 创建Compiler
func NewCompiler() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
		symbolTable:  NewSymbolTable(),
	}
}

// NewWithState 创建Compiler并设置状态存储
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := NewCompiler()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

// Compile 递归遍历AST并生成指令序列
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)
		default:
			return fmt.Errorf("unknown operator:%s", node.Operator)
		}
	case *ast.InfixExpression:
		if node.Operator == "<" {
			// 将'<'转换为'>'
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		// 转换为object.Integer对象，并将该对象转换为指令添加到指令序列中
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.IfExpression:
		// 处理条件condition
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// 设置JumpNotTruthy指令，先使用虚假偏移量
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		// 处理结果consequence
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// 发出带有虚假偏移的OpJump
		jumpPos := c.emit(code.OpJump, 9999)

		// 修正OpJumpNotTruthy指令的跳转位置
		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		// 处理else部分
		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			// 处理alternative
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			// 移除末尾的OpPop
			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}
		// 修正OpJump指令的跳转位置
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.BlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// 为变量创建符号
		symbol := c.symbolTable.Define(node.Name.Value)
		// 将符号的Index作为OpSetGlobal的参数添加到指令序列中
		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.Identifier:
		sym, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable:%s\n", node.Value)
		}
		c.emit(code.OpGetGlobal, sym.Index)
	}
	return nil
}

// Bytecode
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// addConstant 将对象obj添加到常量池中并返回其在常量池中的下标作为引用
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// addInstruction 将指令添加到指令序列中
func (c *Compiler) addInstruction(ins code.Instructions) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

// emit 生成指令并将指令添加到指令序列中
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

// 记录最近一次指令和倒数第二次指令
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{
		Opcode:   op,
		Position: pos,
	}
}

// 判断最后一条指令是否是OpPop
func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

// 移除最后一个OpPop
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:len(c.instructions)-1]
	c.lastInstruction = c.previousInstruction
}

// replaceInstruction 回填操作，修正OpJumpNotTruthy的偏移量
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// changeOperand 修改操作数
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}
