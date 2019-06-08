package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var statusPanelInstance *statusPanel

type statusPanel struct {
	CommonPanel
	step *StepButton
	pb   *gtk.ProgressBar

	bed, tool0, tool1, tool2, tool3 *LabelWithImage
	file, left                      *LabelWithImage
	print, pause, stop              *gtk.Button
}

func StatusPanel(ui *UI, parent Panel) Panel {
	if statusPanelInstance == nil {
		m := &statusPanel{CommonPanel: NewCommonPanel(ui, parent)}
		m.panelH = 3
		m.b = NewBackgroundTask(time.Second*2, m.update)
		m.initialize()

		statusPanelInstance = m
	}

	return statusPanelInstance
}

func (m *statusPanel) initialize() {
	defer m.Initialize()

	m.Grid().Attach(m.createMainBox(), 1, 0, 4, 2)
	m.Grid().Attach(m.createPrintButton(), 1, 2, 1, 1)
	m.Grid().Attach(m.createPauseButton(), 2, 2, 1, 1)
	m.Grid().Attach(m.createStopButton(), 3, 2, 1, 1)
}

func (m *statusPanel) createProgressBar() *gtk.ProgressBar {
	m.pb = MustProgressBar()
	m.pb.SetShowText(true)
	m.pb.SetMarginTop(10)
	m.pb.SetMarginStart(10)
	m.pb.SetMarginEnd(10)

	return m.pb
}

func (m *statusPanel) createMainBox() *gtk.Box {
	grid := MustGrid()
	grid.SetHExpand(true)
	grid.Add(m.createInfoBox())
	grid.Add(m.createTemperatureBox())

	box := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetVExpand(true)
	box.Add(grid)
	box.Add(m.createProgressBar())

	return box
}

func (m *statusPanel) createInfoBox() *gtk.Box {
	m.file = MustLabelWithImage("file.svg", "")
	m.left = MustLabelWithImage("speed-step.svg", "")

	info := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	info.SetHAlign(gtk.ALIGN_START)
	info.SetHExpand(true)
	info.SetVExpand(true)
	info.Add(m.file)
	info.Add(m.left)
	info.SetMarginStart(10)

	return info
}

func (m *statusPanel) createTemperatureBox() *gtk.Box {
	m.bed = MustLabelWithImage("bed.svg", "")
	m.tool0 = MustLabelWithImage("extruder.svg", "")
	m.tool1 = MustLabelWithImage("extruder.svg", "")
	m.tool2 = MustLabelWithImage("extruder.svg", "")
	m.tool3 = MustLabelWithImage("extruder.svg", "")

	temp := MustBox(gtk.ORIENTATION_VERTICAL, 5)
	temp.SetHAlign(gtk.ALIGN_START)
	temp.SetHExpand(true)
	temp.SetVExpand(true)
	temp.Add(m.bed)
	temp.Add(m.tool0)
	temp.Add(m.tool1)
	temp.Add(m.tool2)
	temp.Add(m.tool3)

	return temp
}

func (m *statusPanel) createPrintButton() gtk.IWidget {
	m.print = MustButtonImage("Print", "status.svg", func() {
		defer m.updateTemperature()

		Logger.Warning("Starting a new job")
		if err := (&octoprint.StartRequest{}).Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})

	return m.print
}

func (m *statusPanel) createPauseButton() gtk.IWidget {
	m.pause = MustButtonImage("Pause", "pause.svg", func() {
		defer m.updateTemperature()

		Logger.Warning("Pausing/Resuming job")
		cmd := &octoprint.PauseRequest{Action: octoprint.Toggle}
		if err := cmd.Do(m.UI.Printer); err != nil {
			Logger.Error(err)
			return
		}
	})

	return m.pause
}

func (m *statusPanel) createStopButton() gtk.IWidget {
	m.stop = MustButtonImage("Stop", "stop.svg",
		ConfirmStopDialog(m.UI.w, "Are you sure you want to stop current print?", m),
	)
	return m.stop
}

func (m *statusPanel) update() {
	m.updateTemperature()
	m.updateJob()
}

func (m *statusPanel) updateTemperature() {
	s, err := (&octoprint.StateRequest{Exclude: []string{"sd"}}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	m.doUpdateState(&s.State)

	m.tool1.Hide()
	m.tool2.Hide()
	m.tool3.Hide()

	for tool, s := range s.Temperature.Current {
		text := fmt.Sprintf("%s: %.0f°C / %.0f°C", strings.Title(tool), s.Actual, s.Target)
		switch tool {
		case "bed":
			m.bed.Label.SetLabel(text)
		case "tool0":
			m.tool0.Label.SetLabel(text)
		case "tool1":
			m.tool1.Label.SetLabel(text)
			m.tool1.Show()
		case "tool2":
			m.tool2.Label.SetLabel(text)
			m.tool2.Show()
		case "tool3":
			m.tool3.Label.SetLabel(text)
			m.tool3.Show()
		}
	}
}

func (m *statusPanel) doUpdateState(s *octoprint.PrinterState) {
	switch {
	case s.Flags.Printing:
		m.print.SetSensitive(false)
		m.pause.SetSensitive(true)
		m.stop.SetSensitive(true)
	case s.Flags.Paused:
		m.print.SetSensitive(false)
		m.pause.SetLabel("Resume")
		m.pause.SetImage(MustImageFromFile("resume.svg"))
		m.pause.SetSensitive(true)
		m.stop.SetSensitive(true)
		return
	case s.Flags.Ready:
		m.print.SetSensitive(true)
		m.pause.SetSensitive(false)
		m.stop.SetSensitive(false)
	default:
		m.print.SetSensitive(false)
		m.pause.SetSensitive(false)
		m.stop.SetSensitive(false)
	}

	m.pause.SetLabel("Pause")
	m.pause.SetImage(MustImageFromFile("pause.svg"))
}

func (m *statusPanel) updateJob() {
	s, err := (&octoprint.JobRequest{}).Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}

	file := "<i>not-set</i>"
	if s.Job.File.Name != "" {
		file = filenameEllipsis(s.Job.File.Name)
	}

	m.file.Label.SetLabel(fmt.Sprintf("File: %s", file))
	m.pb.SetFraction(s.Progress.Completion / 100)

	if m.UI.State.IsOperational() {
		m.left.Label.SetLabel("Printer is ready")
		return
	}

	var text string
	switch s.Progress.Completion {
	case 100:
		text = fmt.Sprintf("Completed in %s", time.Duration(int64(s.Job.LastPrintTime)*1e9))
	case 0:
		text = "Warming up ..."
	default:
		Logger.Info(s.Progress.PrintTime)

		e := time.Duration(int64(s.Progress.PrintTime) * 1e9)
		l := time.Duration(int64(s.Progress.PrintTimeLeft) * 1e9)
		// eta := time.Now().Add(l).Format("3:04 PM")
		if l == 0 {
			text = fmt.Sprintf("Print Time: %s", e)
		} else {
			text = fmt.Sprintf("Print Time: %s | Left: %s", e, l)
		}
	}

	m.left.Label.SetLabel(text)
}

func filenameEllipsis(name string) string {
	l := len(name)
	if l > 32 {
		return name[:12] + "..." + name[l-17:l]
	}

	return name
}

func btou(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func ConfirmStopDialog(parent *gtk.Window, msg string, ma *statusPanel) func() {
	return func() {
		win := gtk.MessageDialogNewWithMarkup(
			parent,
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_YES_NO,
			"",
		)

		win.SetMarkup(CleanHTML(msg))
		defer win.Destroy()

		box, _ := win.GetContentArea()
		box.SetMarginStart(15)
		box.SetMarginEnd(15)
		box.SetMarginTop(15)
		box.SetMarginBottom(15)

		ctx, _ := win.GetStyleContext()
		ctx.AddClass("dialog")

		ergebnis := win.Run()

		if ergebnis == int(gtk.RESPONSE_YES) {

			Logger.Warning("Stopping job")
			if err := (&octoprint.CancelRequest{}).Do(ma.UI.Printer); err != nil {
				Logger.Error(err)
				return
			}
		}
	}
}
