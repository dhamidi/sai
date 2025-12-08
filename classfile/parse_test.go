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
		desc        string
		numParams   int
		returnsVoid bool
		returnType  string
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

func TestAnnotatedClassAttributes(t *testing.T) {
	cf, err := ParseFile("testdata/AnnotatedClass.class")
	if err != nil {
		t.Fatalf("Failed to parse AnnotatedClass.class: %v", err)
	}

	t.Run("RuntimeVisibleAnnotations", func(t *testing.T) {
		attr := cf.GetAttribute("RuntimeVisibleAnnotations")
		if attr == nil {
			t.Fatal("Expected RuntimeVisibleAnnotations attribute")
		}
		rva := attr.AsRuntimeVisibleAnnotations()
		if rva == nil {
			t.Fatal("Expected parsed RuntimeVisibleAnnotations")
		}
		if len(rva.Annotations) == 0 {
			t.Error("Expected at least one runtime visible annotation")
		}
	})

	t.Run("RuntimeInvisibleAnnotations", func(t *testing.T) {
		attr := cf.GetAttribute("RuntimeInvisibleAnnotations")
		if attr == nil {
			t.Fatal("Expected RuntimeInvisibleAnnotations attribute")
		}
		ria := attr.AsRuntimeInvisibleAnnotations()
		if ria == nil {
			t.Fatal("Expected parsed RuntimeInvisibleAnnotations")
		}
		if len(ria.Annotations) == 0 {
			t.Error("Expected at least one runtime invisible annotation")
		}
	})

	t.Run("Signature", func(t *testing.T) {
		attr := cf.GetAttribute("Signature")
		if attr == nil {
			t.Fatal("Expected Signature attribute (generics)")
		}
		sig := attr.AsSignature()
		if sig == nil {
			t.Fatal("Expected parsed Signature")
		}
		sigValue := cf.ConstantPool.GetUtf8(sig.SignatureIndex)
		if sigValue == "" {
			t.Error("Expected non-empty signature")
		}
	})

	t.Run("InnerClasses", func(t *testing.T) {
		attr := cf.GetAttribute("InnerClasses")
		if attr == nil {
			t.Fatal("Expected InnerClasses attribute")
		}
		ic := attr.AsInnerClasses()
		if ic == nil {
			t.Fatal("Expected parsed InnerClasses")
		}
		if len(ic.Classes) == 0 {
			t.Error("Expected at least one inner class entry")
		}
	})

	t.Run("Deprecated attribute", func(t *testing.T) {
		found := false
		for _, attr := range cf.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "Deprecated" {
				found = true
				dep := attr.AsDeprecated()
				if dep == nil {
					t.Error("Expected parsed Deprecated attribute")
				}
				break
			}
		}
		if !found {
			t.Log("Deprecated attribute not found on class level (may be on field/method)")
		}
	})

	t.Run("Method LineNumberTable", func(t *testing.T) {
		method := cf.GetMethod("getValue", "()Ljava/lang/Comparable;")
		if method == nil {
			t.Fatal("Expected to find getValue method")
		}
		codeAttr := method.GetCodeAttribute(cf.ConstantPool)
		if codeAttr == nil {
			t.Fatal("Expected Code attribute on getValue")
		}
		lntFound := false
		for _, attr := range codeAttr.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "LineNumberTable" {
				lntFound = true
				lnt := attr.AsLineNumberTable()
				if lnt == nil {
					t.Error("Expected parsed LineNumberTable")
				} else if len(lnt.LineNumberTable) == 0 {
					t.Error("Expected non-empty LineNumberTable")
				}
				break
			}
		}
		if !lntFound {
			t.Error("Expected LineNumberTable in Code attribute")
		}
	})

	t.Run("Method LocalVariableTable", func(t *testing.T) {
		method := cf.GetMethod("getValue", "()Ljava/lang/Comparable;")
		if method == nil {
			t.Fatal("Expected to find getValue method")
		}
		codeAttr := method.GetCodeAttribute(cf.ConstantPool)
		if codeAttr == nil {
			t.Fatal("Expected Code attribute on getValue")
		}
		lvtFound := false
		for _, attr := range codeAttr.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "LocalVariableTable" {
				lvtFound = true
				lvt := attr.AsLocalVariableTable()
				if lvt == nil {
					t.Error("Expected parsed LocalVariableTable")
				} else if len(lvt.LocalVariableTable) == 0 {
					t.Error("Expected non-empty LocalVariableTable")
				}
				break
			}
		}
		if !lvtFound {
			t.Log("LocalVariableTable not found (may require debug compilation)")
		}
	})

	t.Run("Method StackMapTable", func(t *testing.T) {
		method := cf.GetMethod("getValue", "()Ljava/lang/Comparable;")
		if method == nil {
			t.Fatal("Expected to find getValue method")
		}
		codeAttr := method.GetCodeAttribute(cf.ConstantPool)
		if codeAttr == nil {
			t.Fatal("Expected Code attribute on getValue")
		}
		smtFound := false
		for _, attr := range codeAttr.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "StackMapTable" {
				smtFound = true
				smt := attr.AsStackMapTable()
				if smt == nil {
					t.Error("Expected parsed StackMapTable")
				}
				break
			}
		}
		if !smtFound {
			t.Log("StackMapTable not found (may depend on method complexity)")
		}
	})

	t.Run("Method Exceptions", func(t *testing.T) {
		method := cf.GetMethod("methodWithException", "()V")
		if method == nil {
			t.Fatal("Expected to find methodWithException")
		}
		var exceptionsAttr *ExceptionsAttribute
		for _, attr := range method.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "Exceptions" {
				exceptionsAttr = attr.AsExceptions()
				break
			}
		}
		if exceptionsAttr == nil {
			t.Fatal("Expected Exceptions attribute on methodWithException")
		}
		if len(exceptionsAttr.ExceptionIndexTable) < 2 {
			t.Errorf("Expected at least 2 declared exceptions, got %d", len(exceptionsAttr.ExceptionIndexTable))
		}
	})

	t.Run("Constructor ParameterAnnotations", func(t *testing.T) {
		for _, method := range cf.Methods {
			name := cf.ConstantPool.GetUtf8(method.NameIndex)
			if name == "<init>" {
				for _, attr := range method.Attributes {
					attrName := cf.ConstantPool.GetUtf8(attr.NameIndex)
					if attrName == "RuntimeVisibleParameterAnnotations" {
						rvpa := attr.AsRuntimeVisibleParameterAnnotations()
						if rvpa == nil {
							t.Error("Expected parsed RuntimeVisibleParameterAnnotations")
						} else if len(rvpa.ParameterAnnotations) == 0 {
							t.Error("Expected parameter annotations")
						}
						return
					}
				}
			}
		}
		t.Log("No RuntimeVisibleParameterAnnotations found on constructors")
	})
}

func TestInnerClassNestAttributes(t *testing.T) {
	cf, err := ParseFile("testdata/AnnotatedClass$InnerClass.class")
	if err != nil {
		t.Fatalf("Failed to parse inner class: %v", err)
	}

	t.Run("NestHost", func(t *testing.T) {
		attr := cf.GetAttribute("NestHost")
		if attr == nil {
			t.Fatal("Expected NestHost attribute on inner class")
		}
		nh := attr.AsNestHost()
		if nh == nil {
			t.Fatal("Expected parsed NestHost")
		}
		hostName := cf.ConstantPool.GetClassName(nh.HostClassIndex)
		if hostName != "testdata/AnnotatedClass" {
			t.Errorf("NestHost = %q, want %q", hostName, "testdata/AnnotatedClass")
		}
	})

	t.Run("EnclosingMethod", func(t *testing.T) {
		anonCf, err := ParseFile("testdata/AnnotatedClass$1.class")
		if err != nil {
			t.Fatalf("Failed to parse anonymous class: %v", err)
		}
		attr := anonCf.GetAttribute("EnclosingMethod")
		if attr == nil {
			t.Fatal("Expected EnclosingMethod attribute on anonymous class")
		}
		em := attr.AsEnclosingMethod()
		if em == nil {
			t.Fatal("Expected parsed EnclosingMethod")
		}
		className := anonCf.ConstantPool.GetClassName(em.ClassIndex)
		if className != "testdata/AnnotatedClass" {
			t.Errorf("EnclosingMethod class = %q, want %q", className, "testdata/AnnotatedClass")
		}
	})
}

func TestNestMembersAttribute(t *testing.T) {
	cf, err := ParseFile("testdata/AnnotatedClass.class")
	if err != nil {
		t.Fatalf("Failed to parse AnnotatedClass.class: %v", err)
	}

	attr := cf.GetAttribute("NestMembers")
	if attr == nil {
		t.Fatal("Expected NestMembers attribute on outer class")
	}
	nm := attr.AsNestMembers()
	if nm == nil {
		t.Fatal("Expected parsed NestMembers")
	}
	if len(nm.Classes) == 0 {
		t.Error("Expected at least one nest member")
	}
}

func TestRecordAttribute(t *testing.T) {
	cf, err := ParseFile("testdata/RecordClass.class")
	if err != nil {
		t.Fatalf("Failed to parse RecordClass.class: %v", err)
	}

	t.Run("Record attribute", func(t *testing.T) {
		attr := cf.GetAttribute("Record")
		if attr == nil {
			t.Fatal("Expected Record attribute")
		}
		rec := attr.AsRecord()
		if rec == nil {
			t.Fatal("Expected parsed Record")
		}
		if len(rec.Components) != 2 {
			t.Errorf("Expected 2 record components, got %d", len(rec.Components))
		}
		names := make([]string, len(rec.Components))
		for i, comp := range rec.Components {
			names[i] = cf.ConstantPool.GetUtf8(comp.NameIndex)
		}
		expectedNames := []string{"name", "value"}
		for i, expected := range expectedNames {
			if i < len(names) && names[i] != expected {
				t.Errorf("Component %d: got %q, want %q", i, names[i], expected)
			}
		}
	})
}

func TestPermittedSubclassesAttribute(t *testing.T) {
	cf, err := ParseFile("testdata/SealedClass.class")
	if err != nil {
		t.Fatalf("Failed to parse SealedClass.class: %v", err)
	}

	attr := cf.GetAttribute("PermittedSubclasses")
	if attr == nil {
		t.Fatal("Expected PermittedSubclasses attribute on sealed class")
	}
	ps := attr.AsPermittedSubclasses()
	if ps == nil {
		t.Fatal("Expected parsed PermittedSubclasses")
	}
	if len(ps.Classes) != 2 {
		t.Errorf("Expected 2 permitted subclasses, got %d", len(ps.Classes))
	}
	subclassNames := make([]string, len(ps.Classes))
	for i, classIdx := range ps.Classes {
		subclassNames[i] = cf.ConstantPool.GetClassName(classIdx)
	}
	hasSubClass1 := false
	hasSubClass2 := false
	for _, name := range subclassNames {
		if name == "testdata/SubClass1" {
			hasSubClass1 = true
		}
		if name == "testdata/SubClass2" {
			hasSubClass2 = true
		}
	}
	if !hasSubClass1 || !hasSubClass2 {
		t.Errorf("Expected SubClass1 and SubClass2 in permitted subclasses, got %v", subclassNames)
	}
}

func TestConstantPoolAdvanced(t *testing.T) {
	cf, err := ParseFile("testdata/ConstantPoolTest.class")
	if err != nil {
		t.Fatalf("Failed to parse ConstantPoolTest.class: %v", err)
	}

	t.Run("Long constant", func(t *testing.T) {
		field := cf.GetField("LONG_CONST")
		if field == nil {
			t.Fatal("Expected LONG_CONST field")
		}
		for _, attr := range field.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "ConstantValue" {
				cv := attr.AsConstantValue()
				if cv == nil {
					t.Fatal("Expected parsed ConstantValue")
				}
				val, ok := cf.ConstantPool.GetLong(cv.ConstantValueIndex)
				if !ok {
					t.Error("Expected GetLong to succeed")
				} else if val != 9223372036854775807 {
					t.Errorf("Long value = %d, want 9223372036854775807", val)
				}
				return
			}
		}
		t.Error("ConstantValue attribute not found on LONG_CONST")
	})

	t.Run("Double constant", func(t *testing.T) {
		field := cf.GetField("DOUBLE_CONST")
		if field == nil {
			t.Fatal("Expected DOUBLE_CONST field")
		}
		for _, attr := range field.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "ConstantValue" {
				cv := attr.AsConstantValue()
				if cv == nil {
					t.Fatal("Expected parsed ConstantValue")
				}
				val, ok := cf.ConstantPool.GetDouble(cv.ConstantValueIndex)
				if !ok {
					t.Error("Expected GetDouble to succeed")
				} else if val < 1.0e308 {
					t.Errorf("Double value = %e, expected around 1.7976931348623157E308", val)
				}
				return
			}
		}
		t.Error("ConstantValue attribute not found on DOUBLE_CONST")
	})

	t.Run("Float constant", func(t *testing.T) {
		field := cf.GetField("FLOAT_CONST")
		if field == nil {
			t.Fatal("Expected FLOAT_CONST field")
		}
		for _, attr := range field.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "ConstantValue" {
				cv := attr.AsConstantValue()
				if cv == nil {
					t.Fatal("Expected parsed ConstantValue")
				}
				val, ok := cf.ConstantPool.GetFloat(cv.ConstantValueIndex)
				if !ok {
					t.Error("Expected GetFloat to succeed")
				} else if val < 3.0e38 {
					t.Errorf("Float value = %e, expected around 3.4028235E38", val)
				}
				return
			}
		}
		t.Error("ConstantValue attribute not found on FLOAT_CONST")
	})

	t.Run("Integer constant", func(t *testing.T) {
		field := cf.GetField("INT_CONST")
		if field == nil {
			t.Fatal("Expected INT_CONST field")
		}
		for _, attr := range field.Attributes {
			name := cf.ConstantPool.GetUtf8(attr.NameIndex)
			if name == "ConstantValue" {
				cv := attr.AsConstantValue()
				if cv == nil {
					t.Fatal("Expected parsed ConstantValue")
				}
				val, ok := cf.ConstantPool.GetInteger(cv.ConstantValueIndex)
				if !ok {
					t.Error("Expected GetInteger to succeed")
				} else if val != 2147483647 {
					t.Errorf("Integer value = %d, want 2147483647", val)
				}
				return
			}
		}
		t.Error("ConstantValue attribute not found on INT_CONST")
	})

	t.Run("BootstrapMethods attribute (lambda)", func(t *testing.T) {
		attr := cf.GetAttribute("BootstrapMethods")
		if attr == nil {
			t.Fatal("Expected BootstrapMethods attribute (from lambda)")
		}
		bm := attr.AsBootstrapMethods()
		if bm == nil {
			t.Fatal("Expected parsed BootstrapMethods")
		}
		if len(bm.BootstrapMethods) == 0 {
			t.Error("Expected at least one bootstrap method")
		}
	})

	t.Run("Constant pool entry types", func(t *testing.T) {
		tagCounts := make(map[ConstantTag]int)
		for _, entry := range cf.ConstantPool {
			if entry != nil {
				tagCounts[entry.Tag()]++
			}
		}

		requiredTags := []ConstantTag{
			ConstantUtf8,
			ConstantClass,
			ConstantMethodref,
			ConstantFieldref,
			ConstantNameAndType,
			ConstantString,
		}

		for _, tag := range requiredTags {
			if tagCounts[tag] == 0 {
				t.Errorf("Expected at least one constant pool entry with tag %d", tag)
			}
		}
	})
}

func TestConstantPoolAccessorBoundaryConditions(t *testing.T) {
	cf, err := ParseFile("testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to parse TestClass.class: %v", err)
	}

	t.Run("GetUtf8 with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetUtf8(0)
		if result != "" {
			t.Error("Expected empty string for index 0")
		}
		result = cf.ConstantPool.GetUtf8(65535)
		if result != "" {
			t.Error("Expected empty string for out-of-bounds index")
		}
	})

	t.Run("GetClassName with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetClassName(0)
		if result != "" {
			t.Error("Expected empty string for index 0")
		}
		result = cf.ConstantPool.GetClassName(65535)
		if result != "" {
			t.Error("Expected empty string for out-of-bounds index")
		}
	})

	t.Run("GetNameAndType with invalid index", func(t *testing.T) {
		name, desc := cf.ConstantPool.GetNameAndType(0)
		if name != "" || desc != "" {
			t.Error("Expected empty strings for index 0")
		}
	})

	t.Run("GetString with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetString(0)
		if result != "" {
			t.Error("Expected empty string for index 0")
		}
	})

	t.Run("GetInteger with invalid index", func(t *testing.T) {
		_, ok := cf.ConstantPool.GetInteger(0)
		if ok {
			t.Error("Expected false for index 0")
		}
	})

	t.Run("GetLong with invalid index", func(t *testing.T) {
		_, ok := cf.ConstantPool.GetLong(0)
		if ok {
			t.Error("Expected false for index 0")
		}
	})

	t.Run("GetFloat with invalid index", func(t *testing.T) {
		_, ok := cf.ConstantPool.GetFloat(0)
		if ok {
			t.Error("Expected false for index 0")
		}
	})

	t.Run("GetDouble with invalid index", func(t *testing.T) {
		_, ok := cf.ConstantPool.GetDouble(0)
		if ok {
			t.Error("Expected false for index 0")
		}
	})

	t.Run("GetFieldref with invalid index", func(t *testing.T) {
		cn, n, d := cf.ConstantPool.GetFieldref(0)
		if cn != "" || n != "" || d != "" {
			t.Error("Expected empty strings for index 0")
		}
	})

	t.Run("GetMethodref with invalid index", func(t *testing.T) {
		cn, n, d := cf.ConstantPool.GetMethodref(0)
		if cn != "" || n != "" || d != "" {
			t.Error("Expected empty strings for index 0")
		}
	})

	t.Run("GetInterfaceMethodref with invalid index", func(t *testing.T) {
		cn, n, d := cf.ConstantPool.GetInterfaceMethodref(0)
		if cn != "" || n != "" || d != "" {
			t.Error("Expected empty strings for index 0")
		}
	})

	t.Run("GetMethodHandle with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetMethodHandle(0)
		if result != nil {
			t.Error("Expected nil for index 0")
		}
	})

	t.Run("GetMethodType with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetMethodType(0)
		if result != "" {
			t.Error("Expected empty string for index 0")
		}
	})

	t.Run("GetDynamic with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetDynamic(0)
		if result != nil {
			t.Error("Expected nil for index 0")
		}
	})

	t.Run("GetInvokeDynamic with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetInvokeDynamic(0)
		if result != nil {
			t.Error("Expected nil for index 0")
		}
	})

	t.Run("GetModuleName with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetModuleName(0)
		if result != "" {
			t.Error("Expected empty string for index 0")
		}
	})

	t.Run("GetPackageName with invalid index", func(t *testing.T) {
		result := cf.ConstantPool.GetPackageName(0)
		if result != "" {
			t.Error("Expected empty string for index 0")
		}
	})
}

func TestAttributeAsMethodsReturnNil(t *testing.T) {
	cf, err := ParseFile("testdata/TestClass.class")
	if err != nil {
		t.Fatalf("Failed to parse TestClass.class: %v", err)
	}

	sourceFileAttr := cf.GetAttribute("SourceFile")
	if sourceFileAttr == nil {
		t.Fatal("Expected SourceFile attribute")
	}

	if sourceFileAttr.AsCode() != nil {
		t.Error("AsCode should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsLineNumberTable() != nil {
		t.Error("AsLineNumberTable should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsLocalVariableTable() != nil {
		t.Error("AsLocalVariableTable should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsConstantValue() != nil {
		t.Error("AsConstantValue should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsExceptions() != nil {
		t.Error("AsExceptions should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsInnerClasses() != nil {
		t.Error("AsInnerClasses should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsSignature() != nil {
		t.Error("AsSignature should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsBootstrapMethods() != nil {
		t.Error("AsBootstrapMethods should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsEnclosingMethod() != nil {
		t.Error("AsEnclosingMethod should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsSynthetic() != nil {
		t.Error("AsSynthetic should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsDeprecated() != nil {
		t.Error("AsDeprecated should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsSourceDebugExtension() != nil {
		t.Error("AsSourceDebugExtension should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsLocalVariableTypeTable() != nil {
		t.Error("AsLocalVariableTypeTable should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsMethodParameters() != nil {
		t.Error("AsMethodParameters should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsNestHost() != nil {
		t.Error("AsNestHost should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsNestMembers() != nil {
		t.Error("AsNestMembers should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRecord() != nil {
		t.Error("AsRecord should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsPermittedSubclasses() != nil {
		t.Error("AsPermittedSubclasses should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsStackMapTable() != nil {
		t.Error("AsStackMapTable should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRuntimeVisibleAnnotations() != nil {
		t.Error("AsRuntimeVisibleAnnotations should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRuntimeInvisibleAnnotations() != nil {
		t.Error("AsRuntimeInvisibleAnnotations should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRuntimeVisibleParameterAnnotations() != nil {
		t.Error("AsRuntimeVisibleParameterAnnotations should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRuntimeInvisibleParameterAnnotations() != nil {
		t.Error("AsRuntimeInvisibleParameterAnnotations should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRuntimeVisibleTypeAnnotations() != nil {
		t.Error("AsRuntimeVisibleTypeAnnotations should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsRuntimeInvisibleTypeAnnotations() != nil {
		t.Error("AsRuntimeInvisibleTypeAnnotations should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsAnnotationDefault() != nil {
		t.Error("AsAnnotationDefault should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsModule() != nil {
		t.Error("AsModule should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsModulePackages() != nil {
		t.Error("AsModulePackages should return nil for SourceFile attribute")
	}
	if sourceFileAttr.AsModuleMainClass() != nil {
		t.Error("AsModuleMainClass should return nil for SourceFile attribute")
	}
}

func TestConstantPoolTagMethods(t *testing.T) {
	tests := []struct {
		entry ConstantPoolEntry
		tag   ConstantTag
	}{
		{&ConstantUtf8Info{Value: "test"}, ConstantUtf8},
		{&ConstantIntegerInfo{Value: 42}, ConstantInteger},
		{&ConstantFloatInfo{Value: 3.14}, ConstantFloat},
		{&ConstantLongInfo{Value: 12345}, ConstantLong},
		{&ConstantDoubleInfo{Value: 2.718}, ConstantDouble},
		{&ConstantClassInfo{NameIndex: 1}, ConstantClass},
		{&ConstantStringInfo{StringIndex: 1}, ConstantString},
		{&ConstantFieldrefInfo{ClassIndex: 1, NameAndTypeIndex: 2}, ConstantFieldref},
		{&ConstantMethodrefInfo{ClassIndex: 1, NameAndTypeIndex: 2}, ConstantMethodref},
		{&ConstantInterfaceMethodrefInfo{ClassIndex: 1, NameAndTypeIndex: 2}, ConstantInterfaceMethodref},
		{&ConstantNameAndTypeInfo{NameIndex: 1, DescriptorIndex: 2}, ConstantNameAndType},
		{&ConstantMethodHandleInfo{ReferenceKind: RefInvokeVirtual, ReferenceIndex: 1}, ConstantMethodHandle},
		{&ConstantMethodTypeInfo{DescriptorIndex: 1}, ConstantMethodType},
		{&ConstantDynamicInfo{BootstrapMethodAttrIndex: 0, NameAndTypeIndex: 1}, ConstantDynamic},
		{&ConstantInvokeDynamicInfo{BootstrapMethodAttrIndex: 0, NameAndTypeIndex: 1}, ConstantInvokeDynamic},
		{&ConstantModuleInfo{NameIndex: 1}, ConstantModule},
		{&ConstantPackageInfo{NameIndex: 1}, ConstantPackage},
	}

	for _, tt := range tests {
		if got := tt.entry.Tag(); got != tt.tag {
			t.Errorf("Tag() = %d, want %d for %T", got, tt.tag, tt.entry)
		}
	}
}

func TestSyntheticAndBridgeMethods(t *testing.T) {
	cf, err := ParseFile("testdata/AnnotatedClass.class")
	if err != nil {
		t.Fatalf("Failed to parse AnnotatedClass.class: %v", err)
	}

	hasSyntheticOrBridge := false
	for _, method := range cf.Methods {
		if method.AccessFlags.IsSynthetic() || method.AccessFlags.IsBridge() {
			hasSyntheticOrBridge = true
			break
		}
	}
	t.Logf("Has synthetic or bridge method: %v", hasSyntheticOrBridge)
}
