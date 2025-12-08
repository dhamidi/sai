package testdata;

import java.util.function.*;

public class ConstantPoolTest {
    public static final byte BYTE_CONST = 127;
    public static final short SHORT_CONST = 32767;
    public static final int INT_CONST = 2147483647;
    public static final long LONG_CONST = 9223372036854775807L;
    public static final float FLOAT_CONST = 3.4028235E38f;
    public static final double DOUBLE_CONST = 1.7976931348623157E308;
    public static final String STRING_CONST = "Hello, World!";
    public static final Class<?> CLASS_CONST = String.class;
    
    public static int staticMethod(int x) {
        return x * 2;
    }
    
    public int instanceMethod(int x) {
        return x + 1;
    }
    
    public void testMethodReferences() {
        IntUnaryOperator staticRef = ConstantPoolTest::staticMethod;
        IntUnaryOperator instanceRef = this::instanceMethod;
        
        Supplier<ConstantPoolTest> constructorRef = ConstantPoolTest::new;
        
        Function<String, Integer> parseRef = Integer::parseInt;
        
        int result1 = staticRef.applyAsInt(5);
        int result2 = instanceRef.applyAsInt(5);
    }
    
    public void testLambdas() {
        Runnable r = () -> System.out.println("simple lambda");
        
        IntBinaryOperator add = (a, b) -> a + b;
        
        Function<Integer, Integer> square = x -> x * x;
        
        Consumer<String> printer = s -> {
            String upper = s.toUpperCase();
            System.out.println(upper);
        };
        
        r.run();
        int sum = add.applyAsInt(3, 4);
        int squared = square.apply(5);
        printer.accept("test");
    }
    
    public void testAllPrimitives() {
        boolean boolVal = true;
        byte byteVal = 42;
        char charVal = 'A';
        short shortVal = 1000;
        int intVal = 100000;
        long longVal = 10000000000L;
        float floatVal = 1.5f;
        double doubleVal = 2.5;
        
        int[] intArray = new int[]{1, 2, 3};
        String[] stringArray = new String[]{"a", "b", "c"};
        Object[][] multiArray = new Object[2][3];
    }
    
    public interface NestedInterface {
        void interfaceMethod();
    }
    
    public void testInterfaceCall(NestedInterface ni) {
        ni.interfaceMethod();
    }
}
