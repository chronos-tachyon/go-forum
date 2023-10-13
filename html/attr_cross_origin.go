package html

import (
	"encoding"
	"fmt"
	"strings"
)

type CrossOriginAttr byte

const (
	CrossOrigin_None CrossOriginAttr = iota
	CrossOrigin_Anonymous
	CrossOrigin_UseCredentials
)

const crossOriginAttrSize = 3

var crossOriginAttrGoNames = [crossOriginAttrSize]string{
	"html.CrossOrigin_None",
	"html.CrossOrigin_Anonymous",
	"html.CrossOrigin_UseCredentials",
}

var crossOriginAttrNames = [crossOriginAttrSize]string{
	"none",
	"anonymous",
	"use-credentials",
}

func (enum CrossOriginAttr) IsValid() bool {
	return enum < crossOriginAttrSize
}

func (enum CrossOriginAttr) GoString() string {
	if enum.IsValid() {
		return crossOriginAttrGoNames[enum]
	}
	return fmt.Sprintf("html.CrossOriginAttr(%d)", uint(enum))
}

func (enum CrossOriginAttr) String() string {
	if enum.IsValid() {
		return crossOriginAttrNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid CrossOriginAttr %d]", uint(enum))
}

func (enum CrossOriginAttr) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}

func (enum *CrossOriginAttr) Parse(input string) error {
	*enum = ^CrossOriginAttr(0)
	for i, str := range crossOriginAttrNames {
		if strings.EqualFold(str, input) {
			*enum = CrossOriginAttr(i)
			return nil
		}
	}
	return fmt.Errorf("failed to parse %q as CrossOriginAttr", input)
}

func (enum *CrossOriginAttr) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

func (enum CrossOriginAttr) IsPresent() bool {
	return enum > 0 && enum.IsValid()
}

func (enum CrossOriginAttr) Append(out []Attr, attrName string) []Attr {
	if !enum.IsPresent() {
		return out
	}

	attrValue := enum.String()
	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = CrossOriginAttr(0)
	_ fmt.Stringer             = CrossOriginAttr(0)
	_ encoding.TextMarshaler   = CrossOriginAttr(0)
	_ encoding.TextUnmarshaler = (*CrossOriginAttr)(nil)
)
