package org.eclipse.jetty.client;

public class Authentication {
    public static class HeaderInfo {
        private String name;
        private String value;
        
        public HeaderInfo(String name, String value) {
            this.name = name;
            this.value = value;
        }
    }
    
    public HeaderInfo createHeader(String n, String v) {
        return new HeaderInfo(n, v);
    }
}
