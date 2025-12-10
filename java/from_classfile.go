package java

import (
	"io"

	"github.com/dhamidi/sai/classfile"
)

func ClassModelFromReader(r io.Reader) (*ClassModel, error) {
	cf, err := classfile.Parse(r)
	if err != nil {
		return nil, err
	}
	return ClassModelFromClassFile(cf), nil
}

func ClassModelFromFile(path string) (*ClassModel, error) {
	cf, err := classfile.ParseFile(path)
	if err != nil {
		return nil, err
	}
	model := ClassModelFromClassFile(cf)
	model.SourceURL = FileURL(path)
	return model, nil
}

func ClassModelFromClassFile(cf *classfile.ClassFile) *ClassModel {
	name := classfile.InternalToSourceName(cf.ClassName())
	simpleName := extractSimpleName(name)
	pkg := extractPackage(name)

	model := &ClassModel{
		Name:                name,
		SimpleName:          simpleName,
		Package:             pkg,
		SuperClass:          superClassFromClassFile(cf),
		Interfaces:          interfacesFromClassFile(cf),
		Visibility:          visibilityFromClassAccessFlags(cf.AccessFlags),
		Kind:                kindFromClassFile(cf),
		IsFinal:             cf.AccessFlags.IsFinal(),
		IsAbstract:          cf.AccessFlags.IsAbstract(),
		IsStatic:            cf.AccessFlags.IsStatic(),
		IsSynthetic:         cf.AccessFlags.IsSynthetic(),
		MajorVersion:        cf.MajorVersion,
		MinorVersion:        cf.MinorVersion,
		Signature:           signatureFromClassFile(cf),
		SourceFile:          sourceFileFromClassFile(cf),
		IsDeprecated:        cf.GetAttribute("Deprecated") != nil,
		Annotations:         annotationModelsFromClassFile(cf),
		RecordComponents:    recordComponentsFromClassFile(cf),
		PermittedSubclasses: permittedSubclassesFromClassFile(cf),
		NestHost:            nestHostFromClassFile(cf),
		NestMembers:         nestMembersFromClassFile(cf),
		EnclosingClass:      enclosingClassFromClassFile(cf),
		InnerClasses:        innerClassesFromClassFile(cf),
		Fields:              fieldModelsFromClassFile(cf),
		Methods:             methodModelsFromClassFile(cf),
	}
	model.IsSealed = len(model.PermittedSubclasses) > 0
	return model
}

func extractSimpleName(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return name
}

func extractPackage(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[:i]
		}
	}
	return ""
}

func superClassFromClassFile(cf *classfile.ClassFile) string {
	super := cf.SuperClassName()
	if super == "" {
		return ""
	}
	return classfile.InternalToSourceName(super)
}

func interfacesFromClassFile(cf *classfile.ClassFile) []string {
	internal := cf.InterfaceNames()
	result := make([]string, len(internal))
	for i, name := range internal {
		result[i] = classfile.InternalToSourceName(name)
	}
	return result
}

func visibilityFromClassAccessFlags(flags classfile.AccessFlags) Visibility {
	if flags.IsPublic() {
		return VisibilityPublic
	}
	return VisibilityPackage
}

func kindFromClassFile(cf *classfile.ClassFile) ClassKind {
	switch {
	case cf.IsAnnotation():
		return ClassKindAnnotation
	case cf.IsEnum():
		return ClassKindEnum
	case cf.IsModule():
		return ClassKindModule
	case cf.IsInterface():
		return ClassKindInterface
	case cf.GetAttribute("Record") != nil:
		return ClassKindRecord
	default:
		return ClassKindClass
	}
}

func signatureFromClassFile(cf *classfile.ClassFile) string {
	attr := cf.GetAttribute("Signature")
	if attr == nil {
		return ""
	}
	if sig := attr.AsSignature(); sig != nil {
		return cf.ConstantPool.GetUtf8(sig.SignatureIndex)
	}
	return ""
}

func sourceFileFromClassFile(cf *classfile.ClassFile) string {
	attr := cf.GetAttribute("SourceFile")
	if attr == nil {
		return ""
	}
	if sf := attr.AsSourceFile(); sf != nil {
		return cf.ConstantPool.GetUtf8(sf.SourceFileIndex)
	}
	return ""
}

func annotationModelsFromClassFile(cf *classfile.ClassFile) []AnnotationModel {
	attr := cf.GetAttribute("RuntimeVisibleAnnotations")
	if attr == nil {
		return nil
	}
	if rva := attr.AsRuntimeVisibleAnnotations(); rva != nil {
		return annotationModelsFromClassfileAnnotations(rva.Annotations, cf.ConstantPool)
	}
	return nil
}

func annotationModelsFromClassfileAnnotations(anns []classfile.Annotation, cp classfile.ConstantPool) []AnnotationModel {
	if len(anns) == 0 {
		return nil
	}
	result := make([]AnnotationModel, len(anns))
	for i, a := range anns {
		result[i] = AnnotationModel{
			Type:   descriptorToTypeName(cp.GetUtf8(a.TypeIndex)),
			Values: elementValuePairsToMap(a.ElementValuePairs, cp),
		}
	}
	return result
}

func elementValuePairsToMap(pairs []classfile.ElementValuePair, cp classfile.ConstantPool) map[string]interface{} {
	if len(pairs) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	for _, p := range pairs {
		result[cp.GetUtf8(p.ElementNameIndex)] = elementValueToGo(p.Value, cp)
	}
	return result
}

func recordComponentsFromClassFile(cf *classfile.ClassFile) []RecordComponentModel {
	attr := cf.GetAttribute("Record")
	if attr == nil {
		return nil
	}
	if rec := attr.AsRecord(); rec != nil {
		result := make([]RecordComponentModel, len(rec.Components))
		for i, comp := range rec.Components {
			descriptor := cf.ConstantPool.GetUtf8(comp.DescriptorIndex)
			ft := classfile.ParseFieldDescriptor(descriptor)
			result[i] = RecordComponentModel{
				Name: cf.ConstantPool.GetUtf8(comp.NameIndex),
				Type: typeModelFromFieldType(ft),
			}
		}
		return result
	}
	return nil
}

func permittedSubclassesFromClassFile(cf *classfile.ClassFile) []string {
	attr := cf.GetAttribute("PermittedSubclasses")
	if attr == nil {
		return nil
	}
	if ps := attr.AsPermittedSubclasses(); ps != nil {
		result := make([]string, len(ps.Classes))
		for i, idx := range ps.Classes {
			result[i] = classfile.InternalToSourceName(cf.ConstantPool.GetClassName(idx))
		}
		return result
	}
	return nil
}

func nestHostFromClassFile(cf *classfile.ClassFile) string {
	attr := cf.GetAttribute("NestHost")
	if attr == nil {
		return ""
	}
	if nh := attr.AsNestHost(); nh != nil {
		return classfile.InternalToSourceName(cf.ConstantPool.GetClassName(nh.HostClassIndex))
	}
	return ""
}

func nestMembersFromClassFile(cf *classfile.ClassFile) []string {
	attr := cf.GetAttribute("NestMembers")
	if attr == nil {
		return nil
	}
	if nm := attr.AsNestMembers(); nm != nil {
		result := make([]string, len(nm.Classes))
		for i, idx := range nm.Classes {
			result[i] = classfile.InternalToSourceName(cf.ConstantPool.GetClassName(idx))
		}
		return result
	}
	return nil
}

func enclosingClassFromClassFile(cf *classfile.ClassFile) string {
	attr := cf.GetAttribute("EnclosingMethod")
	if attr == nil {
		return ""
	}
	if em := attr.AsEnclosingMethod(); em != nil {
		return classfile.InternalToSourceName(cf.ConstantPool.GetClassName(em.ClassIndex))
	}
	return ""
}

func innerClassesFromClassFile(cf *classfile.ClassFile) []InnerClassModel {
	attr := cf.GetAttribute("InnerClasses")
	if attr == nil {
		return nil
	}
	if ic := attr.AsInnerClasses(); ic != nil {
		result := make([]InnerClassModel, len(ic.Classes))
		for i, entry := range ic.Classes {
			result[i] = InnerClassModel{
				InnerClass: classfile.InternalToSourceName(cf.ConstantPool.GetClassName(entry.InnerClassInfoIndex)),
				OuterClass: classfile.InternalToSourceName(cf.ConstantPool.GetClassName(entry.OuterClassInfoIndex)),
				InnerName:  cf.ConstantPool.GetUtf8(entry.InnerNameIndex),
				Visibility: visibilityFromInnerClassAccessFlags(entry.InnerClassAccessFlags),
				IsStatic:   entry.InnerClassAccessFlags.IsStatic(),
				IsFinal:    entry.InnerClassAccessFlags.IsFinal(),
				IsAbstract: entry.InnerClassAccessFlags.IsAbstract(),
			}
		}
		return result
	}
	return nil
}

func visibilityFromInnerClassAccessFlags(flags classfile.AccessFlags) Visibility {
	if flags.IsPublic() {
		return VisibilityPublic
	}
	if flags.IsPrivate() {
		return VisibilityPrivate
	}
	if flags.IsProtected() {
		return VisibilityProtected
	}
	return VisibilityPackage
}

func fieldModelsFromClassFile(cf *classfile.ClassFile) []FieldModel {
	fields := make([]FieldModel, len(cf.Fields))
	for i := range cf.Fields {
		fields[i] = fieldModelFromFieldInfo(&cf.Fields[i], cf.ConstantPool)
	}
	return fields
}

func fieldModelFromFieldInfo(info *classfile.FieldInfo, cp classfile.ConstantPool) FieldModel {
	ft := info.ParsedDescriptor(cp)
	model := FieldModel{
		Name:         info.Name(cp),
		Type:         typeModelFromFieldType(ft),
		Visibility:   visibilityFromFieldInfo(info),
		IsStatic:     info.IsStatic(),
		IsFinal:      info.IsFinal(),
		IsVolatile:   info.IsVolatile(),
		IsTransient:  info.IsTransient(),
		IsSynthetic:  info.IsSynthetic(),
		IsEnum:       info.IsEnum(),
		Signature:    signatureFromFieldInfo(info, cp),
		IsDeprecated: info.GetAttribute(cp, "Deprecated") != nil,
		Annotations:  annotationModelsFromFieldInfo(info, cp),
	}
	model.ConstantValue = constantValueFromFieldInfo(info, cp)
	return model
}

func visibilityFromFieldInfo(info *classfile.FieldInfo) Visibility {
	if info.IsPublic() {
		return VisibilityPublic
	}
	if info.IsPrivate() {
		return VisibilityPrivate
	}
	if info.IsProtected() {
		return VisibilityProtected
	}
	return VisibilityPackage
}

func signatureFromFieldInfo(info *classfile.FieldInfo, cp classfile.ConstantPool) string {
	attr := info.GetAttribute(cp, "Signature")
	if attr == nil {
		return ""
	}
	if sig := attr.AsSignature(); sig != nil {
		return cp.GetUtf8(sig.SignatureIndex)
	}
	return ""
}

func annotationModelsFromFieldInfo(info *classfile.FieldInfo, cp classfile.ConstantPool) []AnnotationModel {
	attr := info.GetAttribute(cp, "RuntimeVisibleAnnotations")
	if attr == nil {
		return nil
	}
	if rva := attr.AsRuntimeVisibleAnnotations(); rva != nil {
		return annotationModelsFromClassfileAnnotations(rva.Annotations, cp)
	}
	return nil
}

func constantValueFromFieldInfo(info *classfile.FieldInfo, cp classfile.ConstantPool) interface{} {
	attr := info.GetAttribute(cp, "ConstantValue")
	if attr == nil {
		return nil
	}
	if cv := attr.AsConstantValue(); cv != nil {
		if val, ok := cp.GetInteger(cv.ConstantValueIndex); ok {
			return val
		}
		if val, ok := cp.GetLong(cv.ConstantValueIndex); ok {
			return val
		}
		if val, ok := cp.GetFloat(cv.ConstantValueIndex); ok {
			return val
		}
		if val, ok := cp.GetDouble(cv.ConstantValueIndex); ok {
			return val
		}
		if val := cp.GetString(cv.ConstantValueIndex); val != "" {
			return val
		}
	}
	return nil
}

func methodModelsFromClassFile(cf *classfile.ClassFile) []MethodModel {
	methods := make([]MethodModel, len(cf.Methods))
	for i := range cf.Methods {
		methods[i] = methodModelFromMethodInfo(&cf.Methods[i], cf.ConstantPool)
	}
	return methods
}

func methodModelFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool) MethodModel {
	desc := info.ParsedDescriptor(cp)
	model := MethodModel{
		Name:                 info.Name(cp),
		ReturnType:           returnTypeFromDescriptor(desc),
		Parameters:           parametersFromMethodInfo(info, cp, desc),
		Visibility:           visibilityFromMethodInfo(info),
		IsStatic:             info.IsStatic(),
		IsFinal:              info.IsFinal(),
		IsAbstract:           info.IsAbstract(),
		IsSynchronized:       info.IsSynchronized(),
		IsNative:             info.IsNative(),
		IsBridge:             info.IsBridge(),
		IsVarargs:            info.IsVarargs(),
		IsSynthetic:          info.IsSynthetic(),
		Signature:            signatureFromMethodInfo(info, cp),
		IsDeprecated:         info.GetAttribute(cp, "Deprecated") != nil,
		Annotations:          annotationModelsFromMethodInfo(info, cp),
		ParameterAnnotations: parameterAnnotationsFromMethodInfo(info, cp),
		Exceptions:           exceptionsFromMethodInfo(info, cp),
	}
	return model
}

func returnTypeFromDescriptor(desc *classfile.MethodDescriptor) TypeModel {
	if desc == nil {
		return TypeModel{Name: "void"}
	}
	return typeModelFromFieldType(desc.ReturnType)
}

func typeModelFromFieldType(ft *classfile.FieldType) TypeModel {
	if ft == nil {
		return TypeModel{Name: "void"}
	}
	name := ft.BaseType
	if name == "" {
		name = classfile.InternalToSourceName(ft.ClassName)
	}
	return TypeModel{
		Name:       name,
		ArrayDepth: ft.ArrayDepth,
	}
}

func parametersFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool, desc *classfile.MethodDescriptor) []ParameterModel {
	if desc == nil {
		return nil
	}
	params := make([]ParameterModel, len(desc.Parameters))
	for i, p := range desc.Parameters {
		params[i] = ParameterModel{
			Type: typeModelFromFieldType(&p),
		}
	}
	populateParameterNamesFromMethodInfo(info, cp, params)
	return params
}

func populateParameterNamesFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool, params []ParameterModel) {
	code := info.GetCodeAttribute(cp)
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
	if !info.IsStatic() {
		startSlot = 1
	}

	for i := range params {
		slot := startSlot + i
		for _, lv := range lvt.LocalVariableTable {
			if int(lv.Index) == slot && lv.StartPC == 0 {
				params[i].Name = cp.GetUtf8(lv.NameIndex)
				break
			}
		}
	}
}

func visibilityFromMethodInfo(info *classfile.MethodInfo) Visibility {
	if info.IsPublic() {
		return VisibilityPublic
	}
	if info.IsPrivate() {
		return VisibilityPrivate
	}
	if info.IsProtected() {
		return VisibilityProtected
	}
	return VisibilityPackage
}

func signatureFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool) string {
	attr := info.GetAttribute(cp, "Signature")
	if attr == nil {
		return ""
	}
	if sig := attr.AsSignature(); sig != nil {
		return cp.GetUtf8(sig.SignatureIndex)
	}
	return ""
}

func annotationModelsFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool) []AnnotationModel {
	attr := info.GetAttribute(cp, "RuntimeVisibleAnnotations")
	if attr == nil {
		return nil
	}
	if rva := attr.AsRuntimeVisibleAnnotations(); rva != nil {
		return annotationModelsFromClassfileAnnotations(rva.Annotations, cp)
	}
	return nil
}

func parameterAnnotationsFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool) [][]AnnotationModel {
	attr := info.GetAttribute(cp, "RuntimeVisibleParameterAnnotations")
	if attr == nil {
		return nil
	}
	if rvpa := attr.AsRuntimeVisibleParameterAnnotations(); rvpa != nil {
		result := make([][]AnnotationModel, len(rvpa.ParameterAnnotations))
		for i, anns := range rvpa.ParameterAnnotations {
			result[i] = annotationModelsFromClassfileAnnotations(anns, cp)
		}
		return result
	}
	return nil
}

func exceptionsFromMethodInfo(info *classfile.MethodInfo, cp classfile.ConstantPool) []string {
	attr := info.GetAttribute(cp, "Exceptions")
	if attr == nil {
		return nil
	}
	if ex := attr.AsExceptions(); ex != nil {
		result := make([]string, len(ex.ExceptionIndexTable))
		for i, idx := range ex.ExceptionIndexTable {
			result[i] = classfile.InternalToSourceName(cp.GetClassName(idx))
		}
		return result
	}
	return nil
}
