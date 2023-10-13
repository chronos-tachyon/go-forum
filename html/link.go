package html

import (
	"html/template"
)

type Link struct {
	Href           StringAttr         `json:"href"`
	HrefLang       StringAttr         `json:"hrefLang"`
	Type           StringAttr         `json:"type"`
	Rel            StringAttr         `json:"rel"`
	As             LinkAsAttr         `json:"as"`
	Title          StringAttr         `json:"title"`
	Media          StringAttr         `json:"media"`
	Sizes          SizesAttr          `json:"sizes"`
	ImageSizes     SizesAttr          `json:"imageSizes"`
	ImageSrcSet    StringAttr         `json:"imageSrcSet"`
	FetchPriority  FetchPriorityAttr  `json:"fetchPriority"`
	CrossOrigin    CrossOriginAttr    `json:"crossOrigin"`
	ReferrerPolicy ReferrerPolicyAttr `json:"referrerPolicy"`
	Integrity      StringAttr         `json:"integrity"`
	Blocking       BoolAttr           `json:"blocking"`
}

func (link Link) Render() template.HTML {
	return Render("link", false, link.Attributes())
}

func (link Link) Attributes() []Attr {
	out := make([]Attr, 0, 16)
	out = link.Href.Append(out, "href", true)
	out = link.HrefLang.Append(out, "hreflang", false)
	out = link.Type.Append(out, "type", false)
	out = link.Rel.Append(out, "rel", false)
	out = link.As.Append(out, "as")
	out = link.Title.Append(out, "title", false)
	out = link.Media.Append(out, "media", false)
	out = link.Sizes.Append(out, "sizes")
	out = link.ImageSizes.Append(out, "imagesizes")
	out = link.ImageSrcSet.Append(out, "imagesrcset", false)
	out = link.FetchPriority.Append(out, "fetchpriority")
	out = link.CrossOrigin.Append(out, "crossorigin")
	out = link.ReferrerPolicy.Append(out, "referrerpolicy")
	out = link.Integrity.Append(out, "integrity", false)
	out = link.Blocking.AppendText(out, "blocking", "render")
	return out
}
