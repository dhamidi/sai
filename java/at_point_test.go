package java

import (
	"bytes"
	"testing"

	"github.com/dhamidi/javalyzer/java/parser"
)

func TestTypeAtPoint(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		line     int
		column   int
		wantType string
	}{
		{
			name: "local var with explicit type",
			source: `import java.util.ArrayList;
public class Example {
  public static void main(String[] args) {
    ArrayList<String> list = new ArrayList<>();
    list.
  }
}`,
			line:     5,
			column:   5, // on 'l' of list.
			wantType: "java.util.ArrayList",
		},
		{
			name: "local var with var keyword",
			source: `import java.util.ArrayList;
public class Example {
  public static void main(String[] args) {
    var list = new ArrayList<String>();
    list.
  }
}`,
			line:     5,
			column:   5, // on 'l' of list.
			wantType: "java.util.ArrayList",
		},
		{
			name: "parameter type resolution",
			source: `public class Example {
  public static void main(String[] args) {
    args.
  }
}`,
			line:     3,
			column:   5, // on 'a' of args.
			wantType: "java.lang.String[]",
		},
		{
			name: "nested inner class",
			source: `public class Outer {
  public class Inner {
    public void doSomething() {}
  }
  public void test() {
    Inner inner = new Inner();
    inner.
  }
}`,
			line:     7,
			column:   5, // on 'i' of inner.
			wantType: "Outer.Inner",
		},
		{
			name: "field via this",
			source: `import java.util.List;
public class Example {
  private List<String> items;
  public void test() {
    this.items.
  }
}`,
			line:     5,
			column:   10, // on 'i' of items
			wantType: "java.util.List",
		},
		{
			name: "variable shadowing",
			source: `import java.util.List;
import java.util.ArrayList;
public class Example {
  private List<String> items;
  public void test() {
    ArrayList<Integer> items = new ArrayList<>();
    items.
  }
}`,
			line:     7,
			column:   5, // on 'i' of items (local)
			wantType: "java.util.ArrayList",
		},
		{
			name: "for-each loop variable",
			source: `import java.util.List;
public class Example {
  public void test(List<String> items) {
    for (String item : items) {
      item.
    }
  }
}`,
			line:     5,
			column:   7, // on 'i' of item
			wantType: "java.lang.String",
		},
		{
			name: "catch clause exception",
			source: `import java.io.IOException;
public class Example {
  public void test() {
    try {
    } catch (IOException e) {
      e.
    }
  }
}`,
			line:     6,
			column:   7, // on 'e' of e.
			wantType: "java.io.IOException",
		},
		{
			name: "try-with-resources variable",
			source: `import java.io.FileInputStream;
public class Example {
  public void test() throws Exception {
    try (FileInputStream fis = new FileInputStream("test")) {
      fis.
    }
  }
}`,
			line:     5,
			column:   7, // on 'f' of fis.
			wantType: "java.io.FileInputStream",
		},
		{
			name: "lambda parameter explicit type",
			source: `import java.util.function.Consumer;
public class Example {
  Consumer<String> c = (String s) -> {
    s.
  };
}`,
			line:     4,
			column:   5, // on 's' of s.
			wantType: "java.lang.String",
		},
		{
			name: "multi-dimensional array",
			source: `public class Example {
  public void test() {
    int[][] matrix = new int[3][3];
    matrix.
  }
}`,
			line:     4,
			column:   5, // on 'm' of matrix.
			wantType: "int[][]",
		},
		{
			name: "constructor parameter",
			source: `import java.util.List;
public class Example {
  public Example(List<String> items) {
    items.
  }
}`,
			line:     4,
			column:   5, // on 'i' of items.
			wantType: "java.util.List",
		},
		{
			name: "record component",
			source: `import java.util.List;
public record Person(String name, List<String> tags) {
  public void test() {
    name.
  }
}`,
			line:     4,
			column:   5, // on 'n' of name.
			wantType: "java.lang.String",
		},
		{
			name: "pattern matching instanceof",
			source: `public class Example {
  public void test(Object obj) {
    if (obj instanceof String s) {
      s.
    }
  }
}`,
			line:     4,
			column:   7, // on 's' of s.
			wantType: "java.lang.String",
		},
		{
			name: "traditional for loop variable",
			source: `public class Example {
  public void test() {
    for (Integer i = 0; i < 10; i++) {
      i.
    }
  }
}`,
			line:     4,
			column:   7, // on 'i' of i.
			wantType: "java.lang.Integer",
		},
		{
			name: "multiple variables same declaration",
			source: `public class Example {
  public void test() {
    String a, b, c;
    b.
  }
}`,
			line:     4,
			column:   5, // on 'b' of b.
			wantType: "java.lang.String",
		},
		{
			name: "varargs parameter",
			source: `public class Example {
  public void test(String... args) {
    args.
  }
}`,
			line:     3,
			column:   5, // on 'a' of args.
			wantType: "java.lang.String[]",
		},
		{
			name: "switch pattern variable",
			source: `public class Example {
  public void test(Object obj) {
    switch (obj) {
      case String s -> {
        s.
      }
      default -> {}
    }
  }
}`,
			line:     5,
			column:   9, // on 's' of s.
			wantType: "java.lang.String",
		},
		{
			name: "static field same class",
			source: `import java.util.concurrent.atomic.AtomicInteger;
public class Example {
  private static AtomicInteger counter = new AtomicInteger();
  public void test() {
    counter.
  }
}`,
			line:     5,
			column:   5, // on 'c' of counter.
			wantType: "java.util.concurrent.atomic.AtomicInteger",
		},
		{
			name: "array element access",
			source: `public class Example {
  public void test(String[] args) {
    args[0].
  }
}`,
			line:     3,
			column:   5, // on 'a' of args[0].
			wantType: "java.lang.String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.ParseCompilationUnit(bytes.NewReader([]byte(tt.source)))
			node := p.Finish()
			if node == nil {
				t.Fatalf("failed to parse source")
			}

			pos := parser.Position{Line: tt.line, Column: tt.column}
			gotType := TypeAtPoint(node, pos)
			if gotType != tt.wantType {
				t.Errorf("TypeAtPoint() = %q, want %q", gotType, tt.wantType)
			}
		})
	}
}
