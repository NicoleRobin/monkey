package vm

import (
	"fmt"
	"github.com/nicolerobin/monkey/code"
	"github.com/nicolerobin/monkey/compiler"
	"github.com/nicolerobin/monkey/object"
)

const (
	StackSize  = 2048  // 栈大小
	GlobalSize = 65536 // 全局符号表大小
	MaxFrame   = 1024  // 帧栈大小
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

type VM struct {
	constants []object.Object // 常量池

	stack []object.Object // 栈
	sp    int             // 始终指向栈中的下一个空闲槽，栈顶的值是stack[sp-1]

	globals []object.Object // 存储全局变量

	frames     []*Frame // 用于保存帧的栈
	frameIndex int      //
}

func NewVm(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn)

	frames := make([]*Frame, MaxFrame)
	frames[0] = mainFrame

	return &VM{
		constants:  bytecode.Constants,
		stack:      make([]object.Object, StackSize),
		sp:         0,
		globals:    make([]object.Object, GlobalSize),
		frames:     frames,
		frameIndex: 1,
	}
}

func NewVmWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := NewVm(bytecode)
	vm.globals = s
	return vm
}

// StackTop 获取栈顶指令
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	// log.Debug("vm.sp:%d", vm.sp)
	return vm.stack[vm.sp]
}

// Run 运行虚拟机主循环：取指令、解码、执行
func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])
		// log.Debug("ip:%d, ins:%+v, op:%d", ip, ins, op)

		switch op {
		case code.OpConstant:
			// 读取到引用指令
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			// 出栈指令
			vm.pop()
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			// 获取绑定到名称的值，并将其存储到globals中
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			// 从globals中取出值并将其入栈
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpArray:
			arrayLen := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-arrayLen, vm.sp)
			vm.sp = vm.sp - arrayLen

			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			hashLen := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-hashLen, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - hashLen

			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.OpCall:
			// 调用函数，创建新的栈帧
			fn, ok := vm.stack[vm.sp-1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non-function")
			}
			frame := NewFrame(fn)
			// log.Debug("frame:%+v", frame)
			vm.pushFrame(frame)
		case code.OpReturnValue:
			// 在函数栈帧中获取返回值
			returnValue := vm.pop()

			// 销毁函数栈帧
			vm.popFrame()
			vm.pop()

			// 将返回值压入栈
			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			vm.popFrame()
			vm.pop()

			err := vm.push(Null)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown opcode:%d", op)
		}
	}
	return nil
}

func isTruthy(condition object.Object) bool {
	switch obj := condition.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation:%s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func nativeBoolToBooleanObject(result bool) object.Object {
	if result {
		return True
	}
	return False
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator:%d", op)
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	if leftType != rightType {
		return fmt.Errorf("leftType:%s and rightType:%s not equal", leftType, rightType)
	}

	switch leftType {
	case object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %d %s",
			leftType, op, rightType)
	}
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64 = 0
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator:%d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{Value: leftValue + rightValue})
}

// push 将obj入栈
func (vm *VM) push(o object.Object) error {
	// log.Debug("vm.sp:%d, o:%+v", vm.sp, o)
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	// log.Debug("vm.sp:%d", vm.sp)
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

// buildArray 构建数组
func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key:%s", key.Type())
		}

		hashPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashPairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported:%s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(left, index object.Object) error {
	arrayObj := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(Null)
	}

	return vm.push(arrayObj.Elements[i])
}

func (vm *VM) executeHashIndex(left, index object.Object) error {
	hashObj := left.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return vm.push(Null)
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}
	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.frameIndex] = f
	vm.frameIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.frameIndex--
	return vm.frames[vm.frameIndex]
}
