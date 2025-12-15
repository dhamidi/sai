package test;

import java.util.concurrent.Flow;
import java.util.function.Function;
import java.util.function.Consumer;
import java.util.function.Predicate;
import java.util.stream.Stream;
import java.util.List;
import java.util.Optional;

public class ReactiveStreamsLambdaTest {

    // Simple lambda expressions
    private Consumer<String> printer = s -> System.out.println(s);
    private Function<String, Integer> lengthFunc = s -> s.length();
    private Predicate<Integer> isPositive = n -> n > 0;

    // Lambda with block body
    private Function<String, String> transformer = s -> {
        String result = s.trim();
        return result.toUpperCase();
    };

    // Multi-parameter lambda
    private java.util.Comparator<String> comparator = (a, b) -> a.compareTo(b);

    // Method reference
    private Consumer<String> methodRef = System.out::println;

    // Reactive streams subscriber
    public void createSubscriber() {
        Flow.Subscriber<String> subscriber = new Flow.Subscriber<>() {
            private Flow.Subscription subscription;

            @Override
            public void onSubscribe(Flow.Subscription subscription) {
                this.subscription = subscription;
                subscription.request(1);
            }

            @Override
            public void onNext(String item) {
                System.out.println("Received: " + item);
                subscription.request(1);
            }

            @Override
            public void onError(Throwable throwable) {
                throwable.printStackTrace();
            }

            @Override
            public void onComplete() {
                System.out.println("Done");
            }
        };
    }

    // Stream operations with lambdas
    public void streamOperations(List<String> items) {
        items.stream()
            .filter(s -> s != null)
            .filter(s -> !s.isEmpty())
            .map(s -> s.toLowerCase())
            .map(String::trim)
            .forEach(s -> System.out.println(s));

        // Chained operations with method references
        long count = items.stream()
            .filter(java.util.Objects::nonNull)
            .count();

        // Reduce with lambda
        Optional<String> reduced = items.stream()
            .reduce((a, b) -> a + ", " + b);

        // Collect with groupingBy
        var grouped = items.stream()
            .collect(java.util.stream.Collectors.groupingBy(
                s -> s.charAt(0),
                java.util.stream.Collectors.counting()
            ));
    }

    // Nested lambdas
    public Function<Integer, Function<Integer, Integer>> curriedAdd() {
        return a -> b -> a + b;
    }

    // Lambda returning lambda
    public void nestedLambdas() {
        Function<Integer, Function<Integer, Integer>> adder = x -> y -> x + y;
        int result = adder.apply(3).apply(4);
    }

    // Lambda with explicit types
    public void explicitTypeLambda() {
        java.util.function.BiFunction<String, Integer, String> repeat =
            (String s, Integer n) -> s.repeat(n);
    }

    // Lambda in conditional
    public Consumer<String> getConsumer(boolean verbose) {
        return verbose ? s -> System.out.println("VERBOSE: " + s) : s -> {};
    }

    // Lambda with exception handling
    public void lambdaWithTryCatch() {
        Function<String, Integer> parser = s -> {
            try {
                return Integer.parseInt(s);
            } catch (NumberFormatException e) {
                return 0;
            }
        };
    }
}
