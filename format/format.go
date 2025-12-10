package format

import (
	"encoding"

	"github.com/dhamidi/sai/java"
)

type Encoder interface {
	encoding.TextMarshaler
	Encode(class *java.Class) error
}
