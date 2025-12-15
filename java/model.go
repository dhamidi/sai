package java

type Visibility string

const (
	VisibilityPublic    Visibility = "public"
	VisibilityProtected Visibility = "protected"
	VisibilityPrivate   Visibility = "private"
	VisibilityPackage   Visibility = "package"
)

type ClassKind string

const (
	ClassKindClass      ClassKind = "class"
	ClassKindInterface  ClassKind = "interface"
	ClassKindEnum       ClassKind = "enum"
	ClassKindAnnotation ClassKind = "annotation"
	ClassKindRecord     ClassKind = "record"
	ClassKindModule     ClassKind = "module"
)

type ClassModel struct {
	Name                string
	SimpleName          string
	Package             string
	Module              string
	SuperClass          string
	Interfaces          []string
	Visibility          Visibility
	Kind                ClassKind
	IsFinal             bool
	IsAbstract          bool
	IsStatic            bool
	IsSynthetic         bool
	IsSealed            bool
	MajorVersion        uint16
	MinorVersion        uint16
	Signature           string
	SourceFile          string
	SourceURL           URLString
	IsDeprecated        bool
	Javadoc             string
	Annotations         []AnnotationModel
	RecordComponents    []RecordComponentModel
	PermittedSubclasses []string
	NestHost            string
	NestMembers         []string
	EnclosingClass      string
	InnerClasses        []InnerClassModel
	EnumConstants       []EnumConstantModel
	Fields              []FieldModel
	Methods             []MethodModel
	TypeParameters      []TypeParameterModel
	Initializers        []InitializerModel
}

// InitializerModel represents a static or instance initializer block
type InitializerModel struct {
	IsStatic bool
	Body     string // The block body as source code (including braces)
}

type EnumConstantModel struct {
	Name      string
	Arguments []string
}

type FieldModel struct {
	Name          string
	Type          TypeModel
	Visibility    Visibility
	IsStatic      bool
	IsFinal       bool
	IsVolatile    bool
	IsTransient   bool
	IsSynthetic   bool
	IsEnum        bool
	Signature     string
	IsDeprecated  bool
	Javadoc       string
	Annotations   []AnnotationModel
	ConstantValue interface{}
}

type MethodModel struct {
	Name                 string
	ReturnType           TypeModel
	Parameters           []ParameterModel
	Visibility           Visibility
	IsStatic             bool
	IsFinal              bool
	IsAbstract           bool
	IsSynchronized       bool
	IsNative             bool
	IsBridge             bool
	IsVarargs            bool
	IsSynthetic          bool
	IsDefault            bool
	Signature            string
	IsDeprecated         bool
	Javadoc              string
	Annotations          []AnnotationModel
	ParameterAnnotations [][]AnnotationModel
	Exceptions           []string
	TypeParameters       []TypeParameterModel
}

type ParameterModel struct {
	Name        string
	Type        TypeModel
	IsFinal     bool
	Annotations []AnnotationModel
}

type TypeModel struct {
	Name           string
	ArrayDepth     int
	TypeArguments  []TypeArgumentModel
	TypeParameters []TypeParameterModel
}

func (t TypeModel) IsPrimitive() bool {
	if t.ArrayDepth > 0 {
		return false
	}
	switch t.Name {
	case "boolean", "byte", "char", "short", "int", "long", "float", "double":
		return true
	}
	return false
}

func (t TypeModel) IsArray() bool {
	return t.ArrayDepth > 0
}

func (t TypeModel) IsVoid() bool {
	return t.Name == "void" && t.ArrayDepth == 0
}

type TypeArgumentModel struct {
	Type       *TypeModel
	IsWildcard bool
	BoundKind  string // "extends", "super", or "" for unbounded
	Bound      *TypeModel
}

type TypeParameterModel struct {
	Name   string
	Bounds []TypeModel
}

type AnnotationModel struct {
	Type   string
	Values map[string]interface{}
}

type ElementValuePairModel struct {
	Name  string
	Value interface{}
}

type RecordComponentModel struct {
	Name        string
	Type        TypeModel
	Annotations []AnnotationModel
}

type InnerClassModel struct {
	InnerClass string
	OuterClass string
	InnerName  string
	Visibility Visibility
	IsStatic   bool
	IsFinal    bool
	IsAbstract bool
}

type ModuleModel struct {
	Name        string
	IsOpen      bool
	SourceFile  string
	SourceURL   URLString
	Javadoc     string
	Annotations []AnnotationModel
	Requires    []RequiresDirective
	Exports     []ExportsDirective
	Opens       []OpensDirective
	Uses        []string
	Provides    []ProvidesDirective
}

type RequiresDirective struct {
	ModuleName   string
	IsTransitive bool
	IsStatic     bool
}

type ExportsDirective struct {
	PackageName string
	ToModules   []string
}

type OpensDirective struct {
	PackageName string
	ToModules   []string
}

type ProvidesDirective struct {
	ServiceName         string
	ImplementationNames []string
}

type PackageInfoModel struct {
	Name        string
	SourceFile  string
	SourceURL   URLString
	Javadoc     string
	Annotations []AnnotationModel
}
