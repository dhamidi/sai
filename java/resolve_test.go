package java

import (
	"testing"
)

func TestResolveInnerClassReferencesCrossFile(t *testing.T) {
	// Simulate the cross-file scenario from the thread:
	// - Authentication.java defines Authentication with inner class HeaderInfo
	// - Some other file references HeaderInfo and it gets incorrectly resolved

	// These would be scanned from separate files
	authenticationClass := &ClassModel{
		Name:       "org.eclipse.jetty.client.Authentication",
		SimpleName: "Authentication",
		Package:    "org.eclipse.jetty.client",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
		InnerClasses: []InnerClassModel{
			{
				InnerClass: "org.eclipse.jetty.client.Authentication.HeaderInfo",
				OuterClass: "org.eclipse.jetty.client.Authentication",
				InnerName:  "HeaderInfo",
				Visibility: VisibilityPublic,
				IsStatic:   true,
			},
		},
	}

	// Another class from another file that mistakenly resolved HeaderInfo to pkg.HeaderInfo
	consumerClass := &ClassModel{
		Name:       "org.eclipse.jetty.client.SomeConsumer",
		SimpleName: "SomeConsumer",
		Package:    "org.eclipse.jetty.client",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
		Methods: []MethodModel{
			{
				Name:       "process",
				Visibility: VisibilityPublic,
				ReturnType: TypeModel{Name: "void"},
				Parameters: []ParameterModel{
					{
						Name: "header",
						Type: TypeModel{
							Name: "org.eclipse.jetty.client.HeaderInfo", // Incorrectly resolved
						},
					},
				},
			},
		},
	}

	classes := []*ClassModel{authenticationClass, consumerClass}

	// Before fixup
	if consumerClass.Methods[0].Parameters[0].Type.Name != "org.eclipse.jetty.client.HeaderInfo" {
		t.Errorf("Expected incorrectly resolved type before fixup")
	}

	// Apply the fixup
	ResolveInnerClassReferences(classes)

	// After fixup, the parameter type should be corrected
	expectedParamType := "org.eclipse.jetty.client.Authentication.HeaderInfo"
	if consumerClass.Methods[0].Parameters[0].Type.Name != expectedParamType {
		t.Errorf("process parameter type = %q, want %q",
			consumerClass.Methods[0].Parameters[0].Type.Name,
			expectedParamType)
	}
}

func TestResolveInnerClassReferences(t *testing.T) {
	// Create mock class models that simulate the cross-file scenario
	headerInfoClass := &ClassModel{
		Name:       "org.eclipse.jetty.client.HeaderInfo",
		SimpleName: "HeaderInfo",
		Package:    "org.eclipse.jetty.client",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
	}

	authenticationClass := &ClassModel{
		Name:       "org.eclipse.jetty.client.Authentication",
		SimpleName: "Authentication",
		Package:    "org.eclipse.jetty.client",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
		InnerClasses: []InnerClassModel{
			{
				InnerClass: "org.eclipse.jetty.client.Authentication.HeaderInfo",
				OuterClass: "org.eclipse.jetty.client.Authentication",
				InnerName:  "HeaderInfo",
				Visibility: VisibilityPublic,
				IsStatic:   true,
			},
		},
		Methods: []MethodModel{
			{
				Name:       "createHeader",
				Visibility: VisibilityPublic,
				ReturnType: TypeModel{
					Name: "org.eclipse.jetty.client.HeaderInfo", // Incorrectly resolved
				},
				Parameters: []ParameterModel{
					{
						Name: "name",
						Type: TypeModel{Name: "java.lang.String"},
					},
				},
			},
		},
	}

	classes := []*ClassModel{headerInfoClass, authenticationClass}

	// Before fixup, the return type is incorrect
	if authenticationClass.Methods[0].ReturnType.Name != "org.eclipse.jetty.client.HeaderInfo" {
		t.Errorf("Expected incorrectly resolved type before fixup")
	}

	// Apply the fixup
	ResolveInnerClassReferences(classes)

	// After fixup, the return type should be the fully qualified inner class name
	expectedReturnType := "org.eclipse.jetty.client.Authentication.HeaderInfo"
	if authenticationClass.Methods[0].ReturnType.Name != expectedReturnType {
		t.Errorf("createHeader return type = %q, want %q",
			authenticationClass.Methods[0].ReturnType.Name,
			expectedReturnType)
	}
}

func TestResolveInnerClassReferencesMultipleInnerClasses(t *testing.T) {
	// Test with multiple inner classes in the same package
	outerClass := &ClassModel{
		Name:       "org.example.Outer",
		SimpleName: "Outer",
		Package:    "org.example",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
		InnerClasses: []InnerClassModel{
			{
				InnerClass: "org.example.Outer.Inner1",
				OuterClass: "org.example.Outer",
				InnerName:  "Inner1",
				Visibility: VisibilityPublic,
			},
			{
				InnerClass: "org.example.Outer.Inner2",
				OuterClass: "org.example.Outer",
				InnerName:  "Inner2",
				Visibility: VisibilityPublic,
			},
		},
		Methods: []MethodModel{
			{
				Name:       "getInner1",
				Visibility: VisibilityPublic,
				ReturnType: TypeModel{
					Name: "org.example.Inner1", // Incorrectly resolved
				},
			},
			{
				Name:       "getInner2",
				Visibility: VisibilityPublic,
				ReturnType: TypeModel{
					Name: "org.example.Inner2", // Incorrectly resolved
				},
			},
		},
	}

	classes := []*ClassModel{outerClass}

	// Apply the fixup
	ResolveInnerClassReferences(classes)

	// Check first method
	if outerClass.Methods[0].ReturnType.Name != "org.example.Outer.Inner1" {
		t.Errorf("getInner1 return type = %q, want %q",
			outerClass.Methods[0].ReturnType.Name,
			"org.example.Outer.Inner1")
	}

	// Check second method
	if outerClass.Methods[1].ReturnType.Name != "org.example.Outer.Inner2" {
		t.Errorf("getInner2 return type = %q, want %q",
			outerClass.Methods[1].ReturnType.Name,
			"org.example.Outer.Inner2")
	}
}

func TestResolveInnerClassReferencesInFields(t *testing.T) {
	outerClass := &ClassModel{
		Name:       "org.example.Container",
		SimpleName: "Container",
		Package:    "org.example",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
		InnerClasses: []InnerClassModel{
			{
				InnerClass: "org.example.Container.Item",
				OuterClass: "org.example.Container",
				InnerName:  "Item",
				Visibility: VisibilityPublic,
			},
		},
		Fields: []FieldModel{
			{
				Name:       "item",
				Visibility: VisibilityPrivate,
				Type: TypeModel{
					Name: "org.example.Item", // Incorrectly resolved
				},
			},
		},
	}

	classes := []*ClassModel{outerClass}

	// Apply the fixup
	ResolveInnerClassReferences(classes)

	// Check field type
	if outerClass.Fields[0].Type.Name != "org.example.Container.Item" {
		t.Errorf("item field type = %q, want %q",
			outerClass.Fields[0].Type.Name,
			"org.example.Container.Item")
	}
}

func TestResolveInnerClassReferencesInParameters(t *testing.T) {
	outerClass := &ClassModel{
		Name:       "org.example.Processor",
		SimpleName: "Processor",
		Package:    "org.example",
		Kind:       ClassKindClass,
		Visibility: VisibilityPublic,
		SuperClass: "java.lang.Object",
		InnerClasses: []InnerClassModel{
			{
				InnerClass: "org.example.Processor.Task",
				OuterClass: "org.example.Processor",
				InnerName:  "Task",
				Visibility: VisibilityPublic,
			},
		},
		Methods: []MethodModel{
			{
				Name:       "process",
				Visibility: VisibilityPublic,
				ReturnType: TypeModel{Name: "void"},
				Parameters: []ParameterModel{
					{
						Name: "task",
						Type: TypeModel{
							Name: "org.example.Task", // Incorrectly resolved
						},
					},
				},
			},
		},
	}

	classes := []*ClassModel{outerClass}

	// Apply the fixup
	ResolveInnerClassReferences(classes)

	// Check parameter type
	if outerClass.Methods[0].Parameters[0].Type.Name != "org.example.Processor.Task" {
		t.Errorf("process parameter type = %q, want %q",
			outerClass.Methods[0].Parameters[0].Type.Name,
			"org.example.Processor.Task")
	}
}
