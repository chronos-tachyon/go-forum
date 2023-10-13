package html

import (
	"encoding"
	"fmt"
	"strings"
)

type ReferrerPolicyAttr byte

const (
	ReferrerPolicy_None ReferrerPolicyAttr = iota
	ReferrerPolicy_NoReferrer
	ReferrerPolicy_NoReferrerWhenDowngrade
	ReferrerPolicy_Origin
	ReferrerPolicy_OriginWhenCrossOrigin
	ReferrerPolicy_SameOrigin
	ReferrerPolicy_StrictOrigin
	ReferrerPolicy_StrictOriginWhenCrossOrigin
	ReferrerPolicy_UnsafeURL
)

const referrerPolicyAttrSize = 9

var referrerPolicyAttrGoNames = [referrerPolicyAttrSize]string{
	"html.ReferrerPolicy_None",
	"html.ReferrerPolicy_NoReferrer",
	"html.ReferrerPolicy_NoReferrerWhenDowngrade",
	"html.ReferrerPolicy_Origin",
	"html.ReferrerPolicy_OriginWhenCrossOrigin",
	"html.ReferrerPolicy_SameOrigin",
	"html.ReferrerPolicy_StrictOrigin",
	"html.ReferrerPolicy_StrictOriginWhenCrossOrigin",
	"html.ReferrerPolicy_UnsafeURL",
}

var referrerPolicyAttrNames = [referrerPolicyAttrSize]string{
	"none",
	"no-referrer",
	"no-referrer-when-downgrade",
	"origin",
	"origin-when-cross-origin",
	"same-origin",
	"strict-origin",
	"strict-origin-when-cross-origin",
	"unsafe-url",
}

func (enum ReferrerPolicyAttr) IsValid() bool {
	return enum < referrerPolicyAttrSize
}

func (enum ReferrerPolicyAttr) GoString() string {
	if enum.IsValid() {
		return referrerPolicyAttrGoNames[enum]
	}
	return fmt.Sprintf("html.ReferrerPolicyAttr(%d)", uint(enum))
}

func (enum ReferrerPolicyAttr) String() string {
	if enum.IsValid() {
		return referrerPolicyAttrNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid ReferrerPolicyAttr %d]", uint(enum))
}

func (enum ReferrerPolicyAttr) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}

func (enum *ReferrerPolicyAttr) Parse(input string) error {
	*enum = ^ReferrerPolicyAttr(0)
	for i, str := range referrerPolicyAttrNames {
		if strings.EqualFold(str, input) {
			*enum = ReferrerPolicyAttr(i)
			return nil
		}
	}
	return fmt.Errorf("failed to parse %q as ReferrerPolicyAttr", input)
}

func (enum *ReferrerPolicyAttr) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

func (enum ReferrerPolicyAttr) IsPresent() bool {
	return enum > 0 && enum.IsValid()
}

func (enum ReferrerPolicyAttr) Append(out []Attr, attrName string) []Attr {
	if !enum.IsPresent() {
		return out
	}

	attrValue := enum.String()
	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = ReferrerPolicyAttr(0)
	_ fmt.Stringer             = ReferrerPolicyAttr(0)
	_ encoding.TextMarshaler   = ReferrerPolicyAttr(0)
	_ encoding.TextUnmarshaler = (*ReferrerPolicyAttr)(nil)
)
