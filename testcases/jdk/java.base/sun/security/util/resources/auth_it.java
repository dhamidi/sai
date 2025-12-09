package sun.security.util.resources;

import java.util.ListResourceBundle;

public final class auth_it extends ListResourceBundle {
    protected final Object[][] getContents() {
        return new Object[][] {
            { "Configuration.Error.Can.not.specify.multiple.entries.for.appName", "Errore di configurazione:\n\timpossibile specificare pi\u00F9 valori per {0}" },
            { "Configuration.Error.Invalid.control.flag.flag", "Errore di configurazione:\n\tflag di controllo non valido, {0}" },
            { "Configuration.Error.Line.line.expected.expect.", "Errore di configurazione:\n\triga {0}: previsto [{1}]" },
            { "Configuration.Error.Line.line.expected.expect.found.value.", "Errore di configurazione:\n\triga {0}: previsto [{1}], trovato [{2}]" },
            { "Configuration.Error.Line.line.system.property.value.expanded.to.empty.value", "Errore di configurazione:\n\triga {0}: propriet\u00E0 di sistema [{1}] espansa a valore vuoto" },
            { "Configuration.Error.No.such.file.or.directory", "Errore di configurazione:\n\tFile o directory inesistente" },
            { "Configuration.Error.expected.expect.read.end.of.file.", "Errore di configurazione:\n\tprevisto [{0}], letto [end of file]" },
            { "Invalid.NTSid.value", "Valore NTSid non valido" },
            { "Kerberos.password.for.username.", "Password Kerberos per {0}: " },
            { "Kerberos.username.defUsername.", "Nome utente Kerberos [{0}]: " },
            { "Keystore.alias.", "Alias keystore: " },
            { "Keystore.password.", "Password keystore: " },
            { "NTDomainPrincipal.name", "NTDomainPrincipal: {0}" },
            { "NTNumericCredential.name", "NTNumericCredential: {0}" },
            { "NTSid.name", "NTSid: {0}" },
            { "NTSidDomainPrincipal.name", "NTSidDomainPrincipal: {0}" },
            { "NTSidGroupPrincipal.name", "NTSidGroupPrincipal: {0}" },
            { "NTSidPrimaryGroupPrincipal.name", "NTSidPrimaryGroupPrincipal: {0}" },
            { "NTSidUserPrincipal.name", "NTSidUserPrincipal: {0}" },
            { "NTUserPrincipal.name", "NTUserPrincipal: {0}" },
            { "Please.enter.keystore.information", "Immettere le informazioni per il keystore" },
            { "Private.key.password.optional.", "Password chiave privata (opzionale): " },
            { "Unable.to.properly.expand.config", "Impossibile espandere correttamente {0}" },
            { "UnixNumericGroupPrincipal.Primary.Group.name", "UnixNumericGroupPrincipal [gruppo primario]: {0}" },
            { "UnixNumericGroupPrincipal.Supplementary.Group.name", "UnixNumericGroupPrincipal [gruppo supplementare]: {0}" },
            { "UnixNumericUserPrincipal.name", "UnixNumericUserPrincipal: {0}" },
            { "UnixPrincipal.name", "UnixPrincipal: {0}" },
            { "extra.config.No.such.file.or.directory.", "{0} (file o directory inesistente)" },
            { "invalid.null.input.value", "input nullo non valido: {0}" },
            { "password.", "Password: " },
            { "username.", "Nome utente: " },
        };
    }
}
