package main

import "ireul.com/web"

// Render wraps web.Context
type Render struct {
	ctx *web.Context
}

// Data set data to ctx.Data
func (r Render) Data(key string, val interface{}) {
	r.ctx.Data[key] = val
}

// HTML renders a HTML
func (r Render) HTML(code int, t string) {
	r.ctx.HTML(code, t)
}

// JSON renders a JSON
func (r Render) JSON(code int, t interface{}) {
	r.ctx.JSON(code, t)
}

// Error renders a error string
func (r Render) Error(code int, t string) {
	r.ctx.PlainText(code, []byte("ERROR: "+t))
}

// Renderer mount renderer
func Renderer() interface{} {
	return func(ctx *web.Context) {
		ctx.Map(Render{ctx: ctx})
	}
}
