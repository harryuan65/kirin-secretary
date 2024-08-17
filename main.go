package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/harryuan65/kirin_secretary/pages"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabContainer Widget")

	var tabs []*container.TabItem = []*container.TabItem{
		pages.NewYtDlpTab(),
	}
	appTabs := container.NewAppTabs(
		tabs...,
	)

	appTabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(appTabs)
	myWindow.ShowAndRun()
}
