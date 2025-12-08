package classfile

type MethodInfo struct {
	AccessFlags     AccessFlags
	NameIndex       uint16
	DescriptorIndex uint16
	Attributes      []AttributeInfo
}

func (m *MethodInfo) Name(cp ConstantPool) string {
	return cp.GetUtf8(m.NameIndex)
}

func (m *MethodInfo) Descriptor(cp ConstantPool) string {
	return cp.GetUtf8(m.DescriptorIndex)
}

func (m *MethodInfo) GetAttribute(cp ConstantPool, name string) *AttributeInfo {
	for i := range m.Attributes {
		if cp.GetUtf8(m.Attributes[i].NameIndex) == name {
			return &m.Attributes[i]
		}
	}
	return nil
}

func (m *MethodInfo) GetCodeAttribute(cp ConstantPool) *CodeAttribute {
	attr := m.GetAttribute(cp, "Code")
	if attr == nil {
		return nil
	}
	return attr.AsCode()
}

func (m *MethodInfo) IsPublic() bool       { return m.AccessFlags.IsPublic() }
func (m *MethodInfo) IsPrivate() bool      { return m.AccessFlags.IsPrivate() }
func (m *MethodInfo) IsProtected() bool    { return m.AccessFlags.IsProtected() }
func (m *MethodInfo) IsStatic() bool       { return m.AccessFlags.IsStatic() }
func (m *MethodInfo) IsFinal() bool        { return m.AccessFlags.IsFinal() }
func (m *MethodInfo) IsSynchronized() bool { return m.AccessFlags.IsSynchronized() }
func (m *MethodInfo) IsBridge() bool       { return m.AccessFlags.IsBridge() }
func (m *MethodInfo) IsVarargs() bool      { return m.AccessFlags.IsVarargs() }
func (m *MethodInfo) IsNative() bool       { return m.AccessFlags.IsNative() }
func (m *MethodInfo) IsAbstract() bool     { return m.AccessFlags.IsAbstract() }
func (m *MethodInfo) IsStrict() bool       { return m.AccessFlags.IsStrict() }
func (m *MethodInfo) IsSynthetic() bool    { return m.AccessFlags.IsSynthetic() }

func (m *MethodInfo) IsConstructor(cp ConstantPool) bool {
	return m.Name(cp) == "<init>"
}

func (m *MethodInfo) IsStaticInitializer(cp ConstantPool) bool {
	return m.Name(cp) == "<clinit>"
}

func (m *MethodInfo) ParsedDescriptor(cp ConstantPool) *MethodDescriptor {
	return ParseMethodDescriptor(m.Descriptor(cp))
}
