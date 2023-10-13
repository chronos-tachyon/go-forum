package main

import (
	"github.com/chronos-tachyon/go-forum/html"
)

type DataModel struct {
	Metas   []html.Meta
	Links   []html.Link
	Scripts []html.Script

	Site    *SiteModel
	User    *UserModel
	Board   *BoardModel
	Boards  []*BoardModel
	Topic   *TopicModel
	Topics  []*TopicModel
	Thread  *ThreadModel
	Threads []*ThreadModel
	Posts   []*PostModel
}
