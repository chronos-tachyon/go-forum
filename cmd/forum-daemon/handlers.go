package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func AddRoutes(router *mux.Router) {
	gRouter.
		Methods("GET").
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.HandlerFunc(StaticHandler)))
	gRouter.
		Methods("GET").
		Path("/robots.txt").
		HandlerFunc(StaticHandler)
	gRouter.
		Methods("GET").
		Path("/favicon.ico").
		HandlerFunc(StaticHandler)
	gRouter.
		Methods("GET").
		Path("/").
		HandlerFunc(RedirHomeHandler)
	gRouter.
		Methods("GET").
		Path("/u").
		HandlerFunc(RedirHomeHandler)

	gRouter.
		Methods("GET").
		Path("/login").
		HandlerFunc(LogInHandler).
		Name("LogIn")
	gRouter.
		Methods("GET").
		Path("/logout").
		HandlerFunc(LogOutHandler).
		Name("LogOut")

	gRouter.
		Methods("GET").
		Path("/f").
		HandlerFunc(HomeHandler).
		Name("Home")
	gRouter.
		Methods("GET").
		Path("/f/{boardSlug}").
		HandlerFunc(BoardHandler).
		Name("Board")
	gRouter.
		Methods("GET").
		Path("/f/{boardSlug}/{topicSlug}").
		HandlerFunc(TopicHandler).
		Name("Topic")
	gRouter.
		Methods("GET").
		Path("/f/{boardSlug}/{topicSlug}/{threadSlug}").
		HandlerFunc(ThreadHandler).
		Name("Thread")
	gRouter.
		Methods("GET").
		Path("/u/{userHandle}").
		HandlerFunc(UserHandler).
		Name("User")
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	fileName := path.Join("static", path.Clean(r.URL.Path))
	mediaType := gConfig.ResolveMediaType(fileName)

	file, err := OpenFile(fileName)
	if err != nil {
		log := zerolog.Ctx(r.Context())
		log.Error().
			Err(err).
			Str("fs", gDataDir).
			Str("path", fileName).
			Msg("failed to open file")
		const code = http.StatusInternalServerError
		text := http.StatusText(code)
		http.Error(w, text, code)
		return
	}

	w.Header().Set("Content-Type", mediaType)
	ServeFile(w, r, http.StatusOK, file)
}

func RedirHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/f", http.StatusMovedPermanently)
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	ServeTemplate(w, r, http.StatusOK, gConfig.Templates.LogIn, makeData())
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	user := makeUser("fred")
	ServeTemplate(w, r, http.StatusOK, gConfig.Templates.LogOut, makeData(user))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	user := makeUser("fred")
	boards := []*BoardModel{
		makeBoard("a"),
		makeBoard("b"),
		makeBoard("c"),
	}

	ServeTemplate(w, r, http.StatusOK, gConfig.Templates.Home, makeData(user, boards))
}

func BoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardSlug := vars["boardSlug"]
	user := makeUser("fred")
	board := makeBoard(boardSlug)
	topics := []*TopicModel{
		makeTopic(board, "d"),
		makeTopic(board, "e"),
		makeTopic(board, "f"),
	}

	ServeTemplate(w, r, http.StatusOK, gConfig.Templates.Board, makeData(user, board, topics))
}

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardSlug := vars["boardSlug"]
	topicSlug := vars["topicSlug"]
	user := makeUser("fred")
	board := makeBoard(boardSlug)
	topic := makeTopic(board, topicSlug)
	threads := []*ThreadModel{
		makeThread(topic, "g"),
		makeThread(topic, "h"),
		makeThread(topic, "i"),
	}

	ServeTemplate(w, r, http.StatusOK, gConfig.Templates.Topic, makeData(user, board, topic, threads))
}

func ThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardSlug := vars["boardSlug"]
	topicSlug := vars["topicSlug"]
	threadSlug := vars["threadSlug"]
	user := makeUser("fred")
	board := makeBoard(boardSlug)
	topic := makeTopic(board, topicSlug)
	thread := makeThread(topic, threadSlug)
	posts := []*PostModel{
		makePost(thread),
		makePost(thread),
		makePost(thread),
	}

	ServeTemplate(w, r, http.StatusOK, gConfig.Templates.Thread, makeData(user, board, topic, thread, posts))
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userHandle := vars["userHandle"]
	message := fmt.Sprintf("Profile %q not found", userHandle)
	http.Error(w, message, http.StatusNotFound)
}

func makeData(list ...any) *DataModel {
	data := &DataModel{Site: &gConfig.Site}
	for _, item := range list {
		switch x := item.(type) {
		case *SiteModel:
			data.Site = x
		case *UserModel:
			data.User = x
		case *BoardModel:
			data.Board = x
		case []*BoardModel:
			data.Boards = x
		case *TopicModel:
			data.Topic = x
		case []*TopicModel:
			data.Topics = x
		case *ThreadModel:
			data.Thread = x
		case []*ThreadModel:
			data.Threads = x
		case []*PostModel:
			data.Posts = x
		default:
			panic(fmt.Errorf("datum of unexpected type %T", item))
		}
	}
	return data
}

func makeUser(handle string) *UserModel {
	return &UserModel{
		Handle:      handle,
		Email:       handle + "@example.com",
		DisplayName: handle,
	}
}

func makeBoard(slug string) *BoardModel {
	return &BoardModel{
		Slug:        slug,
		Title:       slug,
		Description: "generic board description",
	}
}

func makeTopic(board *BoardModel, slug string) *TopicModel {
	return &TopicModel{
		Board:       board,
		Slug:        slug,
		Title:       slug,
		Description: "generic topic description",
	}
}

func makeThread(topic *TopicModel, slug string) *ThreadModel {
	return &ThreadModel{
		Board: topic.Board,
		Topic: topic,
		Slug:  slug,
		Title: slug,
	}
}

func makePost(thread *ThreadModel) *PostModel {
	return &PostModel{
		Board:  thread.Board,
		Topic:  thread.Topic,
		Thread: thread,
		Body:   "generic post body",
	}
}
