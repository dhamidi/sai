package classfile

type ClassFile struct {
	MinorVersion uint16
	MajorVersion uint16
	ConstantPool ConstantPool
	AccessFlags  AccessFlags
	ThisClass    uint16
	SuperClass   uint16
	Interfaces   []uint16
	Fields       []FieldInfo
	Methods      []MethodInfo
	Attributes   []AttributeInfo
}

func (cf *ClassFile) ClassName() string {
	return cf.ConstantPool.GetClassName(cf.ThisClass)
}

func (cf *ClassFile) SuperClassName() string {
	if cf.SuperClass == 0 {
		return ""
	}
	return cf.ConstantPool.GetClassName(cf.SuperClass)
}

func (cf *ClassFile) InterfaceNames() []string {
	names := make([]string, len(cf.Interfaces))
	for i, idx := range cf.Interfaces {
		names[i] = cf.ConstantPool.GetClassName(idx)
	}
	return names
}

func (cf *ClassFile) IsClass() bool {
	return !cf.AccessFlags.IsInterface() && !cf.AccessFlags.IsModule()
}

func (cf *ClassFile) IsInterface() bool {
	return cf.AccessFlags.IsInterface() && !cf.AccessFlags.IsAnnotation()
}

func (cf *ClassFile) IsAnnotation() bool {
	return cf.AccessFlags.IsAnnotation()
}

func (cf *ClassFile) IsEnum() bool {
	return cf.AccessFlags.IsEnum()
}

func (cf *ClassFile) IsModule() bool {
	return cf.AccessFlags.IsModule()
}

func (cf *ClassFile) GetField(name string) *FieldInfo {
	for i := range cf.Fields {
		if cf.Fields[i].Name(cf.ConstantPool) == name {
			return &cf.Fields[i]
		}
	}
	return nil
}

func (cf *ClassFile) GetMethod(name, descriptor string) *MethodInfo {
	for i := range cf.Methods {
		if cf.Methods[i].Name(cf.ConstantPool) == name {
			if descriptor == "" || cf.Methods[i].Descriptor(cf.ConstantPool) == descriptor {
				return &cf.Methods[i]
			}
		}
	}
	return nil
}

func (cf *ClassFile) GetMethods(name string) []*MethodInfo {
	var methods []*MethodInfo
	for i := range cf.Methods {
		if cf.Methods[i].Name(cf.ConstantPool) == name {
			methods = append(methods, &cf.Methods[i])
		}
	}
	return methods
}

func (cf *ClassFile) GetAttribute(name string) *AttributeInfo {
	for i := range cf.Attributes {
		if cf.ConstantPool.GetUtf8(cf.Attributes[i].NameIndex) == name {
			return &cf.Attributes[i]
		}
	}
	return nil
}
