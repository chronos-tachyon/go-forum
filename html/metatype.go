package html

import (
	"encoding"
	"fmt"
	"strings"
)

type MetaType uint

const (
	Meta_None MetaType = iota
	Meta_HTTPEquiv
	Meta_Name
	Meta_ItemProp
)

const metaTypeSize = 4

var metaTypeGoNames = [metaTypeSize]string{
	"html.Meta_None",
	"html.Meta_HTTPEquiv",
	"html.Meta_Name",
	"html.Meta_ItemProp",
}

var metaTypeNames = [metaTypeSize]string{
	"none",
	"http-equiv",
	"name",
	"itemprop",
}

func (enum MetaType) IsValid() bool {
	return enum < metaTypeSize
}

func (enum MetaType) GoString() string {
	if enum.IsValid() {
		return metaTypeGoNames[enum]
	}
	return fmt.Sprintf("html.MetaType(%d)", uint(enum))
}

func (enum MetaType) String() string {
	if enum.IsValid() {
		return metaTypeNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid MetaType %d]", uint(enum))
}

func (enum MetaType) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}

func (enum *MetaType) Parse(input string) error {
	*enum = ^MetaType(0)
	for i, str := range metaTypeNames {
		if strings.EqualFold(str, input) {
			*enum = MetaType(i)
			return nil
		}
	}
	return fmt.Errorf("failed to parse %q as MetaType", input)
}

func (enum *MetaType) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

func (enum MetaType) IsPresent() bool {
	return enum > 0 && enum.IsValid()
}

func (enum MetaType) Append(out []Attr, name string) []Attr {
	if !enum.IsPresent() {
		return out
	}

	attrName := enum.String()
	return append(out, Attr{Name: attrName, Value: name})
}

var (
	_ fmt.GoStringer           = MetaType(0)
	_ fmt.Stringer             = MetaType(0)
	_ encoding.TextMarshaler   = MetaType(0)
	_ encoding.TextUnmarshaler = (*MetaType)(nil)
)
