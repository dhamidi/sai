package parser

type Position struct {
	File   string
	Offset int
	Line   int
	Column int
}

type Span struct {
	Start Position
	End   Position
}

type TokenKind int

const (
	TokenEOF TokenKind = iota
	TokenError
	TokenWhitespace
	TokenComment
	TokenLineComment

	// Literals
	TokenIdent
	TokenIntLiteral
	TokenFloatLiteral
	TokenCharLiteral
	TokenStringLiteral
	TokenTextBlock
	TokenTrue
	TokenFalse
	TokenNull

	// Keywords
	TokenAbstract
	TokenAssert
	TokenBoolean
	TokenBreak
	TokenByte
	TokenCase
	TokenCatch
	TokenChar
	TokenClass
	TokenConst
	TokenContinue
	TokenDefault
	TokenDo
	TokenDouble
	TokenElse
	TokenEnum
	TokenExtends
	TokenFinal
	TokenFinally
	TokenFloat
	TokenFor
	TokenGoto
	TokenIf
	TokenImplements
	TokenImport
	TokenInstanceof
	TokenInt
	TokenInterface
	TokenLong
	TokenNative
	TokenNew
	TokenPackage
	TokenPrivate
	TokenProtected
	TokenPublic
	TokenReturn
	TokenShort
	TokenStatic
	TokenStrictfp
	TokenSuper
	TokenSwitch
	TokenSynchronized
	TokenThis
	TokenThrow
	TokenThrows
	TokenTransient
	TokenTry
	TokenVoid
	TokenVolatile
	TokenWhile

	// Contextual keywords
	TokenVar
	TokenYield
	TokenRecord
	TokenSealed
	TokenNonSealed
	TokenPermits
	TokenWhen
	TokenModule
	TokenOpen
	TokenRequires
	TokenExports
	TokenOpens
	TokenUses
	TokenProvides
	TokenTo
	TokenWith
	TokenTransitive

	// Operators and punctuation
	TokenLParen
	TokenRParen
	TokenLBrace
	TokenRBrace
	TokenLBracket
	TokenRBracket
	TokenSemicolon
	TokenComma
	TokenDot
	TokenEllipsis
	TokenAt
	TokenColonColon

	TokenAssign
	TokenEQ
	TokenNE
	TokenLT
	TokenLE
	TokenGT
	TokenGE
	TokenAnd
	TokenOr
	TokenNot
	TokenBitAnd
	TokenBitOr
	TokenBitXor
	TokenBitNot
	TokenShl
	TokenShr
	TokenUShr
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenPercent
	TokenIncrement
	TokenDecrement
	TokenQuestion
	TokenColon
	TokenArrow
	TokenPlusAssign
	TokenMinusAssign
	TokenStarAssign
	TokenSlashAssign
	TokenPercentAssign
	TokenAndAssign
	TokenOrAssign
	TokenXorAssign
	TokenShlAssign
	TokenShrAssign
	TokenUShrAssign
)

var tokenKindNames = map[TokenKind]string{
	TokenEOF:           "EOF",
	TokenError:         "Error",
	TokenWhitespace:    "Whitespace",
	TokenComment:       "Comment",
	TokenLineComment:   "LineComment",
	TokenIdent:         "Identifier",
	TokenIntLiteral:    "IntLiteral",
	TokenFloatLiteral:  "FloatLiteral",
	TokenCharLiteral:   "CharLiteral",
	TokenStringLiteral: "StringLiteral",
	TokenTextBlock:     "TextBlock",
	TokenTrue:          "true",
	TokenFalse:         "false",
	TokenNull:          "null",
	TokenAbstract:      "abstract",
	TokenAssert:        "assert",
	TokenBoolean:       "boolean",
	TokenBreak:         "break",
	TokenByte:          "byte",
	TokenCase:          "case",
	TokenCatch:         "catch",
	TokenChar:          "char",
	TokenClass:         "class",
	TokenConst:         "const",
	TokenContinue:      "continue",
	TokenDefault:       "default",
	TokenDo:            "do",
	TokenDouble:        "double",
	TokenElse:          "else",
	TokenEnum:          "enum",
	TokenExtends:       "extends",
	TokenFinal:         "final",
	TokenFinally:       "finally",
	TokenFloat:         "float",
	TokenFor:           "for",
	TokenGoto:          "goto",
	TokenIf:            "if",
	TokenImplements:    "implements",
	TokenImport:        "import",
	TokenInstanceof:    "instanceof",
	TokenInt:           "int",
	TokenInterface:     "interface",
	TokenLong:          "long",
	TokenNative:        "native",
	TokenNew:           "new",
	TokenPackage:       "package",
	TokenPrivate:       "private",
	TokenProtected:     "protected",
	TokenPublic:        "public",
	TokenReturn:        "return",
	TokenShort:         "short",
	TokenStatic:        "static",
	TokenStrictfp:      "strictfp",
	TokenSuper:         "super",
	TokenSwitch:        "switch",
	TokenSynchronized:  "synchronized",
	TokenThis:          "this",
	TokenThrow:         "throw",
	TokenThrows:        "throws",
	TokenTransient:     "transient",
	TokenTry:           "try",
	TokenVoid:          "void",
	TokenVolatile:      "volatile",
	TokenWhile:         "while",
	TokenVar:           "var",
	TokenYield:         "yield",
	TokenRecord:        "record",
	TokenSealed:        "sealed",
	TokenNonSealed:     "non-sealed",
	TokenPermits:       "permits",
	TokenWhen:          "when",
	TokenModule:        "module",
	TokenOpen:          "open",
	TokenRequires:      "requires",
	TokenExports:       "exports",
	TokenOpens:         "opens",
	TokenUses:          "uses",
	TokenProvides:      "provides",
	TokenTo:            "to",
	TokenWith:          "with",
	TokenTransitive:    "transitive",
	TokenLParen:        "(",
	TokenRParen:        ")",
	TokenLBrace:        "{",
	TokenRBrace:        "}",
	TokenLBracket:      "[",
	TokenRBracket:      "]",
	TokenSemicolon:     ";",
	TokenComma:         ",",
	TokenDot:           ".",
	TokenEllipsis:      "...",
	TokenAt:            "@",
	TokenColonColon:    "::",
	TokenAssign:        "=",
	TokenEQ:            "==",
	TokenNE:            "!=",
	TokenLT:            "<",
	TokenLE:            "<=",
	TokenGT:            ">",
	TokenGE:            ">=",
	TokenAnd:           "&&",
	TokenOr:            "||",
	TokenNot:           "!",
	TokenBitAnd:        "&",
	TokenBitOr:         "|",
	TokenBitXor:        "^",
	TokenBitNot:        "~",
	TokenShl:           "<<",
	TokenShr:           ">>",
	TokenUShr:          ">>>",
	TokenPlus:          "+",
	TokenMinus:         "-",
	TokenStar:          "*",
	TokenSlash:         "/",
	TokenPercent:       "%",
	TokenIncrement:     "++",
	TokenDecrement:     "--",
	TokenQuestion:      "?",
	TokenColon:         ":",
	TokenArrow:         "->",
	TokenPlusAssign:    "+=",
	TokenMinusAssign:   "-=",
	TokenStarAssign:    "*=",
	TokenSlashAssign:   "/=",
	TokenPercentAssign: "%=",
	TokenAndAssign:     "&=",
	TokenOrAssign:      "|=",
	TokenXorAssign:     "^=",
	TokenShlAssign:     "<<=",
	TokenShrAssign:     ">>=",
	TokenUShrAssign:    ">>>=",
}

func (k TokenKind) String() string {
	if name, ok := tokenKindNames[k]; ok {
		return name
	}
	return "Unknown"
}

type Token struct {
	Kind    TokenKind
	Span    Span
	Literal string
}

var keywords = map[string]TokenKind{
	"abstract":     TokenAbstract,
	"assert":       TokenAssert,
	"boolean":      TokenBoolean,
	"break":        TokenBreak,
	"byte":         TokenByte,
	"case":         TokenCase,
	"catch":        TokenCatch,
	"char":         TokenChar,
	"class":        TokenClass,
	"const":        TokenConst,
	"continue":     TokenContinue,
	"default":      TokenDefault,
	"do":           TokenDo,
	"double":       TokenDouble,
	"else":         TokenElse,
	"enum":         TokenEnum,
	"extends":      TokenExtends,
	"final":        TokenFinal,
	"finally":      TokenFinally,
	"float":        TokenFloat,
	"for":          TokenFor,
	"goto":         TokenGoto,
	"if":           TokenIf,
	"implements":   TokenImplements,
	"import":       TokenImport,
	"instanceof":   TokenInstanceof,
	"int":          TokenInt,
	"interface":    TokenInterface,
	"long":         TokenLong,
	"native":       TokenNative,
	"new":          TokenNew,
	"package":      TokenPackage,
	"private":      TokenPrivate,
	"protected":    TokenProtected,
	"public":       TokenPublic,
	"return":       TokenReturn,
	"short":        TokenShort,
	"static":       TokenStatic,
	"strictfp":     TokenStrictfp,
	"super":        TokenSuper,
	"switch":       TokenSwitch,
	"synchronized": TokenSynchronized,
	"this":         TokenThis,
	"throw":        TokenThrow,
	"throws":       TokenThrows,
	"transient":    TokenTransient,
	"try":          TokenTry,
	"void":         TokenVoid,
	"volatile":     TokenVolatile,
	"while":        TokenWhile,
	"true":         TokenTrue,
	"false":        TokenFalse,
	"null":         TokenNull,
	"var":          TokenVar,
	"yield":        TokenYield,
	"record":       TokenRecord,
	"sealed":       TokenSealed,
	"permits":      TokenPermits,
	"module":       TokenModule,
	"open":         TokenOpen,
	"requires":     TokenRequires,
	"exports":      TokenExports,
	"opens":        TokenOpens,
	"uses":         TokenUses,
	"provides":     TokenProvides,
	"to":           TokenTo,
	"with":         TokenWith,
	"transitive":   TokenTransitive,
}

func LookupKeyword(ident string) TokenKind {
	if kind, ok := keywords[ident]; ok {
		return kind
	}
	return TokenIdent
}
