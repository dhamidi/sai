package classfile

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

type reader struct {
	r   io.Reader
	err error
}

func (r *reader) readU1() uint8 {
	if r.err != nil {
		return 0
	}
	var buf [1]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	return buf[0]
}

func (r *reader) readU2() uint16 {
	if r.err != nil {
		return 0
	}
	var buf [2]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	return binary.BigEndian.Uint16(buf[:])
}

func (r *reader) readU4() uint32 {
	if r.err != nil {
		return 0
	}
	var buf [4]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	return binary.BigEndian.Uint32(buf[:])
}

func (r *reader) readBytes(n int) []byte {
	if r.err != nil {
		return nil
	}
	buf := make([]byte, n)
	_, r.err = io.ReadFull(r.r, buf)
	return buf
}

func ParseFile(path string) (*ClassFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open class file: %w", err)
	}
	defer f.Close()
	return Parse(f)
}

func Parse(rd io.Reader) (*ClassFile, error) {
	r := &reader{r: rd}

	magic := r.readU4()
	if r.err != nil {
		return nil, fmt.Errorf("failed to read magic: %w", r.err)
	}
	if magic != Magic {
		return nil, fmt.Errorf("invalid magic number: 0x%X (expected 0xCAFEBABE)", magic)
	}

	cf := &ClassFile{
		MinorVersion: r.readU2(),
		MajorVersion: r.readU2(),
	}
	if r.err != nil {
		return nil, fmt.Errorf("failed to read version: %w", r.err)
	}

	constantPoolCount := r.readU2()
	if r.err != nil {
		return nil, fmt.Errorf("failed to read constant pool count: %w", r.err)
	}

	cf.ConstantPool = make(ConstantPool, constantPoolCount-1)
	for i := uint16(1); i < constantPoolCount; i++ {
		entry, skip, err := readConstantPoolEntry(r)
		if err != nil {
			return nil, fmt.Errorf("failed to read constant pool entry %d: %w", i, err)
		}
		cf.ConstantPool[i-1] = entry
		if skip {
			i++
			if i < constantPoolCount {
				cf.ConstantPool[i-1] = nil
			}
		}
	}

	cf.AccessFlags = AccessFlags(r.readU2())
	cf.ThisClass = r.readU2()
	cf.SuperClass = r.readU2()

	interfacesCount := r.readU2()
	if r.err != nil {
		return nil, fmt.Errorf("failed to read class info: %w", r.err)
	}

	cf.Interfaces = make([]uint16, interfacesCount)
	for i := uint16(0); i < interfacesCount; i++ {
		cf.Interfaces[i] = r.readU2()
	}
	if r.err != nil {
		return nil, fmt.Errorf("failed to read interfaces: %w", r.err)
	}

	fieldsCount := r.readU2()
	if r.err != nil {
		return nil, fmt.Errorf("failed to read fields count: %w", r.err)
	}

	cf.Fields = make([]FieldInfo, fieldsCount)
	for i := uint16(0); i < fieldsCount; i++ {
		field, err := readFieldInfo(r, cf.ConstantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to read field %d: %w", i, err)
		}
		cf.Fields[i] = *field
	}

	methodsCount := r.readU2()
	if r.err != nil {
		return nil, fmt.Errorf("failed to read methods count: %w", r.err)
	}

	cf.Methods = make([]MethodInfo, methodsCount)
	for i := uint16(0); i < methodsCount; i++ {
		method, err := readMethodInfo(r, cf.ConstantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to read method %d: %w", i, err)
		}
		cf.Methods[i] = *method
	}

	attributesCount := r.readU2()
	if r.err != nil {
		return nil, fmt.Errorf("failed to read attributes count: %w", r.err)
	}

	cf.Attributes = make([]AttributeInfo, attributesCount)
	for i := uint16(0); i < attributesCount; i++ {
		attr, err := readAttributeInfo(r, cf.ConstantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to read attribute %d: %w", i, err)
		}
		cf.Attributes[i] = *attr
	}

	return cf, nil
}

func readConstantPoolEntry(r *reader) (ConstantPoolEntry, bool, error) {
	tag := ConstantTag(r.readU1())
	if r.err != nil {
		return nil, false, r.err
	}

	switch tag {
	case ConstantUtf8:
		length := r.readU2()
		bytes := r.readBytes(int(length))
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantUtf8Info{Value: decodeModifiedUtf8(bytes)}, false, nil

	case ConstantInteger:
		bytes := r.readU4()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantIntegerInfo{Value: int32(bytes)}, false, nil

	case ConstantFloat:
		bytes := r.readU4()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantFloatInfo{Value: math.Float32frombits(bytes)}, false, nil

	case ConstantLong:
		high := r.readU4()
		low := r.readU4()
		if r.err != nil {
			return nil, false, r.err
		}
		value := (int64(high) << 32) | int64(low)
		return &ConstantLongInfo{Value: value}, true, nil

	case ConstantDouble:
		high := r.readU4()
		low := r.readU4()
		if r.err != nil {
			return nil, false, r.err
		}
		bits := (uint64(high) << 32) | uint64(low)
		return &ConstantDoubleInfo{Value: math.Float64frombits(bits)}, true, nil

	case ConstantClass:
		nameIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantClassInfo{NameIndex: nameIndex}, false, nil

	case ConstantString:
		stringIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantStringInfo{StringIndex: stringIndex}, false, nil

	case ConstantFieldref:
		classIndex := r.readU2()
		nameAndTypeIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantFieldrefInfo{
			ClassIndex:       classIndex,
			NameAndTypeIndex: nameAndTypeIndex,
		}, false, nil

	case ConstantMethodref:
		classIndex := r.readU2()
		nameAndTypeIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantMethodrefInfo{
			ClassIndex:       classIndex,
			NameAndTypeIndex: nameAndTypeIndex,
		}, false, nil

	case ConstantInterfaceMethodref:
		classIndex := r.readU2()
		nameAndTypeIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantInterfaceMethodrefInfo{
			ClassIndex:       classIndex,
			NameAndTypeIndex: nameAndTypeIndex,
		}, false, nil

	case ConstantNameAndType:
		nameIndex := r.readU2()
		descriptorIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantNameAndTypeInfo{
			NameIndex:       nameIndex,
			DescriptorIndex: descriptorIndex,
		}, false, nil

	case ConstantMethodHandle:
		referenceKind := MethodHandleKind(r.readU1())
		referenceIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantMethodHandleInfo{
			ReferenceKind:  referenceKind,
			ReferenceIndex: referenceIndex,
		}, false, nil

	case ConstantMethodType:
		descriptorIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantMethodTypeInfo{
			DescriptorIndex: descriptorIndex,
		}, false, nil

	case ConstantDynamic:
		bootstrapMethodAttrIndex := r.readU2()
		nameAndTypeIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantDynamicInfo{
			BootstrapMethodAttrIndex: bootstrapMethodAttrIndex,
			NameAndTypeIndex:         nameAndTypeIndex,
		}, false, nil

	case ConstantInvokeDynamic:
		bootstrapMethodAttrIndex := r.readU2()
		nameAndTypeIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantInvokeDynamicInfo{
			BootstrapMethodAttrIndex: bootstrapMethodAttrIndex,
			NameAndTypeIndex:         nameAndTypeIndex,
		}, false, nil

	case ConstantModule:
		nameIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantModuleInfo{NameIndex: nameIndex}, false, nil

	case ConstantPackage:
		nameIndex := r.readU2()
		if r.err != nil {
			return nil, false, r.err
		}
		return &ConstantPackageInfo{NameIndex: nameIndex}, false, nil

	default:
		return nil, false, fmt.Errorf("unknown constant pool tag: %d", tag)
	}
}

func readFieldInfo(r *reader, cp ConstantPool) (*FieldInfo, error) {
	field := &FieldInfo{
		AccessFlags:     AccessFlags(r.readU2()),
		NameIndex:       r.readU2(),
		DescriptorIndex: r.readU2(),
	}

	attributesCount := r.readU2()
	if r.err != nil {
		return nil, r.err
	}

	field.Attributes = make([]AttributeInfo, attributesCount)
	for i := uint16(0); i < attributesCount; i++ {
		attr, err := readAttributeInfo(r, cp)
		if err != nil {
			return nil, err
		}
		field.Attributes[i] = *attr
	}

	return field, nil
}

func readMethodInfo(r *reader, cp ConstantPool) (*MethodInfo, error) {
	method := &MethodInfo{
		AccessFlags:     AccessFlags(r.readU2()),
		NameIndex:       r.readU2(),
		DescriptorIndex: r.readU2(),
	}

	attributesCount := r.readU2()
	if r.err != nil {
		return nil, r.err
	}

	method.Attributes = make([]AttributeInfo, attributesCount)
	for i := uint16(0); i < attributesCount; i++ {
		attr, err := readAttributeInfo(r, cp)
		if err != nil {
			return nil, err
		}
		method.Attributes[i] = *attr
	}

	return method, nil
}

func readAttributeInfo(r *reader, cp ConstantPool) (*AttributeInfo, error) {
	nameIndex := r.readU2()
	length := r.readU4()
	info := r.readBytes(int(length))
	if r.err != nil {
		return nil, r.err
	}

	attr := &AttributeInfo{
		NameIndex: nameIndex,
		Info:      info,
	}

	attrName := cp.GetUtf8(nameIndex)
	switch attrName {
	case "Code":
		attr.Parsed = parseCodeAttribute(info, cp)
	case "SourceFile":
		attr.Parsed = parseSourceFileAttribute(info)
	case "ConstantValue":
		attr.Parsed = parseConstantValueAttribute(info)
	case "Exceptions":
		attr.Parsed = parseExceptionsAttribute(info)
	case "InnerClasses":
		attr.Parsed = parseInnerClassesAttribute(info)
	case "Signature":
		attr.Parsed = parseSignatureAttribute(info)
	case "BootstrapMethods":
		attr.Parsed = parseBootstrapMethodsAttribute(info)
	case "LineNumberTable":
		attr.Parsed = parseLineNumberTableAttribute(info)
	case "LocalVariableTable":
		attr.Parsed = parseLocalVariableTableAttribute(info)
	case "EnclosingMethod":
		attr.Parsed = parseEnclosingMethodAttribute(info)
	case "Synthetic":
		attr.Parsed = parseSyntheticAttribute(info)
	case "Deprecated":
		attr.Parsed = parseDeprecatedAttribute(info)
	case "SourceDebugExtension":
		attr.Parsed = parseSourceDebugExtensionAttribute(info)
	case "LocalVariableTypeTable":
		attr.Parsed = parseLocalVariableTypeTableAttribute(info)
	case "MethodParameters":
		attr.Parsed = parseMethodParametersAttribute(info)
	case "NestHost":
		attr.Parsed = parseNestHostAttribute(info)
	case "NestMembers":
		attr.Parsed = parseNestMembersAttribute(info)
	case "Record":
		attr.Parsed = parseRecordAttribute(info, cp)
	case "PermittedSubclasses":
		attr.Parsed = parsePermittedSubclassesAttribute(info)
	case "StackMapTable":
		attr.Parsed = parseStackMapTableAttribute(info)
	case "RuntimeVisibleAnnotations":
		attr.Parsed = parseRuntimeVisibleAnnotationsAttribute(info)
	case "RuntimeInvisibleAnnotations":
		attr.Parsed = parseRuntimeInvisibleAnnotationsAttribute(info)
	case "RuntimeVisibleParameterAnnotations":
		attr.Parsed = parseRuntimeVisibleParameterAnnotationsAttribute(info)
	case "RuntimeInvisibleParameterAnnotations":
		attr.Parsed = parseRuntimeInvisibleParameterAnnotationsAttribute(info)
	case "RuntimeVisibleTypeAnnotations":
		attr.Parsed = parseRuntimeVisibleTypeAnnotationsAttribute(info)
	case "RuntimeInvisibleTypeAnnotations":
		attr.Parsed = parseRuntimeInvisibleTypeAnnotationsAttribute(info)
	case "AnnotationDefault":
		attr.Parsed = parseAnnotationDefaultAttribute(info)
	case "Module":
		attr.Parsed = parseModuleAttribute(info)
	case "ModulePackages":
		attr.Parsed = parseModulePackagesAttribute(info)
	case "ModuleMainClass":
		attr.Parsed = parseModuleMainClassAttribute(info)
	}

	return attr, nil
}

func decodeModifiedUtf8(bytes []byte) string {
	runes := make([]rune, 0, len(bytes))
	i := 0
	for i < len(bytes) {
		b := bytes[i]
		if b&0x80 == 0 {
			runes = append(runes, rune(b))
			i++
		} else if b&0xE0 == 0xC0 {
			if i+1 >= len(bytes) {
				break
			}
			r := rune(b&0x1F)<<6 | rune(bytes[i+1]&0x3F)
			runes = append(runes, r)
			i += 2
		} else if b&0xF0 == 0xE0 {
			if i+2 >= len(bytes) {
				break
			}
			r := rune(b&0x0F)<<12 | rune(bytes[i+1]&0x3F)<<6 | rune(bytes[i+2]&0x3F)
			if r >= 0xD800 && r <= 0xDBFF {
				if i+5 < len(bytes) && bytes[i+3] == 0xED {
					high := r
					low := rune(bytes[i+4]&0x0F)<<12 | rune(bytes[i+5]&0x3F)<<6 | rune(bytes[i+6]&0x3F)
					if low >= 0xDC00 && low <= 0xDFFF {
						r = 0x10000 + ((high - 0xD800) << 10) + (low - 0xDC00)
						runes = append(runes, r)
						i += 6
						continue
					}
				}
			}
			runes = append(runes, r)
			i += 3
		} else {
			runes = append(runes, rune(b))
			i++
		}
	}
	return string(runes)
}
