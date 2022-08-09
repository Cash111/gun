package gun

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Params     map[string]string
	Path       string
	Method     string
	StatusCode int
	handlers    []HandleFunc
	index      int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:-1,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) GetHeader(key string) string {
	value := c.Req.Header.Get(key)
	return value
}

func (c *Context) Json(statusCode int, v interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(statusCode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(v); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(statusCode int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(statusCode)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) HTML(statusCode int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(statusCode)
	c.Writer.Write([]byte(html))
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
