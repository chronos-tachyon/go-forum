package html

import (
	"encoding"
	"fmt"
	"strings"
)

type FetchPriorityAttr byte

const (
	FetchPriority_None FetchPriorityAttr = iota
	FetchPriority_High
	FetchPriority_Low
	FetchPriority_Auto
)

const fetchPriorityAttrSize = 4

var fetchPriorityAttrGoNames = [fetchPriorityAttrSize]string{
	"html.FetchPriority_None",
	"html.FetchPriority_High",
	"html.FetchPriority_Low",
	"html.FetchPriority_Auto",
}

var fetchPriorityAttrNames = [fetchPriorityAttrSize]string{
	"none",
	"high",
	"low",
	"auto",
}

func (enum FetchPriorityAttr) IsValid() bool {
	return enum < fetchPriorityAttrSize
}

func (enum FetchPriorityAttr) GoString() string {
	if enum.IsValid() {
		return fetchPriorityAttrGoNames[enum]
	}
	return fmt.Sprintf("html.FetchPriorityAttr(%d)", uint(enum))
}

func (enum FetchPriorityAttr) String() string {
	if enum.IsValid() {
		return fetchPriorityAttrNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid FetchPriorityAttr %d]", uint(enum))
}

func (enum FetchPriorityAttr) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}

func (enum *FetchPriorityAttr) Parse(input string) error {
	*enum = ^FetchPriorityAttr(0)
	for i, str := range fetchPriorityAttrNames {
		if strings.EqualFold(str, input) {
			*enum = FetchPriorityAttr(i)
			return nil
		}
	}
	return fmt.Errorf("failed to parse %q as FetchPriorityAttr", input)
}

func (enum *FetchPriorityAttr) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

func (enum FetchPriorityAttr) IsPresent() bool {
	return enum > 0 && enum.IsValid()
}

func (enum FetchPriorityAttr) Append(out []Attr, attrName string) []Attr {
	if !enum.IsPresent() {
		return out
	}

	attrValue := enum.String()
	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = FetchPriorityAttr(0)
	_ fmt.Stringer             = FetchPriorityAttr(0)
	_ encoding.TextMarshaler   = FetchPriorityAttr(0)
	_ encoding.TextUnmarshaler = (*FetchPriorityAttr)(nil)
)
