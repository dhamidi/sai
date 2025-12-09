package sun.security.util.resources;

import java.util.ListResourceBundle;

public final class security_zh_CN extends ListResourceBundle {
    protected final Object[][] getContents() {
        return new Object[][] {
            { ".Principal.", "\t\u4E3B\u7528\u6237\uFF1A " },
            { ".Private.Credential.", "\t\u4E13\u7528\u8EAB\u4EFD\u8BC1\u660E\uFF1A " },
            { ".Private.Credential.inaccessible.", "\t\u65E0\u6CD5\u8BBF\u95EE\u4E13\u7528\u8EAB\u4EFD\u8BC1\u660E\n" },
            { ".Public.Credential.", "\t\u516C\u5171\u8EAB\u4EFD\u8BC1\u660E\uFF1A " },
            { "CredOwner.Principal.Class.class.Principal.Name.name", "CredOwner:\n\t\u4E3B\u7528\u6237\u7C7B = {0}\n\t\u4E3B\u7528\u6237\u540D\u79F0 = {1}" },
            { "Credential.Class.not.followed.by.a.Principal.Class.and.Name", "\u8EAB\u4EFD\u8BC1\u660E\u7C7B\u540E\u9762\u672A\u8DDF\u968F\u4E3B\u7528\u6237\u7C7B\u53CA\u540D\u79F0" },
            { "Invalid.null.input.name", "\u65E0\u6548\u7A7A\u8F93\u5165: \u540D\u79F0" },
            { "Login.Failure.all.modules.ignored", "\u767B\u5F55\u5931\u8D25: \u5FFD\u7565\u6240\u6709\u6A21\u5757" },
            { "LoginModuleControlFlag.", "LoginModuleControlFlag\uFF1A " },
            { "NEWLINE", "\n" },
            { "No.LoginModules.configured.for.name", "\u6CA1\u6709\u4E3A{0}\u914D\u7F6E LoginModules" },
            { "PKCS11.Token.providerName.Password.", "PKCS11 \u6807\u8BB0 [{0}] \u5BC6\u7801\uFF1A " },
            { "Principal.Class.not.followed.by.a.Principal.Name", "\u4E3B\u7528\u6237\u7C7B\u540E\u9762\u672A\u8DDF\u968F\u4E3B\u7528\u6237\u540D\u79F0" },
            { "Principal.Name.missing.end.quote", "\u4E3B\u7528\u6237\u540D\u79F0\u7F3A\u5C11\u53F3\u5F15\u53F7" },
            { "Principal.Name.must.be.surrounded.by.quotes", "\u4E3B\u7528\u6237\u540D\u79F0\u5FC5\u987B\u653E\u5728\u5F15\u53F7\u5185" },
            { "PrivateCredentialPermission.Principal.Class.can.not.be.a.wildcard.value.if.Principal.Name.is.not.a.wildcard.value", "\u5982\u679C\u4E3B\u7528\u6237\u540D\u79F0\u4E0D\u662F\u901A\u914D\u7B26 (*) \u503C, \u90A3\u4E48 PrivateCredentialPermission \u4E3B\u7528\u6237\u7C7B\u4E0D\u80FD\u662F\u901A\u914D\u7B26 (*) \u503C" },
            { "Subject.", "\u4E3B\u4F53\uFF1A\n" },
            { "Subject.is.read.only", "\u4E3B\u4F53\u4E3A\u53EA\u8BFB" },
            { "actions.can.only.be.read.", "\u64CD\u4F5C\u53EA\u80FD\u4E3A '\u8BFB\u53D6'" },
            { "attempting.to.add.an.object.which.is.not.an.instance.of.class", "\u6B63\u5728\u5C1D\u8BD5\u6DFB\u52A0\u4E00\u4E2A\u975E{0}\u5B9E\u4F8B\u7684\u5BF9\u8C61" },
            { "attempting.to.add.an.object.which.is.not.an.instance.of.java.security.Principal.to.a.Subject.s.Principal.Set", "\u6B63\u5728\u5C1D\u8BD5\u5C06\u4E00\u4E2A\u975E java.security.Principal \u5B9E\u4F8B\u7684\u5BF9\u8C61\u6DFB\u52A0\u5230\u4E3B\u4F53\u7684\u4E3B\u7528\u6237\u96C6\u4E2D" },
            { "duplicate.keystore.domain.name", "\u5BC6\u94A5\u5E93\u57DF\u540D\u91CD\u590D: {0}" },
            { "duplicate.keystore.name", "\u5BC6\u94A5\u5E93\u540D\u79F0\u91CD\u590D: {0}" },
            { "expected.expect.read.end.of.file.", "\u5E94\u4E3A [{0}], \u8BFB\u53D6\u7684\u662F [\u6587\u4EF6\u7ED3\u5C3E]" },
            { "expected.read.end.of.file.", "\u5E94\u4E3A [;], \u8BFB\u53D6\u7684\u662F [\u6587\u4EF6\u7ED3\u5C3E]" },
            { "invalid.null.CallbackHandler.provided", "\u63D0\u4F9B\u4E86\u65E0\u6548\u7684\u7A7A CallbackHandler" },
            { "invalid.null.Class.provided", "\u63D0\u4F9B\u4E86\u65E0\u6548\u7684\u7A7A\u7C7B" },
            { "invalid.null.Subject.provided", "\u63D0\u4F9B\u4E86\u65E0\u6548\u7684\u7A7A\u4E3B\u4F53" },
            { "invalid.null.action.provided", "\u63D0\u4F9B\u4E86\u65E0\u6548\u7684\u7A7A\u64CD\u4F5C" },
            { "invalid.null.input.s.", "\u65E0\u6548\u7684\u7A7A\u8F93\u5165" },
            { "line.number.expected.expect.found.actual.", "\u884C\u53F7 {0}: \u5E94\u4E3A [{1}], \u627E\u5230 [{2}]" },
            { "line.number.msg", "\u7B2C {0} \u884C\uFF1A{1}" },
            { "null.subject.logout.called.before.login", "\u7A7A\u4E3B\u4F53 - \u5728\u767B\u5F55\u4E4B\u524D\u8C03\u7528\u4E86\u6CE8\u9500" },
            { "number.", "\u7F16\u53F7 " },
            { "permission.name.name.syntax.invalid.", "\u6743\u9650\u540D\u79F0 [{0}] \u8BED\u6CD5\u65E0\u6548: " },
            { "provided.null.OID.map", "\u63D0\u4F9B\u7684 OID \u6620\u5C04\u4E3A\u7A7A\u503C" },
            { "provided.null.keyword.map", "\u63D0\u4F9B\u7684\u5173\u952E\u5B57\u6620\u5C04\u4E3A\u7A7A\u503C" },
            { "provided.null.name", "\u63D0\u4F9B\u7684\u540D\u79F0\u4E3A\u7A7A\u503C" },
        };
    }
}
