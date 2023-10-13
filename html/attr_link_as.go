package html

import (
	"encoding"
	"fmt"
	"strings"
)

type LinkAsAttr byte

const (
	LinkAs_None LinkAsAttr = iota
	LinkAs_Audio
	LinkAs_Document
	LinkAs_Embed
	LinkAs_Fetch
	LinkAs_Font
	LinkAs_Image
	LinkAs_Object
	LinkAs_Script
	LinkAs_Style
	LinkAs_Track
	LinkAs_Video
	LinkAs_Worker
)

const linkAsAttrSize = 13

var linkAsAttrGoNames = [linkAsAttrSize]string{
	"html.LinkAs_None",
	"html.LinkAs_Audio",
	"html.LinkAs_Document",
	"html.LinkAs_Embed",
	"html.LinkAs_Fetch",
	"html.LinkAs_Font",
	"html.LinkAs_Image",
	"html.LinkAs_Object",
	"html.LinkAs_Script",
	"html.LinkAs_Style",
	"html.LinkAs_Track",
	"html.LinkAs_Video",
	"html.LinkAs_Worker",
}

var linkAsAttrNames = [linkAsAttrSize]string{
	"none",
	"audio",
	"document",
	"embed",
	"fetch",
	"font",
	"image",
	"object",
	"script",
	"style",
	"track",
	"video",
	"worker",
}

func (enum LinkAsAttr) IsValid() bool {
	return enum < linkAsAttrSize
}

func (enum LinkAsAttr) GoString() string {
	if enum.IsValid() {
		return linkAsAttrGoNames[enum]
	}
	return fmt.Sprintf("html.LinkAsAttr(%d)", uint(enum))
}

func (enum LinkAsAttr) String() string {
	if enum.IsValid() {
		return linkAsAttrNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid LinkAsAttr %d]", uint(enum))
}

func (enum LinkAsAttr) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}

func (enum *LinkAsAttr) Parse(input string) error {
	*enum = ^LinkAsAttr(0)
	for i, str := range linkAsAttrNames {
		if strings.EqualFold(str, input) {
			*enum = LinkAsAttr(i)
			return nil
		}
	}
	return fmt.Errorf("failed to parse %q as LinkAsAttr", input)
}

func (enum *LinkAsAttr) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

func (enum LinkAsAttr) IsPresent() bool {
	return enum > 0 && enum.IsValid()
}

func (enum LinkAsAttr) Append(out []Attr, attrName string) []Attr {
	if !enum.IsPresent() {
		return out
	}

	attrValue := enum.String()
	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = LinkAsAttr(0)
	_ fmt.Stringer             = LinkAsAttr(0)
	_ encoding.TextMarshaler   = LinkAsAttr(0)
	_ encoding.TextUnmarshaler = (*LinkAsAttr)(nil)
)
