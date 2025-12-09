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

public class FormatData_si extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] metaValue_MonthNames = new String[] {
            "\u0da2\u0db1\u0dc0\u0dcf\u0dbb\u0dd2",
            "\u0db4\u0dd9\u0db6\u0dbb\u0dc0\u0dcf\u0dbb\u0dd2",
            "\u0db8\u0dcf\u0dbb\u0dca\u0dad\u0dd4",
            "\u0d85\u0db4\u0dca\u200d\u0dbb\u0dda\u0dbd\u0dca",
            "\u0db8\u0dd0\u0dba\u0dd2",
            "\u0da2\u0dd6\u0db1\u0dd2",
            "\u0da2\u0dd6\u0dbd\u0dd2",
            "\u0d85\u0d9c\u0ddd\u0dc3\u0dca\u0dad\u0dd4",
            "\u0dc3\u0dd0\u0db4\u0dca\u0dad\u0dd0\u0db8\u0dca\u0db6\u0dbb\u0dca",
            "\u0d94\u0d9a\u0dca\u0dad\u0ddd\u0db6\u0dbb\u0dca",
            "\u0db1\u0ddc\u0dc0\u0dd0\u0db8\u0dca\u0db6\u0dbb\u0dca",
            "\u0daf\u0dd9\u0dc3\u0dd0\u0db8\u0dca\u0db6\u0dbb\u0dca",
            "",
        };
        final String[] metaValue_MonthNarrows = new String[] {
            "\u0da2",
            "\u0db4\u0dd9",
            "\u0db8\u0dcf",
            "\u0d85",
            "\u0db8\u0dd0",
            "\u0da2\u0dd6",
            "\u0da2\u0dd6",
            "\u0d85",
            "\u0dc3\u0dd0",
            "\u0d94",
            "\u0db1\u0dd9",
            "\u0daf\u0dd9",
            "",
        };
        final String[] metaValue_generic_DayNames = new String[] {
            "\u0d89\u0dbb\u0dd2\u0daf\u0dcf",
            "\u0dc3\u0db3\u0dd4\u0daf\u0dcf",
            "\u0d85\u0d9f\u0dc4\u0dbb\u0dd4\u0dc0\u0dcf\u0daf\u0dcf",
            "\u0db6\u0daf\u0dcf\u0daf\u0dcf",
            "\u0db6\u0dca\u200d\u0dbb\u0dc4\u0dc3\u0dca\u0db4\u0dad\u0dd2\u0db1\u0dca\u0daf\u0dcf",
            "\u0dc3\u0dd2\u0d9a\u0dd4\u0dbb\u0dcf\u0daf\u0dcf",
            "\u0dc3\u0dd9\u0db1\u0dc3\u0dd4\u0dbb\u0dcf\u0daf\u0dcf",
        };
        final String[] metaValue_generic_DayAbbreviations = new String[] {
            "\u0d89\u0dbb\u0dd2\u0daf\u0dcf",
            "\u0dc3\u0db3\u0dd4\u0daf\u0dcf",
            "\u0d85\u0d9f\u0dc4",
            "\u0db6\u0daf\u0dcf\u0daf\u0dcf",
            "\u0db6\u0dca\u200d\u0dbb\u0dc4\u0dc3\u0dca",
            "\u0dc3\u0dd2\u0d9a\u0dd4",
            "\u0dc3\u0dd9\u0db1",
        };
        final String[] metaValue_generic_DayNarrows = new String[] {
            "\u0d89",
            "\u0dc3",
            "\u0d85",
            "\u0db6",
            "\u0db6\u0dca\u200d\u0dbb",
            "\u0dc3\u0dd2",
            "\u0dc3\u0dd9",
        };
        final String[] metaValue_generic_QuarterNames = new String[] {
            "1 \u0dc0\u0db1 \u0d9a\u0dcf\u0dbb\u0dca\u0dad\u0dd4\u0dc0",
            "2 \u0dc0\u0db1 \u0d9a\u0dcf\u0dbb\u0dca\u0dad\u0dd4\u0dc0",
            "3 \u0dc0\u0db1 \u0d9a\u0dcf\u0dbb\u0dca\u0dad\u0dd4\u0dc0",
            "4 \u0dc0\u0db1 \u0d9a\u0dcf\u0dbb\u0dca\u0dad\u0dd4\u0dc0",
        };
        final String[] metaValue_generic_QuarterAbbreviations = new String[] {
            "\u0d9a\u0dcf\u0dbb\u0dca:1",
            "\u0d9a\u0dcf\u0dbb\u0dca:2",
            "\u0d9a\u0dcf\u0dbb\u0dca:3",
            "\u0d9a\u0dcf\u0dbb\u0dca:4",
        };
        final String[] metaValue_generic_abbreviated_AmPmMarkers = new String[] {
            "\u0db4\u0dd9.\u0dc0.",
            "\u0db4.\u0dc0.",
            "\u0db8\u0dd0\u0daf\u0dd2\u0dba\u0db8",
            "\u0db8\u0db0\u0dca\u200d\u0dba\u0dcf\u0dc4\u0dca\u0db1\u0dba",
            "\u0db4\u0dcf\u0db1\u0dca\u0daf\u0dbb",
            "\u0d8b\u0daf\u0dda",
            "\u0daf\u0dc0\u0dbd\u0dca",
            "",
            "\u0dc4\u0dc0\u0dc3",
            "",
            "\u0dbb\u0dd1",
            "\u0db8\u0dd0\u0daf\u0dd2\u0dba\u0db8\u0da7 \u0db4\u0dc3\u0dd4",
        };
        final String[] metaValue_generic_narrow_AmPmMarkers = new String[] {
            "\u0db4\u0dd9",
            "\u0db4",
            "\u0db8\u0dd0",
            "\u0db8",
            "\u0db4\u0dcf",
            "\u0d8b",
            "\u0daf",
            "",
            "\u0dc4",
            "",
            "\u0dbb\u0dd1",
            "\u0db8\u0dd0",
        };
        final String[] metaValue_java_time_long_Eras = new String[] {
            "\u0d9a\u0dca\u200d\u0dbb\u0dd2\u0dc3\u0dca\u0dad\u0dd4 \u0db4\u0dd6\u0dbb\u0dca\u0dc0",
            "\u0d9a\u0dca\u200d\u0dbb\u0dd2\u0dc3\u0dca\u0dad\u0dd4 \u0dc0\u0dbb\u0dca\u0dc2",
        };
        final String[] metaValue_java_time_Eras = new String[] {
            "\u0d9a\u0dca\u200d\u0dbb\u0dd2.\u0db4\u0dd6.",
            "\u0d9a\u0dca\u200d\u0dbb\u0dd2.\u0dc0.",
        };
        final String[] metaValue_generic_TimePatterns = new String[] {
            "HH.mm.ss zzzz",
            "HH.mm.ss z",
            "HH.mm.ss",
            "HH.mm",
        };
        final String[] metaValue_generic_DateTimePatterns = new String[] {
            "{1} {0}",
            "{1} {0}",
            "{1} {0}",
            "{1} {0}",
        };
        final String metaValue_generic_DateFormatItem_hm = "a h.mm";
        final String metaValue_generic_DateFormatItem_Ehm = "E a h.mm";
        final String metaValue_generic_DateFormatItem_hms = "a h.mm.ss";
        final String metaValue_generic_DateFormatItem_hmv = "h.mm a v";
        final String metaValue_generic_DateFormatItem_hmsv = "h.mm.ss a v";
        final String metaValue_generic_DateFormatItem_yw = "Y \u0dc3\u0dad\u0dd2\u0dba w";
        final String metaValue_generic_DateFormatItem_Hm = "HH.mm";
        final String metaValue_generic_DateFormatItem_EHm = "E HH.mm";
        final String metaValue_generic_DateFormatItem_MMMEd = "MMM d E";
        final String metaValue_generic_DateFormatItem_h = "a h";
        final String metaValue_generic_DateFormatItem_Hmsv = "HH.mm.ss v";
        final String metaValue_generic_DateFormatItem_MMMMW = "MMMM \u0dc3\u0dad\u0dd2\u0dba W";
        final String metaValue_generic_DateFormatItem_Ehms = "E a h.mm.ss";
        final String metaValue_generic_DateFormatItem_yMEd = "y-M-d, E";
        final String metaValue_generic_DateFormatItem_ms = "mm.ss";
        final String metaValue_generic_DateFormatItem_EHms = "E HH.mm.ss";
        final String metaValue_generic_DateFormatItem_Md = "M-d";
        final String metaValue_generic_DateFormatItem_Ed = "d E";
        final String metaValue_generic_DateFormatItem_MEd = "M-d, E";
        final String metaValue_generic_DateFormatItem_Hms = "HH.mm.ss";
        final String metaValue_generic_DateFormatItem_Hmv = "HH.mm v";
        final String metaValue_generic_DateFormatItem_yMd = "y-M-d";
        final String metaValue_generic_DateFormatItem_yM = "y-M";
        final String[] metaValue_MonthAbbreviations = new String[] {
            "\u0da2\u0db1",
            "\u0db4\u0dd9\u0db6",
            "\u0db8\u0dcf\u0dbb\u0dca\u0dad\u0dd4",
            "\u0d85\u0db4\u0dca\u200d\u0dbb\u0dda\u0dbd\u0dca",
            "\u0db8\u0dd0\u0dba\u0dd2",
            "\u0da2\u0dd6\u0db1\u0dd2",
            "\u0da2\u0dd6\u0dbd\u0dd2",
            "\u0d85\u0d9c\u0ddd",
            "\u0dc3\u0dd0\u0db4\u0dca",
            "\u0d94\u0d9a\u0dca",
            "\u0db1\u0ddc\u0dc0\u0dd0",
            "\u0daf\u0dd9\u0dc3\u0dd0",
            "",
        };
        final String[] metaValue_java_time_generic_DatePatterns = new String[] {
            "G y MMMM d, EEEE",
            "G y MMMM d",
            "G y MMM d",
            "GGGGG y-MM-dd",
        };
        final String[] metaValue_generic_DatePatterns = new String[] {
            "GGGG y MMMM d, EEEE",
            "GGGG y MMMM d",
            "GGGG y MMM d",
            "G y-MM-dd",
        };
        final String metaValue_calendarname_gregorian = "\u0d9c\u0dca\u200d\u0dbb\u0dd9\u0d9c\u0dbb\u0dd2\u0dba\u0dcf\u0db1\u0dd4 \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba";
        final String[] metaValue_ListPatterns_standard = new String[] {
            "",
            "",
            "{0}, \u0dc3\u0dc4 {1}",
            "{0} \u0dc3\u0dc4 {1}",
            "",
        };
        final String[] metaValue_ListPatterns_or = new String[] {
            "",
            "",
            "{0}, \u0dc4\u0ddd {1}",
            "{0} \u0dc4\u0ddd {1}",
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
            { "generic.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "generic.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "generic.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "generic.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "generic.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "generic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "generic.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "generic.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "generic.DateFormatItem.EBhm", "E B h:mm" },
            { "generic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "generic.DateFormatItem.Bh", "B h" },
            { "generic.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "generic.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "generic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "generic.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "generic.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "generic.DateFormatItem.Bhms", "B h:mm:ss" },
            { "generic.DateFormatItem.yyyyMEd", "G y-M-d, E" },
            { "generic.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "generic.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "generic.DateFormatItem.yyyyMd", "G y-M-d" },
            { "generic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "generic.DateFormatItem.EBhms", "E B h:mm:ss" },
            { "generic.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "generic.DateFormatItem.Bhm", "B h:mm" },
            { "generic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "generic.DateFormatItem.yyyyM", "G y-M" },
            { "generic.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "generic.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "generic.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "generic.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "MonthNames", metaValue_MonthNames },
            { "standalone.MonthNames", metaValue_MonthNames },
            { "MonthAbbreviations", metaValue_MonthAbbreviations },
            { "standalone.MonthAbbreviations",
                new String[] {
                    "\u0da2\u0db1",
                    "\u0db4\u0dd9\u0db6",
                    "\u0db8\u0dcf\u0dbb\u0dca",
                    "\u0d85\u0db4\u0dca\u200d\u0dbb\u0dda\u0dbd\u0dca",
                    "\u0db8\u0dd0\u0dba\u0dd2",
                    "\u0da2\u0dd6\u0db1\u0dd2",
                    "\u0da2\u0dd6\u0dbd\u0dd2",
                    "\u0d85\u0d9c\u0ddd",
                    "\u0dc3\u0dd0\u0db4\u0dca",
                    "\u0d94\u0d9a\u0dca",
                    "\u0db1\u0ddc\u0dc0\u0dd0",
                    "\u0daf\u0dd9\u0dc3\u0dd0",
                    "",
                }
            },
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
            { "field.era", "\u0dba\u0dd4\u0d9c\u0dba" },
            { "field.year", "\u0dc0\u0dbb\u0dca\u0dc2\u0dba" },
            { "field.month", "\u0db8\u0dcf\u0dc3\u0dba" },
            { "field.week", "\u0dc3\u0dad\u0dd2\u0dba" },
            { "field.weekday", "\u0dc3\u0dad\u0dd2\u0dba\u0dda \u0daf\u0dd2\u0db1\u0dba" },
            { "field.dayperiod", "\u0db4\u0dd9.\u0dc0/\u0db4.\u0dc0" },
            { "field.hour", "\u0db4\u0dd0\u0dba" },
            { "timezone.hourFormat", "+HH.mm;-HH.mm" },
            { "timezone.gmtFormat", "\u0d9c\u0dca\u200d\u0dbb\u0dd2\u0db8\u0dc0\u0dda{0}" },
            { "timezone.gmtZeroFormat", "\u0d9c\u0dca\u200d\u0dbb\u0dd2\u0db8\u0dc0\u0dda" },
            { "timezone.regionFormat", "{0} \u0dc0\u0dda\u0dbd\u0dcf\u0dc0" },
            { "timezone.regionFormat.daylight", "{0} \u0daf\u0dd2\u0dc0\u0dcf\u0d86\u0dbd\u0ddd\u0d9a \u0dc0\u0dda\u0dbd\u0dcf\u0dc0" },
            { "timezone.regionFormat.standard", "{0} \u0dc3\u0db8\u0dca\u0db8\u0dad \u0dc0\u0dda\u0dbd\u0dcf\u0dc0" },
            { "field.minute", "\u0db8\u0dd2\u0db1\u0dd2\u0dad\u0dca\u0dad\u0dd4\u0dc0" },
            { "field.second", "\u0dad\u0dad\u0dca\u0db4\u0dbb\u0dba" },
            { "field.zone", "\u0dc0\u0dda\u0dbd\u0dcf \u0d9a\u0dbd\u0dcf\u0db4\u0dba" },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "y MMMM d, EEEE",
                    "y MMMM d",
                    "y MMM d",
                    "y-MM-dd",
                }
            },
            { "DateTimePatterns", metaValue_generic_DateTimePatterns },
            { "PluralRules", "one:n = 0,1 or i = 0 and f = 1" },
            { "DayPeriodRules", "midnight:00:00;noon:12:00;night1:18:00-24:00;morning2:06:00-12:00;afternoon1:12:00-14:00;morning1:01:00-06:00;evening1:14:00-18:00;night2:00:00-01:00" },
            { "DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
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
            { "buddhist.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "buddhist.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "buddhist.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "buddhist.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "buddhist.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "buddhist.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "buddhist.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "buddhist.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "buddhist.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "buddhist.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "buddhist.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "buddhist.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "buddhist.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "buddhist.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "buddhist.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "buddhist.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "buddhist.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "buddhist.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "buddhist.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "buddhist.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "buddhist.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "buddhist.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "buddhist.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
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
                    "GGGG yyyy MMMM d, EEEE",
                    "GGGG yyyy MMMM d",
                    "GGGG y MMM d",
                    "G y-MM-dd",
                }
            },
            { "japanese.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "japanese.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "japanese.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "japanese.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "japanese.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "japanese.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "japanese.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "japanese.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "japanese.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "japanese.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "japanese.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "japanese.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "japanese.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "japanese.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "japanese.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "japanese.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "japanese.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "japanese.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "japanese.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "japanese.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "japanese.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "japanese.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "japanese.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "roc.MonthNames", metaValue_MonthNames },
            { "roc.MonthAbbreviations", metaValue_MonthAbbreviations },
            { "roc.DayNames", metaValue_generic_DayNames },
            { "roc.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "roc.DayNarrows", metaValue_generic_DayNarrows },
            { "roc.QuarterNames", metaValue_generic_QuarterNames },
            { "roc.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "roc.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "roc.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.roc.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "roc.DatePatterns", metaValue_generic_DatePatterns },
            { "roc.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "roc.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "roc.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "roc.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "roc.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "roc.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "roc.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "roc.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "roc.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "roc.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "roc.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "roc.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "roc.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "roc.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "roc.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "roc.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "roc.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "roc.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "roc.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "roc.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "roc.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "roc.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "roc.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic.DayNames", metaValue_generic_DayNames },
            { "islamic.DayAbbreviations", metaValue_generic_DayAbbreviations },
            { "islamic.QuarterNames", metaValue_generic_QuarterNames },
            { "islamic.QuarterAbbreviations", metaValue_generic_QuarterAbbreviations },
            { "islamic.narrow.AmPmMarkers", metaValue_generic_narrow_AmPmMarkers },
            { "islamic.abbreviated.AmPmMarkers", metaValue_generic_abbreviated_AmPmMarkers },
            { "islamic.TimePatterns", metaValue_generic_TimePatterns },
            { "java.time.islamic.DatePatterns", metaValue_java_time_generic_DatePatterns },
            { "islamic.DatePatterns", metaValue_generic_DatePatterns },
            { "islamic.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "islamic.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "islamic.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "islamic.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "islamic.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "islamic.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "islamic.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "islamic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "islamic.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "islamic.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "islamic.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "islamic.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "islamic.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "islamic.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "islamic.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic-civil.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "islamic-civil.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "islamic-civil.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "islamic-civil.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "islamic-civil.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "islamic-civil.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-civil.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "islamic-civil.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "islamic-civil.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-civil.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "islamic-civil.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "islamic-civil.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-civil.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "islamic-civil.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic-civil.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "islamic-civil.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "islamic-civil.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-civil.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic-civil.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-civil.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "islamic-civil.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "islamic-civil.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic-civil.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "islamic-umalqura.DateFormatItem.hm", metaValue_generic_DateFormatItem_hm },
            { "islamic-umalqura.DateFormatItem.Ehm", metaValue_generic_DateFormatItem_Ehm },
            { "islamic-umalqura.DateFormatItem.hms", metaValue_generic_DateFormatItem_hms },
            { "islamic-umalqura.DateFormatItem.hmv", metaValue_generic_DateFormatItem_hmv },
            { "islamic-umalqura.DateFormatItem.hmsv", metaValue_generic_DateFormatItem_hmsv },
            { "islamic-umalqura.DateFormatItem.yw", metaValue_generic_DateFormatItem_yw },
            { "islamic-umalqura.DateFormatItem.Hm", metaValue_generic_DateFormatItem_Hm },
            { "islamic-umalqura.DateFormatItem.EHm", metaValue_generic_DateFormatItem_EHm },
            { "islamic-umalqura.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-umalqura.DateFormatItem.h", metaValue_generic_DateFormatItem_h },
            { "islamic-umalqura.DateFormatItem.Hmsv", metaValue_generic_DateFormatItem_Hmsv },
            { "islamic-umalqura.DateFormatItem.MMMMW", metaValue_generic_DateFormatItem_MMMMW },
            { "islamic-umalqura.DateFormatItem.Ehms", metaValue_generic_DateFormatItem_Ehms },
            { "islamic-umalqura.DateFormatItem.yMEd", metaValue_generic_DateFormatItem_yMEd },
            { "islamic-umalqura.DateFormatItem.ms", metaValue_generic_DateFormatItem_ms },
            { "islamic-umalqura.DateFormatItem.EHms", metaValue_generic_DateFormatItem_EHms },
            { "islamic-umalqura.DateFormatItem.Md", metaValue_generic_DateFormatItem_Md },
            { "islamic-umalqura.DateFormatItem.Ed", metaValue_generic_DateFormatItem_Ed },
            { "islamic-umalqura.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-umalqura.DateFormatItem.Hms", metaValue_generic_DateFormatItem_Hms },
            { "islamic-umalqura.DateFormatItem.Hmv", metaValue_generic_DateFormatItem_Hmv },
            { "islamic-umalqura.DateFormatItem.yMd", metaValue_generic_DateFormatItem_yMd },
            { "islamic-umalqura.DateFormatItem.yM", metaValue_generic_DateFormatItem_yM },
            { "calendarname.islamic-civil", "\u0dc4\u0dd2\u0da2\u0dca\u0dbb\u0dd2 \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba (\u0dc0\u0d9c\u0dd4\u0dc0, \u0dc3\u0dd2\u0dc0\u0dd2\u0dbd\u0dca \u0dc0\u0dd3\u0dbb \u0d9a\u0dcf\u0dc0\u0dca\u200d\u0dba\u0dba)" },
            { "calendarname.islamic", "\u0dc4\u0dd2\u0da2\u0dca\u0dbb\u0dd2 \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba" },
            { "calendarname.buddhist", "\u0db6\u0ddc\u0daf\u0dd4 \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba" },
            { "calendarname.japanese", "\u0da2\u0db4\u0db1\u0dca \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba" },
            { "calendarname.roc", "\u0db8\u0dd2\u0db1\u0dca\u0d9c\u0dcf \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba" },
            { "calendarname.islamic-umalqura", "\u0dc4\u0dd2\u0da2\u0dca\u0dbb\u0dd2 \u0daf\u0dd2\u0db1 \u0daf\u0dbb\u0dca\u0dc1\u0db1\u0dba (\u0d8b\u0db8\u0dca \u0d85\u0dbd\u0dca \u0d9a\u0dd4\u0dbb\u0dcf)" },
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
                    "#,##0.###",
                    "\u00a4#,##0.00",
                    "#,##0%",
                    "\u00a4#,##0.00;(\u00a4#,##0.00)",
                }
            },
            { "short.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:\u0daf0 other:\u0daf0}",
                    "{one:\u0daf00 other:\u0daf00}",
                    "{one:\u0daf000 other:\u0daf000}",
                    "{one:\u0db8\u0dd20 other:\u0db8\u0dd20}",
                    "{one:\u0db8\u0dd200 other:\u0db8\u0dd200}",
                    "{one:\u0db8\u0dd2000 other:\u0db8\u0dd2000}",
                    "{one:\u0db6\u0dd20 other:\u0db6\u0dd20}",
                    "{one:\u0db6\u0dd200 other:\u0db6\u0dd200}",
                    "{one:\u0db6\u0dd2000 other:\u0db6\u0dd2000}",
                    "{one:\u0da7\u0dca\u200d\u0dbb\u0dd20 other:\u0da7\u0dca\u200d\u0dbb\u0dd20}",
                    "{one:\u0da7\u0dca\u200d\u0dbb\u0dd200 other:\u0da7\u0dca\u200d\u0dbb\u0dd200}",
                    "{one:\u0da7\u0dca\u200d\u0dbb\u0dd2000 other:\u0da7\u0dca\u200d\u0dbb\u0dd2000}",
                }
            },
            { "long.CompactNumberPatterns",
                new String[] {
                    "",
                    "",
                    "",
                    "{one:\u0daf\u0dc4\u0dc3' '0 other:\u0daf\u0dc4\u0dc3' '0}",
                    "{one:\u0daf\u0dc4\u0dc3' '00 other:\u0daf\u0dc4\u0dc3' '00}",
                    "{one:\u0daf\u0dc4\u0dc3' '000 other:\u0daf\u0dc4\u0dc3' '000}",
                    "{one:\u0db8\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '0 other:\u0db8\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '0}",
                    "{one:\u0db8\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '00 other:\u0db8\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '00}",
                    "{one:\u0db8\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '000 other:\u0db8\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '000}",
                    "{one:\u0db6\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '0 other:\u0db6\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '0}",
                    "{one:\u0db6\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '00 other:\u0db6\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '00}",
                    "{one:\u0db6\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '000 other:\u0db6\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '000}",
                    "{one:\u0da7\u0dca\u200d\u0dbb\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '0 other:\u0da7\u0dca\u200d\u0dbb\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '0}",
                    "{one:\u0da7\u0dca\u200d\u0dbb\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '00 other:\u0da7\u0dca\u200d\u0dbb\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '00}",
                    "{one:\u0da7\u0dca\u200d\u0dbb\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '000 other:\u0da7\u0dca\u200d\u0dbb\u0dd2\u0dbd\u0dd2\u0dba\u0db1' '000}",
                }
            },
            { "ListPatterns_standard", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-short", metaValue_ListPatterns_standard },
            { "ListPatterns_standard-narrow", metaValue_ListPatterns_standard },
            { "ListPatterns_or", metaValue_ListPatterns_or },
            { "ListPatterns_or-short", metaValue_ListPatterns_or },
            { "ListPatterns_or-narrow", metaValue_ListPatterns_or },
            { "ListPatterns_unit", metaValue_ListPatterns_standard },
            { "ListPatterns_unit-short", metaValue_ListPatterns_standard },
            { "ListPatterns_unit-narrow", metaValue_ListPatterns_standard },
        };
        return data;
    }
}
