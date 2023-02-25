package middleware

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
)

type decompressor struct {
}

// Decompress возвращает middleware для восстановления сжатого тела запроса,
// согласно заголовку Content-Encoding.
// Если переданный тип сжатия не поддерживается, middleware остановит обработку запроса
// со статусом 415 Unsupported Media Type.
func Decompress() func(next http.Handler) http.Handler {
	d := decompressor{}

	return d.handler
}

func (d decompressor) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "" {
			reader, err := d.createReader(r)
			if err != nil {
				w.WriteHeader(http.StatusUnsupportedMediaType)

				return
			}

			r.Body = reader
			defer reader.Close()
		}

		next.ServeHTTP(w, r)
	})
}

func (d decompressor) createReader(r *http.Request) (io.ReadCloser, error) {
	enc := d.chooseEncoding(r.Header.Get("Content-Encoding"))

	switch enc {
	case "gzip":
		return gzip.NewReader(r.Body)
	case "identity":
		return r.Body, nil
	}

	return nil, errors.New("not supported")
}

func (d decompressor) chooseEncoding(encodingsList string) string {
	encodings := strings.Split(strings.ReplaceAll(
		encodingsList,
		" ",
		"",
	), ",")
	if len(encodings) == 0 {
		return ""
	}

	return encodings[0]
}
