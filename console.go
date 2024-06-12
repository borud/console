// Package console is a simple Fyne widget for scrolling log messages.
package console

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Console widget
type Console struct {
	widget.BaseWidget
	list     *widget.List
	mu       sync.Mutex
	messages []messageItem
	// MaxLines is the maximum number of messages to keep.  Once we hit this
	// number of messages, older messages are removed.
	MaxLines int
	// ScrollToBottom will scroll the window to the latest message (bottom) when
	// a new message is appended. The default number of lines is 1000.
	ScrollToBottom bool
	// BackgroundColor is the default background color of messages.
	BackgroundColor color.Color
}

type messageItem struct {
	msg   string
	color color.Color
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
		messages:        []messageItem{},
		MaxLines:        defaultMaxLines,
		BackgroundColor: defaultBackgroundColor,
	}

	// The way we create and update list items feels a bit naughty since
	// create() and update() have to be in sync with regard to the inner
	// structure of the CanvasObject.
	var list *widget.List
	list = widget.NewList(
		func() int {
			return len(console.messages)
		},

		func() fyne.CanvasObject {
			label := widget.NewLabel("<template>")
			label.TextStyle.Monospace = true
			label.TextStyle.Bold = true
			label.Wrapping = fyne.TextWrapBreak

			background := canvas.NewRectangle(defaultBackgroundColor)
			background.CornerRadius = 5

			return container.NewStack(background, container.NewPadded(label))
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {
			stack, ok := o.(*fyne.Container)
			if !ok {
				panic("expected a *fyne.Container (stack)")
			}

			background, ok := stack.Objects[0].(*canvas.Rectangle)
			if !ok {
				panic("expected a *canvas.Rectangle (background)")
			}

			padded, ok := stack.Objects[1].(*fyne.Container)
			if !ok {
				panic("expected a *fyne.Container (padding)")
			}

			label, ok := padded.Objects[0].(*widget.Label)
			if !ok {
				panic("expected a *widget.Label")
			}

			// limit lock scope to where we interact with the messages slice.
			console.mu.Lock()
			background.FillColor = console.messages[i].color
			label.Text = console.messages[i].msg
			console.mu.Unlock()

			list.SetItemHeight(i, stack.MinSize().Height)
			o.Refresh()
		})

	console.list = list

	console.ExtendBaseWidget(console)
	return console
}

// AppendWithColor appends a message and applies specified color to the message background.
func (c *Console) AppendWithColor(msg string, bgColor color.Color) {
	c.mu.Lock()

	if len(c.messages) >= c.MaxLines {
		c.messages = c.messages[1:]
	}

	c.messages = append(c.messages, messageItem{
		msg:   msg,
		color: bgColor,
	})

	// Note that we have to release the lock before we make a call that will
	// update the list state.
	c.mu.Unlock()

	if c.ScrollToBottom {
		c.list.ScrollToBottom()
	}
}

// Append appends a message to the console.
func (c *Console) Append(msg string) {
	c.AppendWithColor(msg, c.BackgroundColor)
}

// Clear the console.
func (c *Console) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.messages = c.messages[:0]
}

// CreateRenderer returns a new renderer.  Uses SimpleRenderer for now since we
// don't do anything fancy.
func (c *Console) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewPadded(c.list))
}
