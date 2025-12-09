package java

import (
	"testing"
)

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
