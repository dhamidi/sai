package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/dhamidi/sai/java"
)

type LineEncoder struct {
	w     io.Writer
	class *java.Class
}

func NewLineEncoder(w io.Writer) *LineEncoder {
	return &LineEncoder{w: w}
}

func (e *LineEncoder) Encode(class *java.Class) error {
	e.class = class
	text, err := e.MarshalText()
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *LineEncoder) MarshalText() ([]byte, error) {
	var sb strings.Builder
	c := e.class

	fmt.Fprintf(&sb, "%s\t%s\t%s\n", e.classKind(), c.Name(), e.classModifiersStr())

	for _, f := range c.Fields() {
		fmt.Fprintf(&sb, "field\t%s\t%s\t%s\t%s\n",
			f.Name(),
			f.Type().String(),
			f.Visibility(),
			e.fieldModifiersStr(f),
		)
	}

	for _, m := range c.Methods() {
		fmt.Fprintf(&sb, "method\t%s\t%s\t%s\t%s\t%s\n",
			m.Name(),
			m.ReturnType().String(),
			e.parametersStr(m.Parameters()),
			m.Visibility(),
			e.methodModifiersStr(m),
		)
	}

	return []byte(sb.String()), nil
}

func (e *LineEncoder) classKind() string {
	c := e.class
	switch {
	case c.IsAnnotation():
		return "annotation"
	case c.IsEnum():
		return "enum"
	case c.IsInterface():
		return "interface"
	case c.IsModule():
		return "module"
	case c.IsRecord():
		return "record"
	default:
		return "class"
	}
}

func (e *LineEncoder) classModifiersStr() string {
	c := e.class
	var mods []string
	mods = append(mods, c.Visibility())
	if c.IsFinal() {
		mods = append(mods, "final")
	}
	if c.IsAbstract() {
		mods = append(mods, "abstract")
	}
	if c.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	if c.IsSealed() {
		mods = append(mods, "sealed")
	}
	return strings.Join(mods, ",")
}

func (e *LineEncoder) fieldModifiersStr(f java.Field) string {
	var mods []string
	if f.IsStatic() {
		mods = append(mods, "static")
	}
	if f.IsFinal() {
		mods = append(mods, "final")
	}
	if f.IsVolatile() {
		mods = append(mods, "volatile")
	}
	if f.IsTransient() {
		mods = append(mods, "transient")
	}
	if f.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	if f.IsEnum() {
		mods = append(mods, "enum")
	}
	if len(mods) == 0 {
		return "-"
	}
	return strings.Join(mods, ",")
}

func (e *LineEncoder) methodModifiersStr(m java.Method) string {
	var mods []string
	if m.IsStatic() {
		mods = append(mods, "static")
	}
	if m.IsFinal() {
		mods = append(mods, "final")
	}
	if m.IsAbstract() {
		mods = append(mods, "abstract")
	}
	if m.IsSynchronized() {
		mods = append(mods, "synchronized")
	}
	if m.IsNative() {
		mods = append(mods, "native")
	}
	if m.IsBridge() {
		mods = append(mods, "bridge")
	}
	if m.IsVarargs() {
		mods = append(mods, "varargs")
	}
	if m.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	if len(mods) == 0 {
		return "-"
	}
	return strings.Join(mods, ",")
}

func (e *LineEncoder) parametersStr(params []java.Parameter) string {
	if len(params) == 0 {
		return "-"
	}
	var parts []string
	for _, p := range params {
		parts = append(parts, p.Type.String())
	}
	return strings.Join(parts, ",")
}

type LineModelEncoder struct {
	w     io.Writer
	model *java.ClassModel
}

func NewLineModelEncoder(w io.Writer) *LineModelEncoder {
	return &LineModelEncoder{w: w}
}

func (e *LineModelEncoder) Encode(model *java.ClassModel) error {
	e.model = model
	text, err := e.MarshalText()
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *LineModelEncoder) MarshalText() ([]byte, error) {
	var sb strings.Builder
	m := e.model

	fmt.Fprintf(&sb, "%s\t%s\t%s\n", m.Kind, m.Name, e.classModifiersStr())

	for _, f := range m.Fields {
		fmt.Fprintf(&sb, "field\t%s\t%s\t%s\t%s\n",
			f.Name,
			typeModelStr(f.Type),
			f.Visibility,
			e.fieldModifiersStr(f),
		)
	}

	for _, method := range m.Methods {
		fmt.Fprintf(&sb, "method\t%s\t%s\t%s\t%s\t%s\n",
			method.Name,
			typeModelStr(method.ReturnType),
			e.parametersStr(method.Parameters),
			method.Visibility,
			e.methodModifiersStr(method),
		)
	}

	for _, rc := range m.RecordComponents {
		fmt.Fprintf(&sb, "component\t%s\t%s\n",
			rc.Name,
			typeModelStr(rc.Type),
		)
	}

	return []byte(sb.String()), nil
}

func (e *LineModelEncoder) classModifiersStr() string {
	m := e.model
	var mods []string
	mods = append(mods, string(m.Visibility))
	if m.IsFinal {
		mods = append(mods, "final")
	}
	if m.IsAbstract {
		mods = append(mods, "abstract")
	}
	if m.IsSynthetic {
		mods = append(mods, "synthetic")
	}
	if m.IsSealed {
		mods = append(mods, "sealed")
	}
	return strings.Join(mods, ",")
}

func (e *LineModelEncoder) fieldModifiersStr(f java.FieldModel) string {
	var mods []string
	if f.IsStatic {
		mods = append(mods, "static")
	}
	if f.IsFinal {
		mods = append(mods, "final")
	}
	if f.IsVolatile {
		mods = append(mods, "volatile")
	}
	if f.IsTransient {
		mods = append(mods, "transient")
	}
	if f.IsSynthetic {
		mods = append(mods, "synthetic")
	}
	if f.IsEnum {
		mods = append(mods, "enum")
	}
	if len(mods) == 0 {
		return "-"
	}
	return strings.Join(mods, ",")
}

func (e *LineModelEncoder) methodModifiersStr(m java.MethodModel) string {
	var mods []string
	if m.IsStatic {
		mods = append(mods, "static")
	}
	if m.IsFinal {
		mods = append(mods, "final")
	}
	if m.IsAbstract {
		mods = append(mods, "abstract")
	}
	if m.IsSynchronized {
		mods = append(mods, "synchronized")
	}
	if m.IsNative {
		mods = append(mods, "native")
	}
	if m.IsBridge {
		mods = append(mods, "bridge")
	}
	if m.IsVarargs {
		mods = append(mods, "varargs")
	}
	if m.IsSynthetic {
		mods = append(mods, "synthetic")
	}
	if m.IsDefault {
		mods = append(mods, "default")
	}
	if len(mods) == 0 {
		return "-"
	}
	return strings.Join(mods, ",")
}

func (e *LineModelEncoder) parametersStr(params []java.ParameterModel) string {
	if len(params) == 0 {
		return "-"
	}
	var parts []string
	for _, p := range params {
		parts = append(parts, typeModelStr(p.Type))
	}
	return strings.Join(parts, ",")
}

func typeModelStr(t java.TypeModel) string {
	s := t.Name
	for i := 0; i < t.ArrayDepth; i++ {
		s += "[]"
	}
	return s
}
