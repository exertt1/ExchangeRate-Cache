package router

import (
	"Excnahge-Cacher/handler"
	"context"
	"net/http"
	"strings"
)

type Router struct {
	routes  map[string]map[string]http.HandlerFunc
	handler *handler.HTTPHandler
}

func NewRouter(httpHandler *handler.HTTPHandler) *Router {
	router := &Router{
		routes:  make(map[string]map[string]http.HandlerFunc),
		handler: httpHandler,
	}
	router.registerRoutes()
	return router
}

func (r *Router) registerRoutes() {
	r.HandleFunc("GET", "rate", r.handler.GetCurrency)
}

func (r *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) matchPath(pattern, path string) bool {
	patternPart := strings.Split(strings.Trim(pattern, "/"), "/")
	pathPart := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternPart) != len(pathPart) {
		return false
	}
	for i := 0; i < len(pathPart); i++ {
		if strings.HasPrefix(patternPart[i], "{") && strings.HasSuffix(patternPart[i], "}") {
			continue
		}
		if strings.Contains(pathPart[i], "?") {
			pathPart[i] = pathPart[i][:strings.Index(pathPart[i], "?")]
		}
		if pathPart[i] != patternPart[i] {
			return false
		}
	}
	return true
}

func (r *Router) addPathParams(req *http.Request, path, pattern string) {
	patternPart := strings.Split(strings.Trim(pattern, "/"), "/")
	pathPart := strings.Split(strings.Trim(path, "/"), "/")

	ctx := req.Context()

	for i := 0; i < len(pathPart); i++ {
		if strings.HasPrefix(pathPart[i], "{") && strings.HasSuffix(pathPart[i], "}") {
			paramName := strings.Trim(patternPart[i], "{}")
			paramValue := pathPart[i]

			ctx = context.WithValue(ctx, paramName, paramValue)
		}
	}

	*req = *req.WithContext(ctx)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := req.URL.Path
	if handlers, ok := r.routes[req.Method]; ok {
		for routePath, handlerFunc := range handlers {
			if r.matchPath(routePath, path) {
				r.addPathParams(req, path, routePath)
				handlerFunc(w, req)
				return
			}
		}
	}
	http.NotFound(w, req)
}
