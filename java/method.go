package java

import "github.com/dhamidi/javalyzer/classfile"

type Method struct {
	info   *classfile.MethodInfo
	cp     classfile.ConstantPool
	source *sourceMethod
}

func (m Method) Name() string {
	if m.source != nil {
		return m.source.name
	}
	return m.info.Name(m.cp)
}

func (m Method) Descriptor() string {
	if m.source != nil {
		return ""
	}
	return m.info.Descriptor(m.cp)
}

func (m Method) ReturnType() Type {
	if m.source != nil {
		return m.source.returnType
	}
	desc := m.info.ParsedDescriptor(m.cp)
	if desc == nil {
		return Type{Name: "void"}
	}
	return typeFromFieldType(desc.ReturnType)
}

func (m Method) Parameters() []Parameter {
	if m.source != nil {
		return m.source.parameters
	}
	desc := m.info.ParsedDescriptor(m.cp)
	if desc == nil {
		return nil
	}

	params := make([]Parameter, len(desc.Parameters))
	for i, p := range desc.Parameters {
		params[i] = Parameter{
			Type:  typeFromFieldType(&p),
			Index: i,
		}
	}

	m.populateParameterNames(params)
	return params
}

func (m Method) populateParameterNames(params []Parameter) {
	code := m.info.GetCodeAttribute(m.cp)
	if code == nil {
		return
	}

	var lvt *classfile.LocalVariableTableAttribute
	for i := range code.Attributes {
		if attr := code.Attributes[i].AsLocalVariableTable(); attr != nil {
			lvt = attr
			break
		}
	}
	if lvt == nil {
		return
	}

	startSlot := 0
	if !m.IsStatic() {
		startSlot = 1
	}

	for i := range params {
		slot := startSlot + i
		for _, lv := range lvt.LocalVariableTable {
			if int(lv.Index) == slot && lv.StartPC == 0 {
				params[i].Name = m.cp.GetUtf8(lv.NameIndex)
				break
			}
		}
	}
}

func (m Method) ParameterCount() int {
	desc := m.info.ParsedDescriptor(m.cp)
	if desc == nil {
		return 0
	}
	return len(desc.Parameters)
}

func (m Method) IsPublic() bool {
	if m.source != nil {
		return m.source.visibility == "public"
	}
	return m.info.IsPublic()
}
func (m Method) IsPrivate() bool {
	if m.source != nil {
		return m.source.visibility == "private"
	}
	return m.info.IsPrivate()
}
func (m Method) IsProtected() bool {
	if m.source != nil {
		return m.source.visibility == "protected"
	}
	return m.info.IsProtected()
}
func (m Method) IsStatic() bool {
	if m.source != nil {
		return m.source.isStatic
	}
	return m.info.IsStatic()
}
func (m Method) IsFinal() bool {
	if m.source != nil {
		return m.source.isFinal
	}
	return m.info.IsFinal()
}
func (m Method) IsSynchronized() bool {
	if m.source != nil {
		return false
	}
	return m.info.IsSynchronized()
}
func (m Method) IsBridge() bool {
	if m.source != nil {
		return false
	}
	return m.info.IsBridge()
}
func (m Method) IsVarargs() bool {
	if m.source != nil {
		return false
	}
	return m.info.IsVarargs()
}
func (m Method) IsNative() bool {
	if m.source != nil {
		return false
	}
	return m.info.IsNative()
}
func (m Method) IsAbstract() bool {
	if m.source != nil {
		return m.source.isAbstract
	}
	return m.info.IsAbstract()
}
func (m Method) IsSynthetic() bool {
	if m.source != nil {
		return false
	}
	return m.info.IsSynthetic()
}
func (m Method) IsConstructor() bool {
	if m.source != nil {
		return m.source.name == "<init>"
	}
	return m.info.IsConstructor(m.cp)
}
func (m Method) IsStaticInitializer() bool {
	if m.source != nil {
		return m.source.name == "<clinit>"
	}
	return m.info.IsStaticInitializer(m.cp)
}

func (m Method) Visibility() string {
	if m.IsPublic() {
		return "public"
	}
	if m.IsPrivate() {
		return "private"
	}
	if m.IsProtected() {
		return "protected"
	}
	return "package"
}

func (m Method) String() string {
	var result string
	if v := m.Visibility(); v != "package" {
		result = v + " "
	}
	if m.IsStatic() {
		result += "static "
	}
	if m.IsFinal() {
		result += "final "
	}
	if m.IsAbstract() {
		result += "abstract "
	}
	if m.IsSynchronized() {
		result += "synchronized "
	}
	if m.IsNative() {
		result += "native "
	}

	result += m.ReturnType().String() + " " + m.Name() + "("
	params := m.Parameters()
	for i, p := range params {
		if i > 0 {
			result += ", "
		}
		result += p.String()
	}
	result += ")"
	return result
}

func (m Method) Signature() string {
	if m.source != nil {
		return ""
	}
	attr := m.info.GetAttribute(m.cp, "Signature")
	if attr == nil {
		return ""
	}
	if sig := attr.AsSignature(); sig != nil {
		return m.cp.GetUtf8(sig.SignatureIndex)
	}
	return ""
}

func (m Method) IsDeprecated() bool {
	if m.source != nil {
		return false
	}
	return m.info.GetAttribute(m.cp, "Deprecated") != nil
}

func (m Method) Annotations() []Annotation {
	if m.source != nil {
		return nil
	}
	attr := m.info.GetAttribute(m.cp, "RuntimeVisibleAnnotations")
	if attr == nil {
		return nil
	}
	if rva := attr.AsRuntimeVisibleAnnotations(); rva != nil {
		return annotationsFromClassfile(rva.Annotations, m.cp)
	}
	return nil
}

func (m Method) InvisibleAnnotations() []Annotation {
	if m.source != nil {
		return nil
	}
	attr := m.info.GetAttribute(m.cp, "RuntimeInvisibleAnnotations")
	if attr == nil {
		return nil
	}
	if ria := attr.AsRuntimeInvisibleAnnotations(); ria != nil {
		return annotationsFromClassfile(ria.Annotations, m.cp)
	}
	return nil
}

func (m Method) ParameterAnnotations() [][]Annotation {
	if m.source != nil {
		return nil
	}
	attr := m.info.GetAttribute(m.cp, "RuntimeVisibleParameterAnnotations")
	if attr == nil {
		return nil
	}
	if rvpa := attr.AsRuntimeVisibleParameterAnnotations(); rvpa != nil {
		result := make([][]Annotation, len(rvpa.ParameterAnnotations))
		for i, anns := range rvpa.ParameterAnnotations {
			result[i] = annotationsFromClassfile(anns, m.cp)
		}
		return result
	}
	return nil
}

func (m Method) Exceptions() []string {
	if m.source != nil {
		return nil
	}
	attr := m.info.GetAttribute(m.cp, "Exceptions")
	if attr == nil {
		return nil
	}
	if ex := attr.AsExceptions(); ex != nil {
		result := make([]string, len(ex.ExceptionIndexTable))
		for i, idx := range ex.ExceptionIndexTable {
			result[i] = classfile.InternalToSourceName(m.cp.GetClassName(idx))
		}
		return result
	}
	return nil
}
