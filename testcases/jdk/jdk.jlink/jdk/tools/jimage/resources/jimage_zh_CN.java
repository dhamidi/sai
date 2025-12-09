package jdk.tools.jimage.resources;

import java.util.ListResourceBundle;

public final class jimage_zh_CN extends ListResourceBundle {
    protected final Object[][] getContents() {
        return new Object[][] {
            { "err.ambiguous.arg", "\u9009\u9879 {0} \u7684\u503C\u4EE5 \"--\" \u5F00\u5934\uFF0C\u5E94\u4F7F\u7528 {0}=<value> \u683C\u5F0F" },
            { "err.cannot.create.dir", "\u65E0\u6CD5\u521B\u5EFA\u76EE\u5F55 {0}" },
            { "err.invalid.jimage", "\u65E0\u6CD5\u6253\u5F00 {0}\uFF1A {1}" },
            { "err.missing.arg", "\u6CA1\u6709\u4E3A{0}\u6307\u5B9A\u503C" },
            { "err.no.jimage", "\u672A\u63D0\u4F9B jimage" },
            { "err.not.a.dir", "\u4E0D\u662F\u76EE\u5F55\uFF1A{0}" },
            { "err.not.a.jimage", "\u4E0D\u662F jimage \u6587\u4EF6\uFF1A{0}" },
            { "err.not.a.task", "\u4EFB\u52A1\u5FC5\u987B\u662F <extract | info | list | verify> \u4E4B\u4E00\uFF1A{0}" },
            { "err.option.unsupported", "\u4E0D\u652F\u6301{0}: {1}" },
            { "err.unknown.option", "\u672A\u77E5\u9009\u9879: {0}" },
            { "error.prefix", "\u9519\u8BEF:" },
            { "main.command.files", "       @<\u6587\u4EF6\u540D>                           \u4ECE\u6587\u4EF6\u4E2D\u8BFB\u53D6\u9009\u9879" },
            { "main.opt.dir", "          --dir                        \u63D0\u53D6\u6307\u4EE4\u7684\u76EE\u6807\u76EE\u5F55" },
            { "main.opt.footer", "\n\u5BF9\u4E8E\u9700\u8981 <pattern-list> \u7684\u9009\u9879\uFF0C\u503C\u5C06\u4E3A\u9017\u53F7\u5206\u9694\u7684\u5143\u7D20\u5217\u8868\uFF0C\n\u6BCF\u4E2A\u5143\u7D20\u4F7F\u7528\u4EE5\u4E0B\u683C\u5F0F\u4E4B\u4E00\uFF1A\n  <glob-pattern>\n  glob:<glob-pattern>\n  regex:<regex-pattern>" },
            { "main.opt.full-version", "          --full-version               \u8F93\u51FA\u5B8C\u6574\u7248\u672C\u4FE1\u606F" },
            { "main.opt.help", "  -?, -h, --help                       \u8F93\u51FA\u6B64\u5E2E\u52A9\u6D88\u606F" },
            { "main.opt.include", "          --include <pattern-list>     \u7528\u4E8E\u7B5B\u9009\u6761\u76EE\u7684\u6A21\u5F0F\u5217\u8868\u3002" },
            { "main.opt.verbose", "          --verbose                    \u5217\u8868\u8F93\u51FA\u6761\u76EE\u5927\u5C0F\u548C\u504F\u79FB\u91CF\n                                       \u5C5E\u6027" },
            { "main.opt.version", "          --version                    \u8F93\u51FA\u7248\u672C\u4FE1\u606F" },
            { "main.usage", "\u7528\u6CD5\uFF1A{0} <extract | info | list | verify> <options> jimage...\n\n  extract  - \u63D0\u53D6\u6240\u6709 jimage \u6761\u76EE\uFF0C\u5E76\u5C06\u5176\u653E\u7F6E\u5728\n             \u7531 --dir=<directory> (default=''.'') \u9009\u9879\u6307\u5B9A\u7684\u76EE\u5F55\u4E2D\u3002\n\n  info     - \u8F93\u51FA jimage \u6807\u5934\u4E2D\u5305\u542B\u7684\u8BE6\u7EC6\u4FE1\u606F\u3002\n\n  list     - \u8F93\u51FA jimage \u4E2D\u6240\u6709\u6761\u76EE\u7684\u540D\u79F0\u3002\u4E0E\n             --verbose \u4E00\u8D77\u4F7F\u7528\u65F6\uFF0Clist \u8FD8\u5C06\u8F93\u51FA\u6761\u76EE\u5927\u5C0F\u548C\u504F\u79FB\u91CF\u5C5E\u6027\u3002\n\n  verify   - \u62A5\u544A\u4EFB\u4F55\u672A\u9A8C\u8BC1\u4E3A\u7C7B\u7684 .class \u6761\u76EE\u3002\n\n\u53EF\u80FD\u7684\u9009\u9879\u5305\u62EC\uFF1A" },
            { "main.usage.extract", "  extract  - \u63D0\u53D6\u6240\u6709 jimage \u6761\u76EE\uFF0C\u5E76\u5C06\u5176\u653E\u7F6E\u5728\n             \u7531 --dir=<directory> (default='.') \u9009\u9879\u6307\u5B9A\u7684\u76EE\u5F55\u4E2D\u3002" },
            { "main.usage.info", "  info     - \u8F93\u51FA jimage \u6807\u5934\u4E2D\u5305\u542B\u7684\u8BE6\u7EC6\u4FE1\u606F\u3002" },
            { "main.usage.list", "  list     - \u8F93\u51FA jimage \u4E2D\u6240\u6709\u6761\u76EE\u7684\u540D\u79F0\u3002\u4E0E \n             --verbose \u4E00\u8D77\u4F7F\u7528\u65F6\uFF0Clist \u8FD8\u5C06\u8F93\u51FA\u6761\u76EE\u5927\u5C0F\u548C\u504F\u79FB\u91CF\u5C5E\u6027\u3002" },
            { "main.usage.summary", "\u7528\u6CD5\uFF1A{0} <extract | info | list | verify> <options> jimage...\n\u4F7F\u7528 -h \u6216 --help \u5217\u51FA\u53EF\u80FD\u7684\u9009\u9879\u3002" },
            { "main.usage.verify", "  verify   - \u62A5\u544A\u4EFB\u4F55\u672A\u9A8C\u8BC1\u4E3A\u7C7B\u7684 .class \u6761\u76EE\u7684\u9519\u8BEF\u3002" },
            { "warn.prefix", "\u8B66\u544A:" },
        };
    }
}
