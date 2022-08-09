package gun

import (
	"log"
	"net/http"
	"strings"
)

const (
	GET  = "GET"
	POST = "POST"
)

type HandleFunc func(ctx *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc
	parent      *RouterGroup
	engine      *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: nil,
		parent:      group,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) GET(pattern string, handleFunc HandleFunc) {
	group.addRoute(GET, pattern, handleFunc)
}

func (group *RouterGroup) POST(pattern string, handleFunc HandleFunc) {
	group.addRoute(POST, pattern, handleFunc)
}

func (group *RouterGroup) addRoute(method string, comp string, handleFunc HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handleFunc)
}

func (engine *Engine) GET(pattern string, handleFunc HandleFunc) {
	engine.router.addRoute(GET, pattern, handleFunc)
}

func (engine *Engine) POST(pattern string, handleFunc HandleFunc) {
	engine.router.addRoute(POST, pattern, handleFunc)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandleFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context := newContext(w, req)
	context.handlers = middlewares
	engine.router.handle(context)
}

func New() *Engine {
	engine := &Engine{
		router: &router{
			roots:    make(map[string]*node),
			handlers: make(map[string]HandleFunc),
		},
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
