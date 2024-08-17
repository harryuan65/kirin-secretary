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
			widget.NewLabel("Version"),
			container.New(layout.NewGridLayout(3), widget.NewLabel("ffmpeg status:"), p.ffmpegStatusLabel(), p.updateButton()),
		),
	)
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
			Label: "check ffmpeg installation",
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

// func (p *YtDlpTab) ffmpegStatusLabel() *widget.Label {
// 	// Binding String
// 	// Label with the binding string
// 	// State Binding Bool, changes change the string

// 	ffmpegStatus := binding.NewString()
// 	ffmpegStatus.Set(ctrl.LoadingText)
// 	ffmpegStatusLabel := widget.NewLabelWithData(ffmpegStatus)
// 	p.state.ffmpegInstalled.AddListener(binding.NewDataListener(func() {
// 		if installed, _ := p.state.ffmpegInstalled.Get(); installed {
// 			ffmpegStatus.Set("✅")
// 		} else {
// 			ffmpegStatus.Set("❌")
// 		}
// 	}))

// 	return ffmpegStatusLabel
// }

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
