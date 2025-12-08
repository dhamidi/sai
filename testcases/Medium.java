public class Medium {
    private int counter;
    private String name;

    public Medium(String name) {
        this.name = name;
        this.counter = 0;
    }

    public void increment() {
        counter++;
    }

    public int getCounter() {
        return counter;
    }

    public String getName() {
        return name;
    }

    public static int factorial(int n) {
        if (n <= 1) return 1;
        return n * factorial(n - 1);
    }

    public static void main(String[] args) {
        Medium obj = new Medium("TestObject");
        for (int i = 0; i < 5; i++) {
            obj.increment();
        }
        System.out.println(obj.getName() + " counter: " + obj.getCounter());
        System.out.println("Factorial of 6: " + factorial(6));
    }
}
