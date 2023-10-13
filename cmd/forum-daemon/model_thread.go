package main

import (
	"fmt"
	"net/url"
)

type ThreadModel struct {
	Board *BoardModel
	Topic *TopicModel
	Slug  string
	Title string
}

func (thread *ThreadModel) URL() (*url.URL, error) {
	if thread == nil {
		return nil, fmt.Errorf("*ThreadModel is nil")
	}
	return gRouter.Get("Thread").URL(
		"boardSlug", thread.Board.Slug,
		"topicSlug", thread.Topic.Slug,
		"threadSlug", thread.Slug,
	)
}
