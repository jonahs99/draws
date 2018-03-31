`draws` is a simple go library for drawing to an HTML5 canvas over websockets.

Example:
```go
import "github.com/jonahs99/draws"

func app(ctx draws.Context) {
	ctx.BackgroundStyle("#333")
	ctx.StrokeStyle("#f06")
	ctx.Rect(100, 100, 200, 200)
	ctx.Stroke()
}

func main() {
	canvas := draws.New(app)
	canvas.Serve()
}
```