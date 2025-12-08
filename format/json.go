package format

import (
	"encoding/json"
	"io"

	"github.com/dhamidi/javalyzer/java"
)

type JSONEncoder struct {
	w     io.Writer
	class *java.Class
}

func NewJSONEncoder(w io.Writer) *JSONEncoder {
	return &JSONEncoder{w: w}
}

func (e *JSONEncoder) Encode(class *java.Class) error {
	e.class = class
	text, err := e.MarshalText()
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *JSONEncoder) MarshalText() ([]byte, error) {
	data := e.buildClassData()
	return json.MarshalIndent(data, "", "  ")
}

type jsonClass struct {
	Name       string       `json:"name"`
	SimpleName string       `json:"simpleName"`
	Package    string       `json:"package"`
	SuperClass string       `json:"superClass,omitempty"`
	Interfaces []string     `json:"interfaces,omitempty"`
	Visibility string       `json:"visibility"`
	Kind       string       `json:"kind"`
	Modifiers  []string     `json:"modifiers,omitempty"`
	Version    jsonVersion  `json:"version"`
	Fields     []jsonField  `json:"fields,omitempty"`
	Methods    []jsonMethod `json:"methods,omitempty"`
}

type jsonVersion struct {
	Major uint16 `json:"major"`
	Minor uint16 `json:"minor"`
}

type jsonField struct {
	Name       string   `json:"name"`
	Type       jsonType `json:"type"`
	Visibility string   `json:"visibility"`
	Modifiers  []string `json:"modifiers,omitempty"`
}

type jsonMethod struct {
	Name       string          `json:"name"`
	ReturnType jsonType        `json:"returnType"`
	Parameters []jsonParameter `json:"parameters,omitempty"`
	Visibility string          `json:"visibility"`
	Modifiers  []string        `json:"modifiers,omitempty"`
}

type jsonParameter struct {
	Name string   `json:"name,omitempty"`
	Type jsonType `json:"type"`
}

type jsonType struct {
	Name       string `json:"name"`
	ArrayDepth int    `json:"arrayDepth,omitempty"`
}

func (e *JSONEncoder) buildClassData() jsonClass {
	c := e.class
	data := jsonClass{
		Name:       c.Name(),
		SimpleName: c.SimpleName(),
		Package:    c.Package(),
		SuperClass: c.SuperClass(),
		Interfaces: c.Interfaces(),
		Visibility: c.Visibility(),
		Kind:       e.classKind(),
		Modifiers:  e.classModifiers(),
		Version: jsonVersion{
			Major: c.MajorVersion(),
			Minor: c.MinorVersion(),
		},
		Fields:  e.buildFields(),
		Methods: e.buildMethods(),
	}
	return data
}

func (e *JSONEncoder) classKind() string {
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
	default:
		return "class"
	}
}

func (e *JSONEncoder) classModifiers() []string {
	c := e.class
	var mods []string
	if c.IsFinal() {
		mods = append(mods, "final")
	}
	if c.IsAbstract() {
		mods = append(mods, "abstract")
	}
	if c.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	return mods
}

func (e *JSONEncoder) buildFields() []jsonField {
	fields := e.class.Fields()
	result := make([]jsonField, len(fields))
	for i, f := range fields {
		t := f.Type()
		result[i] = jsonField{
			Name: f.Name(),
			Type: jsonType{
				Name:       t.Name,
				ArrayDepth: t.ArrayDepth,
			},
			Visibility: f.Visibility(),
			Modifiers:  fieldModifiers(f),
		}
	}
	return result
}

func fieldModifiers(f java.Field) []string {
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
	return mods
}

func (e *JSONEncoder) buildMethods() []jsonMethod {
	methods := e.class.Methods()
	result := make([]jsonMethod, len(methods))
	for i, m := range methods {
		rt := m.ReturnType()
		result[i] = jsonMethod{
			Name: m.Name(),
			ReturnType: jsonType{
				Name:       rt.Name,
				ArrayDepth: rt.ArrayDepth,
			},
			Parameters: buildParameters(m.Parameters()),
			Visibility: m.Visibility(),
			Modifiers:  methodModifiers(m),
		}
	}
	return result
}

func buildParameters(params []java.Parameter) []jsonParameter {
	result := make([]jsonParameter, len(params))
	for i, p := range params {
		result[i] = jsonParameter{
			Name: p.Name,
			Type: jsonType{
				Name:       p.Type.Name,
				ArrayDepth: p.Type.ArrayDepth,
			},
		}
	}
	return result
}

func methodModifiers(m java.Method) []string {
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
	return mods
}
