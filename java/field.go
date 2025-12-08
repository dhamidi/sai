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

func (f Field) Signature() string {
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
	return f.info.GetAttribute(f.cp, "Deprecated") != nil
}

func (f Field) Annotations() []Annotation {
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
