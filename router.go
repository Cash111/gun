package gun

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandleFunc
}

func splitPattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handleFunc HandleFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handleFunc
	parts := splitPattern(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
}

func (r *router) getRoute(method string, path string) (string, map[string]string) {
	params := make(map[string]string)
	pathParts := splitPattern(path)
	if root, ok := r.roots[method]; ok {
		if n := root.search(pathParts, 0); n != nil {
			patternParts := splitPattern(n.pattern)
			for index, part := range patternParts {
				if strings.HasPrefix(part, ":") {
					params[part[1:]] = pathParts[index]
				}
				if part[0] == '*' && len(part) > 1 {
					params[part[1:]] = strings.Join(pathParts[index:], "/")
					break
				}
			}
			return n.pattern, params
		}
	}
	return "", params
}

func (r *router) handle(ctx *Context) {
	route, params := r.getRoute(ctx.Method, ctx.Path)
	if route != "" {
		ctx.Params = params
		key := ctx.Method + "-" + route
		ctx.handlers = append(ctx.handlers, r.handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	ctx.Next()
}

func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}