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

public class FormatData_en_AU extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] metaValue_generic_AmPmMarkers = new String[] {
            "am",
            "pm",
            "midnight",
            "midday",
            "in the morning",
            "",
            "in the afternoon",
            "",
            "in the evening",
            "",
            "at night",
            "",
        };
        final String metaValue_generic_DateFormatItem_GyMMMEEEEd = "EEEE d MMM y G";
        final String[] metaValue_MonthAbbreviations = new String[] {
            "Jan",
            "Feb",
            "Mar",
            "Apr",
            "May",
            "June",
            "July",
            "Aug",
            "Sept",
            "Oct",
            "Nov",
            "Dec",
            "",
        };
        final String[] metaValue_generic_DayNarrows = new String[] {
            "Su.",
            "M.",
            "Tu.",
            "W.",
            "Th.",
            "F.",
            "Sa.",
        };
        final String[] metaValue_generic_TimePatterns = new String[] {
            "h:mm:ss\u202fa zzzz",
            "h:mm:ss\u202fa z",
            "h:mm:ss\u202fa",
            "h:mm\u202fa",
        };
        final String[] metaValue_generic_DateTimePatterns = new String[] {
            "{1}, {0}",
            "{1}, {0}",
            "{1}, {0}",
            "{1}, {0}",
        };
        final String metaValue_generic_DateFormatItem_MMMMEEEEd = "EEEE d MMMM";
        final String metaValue_generic_DateFormatItem_yMMMMEEEEd = "EEEE d MMMM y";
        final String metaValue_generic_DateFormatItem_yMMMEEEEd = "EEEE d MMM y";
        final String metaValue_generic_DateFormatItem_Md = "d/M";
        final String metaValue_generic_DateFormatItem_MEd = "E, d/M";
        final String metaValue_generic_DateFormatItem_MMMEEEEd = "EEEE d MMM";
        final String[] metaValue_java_time_generic_DatePatterns = new String[] {
            "EEEE d MMMM y G",
            "d MMMM y G",
            "d MMM y G",
            "dd/MM/y GGGGG",
        };
        final String[] metaValue_generic_DatePatterns = new String[] {
            "EEEE d MMMM y GGGG",
            "d MMMM y GGGG",
            "d MMM y GGGG",
            "dd/MM/y G",
        };
        final Object[][] data = new Object[][] {
            { "generic.DayNarrows", metaValue_generic_DayNarrows },
            { "generic.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "generic.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "generic.abbreviated.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "generic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.generic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "generic.DatePatterns", metaValue_generic_DatePatterns },
            { "generic.DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "generic.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "generic.DateFormatItem.yMEd", "E, dd/MM/y" },
            { "generic.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "generic.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "generic.DateFormatItem.yyyyMMMMEEEEd", "EEEE d MMMM y G" },
            { "generic.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "generic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "generic.DateFormatItem.yyyyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "generic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "generic.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "generic.DateFormatItem.yMd", "dd/MM/y" },
            { "MonthAbbreviations", metaValue_MonthAbbreviations },
            { "standalone.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "DayNarrows", metaValue_generic_DayNarrows },
            { "AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "abbreviated.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "EEEE, d MMMM y",
                    "d MMMM y",
                    "d MMM y",
                    "d/M/yy",
                }
            },
            { "DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "buddhist.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "buddhist.DayNarrows", metaValue_generic_DayNarrows },
            { "buddhist.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "buddhist.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "buddhist.abbreviated.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "buddhist.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.buddhist.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "buddhist.DatePatterns", metaValue_generic_DatePatterns },
            { "buddhist.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "buddhist.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "buddhist.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "buddhist.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "buddhist.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "buddhist.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "buddhist.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "japanese.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "japanese.DayNarrows", metaValue_generic_DayNarrows },
            { "japanese.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "japanese.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "japanese.abbreviated.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "japanese.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.japanese.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "japanese.DatePatterns",
                new String[] {
                    "EEEE d MMMM yyyy GGGG",
                    "d MMMM yyyy GGGG",
                    "d MMM y GGGG",
                    "dd/MM/y G",
                }
            },
            { "japanese.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "japanese.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "japanese.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "japanese.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "japanese.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "japanese.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "japanese.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "roc.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "roc.DayNarrows", metaValue_generic_DayNarrows },
            { "roc.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "roc.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "roc.abbreviated.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.roc.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "roc.DatePatterns", metaValue_generic_DatePatterns },
            { "roc.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "roc.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "roc.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "roc.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "roc.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "roc.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "roc.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "islamic.DayNarrows", metaValue_generic_DayNarrows },
            { "islamic.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "islamic.narrow.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "islamic.abbreviated.AmPmMarkers", metaValue_generic_AmPmMarkers },
            { "islamic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.islamic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "islamic.DatePatterns", metaValue_generic_DatePatterns },
            { "islamic.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "islamic.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "islamic.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "islamic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic.DateFormatItem.Ed", "E d" },
            { "islamic.DateFormatItem.M", "LL" },
            { "islamic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic.DateFormatItem.MMMEd", "E, d MMM" },
            { "islamic.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "islamic.DateFormatItem.yyyyM", "MM/y GGGGG" },
            { "islamic-civil.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "islamic-civil.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "islamic-civil.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic-civil.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "islamic-civil.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-civil.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-civil.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "islamic-umalqura.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-umalqura.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-umalqura.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "latn.NumberElements",
                new String[] {
                    ".",
                    ",",
                    ";",
                    "%",
                    "0",
                    "#",
                    "-",
                    "e",
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
                    "\u00a4#,##0.00",
                    "#,##0%",
                    "\u00a4#,##0.00;(\u00a4#,##0.00)",
                }
            },
        };
        return data;
    }
}
