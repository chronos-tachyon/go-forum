package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

type Config struct {
	MediaTypes   map[string]string
	Translations map[string]string
	Site         SiteModel
	Templates    Templates
}

func (c *Config) Load(fsysName string, fsys fs.FS) error {
	*c = Config{}

	mediaTypes, err := LoadMap[string, string]("config/mediaTypes.json")
	if err != nil {
		return err
	}

	translations, err := LoadMap[string, string]("config/translations.json")
	if err != nil {
		return err
	}

	var site SiteModel
	err = LoadJSON(&site, "<built-in>", gBuiltInFS, "config/site.json")
	if err != nil {
		return err
	}

	if fsys != nil {
		var moreSite SiteModel
		err = LoadJSON(&moreSite, fsysName, fsys, "config/site.json")
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		site.Merge(moreSite)
	}

	var templates Templates
	err = templates.Load(fsys)
	if err != nil {
		return err
	}

	*c = Config{
		MediaTypes:   mediaTypes,
		Translations: translations,
		Site:         site,
		Templates:    templates,
	}
	return nil
}

func (c *Config) ResolveMediaType(fileName string) string {
	fileName = path.Base(fileName)
	for {
		if mediaType, ok := c.MediaTypes[fileName]; ok {
			return mediaType
		}

		if fileName == "" {
			break
		}

		// this can only happen on the first iteration
		if fileName[0] != '.' {
			i := strings.IndexByte(fileName, '.')
			if i < 0 {
				break
			}
			fileName = fileName[i:]
			continue
		}

		i := strings.IndexByte(fileName[1:], '.')
		if i < 0 {
			i = len(fileName)
		}

		fileName = fileName[i:]
	}
	return "application/octet-stream"
}

func (c *Config) Translate(key string) string {
	if str, ok := c.Translations[key]; ok {
		return str
	}
	return "${" + key + "}"
}

func LoadMap[K comparable, V any](fileName string) (map[K]V, error) {
	var data map[K]V
	err := LoadJSON(&data, "<built-in>", gBuiltInFS, fileName)
	if err != nil {
		return nil, err
	}

	var moreData map[K]V
	if gDataFS != nil {
		err = LoadJSON(&moreData, gDataDir, gDataFS, fileName)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range moreData {
		data[key] = value
	}

	return data, nil
}

func LoadJSON[T any](ptr *T, fsysName string, fsys fs.FS, fileName string) error {
	var empty T
	*ptr = empty

	file, err := fsys.Open(fileName)
	if err != nil {
		return fmt.Errorf("%s: %q: failed to open file: %w", fsysName, fileName, err)
	}

	raw, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("%s: %q: failed to read file: %w", fsysName, fileName, err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("%s: %q: failed to close file: %w", fsysName, fileName, err)
	}

	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	d.DisallowUnknownFields()
	err = d.Decode(ptr)
	if err != nil {
		*ptr = empty
		return fmt.Errorf("%s: %q: failed to decode JSON: %w", fsysName, fileName, err)
	}

	return nil
}
