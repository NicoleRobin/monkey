package code

import "fmt"

// Opcode 操作码
type Opcode byte

const (
	OpConstant Opcode = iota // 引用常量
)

// Definition 操作指令定义
type Definition struct {
	Name          string // 操作码的可读名称
	OperandWidths []int  // 每个操作数占用的字节数
}

// definitions 操作指令和其名字及操作数个数的映射关系
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup 根据操作码查询对应的操作指令定义
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}