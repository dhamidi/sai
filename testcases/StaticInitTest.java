package test;

public class StaticInitTest {
    private static final int VALUE;

    static {
        // initialize the value
        VALUE = 42;
    }

    public int getValue() {
        return VALUE;
    }
}
