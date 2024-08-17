package pages

import (
	"os/exec"

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
			widget.NewLabel("YT-DLP"),
			container.New(layout.NewGridLayout(3), widget.NewLabel("yt-dlp version:"), p.VersionLabel(), p.updateButton()),
			container.New(layout.NewGridLayout(2), widget.NewLabel("ffmpeg:"), p.ffmpegStatusLabel()),
		),
	)
}

func (p *YtDlpTab) VersionLabel() *widget.Label {
	versionString := binding.NewString()
	versionString.Set(ctrl.LoadingText)
	versionLabel := widget.NewLabelWithData(versionString)

	go func() {
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
	}()

	return versionLabel
}

func (p *YtDlpTab) ffmpegStatusLabel() *widget.Label {
	// Status string is a bind.string, listen to "installed" state changes and update it
	// Set status string to the label
	loadStatusString := binding.NewString()
	loadStatusString.Set(ctrl.LoadingText)
	label := widget.NewLabelWithData(loadStatusString)

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

	return label
}

func (p *YtDlpTab) updateButton() *widget.Button {
	return widget.NewButton("Update", func() {
		installed, _ := p.state.ffmpegInstalled.Get()
		p.state.ffmpegInstalled.Set(!installed)
	})
}

func NewYtDlpTab() *container.TabItem {
	tab := &YtDlpTab{
		state: &YtDlpState{
			ffmpegInstalled: binding.NewBool(),
		},
	}

	return tab.GetTab()
}
