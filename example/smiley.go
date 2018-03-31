package main

import (
	"math"

	"github.com/jonahs99/draws"
)

func app(ctx draws.Context) {
	ctx.Size(200, 200)
	ctx.BackgroundStyle("#eee")

	ctx.TranslateCenter()

	ctx.FillStyle("yellow")
	ctx.StrokeStyle("black")
	ctx.LineWidth(5)

	ctx.BeginPath()
	ctx.Circle(0, 0, 60)
	ctx.FillStroke()

	ctx.BeginPath()
	ctx.Arc(0, 0, 40, 0, math.Pi)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.Circle(-22, -15, 10)
	ctx.Circle(22, -15, 10)
	ctx.FillStyle("black")
	ctx.Fill()
}

func main() {
	canvas := draws.New(app)
	canvas.Serve(":3000")
}