package routes

import (
  "golang.org/x/net/context"
  "github.com/unrolled/render"
)

func NewRendererContext () context.Context {
  ctx := context.Background()
  render := render.New()
  return context.WithValue(ctx, "renderer", render)
}
