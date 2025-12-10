package format

import (
	"encoding/json"
	"io"

	"github.com/dhamidi/sai/java"
)

type JSONEncoder struct {
	w     io.Writer
	class *java.Class
}

func NewJSONEncoder(w io.Writer) *JSONEncoder {
	return &JSONEncoder{w: w}
}

func (e *JSONEncoder) Encode(class *java.Class) error {
	e.class = class
	text, err := e.MarshalText()
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *JSONEncoder) MarshalText() ([]byte, error) {
	data := e.buildClassData()
	return json.MarshalIndent(data, "", "  ")
}

type jsonClass struct {
	Name                string           `json:"name"`
	SimpleName          string           `json:"simpleName"`
	Package             string           `json:"package"`
	SuperClass          string           `json:"superClass,omitempty"`
	Interfaces          []string         `json:"interfaces,omitempty"`
	Visibility          string           `json:"visibility"`
	Kind                string           `json:"kind"`
	Modifiers           []string         `json:"modifiers,omitempty"`
	Version             jsonVersion      `json:"version"`
	Signature           string           `json:"signature,omitempty"`
	SourceFile          string           `json:"sourceFile,omitempty"`
	SourceURL           java.URLString   `json:"sourceURL,omitempty"`
	Deprecated          bool             `json:"deprecated,omitempty"`
	Annotations         []jsonAnnotation `json:"annotations,omitempty"`
	RecordComponents    []jsonRecordComp `json:"recordComponents,omitempty"`
	PermittedSubclasses []string         `json:"permittedSubclasses,omitempty"`
	NestHost            string           `json:"nestHost,omitempty"`
	NestMembers         []string         `json:"nestMembers,omitempty"`
	EnclosingClass      string           `json:"enclosingClass,omitempty"`
	InnerClasses        []jsonInnerClass `json:"innerClasses,omitempty"`
	Fields              []jsonField      `json:"fields,omitempty"`
	Methods             []jsonMethod     `json:"methods,omitempty"`
}

type jsonVersion struct {
	Major uint16 `json:"major"`
	Minor uint16 `json:"minor"`
}

type jsonField struct {
	Name          string           `json:"name"`
	Type          jsonType         `json:"type"`
	Visibility    string           `json:"visibility"`
	Modifiers     []string         `json:"modifiers,omitempty"`
	Signature     string           `json:"signature,omitempty"`
	Deprecated    bool             `json:"deprecated,omitempty"`
	Annotations   []jsonAnnotation `json:"annotations,omitempty"`
	ConstantValue interface{}      `json:"constantValue,omitempty"`
}

type jsonMethod struct {
	Name                 string             `json:"name"`
	ReturnType           jsonType           `json:"returnType"`
	Parameters           []jsonParameter    `json:"parameters,omitempty"`
	Visibility           string             `json:"visibility"`
	Modifiers            []string           `json:"modifiers,omitempty"`
	Signature            string             `json:"signature,omitempty"`
	Deprecated           bool               `json:"deprecated,omitempty"`
	Annotations          []jsonAnnotation   `json:"annotations,omitempty"`
	ParameterAnnotations [][]jsonAnnotation `json:"parameterAnnotations,omitempty"`
	Exceptions           []string           `json:"exceptions,omitempty"`
}

type jsonParameter struct {
	Name string   `json:"name,omitempty"`
	Type jsonType `json:"type"`
}

type jsonType struct {
	Name       string `json:"name"`
	ArrayDepth int    `json:"arrayDepth,omitempty"`
}

type jsonAnnotation struct {
	Type   string                 `json:"type"`
	Values map[string]interface{} `json:"values,omitempty"`
}

type jsonRecordComp struct {
	Name string   `json:"name"`
	Type jsonType `json:"type"`
}

type jsonInnerClass struct {
	InnerClass string   `json:"innerClass"`
	OuterClass string   `json:"outerClass,omitempty"`
	InnerName  string   `json:"innerName,omitempty"`
	Modifiers  []string `json:"modifiers,omitempty"`
}

func (e *JSONEncoder) buildClassData() jsonClass {
	c := e.class
	data := jsonClass{
		Name:                c.Name(),
		SimpleName:          c.SimpleName(),
		Package:             c.Package(),
		SuperClass:          c.SuperClass(),
		Interfaces:          c.Interfaces(),
		Visibility:          c.Visibility(),
		Kind:                e.classKind(),
		Modifiers:           e.classModifiers(),
		Version:             jsonVersion{Major: c.MajorVersion(), Minor: c.MinorVersion()},
		Signature:           c.Signature(),
		SourceFile:          c.SourceFile(),
		Deprecated:          c.IsDeprecated(),
		Annotations:         buildAnnotations(c.Annotations()),
		RecordComponents:    e.buildRecordComponents(),
		PermittedSubclasses: c.PermittedSubclasses(),
		NestHost:            c.NestHost(),
		NestMembers:         c.NestMembers(),
		EnclosingClass:      c.EnclosingClass(),
		InnerClasses:        e.buildInnerClasses(),
		Fields:              e.buildFields(),
		Methods:             e.buildMethods(),
	}
	return data
}

func (e *JSONEncoder) classKind() string {
	c := e.class
	switch {
	case c.IsAnnotation():
		return "annotation"
	case c.IsEnum():
		return "enum"
	case c.IsInterface():
		return "interface"
	case c.IsModule():
		return "module"
	default:
		return "class"
	}
}

func (e *JSONEncoder) classModifiers() []string {
	c := e.class
	var mods []string
	if c.IsFinal() {
		mods = append(mods, "final")
	}
	if c.IsAbstract() {
		mods = append(mods, "abstract")
	}
	if c.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	return mods
}

func (e *JSONEncoder) buildFields() []jsonField {
	fields := e.class.Fields()
	result := make([]jsonField, len(fields))
	for i, f := range fields {
		t := f.Type()
		result[i] = jsonField{
			Name:          f.Name(),
			Type:          jsonType{Name: t.Name, ArrayDepth: t.ArrayDepth},
			Visibility:    f.Visibility(),
			Modifiers:     fieldModifiers(f),
			Signature:     f.Signature(),
			Deprecated:    f.IsDeprecated(),
			Annotations:   buildAnnotations(f.Annotations()),
			ConstantValue: f.ConstantValue(),
		}
	}
	return result
}

func fieldModifiers(f java.Field) []string {
	var mods []string
	if f.IsStatic() {
		mods = append(mods, "static")
	}
	if f.IsFinal() {
		mods = append(mods, "final")
	}
	if f.IsVolatile() {
		mods = append(mods, "volatile")
	}
	if f.IsTransient() {
		mods = append(mods, "transient")
	}
	if f.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	if f.IsEnum() {
		mods = append(mods, "enum")
	}
	return mods
}

func (e *JSONEncoder) buildMethods() []jsonMethod {
	methods := e.class.Methods()
	result := make([]jsonMethod, len(methods))
	for i, m := range methods {
		rt := m.ReturnType()
		result[i] = jsonMethod{
			Name:                 m.Name(),
			ReturnType:           jsonType{Name: rt.Name, ArrayDepth: rt.ArrayDepth},
			Parameters:           buildParameters(m.Parameters()),
			Visibility:           m.Visibility(),
			Modifiers:            methodModifiers(m),
			Signature:            m.Signature(),
			Deprecated:           m.IsDeprecated(),
			Annotations:          buildAnnotations(m.Annotations()),
			ParameterAnnotations: buildParameterAnnotations(m.ParameterAnnotations()),
			Exceptions:           m.Exceptions(),
		}
	}
	return result
}

func buildParameters(params []java.Parameter) []jsonParameter {
	result := make([]jsonParameter, len(params))
	for i, p := range params {
		result[i] = jsonParameter{
			Name: p.Name,
			Type: jsonType{
				Name:       p.Type.Name,
				ArrayDepth: p.Type.ArrayDepth,
			},
		}
	}
	return result
}

func methodModifiers(m java.Method) []string {
	var mods []string
	if m.IsStatic() {
		mods = append(mods, "static")
	}
	if m.IsFinal() {
		mods = append(mods, "final")
	}
	if m.IsAbstract() {
		mods = append(mods, "abstract")
	}
	if m.IsSynchronized() {
		mods = append(mods, "synchronized")
	}
	if m.IsNative() {
		mods = append(mods, "native")
	}
	if m.IsBridge() {
		mods = append(mods, "bridge")
	}
	if m.IsVarargs() {
		mods = append(mods, "varargs")
	}
	if m.IsSynthetic() {
		mods = append(mods, "synthetic")
	}
	return mods
}

func buildAnnotations(anns []java.Annotation) []jsonAnnotation {
	if len(anns) == 0 {
		return nil
	}
	result := make([]jsonAnnotation, len(anns))
	for i, a := range anns {
		result[i] = jsonAnnotation{
			Type:   a.Type,
			Values: buildAnnotationValues(a.ElementValuePairs),
		}
	}
	return result
}

func buildAnnotationValues(pairs []java.ElementValuePair) map[string]interface{} {
	if len(pairs) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	for _, p := range pairs {
		result[p.Name] = p.Value
	}
	return result
}

func buildParameterAnnotations(paramAnns [][]java.Annotation) [][]jsonAnnotation {
	if len(paramAnns) == 0 {
		return nil
	}
	result := make([][]jsonAnnotation, len(paramAnns))
	for i, anns := range paramAnns {
		result[i] = buildAnnotations(anns)
	}
	return result
}

func (e *JSONEncoder) buildRecordComponents() []jsonRecordComp {
	comps := e.class.RecordComponents()
	if len(comps) == 0 {
		return nil
	}
	result := make([]jsonRecordComp, len(comps))
	for i, c := range comps {
		t := c.Type()
		result[i] = jsonRecordComp{
			Name: c.Name,
			Type: jsonType{Name: t.Name, ArrayDepth: t.ArrayDepth},
		}
	}
	return result
}

func (e *JSONEncoder) buildInnerClasses() []jsonInnerClass {
	classes := e.class.InnerClasses()
	if len(classes) == 0 {
		return nil
	}
	result := make([]jsonInnerClass, len(classes))
	for i, c := range classes {
		result[i] = jsonInnerClass{
			InnerClass: c.InnerClass,
			OuterClass: c.OuterClass,
			InnerName:  c.InnerName,
			Modifiers:  innerClassModifiers(c),
		}
	}
	return result
}

func innerClassModifiers(c java.InnerClass) []string {
	var mods []string
	if c.AccessFlags.IsPublic() {
		mods = append(mods, "public")
	}
	if c.AccessFlags.IsPrivate() {
		mods = append(mods, "private")
	}
	if c.AccessFlags.IsProtected() {
		mods = append(mods, "protected")
	}
	if c.AccessFlags.IsStatic() {
		mods = append(mods, "static")
	}
	if c.AccessFlags.IsFinal() {
		mods = append(mods, "final")
	}
	if c.AccessFlags.IsAbstract() {
		mods = append(mods, "abstract")
	}
	return mods
}

type JSONModelEncoder struct {
	w     io.Writer
	model *java.ClassModel
}

func NewJSONModelEncoder(w io.Writer) *JSONModelEncoder {
	return &JSONModelEncoder{w: w}
}

func (e *JSONModelEncoder) Encode(model *java.ClassModel) error {
	e.model = model
	text, err := e.MarshalText()
	if err != nil {
		return err
	}
	_, err = e.w.Write(text)
	return err
}

func (e *JSONModelEncoder) MarshalText() ([]byte, error) {
	data := e.buildClassData()
	return json.MarshalIndent(data, "", "  ")
}

func (e *JSONModelEncoder) buildClassData() jsonClass {
	m := e.model
	data := jsonClass{
		Name:                m.Name,
		SimpleName:          m.SimpleName,
		Package:             m.Package,
		SuperClass:          m.SuperClass,
		Interfaces:          m.Interfaces,
		Visibility:          string(m.Visibility),
		Kind:                string(m.Kind),
		Modifiers:           e.classModifiers(),
		Version:             jsonVersion{Major: m.MajorVersion, Minor: m.MinorVersion},
		Signature:           m.Signature,
		SourceFile:          m.SourceFile,
		SourceURL:           m.SourceURL,
		Deprecated:          m.IsDeprecated,
		Annotations:         buildModelAnnotations(m.Annotations),
		RecordComponents:    e.buildRecordComponents(),
		PermittedSubclasses: m.PermittedSubclasses,
		NestHost:            m.NestHost,
		NestMembers:         m.NestMembers,
		EnclosingClass:      m.EnclosingClass,
		InnerClasses:        e.buildInnerClasses(),
		Fields:              e.buildFields(),
		Methods:             e.buildMethods(),
	}
	return data
}

func (e *JSONModelEncoder) classModifiers() []string {
	m := e.model
	var mods []string
	if m.IsFinal {
		mods = append(mods, "final")
	}
	if m.IsAbstract {
		mods = append(mods, "abstract")
	}
	if m.IsSynthetic {
		mods = append(mods, "synthetic")
	}
	if m.IsSealed {
		mods = append(mods, "sealed")
	}
	return mods
}

func (e *JSONModelEncoder) buildFields() []jsonField {
	fields := e.model.Fields
	result := make([]jsonField, len(fields))
	for i, f := range fields {
		result[i] = jsonField{
			Name:          f.Name,
			Type:          jsonType{Name: f.Type.Name, ArrayDepth: f.Type.ArrayDepth},
			Visibility:    string(f.Visibility),
			Modifiers:     fieldModelModifiers(f),
			Signature:     f.Signature,
			Deprecated:    f.IsDeprecated,
			Annotations:   buildModelAnnotations(f.Annotations),
			ConstantValue: f.ConstantValue,
		}
	}
	return result
}

func fieldModelModifiers(f java.FieldModel) []string {
	var mods []string
	if f.IsStatic {
		mods = append(mods, "static")
	}
	if f.IsFinal {
		mods = append(mods, "final")
	}
	if f.IsVolatile {
		mods = append(mods, "volatile")
	}
	if f.IsTransient {
		mods = append(mods, "transient")
	}
	if f.IsSynthetic {
		mods = append(mods, "synthetic")
	}
	if f.IsEnum {
		mods = append(mods, "enum")
	}
	return mods
}

func (e *JSONModelEncoder) buildMethods() []jsonMethod {
	methods := e.model.Methods
	result := make([]jsonMethod, len(methods))
	for i, m := range methods {
		result[i] = jsonMethod{
			Name:                 m.Name,
			ReturnType:           jsonType{Name: m.ReturnType.Name, ArrayDepth: m.ReturnType.ArrayDepth},
			Parameters:           buildModelParameters(m.Parameters),
			Visibility:           string(m.Visibility),
			Modifiers:            methodModelModifiers(m),
			Signature:            m.Signature,
			Deprecated:           m.IsDeprecated,
			Annotations:          buildModelAnnotations(m.Annotations),
			ParameterAnnotations: buildModelParameterAnnotations(m.ParameterAnnotations),
			Exceptions:           m.Exceptions,
		}
	}
	return result
}

func buildModelParameters(params []java.ParameterModel) []jsonParameter {
	result := make([]jsonParameter, len(params))
	for i, p := range params {
		result[i] = jsonParameter{
			Name: p.Name,
			Type: jsonType{
				Name:       p.Type.Name,
				ArrayDepth: p.Type.ArrayDepth,
			},
		}
	}
	return result
}

func methodModelModifiers(m java.MethodModel) []string {
	var mods []string
	if m.IsStatic {
		mods = append(mods, "static")
	}
	if m.IsFinal {
		mods = append(mods, "final")
	}
	if m.IsAbstract {
		mods = append(mods, "abstract")
	}
	if m.IsSynchronized {
		mods = append(mods, "synchronized")
	}
	if m.IsNative {
		mods = append(mods, "native")
	}
	if m.IsBridge {
		mods = append(mods, "bridge")
	}
	if m.IsVarargs {
		mods = append(mods, "varargs")
	}
	if m.IsSynthetic {
		mods = append(mods, "synthetic")
	}
	if m.IsDefault {
		mods = append(mods, "default")
	}
	return mods
}

func buildModelAnnotations(anns []java.AnnotationModel) []jsonAnnotation {
	if len(anns) == 0 {
		return nil
	}
	result := make([]jsonAnnotation, len(anns))
	for i, a := range anns {
		result[i] = jsonAnnotation{
			Type:   a.Type,
			Values: a.Values,
		}
	}
	return result
}

func buildModelParameterAnnotations(paramAnns [][]java.AnnotationModel) [][]jsonAnnotation {
	if len(paramAnns) == 0 {
		return nil
	}
	result := make([][]jsonAnnotation, len(paramAnns))
	for i, anns := range paramAnns {
		result[i] = buildModelAnnotations(anns)
	}
	return result
}

func (e *JSONModelEncoder) buildRecordComponents() []jsonRecordComp {
	comps := e.model.RecordComponents
	if len(comps) == 0 {
		return nil
	}
	result := make([]jsonRecordComp, len(comps))
	for i, c := range comps {
		result[i] = jsonRecordComp{
			Name: c.Name,
			Type: jsonType{Name: c.Type.Name, ArrayDepth: c.Type.ArrayDepth},
		}
	}
	return result
}

func (e *JSONModelEncoder) buildInnerClasses() []jsonInnerClass {
	classes := e.model.InnerClasses
	if len(classes) == 0 {
		return nil
	}
	result := make([]jsonInnerClass, len(classes))
	for i, c := range classes {
		result[i] = jsonInnerClass{
			InnerClass: c.InnerClass,
			OuterClass: c.OuterClass,
			InnerName:  c.InnerName,
			Modifiers:  innerClassModelModifiers(c),
		}
	}
	return result
}

func innerClassModelModifiers(c java.InnerClassModel) []string {
	var mods []string
	if c.Visibility == java.VisibilityPublic {
		mods = append(mods, "public")
	}
	if c.Visibility == java.VisibilityPrivate {
		mods = append(mods, "private")
	}
	if c.Visibility == java.VisibilityProtected {
		mods = append(mods, "protected")
	}
	if c.IsStatic {
		mods = append(mods, "static")
	}
	if c.IsFinal {
		mods = append(mods, "final")
	}
	if c.IsAbstract {
		mods = append(mods, "abstract")
	}
	return mods
}
