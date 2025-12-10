package java

import (
	"io"

	"github.com/dhamidi/sai/classfile"
)

type Class struct {
	cf *classfile.ClassFile
}

func ParseClass(r io.Reader) (*Class, error) {
	cf, err := classfile.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Class{cf: cf}, nil
}

func ParseClassFile(path string) (*Class, error) {
	cf, err := classfile.ParseFile(path)
	if err != nil {
		return nil, err
	}
	return &Class{cf: cf}, nil
}

func (c *Class) Name() string {
	return classfile.InternalToSourceName(c.cf.ClassName())
}

func (c *Class) SimpleName() string {
	name := c.Name()
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return name
}

func (c *Class) Package() string {
	name := c.Name()
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[:i]
		}
	}
	return ""
}

func (c *Class) SuperClass() string {
	super := c.cf.SuperClassName()
	if super == "" {
		return ""
	}
	return classfile.InternalToSourceName(super)
}

func (c *Class) Interfaces() []string {
	internal := c.cf.InterfaceNames()
	result := make([]string, len(internal))
	for i, name := range internal {
		result[i] = classfile.InternalToSourceName(name)
	}
	return result
}

func (c *Class) IsClass() bool {
	return c.cf.IsClass()
}

func (c *Class) IsInterface() bool {
	return c.cf.IsInterface()
}

func (c *Class) IsAnnotation() bool {
	return c.cf.IsAnnotation()
}

func (c *Class) IsEnum() bool {
	return c.cf.IsEnum()
}

func (c *Class) IsModule() bool {
	return c.cf.IsModule()
}

func (c *Class) IsPublic() bool {
	return c.cf.AccessFlags.IsPublic()
}

func (c *Class) IsFinal() bool {
	return c.cf.AccessFlags.IsFinal()
}

func (c *Class) IsAbstract() bool {
	return c.cf.AccessFlags.IsAbstract()
}

func (c *Class) IsSynthetic() bool {
	return c.cf.AccessFlags.IsSynthetic()
}

func (c *Class) Methods() []Method {
	methods := make([]Method, len(c.cf.Methods))
	for i := range c.cf.Methods {
		methods[i] = Method{
			info: &c.cf.Methods[i],
			cp:   c.cf.ConstantPool,
		}
	}
	return methods
}

func (c *Class) Method(name string) *Method {
	for i := range c.cf.Methods {
		if c.cf.Methods[i].Name(c.cf.ConstantPool) == name {
			return &Method{
				info: &c.cf.Methods[i],
				cp:   c.cf.ConstantPool,
			}
		}
	}
	return nil
}

func (c *Class) MethodsByName(name string) []Method {
	var methods []Method
	for i := range c.cf.Methods {
		if c.cf.Methods[i].Name(c.cf.ConstantPool) == name {
			methods = append(methods, Method{
				info: &c.cf.Methods[i],
				cp:   c.cf.ConstantPool,
			})
		}
	}
	return methods
}

func (c *Class) Constructors() []Method {
	return c.MethodsByName("<init>")
}

func (c *Class) Fields() []Field {
	fields := make([]Field, len(c.cf.Fields))
	for i := range c.cf.Fields {
		fields[i] = Field{
			info: &c.cf.Fields[i],
			cp:   c.cf.ConstantPool,
		}
	}
	return fields
}

func (c *Class) Field(name string) *Field {
	for i := range c.cf.Fields {
		if c.cf.Fields[i].Name(c.cf.ConstantPool) == name {
			return &Field{
				info: &c.cf.Fields[i],
				cp:   c.cf.ConstantPool,
			}
		}
	}
	return nil
}

func (c *Class) Visibility() string {
	if c.IsPublic() {
		return "public"
	}
	return "package"
}

func (c *Class) MajorVersion() uint16 {
	return c.cf.MajorVersion
}

func (c *Class) MinorVersion() uint16 {
	return c.cf.MinorVersion
}

func (c *Class) ClassFile() *classfile.ClassFile {
	return c.cf
}

func (c *Class) Signature() string {
	attr := c.cf.GetAttribute("Signature")
	if attr == nil {
		return ""
	}
	if sig := attr.AsSignature(); sig != nil {
		return c.cf.ConstantPool.GetUtf8(sig.SignatureIndex)
	}
	return ""
}

func (c *Class) IsDeprecated() bool {
	return c.cf.GetAttribute("Deprecated") != nil
}

func (c *Class) SourceFile() string {
	attr := c.cf.GetAttribute("SourceFile")
	if attr == nil {
		return ""
	}
	if sf := attr.AsSourceFile(); sf != nil {
		return c.cf.ConstantPool.GetUtf8(sf.SourceFileIndex)
	}
	return ""
}

func (c *Class) Annotations() []Annotation {
	attr := c.cf.GetAttribute("RuntimeVisibleAnnotations")
	if attr == nil {
		return nil
	}
	if rva := attr.AsRuntimeVisibleAnnotations(); rva != nil {
		return annotationsFromClassfile(rva.Annotations, c.cf.ConstantPool)
	}
	return nil
}

func (c *Class) InvisibleAnnotations() []Annotation {
	attr := c.cf.GetAttribute("RuntimeInvisibleAnnotations")
	if attr == nil {
		return nil
	}
	if ria := attr.AsRuntimeInvisibleAnnotations(); ria != nil {
		return annotationsFromClassfile(ria.Annotations, c.cf.ConstantPool)
	}
	return nil
}

func (c *Class) IsRecord() bool {
	return c.cf.GetAttribute("Record") != nil
}

func (c *Class) RecordComponents() []RecordComponent {
	attr := c.cf.GetAttribute("Record")
	if attr == nil {
		return nil
	}
	if rec := attr.AsRecord(); rec != nil {
		result := make([]RecordComponent, len(rec.Components))
		for i, comp := range rec.Components {
			result[i] = RecordComponent{
				Name:       c.cf.ConstantPool.GetUtf8(comp.NameIndex),
				Descriptor: c.cf.ConstantPool.GetUtf8(comp.DescriptorIndex),
			}
		}
		return result
	}
	return nil
}

type RecordComponent struct {
	Name       string
	Descriptor string
}

func (rc RecordComponent) Type() Type {
	ft := classfile.ParseFieldDescriptor(rc.Descriptor)
	return typeFromFieldType(ft)
}

func (c *Class) IsSealed() bool {
	return c.cf.GetAttribute("PermittedSubclasses") != nil
}

func (c *Class) PermittedSubclasses() []string {
	attr := c.cf.GetAttribute("PermittedSubclasses")
	if attr == nil {
		return nil
	}
	if ps := attr.AsPermittedSubclasses(); ps != nil {
		result := make([]string, len(ps.Classes))
		for i, idx := range ps.Classes {
			result[i] = classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(idx))
		}
		return result
	}
	return nil
}

func (c *Class) NestHost() string {
	attr := c.cf.GetAttribute("NestHost")
	if attr == nil {
		return ""
	}
	if nh := attr.AsNestHost(); nh != nil {
		return classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(nh.HostClassIndex))
	}
	return ""
}

func (c *Class) NestMembers() []string {
	attr := c.cf.GetAttribute("NestMembers")
	if attr == nil {
		return nil
	}
	if nm := attr.AsNestMembers(); nm != nil {
		result := make([]string, len(nm.Classes))
		for i, idx := range nm.Classes {
			result[i] = classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(idx))
		}
		return result
	}
	return nil
}

func (c *Class) EnclosingClass() string {
	attr := c.cf.GetAttribute("EnclosingMethod")
	if attr == nil {
		return ""
	}
	if em := attr.AsEnclosingMethod(); em != nil {
		return classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(em.ClassIndex))
	}
	return ""
}

func (c *Class) EnclosingMethod() (className, methodName, methodDescriptor string) {
	attr := c.cf.GetAttribute("EnclosingMethod")
	if attr == nil {
		return "", "", ""
	}
	if em := attr.AsEnclosingMethod(); em != nil {
		className = classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(em.ClassIndex))
		if em.MethodIndex != 0 {
			methodName, methodDescriptor = c.cf.ConstantPool.GetNameAndType(em.MethodIndex)
		}
		return
	}
	return "", "", ""
}

func (c *Class) InnerClasses() []InnerClass {
	attr := c.cf.GetAttribute("InnerClasses")
	if attr == nil {
		return nil
	}
	if ic := attr.AsInnerClasses(); ic != nil {
		result := make([]InnerClass, len(ic.Classes))
		for i, entry := range ic.Classes {
			result[i] = InnerClass{
				InnerClass:  classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(entry.InnerClassInfoIndex)),
				OuterClass:  classfile.InternalToSourceName(c.cf.ConstantPool.GetClassName(entry.OuterClassInfoIndex)),
				InnerName:   c.cf.ConstantPool.GetUtf8(entry.InnerNameIndex),
				AccessFlags: entry.InnerClassAccessFlags,
			}
		}
		return result
	}
	return nil
}

type InnerClass struct {
	InnerClass  string
	OuterClass  string
	InnerName   string
	AccessFlags classfile.AccessFlags
}
