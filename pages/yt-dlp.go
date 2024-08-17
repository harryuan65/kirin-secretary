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
			p.loadLabel(),
		),
	)
}

func (p *YtDlpTab) ffmpegStatusLabel() *widget.Label {
	// Binding String
	// Label with the binding string
	// State Binding Bool, changes change the string

	ffmpegStatus := binding.NewString()
	ffmpegStatus.Set("❌")
	ffmpegStatusLabel := widget.NewLabelWithData(ffmpegStatus)
	p.state.ffmpegInstalled.AddListener(binding.NewDataListener(func() {
		if installed, _ := p.state.ffmpegInstalled.Get(); installed {
			ffmpegStatus.Set("✅")
		} else {
			ffmpegStatus.Set("❌")
		}
	}))

	return ffmpegStatusLabel
}

func (p *YtDlpTab) updateButton() *widget.Button {
	return widget.NewButton("Update", func() {
		installed, _ := p.state.ffmpegInstalled.Get()
		p.state.ffmpegInstalled.Set(!installed)
	})
}

func (p *YtDlpTab) loadLabel() *widget.Label {
	loadStatusString := binding.NewString()
	loadStatusString.Set(ctrl.LoadingText)

	label := widget.NewLabelWithData(loadStatusString)
	// ctrl.BindExecOutput(exec.Command("bash", "-c", "for i in {1..5}; do echo $i; sleep 1; done"), loadStatusString)
	ctrl.BindExecOutput(exec.Command("bash", "-c", "for i in {1..5}; do echo $i; sleep 1; done"), loadStatusString)

	return label
}

func NewYtDlpTab() *container.TabItem {
	tab := &YtDlpTab{
		state: &YtDlpState{
			ffmpegInstalled: binding.NewBool(),
		},
	}

	return tab.GetTab()
}
