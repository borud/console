// Package console is a simple Fyne widget for scrolling log messages. The intended use for this is
package console

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Console widget
type Console struct {
	widget.BaseWidget
	vbox               *fyne.Container
	list               *fyne.Widget
	scrollToBottomFunc func()
	// MaxLines is the maximum number of messages to keep.  Once we hit this
	// number of messages, older messages are removed.
	MaxLines int
	// ScrollToBottom will scroll the window to the latest message (bottom) when
	// a new message is appended. The default number of lines is 1000.
	ScrollToBottom bool
	// BackgroundColor is the default background color of messages.
	BackgroundColor color.Color
}

const (
	defaultMaxLines = 1000
)

var (
	defaultBackgroundColor = color.RGBA{R: 0xee, G: 0xee, B: 0xee, A: 0xff}
)

// NewConsole creates a new console widget with default MaxLines and the default
// message background color.
func NewConsole() *Console {
	console := &Console{
		MaxLines:        defaultMaxLines,
		vbox:            container.NewVBox(),
		BackgroundColor: defaultBackgroundColor,
	}
	console.ExtendBaseWidget(console)
	return console
}

// AppendWithColor appends a message and applies specified color to the message background.
func (c *Console) AppendWithColor(msg string, bgColor color.Color) {

	if len(c.vbox.Objects) > c.MaxLines {
		target := c.vbox.Objects[0]
		c.vbox.Remove(target)
	}

	text := widget.NewLabel(msg)
	// text := canvas.NewText(msg, color.Black)
	text.TextStyle.Monospace = true
	text.TextStyle.Bold = true
	text.Wrapping = fyne.TextWrapBreak

	rect := canvas.NewRectangle(bgColor)
	rect.CornerRadius = 5

	c.vbox.Add(container.NewStack(rect, container.NewPadded(text)))
	if c.ScrollToBottom {
		c.scrollToBottomFunc()
	}
}

// Append appends a message to the console.
func (c *Console) Append(msg string) {
	c.AppendWithColor(msg, c.BackgroundColor)
}

// Clear the console.
func (c *Console) Clear() {
	c.vbox.RemoveAll()
}

// CreateRenderer returns a new renderer.  Uses SimpleRenderer for now since we
// don't do anything fancy.
func (c *Console) CreateRenderer() fyne.WidgetRenderer {
	scroll := container.NewScroll(container.NewPadded(container.NewVBox(c.vbox)))
	c.scrollToBottomFunc = scroll.ScrollToBottom
	return widget.NewSimpleRenderer(scroll)
}
