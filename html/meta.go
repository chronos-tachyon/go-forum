package html

import (
	"html/template"
)

type Meta struct {
	Type    MetaType   `json:"type"`
	Name    StringAttr `json:"name"`
	Content StringAttr `json:"content"`
	Media   StringAttr `json:"media"`
}

func (meta Meta) Render() template.HTML {
	return Render("meta", false, meta.Attributes())
}

func (meta Meta) Attributes() []Attr {
	out := make([]Attr, 0, 3)
	if meta.Type.IsPresent() {
		attrName := meta.Type.String()
		out = meta.Name.Append(out, attrName, true)
	}
	out = meta.Content.Append(out, "content", false)
	out = meta.Media.Append(out, "media", false)
	return out
}
