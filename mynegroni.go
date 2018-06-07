// This is a simple implementation of the negroni
// middleware handler just to see if i understood the
// sourcecode correctly
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// Handler is an extension of the http.Hander interface. Instead of just
// demanding ServeHTTP with a writer and request it takes a http.Handler as
// the third argument. This handler is called after the Handler did its work.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.Handler)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, http.Handler)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.Handler) {
	h(w, r, next)
}

// middleware is a list of Handlers. It implements the http.Hander interface
// but can calls a Handler from its ServeHTTP function
type middleware struct {
	handler HandlerFunc
	next    *middleware
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r, m.next)
}

// MyNegroni handles a list of Middlewares. It implements the http.Handler.
type MyNegroni struct {
	middleware *middleware
}

// New creates a new Negroni instance
func New(handlers ...HandlerFunc) *MyNegroni {
	if len(handlers) == 0 {
		panic("na handlers given")
	}

	m := &middleware{}
	n := &MyNegroni{middleware: m}

	for i, h := range handlers {
		m.handler = h
		if i != len(handlers)-1 {
			// it it is the last handler let next be null
			m.next = &middleware{}
			m = m.next
		}
	}

	return n
}

func (n *MyNegroni) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n.middleware.ServeHTTP(w, r)
}

// Wraps a http.Handler to be used as a Handler
func WrapHandler(h http.Handler) HandlerFunc {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		h.ServeHTTP(w, r)
	})
}

func newLogger() HandlerFunc {
	return HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
		next http.Handler) {
		log.Println("this is my super logger")
		next.ServeHTTP(w, r)
	})
}

func myHandlefunc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "this is the end")
}

func main() {
	n := New(newLogger(), WrapHandler(http.HandlerFunc(myHandlefunc)))
	spew.Dump(n)

	http.ListenAndServe(":8080", n)
}
