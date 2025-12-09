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

public class FormatData_ps extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] metaValue_generic_DayNames = new String[] {
            "\u064a\u0648\u0646\u06cd",
            "\u062f\u0648\u0646\u06cd",
            "\u062f\u0631\u06d0\u0646\u06cd",
            "\u0685\u0644\u0631\u0646\u06cd",
            "\u067e\u064a\u0646\u0681\u0646\u06cd",
            "\u062c\u0645\u0639\u0647",
            "\u0627\u0648\u0646\u06cd",
        };
        final String[] metaValue_generic_AmPmMarkers = new String[] {
            "\u063a.\u0645.",
            "\u063a.\u0648.",
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
        final String[] metaValue_MonthNames = new String[] {
            "\u062c\u0646\u0648\u0631\u064a",
            "\u0641\u0628\u0631\u0648\u0631\u064a",
            "\u0645\u0627\u0631\u0686",
            "\u0627\u067e\u0631\u06cc\u0644",
            "\u0645\u06cd",
            "\u062c\u0648\u0646",
            "\u062c\u0648\u0644\u0627\u06cc",
            "\u0627\u06ab\u0633\u062a",
            "\u0633\u06d0\u067e\u062a\u0645\u0628\u0631",
            "\u0627\u06a9\u062a\u0648\u0628\u0631",
            "\u0646\u0648\u0645\u0628\u0631",
            "\u062f\u0633\u0645\u0628\u0631",
            "",
        };
        final String[] metaValue_generic_QuarterNames = new String[] {
            "\u0644\u0648\u0645\u0693\u06cd \u0631\u0628\u0639\u0647",
            "\u06f2\u0645\u0647 \u0631\u0628\u0639\u0647",
            "\u06f3\u0645\u0647 \u0631\u0628\u0639\u0647",
            "\u06f4\u0645\u0647 \u0631\u0628\u0639\u0647",
        };
        final String[] metaValue_java_time_long_Eras = new String[] {
            "\u0644\u0647 \u0645\u06cc\u0644\u0627\u062f \u0685\u062e\u0647 \u0648\u0693\u0627\u0646\u062f\u06d0",
            "\u0644\u0647 \u0645\u06cc\u0644\u0627\u062f \u0685\u062e\u0647 \u0648\u0631\u0648\u0633\u062a\u0647",
        };
        final String[] metaValue_java_time_Eras = new String[] {
            "\u0644\u0647 \u0645\u06cc\u0644\u0627\u062f \u0648\u0693\u0627\u0646\u062f\u06d0",
            "\u0645.",
        };
        final String[] metaValue_generic_TimePatterns = new String[] {
            "H:mm:ss (zzzz)",
            "H:mm:ss (z)",
            "H:mm:ss",
            "H:mm",
        };
        final String[] metaValue_generic_DateTimePatterns = new String[] {
            "{1} {0}",
            "{1} {0}",
            "{1} {0}",
            "{1} {0}",
        };
        final String metaValue_generic_DateFormatItem_yw = "\u0627\u0648\u0646\u06cd w \u062f Y";
        final String metaValue_generic_DateFormatItem_MMMEd = "E, MMM d";
        final String metaValue_generic_DateFormatItem_GyMd = "M/d/y GGGGG";
        final String metaValue_generic_DateFormatItem_MMMMW = "\u0627\u0648\u0646\u06cd W \u062f MMMM";
        final String[] metaValue_MonthNarrows = new String[] {
            "\u062c",
            "\u0641",
            "\u0645",
            "\u0627",
            "\u0645",
            "\u062c",
            "\u062c",
            "\u0627",
            "\u0633",
            "\u0627",
            "\u0646",
            "\u062f",
            "",
        };
        final String[] metaValue_java_time_generic_DatePatterns = new String[] {
            "EEEE \u062f G y \u062f MMMM d",
            "\u062f G y \u062f MMMM d",
            "G y MMM d",
            "GGGGG y/M/d",
        };
        final String[] metaValue_generic_DatePatterns = new String[] {
            "EEEE \u062f GGGG y \u062f MMMM d",
            "\u062f GGGG y \u062f MMMM d",
            "GGGG y MMM d",
            "G y/M/d",
        };
        final String metaValue_calendarname_gregorian = "\u06ab\u0631\u064a\u06ab\u0648\u0631\u064a\u0646 \u062c\u0646\u062a\u0631\u064a";
        final String[] metaValue_ListPatterns_standard = new String[] {
            "{0}\u060c {1}",
            "{0}\u060c {1}",
            "{0}\u060c \u0627\u0648 {1}",
            "{0} \u0627\u0648 {1}",
            "",
        };
        final String[] metaValue_ListPatterns_or = new String[] {
            "",
            "",
            "{0}, \u06cc\u0627 {1}",
            "",
            "",
        };
        final Object[][] data = new Object[][] {
            { "generic.DayNames", metaValue_generic_DayNames },
            { "generic.DayAbbreviations", metaValue_generic_DayNames },
            { "generic.QuarterNames", metaValue_generic_QuarterNames },
            { "generic.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "generic.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "generic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.generic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "generic.DatePatterns", metaValue_generic_DatePatterns },
            { "generic.DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "generic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "generic.DateFormatItem.Hm", "H:mm" },
            { "generic.DateFormatItem.MMMMd", "d MMMM" },
            { "generic.DateFormatItem.H", "H" },
            { "generic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "generic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "generic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "generic.DateFormatItem.Md", "M/d" },
            { "generic.DateFormatItem.Hms", "H:mm:ss" },
            { "generic.DateFormatItem.yM", "G y/M" },
            { "generic.DateFormatItem.yMMMM", "\u062f G y \u062f MMMM" },
            { "MonthNames", metaValue_MonthNames },
            { "standalone.MonthNames",
                new String[] {
                    "\u062c\u0646\u0648\u0631\u064a",
                    "\u0641\u06d0\u0628\u0631\u0648\u0631\u064a",
                    "\u0645\u0627\u0631\u0686",
                    "\u0627\u067e\u0631\u06cc\u0644",
                    "\u0645\u06cd",
                    "\u062c\u0648\u0646",
                    "\u062c\u0648\u0644\u0627\u06cc",
                    "\u0627\u06ab\u0633\u062a",
                    "\u0633\u067e\u062a\u0645\u0628\u0631",
                    "\u0627\u06a9\u062a\u0648\u0628\u0631",
                    "\u0646\u0648\u0645\u0628\u0631",
                    "\u062f\u0633\u0645\u0628\u0631",
                    "",
                }
            },
            { "MonthAbbreviations", metaValue_MonthNames },
            { "standalone.MonthAbbreviations",
                new String[] {
                    "\u062c\u0646\u0648\u0631\u064a",
                    "\u0641\u0628\u0631\u0648\u0631\u064a",
                    "\u0645\u0627\u0631\u0686",
                    "\u0627\u067e\u0631\u06cc\u0644",
                    "\u0645\u06cd",
                    "\u062c\u0648\u0646",
                    "\u062c\u0648\u0644\u0627\u06cc",
                    "\u0627\u06ab\u0633\u062a",
                    "\u0633\u067e\u062a\u0645\u0628\u0631",
                    "\u0627\u06a9\u062a\u0648\u0628\u0631",
                    "\u0646\u0648\u0645\u0628\u0631",
                    "\u062f\u0633\u0645\u0628\u0631",
                    "",
                }
            },
            { "MonthNarrows", metaValue_MonthNarrows },
            { "DayNames", metaValue_generic_DayNames },
            { "standalone.DayNames", metaValue_generic_DayNames },
            { "DayAbbreviations", metaValue_generic_DayNames },
            { "standalone.DayAbbreviations", metaValue_generic_DayNames },
            { "QuarterNames", metaValue_generic_QuarterNames },
            { "standalone.QuarterNames", metaValue_generic_QuarterNames },
            { "QuarterAbbreviations", metaValue_generic_QuarterNames },
            { "AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "java.time.long.Eras", metaValue_java_time_long_Eras },
            { "long.Eras", metaValue_java_time_long_Eras },
            { "java.time.Eras", metaValue_java_time_Eras },
            { "Eras", metaValue_java_time_Eras },
            { "java.time.narrow.Eras", metaValue_java_time_Eras },
            { "narrow.Eras", metaValue_java_time_Eras },
            { "field.era", "\u062f\u0648\u0631" },
            { "field.year", "\u06a9\u0627\u0644" },
            { "field.month", "\u0645\u064a\u0627\u0634\u062a" },
            { "field.week", "\u0627\u0648\u0646\u06cd" },
            { "field.weekday", "\u062f \u0627\u0648\u0646\u06cd \u0648\u0631\u0681" },
            { "field.dayperiod", "\u0648\u0631\u0681 \u0634\u06d0\u0628\u0647" },
            { "field.hour", "\u0633\u0627\u0639\u062a" },
            { "timezone.regionFormat", "\u062f {0} \u067e\u0647 \u0648\u062e\u062a" },
            { "timezone.regionFormat.daylight", "{0} \u0631\u06bc\u0627 \u0648\u0631\u0681 \u0648\u062e\u062a" },
            { "timezone.regionFormat.standard", "{0} \u0645\u0639\u06cc\u0627\u0631\u06cc \u0648\u062e\u062a" },
            { "field.minute", "\u062f\u0642\u064a\u0642\u0647" },
            { "field.second", "\u062b\u0627\u0646\u064a\u0647" },
            { "field.zone", "\u0648\u062e\u062a \u0633\u064a\u0645\u0647" },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "EEEE \u062f y \u062f MMMM d",
                    "y MMMM d",
                    "y MMM d",
                    "y/M/d",
                }
            },
            { "DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "PluralRules", "one:n = 1" },
            { "DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "buddhist.MonthNames", metaValue_MonthNames },
            { "buddhist.MonthNarrows", metaValue_MonthNarrows },
            { "buddhist.DayNames", metaValue_generic_DayNames },
            { "buddhist.DayAbbreviations", metaValue_generic_DayNames },
            { "buddhist.QuarterNames", metaValue_generic_QuarterNames },
            { "buddhist.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "buddhist.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "buddhist.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.buddhist.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "buddhist.DatePatterns", metaValue_generic_DatePatterns },
            { "buddhist.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "buddhist.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "buddhist.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "buddhist.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "japanese.MonthNames", metaValue_MonthNames },
            { "japanese.MonthAbbreviations", metaValue_MonthNames },
            { "japanese.MonthNarrows", metaValue_MonthNarrows },
            { "japanese.DayNames", metaValue_generic_DayNames },
            { "japanese.DayAbbreviations", metaValue_generic_DayNames },
            { "japanese.QuarterNames", metaValue_generic_QuarterNames },
            { "japanese.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "japanese.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "japanese.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.japanese.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "japanese.DatePatterns",
                new String[] {
                    "EEEE \u062f GGGG yyyy \u062f MMMM d",
                    "\u062f GGGG yyyy \u062f MMMM d",
                    "GGGG y MMM d",
                    "G y/M/d",
                }
            },
            { "japanese.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "japanese.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "japanese.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "japanese.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "roc.MonthNames", metaValue_MonthNames },
            { "roc.MonthAbbreviations", metaValue_MonthNames },
            { "roc.MonthNarrows", metaValue_MonthNarrows },
            { "roc.DayNames", metaValue_generic_DayNames },
            { "roc.QuarterNames", metaValue_generic_QuarterNames },
            { "roc.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "roc.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.roc.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "roc.DatePatterns", metaValue_generic_DatePatterns },
            { "roc.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "roc.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "roc.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "roc.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic.MonthNames",
                new String[] {
                    "\u0645\u062d\u0631\u0645",
                    "\u0635\u0641\u0631",
                    "\u0631\u0628\u064a\u0639",
                    "\u0631\u0628\u064a\u0639 II",
                    "\u062c\u0645\u0627\u0639\u0647",
                    "\u062c\u0645\u0627\u0639\u0647 II",
                    "\u0631\u062c\u0628",
                    "\u0634\u0639\u0628\u0627\u0646",
                    "\u0631\u0645\u0636\u0627\u0646",
                    "\u0634\u0648\u0627\u0644",
                    "\u0630\u064a \u0627\u0644\u0642\u0639\u062f\u0647",
                    "\u0630\u064a \u0627\u0644\u062d\u062c",
                    "",
                }
            },
            { "islamic.MonthAbbreviations",
                new String[] {
                    "\u0645\u062d\u0631\u0645",
                    "\u0635\u0641\u0631",
                    "\u0631\u0628\u064a\u0639",
                    "\u0631\u0628\u064a\u0639 II",
                    "\u062c\u0645\u0627\u062f",
                    "\u062c\u0645\u0627\u062f \u06f2",
                    "\u0631\u062c\u0628",
                    "\u0634\u0639\u0628\u0627\u0646",
                    "\u0631\u0645\u0636\u0627\u0646",
                    "\u0634\u0648\u0627\u0644",
                    "\u062f\u0627\u0644\u0642\u0627\u0639\u062f\u0647",
                    "\u0630\u064a \u0627\u0644\u062d\u062c",
                    "",
                }
            },
            { "islamic.DayNames", metaValue_generic_DayNames },
            { "islamic.DayAbbreviations", metaValue_generic_DayNames },
            { "islamic.QuarterNames", metaValue_generic_QuarterNames },
            { "islamic.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "islamic.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "islamic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.islamic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "islamic.DatePatterns", metaValue_generic_DatePatterns },
            { "islamic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-civil.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-civil.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-civil.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-civil.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-umalqura.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-umalqura.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-umalqura.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-umalqura.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "calendarname.islamic-civil", "\u0627\u0633\u0644\u0627\u0645\u064a \u062c\u0646\u062a\u0631\u064a (\u062c\u062f\u0648\u0644\u064a\u060c \u0645\u062f\u0646\u064a \u0639\u0635\u0631)" },
            { "calendarname.islamic", "\u0627\u0633\u0644\u0627\u0645\u064a \u062c\u0646\u062a\u0631\u064a" },
            { "calendarname.buddhist", "\u0628\u0648\u062f\u0627\u064a\u064a \u062c\u0646\u062a\u0631\u064a" },
            { "calendarname.japanese", "\u062c\u0627\u067e\u0627\u0646\u064a \u062c\u0646\u062a\u0631\u064a" },
            { "calendarname.roc", "\u0645\u0646\u06af\u0648\u0648 \u062c\u0646\u062a\u0631\u064a" },
            { "calendarname.islamic-umalqura", "\u0627\u0633\u0644\u0627\u0645\u064a \u06a9\u0644\u06cc\u0632\u0647 (\u0627\u0645 \u0627\u0644\u0642\u0631\u06cc)" },
            { "calendarname.gregorian", metaValue_calendarname_gregorian },
            { "calendarname.gregory", metaValue_calendarname_gregorian },
            { "DefaultNumberingSystem", "arabext" },
            { "arab.NumberElements",
                new String[] {
                    "\u066b",
                    "\u066c",
                    "\u061b",
                    "\u066a\u061c",
                    "\u0660",
                    "#",
                    "\u061c-",
                    "^",
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
            { "latn.NumberElements",
                new String[] {
                    ",",
                    ".",
                    ";",
                    "%",
                    "0",
                    "#",
                    "\u200e\u2212",
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
                    "{other:0K}",
                    "{other:00K}",
                    "{other:000K}",
                    "{other:0M}",
                    "{other:00M}",
                    "{other:000M}",
                    "{one:0B other:0B}",
                    "{one:00B other:00B}",
                    "{one:000G other:000B}",
                    "{other:0T}",
                    "{other:00T}",
                    "{other:000T}",
                }
            },
            { "long.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "",
                    "",
                    "",
                    "",
                    "",
                    "",
                    "{one:0G other:0G}",
                    "{one:00G other:00G}",
                    "{one:000G other:000G}",
                }
            },
            { "ListPatterns_standard", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-short", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-narrow",
                new String[] {
                    "{0}\u060c {1}",
                    "{0}\u060c {1}",
                    "{0}\u060c \u0627\u0648 {1}",
                    "{0}\u060c {1}",
                    "",
                }
            },
            { "ListPatterns_or", metaValue_ListPatterns_or },
            { "ListPatterns_or-short", metaValue_ListPatterns_or },
            { "ListPatterns_or-narrow", metaValue_ListPatterns_or },
            { "ListPatterns_unit", metaValue_ListPatterns_standard },
            { "ListPatterns_unit-short",
                new String[] {
                    "{0}, {1}",
                    "{0}, {1}",
                    "{0}, {1}",
                    "{0} \u0648 {1}",
                    "",
                }
            },
            { "ListPatterns_unit-narrow",
                new String[] {
                    "{0} {1}",
                    "{0} {1}",
                    "{0} {1}",
                    "{0} {1}",
                    "",
                }
            },
        };
        return data;
    }
}
