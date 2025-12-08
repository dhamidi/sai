package java

import (
	"strings"

	"github.com/dhamidi/javalyzer/classfile"
)

type Type struct {
	Name       string
	ArrayDepth int
}

func (t Type) String() string {
	var sb strings.Builder
	sb.WriteString(t.Name)
	for i := 0; i < t.ArrayDepth; i++ {
		sb.WriteString("[]")
	}
	return sb.String()
}

func (t Type) IsPrimitive() bool {
	if t.ArrayDepth > 0 {
		return false
	}
	switch t.Name {
	case "boolean", "byte", "char", "short", "int", "long", "float", "double":
		return true
	}
	return false
}

func (t Type) IsArray() bool {
	return t.ArrayDepth > 0
}

func (t Type) IsVoid() bool {
	return t.Name == "void" && t.ArrayDepth == 0
}

func (t Type) ElementType() Type {
	if t.ArrayDepth == 0 {
		return t
	}
	return Type{Name: t.Name, ArrayDepth: t.ArrayDepth - 1}
}

func typeFromFieldType(ft *classfile.FieldType) Type {
	if ft == nil {
		return Type{Name: "void"}
	}
	name := ft.BaseType
	if name == "" {
		name = classfile.InternalToSourceName(ft.ClassName)
	}
	return Type{
		Name:       name,
		ArrayDepth: ft.ArrayDepth,
	}
}
