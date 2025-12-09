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

sealed class VarHandleSegmentAsFloats {

    static final int NON_PLAIN_ACCESS_MIN_ALIGN_MASK = Float.BYTES - 1;

    static VarForm selectForm(long alignmentMask, boolean constantOffset) {
        return (alignmentMask & NON_PLAIN_ACCESS_MIN_ALIGN_MASK) != NON_PLAIN_ACCESS_MIN_ALIGN_MASK ?
                (constantOffset ? CONSTANT_OFFSET_FORM : VARIABLE_OFFSET_FORM) :
                (constantOffset ? VarHandleSegmentAsFloatsAligned.CONSTANT_OFFSET_FORM : VarHandleSegmentAsFloatsAligned.VARIABLE_OFFSET_FORM);
    }

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsFloats.class, MemorySegment.class, float.class, long.class);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsFloats.class, MemorySegment.class, float.class, long.class, long.class);

    VarHandleSegmentAsFloats() { throw new AssertionError(); }

    @ForceInline
    static float get(VarHandle ob, Object obb, long base) {
        return get(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static float get(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        int rawValue = SCOPED_MEMORY_ACCESS.getIntUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                handle.be);
        return Float.intBitsToFloat(rawValue);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, float value) {
        set(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void set(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntUnaligned(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                Float.floatToRawIntBits(value),
                handle.be);
    }
}

// This class must be accessed through non-aligned VarHandleSegmentAsFloats
final class VarHandleSegmentAsFloatsAligned extends VarHandleSegmentAsFloats {

    static final VarForm CONSTANT_OFFSET_FORM = new VarForm(VarHandleSegmentAsFloatsAligned.class, VarHandleSegmentAsFloats.CONSTANT_OFFSET_FORM);
    static final VarForm VARIABLE_OFFSET_FORM = new VarForm(VarHandleSegmentAsFloatsAligned.class, VarHandleSegmentAsFloats.VARIABLE_OFFSET_FORM);

    VarHandleSegmentAsFloatsAligned() { throw new AssertionError(); }

    @ForceInline
    static int convEndian(boolean big, float v) {
        int rv = Float.floatToRawIntBits(v);
        return big == BE ? rv : Integer.reverseBytes(rv);
    }

    @ForceInline
    static float convEndian(boolean big, int rv) {
        rv = big == BE ? rv : Integer.reverseBytes(rv);
        return Float.intBitsToFloat(rv);
    }

    @ForceInline
    static float getVolatile(VarHandle ob, Object obb, long base) {
        return getVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static float getVolatile(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getIntVolatile(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, float value) {
        setVolatile(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setVolatile(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntVolatile(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static float getAcquire(VarHandle ob, Object obb, long base) {
        return getAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static float getAcquire(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getIntAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, float value) {
        setRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setRelease(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static float getOpaque(VarHandle ob, Object obb, long base) {
        return getOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset);
    }

    @ForceInline
    static float getOpaque(VarHandle ob, Object obb, long base, long offset) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, true);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getIntOpaque(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset)));
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, float value) {
        setOpaque(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static void setOpaque(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        SCOPED_MEMORY_ACCESS.putIntOpaque(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, value));
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, float expected, float value) {
        return compareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean compareAndSet(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.compareAndSetInt(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static float compareAndExchange(VarHandle ob, Object obb, long base, float expected, float value) {
        return compareAndExchange(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static float compareAndExchange(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeInt(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static float compareAndExchangeAcquire(VarHandle ob, Object obb, long base, float expected, float value) {
        return compareAndExchangeAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static float compareAndExchangeAcquire(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeIntAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static float compareAndExchangeRelease(VarHandle ob, Object obb, long base, float expected, float value) {
        return compareAndExchangeRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static float compareAndExchangeRelease(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.compareAndExchangeIntRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, expected), convEndian(handle.be, value)));
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, float expected, float value) {
        return weakCompareAndSetPlain(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetPlain(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetIntPlain(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, float expected, float value) {
        return weakCompareAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSet(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetInt(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, float expected, float value) {
        return weakCompareAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetAcquire(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetIntAcquire(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, float expected, float value) {
        return weakCompareAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, expected, value);
    }

    @ForceInline
    static boolean weakCompareAndSetRelease(VarHandle ob, Object obb, long base, long offset, float expected, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return SCOPED_MEMORY_ACCESS.weakCompareAndSetIntRelease(bb.sessionImpl(),
                bb.unsafeGetBase(),
                offset(bb, base, offset),
                convEndian(handle.be, expected), convEndian(handle.be, value));
    }

    @ForceInline
    static float getAndSet(VarHandle ob, Object obb, long base, float value) {
        return getAndSet(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static float getAndSet(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetInt(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static float getAndSetAcquire(VarHandle ob, Object obb, long base, float value) {
        return getAndSetAcquire(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static float getAndSetAcquire(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetIntAcquire(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }

    @ForceInline
    static float getAndSetRelease(VarHandle ob, Object obb, long base, float value) {
        return getAndSetRelease(ob, obb, base, ((SegmentVarHandle) ob).offset, value);
    }

    @ForceInline
    static float getAndSetRelease(VarHandle ob, Object obb, long base, long offset, float value) {
        SegmentVarHandle handle = (SegmentVarHandle)ob;
        AbstractMemorySegmentImpl bb = handle.checkSegment(obb, base, false);
        return convEndian(handle.be,
                          SCOPED_MEMORY_ACCESS.getAndSetIntRelease(bb.sessionImpl(),
                                  bb.unsafeGetBase(),
                                  offset(bb, base, offset),
                                  convEndian(handle.be, value)));
    }
}
