package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Test")

	content := widget.NewLabel("Content")
	w.SetContent(container.NewVBox(widget.NewLabel("Test"), content, widget.NewButton("Go", func() {
		content.SetText("Clicked")
	})))

	w.ShowAndRun()
}
