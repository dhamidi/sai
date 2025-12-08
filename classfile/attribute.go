package classfile

import (
	"encoding/binary"
)

type AttributeInfo struct {
	NameIndex uint16
	Info      []byte
	Parsed    interface{}
}

type CodeAttribute struct {
	MaxStack       uint16
	MaxLocals      uint16
	Code           []byte
	ExceptionTable []ExceptionTableEntry
	Attributes     []AttributeInfo
}

type ExceptionTableEntry struct {
	StartPC   uint16
	EndPC     uint16
	HandlerPC uint16
	CatchType uint16
}

type LineNumberTableAttribute struct {
	LineNumberTable []LineNumberEntry
}

type LineNumberEntry struct {
	StartPC    uint16
	LineNumber uint16
}

type LocalVariableTableAttribute struct {
	LocalVariableTable []LocalVariableEntry
}

type LocalVariableEntry struct {
	StartPC         uint16
	Length          uint16
	NameIndex       uint16
	DescriptorIndex uint16
	Index           uint16
}

type SourceFileAttribute struct {
	SourceFileIndex uint16
}

type ConstantValueAttribute struct {
	ConstantValueIndex uint16
}

type ExceptionsAttribute struct {
	ExceptionIndexTable []uint16
}

type InnerClassesAttribute struct {
	Classes []InnerClassEntry
}

type InnerClassEntry struct {
	InnerClassInfoIndex   uint16
	OuterClassInfoIndex   uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags AccessFlags
}

type SignatureAttribute struct {
	SignatureIndex uint16
}

type BootstrapMethodsAttribute struct {
	BootstrapMethods []BootstrapMethod
}

type BootstrapMethod struct {
	BootstrapMethodRef uint16
	BootstrapArguments []uint16
}

func (a *AttributeInfo) AsCode() *CodeAttribute {
	if a.Parsed != nil {
		if code, ok := a.Parsed.(*CodeAttribute); ok {
			return code
		}
	}
	return nil
}

func (a *AttributeInfo) AsLineNumberTable() *LineNumberTableAttribute {
	if a.Parsed != nil {
		if lnt, ok := a.Parsed.(*LineNumberTableAttribute); ok {
			return lnt
		}
	}
	return nil
}

func (a *AttributeInfo) AsLocalVariableTable() *LocalVariableTableAttribute {
	if a.Parsed != nil {
		if lvt, ok := a.Parsed.(*LocalVariableTableAttribute); ok {
			return lvt
		}
	}
	return nil
}

func (a *AttributeInfo) AsSourceFile() *SourceFileAttribute {
	if a.Parsed != nil {
		if sf, ok := a.Parsed.(*SourceFileAttribute); ok {
			return sf
		}
	}
	return nil
}

func (a *AttributeInfo) AsConstantValue() *ConstantValueAttribute {
	if a.Parsed != nil {
		if cv, ok := a.Parsed.(*ConstantValueAttribute); ok {
			return cv
		}
	}
	return nil
}

func (a *AttributeInfo) AsExceptions() *ExceptionsAttribute {
	if a.Parsed != nil {
		if ex, ok := a.Parsed.(*ExceptionsAttribute); ok {
			return ex
		}
	}
	return nil
}

func (a *AttributeInfo) AsInnerClasses() *InnerClassesAttribute {
	if a.Parsed != nil {
		if ic, ok := a.Parsed.(*InnerClassesAttribute); ok {
			return ic
		}
	}
	return nil
}

func (a *AttributeInfo) AsSignature() *SignatureAttribute {
	if a.Parsed != nil {
		if sig, ok := a.Parsed.(*SignatureAttribute); ok {
			return sig
		}
	}
	return nil
}

func (a *AttributeInfo) AsBootstrapMethods() *BootstrapMethodsAttribute {
	if a.Parsed != nil {
		if bm, ok := a.Parsed.(*BootstrapMethodsAttribute); ok {
			return bm
		}
	}
	return nil
}

func parseCodeAttribute(info []byte, cp ConstantPool) *CodeAttribute {
	if len(info) < 8 {
		return nil
	}

	code := &CodeAttribute{
		MaxStack:  binary.BigEndian.Uint16(info[0:2]),
		MaxLocals: binary.BigEndian.Uint16(info[2:4]),
	}

	codeLength := binary.BigEndian.Uint32(info[4:8])
	if len(info) < 8+int(codeLength) {
		return nil
	}
	code.Code = info[8 : 8+codeLength]

	offset := 8 + int(codeLength)
	if len(info) < offset+2 {
		return nil
	}

	exceptionTableLength := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	code.ExceptionTable = make([]ExceptionTableEntry, exceptionTableLength)
	for i := uint16(0); i < exceptionTableLength; i++ {
		if len(info) < offset+8 {
			return nil
		}
		code.ExceptionTable[i] = ExceptionTableEntry{
			StartPC:   binary.BigEndian.Uint16(info[offset : offset+2]),
			EndPC:     binary.BigEndian.Uint16(info[offset+2 : offset+4]),
			HandlerPC: binary.BigEndian.Uint16(info[offset+4 : offset+6]),
			CatchType: binary.BigEndian.Uint16(info[offset+6 : offset+8]),
		}
		offset += 8
	}

	if len(info) < offset+2 {
		return nil
	}
	attributesCount := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	code.Attributes = make([]AttributeInfo, 0, attributesCount)
	for i := uint16(0); i < attributesCount; i++ {
		if len(info) < offset+6 {
			return nil
		}
		nameIndex := binary.BigEndian.Uint16(info[offset : offset+2])
		attrLength := binary.BigEndian.Uint32(info[offset+2 : offset+6])
		offset += 6

		if len(info) < offset+int(attrLength) {
			return nil
		}
		attrInfo := info[offset : offset+int(attrLength)]
		offset += int(attrLength)

		attr := AttributeInfo{
			NameIndex: nameIndex,
			Info:      attrInfo,
		}

		attrName := cp.GetUtf8(nameIndex)
		switch attrName {
		case "LineNumberTable":
			attr.Parsed = parseLineNumberTableAttribute(attrInfo)
		case "LocalVariableTable":
			attr.Parsed = parseLocalVariableTableAttribute(attrInfo)
		}

		code.Attributes = append(code.Attributes, attr)
	}

	return code
}

func parseLineNumberTableAttribute(info []byte) *LineNumberTableAttribute {
	if len(info) < 2 {
		return nil
	}

	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*4 {
		return nil
	}

	lnt := &LineNumberTableAttribute{
		LineNumberTable: make([]LineNumberEntry, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		lnt.LineNumberTable[i] = LineNumberEntry{
			StartPC:    binary.BigEndian.Uint16(info[offset : offset+2]),
			LineNumber: binary.BigEndian.Uint16(info[offset+2 : offset+4]),
		}
		offset += 4
	}

	return lnt
}

func parseLocalVariableTableAttribute(info []byte) *LocalVariableTableAttribute {
	if len(info) < 2 {
		return nil
	}

	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*10 {
		return nil
	}

	lvt := &LocalVariableTableAttribute{
		LocalVariableTable: make([]LocalVariableEntry, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		lvt.LocalVariableTable[i] = LocalVariableEntry{
			StartPC:         binary.BigEndian.Uint16(info[offset : offset+2]),
			Length:          binary.BigEndian.Uint16(info[offset+2 : offset+4]),
			NameIndex:       binary.BigEndian.Uint16(info[offset+4 : offset+6]),
			DescriptorIndex: binary.BigEndian.Uint16(info[offset+6 : offset+8]),
			Index:           binary.BigEndian.Uint16(info[offset+8 : offset+10]),
		}
		offset += 10
	}

	return lvt
}

func parseSourceFileAttribute(info []byte) *SourceFileAttribute {
	if len(info) < 2 {
		return nil
	}
	return &SourceFileAttribute{
		SourceFileIndex: binary.BigEndian.Uint16(info[0:2]),
	}
}

func parseConstantValueAttribute(info []byte) *ConstantValueAttribute {
	if len(info) < 2 {
		return nil
	}
	return &ConstantValueAttribute{
		ConstantValueIndex: binary.BigEndian.Uint16(info[0:2]),
	}
}

func parseExceptionsAttribute(info []byte) *ExceptionsAttribute {
	if len(info) < 2 {
		return nil
	}
	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*2 {
		return nil
	}

	ex := &ExceptionsAttribute{
		ExceptionIndexTable: make([]uint16, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		ex.ExceptionIndexTable[i] = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2
	}

	return ex
}

func parseInnerClassesAttribute(info []byte) *InnerClassesAttribute {
	if len(info) < 2 {
		return nil
	}
	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*8 {
		return nil
	}

	ic := &InnerClassesAttribute{
		Classes: make([]InnerClassEntry, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		ic.Classes[i] = InnerClassEntry{
			InnerClassInfoIndex:   binary.BigEndian.Uint16(info[offset : offset+2]),
			OuterClassInfoIndex:   binary.BigEndian.Uint16(info[offset+2 : offset+4]),
			InnerNameIndex:        binary.BigEndian.Uint16(info[offset+4 : offset+6]),
			InnerClassAccessFlags: AccessFlags(binary.BigEndian.Uint16(info[offset+6 : offset+8])),
		}
		offset += 8
	}

	return ic
}

func parseSignatureAttribute(info []byte) *SignatureAttribute {
	if len(info) < 2 {
		return nil
	}
	return &SignatureAttribute{
		SignatureIndex: binary.BigEndian.Uint16(info[0:2]),
	}
}

func parseBootstrapMethodsAttribute(info []byte) *BootstrapMethodsAttribute {
	if len(info) < 2 {
		return nil
	}
	count := binary.BigEndian.Uint16(info[0:2])

	bm := &BootstrapMethodsAttribute{
		BootstrapMethods: make([]BootstrapMethod, 0, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		if len(info) < offset+4 {
			return nil
		}
		methodRef := binary.BigEndian.Uint16(info[offset : offset+2])
		numArgs := binary.BigEndian.Uint16(info[offset+2 : offset+4])
		offset += 4

		if len(info) < offset+int(numArgs)*2 {
			return nil
		}
		args := make([]uint16, numArgs)
		for j := uint16(0); j < numArgs; j++ {
			args[j] = binary.BigEndian.Uint16(info[offset : offset+2])
			offset += 2
		}

		bm.BootstrapMethods = append(bm.BootstrapMethods, BootstrapMethod{
			BootstrapMethodRef: methodRef,
			BootstrapArguments: args,
		})
	}

	return bm
}
