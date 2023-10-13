package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chronos-tachyon/go-acceptable"
	"github.com/rs/zerolog"
)

func Serve(w http.ResponseWriter, r *http.Request, code int, body []byte) {
	sum := sha256.Sum256(body)
	etagHeader := formatETag(sum)
	body, encodeHeader, hasEncodeHeader := compressData(r, body)
	lenHeader := strconv.FormatUint(uint64(len(body)), 10)

	h := w.Header()
	if hasEncodeHeader {
		h.Set("Content-Encoding", encodeHeader)
	}
	h.Set("Content-Length", lenHeader)
	h.Set("ETag", etagHeader)
	h.Add("Vary", "Accept-Encoding")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(body))
}

func OpenFile(fileName string) (fs.File, error) {
	if gDataFS != nil {
		file, err := gDataFS.Open(fileName)
		if err == nil {
			return file, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	return gBuiltInFS.Open(fileName)
}

func ServeFile(w http.ResponseWriter, r *http.Request, code int, file fs.File) {
	body, err := io.ReadAll(file)
	if err != nil {
		log := zerolog.Ctx(r.Context())
		log.Error().
			Err(err).
			Msg("failed to read file")
		const code = http.StatusInternalServerError
		text := http.StatusText(code)
		http.Error(w, text, code)
		return
	}

	Serve(w, r, code, body)
}

func ServeTemplate(w http.ResponseWriter, r *http.Request, code int, tmpl *template.Template, data *DataModel) {
	body, err := Render(tmpl, data)
	if err != nil {
		log := zerolog.Ctx(r.Context())
		log.Error().
			Err(err).
			Msg("failed to execute template")
		const code = http.StatusInternalServerError
		text := http.StatusText(code)
		http.Error(w, text, code)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	Serve(w, r, code, body)
}

func compressData(r *http.Request, in []byte) (out []byte, header string, ok bool) {
	out = in

	log := zerolog.Ctx(r.Context())

	acceptHeaders := r.Header.Values("Accept-Encoding")
	if len(acceptHeaders) <= 0 {
		return
	}
	acceptHeader := strings.Join(acceptHeaders, ", ")

	var parsedCodecs acceptable.List
	if err := parsedCodecs.Parse(acceptHeader, acceptable.AbsentSubValue); err != nil {
		log.Debug().
			Err(err).
			Str("value", acceptHeader).
			Msg("failed to parse Accept-Encoding header")
		return
	}

	a, found := acceptable.Negotiate(availableCodecs, parsedCodecs)
	if !found {
		return
	}

	name := a.Value
	codec, found := codecMap[name]
	if !found {
		return
	}

	tmp, err := codec.Compress(in)
	if err != nil {
		log.Debug().
			Err(err).
			Str("codec", name).
			Msg("failed to compress data")
		return
	}

	out = tmp
	return out, name, true
}

func formatETag(sum [sha256.Size]byte) string {
	const n = ((sha256.Size + 2) / 3) * 4
	tmp := make([]byte, n+2)
	tmp[0] = '"'
	base64.StdEncoding.Encode(tmp[1:], sum[:])
	tmp[n+1] = '"'
	return string(tmp)
}
