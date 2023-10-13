package html

import "html/template"

type A struct {
	Text           string             `json:"text"`
	Href           StringAttr         `json:"href"`
	HrefLang       StringAttr         `json:"hrefLang"`
	Type           StringAttr         `json:"type"`
	Rel            StringAttr         `json:"rel"`
	Ping           StringAttr         `json:"ping"`
	ReferrerPolicy ReferrerPolicyAttr `json:"referrerPolicy"`
}

func (a A) Render() template.HTML {
	return RenderText("a", a.Attributes(), a.Text)
}

func (a A) Attributes() []Attr {
	out := make([]Attr, 0, 6)
	out = a.Href.Append(out, "href", false)
	out = a.HrefLang.Append(out, "hreflang", false)
	out = a.Type.Append(out, "type", false)
	out = a.Rel.Append(out, "rel", false)
	out = a.Ping.Append(out, "ping", false)
	out = a.ReferrerPolicy.Append(out, "referrerpolicy")
	return out
}
