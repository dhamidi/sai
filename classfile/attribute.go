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

type EnclosingMethodAttribute struct {
	ClassIndex  uint16
	MethodIndex uint16
}

type SyntheticAttribute struct{}

type DeprecatedAttribute struct{}

type SourceDebugExtensionAttribute struct {
	DebugExtension string
}

type LocalVariableTypeTableAttribute struct {
	LocalVariableTypeTable []LocalVariableTypeEntry
}

type LocalVariableTypeEntry struct {
	StartPC        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

type MethodParametersAttribute struct {
	Parameters []MethodParameter
}

type MethodParameter struct {
	NameIndex   uint16
	AccessFlags AccessFlags
}

type NestHostAttribute struct {
	HostClassIndex uint16
}

type NestMembersAttribute struct {
	Classes []uint16
}

type RecordAttribute struct {
	Components []RecordComponentInfo
}

type RecordComponentInfo struct {
	NameIndex       uint16
	DescriptorIndex uint16
	Attributes      []AttributeInfo
}

type PermittedSubclassesAttribute struct {
	Classes []uint16
}

type StackMapTableAttribute struct {
	Entries []StackMapFrame
}

type StackMapFrame struct {
	FrameType uint8
	Data      []byte
}

type Annotation struct {
	TypeIndex         uint16
	ElementValuePairs []ElementValuePair
}

type ElementValuePair struct {
	ElementNameIndex uint16
	Value            ElementValue
}

type ElementValue struct {
	Tag   byte
	Value interface{}
}

type EnumConstValue struct {
	TypeNameIndex  uint16
	ConstNameIndex uint16
}

type ArrayValue struct {
	Values []ElementValue
}

type RuntimeVisibleAnnotationsAttribute struct {
	Annotations []Annotation
}

type RuntimeInvisibleAnnotationsAttribute struct {
	Annotations []Annotation
}

type RuntimeVisibleParameterAnnotationsAttribute struct {
	ParameterAnnotations [][]Annotation
}

type RuntimeInvisibleParameterAnnotationsAttribute struct {
	ParameterAnnotations [][]Annotation
}

type TypeAnnotation struct {
	TargetType        uint8
	TargetInfo        []byte
	TargetPath        []TypePathEntry
	TypeIndex         uint16
	ElementValuePairs []ElementValuePair
}

type TypePathEntry struct {
	TypePathKind      uint8
	TypeArgumentIndex uint8
}

type RuntimeVisibleTypeAnnotationsAttribute struct {
	Annotations []TypeAnnotation
}

type RuntimeInvisibleTypeAnnotationsAttribute struct {
	Annotations []TypeAnnotation
}

type AnnotationDefaultAttribute struct {
	DefaultValue ElementValue
}

type ModuleAttribute struct {
	ModuleNameIndex    uint16
	ModuleFlags        uint16
	ModuleVersionIndex uint16
	Requires           []ModuleRequires
	Exports            []ModuleExports
	Opens              []ModuleOpens
	Uses               []uint16
	Provides           []ModuleProvides
}

type ModuleRequires struct {
	RequiresIndex        uint16
	RequiresFlags        uint16
	RequiresVersionIndex uint16
}

type ModuleExports struct {
	ExportsIndex   uint16
	ExportsFlags   uint16
	ExportsToIndex []uint16
}

type ModuleOpens struct {
	OpensIndex   uint16
	OpensFlags   uint16
	OpensToIndex []uint16
}

type ModuleProvides struct {
	ProvidesIndex     uint16
	ProvidesWithIndex []uint16
}

type ModulePackagesAttribute struct {
	PackageIndex []uint16
}

type ModuleMainClassAttribute struct {
	MainClassIndex uint16
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

func (a *AttributeInfo) AsEnclosingMethod() *EnclosingMethodAttribute {
	if a.Parsed != nil {
		if em, ok := a.Parsed.(*EnclosingMethodAttribute); ok {
			return em
		}
	}
	return nil
}

func (a *AttributeInfo) AsSynthetic() *SyntheticAttribute {
	if a.Parsed != nil {
		if s, ok := a.Parsed.(*SyntheticAttribute); ok {
			return s
		}
	}
	return nil
}

func (a *AttributeInfo) AsDeprecated() *DeprecatedAttribute {
	if a.Parsed != nil {
		if d, ok := a.Parsed.(*DeprecatedAttribute); ok {
			return d
		}
	}
	return nil
}

func (a *AttributeInfo) AsSourceDebugExtension() *SourceDebugExtensionAttribute {
	if a.Parsed != nil {
		if sde, ok := a.Parsed.(*SourceDebugExtensionAttribute); ok {
			return sde
		}
	}
	return nil
}

func (a *AttributeInfo) AsLocalVariableTypeTable() *LocalVariableTypeTableAttribute {
	if a.Parsed != nil {
		if lvtt, ok := a.Parsed.(*LocalVariableTypeTableAttribute); ok {
			return lvtt
		}
	}
	return nil
}

func (a *AttributeInfo) AsMethodParameters() *MethodParametersAttribute {
	if a.Parsed != nil {
		if mp, ok := a.Parsed.(*MethodParametersAttribute); ok {
			return mp
		}
	}
	return nil
}

func (a *AttributeInfo) AsNestHost() *NestHostAttribute {
	if a.Parsed != nil {
		if nh, ok := a.Parsed.(*NestHostAttribute); ok {
			return nh
		}
	}
	return nil
}

func (a *AttributeInfo) AsNestMembers() *NestMembersAttribute {
	if a.Parsed != nil {
		if nm, ok := a.Parsed.(*NestMembersAttribute); ok {
			return nm
		}
	}
	return nil
}

func (a *AttributeInfo) AsRecord() *RecordAttribute {
	if a.Parsed != nil {
		if r, ok := a.Parsed.(*RecordAttribute); ok {
			return r
		}
	}
	return nil
}

func (a *AttributeInfo) AsPermittedSubclasses() *PermittedSubclassesAttribute {
	if a.Parsed != nil {
		if ps, ok := a.Parsed.(*PermittedSubclassesAttribute); ok {
			return ps
		}
	}
	return nil
}

func (a *AttributeInfo) AsStackMapTable() *StackMapTableAttribute {
	if a.Parsed != nil {
		if smt, ok := a.Parsed.(*StackMapTableAttribute); ok {
			return smt
		}
	}
	return nil
}

func (a *AttributeInfo) AsRuntimeVisibleAnnotations() *RuntimeVisibleAnnotationsAttribute {
	if a.Parsed != nil {
		if rva, ok := a.Parsed.(*RuntimeVisibleAnnotationsAttribute); ok {
			return rva
		}
	}
	return nil
}

func (a *AttributeInfo) AsRuntimeInvisibleAnnotations() *RuntimeInvisibleAnnotationsAttribute {
	if a.Parsed != nil {
		if ria, ok := a.Parsed.(*RuntimeInvisibleAnnotationsAttribute); ok {
			return ria
		}
	}
	return nil
}

func (a *AttributeInfo) AsRuntimeVisibleParameterAnnotations() *RuntimeVisibleParameterAnnotationsAttribute {
	if a.Parsed != nil {
		if rvpa, ok := a.Parsed.(*RuntimeVisibleParameterAnnotationsAttribute); ok {
			return rvpa
		}
	}
	return nil
}

func (a *AttributeInfo) AsRuntimeInvisibleParameterAnnotations() *RuntimeInvisibleParameterAnnotationsAttribute {
	if a.Parsed != nil {
		if ripa, ok := a.Parsed.(*RuntimeInvisibleParameterAnnotationsAttribute); ok {
			return ripa
		}
	}
	return nil
}

func (a *AttributeInfo) AsRuntimeVisibleTypeAnnotations() *RuntimeVisibleTypeAnnotationsAttribute {
	if a.Parsed != nil {
		if rvta, ok := a.Parsed.(*RuntimeVisibleTypeAnnotationsAttribute); ok {
			return rvta
		}
	}
	return nil
}

func (a *AttributeInfo) AsRuntimeInvisibleTypeAnnotations() *RuntimeInvisibleTypeAnnotationsAttribute {
	if a.Parsed != nil {
		if rita, ok := a.Parsed.(*RuntimeInvisibleTypeAnnotationsAttribute); ok {
			return rita
		}
	}
	return nil
}

func (a *AttributeInfo) AsAnnotationDefault() *AnnotationDefaultAttribute {
	if a.Parsed != nil {
		if ad, ok := a.Parsed.(*AnnotationDefaultAttribute); ok {
			return ad
		}
	}
	return nil
}

func (a *AttributeInfo) AsModule() *ModuleAttribute {
	if a.Parsed != nil {
		if m, ok := a.Parsed.(*ModuleAttribute); ok {
			return m
		}
	}
	return nil
}

func (a *AttributeInfo) AsModulePackages() *ModulePackagesAttribute {
	if a.Parsed != nil {
		if mp, ok := a.Parsed.(*ModulePackagesAttribute); ok {
			return mp
		}
	}
	return nil
}

func (a *AttributeInfo) AsModuleMainClass() *ModuleMainClassAttribute {
	if a.Parsed != nil {
		if mmc, ok := a.Parsed.(*ModuleMainClassAttribute); ok {
			return mmc
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
		case "LocalVariableTypeTable":
			attr.Parsed = parseLocalVariableTypeTableAttribute(attrInfo)
		case "StackMapTable":
			attr.Parsed = parseStackMapTableAttribute(attrInfo)
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

func parseEnclosingMethodAttribute(info []byte) *EnclosingMethodAttribute {
	if len(info) < 4 {
		return nil
	}
	return &EnclosingMethodAttribute{
		ClassIndex:  binary.BigEndian.Uint16(info[0:2]),
		MethodIndex: binary.BigEndian.Uint16(info[2:4]),
	}
}

func parseSyntheticAttribute(_ []byte) *SyntheticAttribute {
	return &SyntheticAttribute{}
}

func parseDeprecatedAttribute(_ []byte) *DeprecatedAttribute {
	return &DeprecatedAttribute{}
}

func parseSourceDebugExtensionAttribute(info []byte) *SourceDebugExtensionAttribute {
	return &SourceDebugExtensionAttribute{
		DebugExtension: string(info),
	}
}

func parseLocalVariableTypeTableAttribute(info []byte) *LocalVariableTypeTableAttribute {
	if len(info) < 2 {
		return nil
	}

	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*10 {
		return nil
	}

	lvtt := &LocalVariableTypeTableAttribute{
		LocalVariableTypeTable: make([]LocalVariableTypeEntry, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		lvtt.LocalVariableTypeTable[i] = LocalVariableTypeEntry{
			StartPC:        binary.BigEndian.Uint16(info[offset : offset+2]),
			Length:         binary.BigEndian.Uint16(info[offset+2 : offset+4]),
			NameIndex:      binary.BigEndian.Uint16(info[offset+4 : offset+6]),
			SignatureIndex: binary.BigEndian.Uint16(info[offset+6 : offset+8]),
			Index:          binary.BigEndian.Uint16(info[offset+8 : offset+10]),
		}
		offset += 10
	}

	return lvtt
}

func parseMethodParametersAttribute(info []byte) *MethodParametersAttribute {
	if len(info) < 1 {
		return nil
	}

	count := uint8(info[0])
	if len(info) < 1+int(count)*4 {
		return nil
	}

	mp := &MethodParametersAttribute{
		Parameters: make([]MethodParameter, count),
	}

	offset := 1
	for i := uint8(0); i < count; i++ {
		mp.Parameters[i] = MethodParameter{
			NameIndex:   binary.BigEndian.Uint16(info[offset : offset+2]),
			AccessFlags: AccessFlags(binary.BigEndian.Uint16(info[offset+2 : offset+4])),
		}
		offset += 4
	}

	return mp
}

func parseNestHostAttribute(info []byte) *NestHostAttribute {
	if len(info) < 2 {
		return nil
	}
	return &NestHostAttribute{
		HostClassIndex: binary.BigEndian.Uint16(info[0:2]),
	}
}

func parseNestMembersAttribute(info []byte) *NestMembersAttribute {
	if len(info) < 2 {
		return nil
	}
	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*2 {
		return nil
	}

	nm := &NestMembersAttribute{
		Classes: make([]uint16, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		nm.Classes[i] = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2
	}

	return nm
}

func parseRecordAttribute(info []byte, cp ConstantPool) *RecordAttribute {
	if len(info) < 2 {
		return nil
	}

	count := binary.BigEndian.Uint16(info[0:2])
	rec := &RecordAttribute{
		Components: make([]RecordComponentInfo, 0, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		if len(info) < offset+6 {
			return nil
		}
		nameIndex := binary.BigEndian.Uint16(info[offset : offset+2])
		descriptorIndex := binary.BigEndian.Uint16(info[offset+2 : offset+4])
		attributesCount := binary.BigEndian.Uint16(info[offset+4 : offset+6])
		offset += 6

		attrs := make([]AttributeInfo, 0, attributesCount)
		for j := uint16(0); j < attributesCount; j++ {
			if len(info) < offset+6 {
				return nil
			}
			attrNameIndex := binary.BigEndian.Uint16(info[offset : offset+2])
			attrLength := binary.BigEndian.Uint32(info[offset+2 : offset+6])
			offset += 6

			if len(info) < offset+int(attrLength) {
				return nil
			}
			attrInfo := info[offset : offset+int(attrLength)]
			offset += int(attrLength)

			attr := AttributeInfo{
				NameIndex: attrNameIndex,
				Info:      attrInfo,
			}

			attrName := cp.GetUtf8(attrNameIndex)
			switch attrName {
			case "Signature":
				attr.Parsed = parseSignatureAttribute(attrInfo)
			}

			attrs = append(attrs, attr)
		}

		rec.Components = append(rec.Components, RecordComponentInfo{
			NameIndex:       nameIndex,
			DescriptorIndex: descriptorIndex,
			Attributes:      attrs,
		})
	}

	return rec
}

func parsePermittedSubclassesAttribute(info []byte) *PermittedSubclassesAttribute {
	if len(info) < 2 {
		return nil
	}
	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*2 {
		return nil
	}

	ps := &PermittedSubclassesAttribute{
		Classes: make([]uint16, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		ps.Classes[i] = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2
	}

	return ps
}

func parseStackMapTableAttribute(info []byte) *StackMapTableAttribute {
	if len(info) < 2 {
		return nil
	}

	count := binary.BigEndian.Uint16(info[0:2])
	smt := &StackMapTableAttribute{
		Entries: make([]StackMapFrame, 0, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		if len(info) <= offset {
			return nil
		}
		frameType := info[offset]
		frameStart := offset
		offset++

		switch {
		case frameType <= 63:
			// same_frame
		case frameType <= 127:
			// same_locals_1_stack_item_frame
			offset += verificationTypeInfoSize(info, offset)
		case frameType == 247:
			// same_locals_1_stack_item_frame_extended
			offset += 2 // offset_delta
			offset += verificationTypeInfoSize(info, offset)
		case frameType >= 248 && frameType <= 250:
			// chop_frame
			offset += 2 // offset_delta
		case frameType == 251:
			// same_frame_extended
			offset += 2 // offset_delta
		case frameType >= 252 && frameType <= 254:
			// append_frame
			offset += 2 // offset_delta
			numLocals := int(frameType) - 251
			for k := 0; k < numLocals; k++ {
				offset += verificationTypeInfoSize(info, offset)
			}
		case frameType == 255:
			// full_frame
			if len(info) < offset+2 {
				return nil
			}
			offset += 2 // offset_delta
			numLocals := int(binary.BigEndian.Uint16(info[offset : offset+2]))
			offset += 2
			for k := 0; k < numLocals; k++ {
				offset += verificationTypeInfoSize(info, offset)
			}
			if len(info) < offset+2 {
				return nil
			}
			numStack := int(binary.BigEndian.Uint16(info[offset : offset+2]))
			offset += 2
			for k := 0; k < numStack; k++ {
				offset += verificationTypeInfoSize(info, offset)
			}
		}

		smt.Entries = append(smt.Entries, StackMapFrame{
			FrameType: frameType,
			Data:      info[frameStart:offset],
		})
	}

	return smt
}

func verificationTypeInfoSize(info []byte, offset int) int {
	if len(info) <= offset {
		return 1
	}
	tag := info[offset]
	switch tag {
	case 0, 1, 2, 3, 4, 5, 6:
		return 1
	case 7, 8:
		return 3
	default:
		return 1
	}
}

func parseElementValue(info []byte, offset int) (ElementValue, int) {
	if len(info) <= offset {
		return ElementValue{}, offset
	}

	tag := info[offset]
	offset++

	ev := ElementValue{Tag: tag}

	switch tag {
	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 's':
		if len(info) < offset+2 {
			return ev, offset
		}
		ev.Value = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2

	case 'e':
		if len(info) < offset+4 {
			return ev, offset
		}
		ev.Value = EnumConstValue{
			TypeNameIndex:  binary.BigEndian.Uint16(info[offset : offset+2]),
			ConstNameIndex: binary.BigEndian.Uint16(info[offset+2 : offset+4]),
		}
		offset += 4

	case 'c':
		if len(info) < offset+2 {
			return ev, offset
		}
		ev.Value = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2

	case '@':
		var ann Annotation
		ann, offset = parseAnnotation(info, offset)
		ev.Value = ann

	case '[':
		if len(info) < offset+2 {
			return ev, offset
		}
		numValues := binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2
		values := make([]ElementValue, numValues)
		for i := uint16(0); i < numValues; i++ {
			values[i], offset = parseElementValue(info, offset)
		}
		ev.Value = ArrayValue{Values: values}
	}

	return ev, offset
}

func parseAnnotation(info []byte, offset int) (Annotation, int) {
	ann := Annotation{}
	if len(info) < offset+4 {
		return ann, offset
	}

	ann.TypeIndex = binary.BigEndian.Uint16(info[offset : offset+2])
	numPairs := binary.BigEndian.Uint16(info[offset+2 : offset+4])
	offset += 4

	ann.ElementValuePairs = make([]ElementValuePair, numPairs)
	for i := uint16(0); i < numPairs; i++ {
		if len(info) < offset+2 {
			return ann, offset
		}
		pair := ElementValuePair{
			ElementNameIndex: binary.BigEndian.Uint16(info[offset : offset+2]),
		}
		offset += 2
		pair.Value, offset = parseElementValue(info, offset)
		ann.ElementValuePairs[i] = pair
	}

	return ann, offset
}

func parseRuntimeVisibleAnnotationsAttribute(info []byte) *RuntimeVisibleAnnotationsAttribute {
	if len(info) < 2 {
		return nil
	}

	numAnnotations := binary.BigEndian.Uint16(info[0:2])
	rva := &RuntimeVisibleAnnotationsAttribute{
		Annotations: make([]Annotation, numAnnotations),
	}

	offset := 2
	for i := uint16(0); i < numAnnotations; i++ {
		rva.Annotations[i], offset = parseAnnotation(info, offset)
	}

	return rva
}

func parseRuntimeInvisibleAnnotationsAttribute(info []byte) *RuntimeInvisibleAnnotationsAttribute {
	if len(info) < 2 {
		return nil
	}

	numAnnotations := binary.BigEndian.Uint16(info[0:2])
	ria := &RuntimeInvisibleAnnotationsAttribute{
		Annotations: make([]Annotation, numAnnotations),
	}

	offset := 2
	for i := uint16(0); i < numAnnotations; i++ {
		ria.Annotations[i], offset = parseAnnotation(info, offset)
	}

	return ria
}

func parseRuntimeVisibleParameterAnnotationsAttribute(info []byte) *RuntimeVisibleParameterAnnotationsAttribute {
	if len(info) < 1 {
		return nil
	}

	numParameters := uint8(info[0])
	rvpa := &RuntimeVisibleParameterAnnotationsAttribute{
		ParameterAnnotations: make([][]Annotation, numParameters),
	}

	offset := 1
	for i := uint8(0); i < numParameters; i++ {
		if len(info) < offset+2 {
			return nil
		}
		numAnnotations := binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2

		annotations := make([]Annotation, numAnnotations)
		for j := uint16(0); j < numAnnotations; j++ {
			annotations[j], offset = parseAnnotation(info, offset)
		}
		rvpa.ParameterAnnotations[i] = annotations
	}

	return rvpa
}

func parseRuntimeInvisibleParameterAnnotationsAttribute(info []byte) *RuntimeInvisibleParameterAnnotationsAttribute {
	if len(info) < 1 {
		return nil
	}

	numParameters := uint8(info[0])
	ripa := &RuntimeInvisibleParameterAnnotationsAttribute{
		ParameterAnnotations: make([][]Annotation, numParameters),
	}

	offset := 1
	for i := uint8(0); i < numParameters; i++ {
		if len(info) < offset+2 {
			return nil
		}
		numAnnotations := binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2

		annotations := make([]Annotation, numAnnotations)
		for j := uint16(0); j < numAnnotations; j++ {
			annotations[j], offset = parseAnnotation(info, offset)
		}
		ripa.ParameterAnnotations[i] = annotations
	}

	return ripa
}

func parseTypeAnnotation(info []byte, offset int) (TypeAnnotation, int) {
	ta := TypeAnnotation{}
	if len(info) <= offset {
		return ta, offset
	}

	ta.TargetType = info[offset]
	offset++

	targetInfoStart := offset
	switch ta.TargetType {
	case 0x00, 0x01:
		offset += 1
	case 0x10:
		offset += 2
	case 0x11, 0x12:
		offset += 2
	case 0x13, 0x14, 0x15:
		// empty_target
	case 0x16:
		offset += 1
	case 0x17:
		offset += 2
	case 0x40, 0x41:
		if len(info) < offset+2 {
			return ta, offset
		}
		tableLength := binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2 + int(tableLength)*6
	case 0x42:
		offset += 2
	case 0x43, 0x44, 0x45, 0x46:
		offset += 2
	case 0x47, 0x48, 0x49, 0x4A, 0x4B:
		offset += 3
	}
	ta.TargetInfo = info[targetInfoStart:offset]

	if len(info) <= offset {
		return ta, offset
	}
	pathLength := int(info[offset])
	offset++

	ta.TargetPath = make([]TypePathEntry, pathLength)
	for i := 0; i < pathLength; i++ {
		if len(info) < offset+2 {
			return ta, offset
		}
		ta.TargetPath[i] = TypePathEntry{
			TypePathKind:      info[offset],
			TypeArgumentIndex: info[offset+1],
		}
		offset += 2
	}

	if len(info) < offset+4 {
		return ta, offset
	}
	ta.TypeIndex = binary.BigEndian.Uint16(info[offset : offset+2])
	numPairs := binary.BigEndian.Uint16(info[offset+2 : offset+4])
	offset += 4

	ta.ElementValuePairs = make([]ElementValuePair, numPairs)
	for i := uint16(0); i < numPairs; i++ {
		if len(info) < offset+2 {
			return ta, offset
		}
		pair := ElementValuePair{
			ElementNameIndex: binary.BigEndian.Uint16(info[offset : offset+2]),
		}
		offset += 2
		pair.Value, offset = parseElementValue(info, offset)
		ta.ElementValuePairs[i] = pair
	}

	return ta, offset
}

func parseRuntimeVisibleTypeAnnotationsAttribute(info []byte) *RuntimeVisibleTypeAnnotationsAttribute {
	if len(info) < 2 {
		return nil
	}

	numAnnotations := binary.BigEndian.Uint16(info[0:2])
	rvta := &RuntimeVisibleTypeAnnotationsAttribute{
		Annotations: make([]TypeAnnotation, numAnnotations),
	}

	offset := 2
	for i := uint16(0); i < numAnnotations; i++ {
		rvta.Annotations[i], offset = parseTypeAnnotation(info, offset)
	}

	return rvta
}

func parseRuntimeInvisibleTypeAnnotationsAttribute(info []byte) *RuntimeInvisibleTypeAnnotationsAttribute {
	if len(info) < 2 {
		return nil
	}

	numAnnotations := binary.BigEndian.Uint16(info[0:2])
	rita := &RuntimeInvisibleTypeAnnotationsAttribute{
		Annotations: make([]TypeAnnotation, numAnnotations),
	}

	offset := 2
	for i := uint16(0); i < numAnnotations; i++ {
		rita.Annotations[i], offset = parseTypeAnnotation(info, offset)
	}

	return rita
}

func parseAnnotationDefaultAttribute(info []byte) *AnnotationDefaultAttribute {
	if len(info) < 1 {
		return nil
	}

	ad := &AnnotationDefaultAttribute{}
	ad.DefaultValue, _ = parseElementValue(info, 0)
	return ad
}

func parseModuleAttribute(info []byte) *ModuleAttribute {
	if len(info) < 6 {
		return nil
	}

	m := &ModuleAttribute{
		ModuleNameIndex:    binary.BigEndian.Uint16(info[0:2]),
		ModuleFlags:        binary.BigEndian.Uint16(info[2:4]),
		ModuleVersionIndex: binary.BigEndian.Uint16(info[4:6]),
	}

	offset := 6

	if len(info) < offset+2 {
		return m
	}
	requiresCount := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	m.Requires = make([]ModuleRequires, requiresCount)
	for i := uint16(0); i < requiresCount; i++ {
		if len(info) < offset+6 {
			return m
		}
		m.Requires[i] = ModuleRequires{
			RequiresIndex:        binary.BigEndian.Uint16(info[offset : offset+2]),
			RequiresFlags:        binary.BigEndian.Uint16(info[offset+2 : offset+4]),
			RequiresVersionIndex: binary.BigEndian.Uint16(info[offset+4 : offset+6]),
		}
		offset += 6
	}

	if len(info) < offset+2 {
		return m
	}
	exportsCount := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	m.Exports = make([]ModuleExports, exportsCount)
	for i := uint16(0); i < exportsCount; i++ {
		if len(info) < offset+6 {
			return m
		}
		export := ModuleExports{
			ExportsIndex: binary.BigEndian.Uint16(info[offset : offset+2]),
			ExportsFlags: binary.BigEndian.Uint16(info[offset+2 : offset+4]),
		}
		exportsToCount := binary.BigEndian.Uint16(info[offset+4 : offset+6])
		offset += 6

		export.ExportsToIndex = make([]uint16, exportsToCount)
		for j := uint16(0); j < exportsToCount; j++ {
			if len(info) < offset+2 {
				return m
			}
			export.ExportsToIndex[j] = binary.BigEndian.Uint16(info[offset : offset+2])
			offset += 2
		}
		m.Exports[i] = export
	}

	if len(info) < offset+2 {
		return m
	}
	opensCount := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	m.Opens = make([]ModuleOpens, opensCount)
	for i := uint16(0); i < opensCount; i++ {
		if len(info) < offset+6 {
			return m
		}
		opens := ModuleOpens{
			OpensIndex: binary.BigEndian.Uint16(info[offset : offset+2]),
			OpensFlags: binary.BigEndian.Uint16(info[offset+2 : offset+4]),
		}
		opensToCount := binary.BigEndian.Uint16(info[offset+4 : offset+6])
		offset += 6

		opens.OpensToIndex = make([]uint16, opensToCount)
		for j := uint16(0); j < opensToCount; j++ {
			if len(info) < offset+2 {
				return m
			}
			opens.OpensToIndex[j] = binary.BigEndian.Uint16(info[offset : offset+2])
			offset += 2
		}
		m.Opens[i] = opens
	}

	if len(info) < offset+2 {
		return m
	}
	usesCount := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	m.Uses = make([]uint16, usesCount)
	for i := uint16(0); i < usesCount; i++ {
		if len(info) < offset+2 {
			return m
		}
		m.Uses[i] = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2
	}

	if len(info) < offset+2 {
		return m
	}
	providesCount := binary.BigEndian.Uint16(info[offset : offset+2])
	offset += 2

	m.Provides = make([]ModuleProvides, providesCount)
	for i := uint16(0); i < providesCount; i++ {
		if len(info) < offset+4 {
			return m
		}
		provides := ModuleProvides{
			ProvidesIndex: binary.BigEndian.Uint16(info[offset : offset+2]),
		}
		providesWithCount := binary.BigEndian.Uint16(info[offset+2 : offset+4])
		offset += 4

		provides.ProvidesWithIndex = make([]uint16, providesWithCount)
		for j := uint16(0); j < providesWithCount; j++ {
			if len(info) < offset+2 {
				return m
			}
			provides.ProvidesWithIndex[j] = binary.BigEndian.Uint16(info[offset : offset+2])
			offset += 2
		}
		m.Provides[i] = provides
	}

	return m
}

func parseModulePackagesAttribute(info []byte) *ModulePackagesAttribute {
	if len(info) < 2 {
		return nil
	}

	count := binary.BigEndian.Uint16(info[0:2])
	if len(info) < 2+int(count)*2 {
		return nil
	}

	mp := &ModulePackagesAttribute{
		PackageIndex: make([]uint16, count),
	}

	offset := 2
	for i := uint16(0); i < count; i++ {
		mp.PackageIndex[i] = binary.BigEndian.Uint16(info[offset : offset+2])
		offset += 2
	}

	return mp
}

func parseModuleMainClassAttribute(info []byte) *ModuleMainClassAttribute {
	if len(info) < 2 {
		return nil
	}
	return &ModuleMainClassAttribute{
		MainClassIndex: binary.BigEndian.Uint16(info[0:2]),
	}
}
