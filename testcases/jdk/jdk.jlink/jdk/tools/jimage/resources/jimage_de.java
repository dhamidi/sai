package jdk.tools.jimage.resources;

import java.util.ListResourceBundle;

public final class jimage_de extends ListResourceBundle {
    protected final Object[][] getContents() {
        return new Object[][] {
            { "err.ambiguous.arg", "Wert f\u00FCr Option {0} beginnt mit \"--\", aber muss das Format {0}=<Wert> verwenden" },
            { "err.cannot.create.dir", "Verzeichnis {0} kann nicht erstellt werden" },
            { "err.invalid.jimage", "{0} kann nicht ge\u00F6ffnet werden: {1}" },
            { "err.missing.arg", "kein Wert angegeben f\u00FCr {0}" },
            { "err.no.jimage", "Kein jimage angegeben" },
            { "err.not.a.dir", "Kein Verzeichnis: {0}" },
            { "err.not.a.jimage", "Keine jimage-Datei: {0}" },
            { "err.not.a.task", "Aufgabe muss einen der folgenden Werte aufweisen: <extract | info | list | verify>: {0}" },
            { "err.option.unsupported", "{0} nicht unterst\u00FCtzt: {1}" },
            { "err.unknown.option", "unbekannte Option: {0}" },
            { "error.prefix", "Fehler:" },
            { "main.command.files", "       @<filename>                  Liest Optionen aus der Datei" },
            { "main.opt.dir", "          --dir                        Zielverzeichnis f\u00FCr extract-Direktive" },
            { "main.opt.footer", "\nBei Optionen, die eine <Musterliste> erfordern, ist der Wert eine durch Komma getrennte\nListe von Elementen, die jeweils eines der folgenden Formate verwenden:\n  <Glob-Muster>\n  glob:<Glob-Muster>\n  regex:<Regex-Muster>" },
            { "main.opt.full-version", "          --full-version               Gibt vollst\u00E4ndige Versionsinformationen aus" },
            { "main.opt.help", "  -?, -h, --help                       Gibt diese Hilfemeldung aus" },
            { "main.opt.include", "          --include <Musterliste>     Musterliste f\u00FCr das Filtern von Eintr\u00E4gen." },
            { "main.opt.verbose", "          --verbose                    Gibt beim Auflisten die Eintragsgr\u00F6\u00DFe\n                                       und Offset-Attribute aus" },
            { "main.opt.version", "          --version                    Gibt Versionsinformationen aus" },
            { "main.usage", "Verwendung: {0} <extract | info | list | verify> <Optionen> jimage...\n\n  extract  - Extrahiert alle jimage-Eintr\u00E4ge und platziert sie in einem mit der Option\n             --dir=<Verzeichnis> (Standardwert = \".\") angegebenen Verzeichnis.\n\n  info     - Gibt detaillierte Informationen aus, die im jimage-Header enthalten sind.\n\n  list     - Gibt die Namen aller Eintr\u00E4ge im jimage aus. Bei Verwendung mit\n             --verbose gibt der Befehl \"list\" auch die Eintragsgr\u00F6\u00DFe und Offset-Attribute aus.\n\n  verify   - Meldet alle .class-Eintr\u00E4ge, die nicht als Klassen verifiziert werden.\n\nM\u00F6gliche Optionen:" },
            { "main.usage.extract", "  extract  - Extrahiert alle jimage-Eintr\u00E4ge und platziert sie in einem mit der Option\n             --dir=<Verzeichnis> (Standardwert = \".\") angegebenen Verzeichnis." },
            { "main.usage.info", "  info     - Gibt detaillierte Informationen aus, die im jimage-Header enthalten sind." },
            { "main.usage.list", "  list     - Gibt die Namen aller Eintr\u00E4ge im jimage aus. Bei Verwendung mit\n             --verbose gibt der Befehl \"list\" auch die Eintragsgr\u00F6\u00DFe und Offset-Attribute aus." },
            { "main.usage.summary", "Verwendung: {0} <extract | info | list | verify> <options> jimage...\nMit -h oder --help k\u00F6nnen Sie eine Liste der m\u00F6glichen Optionen aufrufen." },
            { "main.usage.verify", "  verify   - Meldet alle .class-Eintr\u00E4ge, die nicht als Klassen verifiziert werden." },
            { "warn.prefix", "Warnung:" },
        };
    }
}
