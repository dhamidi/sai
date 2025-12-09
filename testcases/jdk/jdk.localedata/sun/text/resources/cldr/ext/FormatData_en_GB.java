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

public class FormatData_en_GB extends ListResourceBundle {
    @Override
    protected final Object[][] getContents() {
        final String metaValue_generic_DateFormatItem_GyMMMEEEEd = "EEEE, d MMM y G";
        final String[] metaValue_generic_TimePatterns = new String[] {
            "HH:mm:ss zzzz",
            "HH:mm:ss z",
            "HH:mm:ss",
            "HH:mm",
        };
        final String metaValue_generic_DateFormatItem_MMMMEEEEd = "EEEE d MMMM";
        final String metaValue_generic_DateFormatItem_yMMMMEEEEd = "EEEE, d MMMM y";
        final String metaValue_generic_DateFormatItem_yMMMEEEEd = "EEEE, d MMM y";
        final String metaValue_generic_DateFormatItem_MEd = "E dd/MM";
        final String metaValue_generic_DateFormatItem_MMMEd = "E d MMM";
        final String metaValue_generic_DateFormatItem_MMMEEEEd = "EEEE d MMM";
        final String metaValue_generic_DateFormatItem_GyMd = "dd/MM/y G";
        final String metaValue_generic_DateFormatItem_yyyyMMMMEEEEd = "EEEE, d MMMM y G";
        final Object[][] data = new Object[][] {
            { "generic.TimePatterns", metaValue_generic_TimePatterns },
            { "generic.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "generic.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "generic.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "generic.DateFormatItem.yyyyMMMMEEEEd", metaValue_generic_DateFormatItem_yyyyMMMMEEEEd },
            { "generic.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "generic.DateFormatItem.yyyyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "generic.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "generic.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "generic.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "generic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "TimePatterns", metaValue_generic_TimePatterns },
            { "DatePatterns",
                new String[] {
                    "EEEE, d MMMM y",
                    "d MMMM y",
                    "d MMM y",
                    "dd/MM/y",
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
            { "DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "buddhist.TimePatterns", metaValue_generic_TimePatterns },
            { "buddhist.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "buddhist.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "buddhist.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "buddhist.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "buddhist.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "buddhist.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "buddhist.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "buddhist.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "japanese.TimePatterns", metaValue_generic_TimePatterns },
            { "japanese.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "japanese.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "japanese.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "japanese.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "japanese.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "japanese.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "japanese.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "japanese.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "roc.TimePatterns", metaValue_generic_TimePatterns },
            { "roc.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "roc.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "roc.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "roc.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "roc.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "roc.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "roc.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "roc.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic.TimePatterns", metaValue_generic_TimePatterns },
            { "islamic.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "islamic.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "islamic.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic.DateFormatItem.yyyyMMMMEEEEd", metaValue_generic_DateFormatItem_yyyyMMMMEEEEd },
            { "islamic.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "islamic.DateFormatItem.yyyyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic.DateFormatItem.Ed", "E d" },
            { "islamic.DateFormatItem.M", "LL" },
            { "islamic.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "islamic.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic.DateFormatItem.yyyyM", "MM/y GGGGG" },
            { "islamic-civil.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "islamic-civil.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "islamic-civil.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic-civil.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "islamic-civil.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-civil.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-civil.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "islamic-civil.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
            { "islamic-umalqura.DateFormatItem.MMMMEEEEd", metaValue_generic_DateFormatItem_MMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.yMMMMEEEEd", metaValue_generic_DateFormatItem_yMMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.GyMMMEEEEd", metaValue_generic_DateFormatItem_GyMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.yMMMEEEEd", metaValue_generic_DateFormatItem_yMMMEEEEd },
            { "islamic-umalqura.DateFormatItem.MEd", metaValue_generic_DateFormatItem_MEd },
            { "islamic-umalqura.DateFormatItem.MMMEd", metaValue_generic_DateFormatItem_MMMEd },
            { "islamic-umalqura.DateFormatItem.MMMEEEEd", metaValue_generic_DateFormatItem_MMMEEEEd },
            { "islamic-umalqura.DateFormatItem.GyMd", metaValue_generic_DateFormatItem_GyMd },
        };
        return data;
    }
}
