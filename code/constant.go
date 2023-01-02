package code

import "fmt"

// Opcode 操作码
type Opcode byte

const (
	OpConstant Opcode = iota // 引用常量指令
	OpPop                    // 弹出栈顶元素指令，用于在每个表达式语句后将其结果从栈中弹出
	OpAdd                    // 加法操作指令
	OpSub                    // 减法操作指令
	OpMul                    // 乘法操作指令
	OpDiv                    // 除法操作指令
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
)

// Definition 操作指令定义
type Definition struct {
	Name          string // 操作码的可读名称
	OperandWidths []int  // 每个操作数占用的字节数
}

// definitions 操作指令和其名字及操作数个数的映射关系
var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpPop:           {"OpPop", []int{}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
}

// Lookup 根据操作码查询对应的操作指令定义
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}
