package v2

import (
	"sort"
	"strings"

	"github.com/dhamidi/sai/ebnf/lex"
	"github.com/dhamidi/sai/ebnf/parse"
	"github.com/dhamidi/sai/java"
)

// ClassModelsFromSource parses Java source and returns class models.
func ClassModelsFromSource(source []byte) ([]*java.ClassModel, error) {
	cst, comments, err := ParseBytes(source, "")
	if err != nil {
		return nil, err
	}
	if cst == nil {
		return nil, nil
	}

	return classModelsFromCompilationUnit(cst, comments), nil
}

// javadocFinder helps find Javadoc comments for declarations based on position.
type javadocFinder struct {
	comments []lex.Token
	used     map[int]bool
}

func newJavadocFinder(comments []lex.Token) *javadocFinder {
	var javadocs []lex.Token
	for _, c := range comments {
		if c.Kind == "Comment" && strings.HasPrefix(c.Literal, "/**") {
			javadocs = append(javadocs, c)
		}
	}
	sort.Slice(javadocs, func(i, j int) bool {
		return javadocs[i].Position.Line < javadocs[j].Position.Line
	})
	return &javadocFinder{comments: javadocs, used: make(map[int]bool)}
}

func (jf *javadocFinder) FindForNode(node *parse.Node) string {
	if jf == nil || len(jf.comments) == 0 {
		return ""
	}
	startLine := node.Span.Start.Line

	bestIdx := -1
	bestDistance := 100

	for i, c := range jf.comments {
		if jf.used[i] {
			continue
		}
		endLine := c.Position.Line + strings.Count(c.Literal, "\n")
		endCol := len(c.Literal)
		if nl := strings.LastIndex(c.Literal, "\n"); nl >= 0 {
			endCol = len(c.Literal) - nl
		}

		if endLine > startLine {
			continue
		}
		if endLine == startLine && endCol >= node.Span.Start.Column {
			continue
		}
		distance := startLine - endLine
		if distance < bestDistance {
			bestIdx = i
			bestDistance = distance
		}
	}

	if bestIdx >= 0 {
		jf.used[bestIdx] = true
		return jf.comments[bestIdx].Literal
	}
	return ""
}

func classModelsFromCompilationUnit(cu *parse.Node, comments []lex.Token) []*java.ClassModel {
	var models []*java.ClassModel
	pkg := packageFromCompilationUnit(cu)
	resolver := newTypeResolver(pkg, importsFromCompilationUnit(cu), nil)
	jf := newJavadocFinder(comments)

	walkTopLevelDeclarations(cu, func(child *parse.Node) {
		switch child.Kind {
		case "normalClassDeclaration":
			models = append(models, classModelFromNormalClassDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromClassDecl(child, pkg, resolver, jf)...)
		case "normalInterfaceDeclaration":
			models = append(models, classModelFromInterfaceDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromInterfaceDecl(child, pkg, resolver, jf)...)
		case "enumDeclaration":
			models = append(models, classModelFromEnumDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromEnumDecl(child, pkg, resolver, jf)...)
		case "recordDeclaration":
			models = append(models, classModelFromRecordDecl(child, pkg, resolver, jf))
			models = append(models, innerClassesFromRecordDecl(child, pkg, resolver, jf)...)
		case "annotationInterfaceDeclaration":
			models = append(models, classModelFromAnnotationDecl(child, pkg, resolver, jf))
		case "classDeclaration":
			for _, subchild := range child.Children {
				switch subchild.Kind {
				case "normalClassDeclaration":
					models = append(models, classModelFromNormalClassDecl(subchild, pkg, resolver, jf))
					models = append(models, innerClassesFromClassDecl(subchild, pkg, resolver, jf)...)
				case "enumDeclaration":
					models = append(models, classModelFromEnumDecl(subchild, pkg, resolver, jf))
					models = append(models, innerClassesFromEnumDecl(subchild, pkg, resolver, jf)...)
				case "recordDeclaration":
					models = append(models, classModelFromRecordDecl(subchild, pkg, resolver, jf))
					models = append(models, innerClassesFromRecordDecl(subchild, pkg, resolver, jf)...)
				}
			}
		case "interfaceDeclaration":
			for _, subchild := range child.Children {
				switch subchild.Kind {
				case "normalInterfaceDeclaration":
					models = append(models, classModelFromInterfaceDecl(subchild, pkg, resolver, jf))
					models = append(models, innerClassesFromInterfaceDecl(subchild, pkg, resolver, jf)...)
				case "annotationInterfaceDeclaration":
					models = append(models, classModelFromAnnotationDecl(subchild, pkg, resolver, jf))
				}
			}
		}
	})

	return models
}

func walkTopLevelDeclarations(cu *parse.Node, fn func(*parse.Node)) {
	for _, child := range cu.Children {
		switch child.Kind {
		case "ordinaryCompilationUnit":
			for _, subchild := range child.Children {
				switch subchild.Kind {
				case "topLevelClassOrInterfaceDeclaration":
					for _, decl := range subchild.Children {
						fn(decl)
					}
				}
			}
		case "topLevelClassOrInterfaceDeclaration":
			for _, decl := range child.Children {
				fn(decl)
			}
		case "normalClassDeclaration", "normalInterfaceDeclaration", "enumDeclaration",
			"recordDeclaration", "annotationInterfaceDeclaration", "classDeclaration", "interfaceDeclaration":
			fn(child)
		}
	}
}

func packageFromCompilationUnit(cu *parse.Node) string {
	var parts []string
	var walkPackageDecl func(*parse.Node)
	walkPackageDecl = func(node *parse.Node) {
		for _, child := range node.Children {
			switch child.Kind {
			case "packageDeclaration":
				walkPackageDecl(child)
			case "Identifier":
				if child.Token != nil {
					parts = append(parts, child.Token.Literal)
				}
			case "ordinaryCompilationUnit":
				walkPackageDecl(child)
			}
		}
	}
	walkPackageDecl(cu)
	return strings.Join(parts, ".")
}

type importInfo struct {
	qualifiedName string
	isStatic      bool
	isWildcard    bool
}

func importsFromCompilationUnit(cu *parse.Node) []importInfo {
	var imports []importInfo

	var walkImports func(*parse.Node)
	walkImports = func(node *parse.Node) {
		for _, child := range node.Children {
			switch child.Kind {
			case "ordinaryCompilationUnit", "modularCompilationUnit":
				walkImports(child)
			case "importDeclaration":
				imp := importFromDeclaration(child)
				imports = append(imports, imp)
			case "singleTypeImportDeclaration", "typeImportOnDemandDeclaration",
				"singleStaticImportDeclaration", "staticImportOnDemandDeclaration":
				imp := importFromDeclaration(child)
				imports = append(imports, imp)
			}
		}
	}
	walkImports(cu)

	return imports
}

func importFromDeclaration(node *parse.Node) importInfo {
	imp := importInfo{}
	var parts []string

	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		if n.Token != nil {
			switch n.Token.Literal {
			case "static":
				imp.isStatic = true
			case "*":
				imp.isWildcard = true
			case "import", ".", ";":
				// skip
			default:
				if n.Kind == "Identifier" {
					parts = append(parts, n.Token.Literal)
				}
			}
			return
		}
		for _, child := range n.Children {
			walk(child)
		}
	}
	walk(node)

	imp.qualifiedName = strings.Join(parts, ".")
	return imp
}

type typeResolver struct {
	pkg          string
	imports      []importInfo
	innerClasses map[string]string
	classes      []*java.ClassModel
}

func newTypeResolver(pkg string, imports []importInfo, classes []*java.ClassModel) *typeResolver {
	return &typeResolver{
		pkg:          pkg,
		imports:      imports,
		innerClasses: make(map[string]string),
		classes:      classes,
	}
}

func (r *typeResolver) registerInnerClass(simpleName, fullName string) {
	r.innerClasses[simpleName] = fullName
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

	if fullName, ok := r.innerClasses[simpleName]; ok {
		return fullName
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

	for _, imp := range r.imports {
		if !imp.isWildcard || imp.isStatic {
			continue
		}
		candidate := imp.qualifiedName + "." + simpleName
		for _, cls := range r.classes {
			if cls.Name == candidate {
				return candidate
			}
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

func classModelFromNormalClassDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindClass,
		Package:    pkg,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyClassModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			if pkg != "" {
				model.Name = pkg + "." + name
			} else {
				model.Name = name
			}
		case "typeParameters":
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case "classExtends":
			model.SuperClass = typeFromClassExtends(child, resolver)
		case "classImplements":
			model.Interfaces = typesFromImplements(child, resolver)
		case "classPermits":
			model.PermittedSubclasses = typesFromPermits(child, resolver)
		case "classBody":
			extractClassBodyMembers(child, model, resolver, jf)
		}
	}

	collectAndRegisterInnerClasses(node, model.Name, resolver)

	return model
}

func classModelFromInterfaceDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindInterface,
		Package:    pkg,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyInterfaceModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			if pkg != "" {
				model.Name = pkg + "." + name
			} else {
				model.Name = name
			}
		case "typeParameters":
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case "interfaceExtends":
			model.Interfaces = typesFromInterfaceExtends(child, resolver)
		case "interfacePermits":
			model.PermittedSubclasses = typesFromPermits(child, resolver)
		case "interfaceBody":
			extractInterfaceBodyMembers(child, model, resolver, jf)
		}
	}

	return model
}

func classModelFromEnumDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindEnum,
		Package:    pkg,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyClassModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			if pkg != "" {
				model.Name = pkg + "." + name
			} else {
				model.Name = name
			}
		case "classImplements":
			model.Interfaces = typesFromImplements(child, resolver)
		case "enumBody":
			extractEnumBodyMembers(child, model, resolver, jf)
		}
	}

	return model
}

func classModelFromRecordDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindRecord,
		Package:    pkg,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyClassModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			if pkg != "" {
				model.Name = pkg + "." + name
			} else {
				model.Name = name
			}
		case "typeParameters":
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case "recordHeader":
			model.RecordComponents = recordComponentsFromHeader(child, resolver)
		case "classImplements":
			model.Interfaces = typesFromImplements(child, resolver)
		case "recordBody":
			extractRecordBodyMembers(child, model, resolver, jf)
		}
	}

	return model
}

func classModelFromAnnotationDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindAnnotation,
		Package:    pkg,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyInterfaceModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			if pkg != "" {
				model.Name = pkg + "." + name
			} else {
				model.Name = name
			}
		case "annotationInterfaceBody":
			extractAnnotationBodyMembers(child, model, resolver, jf)
		}
	}

	return model
}

func applyClassModifiers(node *parse.Node, model *java.ClassModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "classModifier":
			applyClassModifier(child, model, resolver)
		case "annotation":
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyClassModifier(node *parse.Node, model *java.ClassModel, resolver *typeResolver) {
	for _, child := range node.Children {
		if child.Token != nil {
			switch child.Token.Literal {
			case "public":
				model.Visibility = java.VisibilityPublic
			case "protected":
				model.Visibility = java.VisibilityProtected
			case "private":
				model.Visibility = java.VisibilityPrivate
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
		if child.Kind == "annotation" {
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyInterfaceModifiers(node *parse.Node, model *java.ClassModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "interfaceModifier":
			applyClassModifier(child, model, resolver)
		case "annotation":
			model.Annotations = append(model.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func extractClassBodyMembers(body *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range body.Children {
		switch child.Kind {
		case "classBodyDeclaration":
			extractClassBodyDeclaration(child, model, resolver, jf)
		case "classMemberDeclaration":
			extractClassMemberDeclaration(child, model, resolver, jf)
		case "fieldDeclaration":
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver, jf)...)
		case "methodDeclaration":
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case "constructorDeclaration":
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		case "classDeclaration", "interfaceDeclaration":
			inner := innerClassFromDeclaration(child, model.Name, resolver, jf)
			if inner != nil {
				model.InnerClasses = append(model.InnerClasses, java.InnerClassModel{
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
}

func extractClassBodyDeclaration(node *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range node.Children {
		switch child.Kind {
		case "classMemberDeclaration":
			extractClassMemberDeclaration(child, model, resolver, jf)
		case "constructorDeclaration":
			model.Methods = append(model.Methods, methodModelFromConstructorDecl(child, model.SimpleName, resolver, jf))
		}
	}
}

func extractClassMemberDeclaration(node *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range node.Children {
		switch child.Kind {
		case "fieldDeclaration":
			model.Fields = append(model.Fields, fieldModelsFromFieldDecl(child, resolver, jf)...)
		case "methodDeclaration":
			model.Methods = append(model.Methods, methodModelFromMethodDecl(child, resolver, jf))
		case "classDeclaration":
			inner := innerClassFromDeclaration(child, model.Name, resolver, jf)
			if inner != nil {
				model.InnerClasses = append(model.InnerClasses, java.InnerClassModel{
					InnerClass: inner.Name,
					OuterClass: model.Name,
					InnerName:  inner.SimpleName,
					Visibility: inner.Visibility,
					IsStatic:   inner.IsStatic,
					IsFinal:    inner.IsFinal,
					IsAbstract: inner.IsAbstract,
				})
			}
		case "interfaceDeclaration":
			inner := innerClassFromDeclaration(child, model.Name, resolver, jf)
			if inner != nil {
				model.InnerClasses = append(model.InnerClasses, java.InnerClassModel{
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
}

func extractInterfaceBodyMembers(body *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range body.Children {
		switch child.Kind {
		case "interfaceMemberDeclaration":
			extractInterfaceMemberDeclaration(child, model, resolver, jf)
		case "constantDeclaration":
			model.Fields = append(model.Fields, fieldModelsFromConstantDecl(child, resolver, jf)...)
		case "interfaceMethodDeclaration":
			model.Methods = append(model.Methods, methodModelFromInterfaceMethodDecl(child, resolver, jf))
		}
	}
}

func extractInterfaceMemberDeclaration(node *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range node.Children {
		switch child.Kind {
		case "constantDeclaration":
			model.Fields = append(model.Fields, fieldModelsFromConstantDecl(child, resolver, jf)...)
		case "interfaceMethodDeclaration":
			model.Methods = append(model.Methods, methodModelFromInterfaceMethodDecl(child, resolver, jf))
		case "classDeclaration", "interfaceDeclaration":
			inner := innerClassFromDeclaration(child, model.Name, resolver, jf)
			if inner != nil {
				model.InnerClasses = append(model.InnerClasses, java.InnerClassModel{
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
}

func extractEnumBodyMembers(body *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range body.Children {
		switch child.Kind {
		case "enumConstantList":
			model.EnumConstants = enumConstantsFromList(child)
		case "enumBodyDeclarations":
			extractEnumBodyDeclarations(child, model, resolver, jf)
		}
	}
}

func enumConstantsFromList(node *parse.Node) []java.EnumConstantModel {
	var constants []java.EnumConstantModel
	for _, child := range node.Children {
		if child.Kind == "enumConstant" {
			constants = append(constants, enumConstantFromNode(child))
		}
	}
	return constants
}

func enumConstantFromNode(node *parse.Node) java.EnumConstantModel {
	ec := java.EnumConstantModel{}
	for _, child := range node.Children {
		if child.Kind == "Identifier" && child.Token != nil {
			ec.Name = child.Token.Literal
		}
	}
	return ec
}

func extractEnumBodyDeclarations(node *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range node.Children {
		if child.Kind == "classBodyDeclaration" {
			extractClassBodyDeclaration(child, model, resolver, jf)
		}
	}
}

func extractRecordBodyMembers(body *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range body.Children {
		switch child.Kind {
		case "recordBodyDeclaration":
			extractRecordBodyDeclaration(child, model, resolver, jf)
		case "classBodyDeclaration":
			extractClassBodyDeclaration(child, model, resolver, jf)
		}
	}
}

func extractRecordBodyDeclaration(node *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range node.Children {
		if child.Kind == "classBodyDeclaration" {
			extractClassBodyDeclaration(child, model, resolver, jf)
		}
	}
}

func extractAnnotationBodyMembers(body *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range body.Children {
		switch child.Kind {
		case "annotationInterfaceMemberDeclaration":
			extractAnnotationMemberDeclaration(child, model, resolver, jf)
		}
	}
}

func extractAnnotationMemberDeclaration(node *parse.Node, model *java.ClassModel, resolver *typeResolver, jf *javadocFinder) {
	for _, child := range node.Children {
		switch child.Kind {
		case "annotationInterfaceElementDeclaration":
			model.Methods = append(model.Methods, methodModelFromAnnotationElement(child, resolver, jf))
		case "constantDeclaration":
			model.Fields = append(model.Fields, fieldModelsFromConstantDecl(child, resolver, jf)...)
		case "classDeclaration", "interfaceDeclaration":
			inner := innerClassFromDeclaration(child, model.Name, resolver, jf)
			if inner != nil {
				model.InnerClasses = append(model.InnerClasses, java.InnerClassModel{
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
}

func innerClassFromDeclaration(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "normalClassDeclaration":
			return classModelFromNormalClassDeclNested(child, outerName, resolver, jf)
		case "enumDeclaration":
			return classModelFromEnumDeclNested(child, outerName, resolver, jf)
		case "recordDeclaration":
			return classModelFromRecordDeclNested(child, outerName, resolver, jf)
		case "normalInterfaceDeclaration":
			return classModelFromInterfaceDeclNested(child, outerName, resolver, jf)
		case "annotationInterfaceDeclaration":
			return classModelFromAnnotationDeclNested(child, outerName, resolver, jf)
		}
	}
	return nil
}

func classModelFromNormalClassDeclNested(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindClass,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyClassModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			model.Name = outerName + "." + name
		case "typeParameters":
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case "classExtends":
			model.SuperClass = typeFromClassExtends(child, resolver)
		case "classImplements":
			model.Interfaces = typesFromImplements(child, resolver)
		case "classBody":
			extractClassBodyMembers(child, model, resolver, jf)
		}
	}

	setPackageFromOuterName(model, outerName)
	return model
}

func classModelFromEnumDeclNested(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindEnum,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyClassModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			model.Name = outerName + "." + name
		case "classImplements":
			model.Interfaces = typesFromImplements(child, resolver)
		case "enumBody":
			extractEnumBodyMembers(child, model, resolver, jf)
		}
	}

	setPackageFromOuterName(model, outerName)
	return model
}

func classModelFromRecordDeclNested(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindRecord,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyClassModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			model.Name = outerName + "." + name
		case "typeParameters":
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case "recordHeader":
			model.RecordComponents = recordComponentsFromHeader(child, resolver)
		case "classImplements":
			model.Interfaces = typesFromImplements(child, resolver)
		case "recordBody":
			extractRecordBodyMembers(child, model, resolver, jf)
		}
	}

	setPackageFromOuterName(model, outerName)
	return model
}

func classModelFromInterfaceDeclNested(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindInterface,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyInterfaceModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			model.Name = outerName + "." + name
		case "typeParameters":
			model.TypeParameters = typeParametersFromNode(child, resolver)
		case "interfaceExtends":
			model.Interfaces = typesFromInterfaceExtends(child, resolver)
		case "interfaceBody":
			extractInterfaceBodyMembers(child, model, resolver, jf)
		}
	}

	setPackageFromOuterName(model, outerName)
	return model
}

func classModelFromAnnotationDeclNested(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) *java.ClassModel {
	model := &java.ClassModel{
		Kind:       java.ClassKindAnnotation,
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyInterfaceModifiers(node, model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			name := identifierText(child)
			model.SimpleName = name
			model.Name = outerName + "." + name
		case "annotationInterfaceBody":
			extractAnnotationBodyMembers(child, model, resolver, jf)
		}
	}

	setPackageFromOuterName(model, outerName)
	return model
}

func setPackageFromOuterName(model *java.ClassModel, outerName string) {
	if idx := strings.Index(outerName, "."); idx != -1 {
		parts := strings.Split(outerName, ".")
		for i, part := range parts {
			if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
				model.Package = strings.Join(parts[:i], ".")
				break
			}
		}
	}
}

func fieldModelsFromFieldDecl(node *parse.Node, resolver *typeResolver, jf *javadocFinder) []java.FieldModel {
	var fields []java.FieldModel
	baseField := java.FieldModel{
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyFieldModifiers(node, &baseField, resolver)

	var fieldType java.TypeModel
	for _, child := range node.Children {
		if child.Kind == "unannType" {
			fieldType = typeModelFromUnannType(child, resolver)
			break
		}
	}

	for _, child := range node.Children {
		if child.Kind == "variableDeclaratorList" {
			for _, varDecl := range child.Children {
				if varDecl.Kind == "variableDeclarator" {
					name := variableDeclaratorName(varDecl)
					if name != "" {
						field := baseField
						field.Name = name
						field.Type = fieldType
						fields = append(fields, field)
					}
				}
			}
		}
	}

	return fields
}

func fieldModelsFromConstantDecl(node *parse.Node, resolver *typeResolver, jf *javadocFinder) []java.FieldModel {
	var fields []java.FieldModel
	baseField := java.FieldModel{
		Visibility: java.VisibilityPublic,
		IsStatic:   true,
		IsFinal:    true,
		Javadoc:    jf.FindForNode(node),
	}

	applyConstantModifiers(node, &baseField, resolver)

	var fieldType java.TypeModel
	for _, child := range node.Children {
		if child.Kind == "unannType" {
			fieldType = typeModelFromUnannType(child, resolver)
			break
		}
	}

	for _, child := range node.Children {
		if child.Kind == "variableDeclaratorList" {
			for _, varDecl := range child.Children {
				if varDecl.Kind == "variableDeclarator" {
					name := variableDeclaratorName(varDecl)
					if name != "" {
						field := baseField
						field.Name = name
						field.Type = fieldType
						fields = append(fields, field)
					}
				}
			}
		}
	}

	return fields
}

func applyFieldModifiers(node *parse.Node, field *java.FieldModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "fieldModifier":
			applyFieldModifier(child, field, resolver)
		case "annotation":
			field.Annotations = append(field.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyFieldModifier(node *parse.Node, field *java.FieldModel, resolver *typeResolver) {
	for _, child := range node.Children {
		if child.Token != nil {
			switch child.Token.Literal {
			case "public":
				field.Visibility = java.VisibilityPublic
			case "protected":
				field.Visibility = java.VisibilityProtected
			case "private":
				field.Visibility = java.VisibilityPrivate
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
		if child.Kind == "annotation" {
			field.Annotations = append(field.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyConstantModifiers(node *parse.Node, field *java.FieldModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "constantModifier":
			applyFieldModifier(child, field, resolver)
		case "annotation":
			field.Annotations = append(field.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func variableDeclaratorName(node *parse.Node) string {
	for _, child := range node.Children {
		if child.Kind == "variableDeclaratorId" {
			return identifierText(child)
		}
		if child.Kind == "Identifier" && child.Token != nil {
			return child.Token.Literal
		}
	}
	return ""
}

func methodModelFromMethodDecl(node *parse.Node, resolver *typeResolver, jf *javadocFinder) java.MethodModel {
	model := java.MethodModel{
		Visibility: java.VisibilityPackage,
		Javadoc:    jf.FindForNode(node),
	}

	applyMethodModifiers(node, &model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "methodHeader":
			extractMethodHeader(child, &model, resolver)
		case "methodBody":
			// body is not needed for the model
		}
	}

	return model
}

func methodModelFromInterfaceMethodDecl(node *parse.Node, resolver *typeResolver, jf *javadocFinder) java.MethodModel {
	model := java.MethodModel{
		Visibility: java.VisibilityPublic,
		Javadoc:    jf.FindForNode(node),
	}

	applyInterfaceMethodModifiers(node, &model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "methodHeader":
			extractMethodHeader(child, &model, resolver)
		}
	}

	return model
}

func methodModelFromConstructorDecl(node *parse.Node, className string, resolver *typeResolver, jf *javadocFinder) java.MethodModel {
	model := java.MethodModel{
		Name:       "<init>",
		Visibility: java.VisibilityPackage,
		ReturnType: java.TypeModel{Name: "void"},
		Javadoc:    jf.FindForNode(node),
	}

	applyConstructorModifiers(node, &model, resolver)

	for _, child := range node.Children {
		switch child.Kind {
		case "constructorDeclarator":
			extractConstructorDeclarator(child, &model, resolver)
		case "throws":
			model.Exceptions = exceptionsFromThrows(child, resolver)
		}
	}

	return model
}

func methodModelFromAnnotationElement(node *parse.Node, resolver *typeResolver, jf *javadocFinder) java.MethodModel {
	model := java.MethodModel{
		Visibility: java.VisibilityPublic,
		IsAbstract: true,
		Javadoc:    jf.FindForNode(node),
	}

	for _, child := range node.Children {
		switch child.Kind {
		case "unannType":
			model.ReturnType = typeModelFromUnannType(child, resolver)
		case "Identifier":
			if child.Token != nil {
				model.Name = child.Token.Literal
			}
		}
	}

	return model
}

func applyMethodModifiers(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "methodModifier":
			applyMethodModifier(child, method, resolver)
		case "annotation":
			method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyMethodModifier(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		if child.Token != nil {
			switch child.Token.Literal {
			case "public":
				method.Visibility = java.VisibilityPublic
			case "protected":
				method.Visibility = java.VisibilityProtected
			case "private":
				method.Visibility = java.VisibilityPrivate
			case "abstract":
				method.IsAbstract = true
			case "static":
				method.IsStatic = true
			case "final":
				method.IsFinal = true
			case "synchronized":
				method.IsSynchronized = true
			case "native":
				method.IsNative = true
			}
		}
		if child.Kind == "annotation" {
			method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyInterfaceMethodModifiers(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "interfaceMethodModifier":
			applyInterfaceMethodModifier(child, method, resolver)
		case "annotation":
			method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyInterfaceMethodModifier(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		if child.Token != nil {
			switch child.Token.Literal {
			case "public":
				method.Visibility = java.VisibilityPublic
			case "private":
				method.Visibility = java.VisibilityPrivate
			case "abstract":
				method.IsAbstract = true
			case "default":
				method.IsDefault = true
			case "static":
				method.IsStatic = true
			}
		}
		if child.Kind == "annotation" {
			method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyConstructorModifiers(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "constructorModifier":
			applyConstructorModifier(child, method, resolver)
		case "annotation":
			method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func applyConstructorModifier(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		if child.Token != nil {
			switch child.Token.Literal {
			case "public":
				method.Visibility = java.VisibilityPublic
			case "protected":
				method.Visibility = java.VisibilityProtected
			case "private":
				method.Visibility = java.VisibilityPrivate
			}
		}
		if child.Kind == "annotation" {
			method.Annotations = append(method.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func extractMethodHeader(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "result":
			method.ReturnType = typeModelFromResult(child, resolver)
		case "methodDeclarator":
			extractMethodDeclarator(child, method, resolver)
		case "typeParameters":
			method.TypeParameters = typeParametersFromNode(child, resolver)
		case "throws":
			method.Exceptions = exceptionsFromThrows(child, resolver)
		}
	}
}

func extractMethodDeclarator(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "Identifier":
			if child.Token != nil {
				method.Name = child.Token.Literal
			}
		case "formalParameterList":
			method.Parameters = parametersFromFormalParameterList(child, resolver)
		}
	}
}

func extractConstructorDeclarator(node *parse.Node, method *java.MethodModel, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "typeParameters":
			method.TypeParameters = typeParametersFromNode(child, resolver)
		case "formalParameterList":
			method.Parameters = parametersFromFormalParameterList(child, resolver)
		}
	}
}

func parametersFromFormalParameterList(node *parse.Node, resolver *typeResolver) []java.ParameterModel {
	var params []java.ParameterModel
	for _, child := range node.Children {
		if child.Kind == "formalParameter" {
			params = append(params, parameterFromFormalParameter(child, resolver))
		} else if child.Kind == "variableArityParameter" {
			params = append(params, parameterFromVariableArityParameter(child, resolver))
		}
	}
	return params
}

func parameterFromFormalParameter(node *parse.Node, resolver *typeResolver) java.ParameterModel {
	param := java.ParameterModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "variableModifier":
			applyVariableModifier(child, &param, resolver)
		case "unannType":
			param.Type = typeModelFromUnannType(child, resolver)
		case "variableDeclaratorId":
			param.Name = identifierText(child)
		}
	}

	return param
}

func parameterFromVariableArityParameter(node *parse.Node, resolver *typeResolver) java.ParameterModel {
	param := java.ParameterModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "variableModifier":
			applyVariableModifier(child, &param, resolver)
		case "unannType":
			param.Type = typeModelFromUnannType(child, resolver)
			param.Type.ArrayDepth++ // varargs are arrays
		case "Identifier":
			if child.Token != nil {
				param.Name = child.Token.Literal
			}
		}
	}

	return param
}

func applyVariableModifier(node *parse.Node, param *java.ParameterModel, resolver *typeResolver) {
	for _, child := range node.Children {
		if child.Token != nil && child.Token.Literal == "final" {
			param.IsFinal = true
		}
		if child.Kind == "annotation" {
			param.Annotations = append(param.Annotations, annotationModelFromNode(child, resolver))
		}
	}
}

func exceptionsFromThrows(node *parse.Node, resolver *typeResolver) []string {
	var exceptions []string
	for _, child := range node.Children {
		if child.Kind == "exceptionTypeList" {
			for _, exc := range child.Children {
				if exc.Kind == "exceptionType" {
					exceptions = append(exceptions, exceptionTypeText(exc, resolver))
				}
			}
		}
	}
	return exceptions
}

func exceptionTypeText(node *parse.Node, resolver *typeResolver) string {
	for _, child := range node.Children {
		if child.Kind == "classType" || child.Kind == "typeVariable" {
			return resolver.resolve(classTypeText(child))
		}
	}
	return ""
}

func typeModelFromResult(node *parse.Node, resolver *typeResolver) java.TypeModel {
	for _, child := range node.Children {
		if child.Kind == "unannType" {
			return typeModelFromUnannType(child, resolver)
		}
		if child.Token != nil && child.Token.Literal == "void" {
			return java.TypeModel{Name: "void"}
		}
	}
	return java.TypeModel{Name: "void"}
}

func typeModelFromUnannType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "unannPrimitiveType":
			return typeModelFromUnannPrimitiveType(child, resolver)
		case "unannReferenceType":
			return typeModelFromUnannReferenceType(child, resolver)
		}
	}
	return java.TypeModel{}
}

func typeModelFromUnannPrimitiveType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "numericType":
			return typeModelFromNumericType(child)
		}
		if child.Token != nil && child.Token.Literal == "boolean" {
			return java.TypeModel{Name: "boolean"}
		}
	}
	return java.TypeModel{}
}

func typeModelFromNumericType(node *parse.Node) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "integralType", "floatingPointType":
			if child.Token != nil {
				return java.TypeModel{Name: child.Token.Literal}
			}
			for _, subchild := range child.Children {
				if subchild.Token != nil {
					return java.TypeModel{Name: subchild.Token.Literal}
				}
			}
		}
		if child.Token != nil {
			return java.TypeModel{Name: child.Token.Literal}
		}
	}
	return java.TypeModel{}
}

func typeModelFromUnannReferenceType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "unannClassOrInterfaceType":
			return typeModelFromUnannClassOrInterfaceType(child, resolver)
		case "unannTypeVariable":
			return java.TypeModel{Name: identifierText(child)}
		case "unannArrayType":
			return typeModelFromUnannArrayType(child, resolver)
		}
	}
	return java.TypeModel{}
}

func typeModelFromUnannClassOrInterfaceType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "unannClassType", "unannInterfaceType":
			return typeModelFromUnannClassType(child, resolver)
		}
	}
	return java.TypeModel{}
}

func typeModelFromUnannClassType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	model := java.TypeModel{}
	var parts []string

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			parts = append(parts, identifierText(child))
		case "packageName":
			parts = append(parts, packageNameText(child))
		case "typeArguments":
			model.TypeArguments = typeArgumentsFromNode(child, resolver)
		case "unannClassOrInterfaceType":
			inner := typeModelFromUnannClassOrInterfaceType(child, resolver)
			parts = append(parts, inner.Name)
		}
	}

	if len(parts) > 0 {
		model.Name = resolver.resolve(strings.Join(parts, "."))
	}
	return model
}

func typeModelFromUnannArrayType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	model := java.TypeModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "unannPrimitiveType":
			inner := typeModelFromUnannPrimitiveType(child, resolver)
			model.Name = inner.Name
		case "unannClassOrInterfaceType":
			inner := typeModelFromUnannClassOrInterfaceType(child, resolver)
			model.Name = inner.Name
			model.TypeArguments = inner.TypeArguments
		case "unannTypeVariable":
			model.Name = identifierText(child)
		case "dims":
			model.ArrayDepth = countDims(child)
		}
	}

	return model
}

func countDims(node *parse.Node) int {
	count := 0
	for _, child := range node.Children {
		if child.Token != nil && child.Token.Literal == "[" {
			count++
		}
	}
	if count == 0 {
		count = 1
	}
	return count
}

func typeParametersFromNode(node *parse.Node, resolver *typeResolver) []java.TypeParameterModel {
	var params []java.TypeParameterModel
	for _, child := range node.Children {
		if child.Kind == "typeParameterList" {
			for _, tp := range child.Children {
				if tp.Kind == "typeParameter" {
					params = append(params, typeParameterFromNode(tp, resolver))
				}
			}
		}
	}
	return params
}

func typeParameterFromNode(node *parse.Node, resolver *typeResolver) java.TypeParameterModel {
	param := java.TypeParameterModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "typeIdentifier":
			param.Name = identifierText(child)
		case "typeBound":
			param.Bounds = boundsFromTypeBound(child, resolver)
		}
	}

	return param
}

func boundsFromTypeBound(node *parse.Node, resolver *typeResolver) []java.TypeModel {
	var bounds []java.TypeModel
	for _, child := range node.Children {
		switch child.Kind {
		case "typeVariable":
			bounds = append(bounds, java.TypeModel{Name: identifierText(child)})
		case "classOrInterfaceType":
			bounds = append(bounds, typeModelFromClassOrInterfaceType(child, resolver))
		case "additionalBound":
			for _, sub := range child.Children {
				if sub.Kind == "interfaceType" {
					bounds = append(bounds, typeModelFromClassOrInterfaceType(sub, resolver))
				}
			}
		}
	}
	return bounds
}

func typeModelFromClassOrInterfaceType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	model := java.TypeModel{}
	var parts []string

	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		for _, child := range n.Children {
			switch child.Kind {
			case "classType", "interfaceType":
				walk(child)
			case "typeIdentifier":
				parts = append(parts, identifierText(child))
			case "packageName":
				parts = append(parts, packageNameText(child))
			case "typeArguments":
				model.TypeArguments = typeArgumentsFromNode(child, resolver)
			}
		}
	}
	walk(node)

	if len(parts) > 0 {
		model.Name = resolver.resolve(strings.Join(parts, "."))
	}
	return model
}

func typeArgumentsFromNode(node *parse.Node, resolver *typeResolver) []java.TypeArgumentModel {
	var args []java.TypeArgumentModel
	for _, child := range node.Children {
		if child.Kind == "typeArgumentList" {
			for _, ta := range child.Children {
				if ta.Kind == "typeArgument" {
					args = append(args, typeArgumentFromNode(ta, resolver))
				}
			}
		}
	}
	return args
}

func typeArgumentFromNode(node *parse.Node, resolver *typeResolver) java.TypeArgumentModel {
	arg := java.TypeArgumentModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "referenceType":
			tm := typeModelFromReferenceType(child, resolver)
			arg.Type = &tm
		case "wildcard":
			arg.IsWildcard = true
			for _, wc := range child.Children {
				if wc.Kind == "wildcardBounds" {
					for _, bound := range wc.Children {
						if bound.Token != nil {
							switch bound.Token.Literal {
							case "extends":
								arg.BoundKind = "extends"
							case "super":
								arg.BoundKind = "super"
							}
						}
						if bound.Kind == "referenceType" {
							tm := typeModelFromReferenceType(bound, resolver)
							arg.Bound = &tm
						}
					}
				}
			}
		}
	}

	return arg
}

func typeModelFromReferenceType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "classOrInterfaceType":
			return typeModelFromClassOrInterfaceType(child, resolver)
		case "typeVariable":
			return java.TypeModel{Name: identifierText(child)}
		case "arrayType":
			return typeModelFromArrayType(child, resolver)
		}
	}
	return java.TypeModel{}
}

func typeModelFromArrayType(node *parse.Node, resolver *typeResolver) java.TypeModel {
	model := java.TypeModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "primitiveType":
			inner := typeModelFromPrimitiveType(child)
			model.Name = inner.Name
		case "classOrInterfaceType":
			inner := typeModelFromClassOrInterfaceType(child, resolver)
			model.Name = inner.Name
			model.TypeArguments = inner.TypeArguments
		case "typeVariable":
			model.Name = identifierText(child)
		case "dims":
			model.ArrayDepth = countDims(child)
		}
	}

	return model
}

func typeModelFromPrimitiveType(node *parse.Node) java.TypeModel {
	for _, child := range node.Children {
		switch child.Kind {
		case "numericType":
			return typeModelFromNumericType(child)
		}
		if child.Token != nil && child.Token.Literal == "boolean" {
			return java.TypeModel{Name: "boolean"}
		}
	}
	return java.TypeModel{}
}

func recordComponentsFromHeader(node *parse.Node, resolver *typeResolver) []java.RecordComponentModel {
	var components []java.RecordComponentModel
	for _, child := range node.Children {
		if child.Kind == "recordComponentList" {
			for _, comp := range child.Children {
				if comp.Kind == "recordComponent" {
					components = append(components, recordComponentFromNode(comp, resolver))
				}
			}
		}
	}
	return components
}

func recordComponentFromNode(node *parse.Node, resolver *typeResolver) java.RecordComponentModel {
	comp := java.RecordComponentModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "unannType":
			comp.Type = typeModelFromUnannType(child, resolver)
		case "Identifier":
			if child.Token != nil {
				comp.Name = child.Token.Literal
			}
		case "recordComponentModifier":
			for _, mod := range child.Children {
				if mod.Kind == "annotation" {
					comp.Annotations = append(comp.Annotations, annotationModelFromNode(mod, resolver))
				}
			}
		}
	}

	return comp
}

func annotationModelFromNode(node *parse.Node, resolver *typeResolver) java.AnnotationModel {
	ann := java.AnnotationModel{}

	for _, child := range node.Children {
		switch child.Kind {
		case "normalAnnotation":
			ann = annotationModelFromNormalAnnotation(child, resolver)
		case "markerAnnotation":
			ann = annotationModelFromMarkerAnnotation(child, resolver)
		case "singleElementAnnotation":
			ann = annotationModelFromSingleElementAnnotation(child, resolver)
		case "typeName":
			ann.Type = resolver.resolve(typeNameText(child))
		}
	}

	return ann
}

func annotationModelFromNormalAnnotation(node *parse.Node, resolver *typeResolver) java.AnnotationModel {
	ann := java.AnnotationModel{Values: make(map[string]interface{})}

	for _, child := range node.Children {
		switch child.Kind {
		case "typeName":
			ann.Type = resolver.resolve(typeNameText(child))
		case "elementValuePairList":
			for _, pair := range child.Children {
				if pair.Kind == "elementValuePair" {
					name, value := elementValuePairFromNode(pair, resolver)
					ann.Values[name] = value
				}
			}
		}
	}

	return ann
}

func annotationModelFromMarkerAnnotation(node *parse.Node, resolver *typeResolver) java.AnnotationModel {
	ann := java.AnnotationModel{}

	for _, child := range node.Children {
		if child.Kind == "typeName" {
			ann.Type = resolver.resolve(typeNameText(child))
		}
	}

	return ann
}

func annotationModelFromSingleElementAnnotation(node *parse.Node, resolver *typeResolver) java.AnnotationModel {
	ann := java.AnnotationModel{Values: make(map[string]interface{})}

	for _, child := range node.Children {
		switch child.Kind {
		case "typeName":
			ann.Type = resolver.resolve(typeNameText(child))
		case "elementValue":
			ann.Values["value"] = elementValueFromNode(child, resolver)
		}
	}

	return ann
}

func elementValuePairFromNode(node *parse.Node, resolver *typeResolver) (string, interface{}) {
	name := "value"
	var value interface{}

	for _, child := range node.Children {
		switch child.Kind {
		case "Identifier":
			if child.Token != nil {
				name = child.Token.Literal
			}
		case "elementValue":
			value = elementValueFromNode(child, resolver)
		}
	}

	return name, value
}

func elementValueFromNode(node *parse.Node, resolver *typeResolver) interface{} {
	for _, child := range node.Children {
		switch child.Kind {
		case "conditionalExpression":
			return conditionalExpressionText(child)
		case "elementValueArrayInitializer":
			return elementValueArrayFromNode(child, resolver)
		case "annotation":
			return annotationModelFromNode(child, resolver)
		}
	}
	return nil
}

func conditionalExpressionText(node *parse.Node) string {
	var parts []string
	collectTokens(node, &parts)
	return strings.Join(parts, "")
}

func collectTokens(node *parse.Node, parts *[]string) {
	if node.Token != nil {
		*parts = append(*parts, node.Token.Literal)
		return
	}
	for _, child := range node.Children {
		collectTokens(child, parts)
	}
}

func elementValueArrayFromNode(node *parse.Node, resolver *typeResolver) []interface{} {
	var values []interface{}
	for _, child := range node.Children {
		if child.Kind == "elementValueList" {
			for _, ev := range child.Children {
				if ev.Kind == "elementValue" {
					values = append(values, elementValueFromNode(ev, resolver))
				}
			}
		}
	}
	return values
}

func typeFromClassExtends(node *parse.Node, resolver *typeResolver) string {
	for _, child := range node.Children {
		if child.Kind == "classType" {
			return resolver.resolve(classTypeText(child))
		}
	}
	return ""
}

func typesFromImplements(node *parse.Node, resolver *typeResolver) []string {
	var types []string
	for _, child := range node.Children {
		if child.Kind == "interfaceTypeList" {
			for _, iface := range child.Children {
				if iface.Kind == "interfaceType" {
					types = append(types, resolver.resolve(classTypeText(iface)))
				}
			}
		}
	}
	return types
}

func typesFromInterfaceExtends(node *parse.Node, resolver *typeResolver) []string {
	var types []string
	for _, child := range node.Children {
		if child.Kind == "interfaceTypeList" {
			for _, iface := range child.Children {
				if iface.Kind == "interfaceType" {
					types = append(types, resolver.resolve(classTypeText(iface)))
				}
			}
		}
	}
	return types
}

func typesFromPermits(node *parse.Node, resolver *typeResolver) []string {
	var types []string
	for _, child := range node.Children {
		if child.Kind == "typeName" {
			types = append(types, resolver.resolve(typeNameText(child)))
		}
	}
	return types
}

func classTypeText(node *parse.Node) string {
	var parts []string
	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		for _, child := range n.Children {
			switch child.Kind {
			case "classType", "interfaceType":
				walk(child)
			case "typeIdentifier":
				parts = append(parts, identifierText(child))
			case "packageName":
				parts = append(parts, packageNameText(child))
			case "Identifier":
				if child.Token != nil {
					parts = append(parts, child.Token.Literal)
				}
			}
		}
	}
	walk(node)
	return strings.Join(parts, ".")
}

func typeNameText(node *parse.Node) string {
	var parts []string
	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		for _, child := range n.Children {
			switch child.Kind {
			case "typeIdentifier", "packageOrTypeName":
				walk(child)
			case "Identifier":
				if child.Token != nil {
					parts = append(parts, child.Token.Literal)
				}
			}
		}
	}
	walk(node)
	return strings.Join(parts, ".")
}

func packageNameText(node *parse.Node) string {
	var parts []string
	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		for _, child := range n.Children {
			if child.Kind == "Identifier" && child.Token != nil {
				parts = append(parts, child.Token.Literal)
			} else if child.Kind == "packageName" {
				walk(child)
			}
		}
	}
	walk(node)
	return strings.Join(parts, ".")
}

func identifierText(node *parse.Node) string {
	if node.Token != nil {
		return node.Token.Literal
	}
	for _, child := range node.Children {
		if child.Kind == "Identifier" && child.Token != nil {
			return child.Token.Literal
		}
		if text := identifierText(child); text != "" {
			return text
		}
	}
	return ""
}

func innerClassesFromClassDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*java.ClassModel {
	var name string
	for _, child := range node.Children {
		if child.Kind == "typeIdentifier" {
			name = identifierText(child)
			if pkg != "" {
				name = pkg + "." + name
			}
			break
		}
	}
	if name == "" {
		return nil
	}
	return collectInnerClasses(node, name, resolver, jf)
}

func innerClassesFromInterfaceDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*java.ClassModel {
	var name string
	for _, child := range node.Children {
		if child.Kind == "typeIdentifier" {
			name = identifierText(child)
			if pkg != "" {
				name = pkg + "." + name
			}
			break
		}
	}
	if name == "" {
		return nil
	}
	return collectInnerClasses(node, name, resolver, jf)
}

func innerClassesFromEnumDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*java.ClassModel {
	var name string
	for _, child := range node.Children {
		if child.Kind == "typeIdentifier" {
			name = identifierText(child)
			if pkg != "" {
				name = pkg + "." + name
			}
			break
		}
	}
	if name == "" {
		return nil
	}
	return collectInnerClasses(node, name, resolver, jf)
}

func innerClassesFromRecordDecl(node *parse.Node, pkg string, resolver *typeResolver, jf *javadocFinder) []*java.ClassModel {
	var name string
	for _, child := range node.Children {
		if child.Kind == "typeIdentifier" {
			name = identifierText(child)
			if pkg != "" {
				name = pkg + "." + name
			}
			break
		}
	}
	if name == "" {
		return nil
	}
	return collectInnerClasses(node, name, resolver, jf)
}

func collectInnerClasses(node *parse.Node, outerName string, resolver *typeResolver, jf *javadocFinder) []*java.ClassModel {
	var models []*java.ClassModel

	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		for _, child := range n.Children {
			switch child.Kind {
			case "classBody", "interfaceBody", "enumBody", "recordBody", "enumBodyDeclarations":
				for _, member := range child.Children {
					switch member.Kind {
					case "classBodyDeclaration":
						for _, sub := range member.Children {
							if sub.Kind == "classMemberDeclaration" {
								for _, decl := range sub.Children {
									if inner := innerClassFromDeclaration(decl, outerName, resolver, jf); inner != nil {
										models = append(models, inner)
										models = append(models, collectInnerClasses(decl, inner.Name, resolver, jf)...)
									}
								}
							}
						}
					case "interfaceMemberDeclaration":
						for _, decl := range member.Children {
							if inner := innerClassFromDeclaration(decl, outerName, resolver, jf); inner != nil {
								models = append(models, inner)
								models = append(models, collectInnerClasses(decl, inner.Name, resolver, jf)...)
							}
						}
					case "classMemberDeclaration":
						for _, decl := range member.Children {
							if inner := innerClassFromDeclaration(decl, outerName, resolver, jf); inner != nil {
								models = append(models, inner)
								models = append(models, collectInnerClasses(decl, inner.Name, resolver, jf)...)
							}
						}
					}
				}
			}
		}
	}
	walk(node)

	return models
}

func collectAndRegisterInnerClasses(node *parse.Node, outerName string, resolver *typeResolver) {
	var walk func(*parse.Node)
	walk = func(n *parse.Node) {
		for _, child := range n.Children {
			switch child.Kind {
			case "classBody", "interfaceBody", "enumBody", "recordBody":
				for _, member := range child.Children {
					switch member.Kind {
					case "classBodyDeclaration", "interfaceMemberDeclaration", "classMemberDeclaration":
						for _, sub := range member.Children {
							registerInnerClassFromDecl(sub, outerName, resolver)
						}
					}
				}
			}
		}
	}
	walk(node)
}

func registerInnerClassFromDecl(node *parse.Node, outerName string, resolver *typeResolver) {
	for _, child := range node.Children {
		switch child.Kind {
		case "normalClassDeclaration", "enumDeclaration", "recordDeclaration",
			"normalInterfaceDeclaration", "annotationInterfaceDeclaration":
			name := ""
			for _, sub := range child.Children {
				if sub.Kind == "typeIdentifier" {
					name = identifierText(sub)
					break
				}
			}
			if name != "" {
				fullName := outerName + "." + name
				resolver.registerInnerClass(name, fullName)
				collectAndRegisterInnerClasses(child, fullName, resolver)
			}
		case "classDeclaration", "interfaceDeclaration":
			registerInnerClassFromDecl(child, outerName, resolver)
		case "classMemberDeclaration":
			registerInnerClassFromDecl(child, outerName, resolver)
		}
	}
}
