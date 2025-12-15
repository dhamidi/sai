package test;

import java.io.IOException;
import java.security.GeneralSecurityException;

public class VarargsMultiCatchTest {

    // Varargs in method parameter
    public void printAll(String... items) {
        for (String item : items) {
            System.out.println(item);
        }
    }

    // Varargs with other parameters
    public int sumAll(int first, int... rest) {
        int total = first;
        for (int n : rest) {
            total += n;
        }
        return total;
    }

    // Varargs with generic type
    @SafeVarargs
    public final <T> void processAll(T... elements) {
        for (T elem : elements) {
            System.out.println(elem);
        }
    }

    // Varargs with array type
    public void processArrays(int[]... arrays) {
        for (int[] arr : arrays) {
            for (int n : arr) {
                System.out.println(n);
            }
        }
    }

    // Multi-catch exception handling
    public void multiCatch() {
        try {
            riskyOperation();
        } catch (IOException | GeneralSecurityException e) {
            System.err.println("Error: " + e.getMessage());
        }
    }

    // Multi-catch with more than two types
    public void multiCatchThree() {
        try {
            riskyOperation();
        } catch (IOException | IllegalArgumentException | IllegalStateException e) {
            e.printStackTrace();
        }
    }

    // Multi-catch followed by another catch
    public void multiCatchWithFallback() {
        try {
            riskyOperation();
        } catch (IOException | GeneralSecurityException e) {
            System.err.println("Known error: " + e);
        } catch (Exception e) {
            System.err.println("Unknown error: " + e);
        }
    }

    // Combined: varargs method that uses multi-catch
    public void processFiles(String... filenames) {
        for (String filename : filenames) {
            try {
                readFile(filename);
            } catch (IOException | SecurityException e) {
                System.err.println("Failed to read " + filename + ": " + e);
            }
        }
    }

    private void riskyOperation() throws IOException, GeneralSecurityException {
        // placeholder
    }

    private void readFile(String name) throws IOException {
        // placeholder
    }
}
