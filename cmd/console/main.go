package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/borud/console"
)

func main() {
	app := app.New()
	win := app.NewWindow("Console test")
	win.Resize(fyne.NewSize(1000, 1000))

	console := console.NewConsole()
	console.ScrollToBottom = true

	okColor := color.RGBA{R: 0x0f, G: 0xff, B: 0x60, A: 0xaa}
	warningColor := color.RGBA{R: 0xff, G: 0xaa, B: 0x0f, A: 0xaa}
	errorColor := color.RGBA{R: 0xfc, G: 0x25, B: 0x12, A: 0xaa}

	go func() {
		for i := 0; ; i++ {
			console.AppendWithColor(fmt.Sprintf("%d this is a message that says everything is OK %s", i, time.Now().Format(time.RFC3339)), okColor)
			console.AppendWithColor(fmt.Sprintf("%d this is a warning message %s", i, time.Now().Format(time.RFC3339)), warningColor)
			console.AppendWithColor(fmt.Sprintf("%d this is an error message %s", i, time.Now().Format(time.RFC3339)), errorColor)
			console.Append(fmt.Sprintf("%d this is a plain message %s", i, time.Now().Format(time.RFC3339)))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	win.SetContent(console)
	win.ShowAndRun()
}
