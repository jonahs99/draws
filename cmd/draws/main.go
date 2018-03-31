package main

import (
	"github.com/jonahs99/draws"
)

func main() {
	canvas := draws.New(app)
	canvas.Serve()
}

func app(ctx draws.Context) {
	ctx.BackgroundStyle("#222")
	ctx.StrokeStyle("#eee")
	ctx.Rect(20, 30, 100, 140)
	ctx.Stroke()
}
