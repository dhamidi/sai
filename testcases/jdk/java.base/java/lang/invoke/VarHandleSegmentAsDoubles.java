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

sealed class VarHandleSegmentAsDoubles {

    static final int NON_PLAIN_ACCESS_MIN_ALIGN_MASK = Double.BYTES - 1;

    static VarForm selectForm(long alignmentMask, boolean constantOffset) {
        return (alignmentMask & NON_PLAIN_ACCESS_MIN_ALIGN_MASK) != NON_PLAIN_ACCESS_MIN_ALIGN_MASK ?
                (constantOffset ? CONSTANT_OFFSET_FORM : VARIABLE_OFFSET_FORM) :
                (constantOffset ? VarHandleSegmentAsDoublesAligned.CONSTANT_OFFSET_FORM : VarHandleSegmentAsDoublesAligned.VARIABLE_OFFSET_FORM);
    }

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsDoubles.class, MemorySegment.class, double.class, long.class);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsDoubles.class, MemorySegment.class, double.class, long.class, long.class);

    VarHandleSegmentAsDoubles() { throw new AssertionError(); }

    @ForceInline
    static double get(VarHandle ob, Object obb, long base) {
        return get(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static double get(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        long rawValue = SCOPED_MEMORY_ACCESS.getLongUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                handle.be);
        return Double.longBitsToDouble(rawValue);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, double value) {
        set(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                Double.doubleToRawLongBits(value),
                handle.be);
    }
}

// This class must be accessed through non-aligned VarHandleSegmentAsDoubles
final class VarHandleSegmentAsDoublesAligned extends VarHandleSegmentAsDoubles {

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsDoublesAligned.class, VarHandleSegmentAsDoubles.CONSTANT_OFFSET_FORM);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsDoublesAligned.class, VarHandleSegmentAsDoubles.VARIABLE_OFFSET_FORM);

    VarHandleSegmentAsDoublesAligned() { throw new AssertionError(); }

    @ForceInline
    static long convEndian(boolean big, double v) {
        long rv = Double.doubleToRawLongBits(v);
        return big == BE ? rv : Long.reverseBytes(rv);
    }

    @ForceInline
    static double convEndian(boolean big, long rv) {
        rv = big == BE ? rv : Long.reverseBytes(rv);
        return Double.longBitsToDouble(rv);
    }

    @ForceInline
    static double getVolatile(VarHandle ob, Object obb, long base) {
        return getVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static double getVolatile(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getLongVolatile(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, double value) {
        setVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongVolatile(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static double getAcquire(VarHandle ob, Object obb, long base) {
        return getAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static double getAcquire(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getLongAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, double value) {
        setRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static double getOpaque(VarHandle ob, Object obb, long base) {
        return getOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static double getOpaque(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getLongOpaque(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, double value) {
        setOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongOpaque(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, double expected, double value) {
        return compareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.compareAndSetLong(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static double compareAndExchange(VarHandle ob, Object obb, long base, double expected, double value) {
        return compareAndExchange(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static double compareAndExchange(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeLong(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static double compareAndExchangeAcquire(VarHandle ob, Object obb, long base, double expected, double value) {
        return compareAndExchangeAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static double compareAndExchangeAcquire(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeLongAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static double compareAndExchangeRelease(VarHandle ob, Object obb, long base, double expected, double value) {
        return compareAndExchangeRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static double compareAndExchangeRelease(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeLongRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, double expected, double value) {
        return weakCompareAndSetPlain(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLongPlain(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, double expected, double value) {
        return weakCompareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLong(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, double expected, double value) {
        return weakCompareAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLongAcquire(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, double expected, double value) {
        return weakCompareAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, long offset, double expected, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLongRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static double getAndSet(VarHandle ob, Object obb, long base, double value) {
        return getAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static double getAndSet(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetLong(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static double getAndSetAcquire(VarHandle ob, Object obb, long base, double value) {
        return getAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static double getAndSetAcquire(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetLongAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static double getAndSetRelease(VarHandle ob, Object obb, long base, double value) {
        return getAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static double getAndSetRelease(VarHandle ob, Object obb, long base, long offset, double value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetLongRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }
}
