# monkey
an interpreter written in go

git tag v1.0：complete monkey interpreter

# directory structure
token/ : lexer token code  
lexer/ : lex analyse code  
parser/ : parser code  
ast/ : Abstract syntax tree code  
object/ : monkey object system  
evaluator/ : evaluator code  
repl/ : read evaluate print loop code  
code/ : instruction code  
compiler/ : compiler code, traverse ast and generate instructions   
vm/ : monkey instruction virtual machine, read instructions and execute   

# 笔记
## 第五章：追踪名称
1、符号表是解释器和编译器中用于将标识符与信息相关联的数据结构。  
它主要负责两件事：1. 将全局范围内的标识符与特定的数字相关联；2. 获取已与给定标识符相关联的数字。  
2、通过新增符号表来支持变量的定义与引用  

## 第六章：字符串、数组和哈希表
1、字符串是不变的，因此可以当做常量来实现，类似与整数字面量一样存储在常量池中。  
2、数组和整数字面量或字符串字面量不同，它的值在运行时才能计算得到，因此数组需要在运行时由虚拟机来构建，
实现方式是：定义一个操作指令OpArray，当编译时遇到数组时，先编译其所有元素，然后发出OpArray指令并且操作数设置为数组的长度N，
当虚拟机运行时遇到OpArray指令时，它从栈中弹出N个元素，构建出数组，并将其压栈。  

## 第七章：一个简单的函数
1、通过作用域实现避免编译函数体时影响其他栈中的指令序列，每开始编译一个函数体前先创建一个作用域，编译完取出函数的指令序列并将其作为一个常量添加到常量池中。
同时销毁刚刚使用的作用域。
2、函数调用通过一个保存函数帧的栈来实现不同函数调用不会互相影响。