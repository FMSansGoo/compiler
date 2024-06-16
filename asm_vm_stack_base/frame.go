package asm_vm_stack_base

// 帧是调用帧或者栈帧的简称，指保存与执行相关的信息的数据结构
// 在物理机上，帧并不独立于栈存在，而是栈的特定部分。它是存储返回地址、当前函数的参数及其局部变量的地方。
// 由于它在栈中，因此帧在函数执行结束后很容易被弹栈
// 在虚拟机上，不必使用栈来存储帧，因为此处不受标准化调用约定或其他真实的内容约束，比如真实的内存地址和位置
// 当然，其实也可以放在栈上存储一些信息，就是会比较繁琐一点
type Frame struct {
	cl          *ClosureObject // fn 指向帧引用的已编译函数
	ip          int            // ip寄存器叫做指令寄存器 instruction pointer
	basePointer int
}

func NewFrame(cl *ClosureObject, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() Instructions {
	return f.cl.Fn.Instructions
}
