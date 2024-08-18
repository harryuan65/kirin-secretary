package pages

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/harryuan65/kirin_secretary/ctrl"
)

var downloadTypes = []string{
	"video",
	"music",
	"sound",
}

type YtDlpState struct {
	ffmpegInstalled binding.Bool
	url             binding.String
	outDir          binding.String
	downloadType    binding.String
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
			p.outDirBlock(),
			p.downloadTypeBlock(),
			p.downloadUrlBlock(),
			p.downloadBlock(),
		),
	)
}

func (p *YtDlpTab) versionBlock() *fyne.Container {
	versionString := binding.NewString()
	versionString.Set(ctrl.LoadingText)
	versionLabel := widget.NewLabelWithData(versionString)
	versionLabel.TextStyle.Bold = true

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

	return container.NewVBox(
		container.New(
			layout.NewGridLayout(4),
			container.New(layout.NewGridLayout(2), widget.NewLabel("yt-dlp version:"), versionLabel),
			widget.NewLabel(""),
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

func (p *YtDlpTab) outDirBlock() *fyne.Container {
	e := widget.NewEntryWithData(p.state.outDir)
	e.PlaceHolder = "Set download directory..."

	return container.New(
		layout.NewGridLayout(2),
		widget.NewLabel("Output Folder"),
		e,
	)
}

func (p *YtDlpTab) downloadTypeBlock() *fyne.Container {
	r := widget.NewRadioGroup(downloadTypes, func(s string) {
		p.state.downloadType.Set(s)
		log.Println("selected: ", s)
	})
	r.Selected, _ = p.state.downloadType.Get()

	return container.New(layout.NewGridLayout(2),
		widget.NewLabel("Download Type"),
		r,
	)
}

func (p *YtDlpTab) downloadUrlBlock() *fyne.Container {
	e := widget.NewEntryWithData(p.state.url)
	e.PlaceHolder = "https://www.youtube.com/watch?v=xxxx"

	return container.New(
		layout.NewGridLayout(2),
		widget.NewLabel("Video URL(https)"),
		e,
	)
}

func (p *YtDlpTab) downloadBlock() *fyne.Container {
	downloadStatusString := binding.NewString()
	downloadStatusLabel := widget.NewLabelWithData(downloadStatusString)
	downloadStatusLabel.Wrapping = fyne.TextWrapBreak

	scrollableStatusLabel := container.NewScroll(downloadStatusLabel)
	scrollableStatusLabel.SetMinSize(fyne.NewSize(scrollableStatusLabel.Size().Width, 150))

	handleDownload := func() {
		downloadStatusString.Set(ctrl.LoadingText)
		ffmpegInstalled, _ := p.state.ffmpegInstalled.Get()
		url, _ := p.state.url.Get()
		outDir, _ := p.state.outDir.Get()
		downloadType, _ := p.state.downloadType.Get()

		fmt.Println("p.state.ffmpegInstalled", ffmpegInstalled)
		fmt.Println("p.state.url", url)
		fmt.Println("p.state.outDir", outDir)
		fmt.Println("p.state.downloadType", downloadType)
		if !ctrl.IsDir(outDir) {
			downloadStatusString.Set(fmt.Sprintf("Error: %s is not a folder", outDir))
			return
		}

		if !strings.HasPrefix(url, "https://www.youtube.com") {
			downloadStatusString.Set(fmt.Sprintf("Error: %s is not a youtube url", url))
			return
		}

		downloadArgs := []string{}
		switch downloadType {
		case "music":
			if !ffmpegInstalled {
				downloadStatusString.Set("You need ffmpeg and ffprobe to downlaod as MP3")
				return
			}
			downloadArgs = append(downloadArgs, "-x", "--audio-format", "mp3")
		case "sound":
			if !ffmpegInstalled {
				downloadStatusString.Set("You need ffmpeg and ffprobe to downlaod as WAV")
				return
			}
			downloadArgs = append(downloadArgs, "-x", "--audio-format", "wav")
		default:
			downloadArgs = append(downloadArgs, "-f", "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best")
		}
		downloadArgs = append(downloadArgs, url)
		cmd := exec.Command("yt-dlp", downloadArgs...)
		log.Println("[os] cd into outDir...", outDir)
		os.Chdir(outDir)
		ctrl.Execute(&ctrl.Command{
			Label: "yt-dlp",
			Cmd:   cmd,
			OnSuccess: func() {
				scrollableStatusLabel.ScrollToBottom()
			},
			OnError: func(err error) {
				errString := fmt.Sprintf("Failed to download: %s \nError:%v", cmd.String(), err.Error())
				fmt.Println(errString)
				downloadStatusString.Set(errString)
			},
			OnOutput: func(line string) {
				acc, _ := downloadStatusString.Get()
				downloadStatusString.Set(acc + "\n" + line)
				scrollableStatusLabel.ScrollToBottom()
			},
		})
	}

	handleClearOutput := func() {
		downloadStatusString.Set("")
	}

	handleOpenOutDir := func() {
		outDir, _ := p.state.outDir.Get()
		log.Println("[yt-dlp] opening dir", outDir)
		ctrl.OpenInExplorer(outDir)
	}

	return container.NewVBox(
		container.New(
			layout.NewGridLayout(2),
			widget.NewLabel(""),
			widget.NewButton("Download", handleDownload),
		),
		container.New(
			layout.NewGridLayout(2),
			widget.NewLabel(""),
			widget.NewButton("Open Download Folder", handleOpenOutDir),
		),
		container.New(
			layout.NewGridLayout(2),
			widget.NewLabel(""),
			widget.NewButton("Clear Output", handleClearOutput),
		),
		widget.NewSeparator(),
		scrollableStatusLabel,
	)
}

func NewYtDlpTab() *container.TabItem {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// outDir: Default to home/downloads
	outDir := binding.NewString()
	outDir.Set(filepath.Join(homeDir, "Downloads"))

	// downloadType: Default to music
	downloadType := binding.NewString()
	downloadType.Set(downloadTypes[0])

	tab := &YtDlpTab{
		state: &YtDlpState{
			ffmpegInstalled: binding.NewBool(),
			url:             binding.NewString(),
			outDir:          outDir,
			downloadType:    downloadType,
		},
	}

	return tab.GetTab()
}
