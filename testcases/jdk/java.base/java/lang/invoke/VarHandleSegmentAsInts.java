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

sealed class VarHandleSegmentAsInts {

    static final int NON_PLAIN_ACCESS_MIN_ALIGN_MASK = Integer.BYTES - 1;

    static VarForm selectForm(long alignmentMask, boolean constantOffset) {
        return (alignmentMask & NON_PLAIN_ACCESS_MIN_ALIGN_MASK) != NON_PLAIN_ACCESS_MIN_ALIGN_MASK ?
                (constantOffset ? CONSTANT_OFFSET_FORM : VARIABLE_OFFSET_FORM) :
                (constantOffset ? VarHandleSegmentAsIntsAligned.CONSTANT_OFFSET_FORM : VarHandleSegmentAsIntsAligned.VARIABLE_OFFSET_FORM);
    }

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsInts.class, MemorySegment.class, int.class, long.class);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsInts.class, MemorySegment.class, int.class, long.class, long.class);

    VarHandleSegmentAsInts() { throw new AssertionError(); }

    @ForceInline
    static int get(VarHandle ob, Object obb, long base) {
        return get(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static int get(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return SCOPED_MEMORY_ACCESS.getIntUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                handle.be);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, int value) {
        set(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                value,
                handle.be);
    }
}

// This class must be accessed through non-aligned VarHandleSegmentAsInts
final class VarHandleSegmentAsIntsAligned extends VarHandleSegmentAsInts {

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsIntsAligned.class, VarHandleSegmentAsInts.CONSTANT_OFFSET_FORM);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsIntsAligned.class, VarHandleSegmentAsInts.VARIABLE_OFFSET_FORM);

    VarHandleSegmentAsIntsAligned() { throw new AssertionError(); }

    @ForceInline
    static int convEndian(boolean big, int n) {
        return big == BE ? n : Integer.reverseBytes(n);
    }

    @ForceInline
    static int getVolatile(VarHandle ob, Object obb, long base) {
        return getVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static int getVolatile(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getIntVolatile(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, int value) {
        setVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntVolatile(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static int getAcquire(VarHandle ob, Object obb, long base) {
        return getAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static int getAcquire(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getIntAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, int value) {
        setRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static int getOpaque(VarHandle ob, Object obb, long base) {
        return getOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static int getOpaque(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getIntOpaque(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, int value) {
        setOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntOpaque(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, int expected, int value) {
        return compareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.compareAndSetInt(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static int compareAndExchange(VarHandle ob, Object obb, long base, int expected, int value) {
        return compareAndExchange(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static int compareAndExchange(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeInt(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static int compareAndExchangeAcquire(VarHandle ob, Object obb, long base, int expected, int value) {
        return compareAndExchangeAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static int compareAndExchangeAcquire(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeIntAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static int compareAndExchangeRelease(VarHandle ob, Object obb, long base, int expected, int value) {
        return compareAndExchangeRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static int compareAndExchangeRelease(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeIntRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, int expected, int value) {
        return weakCompareAndSetPlain(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetIntPlain(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, int expected, int value) {
        return weakCompareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetInt(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, int expected, int value) {
        return weakCompareAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetIntAcquire(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, int expected, int value) {
        return weakCompareAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, long offset, int expected, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetIntRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static int getAndSet(VarHandle ob, Object obb, long base, int value) {
        return getAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndSet(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetInt(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static int getAndSetAcquire(VarHandle ob, Object obb, long base, int value) {
        return getAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndSetAcquire(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetIntAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static int getAndSetRelease(VarHandle ob, Object obb, long base, int value) {
        return getAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndSetRelease(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetIntRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static int getAndAdd(VarHandle ob, Object obb, long base, int value) {
        return getAndAdd(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndAdd(VarHandle ob, Object obb, long base, long offset, int delta) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndAddInt(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    delta);
        } else {
            return getAndAddConvEndianWithCAS(bb, offset(bb, base, offset), delta);
        }
    }

    @ForceInline
    static int getAndAddAcquire(VarHandle ob, Object obb, long base, int value) {
        return getAndAddAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndAddAcquire(VarHandle ob, Object obb, long base, long offset, int delta) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndAddIntAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    delta);
        } else {
            return getAndAddConvEndianWithCAS(bb, offset(bb, base, offset), delta);
        }
    }

    @ForceInline
    static int getAndAddRelease(VarHandle ob, Object obb, long base, int value) {
        return getAndAddRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndAddRelease(VarHandle ob, Object obb, long base, long offset, int delta) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndAddIntRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    delta);
        } else {
            return getAndAddConvEndianWithCAS(bb, offset(bb, base, offset), delta);
        }
    }

    @ForceInline
    static int getAndAddConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, int delta) {
        int nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getIntVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Integer.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetInt(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Integer.reverseBytes((expectedValue + delta))));
        return expectedValue;
    }

    @ForceInline
    static int getAndBitwiseOr(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseOr(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseOr(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseOrInt(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseOrConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseOrRelease(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseOrRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseOrRelease(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseOrIntRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseOrConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseOrAcquire(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseOrAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseOrAcquire(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseOrIntAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseOrConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseOrConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, int value) {
        int nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getIntVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Integer.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetInt(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Integer.reverseBytes((expectedValue | value))));
        return expectedValue;
    }

    @ForceInline
    static int getAndBitwiseAnd(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseAnd(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseAnd(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseAndInt(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseAndConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseAndRelease(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseAndRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseAndRelease(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseAndIntRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseAndConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseAndAcquire(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseAndAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }


    @ForceInline
    static int getAndBitwiseAndAcquire(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseAndIntAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseAndConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseAndConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, int value) {
        int nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getIntVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Integer.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetInt(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Integer.reverseBytes((expectedValue & value))));
        return expectedValue;
    }

    @ForceInline
    static int getAndBitwiseXor(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseXor(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseXor(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseXorInt(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseXorConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseXorRelease(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseXorRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseXorRelease(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseXorIntRelease(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseXorConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseXorAcquire(VarHandle ob, Object obb, long base, int value) {
        return getAndBitwiseXorAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static int getAndBitwiseXorAcquire(VarHandle ob, Object obb, long base, long offset, int value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        if (handle.be == BE) {
            return SCOPED_MEMORY_ACCESS.getAndBitwiseXorIntAcquire(bb.sessionImpl(),
                    bb.unsafeGetBase(),
                    offset(bb, base, offset),
                    value);
        } else {
            return getAndBitwiseXorConvEndianWithCAS(bb, offset(bb, base, offset), value);
        }
    }

    @ForceInline
    static int getAndBitwiseXorConvEndianWithCAS(AbstractMemorySegmentImpl  bb, long offset, int value) {
        int nativeExpectedValue, expectedValue;
        Object base = bb.unsafeGetBase();
        do {
            nativeExpectedValue = SCOPED_MEMORY_ACCESS.getIntVolatile(bb.sessionImpl(),base, offset);
            expectedValue = Integer.reverseBytes(nativeExpectedValue);
        } while (!SCOPED_MEMORY_ACCESS.weakCompareAndSetInt(bb.sessionImpl(),base, offset,
                nativeExpectedValue, Integer.reverseBytes((expectedValue ^ value))));
        return expectedValue;
    }
}
