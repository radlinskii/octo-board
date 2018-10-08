package middleware

import (
	"io"
	"net/http"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
)

// CSSMinify returns a minified version of CSS files
type CSSMinify struct {
	Next http.Handler
}

func (cm CSSMinify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cm.Next == nil {
		cm.Next = http.DefaultServeMux
	}

	if !strings.Contains(r.URL.Path, ".css") {
		cm.Next.ServeHTTP(w, r)
		return
	}
	mf := minifyFilter("text/css", w)
	defer mf.Close()
	cm.Next.ServeHTTP(mf, r)
}

type minifyResponseWriter struct {
	http.ResponseWriter
	io.WriteCloser
}

func (m minifyResponseWriter) Write(b []byte) (int, error) {
	return m.WriteCloser.Write(b)
}

func minifyFilter(mediatype string, res http.ResponseWriter) minifyResponseWriter {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	mw := m.Writer(mediatype, res)
	return minifyResponseWriter{res, mw}
}
