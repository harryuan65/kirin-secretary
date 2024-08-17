package pages

import (
	"bufio"
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const LoadingText = "Loading..."

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
	// Create the command you want to execute
	cmd := exec.Command("bash", "-c", "for i in {1..5}; do echo $i; sleep 1; done")
	loadStatusString := binding.NewString()
	loadStatusString.Set(LoadingText)

	label := widget.NewLabelWithData(loadStatusString)
	// Get the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error obtaining stdout pipe: %v\n", err)
		return label
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return label
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("Output: %s\n", line)
			acc, err := loadStatusString.Get()
			if err != nil {
				loadStatusString.Set(err.Error())
				break
			}

			// Overwrite the loading text
			if acc == LoadingText {
				loadStatusString.Set(line)
			} else {
				loadStatusString.Set(acc + line)
			}
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stdout: %v\n", err)
			loadStatusString.Set(err.Error())
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
			loadStatusString.Set(err.Error())
		}
	}()

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
