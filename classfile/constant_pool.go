package classfile

type ConstantPoolEntry interface {
	Tag() ConstantTag
}

type ConstantUtf8Info struct {
	Value string
}

func (c *ConstantUtf8Info) Tag() ConstantTag { return ConstantUtf8 }

type ConstantIntegerInfo struct {
	Value int32
}

func (c *ConstantIntegerInfo) Tag() ConstantTag { return ConstantInteger }

type ConstantFloatInfo struct {
	Value float32
}

func (c *ConstantFloatInfo) Tag() ConstantTag { return ConstantFloat }

type ConstantLongInfo struct {
	Value int64
}

func (c *ConstantLongInfo) Tag() ConstantTag { return ConstantLong }

type ConstantDoubleInfo struct {
	Value float64
}

func (c *ConstantDoubleInfo) Tag() ConstantTag { return ConstantDouble }

type ConstantClassInfo struct {
	NameIndex uint16
}

func (c *ConstantClassInfo) Tag() ConstantTag { return ConstantClass }

type ConstantStringInfo struct {
	StringIndex uint16
}

func (c *ConstantStringInfo) Tag() ConstantTag { return ConstantString }

type ConstantFieldrefInfo struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantFieldrefInfo) Tag() ConstantTag { return ConstantFieldref }

type ConstantMethodrefInfo struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantMethodrefInfo) Tag() ConstantTag { return ConstantMethodref }

type ConstantInterfaceMethodrefInfo struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantInterfaceMethodrefInfo) Tag() ConstantTag { return ConstantInterfaceMethodref }

type ConstantNameAndTypeInfo struct {
	NameIndex       uint16
	DescriptorIndex uint16
}

func (c *ConstantNameAndTypeInfo) Tag() ConstantTag { return ConstantNameAndType }

type ConstantMethodHandleInfo struct {
	ReferenceKind  MethodHandleKind
	ReferenceIndex uint16
}

func (c *ConstantMethodHandleInfo) Tag() ConstantTag { return ConstantMethodHandle }

type ConstantMethodTypeInfo struct {
	DescriptorIndex uint16
}

func (c *ConstantMethodTypeInfo) Tag() ConstantTag { return ConstantMethodType }

type ConstantDynamicInfo struct {
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (c *ConstantDynamicInfo) Tag() ConstantTag { return ConstantDynamic }

type ConstantInvokeDynamicInfo struct {
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (c *ConstantInvokeDynamicInfo) Tag() ConstantTag { return ConstantInvokeDynamic }

type ConstantModuleInfo struct {
	NameIndex uint16
}

func (c *ConstantModuleInfo) Tag() ConstantTag { return ConstantModule }

type ConstantPackageInfo struct {
	NameIndex uint16
}

func (c *ConstantPackageInfo) Tag() ConstantTag { return ConstantPackage }

type ConstantPool []ConstantPoolEntry

func (cp ConstantPool) GetUtf8(index uint16) string {
	if index == 0 || int(index) > len(cp) {
		return ""
	}
	if entry, ok := cp[index-1].(*ConstantUtf8Info); ok {
		return entry.Value
	}
	return ""
}

func (cp ConstantPool) GetClassName(index uint16) string {
	if index == 0 || int(index) > len(cp) {
		return ""
	}
	if entry, ok := cp[index-1].(*ConstantClassInfo); ok {
		return cp.GetUtf8(entry.NameIndex)
	}
	return ""
}

func (cp ConstantPool) GetNameAndType(index uint16) (name, descriptor string) {
	if index == 0 || int(index) > len(cp) {
		return "", ""
	}
	if entry, ok := cp[index-1].(*ConstantNameAndTypeInfo); ok {
		return cp.GetUtf8(entry.NameIndex), cp.GetUtf8(entry.DescriptorIndex)
	}
	return "", ""
}

func (cp ConstantPool) GetString(index uint16) string {
	if index == 0 || int(index) > len(cp) {
		return ""
	}
	if entry, ok := cp[index-1].(*ConstantStringInfo); ok {
		return cp.GetUtf8(entry.StringIndex)
	}
	return ""
}

func (cp ConstantPool) GetModuleName(index uint16) string {
	if index == 0 || int(index) > len(cp) {
		return ""
	}
	if entry, ok := cp[index-1].(*ConstantModuleInfo); ok {
		return cp.GetUtf8(entry.NameIndex)
	}
	return ""
}

func (cp ConstantPool) GetPackageName(index uint16) string {
	if index == 0 || int(index) > len(cp) {
		return ""
	}
	if entry, ok := cp[index-1].(*ConstantPackageInfo); ok {
		return cp.GetUtf8(entry.NameIndex)
	}
	return ""
}

func (cp ConstantPool) GetInteger(index uint16) (int32, bool) {
	if index == 0 || int(index) > len(cp) {
		return 0, false
	}
	if entry, ok := cp[index-1].(*ConstantIntegerInfo); ok {
		return entry.Value, true
	}
	return 0, false
}

func (cp ConstantPool) GetLong(index uint16) (int64, bool) {
	if index == 0 || int(index) > len(cp) {
		return 0, false
	}
	if entry, ok := cp[index-1].(*ConstantLongInfo); ok {
		return entry.Value, true
	}
	return 0, false
}

func (cp ConstantPool) GetFloat(index uint16) (float32, bool) {
	if index == 0 || int(index) > len(cp) {
		return 0, false
	}
	if entry, ok := cp[index-1].(*ConstantFloatInfo); ok {
		return entry.Value, true
	}
	return 0, false
}

func (cp ConstantPool) GetDouble(index uint16) (float64, bool) {
	if index == 0 || int(index) > len(cp) {
		return 0, false
	}
	if entry, ok := cp[index-1].(*ConstantDoubleInfo); ok {
		return entry.Value, true
	}
	return 0, false
}

func (cp ConstantPool) GetFieldref(index uint16) (className, name, descriptor string) {
	if index == 0 || int(index) > len(cp) {
		return "", "", ""
	}
	if entry, ok := cp[index-1].(*ConstantFieldrefInfo); ok {
		className = cp.GetClassName(entry.ClassIndex)
		name, descriptor = cp.GetNameAndType(entry.NameAndTypeIndex)
		return
	}
	return "", "", ""
}

func (cp ConstantPool) GetMethodref(index uint16) (className, name, descriptor string) {
	if index == 0 || int(index) > len(cp) {
		return "", "", ""
	}
	if entry, ok := cp[index-1].(*ConstantMethodrefInfo); ok {
		className = cp.GetClassName(entry.ClassIndex)
		name, descriptor = cp.GetNameAndType(entry.NameAndTypeIndex)
		return
	}
	return "", "", ""
}

func (cp ConstantPool) GetInterfaceMethodref(index uint16) (className, name, descriptor string) {
	if index == 0 || int(index) > len(cp) {
		return "", "", ""
	}
	if entry, ok := cp[index-1].(*ConstantInterfaceMethodrefInfo); ok {
		className = cp.GetClassName(entry.ClassIndex)
		name, descriptor = cp.GetNameAndType(entry.NameAndTypeIndex)
		return
	}
	return "", "", ""
}

func (cp ConstantPool) GetMethodHandle(index uint16) *ConstantMethodHandleInfo {
	if index == 0 || int(index) > len(cp) {
		return nil
	}
	if entry, ok := cp[index-1].(*ConstantMethodHandleInfo); ok {
		return entry
	}
	return nil
}

func (cp ConstantPool) GetMethodType(index uint16) string {
	if index == 0 || int(index) > len(cp) {
		return ""
	}
	if entry, ok := cp[index-1].(*ConstantMethodTypeInfo); ok {
		return cp.GetUtf8(entry.DescriptorIndex)
	}
	return ""
}

func (cp ConstantPool) GetDynamic(index uint16) *ConstantDynamicInfo {
	if index == 0 || int(index) > len(cp) {
		return nil
	}
	if entry, ok := cp[index-1].(*ConstantDynamicInfo); ok {
		return entry
	}
	return nil
}

func (cp ConstantPool) GetInvokeDynamic(index uint16) *ConstantInvokeDynamicInfo {
	if index == 0 || int(index) > len(cp) {
		return nil
	}
	if entry, ok := cp[index-1].(*ConstantInvokeDynamicInfo); ok {
		return entry
	}
	return nil
}
