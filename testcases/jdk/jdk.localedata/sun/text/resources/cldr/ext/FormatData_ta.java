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

public class FormatData_ta extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String metaValue_generic_DateFormatItem_Gy = "y G";
        final String metaValue_generic_DateFormatItem_yyyyMMM = "MMM y G";
        final String metaValue_generic_DateFormatItem_yyyyMMMd = "d MMM, y G";
        final String metaValue_generic_DateFormatItem_GyMd = "d/M/y G";
        final String metaValue_generic_DateFormatItem_yyyyMMMEd = "E, d MMM, y G";
        final String[] metaValue_MonthNames = new String[] {
            "\u0b9c\u0ba9\u0bb5\u0bb0\u0bbf",
            "\u0baa\u0bbf\u0baa\u0bcd\u0bb0\u0bb5\u0bb0\u0bbf",
            "\u0bae\u0bbe\u0bb0\u0bcd\u0b9a\u0bcd",
            "\u0b8f\u0baa\u0bcd\u0bb0\u0bb2\u0bcd",
            "\u0bae\u0bc7",
            "\u0b9c\u0bc2\u0ba9\u0bcd",
            "\u0b9c\u0bc2\u0bb2\u0bc8",
            "\u0b86\u0b95\u0bb8\u0bcd\u0b9f\u0bcd",
            "\u0b9a\u0bc6\u0baa\u0bcd\u0b9f\u0bae\u0bcd\u0baa\u0bb0\u0bcd",
            "\u0b85\u0b95\u0bcd\u0b9f\u0bcb\u0baa\u0bb0\u0bcd",
            "\u0ba8\u0bb5\u0bae\u0bcd\u0baa\u0bb0\u0bcd",
            "\u0b9f\u0bbf\u0b9a\u0bae\u0bcd\u0baa\u0bb0\u0bcd",
            "",
        };
        final String[] metaValue_MonthAbbreviations = new String[] {
            "\u0b9c\u0ba9.",
            "\u0baa\u0bbf\u0baa\u0bcd.",
            "\u0bae\u0bbe\u0bb0\u0bcd.",
            "\u0b8f\u0baa\u0bcd.",
            "\u0bae\u0bc7",
            "\u0b9c\u0bc2\u0ba9\u0bcd",
            "\u0b9c\u0bc2\u0bb2\u0bc8",
            "\u0b86\u0b95.",
            "\u0b9a\u0bc6\u0baa\u0bcd.",
            "\u0b85\u0b95\u0bcd.",
            "\u0ba8\u0bb5.",
            "\u0b9f\u0bbf\u0b9a.",
            "",
        };
        final String[] metaValue_MonthNarrows = new String[] {
            "\u0b9c",
            "\u0baa\u0bbf",
            "\u0bae\u0bbe",
            "\u0b8f",
            "\u0bae\u0bc7",
            "\u0b9c\u0bc2",
            "\u0b9c\u0bc2",
            "\u0b86",
            "\u0b9a\u0bc6",
            "\u0b85",
            "\u0ba8",
            "\u0b9f\u0bbf",
            "",
        };
        final String[] metaValue_generic_DayNames = new String[] {
            "\u0b9e\u0bbe\u0baf\u0bbf\u0bb1\u0bc1",
            "\u0ba4\u0bbf\u0b99\u0bcd\u0b95\u0bb3\u0bcd",
            "\u0b9a\u0bc6\u0bb5\u0bcd\u0bb5\u0bbe\u0baf\u0bcd",
            "\u0baa\u0bc1\u0ba4\u0ba9\u0bcd",
            "\u0bb5\u0bbf\u0baf\u0bbe\u0bb4\u0ba9\u0bcd",
            "\u0bb5\u0bc6\u0bb3\u0bcd\u0bb3\u0bbf",
            "\u0b9a\u0ba9\u0bbf",
        };
        final String[] metaValue_generic_DayAbbreviations = new String[] {
            "\u0b9e\u0bbe\u0baf\u0bbf.",
            "\u0ba4\u0bbf\u0b99\u0bcd.",
            "\u0b9a\u0bc6\u0bb5\u0bcd.",
            "\u0baa\u0bc1\u0ba4.",
            "\u0bb5\u0bbf\u0baf\u0bbe.",
            "\u0bb5\u0bc6\u0bb3\u0bcd.",
            "\u0b9a\u0ba9\u0bbf",
        };
        final String[] metaValue_generic_DayNarrows = new String[] {
            "\u0b9e\u0bbe",
            "\u0ba4\u0bbf",
            "\u0b9a\u0bc6",
            "\u0baa\u0bc1",
            "\u0bb5\u0bbf",
            "\u0bb5\u0bc6",
            "\u0b9a",
        };
        final String[] metaValue_generic_QuarterNames = new String[] {
            "\u0bae\u0bc1\u0ba4\u0bb2\u0bcd \u0b95\u0bbe\u0bb2\u0bbe\u0ba3\u0bcd\u0b9f\u0bc1",
            "\u0b87\u0bb0\u0ba3\u0bcd\u0b9f\u0bbe\u0bae\u0bcd \u0b95\u0bbe\u0bb2\u0bbe\u0ba3\u0bcd\u0b9f\u0bc1",
            "\u0bae\u0bc2\u0ba9\u0bcd\u0bb1\u0bbe\u0bae\u0bcd \u0b95\u0bbe\u0bb2\u0bbe\u0ba3\u0bcd\u0b9f\u0bc1",
            "\u0ba8\u0bbe\u0ba9\u0bcd\u0b95\u0bbe\u0bae\u0bcd \u0b95\u0bbe\u0bb2\u0bbe\u0ba3\u0bcd\u0b9f\u0bc1",
        };
        final String[] metaValue_generic_QuarterAbbreviations = new String[] {
            "\u0b95\u0bbe.1",
            "\u0b95\u0bbe.2",
            "\u0b95\u0bbe.3",
            "\u0b95\u0bbe.4",
        };
        final String[] metaValue_generic_abbreviated_AmPmMarkers = new String[] {
            "AM",
            "PM",
            "\u0ba8\u0bb3\u0bcd\u0bb3\u0bbf\u0bb0\u0bb5\u0bc1",
            "\u0ba8\u0ba3\u0bcd\u0baa\u0b95\u0bb2\u0bcd",
            "\u0b85\u0ba4\u0bbf\u0b95\u0bbe\u0bb2\u0bc8",
            "\u0b95\u0bbe\u0bb2\u0bc8",
            "\u0bae\u0ba4\u0bbf\u0baf\u0bae\u0bcd",
            "\u0baa\u0bbf\u0bb1\u0bcd\u0baa\u0b95\u0bb2\u0bcd",
            "\u0bae\u0bbe\u0bb2\u0bc8",
            "\u0b85\u0ba8\u0bcd\u0ba4\u0bbf \u0bae\u0bbe\u0bb2\u0bc8",
            "\u0b87\u0bb0\u0bb5\u0bc1",
            "",
        };
        final String[] metaValue_generic_narrow_AmPmMarkers = new String[] {
            "AM",
            "PM",
            "\u0ba8\u0bb3\u0bcd.",
            "\u0ba8\u0ba3\u0bcd.",
            "\u0b85\u0ba4\u0bbf.",
            "\u0b95\u0bbe.",
            "\u0bae\u0ba4\u0bbf.",
            "\u0baa\u0bbf\u0bb1\u0bcd.",
            "\u0bae\u0bbe.",
            "\u0b85\u0ba8\u0bcd\u0ba4\u0bbf \u0bae\u0bbe.",
            "\u0b87\u0bb0.",
            "",
        };
        final String[] metaValue_java_time_long_Eras = new String[] {
            "\u0b95\u0bbf\u0bb1\u0bbf\u0bb8\u0bcd\u0ba4\u0bc1\u0bb5\u0bc1\u0b95\u0bcd\u0b95\u0bc1 \u0bae\u0bc1\u0ba9\u0bcd",
            "\u0b85\u0ba9\u0bcd\u0ba9\u0bcb \u0b9f\u0bcb\u0bae\u0bbf\u0ba9\u0bbf",
        };
        final String[] metaValue_java_time_Eras = new String[] {
            "\u0b95\u0bbf.\u0bae\u0bc1.",
            "\u0b95\u0bbf.\u0baa\u0bbf.",
        };
        final String[] metaValue_generic_TimePatterns = new String[] {
            "h:mm:ss a zzzz",
            "h:mm:ss a z",
            "h:mm:ss a",
            "h:mm a",
        };
        final String metaValue_generic_DateFormatItem_MMdd = "dd-MM";
        final String metaValue_generic_DateFormatItem_yw = "Y \u0b86\u0bae\u0bcd \u0b86\u0ba3\u0bcd\u0b9f\u0bbf\u0ba9\u0bcd w -\u0b86\u0bae\u0bcd \u0bb5\u0bbe\u0bb0\u0bae\u0bcd";
        final String metaValue_generic_DateFormatItem_yQQQ = "QQQ y";
        final String metaValue_generic_DateFormatItem_EBhm = "E B h:mm";
        final String metaValue_generic_DateFormatItem_MMMMd = "d MMMM";
        final String metaValue_generic_DateFormatItem_Bh = "B h";
        final String metaValue_generic_DateFormatItem_MMMd = "d MMM";
        final String metaValue_generic_DateFormatItem_MMMMW = "MMMM W -\u0b86\u0bae\u0bcd \u0bb5\u0bbe\u0bb0\u0bae\u0bcd";
        final String metaValue_generic_DateFormatItem_yMMMEd = "E, d MMM, y";
        final String metaValue_generic_DateFormatItem_yMEd = "E, d/M/y";
        final String metaValue_generic_DateFormatItem_Bhms = "B h:mm:ss";
        final String metaValue_generic_DateFormatItem_yMMMd = "d MMM, y";
        final String metaValue_generic_DateFormatItem_yMMM = "MMM y";
        final String metaValue_generic_DateFormatItem_Md = "d/M";
        final String metaValue_generic_DateFormatItem_EBhms = "E B h:mm:ss";
        final String metaValue_generic_DateFormatItem_Ed = "d E";
        final String metaValue_generic_DateFormatItem_Bhm = "B h:mm";
        final String metaValue_generic_DateFormatItem_MEd = "dd-MM, E";
        final String metaValue_generic_DateFormatItem_yMM = "MM-y";
        final String metaValue_generic_DateFormatItem_yMd = "d/M/y";
        final String metaValue_generic_DateFormatItem_yM = "M/y";
        final String metaValue_generic_DateFormatItem_yMMMM = "MMMM y";
        final String metaValue_generic_DateFormatItem_yQQQQ = "QQQQ y";
        final String[] metaValue_java_time_generic_DatePatterns = new String[] {
            "EEEE, d MMMM, y G",
            "d MMMM, y G",
            "d MMM, y G",
            "d/M/y GGGGG",
        };
        final String[] metaValue_generic_DatePatterns = new String[] {
            "EEEE, d MMMM, y GGGG",
            "d MMMM, y GGGG",
            "d MMM, y GGGG",
            "d/M/y G",
        };
        final String[] metaValue_java_time_roc_long_Eras = new String[] {
            "R.O.C. -\u0b95\u0bcd\u0b95\u0bc1 \u0bae\u0bc1\u0ba9\u0bcd\u0baa\u0bc1",
            "R.O.C.",
        };
        final String[] metaValue_java_time_roc_Eras = new String[] {
            "ROC\u0b95\u0bcd\u0b95\u0bc1 \u0bae\u0bc1\u0ba9\u0bcd",
            "R.O.C.",
        };
        final String metaValue_calendarname_gregorian = "\u0b95\u0bbf\u0bb0\u0bbf\u0b95\u0bcb\u0bb0\u0bbf\u0baf\u0ba9\u0bcd \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf";
        final String[] metaValue_ListPatterns_standard = new String[] {
            "",
            "",
            "{0} \u0bae\u0bb1\u0bcd\u0bb1\u0bc1\u0bae\u0bcd {1}",
            "{0} \u0bae\u0bb1\u0bcd\u0bb1\u0bc1\u0bae\u0bcd {1}",
            "",
        };
        final String[] metaValue_ListPatterns_or = new String[] {
            "",
            "",
            "{0} \u0b85\u0bb2\u0bcd\u0bb2\u0ba4\u0bc1 {1}",
            "{0} \u0b85\u0bb2\u0bcd\u0bb2\u0ba4\u0bc1 {1}",
            "",
        };
        final Object[][] data = new Object[][] {
            { "generic.DayNames", metaValue_generic_DayNames },
            { "generic.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "generic.DayNarrows", metaValue_generic_DayNarrows },
            { "generic.QuarterNames", metaValue_generic_QuarterNames },
            { "generic.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "generic.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "generic.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "generic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.generic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "generic.DatePatterns", metaValue_generic_DatePatterns },
            { "generic.DateTimePatterns",
                new String[] {
                    "{1} {0}",
                    "{1} {0}",
                    "{1}, {0}",
                    "{1}, {0}",
                }
            },
            { "generic.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "generic.DateFormatItem.yyyyMMMM", "MMMM y G" },
            { "generic.DateFormatItem.yyyyQQQQ", "QQQQ y G" },
            { "generic.DateFormatItem.yyyyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "generic.DateFormatItem.Gy", metaValue_generic_DateFormatItem_Gy },
            { "generic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "generic.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "generic.DateFormatItem.yyyyMMM", metaValue_generic_DateFormatItem_yyyyMMM },
            { "generic.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "generic.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "generic.DateFormatItem.MMMEd", "E, d MMM" },
            { "generic.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "generic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "generic.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "generic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "generic.DateFormatItem.yyyy", metaValue_generic_DateFormatItem_Gy },
            { "generic.DateFormatItem.y", metaValue_generic_DateFormatItem_Gy },
            { "generic.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_yyyyMMM },
            { "generic.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "generic.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "generic.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "generic.DateFormatItem.yyyyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "generic.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "generic.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "generic.DateFormatItem.yyyyMEd", "E, d/M/y G" },
            { "generic.DateFormatItem.yyyyMd", metaValue_generic_DateFormatItem_GyMd },
            { "generic.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "generic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "generic.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "generic.DateFormatItem.yyyyMM", "MM-y G" },
            { "generic.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "generic.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "generic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "generic.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "generic.DateFormatItem.yyyyQQQ", "QQQ y G" },
            { "generic.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "generic.DateFormatItem.yyyyM", "M/y G" },
            { "generic.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "generic.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "generic.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "generic.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "MonthNames", metaValue_MonthNames },
            { "standalone.MonthNames", metaValue_MonthNames },
            { "MonthAbbreviations", metaValue_MonthAbbreviations },
            { "standalone.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "MonthNarrows", metaValue_MonthNarrows },
            { "standalone.MonthNarrows", metaValue_MonthNarrows },
            { "DayNames", metaValue_generic_DayNames },
            { "standalone.DayNames", metaValue_generic_DayNames },
            { "DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "standalone.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "DayNarrows", metaValue_generic_DayNarrows },
            { "standalone.DayNarrows", metaValue_generic_DayNarrows },
            { "QuarterNames", metaValue_generic_QuarterNames },
            { "standalone.QuarterNames", metaValue_generic_QuarterNames },
            { "QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "standalone.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "java.time.long.Eras", metaValue_java_time_long_Eras },
            { "long.Eras", metaValue_java_time_long_Eras },
            { "java.time.Eras", metaValue_java_time_Eras },
            { "Eras", metaValue_java_time_Eras },
            { "java.time.narrow.Eras", metaValue_java_time_Eras },
            { "narrow.Eras", metaValue_java_time_Eras },
            { "field.era", "\u0b95\u0bbe\u0bb2\u0bae\u0bcd" },
            { "field.year", "\u0b86\u0ba3\u0bcd\u0b9f\u0bc1" },
            { "field.month", "\u0bae\u0bbe\u0ba4\u0bae\u0bcd" },
            { "field.week", "\u0bb5\u0bbe\u0bb0\u0bae\u0bcd" },
            { "field.weekday", "\u0bb5\u0bbe\u0bb0\u0ba4\u0bcd\u0ba4\u0bbf\u0ba9\u0bcd \u0ba8\u0bbe\u0bb3\u0bcd" },
            { "field.dayperiod", "\u0bae\u0bc1\u0bb1\u0bcd\u0baa\u0b95\u0bb2\u0bcd/\u0baa\u0bbf\u0bb1\u0bcd\u0baa\u0b95\u0bb2\u0bcd" },
            { "field.hour", "\u0bae\u0ba3\u0bbf" },
            { "timezone.regionFormat", "{0} \u0ba8\u0bc7\u0bb0\u0bae\u0bcd" },
            { "timezone.regionFormat.daylight", "{0} \u0baa\u0b95\u0bb2\u0bca\u0bb3\u0bbf \u0ba8\u0bc7\u0bb0\u0bae\u0bcd" },
            { "timezone.regionFormat.standard", "{0} \u0ba8\u0bbf\u0bb2\u0bc8\u0baf\u0bbe\u0ba9 \u0ba8\u0bc7\u0bb0\u0bae\u0bcd" },
            { "field.minute", "\u0ba8\u0bbf\u0bae\u0bbf\u0b9f\u0bae\u0bcd" },
            { "field.second", "\u0bb5\u0bbf\u0ba8\u0bbe\u0b9f\u0bbf" },
            { "field.zone", "\u0ba8\u0bc7\u0bb0 \u0bae\u0ba3\u0bcd\u0b9f\u0bb2\u0bae\u0bcd" },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "EEEE, d MMMM, y",
                    "d MMMM, y",
                    "d MMM, y",
                    "d/M/yy",
                }
            },
            { "DateTimePatterns",
                new String[] {
                    "{1}, {0}",
                    "{1}, {0}",
                    "{1}, {0}",
                    "{1}, {0}",
                }
            },
            { "PluralRules", "one:n = 1" },
            { "DayPeriodRules", "midnight:00:00;noon:12:00;evening2:18:00-21:00;night1:21:00-03:00;morning2:05:00-12:00;afternoon1:12:00-14:00;morning1:03:00-05:00;afternoon2:14:00-16:00;evening1:16:00-18:00" },
            { "DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "buddhist.MonthNames", metaValue_MonthNames },
            { "buddhist.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "buddhist.DayNames", metaValue_generic_DayNames },
            { "buddhist.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "buddhist.QuarterNames", metaValue_generic_QuarterNames },
            { "buddhist.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "buddhist.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "buddhist.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "buddhist.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "buddhist.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.buddhist.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "buddhist.DatePatterns", metaValue_generic_DatePatterns },
            { "buddhist.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "buddhist.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "buddhist.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "buddhist.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "buddhist.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "buddhist.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "buddhist.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "buddhist.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "buddhist.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "buddhist.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "buddhist.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "buddhist.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "buddhist.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "buddhist.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "buddhist.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "buddhist.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "buddhist.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "buddhist.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "buddhist.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "buddhist.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "buddhist.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "buddhist.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "buddhist.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "buddhist.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "japanese.MonthNames", metaValue_MonthNames },
            { "japanese.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "japanese.DayNames", metaValue_generic_DayNames },
            { "japanese.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "japanese.QuarterNames", metaValue_generic_QuarterNames },
            { "japanese.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "japanese.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "japanese.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "japanese.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "japanese.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.japanese.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "japanese.DatePatterns",
                new String[] {
                    "EEEE, d MMMM, yyyy GGGG",
                    "d MMMM, yyyy GGGG",
                    "d MMM, y GGGG",
                    "d/M/y G",
                }
            },
            { "japanese.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "japanese.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "japanese.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "japanese.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "japanese.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "japanese.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "japanese.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "japanese.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "japanese.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "japanese.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "japanese.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "japanese.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "japanese.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "japanese.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "japanese.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "japanese.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "japanese.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "japanese.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "japanese.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "japanese.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "japanese.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "japanese.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "japanese.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "japanese.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "roc.MonthNames", metaValue_MonthNames },
            { "roc.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "roc.DayNames", metaValue_generic_DayNames },
            { "roc.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "roc.DayNarrows", metaValue_generic_DayNarrows },
            { "roc.QuarterNames", metaValue_generic_QuarterNames },
            { "roc.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "roc.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "roc.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "java.time.roc.long.Eras", metaValue_java_time_roc_long_Eras },
            { "roc.long.Eras", metaValue_java_time_roc_long_Eras },
            { "java.time.roc.Eras", metaValue_java_time_roc_Eras },
            { "roc.Eras", metaValue_java_time_roc_Eras },
            { "java.time.roc.narrow.Eras", metaValue_java_time_roc_Eras },
            { "roc.narrow.Eras", metaValue_java_time_roc_Eras },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.roc.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "roc.DatePatterns", metaValue_generic_DatePatterns },
            { "roc.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "roc.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "roc.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "roc.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "roc.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "roc.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "roc.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "roc.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "roc.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "roc.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "roc.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "roc.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "roc.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "roc.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "roc.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "roc.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "roc.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "roc.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "roc.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "roc.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "roc.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "roc.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "roc.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "roc.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "islamic.MonthNames",
                new String[] {
                    "\u0bae\u0bc1\u0bb9\u0bb0\u0bcd\u0bb0\u0bae\u0bcd",
                    "\u0b9a\u0b83\u0baa\u0bb0\u0bcd",
                    "\u0bb0\u0baa\u0bbf 1",
                    "\u0bb0\u0baa\u0bbf 2",
                    "\u0b9c\u0bc1\u0bae\u0ba4\u0bbe 1",
                    "\u0b9c\u0bc1\u0bae\u0ba4\u0bbe 2",
                    "\u0bb0\u0b9c\u0baa\u0bcd",
                    "\u0bb7\u0b83\u0baa\u0bbe\u0ba9\u0bcd",
                    "\u0bb0\u0bae\u0bb2\u0bbe\u0ba9\u0bcd",
                    "\u0bb7\u0bb5\u0bcd\u0bb5\u0bbe\u0bb2\u0bcd",
                    "\u0ba4\u0bc1\u0bb2\u0bcd \u0b95\u0b83\u0ba4\u0bbe",
                    "\u0ba4\u0bc1\u0bb2\u0bcd \u0bb9\u0bbf\u0b9c\u0bcd\u0b9c\u0bbe",
                    "",
                }
            },
            { "islamic.MonthAbbreviations",
                new String[] {
                    "\u0bae\u0bc1\u0bb9.",
                    "\u0b9a\u0b83\u0baa.",
                    "\u0bb0\u0baa\u0bbf 1",
                    "\u0bb0\u0baa\u0bbf 2",
                    "\u0b9c\u0bc1\u0bae. 1",
                    "\u0b9c\u0bc1\u0bae. 2",
                    "\u0bb0\u0b9c.",
                    "\u0bb7\u0b83.",
                    "\u0bb0\u0bae.",
                    "\u0bb7\u0bb5\u0bcd.",
                    "\u0ba4\u0bc1\u0bb2\u0bcd \u0b95\u0b83.",
                    "\u0ba4\u0bc1\u0bb2\u0bcd \u0bb9\u0bbf\u0b9c\u0bcd.",
                    "",
                }
            },
            { "islamic.DayNames", metaValue_generic_DayNames },
            { "islamic.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "islamic.QuarterNames", metaValue_generic_QuarterNames },
            { "islamic.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "islamic.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "islamic.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "islamic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.islamic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "islamic.DatePatterns", metaValue_generic_DatePatterns },
            { "islamic.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "islamic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "islamic.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "islamic.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "islamic.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "islamic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "islamic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "islamic.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "islamic.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "islamic.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "islamic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "islamic.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "islamic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "islamic.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "islamic.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "islamic-civil.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "islamic-civil.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-civil.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "islamic-civil.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "islamic-civil.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "islamic-civil.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "islamic-civil.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-civil.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "islamic-civil.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-civil.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "islamic-civil.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic-civil.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "islamic-civil.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "islamic-civil.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "islamic-civil.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-civil.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "islamic-civil.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic-civil.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "islamic-civil.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-civil.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "islamic-civil.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic-civil.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic-civil.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "islamic-civil.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "islamic-umalqura.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "islamic-umalqura.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-umalqura.DateFormatItem.yQQQ", metaValue_generic_DateFormatItem_yQQQ },
            { "islamic-umalqura.DateFormatItem.EBhm", metaValue_generic_DateFormatItem_EBhm },
            { "islamic-umalqura.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "islamic-umalqura.DateFormatItem.Bh", metaValue_generic_DateFormatItem_Bh },
            { "islamic-umalqura.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-umalqura.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "islamic-umalqura.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-umalqura.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "islamic-umalqura.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic-umalqura.DateFormatItem.Bhms", metaValue_generic_DateFormatItem_Bhms },
            { "islamic-umalqura.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "islamic-umalqura.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "islamic-umalqura.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-umalqura.DateFormatItem.EBhms", metaValue_generic_DateFormatItem_EBhms },
            { "islamic-umalqura.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic-umalqura.DateFormatItem.Bhm", metaValue_generic_DateFormatItem_Bhm },
            { "islamic-umalqura.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-umalqura.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "islamic-umalqura.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic-umalqura.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic-umalqura.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "islamic-umalqura.DateFormatItem.yQQQQ", metaValue_generic_DateFormatItem_yQQQQ },
            { "calendarname.islamic-civil", "\u0b87\u0bb8\u0bcd\u0bb2\u0bbe\u0bae\u0bbf\u0baf \u0b9a\u0bbf\u0bb5\u0bbf\u0bb2\u0bcd \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf" },
            { "calendarname.islamic", "\u0b87\u0bb8\u0bcd\u0bb2\u0bbe\u0bae\u0bbf\u0baf \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf" },
            { "calendarname.buddhist", "\u0baa\u0bc1\u0ba4\u0bcd\u0ba4 \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf" },
            { "calendarname.japanese", "\u0b9c\u0baa\u0bcd\u0baa\u0bbe\u0ba9\u0bbf\u0baf \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf" },
            { "calendarname.roc", "\u0bae\u0bbf\u0ba9\u0bcd\u0b95\u0bcb \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf" },
            { "calendarname.islamic-umalqura", "\u0b87\u0bb8\u0bcd\u0bb2\u0bbe\u0bae\u0bbf\u0baf \u0ba8\u0bbe\u0bb3\u0bcd\u0b95\u0bbe\u0b9f\u0bcd\u0b9f\u0bbf (\u0b89\u0bae\u0bcd-\u0b85\u0bb2\u0bcd-\u0b95\u0bc1\u0bb0\u0bbe)" },
            { "calendarname.gregorian", metaValue_calendarname_gregorian },
            { "calendarname.gregory", metaValue_calendarname_gregorian },
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
                    "#,##,##0.###",
                    "\u00a4#,##,##0.00",
                    "#,##,##0%",
                    "\u00a4#,##0.00;(\u00a4#,##0.00)",
                }
            },
            { "tamldec.NumberElements",
                new String[] {
                    ".",
                    ",",
                    ";",
                    "%",
                    "\u0be6",
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
            { "tamldec.NumberPatterns",
                new String[] {
                    "#,##0.###",
                    "\u00a4#,##0.00",
                    "#,##0%",
                    "",
                }
            },
            { "short.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:0\u0b86 other:0\u0b86}",
                    "{one:00\u0b86 other:00\u0b86}",
                    "{one:000\u0b86 other:000\u0b86}",
                    "{one:0\u0bae\u0bbf other:0\u0bae\u0bbf}",
                    "{one:00\u0bae\u0bbf other:00\u0bae\u0bbf}",
                    "{one:000\u0bae\u0bbf other:000\u0bae\u0bbf}",
                    "{one:0\u0baa\u0bbf other:0\u0baa\u0bbf}",
                    "{one:00\u0baa\u0bbf other:00\u0baa\u0bbf}",
                    "{one:000\u0baa\u0bbf other:000\u0baa\u0bbf}",
                    "{one:0\u0b9f\u0bbf other:0\u0b9f\u0bbf}",
                    "{one:00\u0b9f\u0bbf other:00\u0b9f\u0bbf}",
                    "{one:000\u0b9f\u0bbf other:000\u0b9f\u0bbf}",
                }
            },
            { "long.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:0' '\u0b86\u0baf\u0bbf\u0bb0\u0bae\u0bcd other:0' '\u0b86\u0baf\u0bbf\u0bb0\u0bae\u0bcd}",
                    "{one:00' '\u0b86\u0baf\u0bbf\u0bb0\u0bae\u0bcd other:00' '\u0b86\u0baf\u0bbf\u0bb0\u0bae\u0bcd}",
                    "{one:000' '\u0b86\u0baf\u0bbf\u0bb0\u0bae\u0bcd other:000' '\u0b86\u0baf\u0bbf\u0bb0\u0bae\u0bcd}",
                    "{one:0' '\u0bae\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:0' '\u0bae\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:00' '\u0bae\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:00' '\u0bae\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:000' '\u0bae\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:000' '\u0bae\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:0' '\u0baa\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:0' '\u0baa\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:00' '\u0baa\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:00' '\u0baa\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:000' '\u0baa\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:000' '\u0baa\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:0' '\u0b9f\u0bbf\u0bb0\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:0' '\u0b9f\u0bbf\u0bb0\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:00' '\u0b9f\u0bbf\u0bb0\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:00' '\u0b9f\u0bbf\u0bb0\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                    "{one:000' '\u0b9f\u0bbf\u0bb0\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd other:000' '\u0b9f\u0bbf\u0bb0\u0bbf\u0bb2\u0bcd\u0bb2\u0bbf\u0baf\u0ba9\u0bcd}",
                }
            },
            { "ListPatterns_standard", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-short", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-narrow", metaValue_ListPatterns_standard },
            { "ListPatterns_or", metaValue_ListPatterns_or },
            { "ListPatterns_or-short", metaValue_ListPatterns_or },
            { "ListPatterns_or-narrow", metaValue_ListPatterns_or },
            { "ListPatterns_unit", metaValue_ListPatterns_standard },
            { "ListPatterns_unit-short",
                new String[] {
                    "",
                    "",
                    "{0}, {1}",
                    "{0}, {1}",
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
