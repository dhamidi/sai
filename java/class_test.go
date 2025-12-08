package java

import (
	"testing"
)

func TestParseClass(t *testing.T) {
	c, err := ParseClassFile("../classfile/testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to parse class file: %v", err)
	}

	t.Run("class name", func(t *testing.T) {
		if got := c.Name(); got != "testdata.TestClass" {
			t.Errorf("Name() = %q, want %q", got, "testdata.TestClass")
		}
	})

	t.Run("simple name", func(t *testing.T) {
		if got := c.SimpleName(); got != "TestClass" {
			t.Errorf("SimpleName() = %q, want %q", got, "TestClass")
		}
	})

	t.Run("package", func(t *testing.T) {
		if got := c.Package(); got != "testdata" {
			t.Errorf("Package() = %q, want %q", got, "testdata")
		}
	})

	t.Run("super class", func(t *testing.T) {
		if got := c.SuperClass(); got != "java.lang.Object" {
			t.Errorf("SuperClass() = %q, want %q", got, "java.lang.Object")
		}
	})

	t.Run("interfaces", func(t *testing.T) {
		interfaces := c.Interfaces()
		if len(interfaces) != 1 {
			t.Fatalf("Expected 1 interface, got %d", len(interfaces))
		}
		if interfaces[0] != "java.lang.Runnable" {
			t.Errorf("Interface[0] = %q, want %q", interfaces[0], "java.lang.Runnable")
		}
	})

	t.Run("is class", func(t *testing.T) {
		if !c.IsClass() {
			t.Error("Expected IsClass() to be true")
		}
		if c.IsInterface() {
			t.Error("Expected IsInterface() to be false")
		}
	})

	t.Run("visibility", func(t *testing.T) {
		if c.Visibility() != "public" {
			t.Errorf("Visibility() = %q, want %q", c.Visibility(), "public")
		}
	})
}

func TestClassFields(t *testing.T) {
	c, err := ParseClassFile("../classfile/testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to parse class file: %v", err)
	}

	t.Run("fields count", func(t *testing.T) {
		fields := c.Fields()
		if len(fields) != 3 {
			t.Fatalf("Expected 3 fields, got %d", len(fields))
		}
	})

	t.Run("CONSTANT_VALUE field", func(t *testing.T) {
		f := c.Field("CONSTANT_VALUE")
		if f == nil {
			t.Fatal("Expected to find CONSTANT_VALUE field")
		}
		if f.Name() != "CONSTANT_VALUE" {
			t.Errorf("Name() = %q, want %q", f.Name(), "CONSTANT_VALUE")
		}
		if f.Type().String() != "int" {
			t.Errorf("Type() = %q, want %q", f.Type().String(), "int")
		}
		if !f.IsPublic() || !f.IsStatic() || !f.IsFinal() {
			t.Error("CONSTANT_VALUE should be public static final")
		}
		if f.Visibility() != "public" {
			t.Errorf("Visibility() = %q, want %q", f.Visibility(), "public")
		}
	})

	t.Run("name field", func(t *testing.T) {
		f := c.Field("name")
		if f == nil {
			t.Fatal("Expected to find name field")
		}
		if f.Type().String() != "java.lang.String" {
			t.Errorf("Type() = %q, want %q", f.Type().String(), "java.lang.String")
		}
		if !f.IsPrivate() {
			t.Error("name field should be private")
		}
	})

	t.Run("count field", func(t *testing.T) {
		f := c.Field("count")
		if f == nil {
			t.Fatal("Expected to find count field")
		}
		if !f.IsProtected() {
			t.Error("count field should be protected")
		}
		if f.Visibility() != "protected" {
			t.Errorf("Visibility() = %q, want %q", f.Visibility(), "protected")
		}
	})
}

func TestClassMethods(t *testing.T) {
	c, err := ParseClassFile("../classfile/testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to parse class file: %v", err)
	}

	t.Run("constructors", func(t *testing.T) {
		constructors := c.Constructors()
		if len(constructors) != 2 {
			t.Fatalf("Expected 2 constructors, got %d", len(constructors))
		}
		for _, ctor := range constructors {
			if !ctor.IsConstructor() {
				t.Error("Expected constructor to report IsConstructor() = true")
			}
		}
	})

	t.Run("getName method", func(t *testing.T) {
		m := c.Method("getName")
		if m == nil {
			t.Fatal("Expected to find getName method")
		}
		if m.Name() != "getName" {
			t.Errorf("Name() = %q, want %q", m.Name(), "getName")
		}
		if m.ReturnType().String() != "java.lang.String" {
			t.Errorf("ReturnType() = %q, want %q", m.ReturnType().String(), "java.lang.String")
		}
		if m.ParameterCount() != 0 {
			t.Errorf("ParameterCount() = %d, want %d", m.ParameterCount(), 0)
		}
		if !m.IsPublic() {
			t.Error("getName should be public")
		}
	})

	t.Run("setName method", func(t *testing.T) {
		m := c.Method("setName")
		if m == nil {
			t.Fatal("Expected to find setName method")
		}
		if m.ReturnType().String() != "void" {
			t.Errorf("ReturnType() = %q, want %q", m.ReturnType().String(), "void")
		}
		if m.ParameterCount() != 1 {
			t.Fatalf("ParameterCount() = %d, want %d", m.ParameterCount(), 1)
		}
		params := m.Parameters()
		if params[0].Type.String() != "java.lang.String" {
			t.Errorf("Parameter[0].Type = %q, want %q", params[0].Type.String(), "java.lang.String")
		}
	})

	t.Run("helper method", func(t *testing.T) {
		m := c.Method("helper")
		if m == nil {
			t.Fatal("Expected to find helper method")
		}
		if !m.IsPrivate() || !m.IsStatic() {
			t.Error("helper should be private static")
		}
		if m.ReturnType().String() != "int" {
			t.Errorf("ReturnType() = %q, want %q", m.ReturnType().String(), "int")
		}
		if m.ParameterCount() != 2 {
			t.Fatalf("ParameterCount() = %d, want %d", m.ParameterCount(), 2)
		}
		params := m.Parameters()
		if params[0].Type.String() != "int" || params[1].Type.String() != "int" {
			t.Errorf("Parameters should both be int")
		}
	})

	t.Run("method visibility", func(t *testing.T) {
		getName := c.Method("getName")
		if getName.Visibility() != "public" {
			t.Errorf("getName.Visibility() = %q, want %q", getName.Visibility(), "public")
		}

		helper := c.Method("helper")
		if helper.Visibility() != "private" {
			t.Errorf("helper.Visibility() = %q, want %q", helper.Visibility(), "private")
		}
	})
}

func TestMethodString(t *testing.T) {
	c, err := ParseClassFile("../classfile/testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to parse class file: %v", err)
	}

	t.Run("getName", func(t *testing.T) {
		m := c.Method("getName")
		got := m.String()
		if got != "public java.lang.String getName()" {
			t.Errorf("String() = %q", got)
		}
	})

	t.Run("helper", func(t *testing.T) {
		m := c.Method("helper")
		got := m.String()
		if got != "private static int helper(int, int)" {
			t.Errorf("String() = %q", got)
		}
	})
}

func TestType(t *testing.T) {
	tests := []struct {
		typ       Type
		str       string
		primitive bool
		array     bool
		void      bool
	}{
		{Type{Name: "int"}, "int", true, false, false},
		{Type{Name: "boolean"}, "boolean", true, false, false},
		{Type{Name: "java.lang.String"}, "java.lang.String", false, false, false},
		{Type{Name: "int", ArrayDepth: 1}, "int[]", false, true, false},
		{Type{Name: "java.lang.Object", ArrayDepth: 2}, "java.lang.Object[][]", false, true, false},
		{Type{Name: "void"}, "void", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if got := tt.typ.String(); got != tt.str {
				t.Errorf("String() = %q, want %q", got, tt.str)
			}
			if got := tt.typ.IsPrimitive(); got != tt.primitive {
				t.Errorf("IsPrimitive() = %v, want %v", got, tt.primitive)
			}
			if got := tt.typ.IsArray(); got != tt.array {
				t.Errorf("IsArray() = %v, want %v", got, tt.array)
			}
			if got := tt.typ.IsVoid(); got != tt.void {
				t.Errorf("IsVoid() = %v, want %v", got, tt.void)
			}
		})
	}
}

func TestTypeElementType(t *testing.T) {
	arr := Type{Name: "int", ArrayDepth: 2}
	elem := arr.ElementType()
	if elem.ArrayDepth != 1 || elem.Name != "int" {
		t.Errorf("ElementType() = %v, want int[]", elem)
	}

	single := Type{Name: "int"}
	sameElem := single.ElementType()
	if sameElem.ArrayDepth != 0 || sameElem.Name != "int" {
		t.Errorf("ElementType() on non-array = %v, want int", sameElem)
	}
}
