/*
 * Copyright (c) 2012, 2025, Oracle and/or its affiliates. All rights reserved.
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.
 *
 * This code is free software; you can redistribute it and/or modify it
 * under the terms of the GNU General Public License version 2 only, as
 * published by the Free Software Foundation.  Oracle designates this
 * particular file as subject to the "Classpath" exception as provided
 * by Oracle in the LICENSE file that accompanied this code.
 *
 * This code is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
 * version 2 for more details (a copy is included in the LICENSE file that
 * accompanied this code).
 *
 * You should have received a copy of the GNU General Public License version
 * 2 along with this work; if not, write to the Free Software Foundation,
 * Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301 USA.
 *
 * Please contact Oracle, 500 Oracle Parkway, Redwood Shores, CA 94065 USA
 * or visit www.oracle.com if you need additional information or have any
 * questions.
 */

/*
 * COPYRIGHT AND PERMISSION NOTICE
 *
 * Copyright (c) 1991-2022 Unicode, Inc. All rights reserved.
 * Distributed under the Terms of Use in https://www.unicode.org/copyright.html.
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of the Unicode data files and any associated documentation
 * (the "Data Files") or Unicode software and any associated documentation
 * (the "Software") to deal in the Data Files or Software
 * without restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, and/or sell copies of
 * the Data Files or Software, and to permit persons to whom the Data Files
 * or Software are furnished to do so, provided that either
 * (a) this copyright and permission notice appear with all copies
 * of the Data Files or Software, or
 * (b) this copyright and permission notice appear in associated
 * Documentation.
 *
 * THE DATA FILES AND SOFTWARE ARE PROVIDED "AS IS", WITHOUT WARRANTY OF
 * ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT OF THIRD PARTY RIGHTS.
 * IN NO EVENT SHALL THE COPYRIGHT HOLDER OR HOLDERS INCLUDED IN THIS
 * NOTICE BE LIABLE FOR ANY CLAIM, OR ANY SPECIAL INDIRECT OR CONSEQUENTIAL
 * DAMAGES, OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE,
 * DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
 * TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
 * PERFORMANCE OF THE DATA FILES OR SOFTWARE.
 *
 * Except as contained in this notice, the name of a copyright holder
 * shall not be used in advertising or otherwise to promote the sale,
 * use or other dealings in these Data Files or Software without prior
 * written authorization of the copyright holder.
 */

package sun.text.resources.cldr.ext;

import java.util.ListResourceBundle;

public class FormatData_lrc extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] metaValue_MonthNames = new String[] {
            "\u062c\u0627\u0646\u06a4\u06cc\u06d5",
            "\u0641\u0626\u06a4\u0631\u06cc\u06d5",
            "\u0645\u0627\u0631\u0633",
            "\u0622\u06a4\u0631\u06cc\u0644",
            "\u0645\u0626\u06cc",
            "\u062c\u0648\u0659\u0623\u0646",
            "\u062c\u0648\u0659\u0644\u0627",
            "\u0622\u06af\u0648\u0633\u062a",
            "\u0633\u0626\u067e\u062a\u0627\u0645\u0631",
            "\u0626\u0648\u06a9\u062a\u0648\u06a4\u0631",
            "\u0646\u0648\u06a4\u0627\u0645\u0631",
            "\u062f\u0626\u0633\u0627\u0645\u0631",
            "",
        };
        final String metaValue_calendarname_gregorian = "\u062a\u0623\u0642\u06a4\u06cc\u0645 \u06af\u0623\u0631\u0626\u06af\u0648\u0659\u0631\u06cc";
        final Object[][] data = new Object[][] {
            { "MonthNames", metaValue_MonthNames },
            { "standalone.MonthNames", metaValue_MonthNames },
            { "MonthAbbreviations", metaValue_MonthNames },
            { "standalone.MonthAbbreviations", metaValue_MonthNames },
            { "standalone.QuarterNames",
                new String[] {
                    "\u0686\u0627\u0631\u0623\u06a9 \u0623\u06a4\u0623\u0644",
                    "\u0686\u0627\u0631\u0623\u06a9 \u062f\u0648\u06cc\u0648\u0645",
                    "\u0686\u0627\u0631\u0623\u06a9 \u0633\u0626\u06cc\u0648\u0645",
                    "\u0686\u0627\u0631\u0623\u06a9 \u0686\u0627\u0631\u0648\u0645",
                }
            },
            { "standalone.QuarterAbbreviations",
                new String[] {
                    "\u0686\u0627\u0631\u0623\u06a9 \u0623\u06a4\u0623\u0644",
                    "Q2",
                    "Q3",
                    "Q4",
                }
            },
            { "field.era", "\u0633\u0623\u0631\u06c9" },
            { "field.year", "\u0633\u0627\u0644" },
            { "field.month", "\u0645\u0627" },
            { "field.week", "\u06be\u0623\u0641\u062a\u06d5" },
            { "field.weekday", "\u0631\u0648\u0659\u0632 \u06be\u0623\u0641\u062a\u06d5" },
            { "field.dayperiod", "\u06af\u0627\u062a \u0631\u0648\u0659\u0632" },
            { "field.hour", "\u0633\u0627\u0623\u062a" },
            { "field.minute", "\u062f\u0626\u06cc\u0642\u06d5" },
            { "field.second", "\u062b\u0627\u0646\u06cc\u06d5" },
            { "field.zone", "\u0631\u0627\u0633\u0627\u06af\u06d5" },
            { "buddhist.MonthNames", metaValue_MonthNames },
            { "japanese.MonthNames", metaValue_MonthNames },
            { "japanese.MonthAbbreviations", metaValue_MonthNames },
            { "roc.MonthNames", metaValue_MonthNames },
            { "roc.MonthAbbreviations", metaValue_MonthNames },
            { "calendarname.gregorian", metaValue_calendarname_gregorian },
            { "calendarname.gregory", metaValue_calendarname_gregorian },
            { "DefaultNumberingSystem", "arabext" },
            { "arabext.NumberElements",
                new String[] {
                    "\u066b",
                    "\u066c",
                    "\u061b",
                    "\u066a",
                    "\u06f0",
                    "#",
                    "\u200e-\u200e",
                    "\u00d7\u06f1\u06f0^",
                    "\u0609",
                    "\u221e",
                    "NaN",
                    "",
                    "",
                }
            },
            { "arabext.NumberPatterns",
                new String[] {
                    "#,##0.###",
                    "\u00a4\u00a0#,##0.00",
                    "#,##0%",
                    "",
                }
            },
        };
        return data;
    }
}
