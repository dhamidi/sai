/*
 * Copyright (c) 2019, 2025, Oracle and/or its affiliates. All rights reserved.
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
package java.lang.invoke;

import jdk.internal.foreign.AbstractMemorySegmentImpl;
import jdk.internal.vm.annotation.ForceInline;

import java.lang.foreign.MemorySegment;

import static java.lang.invoke.SegmentVarHandle.*;

// -- This file was mechanically generated: Do not edit! -- //

final class VarHandleSegmentAsBooleans {

    static VarForm selectForm(long alignmentMask, boolean constantOffset) {
        return constantOffset ? CONSTANT_OFFSET_FORM : VARIABLE_OFFSET_FORM;
    }

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsBooleans.class, MemorySegment.class, boolean.class, long.class);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsBooleans.class, MemorySegment.class, boolean.class, long.class, long.class);

    VarHandleSegmentAsBooleans() { throw new AssertionError(); }

    @ForceInline
    static boolean get(VarHandle ob, Object obb, long base) {
        return get(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static boolean get(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return SCOPED_MEMORY_ACCESS.getBoolean(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset));
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, boolean value) {
        set(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, long offset, boolean value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putBoolean(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                value);
    }

    @ForceInline
    static boolean convEndian(boolean big, boolean n) {
        return n;
    }

    @ForceInline
    static boolean getVolatile(VarHandle ob, Object obb, long base) {
        return getVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static boolean getVolatile(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getBooleanVolatile(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, boolean value) {
        setVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, long offset, boolean value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putBooleanVolatile(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static boolean getAcquire(VarHandle ob, Object obb, long base) {
        return getAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static boolean getAcquire(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getBooleanAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, boolean value) {
        setRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, long offset, boolean value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putBooleanRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static boolean getOpaque(VarHandle ob, Object obb, long base) {
        return getOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static boolean getOpaque(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getBooleanOpaque(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, boolean value) {
        setOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, long offset, boolean value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putBooleanOpaque(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }
}
