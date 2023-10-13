package main

import (
	"github.com/chronos-tachyon/go-forum/html"
)

type SiteModel struct {
	Lang     string        `json:"lang"`
	Viewport string        `json:"viewport"`
	Title    string        `json:"title"`
	Metas    []html.Meta   `json:"metas"`
	Links    []html.Link   `json:"links"`
	Scripts  []html.Script `json:"scripts"`
}

func (site *SiteModel) Merge(other SiteModel) {
	if other.Lang != "" {
		site.Lang = other.Lang
	}
	if other.Viewport != "" {
		site.Viewport = other.Viewport
	}
	if other.Title != "" {
		site.Title = other.Title
	}
	if other.Metas != nil {
		site.Metas = other.Metas
	}
	if other.Links != nil {
		site.Links = other.Links
	}
	if other.Scripts != nil {
		site.Scripts = other.Scripts
	}
}
