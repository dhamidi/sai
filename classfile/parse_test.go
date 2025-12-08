package classfile

import (
	"os"
	"testing"
)

func TestParseClassFile(t *testing.T) {
	f, err := os.Open("testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to open test class file: %v", err)
	}
	defer f.Close()

	cf, err := Parse(f)
	if err != nil {
		t.Fatalf("Failed to parse class file: %v", err)
	}

	t.Run("class name", func(t *testing.T) {
		expected := "testdata/TestClass"
		if got := cf.ClassName(); got != expected {
			t.Errorf("ClassName() = %q, want %q", got, expected)
		}
	})

	t.Run("super class", func(t *testing.T) {
		expected := "java/lang/Object"
		if got := cf.SuperClassName(); got != expected {
			t.Errorf("SuperClassName() = %q, want %q", got, expected)
		}
	})

	t.Run("interfaces", func(t *testing.T) {
		interfaces := cf.InterfaceNames()
		if len(interfaces) != 1 {
			t.Fatalf("Expected 1 interface, got %d", len(interfaces))
		}
		expected := "java/lang/Runnable"
		if interfaces[0] != expected {
			t.Errorf("Interface[0] = %q, want %q", interfaces[0], expected)
		}
	})

	t.Run("is class", func(t *testing.T) {
		if !cf.IsClass() {
			t.Error("Expected IsClass() to be true")
		}
		if cf.IsInterface() {
			t.Error("Expected IsInterface() to be false")
		}
	})

	t.Run("access flags", func(t *testing.T) {
		if !cf.AccessFlags.IsPublic() {
			t.Error("Expected class to be public")
		}
		if cf.AccessFlags.IsFinal() {
			t.Error("Expected class to not be final")
		}
	})

	t.Run("fields", func(t *testing.T) {
		if len(cf.Fields) != 3 {
			t.Fatalf("Expected 3 fields, got %d", len(cf.Fields))
		}

		constantValue := cf.GetField("CONSTANT_VALUE")
		if constantValue == nil {
			t.Fatal("Expected to find CONSTANT_VALUE field")
		}
		if !constantValue.IsPublic() || !constantValue.IsStatic() || !constantValue.IsFinal() {
			t.Error("CONSTANT_VALUE should be public static final")
		}
		if constantValue.Descriptor(cf.ConstantPool) != "I" {
			t.Errorf("CONSTANT_VALUE descriptor = %q, want %q", constantValue.Descriptor(cf.ConstantPool), "I")
		}

		nameField := cf.GetField("name")
		if nameField == nil {
			t.Fatal("Expected to find name field")
		}
		if !nameField.IsPrivate() {
			t.Error("name field should be private")
		}
		if nameField.Descriptor(cf.ConstantPool) != "Ljava/lang/String;" {
			t.Errorf("name descriptor = %q, want %q", nameField.Descriptor(cf.ConstantPool), "Ljava/lang/String;")
		}

		countField := cf.GetField("count")
		if countField == nil {
			t.Fatal("Expected to find count field")
		}
		if !countField.IsProtected() {
			t.Error("count field should be protected")
		}
	})

	t.Run("methods", func(t *testing.T) {
		constructors := cf.GetMethods("<init>")
		if len(constructors) != 2 {
			t.Fatalf("Expected 2 constructors, got %d", len(constructors))
		}

		getNameMethod := cf.GetMethod("getName", "()Ljava/lang/String;")
		if getNameMethod == nil {
			t.Fatal("Expected to find getName method")
		}
		if !getNameMethod.IsPublic() {
			t.Error("getName should be public")
		}

		setNameMethod := cf.GetMethod("setName", "(Ljava/lang/String;)V")
		if setNameMethod == nil {
			t.Fatal("Expected to find setName method")
		}

		helperMethod := cf.GetMethod("helper", "(II)I")
		if helperMethod == nil {
			t.Fatal("Expected to find helper method")
		}
		if !helperMethod.IsPrivate() || !helperMethod.IsStatic() {
			t.Error("helper should be private static")
		}

		runMethod := cf.GetMethod("run", "()V")
		if runMethod == nil {
			t.Fatal("Expected to find run method")
		}
	})

	t.Run("method code attribute", func(t *testing.T) {
		getNameMethod := cf.GetMethod("getName", "()Ljava/lang/String;")
		if getNameMethod == nil {
			t.Fatal("Expected to find getName method")
		}

		codeAttr := getNameMethod.GetCodeAttribute(cf.ConstantPool)
		if codeAttr == nil {
			t.Fatal("Expected getName to have Code attribute")
		}

		if codeAttr.MaxStack == 0 {
			t.Error("MaxStack should be > 0")
		}
		if codeAttr.MaxLocals == 0 {
			t.Error("MaxLocals should be > 0")
		}
		if len(codeAttr.Code) == 0 {
			t.Error("Code should not be empty")
		}
	})

	t.Run("source file attribute", func(t *testing.T) {
		sourceFileAttr := cf.GetAttribute("SourceFile")
		if sourceFileAttr == nil {
			t.Fatal("Expected SourceFile attribute")
		}
		sf := sourceFileAttr.AsSourceFile()
		if sf == nil {
			t.Fatal("Expected parsed SourceFile")
		}
		sourceName := cf.ConstantPool.GetUtf8(sf.SourceFileIndex)
		if sourceName != "TestClass.java" {
			t.Errorf("SourceFile = %q, want %q", sourceName, "TestClass.java")
		}
	})
}

func TestParseFieldDescriptor(t *testing.T) {
	tests := []struct {
		desc       string
		baseType   string
		className  string
		arrayDepth int
	}{
		{"I", "int", "", 0},
		{"Z", "boolean", "", 0},
		{"Ljava/lang/String;", "", "java/lang/String", 0},
		{"[I", "int", "", 1},
		{"[[D", "double", "", 2},
		{"[Ljava/lang/Object;", "", "java/lang/Object", 1},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ft := ParseFieldDescriptor(tt.desc)
			if ft == nil {
				t.Fatalf("ParseFieldDescriptor(%q) returned nil", tt.desc)
			}
			if ft.BaseType != tt.baseType {
				t.Errorf("BaseType = %q, want %q", ft.BaseType, tt.baseType)
			}
			if ft.ClassName != tt.className {
				t.Errorf("ClassName = %q, want %q", ft.ClassName, tt.className)
			}
			if ft.ArrayDepth != tt.arrayDepth {
				t.Errorf("ArrayDepth = %d, want %d", ft.ArrayDepth, tt.arrayDepth)
			}
		})
	}
}

func TestParseMethodDescriptor(t *testing.T) {
	tests := []struct {
		desc         string
		numParams    int
		returnsVoid  bool
		returnType   string
	}{
		{"()V", 0, true, ""},
		{"()I", 0, false, "int"},
		{"(I)V", 1, true, ""},
		{"(II)I", 2, false, "int"},
		{"(Ljava/lang/String;)V", 1, true, ""},
		{"(IDLjava/lang/Thread;)Ljava/lang/Object;", 3, false, "java/lang/Object"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			md := ParseMethodDescriptor(tt.desc)
			if md == nil {
				t.Fatalf("ParseMethodDescriptor(%q) returned nil", tt.desc)
			}
			if len(md.Parameters) != tt.numParams {
				t.Errorf("len(Parameters) = %d, want %d", len(md.Parameters), tt.numParams)
			}
			if tt.returnsVoid {
				if md.ReturnType != nil {
					t.Error("Expected nil ReturnType for void")
				}
			} else {
				if md.ReturnType == nil {
					t.Error("Expected non-nil ReturnType")
				} else {
					if md.ReturnType.BaseType != "" && md.ReturnType.BaseType != tt.returnType {
						t.Errorf("ReturnType.BaseType = %q, want %q", md.ReturnType.BaseType, tt.returnType)
					}
					if md.ReturnType.ClassName != "" && md.ReturnType.ClassName != tt.returnType {
						t.Errorf("ReturnType.ClassName = %q, want %q", md.ReturnType.ClassName, tt.returnType)
					}
				}
			}
		})
	}
}

func TestConstantPoolGetters(t *testing.T) {
	f, err := os.Open("testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to open test class file: %v", err)
	}
	defer f.Close()

	cf, err := Parse(f)
	if err != nil {
		t.Fatalf("Failed to parse class file: %v", err)
	}

	className := cf.ConstantPool.GetClassName(cf.ThisClass)
	if className != "testdata/TestClass" {
		t.Errorf("GetClassName() = %q, want %q", className, "testdata/TestClass")
	}

	superClassName := cf.ConstantPool.GetClassName(cf.SuperClass)
	if superClassName != "java/lang/Object" {
		t.Errorf("GetClassName(super) = %q, want %q", superClassName, "java/lang/Object")
	}
}
