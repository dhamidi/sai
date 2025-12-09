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

public class FormatData_gu extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String metaValue_generic_DateFormatItem_yyyyMMMd = "d MMM, G y";
        final String metaValue_generic_DateFormatItem_yyyyMMMEd = "E, d MMM, G y";
        final String[] metaValue_MonthNames = new String[] {
            "\u0a9c\u0abe\u0aa8\u0acd\u0aaf\u0ac1\u0a86\u0ab0\u0ac0",
            "\u0aab\u0ac7\u0aac\u0acd\u0ab0\u0ac1\u0a86\u0ab0\u0ac0",
            "\u0aae\u0abe\u0ab0\u0acd\u0a9a",
            "\u0a8f\u0aaa\u0acd\u0ab0\u0abf\u0ab2",
            "\u0aae\u0ac7",
            "\u0a9c\u0ac2\u0aa8",
            "\u0a9c\u0ac1\u0ab2\u0abe\u0a88",
            "\u0a91\u0a97\u0ab8\u0acd\u0a9f",
            "\u0ab8\u0aaa\u0acd\u0a9f\u0ac7\u0aae\u0acd\u0aac\u0ab0",
            "\u0a91\u0a95\u0acd\u0a9f\u0acb\u0aac\u0ab0",
            "\u0aa8\u0ab5\u0ac7\u0aae\u0acd\u0aac\u0ab0",
            "\u0aa1\u0abf\u0ab8\u0ac7\u0aae\u0acd\u0aac\u0ab0",
            "",
        };
        final String[] metaValue_MonthAbbreviations = new String[] {
            "\u0a9c\u0abe\u0aa8\u0acd\u0aaf\u0ac1",
            "\u0aab\u0ac7\u0aac\u0acd\u0ab0\u0ac1",
            "\u0aae\u0abe\u0ab0\u0acd\u0a9a",
            "\u0a8f\u0aaa\u0acd\u0ab0\u0abf\u0ab2",
            "\u0aae\u0ac7",
            "\u0a9c\u0ac2\u0aa8",
            "\u0a9c\u0ac1\u0ab2\u0abe\u0a88",
            "\u0a91\u0a97\u0ab8\u0acd\u0a9f",
            "\u0ab8\u0aaa\u0acd\u0a9f\u0ac7",
            "\u0a91\u0a95\u0acd\u0a9f\u0acb",
            "\u0aa8\u0ab5\u0ac7",
            "\u0aa1\u0abf\u0ab8\u0ac7",
            "",
        };
        final String[] metaValue_MonthNarrows = new String[] {
            "\u0a9c\u0abe",
            "\u0aab\u0ac7",
            "\u0aae\u0abe",
            "\u0a8f",
            "\u0aae\u0ac7",
            "\u0a9c\u0ac2",
            "\u0a9c\u0ac1",
            "\u0a91",
            "\u0ab8",
            "\u0a91",
            "\u0aa8",
            "\u0aa1\u0abf",
            "",
        };
        final String[] metaValue_generic_DayNames = new String[] {
            "\u0ab0\u0ab5\u0abf\u0ab5\u0abe\u0ab0",
            "\u0ab8\u0acb\u0aae\u0ab5\u0abe\u0ab0",
            "\u0aae\u0a82\u0a97\u0ab3\u0ab5\u0abe\u0ab0",
            "\u0aac\u0ac1\u0aa7\u0ab5\u0abe\u0ab0",
            "\u0a97\u0ac1\u0ab0\u0ac1\u0ab5\u0abe\u0ab0",
            "\u0ab6\u0ac1\u0a95\u0acd\u0ab0\u0ab5\u0abe\u0ab0",
            "\u0ab6\u0aa8\u0abf\u0ab5\u0abe\u0ab0",
        };
        final String[] metaValue_generic_DayAbbreviations = new String[] {
            "\u0ab0\u0ab5\u0abf",
            "\u0ab8\u0acb\u0aae",
            "\u0aae\u0a82\u0a97\u0ab3",
            "\u0aac\u0ac1\u0aa7",
            "\u0a97\u0ac1\u0ab0\u0ac1",
            "\u0ab6\u0ac1\u0a95\u0acd\u0ab0",
            "\u0ab6\u0aa8\u0abf",
        };
        final String[] metaValue_generic_DayNarrows = new String[] {
            "\u0ab0",
            "\u0ab8\u0acb",
            "\u0aae\u0a82",
            "\u0aac\u0ac1",
            "\u0a97\u0ac1",
            "\u0ab6\u0ac1",
            "\u0ab6",
        };
        final String[] metaValue_generic_QuarterNames = new String[] {
            "1\u0ab2\u0acb \u0aa4\u0acd\u0ab0\u0abf\u0aae\u0abe\u0ab8",
            "2\u0a9c\u0acb \u0aa4\u0acd\u0ab0\u0abf\u0aae\u0abe\u0ab8",
            "3\u0a9c\u0acb \u0aa4\u0acd\u0ab0\u0abf\u0aae\u0abe\u0ab8",
            "4\u0aa5\u0acb \u0aa4\u0acd\u0ab0\u0abf\u0aae\u0abe\u0ab8",
        };
        final String[] metaValue_generic_QuarterAbbreviations = new String[] {
            "Q1",
            "Q2",
            "Q3",
            "Q4",
        };
        final String[] metaValue_generic_abbreviated_AmPmMarkers = new String[] {
            "AM",
            "PM",
            "\u0aae\u0aa7\u0acd\u0aaf\u0ab0\u0abe\u0aa4\u0acd\u0ab0\u0abf",
            "",
            "\u0ab8\u0ab5\u0abe\u0ab0\u0ac7",
            "",
            "\u0aac\u0aaa\u0acb\u0ab0\u0ac7",
            "",
            "\u0ab8\u0abe\u0a82\u0a9c\u0ac7",
            "",
            "\u0ab0\u0abe\u0aa4\u0acd\u0ab0\u0ac7",
            "",
        };
        final String[] metaValue_generic_narrow_AmPmMarkers = new String[] {
            "AM",
            "PM",
            "\u0aae.\u0ab0\u0abe\u0aa4\u0acd\u0ab0\u0abf",
            "",
            "\u0ab8\u0ab5\u0abe\u0ab0\u0ac7",
            "",
            "\u0aac\u0aaa\u0acb\u0ab0\u0ac7",
            "",
            "\u0ab8\u0abe\u0a82\u0a9c\u0ac7",
            "",
            "\u0ab0\u0abe\u0aa4\u0acd\u0ab0\u0ac7",
            "",
        };
        final String[] metaValue_java_time_long_Eras = new String[] {
            "\u0a88\u0ab8\u0ab5\u0ac0\u0ab8\u0aa8 \u0aaa\u0ac2\u0ab0\u0acd\u0ab5\u0ac7",
            "\u0a87\u0ab8\u0ab5\u0ac0\u0ab8\u0aa8",
        };
        final String[] metaValue_java_time_Eras = new String[] {
            "\u0a88.\u0ab8.\u0aaa\u0ac2\u0ab0\u0acd\u0ab5\u0ac7",
            "\u0a88.\u0ab8.",
        };
        final String[] metaValue_java_time_narrow_Eras = new String[] {
            "\u0a87 \u0ab8 \u0aaa\u0ac1",
            "\u0a87\u0ab8",
        };
        final String[] metaValue_generic_TimePatterns = new String[] {
            "hh:mm:ss a zzzz",
            "hh:mm:ss a z",
            "hh:mm:ss a",
            "hh:mm a",
        };
        final String[] metaValue_generic_DateTimePatterns = new String[] {
            "{1} {0}",
            "{1} {0}",
            "{1} {0}",
            "{1} {0}",
        };
        final String metaValue_generic_DateFormatItem_MMdd = "dd-MM";
        final String metaValue_generic_DateFormatItem_yMMMEd = "E, d MMM, y";
        final String metaValue_generic_DateFormatItem_yMEd = "E, d/M/y";
        final String metaValue_generic_DateFormatItem_yMMMd = "d MMM, y";
        final String metaValue_generic_DateFormatItem_yw = "Y \u0aa8\u0abe w \u0a85\u0aa0\u0ab5\u0abe\u0aa1\u0abf\u0aaf\u0abe";
        final String metaValue_generic_DateFormatItem_yMMM = "MMM y";
        final String metaValue_generic_DateFormatItem_Md = "d/M";
        final String metaValue_generic_DateFormatItem_Ed = "E d";
        final String metaValue_generic_DateFormatItem_MMMMd = "d MMMM";
        final String metaValue_generic_DateFormatItem_MEd = "E, d/M";
        final String metaValue_generic_DateFormatItem_MMMEd = "E, d MMM";
        final String metaValue_generic_DateFormatItem_yMM = "MM-y";
        final String metaValue_generic_DateFormatItem_yMd = "d/M/y";
        final String metaValue_generic_DateFormatItem_MMMd = "d MMM";
        final String metaValue_generic_DateFormatItem_yM = "M/y";
        final String metaValue_generic_DateFormatItem_yMMMM = "MMMM y";
        final String metaValue_generic_DateFormatItem_MMMMW = "MMMM\u0aa8\u0abe W \u0a85\u0aa0\u0ab5\u0abe\u0aa1\u0abf\u0aaf\u0abe";
        final String metaValue_generic_DateFormatItem_GyMMM = "MMM, G y";
        final String[] metaValue_java_time_generic_DatePatterns = new String[] {
            "EEEE, d MMMM, G y",
            "d MMMM, G y",
            "d MMM, G y",
            "d-MM- GGGGG y",
        };
        final String[] metaValue_generic_DatePatterns = new String[] {
            "EEEE, d MMMM, GGGG y",
            "d MMMM, GGGG y",
            "d MMM, GGGG y",
            "d-MM- G y",
        };
        final String[] metaValue_java_time_roc_long_Eras = new String[] {
            "\u0a86\u0ab0.\u0a93.\u0ab8\u0ac0. \u0aaa\u0ab9\u0ac7\u0ab2\u0abe\u0a82",
            "\u0a86\u0ab0.\u0a93.\u0ab8\u0ac0.",
        };
        final String metaValue_calendarname_gregorian = "\u0a97\u0acd\u0ab0\u0ac7\u0a97\u0acb\u0ab0\u0abf\u0a85\u0aa8 \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0";
        final String[] metaValue_ListPatterns_standard = new String[] {
            "",
            "",
            "{0} \u0a85\u0aa8\u0ac7 {1}",
            "{0} \u0a85\u0aa8\u0ac7 {1}",
            "",
        };
        final String[] metaValue_ListPatterns_or = new String[] {
            "",
            "",
            "{0}, \u0a85\u0aa5\u0ab5\u0abe {1}",
            "{0} \u0a85\u0aa5\u0ab5\u0abe {1}",
            "",
        };
        final String[] metaValue_ListPatterns_standard_narrow = new String[] {
            "",
            "",
            "{0}, {1}",
            "{0}, {1}",
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
            { "generic.DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "generic.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "generic.DateFormatItem.hm", "hh:mm a" },
            { "generic.DateFormatItem.yyyyMMMM", "MMMM G y" },
            { "generic.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "generic.DateFormatItem.yyyyQQQQ", "QQQQ G y" },
            { "generic.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "generic.DateFormatItem.hms", "hh:mm:ss a" },
            { "generic.DateFormatItem.yyyyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "generic.DateFormatItem.yyyyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "generic.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "generic.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "generic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "generic.DateFormatItem.yyyyMEd", "E, d/M/G y" },
            { "generic.DateFormatItem.yyyyMd", "d/M/G y" },
            { "generic.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "generic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "generic.DateFormatItem.yyyyMM", "MM-y G" },
            { "generic.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "generic.DateFormatItem.yyyyMMM", "MMM G y" },
            { "generic.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "generic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "generic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "generic.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "generic.DateFormatItem.yyyyQQQ", "QQQ G y" },
            { "generic.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "generic.DateFormatItem.yyyyM", "M/G y" },
            { "generic.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "generic.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "generic.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "generic.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "generic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "generic.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
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
            { "standalone.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "java.time.long.Eras", metaValue_java_time_long_Eras },
            { "long.Eras", metaValue_java_time_long_Eras },
            { "java.time.Eras", metaValue_java_time_Eras },
            { "Eras", metaValue_java_time_Eras },
            { "java.time.narrow.Eras", metaValue_java_time_narrow_Eras },
            { "narrow.Eras", metaValue_java_time_narrow_Eras },
            { "field.era", "\u0aaf\u0ac1\u0a97" },
            { "field.year", "\u0ab5\u0ab0\u0acd\u0ab7" },
            { "field.month", "\u0aae\u0ab9\u0abf\u0aa8\u0acb" },
            { "field.week", "\u0a85\u0aa0\u0ab5\u0abe\u0aa1\u0abf\u0aaf\u0ac1\u0a82" },
            { "field.weekday", "\u0a85\u0aa0\u0ab5\u0abe\u0aa1\u0abf\u0aaf\u0abe\u0aa8\u0acb \u0aa6\u0abf\u0ab5\u0ab8" },
            { "field.dayperiod", "AM/PM" },
            { "field.hour", "\u0a95\u0ab2\u0abe\u0a95" },
            { "timezone.regionFormat", "{0} \u0ab8\u0aae\u0aaf" },
            { "timezone.regionFormat.daylight", "{0} \u0aa6\u0abf\u0ab5\u0ab8 \u0ab8\u0aae\u0aaf" },
            { "timezone.regionFormat.standard", "{0} \u0aae\u0abe\u0aa8\u0a95 \u0ab8\u0aae\u0aaf" },
            { "field.minute", "\u0aae\u0abf\u0aa8\u0abf\u0a9f" },
            { "field.second", "\u0ab8\u0ac7\u0a95\u0aa8\u0acd\u0aa1" },
            { "field.zone", "\u0ab8\u0aae\u0aaf \u0a9d\u0acb\u0aa8" },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "EEEE, d MMMM, y",
                    "d MMMM, y",
                    "d MMM, y",
                    "d/M/yy",
                }
            },
            { "DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "PluralRules", "one:i = 0 or n = 1" },
            { "DayPeriodRules", "midnight:00:00;night1:20:00-04:00;afternoon1:12:00-16:00;morning1:04:00-12:00;evening1:16:00-20:00" },
            { "DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
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
            { "buddhist.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "buddhist.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "buddhist.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "buddhist.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "buddhist.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "buddhist.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "buddhist.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "buddhist.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "buddhist.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "buddhist.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "buddhist.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "buddhist.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "buddhist.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "buddhist.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "buddhist.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "buddhist.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "buddhist.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "buddhist.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "buddhist.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
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
                    "EEEE, d MMMM, GGGG yyyy",
                    "d MMMM, GGGG yyyy",
                    "d MMM, GGGG y",
                    "d-MM- G y",
                }
            },
            { "japanese.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "japanese.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "japanese.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "japanese.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "japanese.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "japanese.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "japanese.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "japanese.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "japanese.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "japanese.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "japanese.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "japanese.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "japanese.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "japanese.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "japanese.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "japanese.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "japanese.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "japanese.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "japanese.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "japanese.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
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
            { "java.time.roc.Eras", metaValue_java_time_roc_long_Eras },
            { "roc.Eras", metaValue_java_time_roc_long_Eras },
            { "java.time.roc.narrow.Eras", metaValue_java_time_roc_long_Eras },
            { "roc.narrow.Eras", metaValue_java_time_roc_long_Eras },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.roc.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "roc.DatePatterns", metaValue_generic_DatePatterns },
            { "roc.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "roc.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "roc.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "roc.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "roc.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "roc.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "roc.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "roc.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "roc.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "roc.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "roc.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "roc.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "roc.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "roc.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "roc.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "roc.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "roc.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "roc.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "roc.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "roc.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
            { "islamic.MonthNames",
                new String[] {
                    "\u0aae\u0ac1\u0ab9\u0ab0\u0acd\u0ab0\u0aae",
                    "\u0ab8\u0aab\u0ab0",
                    "\u0ab0\u0abe\u0aac\u0ac0\u02bb I",
                    "\u0ab0\u0abe\u0aac\u0ac0\u02bb II",
                    "\u0a9c\u0ac1\u0aae\u0abe\u0aa6\u0abe I",
                    "\u0a9c\u0ac1\u0aae\u0abe\u0aa6\u0abe II",
                    "\u0ab0\u0a9c\u0aac",
                    "\u0ab6\u0abe\u02bb\u0aac\u0abe\u0aa8",
                    "\u0ab0\u0aae\u0aa6\u0abe\u0aa8",
                    "\u0ab6\u0abe\u0ab5\u0acd\u0ab5\u0ab2",
                    "\u0aa7\u0ac1\u02bb\u0ab2-\u0a95\u0acd\u0ab5\u0ac0\u02bb\u0aa1\u0abe\u0ab9",
                    "\u0aa7\u0ac1\u02bb\u0ab2-\u0ab9\u0abf\u0a9c\u0acd\u0a9c\u0abe\u0ab9",
                    "",
                }
            },
            { "islamic.MonthAbbreviations",
                new String[] {
                    "\u0aae\u0ac1\u0ab9.",
                    "\u0ab8\u0aab.",
                    "\u0ab0\u0aac.I",
                    "\u0ab0\u0aac. II",
                    "\u0a9c\u0ac1\u0aae. I",
                    "\u0a9c\u0ac1\u0aae. II",
                    "\u0ab0\u0abe\u0a9c.",
                    "\u0ab6\u0abe.",
                    "\u0ab0\u0abe\u0aae.",
                    "\u0ab6\u0abe\u0ab5.",
                    "\u0aa7\u0ac1\u02bb\u0ab2-\u0a95\u0acd\u0aaf\u0ac1.",
                    "\u0aa7\u0ac1\u02bb\u0ab2-\u0a8f\u0a9a.",
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
            { "islamic.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "islamic.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "islamic.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "islamic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "islamic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "islamic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "islamic.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "islamic.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "islamic.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "islamic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
            { "islamic-civil.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "islamic-civil.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "islamic-civil.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic-civil.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "islamic-civil.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "islamic-civil.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-civil.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "islamic-civil.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-civil.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic-civil.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "islamic-civil.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-civil.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-civil.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "islamic-civil.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "islamic-civil.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic-civil.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "islamic-civil.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic-civil.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "islamic-civil.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-civil.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
            { "islamic-umalqura.DateFormatItem.MMdd", metaValue_generic_DateFormatItem_MMdd },
            { "islamic-umalqura.DateFormatItem.yMMMEd", metaValue_generic_DateFormatItem_yMMMEd },
            { "islamic-umalqura.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic-umalqura.DateFormatItem.GyMMMd", metaValue_generic_DateFormatItem_yyyyMMMd },
            { "islamic-umalqura.DateFormatItem.yMMMd", metaValue_generic_DateFormatItem_yMMMd },
            { "islamic-umalqura.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-umalqura.DateFormatItem.yMMM", metaValue_generic_DateFormatItem_yMMM },
            { "islamic-umalqura.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-umalqura.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic-umalqura.DateFormatItem.MMMMd", metaValue_generic_DateFormatItem_MMMMd },
            { "islamic-umalqura.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-umalqura.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-umalqura.DateFormatItem.yMM", metaValue_generic_DateFormatItem_yMM },
            { "islamic-umalqura.DateFormatItem.GyMMMEd", metaValue_generic_DateFormatItem_yyyyMMMEd },
            { "islamic-umalqura.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic-umalqura.DateFormatItem.MMMd", metaValue_generic_DateFormatItem_MMMd },
            { "islamic-umalqura.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic-umalqura.DateFormatItem.yMMMM", metaValue_generic_DateFormatItem_yMMMM },
            { "islamic-umalqura.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-umalqura.DateFormatItem.GyMMM", metaValue_generic_DateFormatItem_GyMMM },
            { "calendarname.islamic-civil", "\u0ab9\u0abf\u0a9c\u0ab0\u0ac0-\u0aa8\u0abe\u0a97\u0ab0\u0abf\u0a95 \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0" },
            { "calendarname.islamic", "\u0ab9\u0abf\u0a9c\u0ab0\u0ac0 \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0" },
            { "calendarname.buddhist", "\u0aac\u0ac1\u0aa6\u0acd\u0aa7\u0abf\u0ab8\u0acd\u0a9f \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0" },
            { "calendarname.japanese", "\u0a9c\u0abe\u0aaa\u0abe\u0aa8\u0ac0\u0a9d \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0" },
            { "calendarname.roc", "\u0aae\u0abf\u0a82\u0a97\u0ac1\u0a93 \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0" },
            { "calendarname.islamic-umalqura", "\u0ab9\u0abf\u0a9c\u0ab0\u0ac0 \u0a95\u0ac7\u0ab2\u0ac7\u0aa8\u0acd\u0aa1\u0ab0 (\u0a89\u0aae \u0a85\u0ab2-\u0a95\u0ac1\u0ab0\u0abe)" },
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
                    "\u00a4#,##,##0.00;(\u00a4#,##,##0.00)",
                }
            },
            { "gujr.NumberElements",
                new String[] {
                    ".",
                    ",",
                    ";",
                    "%",
                    "\u0ae6",
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
            { "gujr.NumberPatterns",
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
                    "{one:0\u00a0\u0ab9\u0a9c\u0abe\u0ab0 other:0\u00a0\u0ab9\u0a9c\u0abe\u0ab0}",
                    "{one:00\u00a0\u0ab9\u0a9c\u0abe\u0ab0 other:00\u00a0\u0ab9\u0a9c\u0abe\u0ab0}",
                    "{one:0\u00a0\u0ab2\u0abe\u0a96 other:0\u00a0\u0ab2\u0abe\u0a96}",
                    "{one:00\u00a0\u0ab2\u0abe\u0a96 other:00\u00a0\u0ab2\u0abe\u0a96}",
                    "{one:0\u00a0\u0a95\u0ab0\u0acb\u0aa1 other:0\u00a0\u0a95\u0ab0\u0acb\u0aa1}",
                    "{one:00\u00a0\u0a95\u0ab0\u0acb\u0aa1 other:00\u00a0\u0a95\u0ab0\u0acb\u0aa1}",
                    "{one:0\u00a0\u0a85\u0aac\u0a9c other:0\u00a0\u0a85\u0aac\u0a9c}",
                    "{one:00\u00a0\u0a85\u0aac\u0a9c other:00\u00a0\u0a85\u0aac\u0a9c}",
                    "{one:0\u00a0\u0aa8\u0abf\u0a96\u0ab0\u0acd\u0ab5 other:0\u00a0\u0aa8\u0abf\u0a96\u0ab0\u0acd\u0ab5}",
                    "{one:0\u00a0\u0aae\u0ab9\u0abe\u0aaa\u0aa6\u0acd\u0aae other:0\u00a0\u0aae\u0ab9\u0abe\u0aaa\u0aa6\u0acd\u0aae}",
                    "{one:0\u00a0\u0ab6\u0a82\u0a95\u0ac1 other:0\u00a0\u0ab6\u0a82\u0a95\u0ac1}",
                    "{one:0\u00a0\u0a9c\u0ab2\u0aa7\u0abf other:0\u00a0\u0a9c\u0ab2\u0aa7\u0abf}",
                }
            },
            { "long.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:0' '\u0ab9\u0a9c\u0abe\u0ab0 other:0' '\u0ab9\u0a9c\u0abe\u0ab0}",
                    "{one:00' '\u0ab9\u0a9c\u0abe\u0ab0 other:00' '\u0ab9\u0a9c\u0abe\u0ab0}",
                    "{one:0' '\u0ab2\u0abe\u0a96 other:0' '\u0ab2\u0abe\u0a96}",
                    "{one:00' '\u0ab2\u0abe\u0a96 other:00' '\u0ab2\u0abe\u0a96}",
                    "{one:0' '\u0a95\u0ab0\u0acb\u0aa1 other:0' '\u0a95\u0ab0\u0acb\u0aa1}",
                    "{one:00' '\u0a95\u0ab0\u0acb\u0aa1 other:00' '\u0a95\u0ab0\u0acb\u0aa1}",
                    "{one:0' '\u0a85\u0aac\u0a9c other:0' '\u0a85\u0aac\u0a9c}",
                    "{one:00' '\u0a85\u0aac\u0a9c other:00' '\u0a85\u0aac\u0a9c}",
                    "{one:0' '\u0aa8\u0abf\u0a96\u0ab0\u0acd\u0ab5 other:0' '\u0aa8\u0abf\u0a96\u0ab0\u0acd\u0ab5}",
                    "{one:0' '\u0aae\u0ab9\u0abe\u0aaa\u0aa6\u0acd\u0aae other:0' '\u0aae\u0ab9\u0abe\u0aaa\u0aa6\u0acd\u0aae}",
                    "{one:0' '\u0ab6\u0a82\u0a95\u0ac1 other:0' '\u0ab6\u0a82\u0a95\u0ac1}",
                    "{one:0' '\u0a9c\u0ab2\u0aa7\u0abf other:0' '\u0a9c\u0ab2\u0aa7\u0abf}",
                }
            },
            { "ListPatterns_standard", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-short", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-narrow", metaValue_ListPatterns_standard_narrow },
            { "ListPatterns_or", metaValue_ListPatterns_or },
            { "ListPatterns_or-short",
                new String[] {
                    "",
                    "",
                    "{0} \u0a85\u0aa5\u0ab5\u0abe {1}",
                    "{0} \u0a85\u0aa5\u0ab5\u0abe {1}",
                    "",
                }
            },
            { "ListPatterns_or-narrow", metaValue_ListPatterns_or },
            { "ListPatterns_unit", metaValue_ListPatterns_standard },
            { "ListPatterns_unit-short", metaValue_ListPatterns_standard_narrow },
            { "ListPatterns_unit-narrow", metaValue_ListPatterns_standard },
        };
        return data;
    }
}
