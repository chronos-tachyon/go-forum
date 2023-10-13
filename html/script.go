package html

import (
	"html/template"
)

type Script struct {
	Src            StringAttr         `json:"src"`
	Type           StringAttr         `json:"type"`
	FetchPriority  FetchPriorityAttr  `json:"fetchPriority"`
	CrossOrigin    CrossOriginAttr    `json:"crossOrigin"`
	ReferrerPolicy ReferrerPolicyAttr `json:"referrerPolicy"`
	Nonce          StringAttr         `json:"nonce"`
	Integrity      StringAttr         `json:"integrity"`
	Blocking       BoolAttr           `json:"blocking"`
	Async          BoolAttr           `json:"async"`
	Defer          BoolAttr           `json:"defer"`
	NoModule       BoolAttr           `json:"noModule"`
}

func (script Script) Render() template.HTML {
	return Render("script", true, script.Attributes())
}

func (script Script) Attributes() []Attr {
	out := make([]Attr, 0, 16)
	out = script.Src.Append(out, "src", true)
	out = script.Type.Append(out, "type", false)
	out = script.FetchPriority.Append(out, "fetchpriority")
	out = script.CrossOrigin.Append(out, "crossorigin")
	out = script.ReferrerPolicy.Append(out, "referrerpolicy")
	out = script.Nonce.Append(out, "nonce", false)
	out = script.Integrity.Append(out, "integrity", false)
	out = script.Blocking.AppendText(out, "blocking", "render")
	out = script.Defer.Append(out, "defer")
	out = script.Async.Append(out, "async")
	out = script.NoModule.Append(out, "nomodule")
	return out
}
