package draws

import (
	"fmt"
)

// Context has methods for commands
type Context interface {
	//Batch
	Batch()
	Draw()

	// Canvas
	BackgroundStyle(style string)
	Size(w, h float64)

	// Context
	LineWidth(w float64)
	LineCap(style string)
	LineJoin(style string)

	FillStyle(style string)
	StrokeStyle(style string)

	BeginPath()
	ClosePath()
	MoveTo(x, y float64)
	LineTo(x, y float64)
	Arc(x, y, r, start, end float64)
	ArcTo(x, y, r, start, end float64)
	Ellipse(x, y, rx, ry, start, end float64)
	Rect(x, y, w, h float64)

	Fill()
	Stroke()

	Rotate(angle float64)
	Scale(x, y float64)
	Translate(x, y float64)
	Transform(transform [6]float64)
	ResetTransform()

	GlobalAlpha(alpha float64)
	GlobalCompositeOperation(op string)

	Save()
	Restore()

	// Convieniece
	Clear()
	FillStroke()
	TranslateCenter()
	Circle(x, y, r float64)

	// Mouse
	WatchMouse()
	MousePos() (float64, float64)
}

type context struct {
	command        chan string
	event          chan string
	mousex, mousey float64
}

// Batch

func (c *context) Batch() {
	c.command <- "BATCH"
}

func (c *context) Draw() {
	c.command <- "DRAW"
}

// Context

func (c *context) BackgroundStyle(style string) {
	c.command <- fmt.Sprintf("canvas.style.background='%s'", style)
}

func (c *context) Size(w, h float64) {
	c.command <- fmt.Sprintf("canvas.width=%v;canvas.height=%v", w, h)
}

func (c *context) LineWidth(width float64) {
	c.command <- fmt.Sprintf("context.lineWidth=%v", width)
}

func (c *context) LineCap(style string) {
	c.command <- fmt.Sprintf("context.lineCap='%s'", style)
}

func (c *context) LineJoin(style string) {
	c.command <- fmt.Sprintf("context.linJoin='%s'", style)
}

func (c *context) FillStyle(style string) {
	c.command <- fmt.Sprintf("context.fillStyle='%s'", style)
}

func (c *context) StrokeStyle(style string) {
	c.command <- fmt.Sprintf("context.strokeStyle='%s'", style)
}

func (c *context) BeginPath() {
	c.command <- "context.beginPath()"
}

func (c *context) ClosePath() {
	c.command <- "context.closePath()"
}

func (c *context) MoveTo(x, y float64) {
	c.command <- fmt.Sprintf("context.moveTo(%v,%v)", x, y)
}

func (c *context) LineTo(x, y float64) {
	c.command <- fmt.Sprintf("context.lineTo(%v,%v)", x, y)
}

func (c *context) Arc(x, y, r, start, end float64) {
	c.command <- fmt.Sprintf("context.arc(%v,%v,%v,%v,%v)", x, y, r, start, end)
}

func (c *context) ArcTo(x, y, r, start, end float64) {
	c.command <- fmt.Sprintf("context.arcTo(%v,%v,%v,%v,%v)", x, y, r, start, end)
}

func (c *context) Ellipse(x, y, rx, ry, start, end float64) {
	c.command <- fmt.Sprintf("context.ellipse(%v,%v,%v,%v,%v,%v)", x, y, rx, ry, start, end)
}

func (c *context) Rect(x, y, w, h float64) {
	c.command <- fmt.Sprintf("context.rect(%v,%v,%v,%v)", x, y, w, h)
}

func (c *context) Fill() {
	c.command <- "context.fill()"
}

func (c *context) Stroke() {
	c.command <- "context.stroke()"
}

func (c *context) Rotate(angle float64) {
	c.command <- fmt.Sprintf("context.rotate(%v)", angle)
}

func (c *context) Scale(x, y float64) {
	c.command <- fmt.Sprintf("context.scale(%v,%v)", x, y)
}

func (c *context) Translate(x, y float64) {
	c.command <- fmt.Sprintf("context.translate(%v,%v)", x, y)
}

func (c *context) Transform(matrix [6]float64) {
	c.command <- fmt.Sprintf("context.transform(%v,%v,%v,%v,%v,%v)", []interface{}{matrix}...)
}

func (c *context) ResetTransform() {
	c.command <- "context.resetTransform()"
}

func (c *context) GlobalAlpha(alpha float64) {
	c.command <- fmt.Sprintf("context.globalAlpha=%v", alpha)
}

func (c *context) GlobalCompositeOperation(op string) {
	c.command <- fmt.Sprintf("context.globalCompositeOperation=%s", op)
}

func (c *context) Save() {
	c.command <- "context.save()"
}

func (c *context) Restore() {
	c.command <- "context.restore()"
}

// Convienience

func (c *context) Clear() {
	c.command <- "context.save();context.resetTransform();context.clearRect(0,0,canvas.width,canvas.height);context.restore()"
}

func (c *context) TranslateCenter() {
	c.command <- "context.translate(canvas.width/2,canvas.height/2)"
}

func (c *context) FillStroke() {
	c.command <- "context.fill();context.stroke()"
}

func (c *context) Circle(x, y, r float64) {
	c.command <- fmt.Sprintf("context.arc(%v,%v,%v,0,2*Math.PI)", x, y, r)
}

// Mouse

func (c *context) WatchMouse() {
	c.command <- "sendMouseEvents()"
}

func (c *context) MousePos() (float64, float64) {
	return c.mousex, c.mousey
}
