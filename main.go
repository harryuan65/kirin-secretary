package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/harryuan65/kirin_secretary/pages"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("TabContainer Widget")
	w.Resize(fyne.NewSize(400, 300))
	w.SetTitle("Kirin Secretary")

	var tabs []*container.TabItem = []*container.TabItem{
		pages.NewYtDlpTab(),
	}
	appTabs := container.NewAppTabs(
		tabs...,
	)

	appTabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(appTabs)
	w.ShowAndRun()
}
