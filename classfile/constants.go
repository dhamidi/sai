package classfile

const (
	Magic = 0xCAFEBABE
)

type AccessFlags uint16

const (
	AccPublic       AccessFlags = 0x0001
	AccPrivate      AccessFlags = 0x0002
	AccProtected    AccessFlags = 0x0004
	AccStatic       AccessFlags = 0x0008
	AccFinal        AccessFlags = 0x0010
	AccSuper        AccessFlags = 0x0020
	AccSynchronized AccessFlags = 0x0020
	AccVolatile     AccessFlags = 0x0040
	AccBridge       AccessFlags = 0x0040
	AccTransient    AccessFlags = 0x0080
	AccVarargs      AccessFlags = 0x0080
	AccNative       AccessFlags = 0x0100
	AccInterface    AccessFlags = 0x0200
	AccAbstract     AccessFlags = 0x0400
	AccStrict       AccessFlags = 0x0800
	AccSynthetic    AccessFlags = 0x1000
	AccAnnotation   AccessFlags = 0x2000
	AccEnum         AccessFlags = 0x4000
	AccModule       AccessFlags = 0x8000
)

func (f AccessFlags) IsPublic() bool       { return f&AccPublic != 0 }
func (f AccessFlags) IsPrivate() bool      { return f&AccPrivate != 0 }
func (f AccessFlags) IsProtected() bool    { return f&AccProtected != 0 }
func (f AccessFlags) IsStatic() bool       { return f&AccStatic != 0 }
func (f AccessFlags) IsFinal() bool        { return f&AccFinal != 0 }
func (f AccessFlags) IsSuper() bool        { return f&AccSuper != 0 }
func (f AccessFlags) IsSynchronized() bool { return f&AccSynchronized != 0 }
func (f AccessFlags) IsVolatile() bool     { return f&AccVolatile != 0 }
func (f AccessFlags) IsBridge() bool       { return f&AccBridge != 0 }
func (f AccessFlags) IsTransient() bool    { return f&AccTransient != 0 }
func (f AccessFlags) IsVarargs() bool      { return f&AccVarargs != 0 }
func (f AccessFlags) IsNative() bool       { return f&AccNative != 0 }
func (f AccessFlags) IsInterface() bool    { return f&AccInterface != 0 }
func (f AccessFlags) IsAbstract() bool     { return f&AccAbstract != 0 }
func (f AccessFlags) IsStrict() bool       { return f&AccStrict != 0 }
func (f AccessFlags) IsSynthetic() bool    { return f&AccSynthetic != 0 }
func (f AccessFlags) IsAnnotation() bool   { return f&AccAnnotation != 0 }
func (f AccessFlags) IsEnum() bool         { return f&AccEnum != 0 }
func (f AccessFlags) IsModule() bool       { return f&AccModule != 0 }

type ConstantTag uint8

const (
	ConstantUtf8               ConstantTag = 1
	ConstantInteger            ConstantTag = 3
	ConstantFloat              ConstantTag = 4
	ConstantLong               ConstantTag = 5
	ConstantDouble             ConstantTag = 6
	ConstantClass              ConstantTag = 7
	ConstantString             ConstantTag = 8
	ConstantFieldref           ConstantTag = 9
	ConstantMethodref          ConstantTag = 10
	ConstantInterfaceMethodref ConstantTag = 11
	ConstantNameAndType        ConstantTag = 12
	ConstantMethodHandle       ConstantTag = 15
	ConstantMethodType         ConstantTag = 16
	ConstantDynamic            ConstantTag = 17
	ConstantInvokeDynamic      ConstantTag = 18
	ConstantModule             ConstantTag = 19
	ConstantPackage            ConstantTag = 20
)

type MethodHandleKind uint8

const (
	RefGetField         MethodHandleKind = 1
	RefGetStatic        MethodHandleKind = 2
	RefPutField         MethodHandleKind = 3
	RefPutStatic        MethodHandleKind = 4
	RefInvokeVirtual    MethodHandleKind = 5
	RefInvokeStatic     MethodHandleKind = 6
	RefInvokeSpecial    MethodHandleKind = 7
	RefNewInvokeSpecial MethodHandleKind = 8
	RefInvokeInterface  MethodHandleKind = 9
)
