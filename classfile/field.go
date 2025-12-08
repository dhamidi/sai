package classfile

type FieldInfo struct {
	AccessFlags     AccessFlags
	NameIndex       uint16
	DescriptorIndex uint16
	Attributes      []AttributeInfo
}

func (f *FieldInfo) Name(cp ConstantPool) string {
	return cp.GetUtf8(f.NameIndex)
}

func (f *FieldInfo) Descriptor(cp ConstantPool) string {
	return cp.GetUtf8(f.DescriptorIndex)
}

func (f *FieldInfo) GetAttribute(cp ConstantPool, name string) *AttributeInfo {
	for i := range f.Attributes {
		if cp.GetUtf8(f.Attributes[i].NameIndex) == name {
			return &f.Attributes[i]
		}
	}
	return nil
}

func (f *FieldInfo) IsPublic() bool    { return f.AccessFlags.IsPublic() }
func (f *FieldInfo) IsPrivate() bool   { return f.AccessFlags.IsPrivate() }
func (f *FieldInfo) IsProtected() bool { return f.AccessFlags.IsProtected() }
func (f *FieldInfo) IsStatic() bool    { return f.AccessFlags.IsStatic() }
func (f *FieldInfo) IsFinal() bool     { return f.AccessFlags.IsFinal() }
func (f *FieldInfo) IsVolatile() bool  { return f.AccessFlags.IsVolatile() }
func (f *FieldInfo) IsTransient() bool { return f.AccessFlags.IsTransient() }
func (f *FieldInfo) IsSynthetic() bool { return f.AccessFlags.IsSynthetic() }
func (f *FieldInfo) IsEnum() bool      { return f.AccessFlags.IsEnum() }

func (f *FieldInfo) ParsedDescriptor(cp ConstantPool) *FieldType {
	return ParseFieldDescriptor(f.Descriptor(cp))
}
