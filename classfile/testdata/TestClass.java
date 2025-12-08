package testdata;

public class TestClass extends Object implements Runnable {
    public static final int CONSTANT_VALUE = 42;
    private String name;
    protected int count;

    public TestClass() {
        this.name = "default";
    }

    public TestClass(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    private static int helper(int x, int y) {
        return x + y;
    }

    @Override
    public void run() {
        System.out.println("Running: " + name);
    }
}
