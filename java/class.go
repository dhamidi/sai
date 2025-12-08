package java

import (
	"io"

	"github.com/dhamidi/javalyzer/classfile"
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

func (c *Class) IsClass() bool      { return c.cf.IsClass() }
func (c *Class) IsInterface() bool  { return c.cf.IsInterface() }
func (c *Class) IsAnnotation() bool { return c.cf.IsAnnotation() }
func (c *Class) IsEnum() bool       { return c.cf.IsEnum() }
func (c *Class) IsModule() bool     { return c.cf.IsModule() }

func (c *Class) IsPublic() bool    { return c.cf.AccessFlags.IsPublic() }
func (c *Class) IsFinal() bool     { return c.cf.AccessFlags.IsFinal() }
func (c *Class) IsAbstract() bool  { return c.cf.AccessFlags.IsAbstract() }
func (c *Class) IsSynthetic() bool { return c.cf.AccessFlags.IsSynthetic() }

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

func (c *Class) MajorVersion() uint16 { return c.cf.MajorVersion }
func (c *Class) MinorVersion() uint16 { return c.cf.MinorVersion }

func (c *Class) ClassFile() *classfile.ClassFile {
	return c.cf
}
