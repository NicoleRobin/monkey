package object

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOLLEAN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
