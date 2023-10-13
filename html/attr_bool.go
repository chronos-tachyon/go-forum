package html

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type BoolAttr bool

var (
	kFalse = []byte("false")
	kTrue  = []byte("true")
	k0     = []byte("0")
	k1     = []byte("1")
)

var boolAttrMap = map[string]BoolAttr{
	"0":     false,
	"f":     false,
	"false": false,
	"n":     false,
	"no":    false,
	"off":   false,
	"1":     true,
	"t":     true,
	"true":  true,
	"y":     true,
	"yes":   true,
	"on":    true,
}

func (value BoolAttr) GoString() string {
	return fmt.Sprintf("html.BoolAttr(%t)", bool(value))
}

func (value BoolAttr) String() string {
	return strconv.FormatBool(bool(value))
}

func (value BoolAttr) MarshalText() ([]byte, error) {
	return []byte(value.String()), nil
}

func (value BoolAttr) MarshalJSON() ([]byte, error) {
	if value {
		return kTrue, nil
	}
	return kFalse, nil
}

func (value *BoolAttr) Parse(input string) error {
	*value = false

	if x, found := boolAttrMap[input]; found {
		*value = x
		return nil
	}

	lc := strings.ToLower(input)
	if x, found := boolAttrMap[lc]; found {
		*value = x
		return nil
	}

	return fmt.Errorf("failed to parse %q as BoolAttr", input)
}

func (value *BoolAttr) UnmarshalText(input []byte) error {
	return value.Parse(string(input))
}

func (value *BoolAttr) UnmarshalJSON(input []byte) error {
	*value = false

	if bytes.Equal(input, kFalse) || bytes.Equal(input, k0) {
		return nil
	}

	if bytes.Equal(input, kTrue) || bytes.Equal(input, k1) {
		*value = true
		return nil
	}

	var tmp bool
	if err := json.Unmarshal(input, &tmp); err != nil {
		return err
	}

	*value = BoolAttr(tmp)
	return nil
}

func (value BoolAttr) Append(out []Attr, attrName string) []Attr {
	if !value {
		return out
	}

	return append(out, Attr{Name: attrName, OmitValue: true})
}

func (value BoolAttr) AppendText(out []Attr, attrName string, attrValue string) []Attr {
	if !value {
		return out
	}

	return append(out, Attr{Name: attrName, Value: attrValue})
}

var (
	_ fmt.GoStringer           = BoolAttr(false)
	_ fmt.Stringer             = BoolAttr(false)
	_ encoding.TextMarshaler   = BoolAttr(false)
	_ encoding.TextUnmarshaler = (*BoolAttr)(nil)
	_ json.Marshaler           = BoolAttr(false)
	_ json.Unmarshaler         = (*BoolAttr)(nil)
)
