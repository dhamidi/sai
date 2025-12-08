package java

import (
	"strings"

	"github.com/dhamidi/javalyzer/java/parser"
)

func ClassModelsFromSource(source []byte, opts ...parser.Option) ([]*ClassModel, error) {
	p := parser.ParseCompilationUnit(opts...)
	p.Push(source)
	node := p.Finish()
	if node == nil {
		return nil, nil
	}
	return classModelsFromCompilationUnit(node), nil
}

func classModelsFromCompilationUnit(cu *parser.Node) []*ClassModel {
	var models []*ClassModel
	pkg := packageFromCompilationUnit(cu)
	resolver := newTypeResolver(pkg, importsFromCompilationUnit(cu))

	for _, child := range cu.Children {
		switch child.Kind {
		case parser.KindClassDecl:
			models = append(models, classModelFromClassDecl(child, pkg, resolver))
		case parser.KindInterfaceDecl:
			models = append(models, classModelFromInterfaceDecl(child, pkg, resolver))
		case parser.KindEnumDecl:
			models = append(models, classModelFromEnumDecl(child, pkg, resolver))
		case parser.KindRecordDecl:
			models = append(models, classModelFromRecordDecl(child, pkg, resolver))
		case parser.KindAnnotationDecl:
			models = append(models, classModelFromAnnotationDecl(child, pkg, resolver))
		}
	}
	return models
}

func packageFromCompilationUnit(cu *parser.Node) string {
	pkgDecl := cu.FirstChildOfKind(parser.KindPackageDecl)
	if pkgDecl == nil {
		return ""
	}
	qn := pkgDecl.FirstChildOfKind(parser.KindQualifiedName)
	if qn == nil {
		return ""
	}
	return qualifiedNameToString(qn)
}

func qualifiedNameToString(qn *parser.Node) string {
	var parts []string
	for _, child := range qn.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			parts = append(parts, child.Token.Literal)
		}
	}
	return strings.Join(parts, ".")
}

type importInfo struct {
	qualifiedName string
	isStatic      bool
	isWildcard    bool
}

func importsFromCompilationUnit(cu *parser.Node) []importInfo {
	var imports []importInfo
	for _, child := range cu.Children {
		if child.Kind == parser.KindImportDecl {
			imp := importInfo{}
			for _, ic := range child.Children {
				if ic.Kind == parser.KindIdentifier && ic.Token != nil {
					if ic.Token.Literal == "static" {
						imp.isStatic = true
					} else if ic.Token.Literal == "*" {
						imp.isWildcard = true
					}
				} else if ic.Kind == parser.KindQualifiedName {
					imp.qualifiedName = qualifiedNameToString(ic)
				}
			}
			imports = append(imports, imp)
		}
	}
	return imports
}

type typeResolver struct {
	pkg     string
	imports []importInfo
}

func newTypeResolver(pkg string, imports []importInfo) *typeResolver {
	return &typeResolver{pkg: pkg, imports: imports}
}

var javaLangTypes = map[string]bool{
	"Object": true, "String": true, "Class": true, "System": true,
	"Throwable": true, "Exception": true, "RuntimeException": true, "Error": true,
	"Integer": true, "Long": true, "Short": true, "Byte": true,
	"Float": true, "Double": true, "Character": true, "Boolean": true,
	"Number": true, "Comparable": true, "CharSequence": true,
	"Iterable": true, "Cloneable": true, "Runnable": true,
	"Thread": true, "StringBuilder": true, "StringBuffer": true,
	"Math": true, "Enum": true, "Record": true,
	"Override": true, "Deprecated": true, "SuppressWarnings": true, "FunctionalInterface": true,
}

func (r *typeResolver) resolve(simpleName string) string {
	if simpleName == "" {
		return ""
	}

	if strings.Contains(simpleName, ".") {
		return simpleName
	}

	switch simpleName {
	case "boolean", "byte", "char", "short", "int", "long", "float", "double", "void":
		return simpleName
	}

	for _, imp := range r.imports {
		if imp.isWildcard || imp.isStatic {
			continue
		}
		parts := strings.Split(imp.qualifiedName, ".")
		if len(parts) > 0 && parts[len(parts)-1] == simpleName {
			return imp.qualifiedName
		}
	}

	if javaLangTypes[simpleName] {
		return "java.lang." + simpleName
	}

	if r.pkg != "" {
		return r.pkg + "." + simpleName
	}

	return simpleName
}

func classModelFromClassDecl(node *parser.Node, pkg string, resolver *typeResolver) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindClass,
		Package:    pkg,
		Visibility: VisibilityPackage,
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				model.SimpleName = child.Token.Literal
				if pkg != "" {
					model.Name = pkg + "." + model.SimpleName
				} else {
					model.Name = model.SimpleName
				}
			}
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindType:
			if model.SuperClass == "" {
				model.SuperClass = typeModelFromTypeNode(child, resolver).Name
			} else {
				model.Interfaces = append(model.Interfaces, typeModelFromTypeNode(child, resolver).Name)
			}
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver)
		case parser.KindFieldDecl:
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver)...)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			inner := classModelFromClassDeclNested(child, model.Name, resolver)
			model.InnerClasses = append(model.InnerClasses, InnerClassModel{
				InnerClass: inner.Name,
				OuterClass: model.Name,
				InnerName:  inner.SimpleName,
				Visibility: inner.Visibility,
				IsStatic:   inner.IsStatic,
				IsFinal:    inner.IsFinal,
				IsAbstract: inner.IsAbstract,
			})
		}
	}

	return model
}

func extractClassBodyMembers(block *parser.Node, model *ClassModel, resolver *typeResolver) {
	for _, child := range block.Children {
		switch child.Kind {
		case parser.KindFieldDecl:
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver)...)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver))
		case parser.KindClassDecl, parser.KindInterfaceDecl, parser.KindEnumDecl, parser.KindRecordDecl:
			inner := classModelFromClassDeclNested(child, model.Name, resolver)
			model.InnerClasses = append(model.InnerClasses, InnerClassModel{
				InnerClass: inner.Name,
				OuterClass: model.Name,
				InnerName:  inner.SimpleName,
				Visibility: inner.Visibility,
				IsStatic:   inner.IsStatic,
				IsFinal:    inner.IsFinal,
				IsAbstract: inner.IsAbstract,
			})
		}
	}
}

func classModelFromClassDeclNested(node *parser.Node, outerName string, resolver *typeResolver) *ClassModel {
	model := &ClassModel{
		Visibility: VisibilityPackage,
	}

	switch node.Kind {
	case parser.KindClassDecl:
		model.Kind = ClassKindClass
	case parser.KindInterfaceDecl:
		model.Kind = ClassKindInterface
	case parser.KindEnumDecl:
		model.Kind = ClassKindEnum
	case parser.KindRecordDecl:
		model.Kind = ClassKindRecord
	case parser.KindAnnotationDecl:
		model.Kind = ClassKindAnnotation
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			model.SimpleName = child.Token.Literal
			model.Name = outerName + "." + model.SimpleName
			break
		}
	}

	return model
}

func classModelFromInterfaceDecl(node *parser.Node, pkg string, resolver *typeResolver) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindInterface,
		Package:    pkg,
		Visibility: VisibilityPackage,
		IsAbstract: true,
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				model.SimpleName = child.Token.Literal
				if pkg != "" {
					model.Name = pkg + "." + model.SimpleName
				} else {
					model.Name = model.SimpleName
				}
			}
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindType:
			model.Interfaces = append(model.Interfaces, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		case parser.KindBlock:
			extractInterfaceBodyMembers(child, model, resolver)
		case parser.KindFieldDecl:
			fields := fieldModelsFromFieldDecl(child, resolver)
			for i := range fields {
				fields[i].IsStatic = true
				fields[i].IsFinal = true
				fields[i].Visibility = VisibilityPublic
			}
			model.Fields = append(model.Fields, fields...)
		case parser.KindMethodDecl:
			method := methodModelFromMethodDecl(child, resolver)
			if !method.IsStatic && !method.IsDefault {
				method.IsAbstract = true
			}
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func extractInterfaceBodyMembers(block *parser.Node, model *ClassModel, resolver *typeResolver) {
	for _, child := range block.Children {
		switch child.Kind {
		case parser.KindFieldDecl:
			fields := fieldModelsFromFieldDecl(child, resolver)
			for i := range fields {
				fields[i].IsStatic = true
				fields[i].IsFinal = true
				fields[i].Visibility = VisibilityPublic
			}
			model.Fields = append(model.Fields, fields...)
		case parser.KindMethodDecl:
			method := methodModelFromMethodDecl(child, resolver)
			if !method.IsStatic && !method.IsDefault {
				method.IsAbstract = true
			}
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		}
	}
}

func classModelFromEnumDecl(node *parser.Node, pkg string, resolver *typeResolver) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindEnum,
		Package:    pkg,
		Visibility: VisibilityPackage,
		SuperClass: "java.lang.Enum",
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				model.SimpleName = child.Token.Literal
				if pkg != "" {
					model.Name = pkg + "." + model.SimpleName
				} else {
					model.Name = model.SimpleName
				}
			}
		case parser.KindType:
			model.Interfaces = append(model.Interfaces, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver)
		case parser.KindFieldDecl:
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver)...)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func classModelFromRecordDecl(node *parser.Node, pkg string, resolver *typeResolver) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindRecord,
		Package:    pkg,
		Visibility: VisibilityPackage,
		SuperClass: "java.lang.Record",
		IsFinal:    true,
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				model.SimpleName = child.Token.Literal
				if pkg != "" {
					model.Name = pkg + "." + model.SimpleName
				} else {
					model.Name = model.SimpleName
				}
			}
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			model.RecordComponents = recordComponentsFromParameters(child, resolver)
		case parser.KindType:
			model.Interfaces = append(model.Interfaces, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		case parser.KindBlock:
			extractClassBodyMembers(child, model, resolver)
		case parser.KindMethodDecl:
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver))
		case parser.KindConstructorDecl:
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver))
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func classModelFromAnnotationDecl(node *parser.Node, pkg string, resolver *typeResolver) *ClassModel {
	model := &ClassModel{
		Kind:       ClassKindAnnotation,
		Package:    pkg,
		Visibility: VisibilityPackage,
		IsAbstract: true,
		Interfaces: []string{"java.lang.annotation.Annotation"},
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToClass(modifiers, model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				model.SimpleName = child.Token.Literal
				if pkg != "" {
					model.Name = pkg + "." + model.SimpleName
				} else {
					model.Name = model.SimpleName
				}
			}
		case parser.KindBlock:
			extractAnnotationBodyMembers(child, model, resolver)
		case parser.KindMethodDecl:
			method := methodModelFromMethodDecl(child, resolver)
			method.IsAbstract = true
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		case parser.KindAnnotation:
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}

	return model
}

func extractAnnotationBodyMembers(block *parser.Node, model *ClassModel, resolver *typeResolver) {
	for _, child := range block.Children {
		if child.Kind == parser.KindMethodDecl {
			method := methodModelFromMethodDecl(child, resolver)
			method.IsAbstract = true
			method.Visibility = VisibilityPublic
			model.Methods = append(model.Methods, method)
		}
	}
}

func applyModifiersToClass(modifiers *parser.Node, model *ClassModel, resolver *typeResolver) {
	for _, child := range modifiers.Children {
		if child.Token == nil {
			if child.Kind == parser.KindAnnotation {
				model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
			}
			continue
		}
		switch child.Token.Literal {
		case "public":
			model.Visibility = VisibilityPublic
		case "protected":
			model.Visibility = VisibilityProtected
		case "private":
			model.Visibility = VisibilityPrivate
		case "abstract":
			model.IsAbstract = true
		case "static":
			model.IsStatic = true
		case "final":
			model.IsFinal = true
		case "sealed":
			model.IsSealed = true
		}
	}
}

func typeParametersFromNode(node *parser.Node, resolver *typeResolver) []TypeParameterModel {
	var params []TypeParameterModel
	for _, child := range node.Children {
		if child.Kind == parser.KindTypeParameter {
			params = append(params, typeParameterFromNode(child, resolver))
		}
	}
	return params
}

func typeParameterFromNode(node *parser.Node, resolver *typeResolver) TypeParameterModel {
	param := TypeParameterModel{}
	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				param.Name = child.Token.Literal
			}
		case parser.KindType:
			param.Bounds = append(param.Bounds, typeModelFromTypeNode(child, resolver))
		}
	}
	return param
}

func typeNameFromTypeNode(node *parser.Node, resolver *typeResolver) string {
	tm := typeModelFromTypeNode(node, resolver)
	return tm.Name
}

func typeModelFromTypeNode(node *parser.Node, resolver *typeResolver) TypeModel {
	model := TypeModel{}

	if node.Token != nil {
		model.Name = node.Token.Literal
		if resolver != nil {
			model.Name = resolver.resolve(model.Name)
		}
		return model
	}

	if node.Kind == parser.KindArrayType {
		for _, ac := range node.Children {
			if ac.Kind == parser.KindType || ac.Kind == parser.KindQualifiedName || ac.Kind == parser.KindIdentifier {
				inner := typeModelFromTypeNode(ac, resolver)
				model.Name = inner.Name
				model.ArrayDepth = inner.ArrayDepth + 1
				model.TypeArguments = inner.TypeArguments
				return model
			}
		}
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindIdentifier:
			if child.Token != nil {
				if model.Name != "" {
					model.Name += "."
				}
				model.Name += child.Token.Literal
			}
		case parser.KindQualifiedName:
			model.Name = qualifiedNameToString(child)
		case parser.KindArrayType:
			for _, ac := range child.Children {
				if ac.Kind == parser.KindType || ac.Kind == parser.KindQualifiedName || ac.Kind == parser.KindIdentifier {
					inner := typeModelFromTypeNode(ac, resolver)
					model.Name = inner.Name
					model.ArrayDepth = inner.ArrayDepth + 1
					model.TypeArguments = inner.TypeArguments
					break
				}
			}
		case parser.KindParameterizedType:
			for _, pc := range child.Children {
				switch pc.Kind {
				case parser.KindQualifiedName:
					model.Name = qualifiedNameToString(pc)
				case parser.KindIdentifier:
					if pc.Token != nil {
						if model.Name != "" {
							model.Name += "."
						}
						model.Name += pc.Token.Literal
					}
				case parser.KindTypeArguments:
					model.TypeArguments = typeArgumentsFromNode(pc, resolver)
				}
			}
		case parser.KindType:
			inner := typeModelFromTypeNode(child, resolver)
			model.Name = inner.Name
			model.ArrayDepth = inner.ArrayDepth
			model.TypeArguments = inner.TypeArguments
		}
	}

	if resolver != nil {
		model.Name = resolver.resolve(model.Name)
	}
	return model
}

func typeArgumentsFromNode(node *parser.Node, resolver *typeResolver) []TypeArgumentModel {
	var args []TypeArgumentModel
	for _, child := range node.Children {
		if child.Kind == parser.KindTypeArgument || child.Kind == parser.KindType {
			args = append(args, typeArgumentFromNode(child, resolver))
		}
	}
	return args
}

func typeArgumentFromNode(node *parser.Node, resolver *typeResolver) TypeArgumentModel {
	arg := TypeArgumentModel{}

	if node.Kind == parser.KindType {
		tm := typeModelFromTypeNode(node, resolver)
		arg.Type = &tm
		return arg
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			tm := typeModelFromTypeNode(child, resolver)
			if arg.BoundKind != "" {
				arg.Bound = &tm
			} else {
				arg.Type = &tm
			}
		case parser.KindWildcard:
			arg.IsWildcard = true
			for _, wc := range child.Children {
				if wc.Token != nil {
					switch wc.Token.Literal {
					case "extends":
						arg.BoundKind = "extends"
					case "super":
						arg.BoundKind = "super"
					}
				}
				if wc.Kind == parser.KindType {
					tm := typeModelFromTypeNode(wc, resolver)
					arg.Bound = &tm
				}
			}
		}
	}
	return arg
}

func fieldModelsFromFieldDecl(node *parser.Node, resolver *typeResolver) []FieldModel {
	var fields []FieldModel
	baseField := FieldModel{
		Visibility: VisibilityPackage,
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToField(modifiers, &baseField, resolver)
	}

	var fieldType TypeModel
	for _, child := range node.Children {
		if child.Kind == parser.KindType {
			fieldType = typeModelFromTypeNode(child, resolver)
			break
		}
	}

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			field := baseField
			field.Name = child.Token.Literal
			field.Type = fieldType
			fields = append(fields, field)
		}
	}

	return fields
}

func applyModifiersToField(modifiers *parser.Node, field *FieldModel, resolver *typeResolver) {
	for _, child := range modifiers.Children {
		if child.Token == nil {
			if child.Kind == parser.KindAnnotation {
				field.Annotations = append(field.Annotations, annotationModelFromNode(child, resolver))
			}
			continue
		}
		switch child.Token.Literal {
		case "public":
			field.Visibility = VisibilityPublic
		case "protected":
			field.Visibility = VisibilityProtected
		case "private":
			field.Visibility = VisibilityPrivate
		case "static":
			field.IsStatic = true
		case "final":
			field.IsFinal = true
		case "volatile":
			field.IsVolatile = true
		case "transient":
			field.IsTransient = true
		}
	}
}

func methodModelFromMethodDecl(node *parser.Node, resolver *typeResolver) MethodModel {
	model := MethodModel{
		Visibility: VisibilityPackage,
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToMethod(modifiers, &model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			model.ReturnType = typeModelFromTypeNode(child, resolver)
		case parser.KindIdentifier:
			if child.Token != nil {
				model.Name = child.Token.Literal
			}
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			model.Parameters = parametersFromNode(child, resolver)
		case parser.KindThrowsList:
			model.Exceptions = exceptionsFromThrowsList(child, resolver)
		}
	}

	return model
}

func methodModelFromConstructorDecl(node *parser.Node, className string, resolver *typeResolver) MethodModel {
	model := MethodModel{
		Name:       "<init>",
		Visibility: VisibilityPackage,
		ReturnType: TypeModel{Name: "void"},
	}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		applyModifiersToMethod(modifiers, &model, resolver)
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindTypeParameters:
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case parser.KindParameters:
			model.Parameters = parametersFromNode(child, resolver)
		case parser.KindThrowsList:
			model.Exceptions = exceptionsFromThrowsList(child, resolver)
		}
	}

	return model
}

func applyModifiersToMethod(modifiers *parser.Node, method *MethodModel, resolver *typeResolver) {
	for _, child := range modifiers.Children {
		if child.Token == nil {
			if child.Kind == parser.KindAnnotation {
				method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
			}
			continue
		}
		switch child.Token.Literal {
		case "public":
			method.Visibility = VisibilityPublic
		case "protected":
			method.Visibility = VisibilityProtected
		case "private":
			method.Visibility = VisibilityPrivate
		case "static":
			method.IsStatic = true
		case "final":
			method.IsFinal = true
		case "abstract":
			method.IsAbstract = true
		case "synchronized":
			method.IsSynchronized = true
		case "native":
			method.IsNative = true
		case "default":
			method.IsDefault = true
		}
	}
}

func parametersFromNode(node *parser.Node, resolver *typeResolver) []ParameterModel {
	var params []ParameterModel
	for _, child := range node.Children {
		if child.Kind == parser.KindParameter {
			params = append(params, parameterFromNode(child, resolver))
		}
	}
	return params
}

func parameterFromNode(node *parser.Node, resolver *typeResolver) ParameterModel {
	param := ParameterModel{}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		for _, child := range modifiers.Children {
			if child.Token != nil && child.Token.Literal == "final" {
				param.IsFinal = true
			}
			if child.Kind == parser.KindAnnotation {
				param.Annotations = append(param.Annotations, annotationModelFromNode(child, resolver))
			}
		}
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			param.Type = typeModelFromTypeNode(child, resolver)
		case parser.KindArrayType:
			param.Type = typeModelFromTypeNode(child, resolver)
		case parser.KindIdentifier:
			if child.Token != nil {
				param.Name = child.Token.Literal
			}
		}
	}

	return param
}

func exceptionsFromThrowsList(node *parser.Node, resolver *typeResolver) []string {
	var exceptions []string
	for _, child := range node.Children {
		if child.Kind == parser.KindType {
			exceptions = append(exceptions, resolver.resolve(typeNameFromTypeNode(child, resolver)))
		}
	}
	return exceptions
}

func recordComponentsFromParameters(node *parser.Node, resolver *typeResolver) []RecordComponentModel {
	var components []RecordComponentModel
	for _, child := range node.Children {
		if child.Kind == parser.KindParameter {
			comp := recordComponentFromParameter(child, resolver)
			components = append(components, comp)
		}
	}
	return components
}

func recordComponentFromParameter(node *parser.Node, resolver *typeResolver) RecordComponentModel {
	comp := RecordComponentModel{}

	modifiers := node.FirstChildOfKind(parser.KindModifiers)
	if modifiers != nil {
		for _, child := range modifiers.Children {
			if child.Kind == parser.KindAnnotation {
				comp.Annotations = append(comp.Annotations, annotationModelFromNode(child, resolver))
			}
		}
	}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindType:
			comp.Type = typeModelFromTypeNode(child, resolver)
		case parser.KindIdentifier:
			if child.Token != nil {
				comp.Name = child.Token.Literal
			}
		}
	}

	return comp
}

func annotationModelFromNode(node *parser.Node, resolver *typeResolver) AnnotationModel {
	ann := AnnotationModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case parser.KindQualifiedName:
			ann.Type = qualifiedNameToString(child)
		case parser.KindIdentifier:
			if child.Token != nil {
				ann.Type = child.Token.Literal
			}
		case parser.KindAnnotationElement:
			if ann.Values == nil {
				ann.Values = make(map[string]interface{})
			}
			name, value := annotationElementFromNode(child)
			ann.Values[name] = value
		}
	}

	if resolver != nil {
		ann.Type = resolver.resolve(ann.Type)
	}

	return ann
}

func annotationElementFromNode(node *parser.Node) (string, interface{}) {
	name := "value"
	var value interface{}

	for _, child := range node.Children {
		if child.Kind == parser.KindIdentifier && child.Token != nil {
			name = child.Token.Literal
		} else {
			value = annotationValueFromNode(child)
		}
	}

	return name, value
}

func annotationValueFromNode(node *parser.Node) interface{} {
	if node.Kind == parser.KindLiteral && node.Token != nil {
		return node.Token.Literal
	}
	if node.Kind == parser.KindIdentifier && node.Token != nil {
		return node.Token.Literal
	}
	if node.Kind == parser.KindAnnotation {
		return annotationModelFromNode(node, nil)
	}
	if node.Kind == parser.KindArrayInit {
		var values []interface{}
		for _, child := range node.Children {
			values = append(values, annotationValueFromNode(child))
		}
		return values
	}
	if node.Kind == parser.KindFieldAccess {
		var parts []string
		for _, child := range node.Children {
			if child.Token != nil {
				parts = append(parts, child.Token.Literal)
			}
		}
		return strings.Join(parts, ".")
	}
	return nil
}
