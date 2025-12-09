package java

import (
	"testing"
)

func TestInnerClassTypeResolution(t *testing.T) {
	source := []byte(`
package org.eclipse.jetty.client;

public class Authentication {
    public static class HeaderInfo {
        private String name;
        private String value;
        
        public HeaderInfo(String name, String value) {
            this.name = name;
            this.value = value;
        }
    }
    
    public HeaderInfo createHeader(String n, String v) {
        return new HeaderInfo(n, v);
    }
}
`)

	models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	if len(models) != 2 {
		t.Fatalf("Expected 2 models (Authentication and HeaderInfo), got %d", len(models))
	}

	// Find the Authentication class
	var authClass *ClassModel
	var headerInfoClass *ClassModel
	for _, m := range models {
		if m.SimpleName == "Authentication" {
			authClass = m
		}
		if m.SimpleName == "HeaderInfo" {
			headerInfoClass = m
		}
	}

	if authClass == nil {
		t.Fatal("Expected to find Authentication class")
	}
	if headerInfoClass == nil {
		t.Fatal("Expected to find HeaderInfo class")
	}

	// Check HeaderInfo inner class is properly registered
	if headerInfoClass.Name != "org.eclipse.jetty.client.Authentication.HeaderInfo" {
		t.Errorf("HeaderInfo.Name = %q, want %q", headerInfoClass.Name, "org.eclipse.jetty.client.Authentication.HeaderInfo")
	}

	// Check that Authentication's createHeader method references HeaderInfo correctly
	if len(authClass.Methods) == 0 {
		t.Fatal("Expected Authentication to have methods")
	}

	// Find the createHeader method
	var createHeaderMethod *MethodModel
	for _, m := range authClass.Methods {
		if m.Name == "createHeader" {
			createHeaderMethod = &m
			break
		}
	}

	if createHeaderMethod == nil {
		t.Fatal("Expected to find createHeader method")
	}

	// The return type should be the fully qualified inner class name
	expectedReturnType := "org.eclipse.jetty.client.Authentication.HeaderInfo"
	if createHeaderMethod.ReturnType.Name != expectedReturnType {
		t.Errorf("createHeader return type = %q, want %q", createHeaderMethod.ReturnType.Name, expectedReturnType)
	}
}

func TestNestedInnerClassTypeResolution(t *testing.T) {
	source := []byte(`
package org.example;

public class Outer {
    public class Inner {
        public class DeepInner {
            public String value;
        }
    }
    
    public DeepInner createDeep() {
        return new DeepInner();
    }
}
`)

	models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	if len(models) < 2 {
		t.Fatalf("Expected at least 2 models, got %d", len(models))
	}

	// Find the Outer class
	var outerClass *ClassModel
	for _, m := range models {
		if m.SimpleName == "Outer" {
			outerClass = m
			break
		}
	}

	if outerClass == nil {
		t.Fatal("Expected to find Outer class")
	}

	// Find the createDeep method
	var createDeepMethod *MethodModel
	for _, m := range outerClass.Methods {
		if m.Name == "createDeep" {
			createDeepMethod = &m
			break
		}
	}

	if createDeepMethod == nil {
		t.Fatal("Expected to find createDeep method")
	}

	// The return type should resolve DeepInner to the fully qualified inner class name
	expectedReturnType := "org.example.Outer.Inner.DeepInner"
	if createDeepMethod.ReturnType.Name != expectedReturnType {
		t.Errorf("createDeep return type = %q, want %q", createDeepMethod.ReturnType.Name, expectedReturnType)
	}
}
