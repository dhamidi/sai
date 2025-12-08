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

	e.writeAnnotations(sb, c.Annotations(), "")

	if c.IsPublic() {
		sb.WriteString("public ")
	}
	if c.IsAbstract() && !c.IsInterface() && !c.IsSealed() {
		sb.WriteString("abstract ")
	}
	if c.IsSealed() {
		sb.WriteString("sealed ")
	}
	if c.IsFinal() && !c.IsRecord() {
		sb.WriteString("final ")
	}

	switch {
	case c.IsAnnotation():
		sb.WriteString("@interface ")
	case c.IsEnum():
		sb.WriteString("enum ")
	case c.IsRecord():
		sb.WriteString("record ")
	case c.IsInterface():
		sb.WriteString("interface ")
	default:
		sb.WriteString("class ")
	}

	sb.WriteString(c.SimpleName())

	if c.IsRecord() {
		e.writeRecordComponents(sb)
	}

	if super := c.SuperClass(); super != "" && super != "java.lang.Object" && super != "java.lang.Record" && !c.IsEnum() {
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

	if permitted := c.PermittedSubclasses(); len(permitted) > 0 {
		sb.WriteString(" permits ")
		sb.WriteString(strings.Join(permitted, ", "))
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
	e.writeAnnotations(sb, f.Annotations(), "    ")
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
	e.writeAnnotations(sb, m.Annotations(), "    ")
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

	if exceptions := m.Exceptions(); len(exceptions) > 0 {
		sb.WriteString(" throws ")
		sb.WriteString(strings.Join(exceptions, ", "))
	}
}

func (e *JavaEncoder) writeAnnotations(sb *strings.Builder, anns []java.Annotation, indent string) {
	for _, a := range anns {
		sb.WriteString("@")
		sb.WriteString(a.Type)
		if len(a.ElementValuePairs) > 0 {
			sb.WriteString("(")
			for i, p := range a.ElementValuePairs {
				if i > 0 {
					sb.WriteString(", ")
				}
				if len(a.ElementValuePairs) == 1 && p.Name == "value" {
					e.writeAnnotationValue(sb, p.Value)
				} else {
					sb.WriteString(p.Name)
					sb.WriteString(" = ")
					e.writeAnnotationValue(sb, p.Value)
				}
			}
			sb.WriteString(")")
		}
		sb.WriteString("\n")
		sb.WriteString(indent)
	}
}

func (e *JavaEncoder) writeAnnotationValue(sb *strings.Builder, v interface{}) {
	switch val := v.(type) {
	case string:
		sb.WriteString("\"")
		sb.WriteString(val)
		sb.WriteString("\"")
	case int32:
		sb.WriteString(strings.TrimSpace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						strings.Replace(
							"                                ",
							" ", "", -1),
						"", "", -1),
					"", "", -1),
				"", "", -1),
			"", "", -1),
		"", "", -1), "", "", -1), "", "", -1)))
		sb.WriteString(itoa(int(val)))
	case int64:
		sb.WriteString(itoa64(val))
		sb.WriteString("L")
	case float32:
		sb.WriteString(ftoa32(val))
		sb.WriteString("f")
	case float64:
		sb.WriteString(ftoa64(val))
	case map[string]string:
		if typ, ok := val["type"]; ok {
			sb.WriteString(typ)
			sb.WriteString(".")
			sb.WriteString(val["value"])
		}
	case []interface{}:
		sb.WriteString("{")
		for i, elem := range val {
			if i > 0 {
				sb.WriteString(", ")
			}
			e.writeAnnotationValue(sb, elem)
		}
		sb.WriteString("}")
	case java.Annotation:
		sb.WriteString("@")
		sb.WriteString(val.Type)
		if len(val.ElementValuePairs) > 0 {
			sb.WriteString("(")
			for i, p := range val.ElementValuePairs {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(p.Name)
				sb.WriteString(" = ")
				e.writeAnnotationValue(sb, p.Value)
			}
			sb.WriteString(")")
		}
	default:
		sb.WriteString("?")
	}
}

func (e *JavaEncoder) writeRecordComponents(sb *strings.Builder) {
	comps := e.class.RecordComponents()
	sb.WriteString("(")
	for i, c := range comps {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(c.Type().String())
		sb.WriteString(" ")
		sb.WriteString(c.Name)
	}
	sb.WriteString(")")
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

func itoa64(i int64) string {
	if i == 0 {
		return "0"
	}
	s := ""
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

func ftoa32(f float32) string {
	return strings.TrimRight(strings.TrimRight(formatFloat(float64(f), 6), "0"), ".")
}

func ftoa64(f float64) string {
	return strings.TrimRight(strings.TrimRight(formatFloat(f, 15), "0"), ".")
}

func formatFloat(f float64, prec int) string {
	if f == 0 {
		return "0"
	}
	neg := f < 0
	if neg {
		f = -f
	}
	intPart := int64(f)
	fracPart := f - float64(intPart)
	s := itoa64(intPart)
	if fracPart > 0 {
		s += "."
		for i := 0; i < prec && fracPart > 0; i++ {
			fracPart *= 10
			digit := int(fracPart)
			s += string(rune('0' + digit))
			fracPart -= float64(digit)
		}
	}
	if neg {
		s = "-" + s
	}
	return s
}
