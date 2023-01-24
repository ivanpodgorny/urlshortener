package router

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	routes map[string][]patternHandler
}

func New() *Router {
	return &Router{
		routes: map[string][]patternHandler{},
	}
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func NotAllowed(w http.ResponseWriter) {
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

func ServerError(w http.ResponseWriter) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, ...string)

type patternHandler struct {
	method  string
	handler HandlerFunc
}

var (
	errNotAllowed = errors.New("not allowed")
	errNotFound   = errors.New("not found")
)

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, params, err := r.findHandler(req.URL.Path, req.Method)
	if err == nil {
		handler(w, req, params...)
	} else {
		if errors.Is(err, errNotFound) {
			http.NotFound(w, req)
		} else if errors.Is(err, errNotAllowed) {
			NotAllowed(w)
		} else {
			BadRequest(w)
		}
	}
}

func (r Router) Add(method string, pattern string, handler HandlerFunc) {
	preparedPattern := r.preparePattern(pattern)
	r.routes[preparedPattern] = append(r.routes[preparedPattern], patternHandler{
		method:  method,
		handler: handler,
	})
}

func (r Router) findHandler(url string, method string) (HandlerFunc, []string, error) {
	for p, handlers := range r.routes {
		matched, err := regexp.MatchString(p, url)
		if err != nil || !matched {
			continue
		}

		for _, h := range handlers {
			if h.method == method {
				re := regexp.MustCompile(p)

				return h.handler, re.FindStringSubmatch(url)[1:], nil
			}
		}

		return nil, []string{}, errNotAllowed
	}

	return nil, []string{}, errNotFound
}

func (r Router) preparePattern(pattern string) string {
	return fmt.Sprintf(`^\/(%s)$`, strings.Trim(pattern, "/"))
}
