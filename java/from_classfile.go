package java

import (
	"io"
	"os"
	"strings"

	"github.com/dhamidi/sai/classfile"
)

func ClassModelFromFile(path string) (*ClassModel, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ClassModelFromReader(f)
}

func ClassModelFromReader(r io.Reader) (*ClassModel, error) {
	cf, err := classfile.Parse(r)
	if err != nil {
		return nil, err
	}
	return ClassModelFromClassFile(cf), nil
}

func ClassModelFromClassFile(cf *classfile.ClassFile) *ClassModel {
	className := classfile.InternalToSourceName(cf.ClassName())
	pkg, simpleName := splitClassName(className)

	model := &ClassModel{
		Name:         className,
		SimpleName:   simpleName,
		Package:      pkg,
		MajorVersion: cf.MajorVersion,
		MinorVersion: cf.MinorVersion,
		Visibility:   visibilityFromAccessFlags(cf.AccessFlags),
		Kind:         classKindFromClassFile(cf),
		IsFinal:      cf.AccessFlags.IsFinal(),
		IsAbstract:   cf.AccessFlags.IsAbstract(),
		IsSynthetic:  cf.AccessFlags.IsSynthetic(),
	}

	if cf.SuperClass != 0 {
		model.SuperClass = classfile.InternalToSourceName(cf.SuperClassName())
	}

	for _, iface := range cf.InterfaceNames() {
		model.Interfaces = append(model.Interfaces, classfile.InternalToSourceName(iface))
	}

	for i := range cf.Fields {
		field := &cf.Fields[i]
		if field.IsSynthetic() {
			continue
		}
		model.Fields = append(model.Fields, fieldModelFromFieldInfo(field, cf.ConstantPool))
	}

	for i := range cf.Methods {
		method := &cf.Methods[i]
		if method.IsSynthetic() || method.IsBridge() {
			continue
		}
		if method.IsStaticInitializer(cf.ConstantPool) {
			continue
		}
		model.Methods = append(model.Methods, methodModelFromMethodInfo(method, cf.ConstantPool))
	}

	return model
}

func splitClassName(fullName string) (pkg, simpleName string) {
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 {
		return "", fullName
	}
	return fullName[:lastDot], fullName[lastDot+1:]
}

func extractSimpleName(fullName string) string {
	_, simpleName := splitClassName(fullName)
	return simpleName
}

func visibilityFromAccessFlags(flags classfile.AccessFlags) Visibility {
	if flags.IsPublic() {
		return VisibilityPublic
	}
	if flags.IsProtected() {
		return VisibilityProtected
	}
	if flags.IsPrivate() {
		return VisibilityPrivate
	}
	return VisibilityPackage
}

func classKindFromClassFile(cf *classfile.ClassFile) ClassKind {
	if cf.IsAnnotation() {
		return ClassKindAnnotation
	}
	if cf.IsEnum() {
		return ClassKindEnum
	}
	if cf.IsInterface() {
		return ClassKindInterface
	}
	// Note: records are compiled as regular classes with a Record attribute
	// We could check for Record attribute here if needed
	return ClassKindClass
}

func fieldModelFromFieldInfo(f *classfile.FieldInfo, cp classfile.ConstantPool) FieldModel {
	desc := f.ParsedDescriptor(cp)
	return FieldModel{
		Name:        f.Name(cp),
		Type:        typeModelFromFieldType(desc),
		Visibility:  visibilityFromAccessFlags(f.AccessFlags),
		IsStatic:    f.IsStatic(),
		IsFinal:     f.IsFinal(),
		IsVolatile:  f.IsVolatile(),
		IsTransient: f.IsTransient(),
		IsSynthetic: f.IsSynthetic(),
		IsEnum:      f.IsEnum(),
	}
}

func methodModelFromMethodInfo(m *classfile.MethodInfo, cp classfile.ConstantPool) MethodModel {
	desc := m.ParsedDescriptor(cp)
	model := MethodModel{
		Name:           m.Name(cp),
		Visibility:     visibilityFromAccessFlags(m.AccessFlags),
		IsStatic:       m.IsStatic(),
		IsFinal:        m.IsFinal(),
		IsAbstract:     m.IsAbstract(),
		IsSynchronized: m.IsSynchronized(),
		IsNative:       m.IsNative(),
		IsBridge:       m.IsBridge(),
		IsVarargs:      m.IsVarargs(),
		IsSynthetic:    m.IsSynthetic(),
	}

	if desc != nil {
		if desc.ReturnType != nil {
			model.ReturnType = typeModelFromFieldType(desc.ReturnType)
		} else {
			model.ReturnType = TypeModel{Name: "void"}
		}

		for i, param := range desc.Parameters {
			model.Parameters = append(model.Parameters, ParameterModel{
				Name: parameterName(i),
				Type: typeModelFromFieldType(&param),
			})
		}
	}

	return model
}

func typeModelFromFieldType(ft *classfile.FieldType) TypeModel {
	if ft == nil {
		return TypeModel{Name: "void"}
	}

	model := TypeModel{
		ArrayDepth: ft.ArrayDepth,
	}

	if ft.BaseType != "" {
		model.Name = ft.BaseType
	} else if ft.ClassName != "" {
		model.Name = classfile.InternalToSourceName(ft.ClassName)
	}

	return model
}

func parameterName(index int) string {
	return "" // Class files don't always have parameter names
}
