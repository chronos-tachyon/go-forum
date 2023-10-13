package html

import (
	"bytes"
	"fmt"
	"text/template"
)

type Attr struct {
	Name      string
	Value     string
	OmitValue bool
}

func (attr Attr) WriteTo(buf *bytes.Buffer) {
	buf.WriteByte(' ')
	buf.WriteString(attr.Name)

	if attr.OmitValue {
		return
	}

	buf.WriteByte('=')
	buf.WriteByte('"')
	template.HTMLEscape(buf, []byte(attr.Value))
	buf.WriteByte('"')
}

func (attr Attr) GoString() string {
	return fmt.Sprintf("html.Attr{%q, %q, %t}", attr.Name, attr.Value, attr.OmitValue)
}

func (attr Attr) String() string {
	if attr.OmitValue {
		return attr.Name
	}
	return fmt.Sprintf("%s=%q", attr.Name, attr.Value)
}

var (
	_ fmt.GoStringer = Attr{}
	_ fmt.Stringer   = Attr{}
)
