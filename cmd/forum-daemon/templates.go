package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"
	"sync"
)

var gRenderPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

type Templates struct {
	Base   *template.Template
	Home   *template.Template
	LogIn  *template.Template
	LogOut *template.Template
	Board  *template.Template
	Topic  *template.Template
	Thread *template.Template
}

func Render(t *template.Template, data *DataModel) ([]byte, error) {
	buf := gRenderPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		gRenderPool.Put(buf)
	}()

	if err := t.Execute(buf, data); err != nil {
		return nil, err
	}

	shared := buf.Bytes()
	result := make([]byte, len(shared))
	copy(result, shared)
	return result, nil
}

func (templates *Templates) Load(fsys fs.FS) error {
	*templates = Templates{}

	baseTemplate, err := createBaseTemplate(fsys)
	if err != nil {
		return err
	}

	homeTemplate, err := createTemplate(baseTemplate, fsys, "home.tmpl")
	if err != nil {
		return err
	}

	logInTemplate, err := createTemplate(baseTemplate, fsys, "login.tmpl")
	if err != nil {
		return err
	}

	logOutTemplate, err := createTemplate(baseTemplate, fsys, "logout.tmpl")
	if err != nil {
		return err
	}

	boardTemplate, err := createTemplate(baseTemplate, fsys, "board.tmpl")
	if err != nil {
		return err
	}

	topicTemplate, err := createTemplate(baseTemplate, fsys, "topic.tmpl")
	if err != nil {
		return err
	}

	threadTemplate, err := createTemplate(baseTemplate, fsys, "thread.tmpl")
	if err != nil {
		return err
	}

	*templates = Templates{
		Base:   baseTemplate,
		Home:   homeTemplate,
		LogIn:  logInTemplate,
		LogOut: logOutTemplate,
		Board:  boardTemplate,
		Topic:  topicTemplate,
		Thread: threadTemplate,
	}
	return nil
}

func createBaseTemplate(fsys fs.FS) (*template.Template, error) {
	t := template.New("root")
	t.Option("missingkey=error")
	t.Funcs(kFuncMap)

	if err := mergeTemplateIncludes(t, gBuiltInFS); err != nil {
		return nil, fmt.Errorf("built-in template: %w", err)
	}

	if fsys != nil {
		if err := mergeTemplateIncludes(t, fsys); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func createTemplate(baseTemplate *template.Template, fsys fs.FS, pathName string) (*template.Template, error) {
	pathName = path.Join("templates", pathName)

	t, err := baseTemplate.Clone()
	if err != nil {
		panic(err)
	}

	err = mergeTemplate(t, gBuiltInFS, pathName)
	if err != nil {
		return nil, fmt.Errorf("built-in template: %w", err)
	}

	if fsys != nil {
		err = mergeTemplate(t, fsys, pathName)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	return t, nil
}

func mergeTemplateIncludes(t *template.Template, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(pathName string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasPrefix(pathName, "templates/include/") && strings.HasSuffix(pathName, ".tmpl") {
			return mergeTemplate(t, fsys, pathName)
		}

		return nil
	})
}

func mergeTemplate(t *template.Template, fsys fs.FS, pathName string) error {
	raw, err := fs.ReadFile(fsys, pathName)
	if err != nil {
		return fmt.Errorf("%q: failed to read file: %w", pathName, err)
	}

	_, err = t.Parse(string(raw))
	if err != nil {
		return fmt.Errorf("%q: failed to parse file as Go HTML template: %w", pathName, err)
	}

	return nil
}

var kFuncMap = template.FuncMap(nil)
