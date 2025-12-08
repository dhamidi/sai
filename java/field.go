package java

import "github.com/dhamidi/javalyzer/classfile"

type Field struct {
	info *classfile.FieldInfo
	cp   classfile.ConstantPool
}

func (f Field) Name() string {
	return f.info.Name(f.cp)
}

func (f Field) Descriptor() string {
	return f.info.Descriptor(f.cp)
}

func (f Field) Type() Type {
	ft := f.info.ParsedDescriptor(f.cp)
	return typeFromFieldType(ft)
}

func (f Field) IsPublic() bool    { return f.info.IsPublic() }
func (f Field) IsPrivate() bool   { return f.info.IsPrivate() }
func (f Field) IsProtected() bool { return f.info.IsProtected() }
func (f Field) IsStatic() bool    { return f.info.IsStatic() }
func (f Field) IsFinal() bool     { return f.info.IsFinal() }
func (f Field) IsVolatile() bool  { return f.info.IsVolatile() }
func (f Field) IsTransient() bool { return f.info.IsTransient() }
func (f Field) IsSynthetic() bool { return f.info.IsSynthetic() }
func (f Field) IsEnum() bool      { return f.info.IsEnum() }

func (f Field) Visibility() string {
	if f.IsPublic() {
		return "public"
	}
	if f.IsPrivate() {
		return "private"
	}
	if f.IsProtected() {
		return "protected"
	}
	return "package"
}

func (f Field) String() string {
	var result string
	if v := f.Visibility(); v != "package" {
		result = v + " "
	}
	if f.IsStatic() {
		result += "static "
	}
	if f.IsFinal() {
		result += "final "
	}
	if f.IsVolatile() {
		result += "volatile "
	}
	if f.IsTransient() {
		result += "transient "
	}
	result += f.Type().String() + " " + f.Name()
	return result
}
