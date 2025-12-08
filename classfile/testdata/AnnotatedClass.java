package testdata;

import java.lang.annotation.*;

@Retention(RetentionPolicy.RUNTIME)
@Target({ElementType.TYPE, ElementType.METHOD, ElementType.FIELD, ElementType.PARAMETER, ElementType.CONSTRUCTOR})
@interface MyAnnotation {
    String value() default "";
    int count() default 0;
}

@Retention(RetentionPolicy.CLASS)
@Target({ElementType.TYPE, ElementType.METHOD})
@interface ClassOnlyAnnotation {
    String[] tags() default {};
}

@Deprecated
@MyAnnotation(value = "class level", count = 42)
@ClassOnlyAnnotation(tags = {"test", "example"})
public class AnnotatedClass<T extends Comparable<T>> {
    @Deprecated
    @MyAnnotation(value = "deprecated field")
    public static final long LONG_CONST = 123456789012345L;
    
    public static final double DOUBLE_CONST = 3.14159265359;
    public static final float FLOAT_CONST = 2.71828f;
    
    private T value;
    
    @MyAnnotation(value = "constructor")
    public AnnotatedClass(@MyAnnotation(value = "param") T value) {
        this.value = value;
    }
    
    @Deprecated
    @MyAnnotation(value = "method annotation", count = 100)
    public T getValue() {
        int localVar = 42;
        String message = "test";
        for (int i = 0; i < localVar; i++) {
            if (i % 2 == 0) {
                message = message + i;
            }
        }
        return value;
    }
    
    public void methodWithException() throws Exception, RuntimeException {
        throw new Exception("test");
    }
    
    public void methodWithLambda() {
        Runnable r = () -> System.out.println("lambda");
        r.run();
    }
    
    public class InnerClass {
        public void innerMethod() {}
    }
    
    public static class StaticNestedClass {
        public void nestedMethod() {}
    }
    
    private class PrivateInner {
        public void test() {}
    }
    
    public void methodWithAnonymous() {
        new Runnable() {
            @Override
            public void run() {
                System.out.println("anonymous");
            }
        }.run();
    }
}
