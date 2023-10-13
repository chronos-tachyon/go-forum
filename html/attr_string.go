package html

import (
	"encoding"
	"fmt"
)

type StringAttr string

func (value StringAttr) GoString() string {
	return fmt.Sprintf("html.StringAttr(%q)", string(value))
}

func (value StringAttr) String() string {
	return string(value)
}

func (value StringAttr) MarshalText() ([]byte, error) {
	return []byte(value), nil
}

func (value *StringAttr) Parse(input string) error {
	*value = StringAttr(input)
	return nil
}

func (value *StringAttr) UnmarshalText(input []byte) error {
	return value.Parse(string(input))
}

func (value StringAttr) IsPresent() bool {
	return value != ""
}

func (value StringAttr) Append(out []Attr, attrName string, isMandatory bool) []Attr {
	if !isMandatory && !value.IsPresent() {
		return out
	}

	attrValue := value.String()
	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = StringAttr(0)
	_ fmt.Stringer             = StringAttr(0)
	_ encoding.TextMarshaler   = StringAttr(0)
	_ encoding.TextUnmarshaler = (*StringAttr)(nil)
)
