package testdata;

public sealed class SealedClass permits SubClass1, SubClass2 {
    public void baseMethod() {}
}

final class SubClass1 extends SealedClass {
    public void method1() {}
}

final class SubClass2 extends SealedClass {
    public void method2() {}
}
