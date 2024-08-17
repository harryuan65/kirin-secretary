package pages

import (
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/harryuan65/kirin_secretary/ctrl"
)

type YtDlpState struct {
	version         string
	ffmpegInstalled binding.Bool
}

type YtDlpTab struct {
	state *YtDlpState
}

func (p *YtDlpTab) GetTab() *container.TabItem {
	return container.NewTabItem(
		"yt-dlp",
		container.NewVBox(
			widget.NewLabel("yt-dlp"),
			p.versionBlock(),
			p.ffmpegStatusBlock(),
		),
	)
}

func (p *YtDlpTab) versionBlock() *fyne.Container {
	versionString := binding.NewString()
	versionString.Set(ctrl.LoadingText)
	versionLabel := widget.NewLabelWithData(versionString)

	versionCheck := func() {
		ctrl.Execute(&ctrl.Command{
			Label:     "checking yt-dlp version",
			Cmd:       exec.Command("yt-dlp", "--version"),
			OnSuccess: func() {},
			OnError: func(err error) {
				versionString.Set(err.Error())
			},
			OnOutput: func(s string) {
				versionString.Set(s)
			},
		})
	}
	go versionCheck()

	updateStatusString := binding.NewString()
	updateStatusString.Set("")
	updateStatusLabel := widget.NewLabelWithData(updateStatusString)
	updateStatusLabel.Wrapping = fyne.TextWrapBreak

	updateButton := widget.NewButton("Update", nil)
	updateButton.OnTapped = func() {
		updateStatusString.Set(ctrl.LoadingText)
		ctrl.Execute(&ctrl.Command{
			Label: "update yt-dlp",
			Cmd:   exec.Command("pip", "install", "--upgrade", "yt-dlp"),
			OnSuccess: func() {
				// Run check version again
				versionCheck()
			},
			OnError: func(err error) {},
			OnOutput: func(s string) {
				updateStatusString.Set(s)
			},
		})
	}

	return container.New(
		layout.NewGridLayoutWithRows(2),
		container.New(
			layout.NewGridLayout(3),
			container.New(layout.NewGridLayout(2), widget.NewLabel("yt-dlp version:"), versionLabel),
			updateButton,
		),
		updateStatusLabel,
	)
}

func (p *YtDlpTab) ffmpegStatusBlock() *fyne.Container {
	// Status string is a bind.string, listen to "installed" state changes and update it
	// Set status string to the label
	loadStatusString := binding.NewString()
	loadStatusString.Set(ctrl.LoadingText)
	label := widget.NewLabelWithData(loadStatusString)
	label.Wrapping = fyne.TextWrapBreak

	p.state.ffmpegInstalled.AddListener(binding.NewDataListener(func() {
		if installed, _ := p.state.ffmpegInstalled.Get(); installed {
			loadStatusString.Set("✅")
		} else {
			loadStatusString.Set("❌")
		}
	}))

	// Run command on component creation
	go func() {
		c := &ctrl.Command{
			Label: "checking ffmpeg installation",
			Cmd:   exec.Command("ffmpeg", "-version"),
			OnSuccess: func() {
				p.state.ffmpegInstalled.Set(true)
			},
			OnError: func(err error) {
				p.state.ffmpegInstalled.Set(false)
			},
		}

		ctrl.Execute(c)
	}()

	return container.New(layout.NewGridLayout(2), widget.NewLabel("ffmpeg:"), label)
}

func NewYtDlpTab() *container.TabItem {
	tab := &YtDlpTab{
		state: &YtDlpState{
			ffmpegInstalled: binding.NewBool(),
		},
	}

	return tab.GetTab()
}
