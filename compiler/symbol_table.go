package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol 符号
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable 符号表
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (st *SymbolTable) Define(name string) Symbol {
	sym := Symbol{
		Name:  name,
		Index: st.numDefinitions,
		Scope: GlobalScope,
	}
	st.store[name] = sym
	st.numDefinitions++
	return sym
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	return sym, ok
}
