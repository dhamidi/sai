package testdata;

public record RecordClass(String name, int value) {
    public RecordClass {
        if (value < 0) {
            throw new IllegalArgumentException("value must be non-negative");
        }
    }
    
    public String fullDescription() {
        return name + ": " + value;
    }
}
