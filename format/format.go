package format

import (
	"encoding"

	"github.com/dhamidi/javalyzer/java"
)

type Encoder interface {
	encoding.TextMarshaler
	Encode(class *java.Class) error
}
