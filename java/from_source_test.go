package java

import (
	"strings"
	"testing"
)

func TestSingleLineJavadoc(t *testing.T) {
	source := []byte(`package com.example;

/** A single-line javadoc */ public class Test {
    /** Field doc */ private String name;
    /** Method doc */ public void test() {}
}
`)

	models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	if len(models) != 1 {
		t.Fatalf("Expected 1 class model, got %d", len(models))
	}

	cls := models[0]

	t.Run("single-line class javadoc", func(t *testing.T) {
		if cls.Javadoc != "/** A single-line javadoc */" {
			t.Errorf("Expected class javadoc %q, got %q", "/** A single-line javadoc */", cls.Javadoc)
		}
	})

	t.Run("single-line field javadoc", func(t *testing.T) {
		if len(cls.Fields) != 1 {
			t.Fatalf("Expected 1 field, got %d", len(cls.Fields))
		}
		if cls.Fields[0].Javadoc != "/** Field doc */" {
			t.Errorf("Expected field javadoc %q, got %q", "/** Field doc */", cls.Fields[0].Javadoc)
		}
	})

	t.Run("single-line method javadoc", func(t *testing.T) {
		var method *MethodModel
		for i := range cls.Methods {
			if cls.Methods[i].Name == "test" {
				method = &cls.Methods[i]
				break
			}
		}
		if method == nil {
			t.Fatal("Expected to find test method")
		}
		if method.Javadoc != "/** Method doc */" {
			t.Errorf("Expected method javadoc %q, got %q", "/** Method doc */", method.Javadoc)
		}
	})
}

func TestJavadocExtraction(t *testing.T) {
	source := []byte(`package com.example;

/**
 * This is the class Javadoc.
 */
public class Example {
    /**
     * Field documentation.
     */
    private String name;

    /**
     * Constructor documentation.
     * @param name the name
     */
    public Example(String name) {
        this.name = name;
    }

    /**
     * Method documentation.
     * @return the name
     */
    public String getName() {
        return name;
    }

    // This is a line comment, not Javadoc
    public void setName(String name) {
        this.name = name;
    }

    /* This is a block comment, not Javadoc */
    public void noJavadoc() {
    }
}
`)

	models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	if len(models) != 1 {
		t.Fatalf("Expected 1 class model, got %d", len(models))
	}

	cls := models[0]

	t.Run("class javadoc", func(t *testing.T) {
		if cls.Javadoc == "" {
			t.Error("Expected class to have Javadoc")
		}
		if cls.Javadoc != "/**\n * This is the class Javadoc.\n */" {
			t.Errorf("Unexpected class Javadoc: %q", cls.Javadoc)
		}
	})

	t.Run("field javadoc", func(t *testing.T) {
		if len(cls.Fields) != 1 {
			t.Fatalf("Expected 1 field, got %d", len(cls.Fields))
		}
		field := cls.Fields[0]
		if field.Javadoc == "" {
			t.Error("Expected field to have Javadoc")
		}
		if field.Javadoc != "/**\n     * Field documentation.\n     */" {
			t.Errorf("Unexpected field Javadoc: %q", field.Javadoc)
		}
	})

	t.Run("constructor javadoc", func(t *testing.T) {
		var constructor *MethodModel
		for i := range cls.Methods {
			if cls.Methods[i].Name == "<init>" {
				constructor = &cls.Methods[i]
				break
			}
		}
		if constructor == nil {
			t.Fatal("Expected to find constructor")
		}
		if constructor.Javadoc == "" {
			t.Error("Expected constructor to have Javadoc")
		}
	})

	t.Run("method with javadoc", func(t *testing.T) {
		var method *MethodModel
		for i := range cls.Methods {
			if cls.Methods[i].Name == "getName" {
				method = &cls.Methods[i]
				break
			}
		}
		if method == nil {
			t.Fatal("Expected to find getName method")
		}
		if method.Javadoc == "" {
			t.Error("Expected getName method to have Javadoc")
		}
	})

	t.Run("method without javadoc", func(t *testing.T) {
		var method *MethodModel
		for i := range cls.Methods {
			if cls.Methods[i].Name == "setName" {
				method = &cls.Methods[i]
				break
			}
		}
		if method == nil {
			t.Fatal("Expected to find setName method")
		}
		if method.Javadoc != "" {
			t.Errorf("Expected setName method to have no Javadoc, got: %q", method.Javadoc)
		}
	})

	t.Run("method with block comment but no javadoc", func(t *testing.T) {
		var method *MethodModel
		for i := range cls.Methods {
			if cls.Methods[i].Name == "noJavadoc" {
				method = &cls.Methods[i]
				break
			}
		}
		if method == nil {
			t.Fatal("Expected to find noJavadoc method")
		}
		if method.Javadoc != "" {
			t.Errorf("Expected noJavadoc method to have no Javadoc (block comment without ** is not Javadoc), got: %q", method.Javadoc)
		}
	})
}

func TestPackageInfoModelFromSource(t *testing.T) {
	source := []byte(`/**
 * This is the package documentation.
 * It describes what the package does.
 * @since 1.0
 */
@Deprecated
@SuppressWarnings("unchecked")
package com.example.mypackage;
`)

	pkg, err := PackageInfoModelFromSource(source)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	if pkg == nil {
		t.Fatal("Expected package info model, got nil")
	}

	t.Run("package name", func(t *testing.T) {
		if pkg.Name != "com.example.mypackage" {
			t.Errorf("Expected package name %q, got %q", "com.example.mypackage", pkg.Name)
		}
	})

	t.Run("package javadoc", func(t *testing.T) {
		if pkg.Javadoc == "" {
			t.Error("Expected package javadoc, got empty string")
		}
		if !strings.Contains(pkg.Javadoc, "package documentation") {
			t.Errorf("Expected javadoc to contain 'package documentation', got %q", pkg.Javadoc)
		}
	})

	t.Run("package annotations", func(t *testing.T) {
		if len(pkg.Annotations) != 2 {
			t.Fatalf("Expected 2 annotations, got %d", len(pkg.Annotations))
		}

		var deprecated, suppressWarnings bool
		for _, ann := range pkg.Annotations {
			if ann.Type == "Deprecated" {
				deprecated = true
			}
			if ann.Type == "SuppressWarnings" {
				suppressWarnings = true
			}
		}
		if !deprecated {
			t.Error("Expected @Deprecated annotation")
		}
		if !suppressWarnings {
			t.Error("Expected @SuppressWarnings annotation")
		}
	})
}
