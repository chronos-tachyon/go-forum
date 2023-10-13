package main

import (
	"net/url"
)

type UserModel struct {
	Handle      string
	DisplayName string
	Email       string
}

func (user *UserModel) URL() (*url.URL, error) {
	return gRouter.Get("User").URL("userHandle", user.Handle)
}
