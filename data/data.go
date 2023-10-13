package data

import (
	"embed"
	"io/fs"
)

//go:embed config/*
//go:embed static/*
//go:embed templates/*
var fsys embed.FS

func FS() fs.FS {
	return fsys
}
