package html

import (
	"bytes"
	"html/template"
	"sync"
)

var gPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

func Render(tagName string, tagClose bool, attrs []Attr) template.HTML {
	buf := gPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		gPool.Put(buf)
	}()

	buf.WriteByte('<')
	buf.WriteString(tagName)
	for _, attr := range attrs {
		attr.WriteTo(buf)
	}
	buf.WriteByte('>')
	if tagClose {
		buf.WriteByte('<')
		buf.WriteByte('/')
		buf.WriteString(tagName)
		buf.WriteByte('>')
	}
	return template.HTML(buf.String())
}

func RenderText(tagName string, attrs []Attr, tagContent string) template.HTML {
	buf := gPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		gPool.Put(buf)
	}()

	buf.WriteByte('<')
	buf.WriteString(tagName)
	for _, attr := range attrs {
		attr.WriteTo(buf)
	}
	buf.WriteByte('>')
	template.HTMLEscape(buf, []byte(tagContent))
	buf.WriteByte('<')
	buf.WriteByte('/')
	buf.WriteString(tagName)
	buf.WriteByte('>')
	return template.HTML(buf.String())
}
