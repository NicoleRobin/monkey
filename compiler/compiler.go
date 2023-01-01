package compiler

import (
	"fmt"
	"github.com/nicolerobin/monkey/ast"
	"github.com/nicolerobin/monkey/code"
	"github.com/nicolerobin/monkey/object"
)

type Compiler struct {
	instructions code.Instructions // 指令序列
	constants    []object.Object   // 常量池
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
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
	case *ast.InfixExpression:
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
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		// 转换为object.Integer对象，并将该对象转换为指令添加到指令序列中
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
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

// emit 生成指令并将指令添加到指令序列中
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

// addInstruction 将指令添加到指令序列中
func (c *Compiler) addInstruction(ins code.Instructions) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
