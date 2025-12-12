package v2

import (
	"reflect"
	"testing"

	"github.com/dhamidi/sai/java"
)

func TestV1V2Equivalence_SimpleClass(t *testing.T) {
	source := []byte(`package com.example;

public class Test {
    private String name;
    
    public Test(String name) {
        this.name = name;
    }
    
    public String getName() {
        return name;
    }
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func TestV1V2Equivalence_Interface(t *testing.T) {
	source := []byte(`package com.example;

public interface Runnable {
    void run();
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func TestV1V2Equivalence_Enum(t *testing.T) {
	source := []byte(`package com.example;

public enum Status {
    ACTIVE, INACTIVE, PENDING;
    
    public boolean isActive() {
        return this == ACTIVE;
    }
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func TestV1V2Equivalence_Record(t *testing.T) {
	source := []byte(`package com.example;

public record Point(int x, int y) {
    public double distance() {
        return Math.sqrt(x * x + y * y);
    }
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func TestV1V2Equivalence_Generics(t *testing.T) {
	source := []byte(`package com.example;

import java.util.List;

public class Container<T extends Comparable<T>> {
    private List<T> items;
    
    public void add(T item) {
        items.add(item);
    }
    
    public <U> U transform(java.util.function.Function<T, U> fn) {
        return null;
    }
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func TestV1V2Equivalence_InnerClass(t *testing.T) {
	source := []byte(`package com.example;

public class Outer {
    private String value;
    
    public class Inner {
        public String getValue() {
            return value;
        }
    }
    
    public static class StaticInner {
        private int count;
    }
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func TestV1V2Equivalence_Annotations(t *testing.T) {
	source := []byte(`package com.example;

@Deprecated
public class Legacy {
    @SuppressWarnings("unchecked")
    public void process() {
    }
}
`)

	v1Models, err := java.ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v1 parse failed: %v", err)
	}

	v2Models, err := ClassModelsFromSource(source)
	if err != nil {
		t.Fatalf("v2 parse failed: %v", err)
	}

	if len(v1Models) != len(v2Models) {
		t.Fatalf("model count differs: v1=%d, v2=%d", len(v1Models), len(v2Models))
	}

	for i := range v1Models {
		compareClassModels(t, v1Models[i], v2Models[i])
	}
}

func compareClassModels(t *testing.T, v1, v2 *java.ClassModel) {
	t.Helper()

	if v1.Name != v2.Name {
		t.Errorf("Name differs: v1=%q, v2=%q", v1.Name, v2.Name)
	}
	if v1.SimpleName != v2.SimpleName {
		t.Errorf("SimpleName differs: v1=%q, v2=%q", v1.SimpleName, v2.SimpleName)
	}
	if v1.Package != v2.Package {
		t.Errorf("Package differs: v1=%q, v2=%q", v1.Package, v2.Package)
	}
	if v1.Kind != v2.Kind {
		t.Errorf("Kind differs: v1=%q, v2=%q", v1.Kind, v2.Kind)
	}
	if v1.Visibility != v2.Visibility {
		t.Errorf("Visibility differs: v1=%q, v2=%q", v1.Visibility, v2.Visibility)
	}
	if v1.SuperClass != v2.SuperClass {
		t.Errorf("SuperClass differs: v1=%q, v2=%q", v1.SuperClass, v2.SuperClass)
	}
	if !reflect.DeepEqual(v1.Interfaces, v2.Interfaces) {
		t.Errorf("Interfaces differ: v1=%v, v2=%v", v1.Interfaces, v2.Interfaces)
	}
	if v1.IsFinal != v2.IsFinal {
		t.Errorf("IsFinal differs: v1=%v, v2=%v", v1.IsFinal, v2.IsFinal)
	}
	if v1.IsAbstract != v2.IsAbstract {
		t.Errorf("IsAbstract differs: v1=%v, v2=%v", v1.IsAbstract, v2.IsAbstract)
	}
	if v1.IsStatic != v2.IsStatic {
		t.Errorf("IsStatic differs: v1=%v, v2=%v", v1.IsStatic, v2.IsStatic)
	}

	if len(v1.TypeParameters) != len(v2.TypeParameters) {
		t.Errorf("TypeParameters count differs: v1=%d, v2=%d", len(v1.TypeParameters), len(v2.TypeParameters))
	} else {
		for i := range v1.TypeParameters {
			if v1.TypeParameters[i].Name != v2.TypeParameters[i].Name {
				t.Errorf("TypeParameter[%d].Name differs: v1=%q, v2=%q", i, v1.TypeParameters[i].Name, v2.TypeParameters[i].Name)
			}
		}
	}

	if len(v1.Fields) != len(v2.Fields) {
		t.Errorf("Fields count differs: v1=%d, v2=%d", len(v1.Fields), len(v2.Fields))
	} else {
		for i := range v1.Fields {
			compareFieldModels(t, i, v1.Fields[i], v2.Fields[i])
		}
	}

	if len(v1.Methods) != len(v2.Methods) {
		t.Errorf("Methods count differs: v1=%d, v2=%d", len(v1.Methods), len(v2.Methods))
	} else {
		for i := range v1.Methods {
			compareMethodModels(t, i, v1.Methods[i], v2.Methods[i])
		}
	}

	if len(v1.EnumConstants) != len(v2.EnumConstants) {
		t.Errorf("EnumConstants count differs: v1=%d, v2=%d", len(v1.EnumConstants), len(v2.EnumConstants))
	} else {
		for i := range v1.EnumConstants {
			if v1.EnumConstants[i].Name != v2.EnumConstants[i].Name {
				t.Errorf("EnumConstant[%d].Name differs: v1=%q, v2=%q", i, v1.EnumConstants[i].Name, v2.EnumConstants[i].Name)
			}
		}
	}

	if len(v1.RecordComponents) != len(v2.RecordComponents) {
		t.Errorf("RecordComponents count differs: v1=%d, v2=%d", len(v1.RecordComponents), len(v2.RecordComponents))
	} else {
		for i := range v1.RecordComponents {
			if v1.RecordComponents[i].Name != v2.RecordComponents[i].Name {
				t.Errorf("RecordComponent[%d].Name differs: v1=%q, v2=%q", i, v1.RecordComponents[i].Name, v2.RecordComponents[i].Name)
			}
			if v1.RecordComponents[i].Type.Name != v2.RecordComponents[i].Type.Name {
				t.Errorf("RecordComponent[%d].Type differs: v1=%q, v2=%q", i, v1.RecordComponents[i].Type.Name, v2.RecordComponents[i].Type.Name)
			}
		}
	}

	if len(v1.InnerClasses) != len(v2.InnerClasses) {
		t.Errorf("InnerClasses count differs: v1=%d, v2=%d", len(v1.InnerClasses), len(v2.InnerClasses))
	}
}

func compareFieldModels(t *testing.T, idx int, v1, v2 java.FieldModel) {
	t.Helper()

	if v1.Name != v2.Name {
		t.Errorf("Field[%d].Name differs: v1=%q, v2=%q", idx, v1.Name, v2.Name)
	}
	if v1.Type.Name != v2.Type.Name {
		t.Errorf("Field[%d].Type.Name differs: v1=%q, v2=%q", idx, v1.Type.Name, v2.Type.Name)
	}
	if v1.Type.ArrayDepth != v2.Type.ArrayDepth {
		t.Errorf("Field[%d].Type.ArrayDepth differs: v1=%d, v2=%d", idx, v1.Type.ArrayDepth, v2.Type.ArrayDepth)
	}
	if v1.Visibility != v2.Visibility {
		t.Errorf("Field[%d].Visibility differs: v1=%q, v2=%q", idx, v1.Visibility, v2.Visibility)
	}
	if v1.IsStatic != v2.IsStatic {
		t.Errorf("Field[%d].IsStatic differs: v1=%v, v2=%v", idx, v1.IsStatic, v2.IsStatic)
	}
	if v1.IsFinal != v2.IsFinal {
		t.Errorf("Field[%d].IsFinal differs: v1=%v, v2=%v", idx, v1.IsFinal, v2.IsFinal)
	}
}

func compareMethodModels(t *testing.T, idx int, v1, v2 java.MethodModel) {
	t.Helper()

	if v1.Name != v2.Name {
		t.Errorf("Method[%d].Name differs: v1=%q, v2=%q", idx, v1.Name, v2.Name)
	}
	if v1.ReturnType.Name != v2.ReturnType.Name {
		t.Errorf("Method[%d].ReturnType.Name differs: v1=%q, v2=%q", idx, v1.ReturnType.Name, v2.ReturnType.Name)
	}
	if v1.Visibility != v2.Visibility {
		t.Errorf("Method[%d].Visibility differs: v1=%q, v2=%q", idx, v1.Visibility, v2.Visibility)
	}
	if v1.IsStatic != v2.IsStatic {
		t.Errorf("Method[%d].IsStatic differs: v1=%v, v2=%v", idx, v1.IsStatic, v2.IsStatic)
	}
	if v1.IsFinal != v2.IsFinal {
		t.Errorf("Method[%d].IsFinal differs: v1=%v, v2=%v", idx, v1.IsFinal, v2.IsFinal)
	}
	if v1.IsAbstract != v2.IsAbstract {
		t.Errorf("Method[%d].IsAbstract differs: v1=%v, v2=%v", idx, v1.IsAbstract, v2.IsAbstract)
	}

	if len(v1.Parameters) != len(v2.Parameters) {
		t.Errorf("Method[%d].Parameters count differs: v1=%d, v2=%d", idx, len(v1.Parameters), len(v2.Parameters))
	} else {
		for i := range v1.Parameters {
			if v1.Parameters[i].Name != v2.Parameters[i].Name {
				t.Errorf("Method[%d].Parameter[%d].Name differs: v1=%q, v2=%q", idx, i, v1.Parameters[i].Name, v2.Parameters[i].Name)
			}
			if v1.Parameters[i].Type.Name != v2.Parameters[i].Type.Name {
				t.Errorf("Method[%d].Parameter[%d].Type.Name differs: v1=%q, v2=%q", idx, i, v1.Parameters[i].Type.Name, v2.Parameters[i].Type.Name)
			}
		}
	}
}
