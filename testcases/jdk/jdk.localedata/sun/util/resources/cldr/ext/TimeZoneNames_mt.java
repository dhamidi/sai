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

public class TimeZoneNames_mt extends TimeZoneNamesBundle {
    @Override
    protected final Object[][] getContents() {
        final String[] GMT = new String[] {
               "Greenwich Mean Time",
               "GMT",
               "",
               "BST",
               "",
               "GMT",
            };
        final String[] EMPTY_ZONE = new String[] {
               "",
               "",
               "",
               "",
               "",
               "",
            };
        final String[] Europe_Central = new String[] {
               "\u0126in \u010aentrali Ewropew Standard",
               "CET",
               "\u0126in \u010aentrali Ewropew tas-Sajf",
               "CEST",
               "\u0126in \u010aentrali Ewropew",
               "CET",
            };
        final String[] Europe_Eastern = new String[] {
               "",
               "EET",
               "",
               "EEST",
               "",
               "EET",
            };
        final String[] Europe_Western = new String[] {
               "",
               "WET",
               "",
               "WEST",
               "",
               "WET",
            };
        final Object[][] data = new Object[][] {
            { "Europe/Paris", Europe_Central },
            { "GMT", GMT },
            { "Europe/Bucharest", Europe_Eastern },
            { "UTC", EMPTY_ZONE },
            { "GB", GMT },
            { "ART", Europe_Eastern },
            { "CET", Europe_Central },
            { "ECT", Europe_Central },
            { "EET", Europe_Eastern },
            { "MET", Europe_Central },
            { "WET", Europe_Western },
            { "Eire", GMT },
            { "GMT0", GMT },
            { "Egypt", Europe_Eastern },
            { "Libya", Europe_Eastern },
            { "Poland", Europe_Central },
            { "Etc/GMT", GMT },
            { "GB-Eire", GMT },
            { "Iceland", GMT },
            { "Etc/GMT0", GMT },
            { "Portugal", Europe_Western },
            { "Asia/Gaza", Europe_Eastern },
            { "Etc/GMT+0", GMT },
            { "Etc/GMT-0", GMT },
            { "Greenwich", GMT },
            { "Africa/Lome", GMT },
            { "Asia/Beirut", Europe_Eastern },
            { "Asia/Hebron", Europe_Eastern },
            { "Europe/Kiev", Europe_Eastern },
            { "Europe/Kyiv", Europe_Eastern },
            { "Europe/Oslo", Europe_Central },
            { "Europe/Riga", Europe_Eastern },
            { "Europe/Rome", Europe_Central },
            { "Africa/Accra", GMT },
            { "Africa/Cairo", Europe_Eastern },
            { "Africa/Ceuta", Europe_Central },
            { "Africa/Dakar", GMT },
            { "Africa/Tunis", Europe_Central },
            { "Asia/Nicosia", Europe_Eastern },
            { "Europe/Malta", Europe_Central },
            { "Europe/Sofia", Europe_Eastern },
            { "Europe/Vaduz", Europe_Central },
            { "Africa/Bamako", GMT },
            { "Africa/Banjul", GMT },
            { "Africa/Bissau", GMT },
            { "Etc/Greenwich", GMT },
            { "Europe/Athens", Europe_Eastern },
            { "Europe/Berlin", Europe_Central },
            { "Europe/Dublin", GMT },
            { "Europe/Jersey", GMT },
            { "Europe/Lisbon", Europe_Western },
            { "Europe/London", GMT },
            { "Europe/Madrid", Europe_Central },
            { "Europe/Monaco", Europe_Central },
            { "Europe/Prague", Europe_Central },
            { "Europe/Skopje", Europe_Central },
            { "Europe/Tirane", Europe_Central },
            { "Europe/Vienna", Europe_Central },
            { "Europe/Warsaw", Europe_Central },
            { "Europe/Zagreb", Europe_Central },
            { "Europe/Zurich", Europe_Central },
            { "Africa/Abidjan", GMT },
            { "Africa/Algiers", Europe_Central },
            { "Africa/Conakry", GMT },
            { "Africa/Tripoli", Europe_Eastern },
            { "Atlantic/Faroe", Europe_Western },
            { "Europe/Andorra", Europe_Central },
            { "Europe/Belfast", GMT },
            { "Europe/Nicosia", Europe_Eastern },
            { "Europe/Tallinn", Europe_Eastern },
            { "Europe/Vatican", Europe_Central },
            { "Europe/Vilnius", Europe_Eastern },
            { "Africa/Freetown", GMT },
            { "Africa/Monrovia", GMT },
            { "Africa/Sao_Tome", GMT },
            { "Africa/Timbuktu", GMT },
            { "Atlantic/Canary", Europe_Western },
            { "Atlantic/Faeroe", Europe_Western },
            { "Europe/Belgrade", Europe_Central },
            { "Europe/Brussels", Europe_Central },
            { "Europe/Budapest", Europe_Central },
            { "Europe/Busingen", Europe_Central },
            { "Europe/Chisinau", Europe_Eastern },
            { "Europe/Guernsey", GMT },
            { "Europe/Helsinki", Europe_Eastern },
            { "Europe/Sarajevo", Europe_Central },
            { "Europe/Tiraspol", Europe_Eastern },
            { "Europe/Uzhgorod", Europe_Eastern },
            { "Antarctica/Troll", GMT },
            { "Atlantic/Madeira", Europe_Western },
            { "Europe/Amsterdam", Europe_Central },
            { "Europe/Gibraltar", Europe_Central },
            { "Europe/Ljubljana", Europe_Central },
            { "Europe/Mariehamn", Europe_Eastern },
            { "Europe/Podgorica", Europe_Central },
            { "Europe/Stockholm", Europe_Central },
            { "Africa/Nouakchott", GMT },
            { "Europe/Bratislava", Europe_Central },
            { "Europe/Copenhagen", Europe_Central },
            { "Europe/Luxembourg", Europe_Central },
            { "Europe/San_Marino", Europe_Central },
            { "Europe/Zaporozhye", Europe_Eastern },
            { "Africa/Ouagadougou", GMT },
            { "Atlantic/Jan_Mayen", Europe_Central },
            { "Atlantic/Reykjavik", GMT },
            { "Atlantic/St_Helena", GMT },
            { "Europe/Isle_of_Man", GMT },
            { "Europe/Kaliningrad", Europe_Eastern },
            { "Arctic/Longyearbyen", Europe_Central },
            { "America/Danmarkshavn", GMT },
            { "timezone.excity.Asia/Macau", "Macau" },
            { "timezone.excity.Asia/Beirut", "Bejrut" },
            { "timezone.excity.Asia/Kuwait", "il-Belt tal-Kuwajt" },
            { "timezone.excity.Etc/Unknown", "Belt Mhux Mag\u0127ruf" },
            { "timezone.excity.Europe/Rome", "Ruma" },
            { "timezone.excity.Asia/Nicosia", "Nikosija" },
            { "timezone.excity.Asia/Rangoon", "Rangoon" },
            { "timezone.excity.Europe/Malta", "Valletta" },
            { "timezone.excity.Europe/Paris", "Pari\u0121i" },
            { "timezone.excity.Europe/Sofia", "Sofija" },
            { "timezone.excity.Asia/Damascus", "Damasku" },
            { "timezone.excity.Europe/Athens", "Ateni" },
            { "timezone.excity.Europe/Lisbon", "Li\u017cbona" },
            { "timezone.excity.Europe/London", "Londra" },
            { "timezone.excity.Europe/Moscow", "Moska" },
            { "timezone.excity.Europe/Prague", "Praga" },
            { "timezone.excity.Europe/Tirane", "Tirana" },
            { "timezone.excity.Europe/Vienna", "Vjenna" },
            { "timezone.excity.Europe/Warsaw", "Varsavja" },
            { "timezone.excity.Africa/Algiers", "l-Al\u0121ier" },
            { "timezone.excity.Asia/Jerusalem", "\u0120erusalemm" },
            { "timezone.excity.Europe/Vatican", "il-Belt tal-Vatikan" },
            { "timezone.excity.America/Jamaica", "il-\u0120amajka" },
            { "timezone.excity.Europe/Belgrade", "Belgrad" },
            { "timezone.excity.Europe/Brussels", "Brussell" },
            { "timezone.excity.Indian/Maldives", "il-Maldivi" },
            { "timezone.excity.Europe/Gibraltar", "\u0120ibilt\u00e0" },
            { "timezone.excity.Europe/Stockholm", "Stokkolma" },
            { "timezone.excity.Pacific/Honolulu", "Honolulu" },
            { "timezone.excity.America/Guatemala", "il-Gwatemala" },
            { "timezone.excity.Europe/Luxembourg", "il-Lussemburgu" },
            { "timezone.excity.Pacific/Enderbury", "Enderbury" },
            { "timezone.excity.Atlantic/South_Georgia", "il-Georgia tan-Nofsinhar" },
        };
        return data;
    }
}
