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
2、