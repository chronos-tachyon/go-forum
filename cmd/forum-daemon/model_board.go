package main

import (
	"fmt"
	"net/url"
)

type BoardModel struct {
	Slug        string
	Title       string
	Description string
}

func (board *BoardModel) URL() (*url.URL, error) {
	if board == nil {
		return nil, fmt.Errorf("*BoardModel is nil")
	}
	return gRouter.Get("Board").URL(
		"boardSlug", board.Slug,
	)
}
