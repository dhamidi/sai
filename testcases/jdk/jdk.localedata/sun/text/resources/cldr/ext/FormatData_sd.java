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

public class FormatData_sd extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] metaValue_generic_DayNames = new String[] {
            "\u0622\u0686\u0631",
            "\u0633\u0648\u0645\u0631",
            "\u0627\u06b1\u0627\u0631\u0648",
            "\u0627\u0631\u0628\u0639",
            "\u062e\u0645\u064a\u0633",
            "\u062c\u0645\u0639\u0648",
            "\u0687\u0646\u0687\u0631",
        };
        final String[] metaValue_MonthNames = new String[] {
            "\u062c\u0646\u0648\u0631\u064a",
            "\u0641\u064a\u0628\u0631\u0648\u0631\u064a",
            "\u0645\u0627\u0631\u0686",
            "\u0627\u067e\u0631\u064a\u0644",
            "\u0645\u0626\u064a",
            "\u062c\u0648\u0646",
            "\u062c\u0648\u0644\u0627\u0621\u0650",
            "\u0622\u06af\u0633\u067d",
            "\u0633\u064a\u067e\u067d\u0645\u0628\u0631",
            "\u0622\u06aa\u067d\u0648\u0628\u0631",
            "\u0646\u0648\u0645\u0628\u0631",
            "\u068a\u0633\u0645\u0628\u0631",
            "",
        };
        final String[] metaValue_MonthNarrows = new String[] {
            "J",
            "F",
            "M",
            "A",
            "M",
            "J",
            "J",
            "A",
            "S",
            "O",
            "N",
            "D",
            "",
        };
        final String[] metaValue_generic_DayNarrows = new String[] {
            "\u0622\u0686\u0631",
            "\u0633\u0648",
            "\u0627\u06b1\u0627\u0631\u0648",
            "\u0627\u0631\u0628\u0639",
            "\u062e\u0645",
            "\u062c\u0645\u0639\u0648",
            "\u0687\u0646\u0687\u0631",
        };
        final String[] metaValue_generic_QuarterNames = new String[] {
            "\u067e\u0647\u0631\u064a\u0646 \u067d\u064a \u0645\u0627\u0647\u064a",
            "\u067b\u064a\u0646 \u067d\u064a \u0645\u0627\u0647\u064a",
            "\u067d\u064a\u0646 \u067d\u064a \u0645\u0627\u0647\u064a",
            "\u0686\u0648\u067f\u064a\u0646 \u067d\u064a \u0645\u0627\u0647\u064a",
        };
        final String[] metaValue_generic_QuarterAbbreviations = new String[] {
            "Q1",
            "Q2",
            "Q3",
            "Q5",
        };
        final String[] metaValue_generic_AmPmMarkers = new String[] {
            "\u0635\u0628\u062d\u060c \u0645\u0646\u062c\u0647\u0646\u062f",
            "\u0645\u0646\u062c\u0647\u0646\u062f\u060c \u0634\u0627\u0645",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
        };
        final String[] metaValue_generic_narrow_AmPmMarkers = new String[] {
            "\u0635\u0628\u062d\u060c\u202f\u0645\u0646\u062c\u0647\u0646\u062f",
            "\u0645\u0646\u062c\u0647\u0646\u062f\u060c\u202f\u0634\u0627\u0645",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
        };
        final String[] metaValue_generic_abbreviated_AmPmMarkers = new String[] {
            "\u0635\u0628\u062d\u060c\u202f\u0645\u0646\u062c\u0647\u0646\u062f",
            "\u0634\u0627\u0645\u060c\u202f\u0645\u0646\u062c\u0647\u0646\u062f",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
        };
        final String[] metaValue_java_time_long_Eras = new String[] {
            "\u0642\u0628\u0644 \u0645\u0633\u064a\u062d",
            "\u0639\u064a\u0633\u0648\u064a \u06a9\u0627\u0646 \u067e\u0647\u0631\u064a\u0646",
        };
        final String[] metaValue_java_time_Eras = new String[] {
            "BC",
            "CD",
        };
        final String[] metaValue_generic_TimePatterns = new String[] {
            "h:mm:ss a zzzz",
            "h:mm:ss a z",
            "h:mm:ss a",
            "h:mm a",
        };
        final String metaValue_generic_DateFormatItem_yw = "\u06be\u0641\u062a\u0648 w \u062c\u0648 Y";
        final String metaValue_generic_DateFormatItem_GyMd = "M/d/y G";
        final String metaValue_generic_DateFormatItem_MMMMW = "\u06be\u0641\u062a\u0648 W \u062c\u0648 MMMM";
        final String[] metaValue_islamic_MonthNames = new String[] {
            "\u0645\u062d\u0631\u0645",
            "\u0635\u0641\u0631",
            "\u0631\u0628\u064a\u0639 \u0627\u0644\u0627\u0648\u0644",
            "\u0631\u0628\u06cc\u0639 \u0627\u0644\u0627\u062e\u0631",
            "\u062c\u0645\u0627\u062f\u06cc \u0627\u0644\u0627\u0648\u0644",
            "\u062c\u0645\u0627\u062f\u064a \u0627\u0644\u0627\u062e\u0631",
            "\u0631\u062c\u0628",
            "\u0634\u0639\u0628\u0627\u0646",
            "\u0631\u0645\u0636\u0627\u0646",
            "\u0634\u0648\u0627\u0644",
            "\u0630\u0648\u0627\u0644\u0642\u0639\u062f",
            "\u0630\u0648\u0627\u0644\u062d\u062c\u06c1",
            "",
        };
        final String metaValue_calendarname_gregorian = "\u062c\u0627\u0631\u062c\u064a\u0627\u0626\u064a \u06aa\u0626\u0644\u064a\u0646\u068a\u0631";
        final String[] metaValue_arabext_NumberPatterns = new String[] {
            "#,##0.###",
            "\u00a4\u00a0#,##0.00",
            "#,##0%",
            "",
        };
        final String[] metaValue_ListPatterns_standard = new String[] {
            "{0}\u060c {1}",
            "{0}\u060c {1}",
            "{0}\u060c \u06fd {1}",
            "{0} \u06fd {1}",
            "",
        };
        final String[] metaValue_ListPatterns_or = new String[] {
            "",
            "",
            "{0}, \u064a\u0627 {1}",
            "{0} \u064a\u0627 {1}",
            "",
        };
        final String[] metaValue_ListPatterns_standard_narrow = new String[] {
            "{0}\u060c {1}",
            "{0}\u060c {1}",
            "{0}, {1}",
            "{0}\u060c {1}",
            "",
        };
        final Object[][] data = new Object[][] {
            { "generic.DayNames", metaValue_generic_DayNames },
            { "generic.DayAbbreviations", metaValue_generic_DayNames },
            { "generic.DayNarrows", metaValue_generic_DayNarrows },
            { "generic.QuarterNames", metaValue_generic_QuarterNames },
            { "generic.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "generic.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "generic.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "generic.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "generic.TimePatterns", metaValue_generic_TimePatterns },
            { "generic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "generic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "generic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "MonthNames", metaValue_MonthNames },
            { "standalone.MonthNames", metaValue_MonthNames },
            { "MonthAbbreviations", metaValue_MonthNames },
            { "standalone.MonthAbbreviations", metaValue_MonthNames },
            { "MonthNarrows", metaValue_MonthNarrows },
            { "standalone.MonthNarrows", metaValue_MonthNarrows },
            { "DayNames", metaValue_generic_DayNames },
            { "standalone.DayNames", metaValue_generic_DayNames },
            { "DayAbbreviations", metaValue_generic_DayNames },
            { "standalone.DayAbbreviations", metaValue_generic_DayNames },
            { "DayNarrows", metaValue_generic_DayNarrows },
            { "standalone.DayNarrows", metaValue_generic_DayNarrows },
            { "QuarterNames", metaValue_generic_QuarterNames },
            { "standalone.QuarterNames", metaValue_generic_QuarterNames },
            { "QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "standalone.QuarterAbbreviations",
                new String[] {
                    "Q1",
                    "Q2",
                    "Q3",
                    "Q4",
                }
            },
            { "AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "java.time.long.Eras", metaValue_java_time_long_Eras },
            { "long.Eras", metaValue_java_time_long_Eras },
            { "java.time.Eras", metaValue_java_time_Eras },
            { "Eras", metaValue_java_time_Eras },
            { "java.time.narrow.Eras", metaValue_java_time_Eras },
            { "narrow.Eras", metaValue_java_time_Eras },
            { "field.era", "\u062f\u0648\u0631" },
            { "field.year", "\u0633\u0627\u0644" },
            { "field.month", "\u0645\u0647\u064a\u0646\u0648" },
            { "field.week", "\u0647\u0641\u062a\u0648" },
            { "field.weekday", "\u0647\u0641\u062a\u064a \u062c\u0648 \u068f\u064a\u0646\u0647\u0646" },
            { "field.dayperiod", "\u0635\u0628\u062d/\u0634\u0627\u0645" },
            { "field.hour", "\u06aa\u0644\u0627\u06aa" },
            { "timezone.regionFormat", "{0} \u0648\u0642\u062a" },
            { "timezone.regionFormat.daylight", "{0} \u068f\u064a\u0646\u0647\u0646 \u062c\u0648 \u0648\u0642\u062a" },
            { "timezone.regionFormat.standard", "{0} \u0645\u0639\u064a\u0627\u0631\u064a \u0648\u0642\u062a" },
            { "field.minute", "\u0645\u0646\u067d" },
            { "field.second", "\u0633\u064a\u06aa\u0646\u068a" },
            { "field.zone", "\u067d\u0627\u0626\u064a\u0645 \u0632\u0648\u0646" },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "EEEE, MMMM d, y",
                    "y MMMM d",
                    "y MMM d",
                    "y-MM-dd",
                }
            },
            { "DateTimePatterns",
                new String[] {
                    "{1} {0}",
                    "{1} {0}",
                    "{1} {0}",
                    "{1} {0}",
                }
            },
            { "PluralRules", "one:n = 1" },
            { "DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "buddhist.MonthNames", metaValue_MonthNames },
            { "buddhist.DayNames", metaValue_generic_DayNames },
            { "buddhist.DayAbbreviations", metaValue_generic_DayNames },
            { "buddhist.QuarterNames", metaValue_generic_QuarterNames },
            { "buddhist.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "buddhist.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "buddhist.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "buddhist.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "buddhist.TimePatterns", metaValue_generic_TimePatterns },
            { "buddhist.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "buddhist.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "buddhist.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "japanese.MonthNames", metaValue_MonthNames },
            { "japanese.MonthAbbreviations", metaValue_MonthNames },
            { "japanese.DayNames", metaValue_generic_DayNames },
            { "japanese.DayAbbreviations", metaValue_generic_DayNames },
            { "japanese.QuarterNames", metaValue_generic_QuarterNames },
            { "japanese.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "japanese.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "japanese.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "japanese.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "japanese.TimePatterns", metaValue_generic_TimePatterns },
            { "japanese.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "japanese.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "japanese.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "roc.MonthNames", metaValue_MonthNames },
            { "roc.MonthAbbreviations", metaValue_MonthNames },
            { "roc.DayNames", metaValue_generic_DayNames },
            { "roc.DayNarrows", metaValue_generic_DayNarrows },
            { "roc.QuarterNames", metaValue_generic_QuarterNames },
            { "roc.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "roc.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "roc.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "roc.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "roc.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "roc.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "roc.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic.MonthNames", metaValue_islamic_MonthNames },
            { "islamic.MonthAbbreviations", metaValue_islamic_MonthNames },
            { "islamic.DayNames", metaValue_generic_DayNames },
            { "islamic.DayAbbreviations", metaValue_generic_DayNames },
            { "islamic.QuarterNames", metaValue_generic_QuarterNames },
            { "islamic.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "islamic.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "islamic.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "islamic.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "islamic.TimePatterns", metaValue_generic_TimePatterns },
            { "islamic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-civil.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-civil.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-civil.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-umalqura.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-umalqura.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-umalqura.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "calendarname.islamic-civil", "\u0647\u062c\u0631\u064a \u06aa\u0626\u0644\u064a\u0646\u068a\u0631 (\u067d\u064a\u0628\u0644 \u0648\u0627\u0631\u0648\u060c \u0634\u0647\u0631\u064a \u062f\u0648\u0631)" },
            { "calendarname.islamic", "\u0647\u062c\u0631\u064a \u06aa\u0626\u0644\u064a\u0646\u068a\u0631" },
            { "calendarname.buddhist", "\u067b\u068c \u068c\u0631\u0645 \u062c\u0648 \u06aa\u0626\u0644\u064a\u0646\u068a\u0631" },
            { "calendarname.japanese", "\u062c\u0627\u067e\u0627\u0646\u064a \u06aa\u0626\u0644\u064a\u0646\u068a\u0631" },
            { "calendarname.roc", "\u0645\u0646\u06af\u0648\u0648 \u06aa\u0626\u0644\u064a\u0646\u068a\u0631" },
            { "calendarname.islamic-umalqura", "\u0647\u062c\u0631\u064a \u06aa\u0626\u0644\u064a\u0646\u068a\u0631 (\u0627\u064f\u0645 \u0627\u0644\u0642\u0631\u0627)" },
            { "calendarname.gregorian", metaValue_calendarname_gregorian },
            { "calendarname.gregory", metaValue_calendarname_gregorian },
            { "DefaultNumberingSystem", "arab" },
            { "arab.NumberElements",
                new String[] {
                    ".",
                    "\u066c",
                    "\u061b",
                    "\u066a\u061c",
                    "\u0660",
                    "#",
                    "\u061c-",
                    "\u0627\u0633",
                    "\u0609",
                    "\u221e",
                    "NaN",
                    "",
                    "",
                }
            },
            { "arab.NumberPatterns",
                new String[] {
                    "#,##0.###",
                    "#,##0.00\u00a0\u00a4",
                    "#,##0%",
                    "",
                }
            },
            { "arabext.NumberElements",
                new String[] {
                    ".",
                    "\u060c",
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
            { "arabext.NumberPatterns", metaValue_arabext_NumberPatterns },
            { "deva.NumberElements",
                new String[] {
                    ".",
                    "\u060c",
                    ";",
                    "%",
                    "\u0966",
                    "#",
                    "-",
                    "E",
                    "\u2030",
                    "\u221e",
                    "NaN",
                    "",
                    "",
                }
            },
            { "deva.NumberPatterns", metaValue_arabext_NumberPatterns },
            { "latn.NumberElements",
                new String[] {
                    ".",
                    ",",
                    ";",
                    "%",
                    "0",
                    "#",
                    "-",
                    "E",
                    "\u2030",
                    "\u221e",
                    "NaN",
                    "",
                    "",
                }
            },
            { "latn.NumberPatterns",
                new String[] {
                    "#,##0.###",
                    "\u00a4\u00a0#,##0.00",
                    "#,##0%",
                    "\u00a4#,##0.00;(\u00a4#,##0.00)",
                }
            },
            { "short.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:0\u00a0\u0647\u0632\u0627\u0631 other:0\u00a0\u0647\u0632\u0627\u0631}",
                    "{one:00\u00a0\u0647\u0632\u0627\u0631 other:00\u00a0\u0647\u0632\u0627\u0631}",
                    "{one:000\u00a0\u0647\u0632\u0627\u0631 other:000\u00a0\u0647\u0632\u0627\u0631}",
                    "{one:0\u00a0\u0645\u0644\u064a\u0646 other:0\u00a0\u0645\u0644\u064a\u0646}",
                    "{one:00\u00a0\u0645\u0644\u064a\u0646 other:00\u00a0\u0645\u0644\u064a\u0646}",
                    "{one:000\u00a0\u0645\u0644\u064a\u0646 other:000\u00a0\u0645\u0644\u064a\u0646}",
                    "{one:0\u00a0\u0628\u0644\u064a\u0646 other:0\u00a0\u0628\u0644\u064a\u0646}",
                    "{one:00\u00a0\u0628\u0644\u064a\u0646 other:00\u00a0\u0628\u0644\u064a\u0646}",
                    "{one:000\u00a0\u0628\u0644\u064a\u0646 other:000\u00a0\u0628\u0644\u064a\u0646}",
                    "{one:0\u00a0\u067d\u0631\u0644\u064a\u0646 other:0\u00a0\u067d\u0631\u0644\u064a\u0646}",
                    "{one:00\u00a0\u067d\u0631\u0644\u064a\u0646 other:00\u00a0\u067d\u0631\u0644\u064a\u0646}",
                    "{one:000\u00a0\u067d\u0631\u0644\u064a\u0646 other:000\u00a0\u067d\u0631\u0644\u064a\u0646}",
                }
            },
            { "long.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:0' '\u06be\u0632\u0627\u0631 other:0' '\u06be\u0632\u0627\u0631}",
                    "{one:00' '\u06be\u0632\u0627\u0631 other:00' '\u0647\u0632\u0627\u0631}",
                    "{one:000' '\u0647\u0632\u0627\u0631 other:000' '\u0647\u0632\u0627\u0631}",
                    "{one:0' '\u0645\u0644\u064a\u0646 other:0' '\u0645\u0644\u064a\u0646}",
                    "{one:00' '\u0645\u0644\u064a\u0646 other:00' '\u0645\u0644\u064a\u0646}",
                    "{one:000' '\u0645\u0644\u064a\u0646 other:000' '\u0645\u0644\u064a\u0646}",
                    "{one:0' '\u0628\u0644\u064a\u0646 other:0' '\u0628\u0644\u064a\u0646}",
                    "{one:00' '\u0628\u0644\u064a\u0646 other:00' '\u0628\u0644\u064a\u0646}",
                    "{one:000' '\u0628\u0644\u064a\u0646 other:000' '\u0628\u0644\u064a\u0646}",
                    "{one:0' '\u067d\u0631\u0644\u064a\u0646 other:0' '\u067d\u0631\u0644\u064a\u0646}",
                    "{one:00' '\u067d\u0631\u0644\u064a\u0646 other:00' '\u067d\u0631\u0644\u064a\u0646}",
                    "{one:000' '\u067d\u0631\u0644\u064a\u0646 other:000' '\u067d\u0631\u0644\u064a\u0646}",
                }
            },
            { "ListPatterns_standard", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-short", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-narrow", metaValue_ListPatterns_standard_narrow },
            { "ListPatterns_or", metaValue_ListPatterns_or },
            { "ListPatterns_or-short", metaValue_ListPatterns_or },
            { "ListPatterns_or-narrow", metaValue_ListPatterns_or },
            { "ListPatterns_unit", metaValue_ListPatterns_standard },
            { "ListPatterns_unit-short", metaValue_ListPatterns_standard_narrow },
            { "ListPatterns_unit-narrow", metaValue_ListPatterns_standard },
        };
        return data;
    }
}
