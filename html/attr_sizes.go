package html

import (
	"encoding"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reSpace = regexp.MustCompile(`\s+`)
	reSize  = regexp.MustCompile(`^(?:any|\d+[Xx]\d+)$`)
)

type SizesAttr []string

func (value SizesAttr) IsValid() bool {
	for _, str := range value {
		if !reSize.MatchString(str) {
			return false
		}
	}
	return true
}

func (value SizesAttr) GoString() string {
	var out []byte
	out = append(out, "html.SizesAttr([]string{"...)
	for i, str := range value {
		if i > 0 {
			out = append(out, ", "...)
		}
		out = strconv.AppendQuote(out, str)
	}
	out = append(out, "})"...)
	return string(out)
}

func (value SizesAttr) String() string {
	return strings.Join(value, " ")
}

func (value *SizesAttr) Parse(input string) error {
	*value = nil
	return value.parseImpl(reSpace.Split(input, -1))
}

func (value *SizesAttr) parseImpl(list []string) error {
	for _, str := range list {
		if !reSize.MatchString(str) {
			return fmt.Errorf("failed to parse %q as image size")
		}
	}
	*value = SizesAttr(list)
	return nil
}

func (value SizesAttr) MarshalText() ([]byte, error) {
	return []byte(value.String()), nil
}

func (value *SizesAttr) UnmarshalText(input []byte) error {
	return value.Parse(string(input))
}

func (value SizesAttr) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(value))
}

func (value *SizesAttr) UnmarshalJSON(input []byte) error {
	*value = nil

	var str string
	err := json.Unmarshal(input, &str)
	if err == nil {
		return value.parseImpl(reSpace.Split(str, -1))
	}

	var list []string
	err = json.Unmarshal(input, &list)
	if err == nil {
		return value.parseImpl(list)
	}

	return err
}

func (value SizesAttr) IsPresent() bool {
	return len(value) > 0 && value.IsValid()
}

func (value SizesAttr) Append(out []Attr, attrName string) []Attr {
	if !value.IsPresent() {
		return out
	}

	attrValue := value.String()
	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = SizesAttr(nil)
	_ fmt.Stringer             = SizesAttr(nil)
	_ encoding.TextMarshaler   = SizesAttr(nil)
	_ encoding.TextUnmarshaler = (*SizesAttr)(nil)
	_ json.Marshaler           = SizesAttr(nil)
	_ json.Unmarshaler         = (*SizesAttr)(nil)
)
