package format

import (
	"io"
	"strings"

	"github.com/dhamidi/javalyzer/java"
)

type JavaEncoder struct {
	w     io.Writer
	class *java.Class
}

func NewJavaEncoder(w io.Writer) *JavaEncoder {
	return &JavaEncoder{w: w}
}

func (e *JavaEncoder) Encode(class *java.Class) error {
	e.class = class
	text, err := e.MarshalText()
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *JavaEncoder) MarshalText() ([]byte, error) {
	var sb strings.Builder
	c := e.class

	if pkg := c.Package(); pkg != "" {
		sb.WriteString("package ")
		sb.WriteString(pkg)
		sb.WriteString(";\n\n")
	}

	e.writeClassDeclaration(&sb)
	sb.WriteString(" {\n")

	e.writeFields(&sb)
	e.writeMethods(&sb)

	sb.WriteString("}\n")
	return []byte(sb.String()), nil
}

func (e *JavaEncoder) writeClassDeclaration(sb *strings.Builder) {
	c := e.class

	if c.IsPublic() {
		sb.WriteString("public ")
	}
	if c.IsAbstract() && !c.IsInterface() {
		sb.WriteString("abstract ")
	}
	if c.IsFinal() {
		sb.WriteString("final ")
	}

	switch {
	case c.IsAnnotation():
		sb.WriteString("@interface ")
	case c.IsEnum():
		sb.WriteString("enum ")
	case c.IsInterface():
		sb.WriteString("interface ")
	default:
		sb.WriteString("class ")
	}

	sb.WriteString(c.SimpleName())

	if super := c.SuperClass(); super != "" && super != "java.lang.Object" && !c.IsEnum() {
		sb.WriteString(" extends ")
		sb.WriteString(super)
	}

	if ifaces := c.Interfaces(); len(ifaces) > 0 {
		if c.IsInterface() {
			sb.WriteString(" extends ")
		} else {
			sb.WriteString(" implements ")
		}
		sb.WriteString(strings.Join(ifaces, ", "))
	}
}

func (e *JavaEncoder) writeFields(sb *strings.Builder) {
	fields := e.class.Fields()
	for _, f := range fields {
		if f.IsSynthetic() {
			continue
		}
		sb.WriteString("    ")
		e.writeFieldDeclaration(sb, f)
		sb.WriteString(";\n")
	}
	if len(fields) > 0 {
		sb.WriteString("\n")
	}
}

func (e *JavaEncoder) writeFieldDeclaration(sb *strings.Builder, f java.Field) {
	if f.IsPublic() {
		sb.WriteString("public ")
	} else if f.IsPrivate() {
		sb.WriteString("private ")
	} else if f.IsProtected() {
		sb.WriteString("protected ")
	}
	if f.IsStatic() {
		sb.WriteString("static ")
	}
	if f.IsFinal() {
		sb.WriteString("final ")
	}
	if f.IsVolatile() {
		sb.WriteString("volatile ")
	}
	if f.IsTransient() {
		sb.WriteString("transient ")
	}
	sb.WriteString(f.Type().String())
	sb.WriteString(" ")
	sb.WriteString(f.Name())
}

func (e *JavaEncoder) writeMethods(sb *strings.Builder) {
	methods := e.class.Methods()
	first := true
	for _, m := range methods {
		if m.IsSynthetic() || m.IsBridge() {
			continue
		}
		if m.IsStaticInitializer() {
			continue
		}
		if !first {
			sb.WriteString("\n")
		}
		first = false
		sb.WriteString("    ")
		e.writeMethodDeclaration(sb, m)
		if m.IsAbstract() || m.IsNative() || e.class.IsInterface() {
			sb.WriteString(";\n")
		} else {
			sb.WriteString(" { }\n")
		}
	}
}

func (e *JavaEncoder) writeMethodDeclaration(sb *strings.Builder, m java.Method) {
	if m.IsPublic() {
		sb.WriteString("public ")
	} else if m.IsPrivate() {
		sb.WriteString("private ")
	} else if m.IsProtected() {
		sb.WriteString("protected ")
	}
	if m.IsStatic() {
		sb.WriteString("static ")
	}
	if m.IsFinal() {
		sb.WriteString("final ")
	}
	if m.IsAbstract() && !e.class.IsInterface() {
		sb.WriteString("abstract ")
	}
	if m.IsSynchronized() {
		sb.WriteString("synchronized ")
	}
	if m.IsNative() {
		sb.WriteString("native ")
	}

	if m.IsConstructor() {
		sb.WriteString(e.class.SimpleName())
	} else {
		sb.WriteString(m.ReturnType().String())
		sb.WriteString(" ")
		sb.WriteString(m.Name())
	}

	sb.WriteString("(")
	params := m.Parameters()
	for i, p := range params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.Type.String())
		if p.Name != "" {
			sb.WriteString(" ")
			sb.WriteString(p.Name)
		}
	}
	sb.WriteString(")")
}
