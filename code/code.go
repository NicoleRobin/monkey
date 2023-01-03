package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/nicolerobin/log"
)

type Instructions []byte

// String 指令序列的反汇编函数，返回可读的指令序列
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			_, err := fmt.Fprintf(&out, "ERROR: %s\n", err)
			if err != nil {
				log.Error("fmt.Fprintf() failed, error:%s", err)
			}
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		_, err = fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		if err != nil {
			log.Error("fmt.Fprintf() failed, error:%s", err)
		}
		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match definied %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Make 编码，返回指定操作码和操作数对应的字节码序列
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// 计算指令长度
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instruction[offset] = uint8(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

// ReadOperands 解码，返回
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}

	return operands, offset
}

// ReadUint8 读取单字节
func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

// ReadUint16 以大端序的方式读取两个字节数据
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
