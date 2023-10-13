package main

import (
	"html/template"
)

type PostModel struct {
	Board  *BoardModel
	Topic  *TopicModel
	Thread *ThreadModel
	Body   string
}

func (post *PostModel) BodyHTML() template.HTML {
	return template.HTML(template.HTMLEscapeString(post.Body))
}
