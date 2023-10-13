package main

import (
	"fmt"
	"net/url"
)

type TopicModel struct {
	Board       *BoardModel
	Slug        string
	Title       string
	Description string
}

func (topic *TopicModel) URL() (*url.URL, error) {
	if topic == nil {
		return nil, fmt.Errorf("*TopicModel is nil")
	}
	return gRouter.Get("Topic").URL(
		"boardSlug", topic.Board.Slug,
		"topicSlug", topic.Slug,
	)
}
