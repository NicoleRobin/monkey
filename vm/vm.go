package vm

import (
	"fmt"
	"github.com/nicolerobin/monkey/code"
	"github.com/nicolerobin/monkey/compiler"
	"github.com/nicolerobin/monkey/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // 始终指向栈中的下一个空闲槽，栈顶的值是stack[sp-1]
}

func NewVm(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// StackTop 获取栈顶指令
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// Run 运行虚拟机主循环：取指令、解码、执行
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			// 读取到引用指令
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			// 读取到加法指令
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value
			sum := leftValue + rightValue
			err := vm.push(&object.Integer{Value: sum})
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("")
		}
	}
	return nil
}

// push 将obj入栈
func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}
