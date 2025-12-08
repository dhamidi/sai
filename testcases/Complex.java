import java.util.ArrayList;
import java.util.List;

interface Drawable {
    void draw();
}

interface Resizable {
    void resize(double factor);
}

abstract class Shape implements Drawable {
    protected String color;

    public Shape(String color) {
        this.color = color;
    }

    public abstract double area();

    @Override
    public String toString() {
        return getClass().getSimpleName() + "[color=" + color + ", area=" + area() + "]";
    }
}

class Circle extends Shape implements Resizable {
    private double radius;

    public Circle(String color, double radius) {
        super(color);
        this.radius = radius;
    }

    @Override
    public double area() {
        return Math.PI * radius * radius;
    }

    @Override
    public void draw() {
        System.out.println("Drawing circle with radius " + radius);
    }

    @Override
    public void resize(double factor) {
        radius *= factor;
    }
}

class Rectangle extends Shape implements Resizable {
    private double width;
    private double height;

    public Rectangle(String color, double width, double height) {
        super(color);
        this.width = width;
        this.height = height;
    }

    @Override
    public double area() {
        return width * height;
    }

    @Override
    public void draw() {
        System.out.println("Drawing rectangle " + width + "x" + height);
    }

    @Override
    public void resize(double factor) {
        width *= factor;
        height *= factor;
    }
}

public class Complex {
    public static void main(String[] args) {
        List<Shape> shapes = new ArrayList<>();
        shapes.add(new Circle("red", 5.0));
        shapes.add(new Rectangle("blue", 4.0, 6.0));
        shapes.add(new Circle("green", 3.0));

        System.out.println("Original shapes:");
        for (Shape s : shapes) {
            s.draw();
            System.out.println("  " + s);
        }

        System.out.println("\nAfter resizing by 2x:");
        for (Shape s : shapes) {
            if (s instanceof Resizable) {
                ((Resizable) s).resize(2.0);
            }
            System.out.println("  " + s);
        }
    }
}
