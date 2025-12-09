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

package sun.util.resources.cldr.ext;

import sun.util.resources.TimeZoneNamesBundle;

public class TimeZoneNames_hi_Latn extends TimeZoneNamesBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] Reunion = new String[] {
               "Reunion Time",
               "",
               "",
               "",
               "",
               "",
            };
        final String[] EMPTY_ZONE = new String[] {
               "",
               "",
               "",
               "",
               "",
               "",
            };
        final String[] America_Central = new String[] {
               "North America Central Standard Time",
               "\u2205\u2205\u2205",
               "North America Central Daylight Time",
               "\u2205\u2205\u2205",
               "North America Central Time",
               "\u2205\u2205\u2205",
            };
        final String[] America_Eastern = new String[] {
               "North America Eastern Standard Time",
               "\u2205\u2205\u2205",
               "North America Eastern Daylight Time",
               "\u2205\u2205\u2205",
               "North America Eastern Time",
               "\u2205\u2205\u2205",
            };
        final String[] America_Pacific = new String[] {
               "North America Pacific Standard Time",
               "\u2205\u2205\u2205",
               "North America Pacific Daylight Time",
               "\u2205\u2205\u2205",
               "North America Pacific Time",
               "\u2205\u2205\u2205",
            };
        final String[] Pierre_Miquelon = new String[] {
               "St. Pierre & Miquelon Standard Time",
               "",
               "St. Pierre & Miquelon Daylight Time",
               "",
               "St. Pierre & Miquelon Time",
               "",
            };
        final String[] America_Mountain = new String[] {
               "North America Mountain Standard Time",
               "\u2205\u2205\u2205",
               "North America Mountain Daylight Time",
               "\u2205\u2205\u2205",
               "North America Mountain Time",
               "\u2205\u2205\u2205",
            };
        final Object[][] data = new Object[][] {
            { "America/Los_Angeles", America_Pacific },
            { "America/Denver", America_Mountain },
            { "America/Phoenix", America_Mountain },
            { "America/Chicago", America_Central },
            { "America/New_York", America_Eastern },
            { "America/Indianapolis", America_Eastern },
            { "Pacific/Honolulu",
                new String[] {
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "HST",
                }
            },
            { "UTC", EMPTY_ZONE },
            { "CST", America_Central },
            { "EST", America_Eastern },
            { "HST",
                new String[] {
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "HST",
                }
            },
            { "IET", America_Eastern },
            { "MST", America_Mountain },
            { "PNT", America_Mountain },
            { "PST", America_Pacific },
            { "Navajo", America_Mountain },
            { "CST6CDT", America_Central },
            { "EST5EDT", America_Eastern },
            { "Jamaica", America_Eastern },
            { "MST7MDT", America_Mountain },
            { "PST8PDT", America_Pacific },
            { "US/Hawaii",
                new String[] {
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "HST",
                }
            },
            { "Asia/Dubai", Reunion },
            { "US/Arizona", America_Mountain },
            { "US/Central", America_Central },
            { "US/Eastern", America_Eastern },
            { "US/Pacific", America_Pacific },
            { "US/Michigan", America_Eastern },
            { "US/Mountain", America_Mountain },
            { "SystemV/CST6", America_Central },
            { "SystemV/EST5", America_Eastern },
            { "SystemV/MST7", America_Mountain },
            { "America/Boise", America_Mountain },
            { "SystemV/HST10",
                new String[] {
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "HST",
                }
            },
            { "America/Belize", America_Central },
            { "America/Cancun", America_Eastern },
            { "America/Cayman", America_Eastern },
            { "America/Inuvik", America_Mountain },
            { "America/Merida", America_Central },
            { "America/Nassau", America_Eastern },
            { "America/Panama", America_Eastern },
            { "America/Regina", America_Central },
            { "Canada/Central", America_Central },
            { "Canada/Eastern", America_Eastern },
            { "Canada/Pacific", America_Pacific },
            { "Indian/Reunion", Reunion },
            { "Mexico/General", America_Central },
            { "America/Creston", America_Mountain },
            { "America/Detroit", America_Eastern },
            { "America/Iqaluit", America_Eastern },
            { "America/Jamaica", America_Eastern },
            { "America/Knox_IN", America_Central },
            { "America/Managua", America_Central },
            { "America/Nipigon", America_Eastern },
            { "America/Ojinaga", America_Central },
            { "America/Tijuana", America_Pacific },
            { "America/Toronto", America_Eastern },
            { "Canada/Mountain", America_Mountain },
            { "SystemV/CST6CDT", America_Central },
            { "SystemV/EST5EDT", America_Eastern },
            { "SystemV/MST7MDT", America_Mountain },
            { "SystemV/PST8PDT", America_Pacific },
            { "America/Atikokan", America_Eastern },
            { "America/Edmonton", America_Mountain },
            { "America/Ensenada", America_Pacific },
            { "America/Miquelon", Pierre_Miquelon },
            { "America/Montreal", America_Eastern },
            { "America/Resolute", America_Central },
            { "America/Shiprock", America_Mountain },
            { "America/Winnipeg", America_Central },
            { "Mexico/BajaNorte", America_Pacific },
            { "Pacific/Johnston",
                new String[] {
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "\u2205\u2205\u2205",
                    "",
                    "HST",
                }
            },
            { "America/Chihuahua", America_Central },
            { "America/Guatemala", America_Central },
            { "America/Matamoros", America_Central },
            { "America/Menominee", America_Central },
            { "America/Monterrey", America_Central },
            { "America/Vancouver", America_Pacific },
            { "US/Indiana-Starke", America_Central },
            { "America/Costa_Rica", America_Central },
            { "America/Grand_Turk", America_Eastern },
            { "America/Louisville", America_Eastern },
            { "America/El_Salvador", America_Central },
            { "America/Fort_Nelson", America_Mountain },
            { "America/Mexico_City", America_Central },
            { "America/Pangnirtung", America_Eastern },
            { "America/Rainy_River", America_Central },
            { "America/Tegucigalpa", America_Central },
            { "America/Thunder_Bay", America_Eastern },
            { "America/Yellowknife", America_Mountain },
            { "Canada/Saskatchewan", America_Central },
            { "America/Dawson_Creek", America_Mountain },
            { "America/Indiana/Knox", America_Central },
            { "America/Rankin_Inlet", America_Central },
            { "America/Santa_Isabel", America_Pacific },
            { "America/Cambridge_Bay", America_Mountain },
            { "America/Ciudad_Juarez", America_Mountain },
            { "America/Coral_Harbour", America_Eastern },
            { "America/Indiana/Vevay", America_Eastern },
            { "America/Swift_Current", America_Central },
            { "America/Bahia_Banderas", America_Central },
            { "America/Port-au-Prince", America_Eastern },
            { "America/Indiana/Marengo", America_Eastern },
            { "America/Indiana/Winamac", America_Eastern },
            { "America/Indiana/Tell_City", America_Central },
            { "America/Indiana/Vincennes", America_Eastern },
            { "America/Indiana/Petersburg", America_Eastern },
            { "timezone.excity.Asia/Aqtau", "Aqtau" },
            { "timezone.excity.Asia/Macau", "Macau" },
            { "America/Kentucky/Louisville", America_Eastern },
            { "America/Kentucky/Monticello", America_Eastern },
            { "America/North_Dakota/Beulah", America_Central },
            { "America/North_Dakota/Center", America_Central },
            { "timezone.excity.Asia/Saigon", "Saigon" },
            { "America/Indiana/Indianapolis", America_Eastern },
            { "timezone.excity.Pacific/Truk", "Truk" },
            { "timezone.excity.Africa/Asmera", "Asmera" },
            { "timezone.excity.Asia/Qostanay", "Qostanay" },
            { "America/North_Dakota/New_Salem", America_Central },
            { "timezone.excity.America/Cancun", "Cancun" },
            { "timezone.excity.America/Merida", "Merida" },
            { "timezone.excity.Indian/Reunion", "Reunion" },
            { "timezone.excity.Pacific/Ponape", "Ponape" },
            { "timezone.excity.Africa/Sao_Tome", "Sao Tome" },
            { "timezone.excity.America/Curacao", "Curacao" },
            { "timezone.excity.Atlantic/Faeroe", "Faeroe" },
            { "timezone.excity.America/Asuncion", "Asuncion" },
            { "timezone.excity.Pacific/Honolulu", "Honolulu" },
            { "timezone.excity.Africa/Porto-Novo", "Porto Novo" },
            { "timezone.excity.America/Blanc-Sablon", "Blanc Sablon" },
            { "timezone.excity.America/Ciudad_Juarez", "Ciudad Juarez" },
            { "timezone.excity.America/Coral_Harbour", "Coral Harbour" },
            { "timezone.excity.America/St_Barthelemy", "St Barthelemy" },
            { "timezone.excity.America/Bahia_Banderas", "Bahia Banderas" },
            { "timezone.excity.Antarctica/DumontDUrville", "DumontDUrville" },
        };
        return data;
    }
}
