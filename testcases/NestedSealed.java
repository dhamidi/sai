package testcases;

public class NestedSealed {
    public sealed interface Content permits Content.Text, Content.Image {
        record Text(String text) implements Content {}
        record Image(String url) implements Content {}

        default Map<String, Object> toJson() {
            return switch (this) {
                case Text t -> Map.of("type", "text", "text", t.text());
                case Image i -> Map.of("type", "image", "url", i.url());
            };
        }

        static Content fromJson(Map<String, Object> m) {
            var type = (String) m.get("type");
            return switch (type) {
                case "text" -> new Text((String) m.get("text"));
                case "image" -> new Image((String) m.get("url"));
                default -> throw new IllegalArgumentException("Unknown content type: " + type);
            };
        }
    }
}
