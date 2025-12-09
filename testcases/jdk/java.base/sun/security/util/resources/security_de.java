package sun.security.util.resources;

import java.util.ListResourceBundle;

public final class security_de extends ListResourceBundle {
    protected final Object[][] getContents() {
        return new Object[][] {
            { ".Principal.", "\tPrincipal: " },
            { ".Private.Credential.", "\tPrivate Zugangsdaten: " },
            { ".Private.Credential.inaccessible.", "\tKein Zugriff auf private Zugangsdaten\n" },
            { ".Public.Credential.", "\t\u00D6ffentliche Zugangsdaten: " },
            { "CredOwner.Principal.Class.class.Principal.Name.name", "CredOwner:\n\tPrincipal-Klasse = {0}\n\tPrincipal-Name = {1}" },
            { "Credential.Class.not.followed.by.a.Principal.Class.and.Name", "Nach Zugangsdatenklasse folgt keine Principal-Klasse und kein Name" },
            { "Invalid.null.input.name", "Ung\u00FCltige Nulleingabe: Name" },
            { "Login.Failure.all.modules.ignored", "Anmeldefehler: Alle Module werden ignoriert" },
            { "LoginModuleControlFlag.", "LoginModuleControlFlag: " },
            { "NEWLINE", "\n" },
            { "No.LoginModules.configured.for.name", "F\u00FCr {0} sind keine LoginModules konfiguriert" },
            { "PKCS11.Token.providerName.Password.", "Kennwort f\u00FCr PKCS11-Token [{0}]: " },
            { "Principal.Class.not.followed.by.a.Principal.Name", "Nach Principal-Klasse folgt kein Principal-Name" },
            { "Principal.Name.missing.end.quote", "Abschlie\u00DFendes Anf\u00FChrungszeichen f\u00FCr Principal-Name fehlt" },
            { "Principal.Name.must.be.surrounded.by.quotes", "Principal-Name muss in Anf\u00FChrungszeichen stehen" },
            { "PrivateCredentialPermission.Principal.Class.can.not.be.a.wildcard.value.if.Principal.Name.is.not.a.wildcard.value", "Principal-Klasse PrivateCredentialPermission kann kein Platzhalterwert (*) sein, wenn der Principal-Name kein Platzhalterwert (*) ist" },
            { "Subject.", "Subjekt:\n" },
            { "Subject.is.read.only", "Subjekt ist schreibgesch\u00FCtzt" },
            { "actions.can.only.be.read.", "Aktionen k\u00F6nnen nur \"lesen\" sein" },
            { "attempting.to.add.an.object.which.is.not.an.instance.of.class", "Es wird versucht, ein Objekt hinzuzuf\u00FCgen, das keine Instanz von {0} ist" },
            { "attempting.to.add.an.object.which.is.not.an.instance.of.java.security.Principal.to.a.Subject.s.Principal.Set", "Es wird versucht, ein Objekt hinzuzuf\u00FCgen, das keine Instanz von java.security.Principal f\u00FCr eine Principal-Gruppe eines Subjekts ist" },
            { "duplicate.keystore.domain.name", "Keystore-Domainname doppelt vorhanden: {0}" },
            { "duplicate.keystore.name", "Keystore-Name doppelt vorhanden: {0}" },
            { "expected.expect.read.end.of.file.", "[{0}] erwartet, [Dateiende] gelesen" },
            { "expected.read.end.of.file.", "[;] erwartet, [Dateiende] gelesen" },
            { "invalid.null.CallbackHandler.provided", "Ung\u00FCltiger Nullwert f\u00FCr CallbackHandler angegeben" },
            { "invalid.null.Class.provided", "Ung\u00FCltige Nullklasse angegeben" },
            { "invalid.null.Subject.provided", "Ung\u00FCltiges Nullsubjekt angegeben" },
            { "invalid.null.action.provided", "Ung\u00FCltige Nullaktion angegeben" },
            { "invalid.null.input.s.", "Ung\u00FCltige Nulleingabe(n)" },
            { "line.number.expected.expect.found.actual.", "Zeile {0}: [{1}] erwartet, [{2}] gefunden" },
            { "line.number.msg", "Zeile {0}: {1}" },
            { "null.subject.logout.called.before.login", "Nullsubjekt - Abmeldung vor Anmeldung aufgerufen" },
            { "number.", "Nummer " },
            { "permission.name.name.syntax.invalid.", "Syntax f\u00FCr Berechtigungsnamen [{0}] ung\u00FCltig: " },
            { "provided.null.OID.map", "Null-OID-Zuordnung angegeben" },
            { "provided.null.keyword.map", "Null-Schl\u00FCsselwortzuordnung angegeben" },
            { "provided.null.name", "Nullname angegeben" },
        };
    }
}
