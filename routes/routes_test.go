package routes

import (
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

func NewRendererContext() context.Context {
	ctx := context.Background()
	render := render.New()
	return context.WithValue(ctx, "renderer", render)
}
