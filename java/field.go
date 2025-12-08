package java

import "github.com/dhamidi/javalyzer/classfile"

type Field struct {
	info   *classfile.FieldInfo
	cp     classfile.ConstantPool
	source *sourceField
}

func (f Field) Name() string {
	if f.source != nil {
		return f.source.name
	}
	return f.info.Name(f.cp)
}

func (f Field) Descriptor() string {
	if f.source != nil {
		return ""
	}
	return f.info.Descriptor(f.cp)
}

func (f Field) Type() Type {
	if f.source != nil {
		return Type{Name: f.source.typeName, ArrayDepth: f.source.arrayDepth}
	}
	ft := f.info.ParsedDescriptor(f.cp)
	return typeFromFieldType(ft)
}

func (f Field) IsPublic() bool {
	if f.source != nil {
		return f.source.visibility == "public"
	}
	return f.info.IsPublic()
}
func (f Field) IsPrivate() bool {
	if f.source != nil {
		return f.source.visibility == "private"
	}
	return f.info.IsPrivate()
}
func (f Field) IsProtected() bool {
	if f.source != nil {
		return f.source.visibility == "protected"
	}
	return f.info.IsProtected()
}
func (f Field) IsStatic() bool {
	if f.source != nil {
		return f.source.isStatic
	}
	return f.info.IsStatic()
}
func (f Field) IsFinal() bool {
	if f.source != nil {
		return f.source.isFinal
	}
	return f.info.IsFinal()
}
func (f Field) IsVolatile() bool {
	if f.source != nil {
		return false
	}
	return f.info.IsVolatile()
}
func (f Field) IsTransient() bool {
	if f.source != nil {
		return false
	}
	return f.info.IsTransient()
}
func (f Field) IsSynthetic() bool {
	if f.source != nil {
		return false
	}
	return f.info.IsSynthetic()
}
func (f Field) IsEnum() bool {
	if f.source != nil {
		return false
	}
	return f.info.IsEnum()
}

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

func (f Field) Signature() string {
	if f.source != nil {
		return ""
	}
	attr := f.info.GetAttribute(f.cp, "Signature")
	if attr == nil {
		return ""
	}
	if sig := attr.AsSignature(); sig != nil {
		return f.cp.GetUtf8(sig.SignatureIndex)
	}
	return ""
}

func (f Field) IsDeprecated() bool {
	if f.source != nil {
		return false
	}
	return f.info.GetAttribute(f.cp, "Deprecated") != nil
}

func (f Field) Annotations() []Annotation {
	if f.source != nil {
		return nil
	}
	attr := f.info.GetAttribute(f.cp, "RuntimeVisibleAnnotations")
	if attr == nil {
		return nil
	}
	if rva := attr.AsRuntimeVisibleAnnotations(); rva != nil {
		return annotationsFromClassfile(rva.Annotations, f.cp)
	}
	return nil
}

func (f Field) InvisibleAnnotations() []Annotation {
	if f.source != nil {
		return nil
	}
	attr := f.info.GetAttribute(f.cp, "RuntimeInvisibleAnnotations")
	if attr == nil {
		return nil
	}
	if ria := attr.AsRuntimeInvisibleAnnotations(); ria != nil {
		return annotationsFromClassfile(ria.Annotations, f.cp)
	}
	return nil
}

func (f Field) ConstantValue() interface{} {
	if f.source != nil {
		return nil
	}
	attr := f.info.GetAttribute(f.cp, "ConstantValue")
	if attr == nil {
		return nil
	}
	if cv := attr.AsConstantValue(); cv != nil {
		if val, ok := f.cp.GetInteger(cv.ConstantValueIndex); ok {
			return val
		}
		if val, ok := f.cp.GetLong(cv.ConstantValueIndex); ok {
			return val
		}
		if val, ok := f.cp.GetFloat(cv.ConstantValueIndex); ok {
			return val
		}
		if val, ok := f.cp.GetDouble(cv.ConstantValueIndex); ok {
			return val
		}
		if val := f.cp.GetString(cv.ConstantValueIndex); val != "" {
			return val
		}
	}
	return nil
}
