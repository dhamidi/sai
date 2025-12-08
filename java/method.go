package java

import "github.com/dhamidi/javalyzer/classfile"

type Method struct {
	info *classfile.MethodInfo
	cp   classfile.ConstantPool
}

func (m Method) Name() string {
	return m.info.Name(m.cp)
}

func (m Method) Descriptor() string {
	return m.info.Descriptor(m.cp)
}

func (m Method) ReturnType() Type {
	desc := m.info.ParsedDescriptor(m.cp)
	if desc == nil {
		return Type{Name: "void"}
	}
	return typeFromFieldType(desc.ReturnType)
}

func (m Method) Parameters() []Parameter {
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

func (m Method) IsPublic() bool       { return m.info.IsPublic() }
func (m Method) IsPrivate() bool      { return m.info.IsPrivate() }
func (m Method) IsProtected() bool    { return m.info.IsProtected() }
func (m Method) IsStatic() bool       { return m.info.IsStatic() }
func (m Method) IsFinal() bool        { return m.info.IsFinal() }
func (m Method) IsSynchronized() bool { return m.info.IsSynchronized() }
func (m Method) IsBridge() bool       { return m.info.IsBridge() }
func (m Method) IsVarargs() bool      { return m.info.IsVarargs() }
func (m Method) IsNative() bool       { return m.info.IsNative() }
func (m Method) IsAbstract() bool     { return m.info.IsAbstract() }
func (m Method) IsSynthetic() bool    { return m.info.IsSynthetic() }
func (m Method) IsConstructor() bool  { return m.info.IsConstructor(m.cp) }
func (m Method) IsStaticInitializer() bool { return m.info.IsStaticInitializer(m.cp) }

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
