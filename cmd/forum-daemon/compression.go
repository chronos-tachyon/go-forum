package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"sync"

	"github.com/chronos-tachyon/go-acceptable"
)

type flateWriterPair struct {
	buffer *bytes.Buffer
	writer *flate.Writer
}

type gzipWriterPair struct {
	buffer *bytes.Buffer
	writer *gzip.Writer
}

var gDeflatePool = sync.Pool{
	New: func() any {
		buffer := bytes.NewBuffer(make([]byte, 0, 4096))
		writer, err := flate.NewWriter(buffer, flate.DefaultCompression)
		if err != nil {
			panic(err)
		}
		return flateWriterPair{buffer: buffer, writer: writer}
	},
}

var gGzipPool = sync.Pool{
	New: func() any {
		buffer := bytes.NewBuffer(make([]byte, 0, 4096))
		writer, err := gzip.NewWriterLevel(buffer, gzip.DefaultCompression)
		if err != nil {
			panic(err)
		}
		return gzipWriterPair{buffer: buffer, writer: writer}
	},
}

type compressionCodec interface {
	Compress([]byte) ([]byte, error)
}

type identityCodec struct{}

func (*identityCodec) Compress(in []byte) ([]byte, error) {
	return in, nil
}

type deflateCodec struct{}

func (*deflateCodec) Compress(in []byte) ([]byte, error) {
	pair := gDeflatePool.Get().(flateWriterPair)
	defer func() {
		pair.buffer.Reset()
		pair.writer.Reset(pair.buffer)
		gDeflatePool.Put(pair)
	}()

	_, err := pair.writer.Write(in)
	if err != nil {
		return nil, err
	}

	err = pair.writer.Close()
	if err != nil {
		return nil, err
	}

	data := pair.buffer.Bytes()
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

type gzipCodec struct{}

func (*gzipCodec) Compress(in []byte) ([]byte, error) {
	pair := gGzipPool.Get().(gzipWriterPair)
	defer func() {
		pair.buffer.Reset()
		pair.writer.Reset(pair.buffer)
		gGzipPool.Put(pair)
	}()

	_, err := pair.writer.Write(in)
	if err != nil {
		return nil, err
	}

	err = pair.writer.Close()
	if err != nil {
		return nil, err
	}

	data := pair.buffer.Bytes()
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

var codecMap = map[string]compressionCodec{
	"identity":  (*identityCodec)(nil),
	"deflate":   (*deflateCodec)(nil),
	"x-deflate": (*deflateCodec)(nil),
	"gzip":      (*gzipCodec)(nil),
	"x-gzip":    (*gzipCodec)(nil),
}

var availableCodecs acceptable.List

func init() {
	list := make(acceptable.List, 0, len(codecMap))
	for name := range codecMap {
		list = append(list, acceptable.Acceptable{Value: name, Quality: 1000})
	}
	availableCodecs = list
}
