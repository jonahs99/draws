package draws

import (
	"fmt"
)

// Context has methods for commands
type Context interface {
	BackgroundStyle(string)
	StrokeStyle(string)
	FillStyle(string)

	Clear()

	BeginPath()
	Stroke()

	Rect(x, y, w, h float64)
}

type context struct {
	command chan string
}

func (c *context) BackgroundStyle(style string) {
	c.command <- fmt.Sprintf("canvas.style.background='%s'", style)
}

func (c *context) StrokeStyle(style string) {
	c.command <- fmt.Sprintf("context.strokeStyle='%s'", style)
}

func (c *context) FillStyle(style string) {
	c.command <- fmt.Sprintf("context.fillStyle='%s'", style)
}

func (c *context) Clear() {
	c.command <- fmt.Sprintf("context.save();context.setTransform();context.clearRect(0,0,canvas.width,canvas.height);context.restore()")
}

func (c *context) BeginPath() {
	c.command <- "context.beginPath()"
}

func (c *context) Stroke() {
	c.command <- "context.stroke()"
}

func (c *context) Rect(x, y, w, h float64) {
	c.command <- fmt.Sprintf("context.rect(%v,%v,%v,%v)", x, y, w, h)
}
