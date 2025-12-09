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

sealed class VarHandleSegmentAsLongs {

    static final int NON_PLAIN_ACCESS_MIN_ALIGN_MASK = Long.BYTES - 1;

    static VarForm selectForm(long alignmentMask, boolean constantOffset) {
        return (alignmentMask & NON_PLAIN_ACCESS_MIN_ALIGN_MASK) != NON_PLAIN_ACCESS_MIN_ALIGN_MASK ?
                (constantOffset ? CONSTANT_OFFSET_FORM : VARIABLE_OFFSET_FORM) :
                (constantOffset ? VarHandleSegmentAsLongsAligned.CONSTANT_OFFSET_FORM : VarHandleSegmentAsLongsAligned.VARIABLE_OFFSET_FORM);
    }

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsLongs.class, MemorySegment.class, long.class, long.class);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsLongs.class, MemorySegment.class, long.class, long.class, long.class);

    VarHandleSegmentAsLongs() { throw new AssertionError(); }

    @ForceInline
    static long get(VarHandle ob, Object obb, long base) {
        return get(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static long get(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return SCOPED_MEMORY_ACCESS.getLongUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                handle.be);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, long value) {
        set(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                value,
                handle.be);
    }
}

// This class must be accessed through non-aligned VarHandleSegmentAsLongs
final class VarHandleSegmentAsLongsAligned extends VarHandleSegmentAsLongs {

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsLongsAligned.class, VarHandleSegmentAsLongs.CONSTANT_OFFSET_FORM);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsLongsAligned.class, VarHandleSegmentAsLongs.VARIABLE_OFFSET_FORM);

    VarHandleSegmentAsLongsAligned() { throw new AssertionError(); }

    @ForceInline
    static long convEndian(boolean big, long n) {
        return big == BE ? n : Long.reverseBytes(n);
    }

    @ForceInline
    static long getVolatile(VarHandle ob, Object obb, long base) {
        return getVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static long getVolatile(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getLongVolatile(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, long value) {
        setVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongVolatile(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static long getAcquire(VarHandle ob, Object obb, long base) {
        return getAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static long getAcquire(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getLongAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, long value) {
        setRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static long getOpaque(VarHandle ob, Object obb, long base) {
        return getOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static long getOpaque(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getLongOpaque(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, long value) {
        setOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putLongOpaque(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, long expected, long value) {
        return compareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.compareAndSetLong(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static long compareAndExchange(VarHandle ob, Object obb, long base, long expected, long value) {
        return compareAndExchange(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static long compareAndExchange(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeLong(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static long compareAndExchangeAcquire(VarHandle ob, Object obb, long base, long expected, long value) {
        return compareAndExchangeAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static long compareAndExchangeAcquire(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeLongAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static long compareAndExchangeRelease(VarHandle ob, Object obb, long base, long expected, long value) {
        return compareAndExchangeRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static long compareAndExchangeRelease(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeLongRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, long expected, long value) {
        return weakCompareAndSetPlain(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLongPlain(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, long expected, long value) {
        return weakCompareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLong(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, long expected, long value) {
        return weakCompareAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLongAcquire(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, long expected, long value) {
        return weakCompareAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, long offset, long expected, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetLongRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static long getAndSet(VarHandle ob, Object obb, long base, long value) {
        return getAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndSet(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetLong(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static long getAndSetAcquire(VarHandle ob, Object obb, long base, long value) {
        return getAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndSetAcquire(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetLongAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static long getAndSetRelease(VarHandle ob, Object obb, long base, long value) {
        return getAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndSetRelease(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetLongRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static long getAndAdd(VarHandle ob, Object obb, long base, long value) {
        return getAndAdd(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndAdd(VarHandle ob, Object obb, long base, long offset, long delta) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndAddLong(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    delta);
        } else {
            return getAndAddConvEndianWithCAS(bb, offset(bb, base, offset), delta);
        }
    }

    @ForceInline
    static long getAndAddAcquire(VarHandle ob, Object obb, long base, long value) {
        return getAndAddAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndAddAcquire(VarHandle ob, Object obb, long base, long offset, long delta) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndAddLongAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    delta);
        } else {
            return getAndAddConvEndianWithCAS(bb, offset(bb, base, offset), delta);
        }
    }

    @ForceInline
    static long getAndAddRelease(VarHandle ob, Object obb, long base, long value) {
        return getAndAddRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndAddRelease(VarHandle ob, Object obb, long base, long offset, long delta) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndAddLongRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    delta);
        } else {
            return getAndAddConvEndianWithCAS(bb, offset(bb, base, offset), delta);
        }
    }

    @ForceInline
    static long getAndAddConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, long delta) {
        long nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getLongVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Long.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetLong(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Long.reverseBytes((expectedValue + delta))));
        return expectedValue;
    }

    @ForceInline
    static long getAndBitwiseOr(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseOr(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseOr(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseOrLong(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseOrConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseOrRelease(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseOrRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseOrRelease(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseOrLongRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseOrConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseOrAcquire(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseOrAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseOrAcquire(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseOrLongAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseOrConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseOrConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, long value) {
        long nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getLongVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Long.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetLong(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Long.reverseBytes((expectedValue | value))));
        return expectedValue;
    }

    @ForceInline
    static long getAndBitwiseAnd(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseAnd(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseAnd(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseAndLong(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseAndConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseAndRelease(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseAndRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseAndRelease(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseAndLongRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseAndConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseAndAcquire(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseAndAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }


    @ForceInline
    static long getAndBitwiseAndAcquire(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseAndLongAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseAndConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseAndConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, long value) {
        long nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getLongVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Long.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetLong(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Long.reverseBytes((expectedValue & value))));
        return expectedValue;
    }

    @ForceInline
    static long getAndBitwiseXor(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseXor(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseXor(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseXorLong(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseXorConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseXorRelease(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseXorRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseXorRelease(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseXorLongRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseXorConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseXorAcquire(VarHandle ob, Object obb, long base, long value) {
        return getAndBitwiseXorAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static long getAndBitwiseXorAcquire(VarHandle ob, Object obb, long base, long offset, long value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseXorLongAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseXorConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static long getAndBitwiseXorConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, long value) {
        long nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getLongVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Long.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetLong(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Long.reverseBytes((expectedValue ^ value))));
        return expectedValue;
    }
}
