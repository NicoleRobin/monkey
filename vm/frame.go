package vm

import (
	"github.com/nicolerobin/monkey/code"
	"github.com/nicolerobin/monkey/object"
)

// Frame 栈帧，保存与函数执行相关信息的数据结构
type Frame struct {
	fn *object.CompiledFunction // 栈帧引用的已编译函数
	ip int                      // 栈帧的指令指针
}

// NewFrame 创建新的栈帧
func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

// Instructions 返回该栈帧对应的函数体指令序列
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
