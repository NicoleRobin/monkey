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