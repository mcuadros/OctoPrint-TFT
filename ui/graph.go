package ui

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func init() {
	ColorAlpha := drawing.Color{R: 255, G: 255, B: 255, A: 0}
	chart.DefaultTextColor = chart.ColorWhite
	chart.DefaultBackgroundColor = ColorAlpha
	chart.DefaultBackgroundStrokeColor = chart.ColorBlack
	chart.DefaultCanvasColor = ColorAlpha
	chart.DefaultCanvasStrokeColor = chart.ColorBlack
	chart.DefaultTextColor = chart.ColorWhite
	chart.DefaultAxisColor = chart.ColorWhite
	chart.DefaultStrokeColor = chart.ColorLightGray
	chart.DefaultFillColor = chart.ColorBlue
	chart.DefaultAnnotationFillColor = chart.ColorBlack
	chart.DefaultGridLineColor = chart.ColorLightGray

}

var colors = []drawing.Color{
	chart.ColorBlue,
	chart.ColorRed,
}

type GraphPanel struct {
	CommonPanel
	Label *gtk.Label
}

func NewGraphPanel(ui *UI) *GraphPanel {
	m := &GraphPanel{CommonPanel: NewCommonPanel(ui, nil)}
	m.initialize()
	return m
}

func (m *GraphPanel) initialize() {
	logo := m.drawHistory(m.retrieveHistory())
	//	m.Label = MustLabel("Connecting to OctoPrint...")

	box := MustBox(gtk.ORIENTATION_VERTICAL, 15)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetVExpand(true)
	box.SetHExpand(true)

	box.Add(logo)
	//	box.Add(m.Label)

	m.Grid().Attach(box, 1, 0, 1, 1)
}
func (m *GraphPanel) retrieveHistory() *octoprint.FullStateResponse {
	r := octoprint.StateRequest{History: true}
	s, err := r.Do(m.UI.Printer)
	if err != nil {
		Logger.Error(err)
	}

	return s
}

func (m *GraphPanel) drawHistory(s *octoprint.FullStateResponse) gtk.IWidget {
	w, h := m.UI.w.GetSize()
	graph := chart.Chart{Width: w - 20, Height: h - 20}
	graph.XAxis = chart.XAxis{
		Style:          chart.StyleShow(),
		TickStyle:      chart.Style{FontSize: 8},
		ValueFormatter: timeSinceFormatter,
	}

	graph.YAxis = chart.YAxis{
		Style:          chart.StyleShow(),
		TickStyle:      chart.Style{FontSize: 8},
		Range:          &chart.ContinuousRange{Min: 0.0, Max: 300.0},
		ValueFormatter: temperatureFormatter,
	}

	graph.Elements = []chart.Renderable{
		changeLegendColor,
		chart.LegendLeft(&graph, chart.Style{
			FillColor: chart.DefaultBackgroundColor,
		}),
	}

	graph.Series = m.buildSeries(s)

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Fatal(err)
	}

	fo, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	if _, err := fo.Write(buffer.Bytes()); err != nil {
		panic(err)
	}

	fo.Close()

	p, err := gdk.PixbufNewFromFile("output.png")
	if err != nil {
		Logger.Error(err)
	}

	i, err := gtk.ImageNewFromPixbuf(p)
	if err != nil {
		panic(err)
	}

	return i
}

func (m *GraphPanel) buildSeries(s *octoprint.FullStateResponse) []chart.Series {
	actual, target := map[string]*chart.TimeSeries{}, map[string]*chart.TimeSeries{}

	var i int
	for tool := range s.Temperature.Current {
		actual[tool] = &chart.TimeSeries{
			Name: fmt.Sprintf("Actual %s", tool),
			Style: chart.Style{
				Show:        true,
				StrokeColor: colors[i],
			},
		}

		target[tool] = &chart.TimeSeries{
			Name: fmt.Sprintf("Target %s", tool),
			Style: chart.Style{
				Show:            true,
				StrokeColor:     colors[i],
				StrokeDashArray: []float64{5.0, 5.0},
			},
		}

		i++
	}

	for _, history := range s.Temperature.History {
		for tool, data := range history.Tools {
			actual[tool].XValues = append(actual[tool].XValues, history.Time.Time)
			actual[tool].YValues = append(actual[tool].YValues, data.Actual)

			target[tool].XValues = append(target[tool].XValues, history.Time.Time)
			target[tool].YValues = append(target[tool].YValues, data.Target)
		}
	}

	var series []chart.Series
	for tool := range actual {
		ann := chart.LastValueAnnotation(*actual[tool], temperatureFormatter)
		ann.Style.StrokeColor = actual[tool].Style.StrokeColor

		series = append(series, *actual[tool], *target[tool], ann)
	}

	return series
}

func timeSinceFormatter(v interface{}) string {
	if typed, isTyped := v.(float64); isTyped {
		t := time.Unix(0, int64(typed))

		return humanize.CustomRelTime(t, time.Now(), "ago", "from now", magnitudes)
	}

	return ""
}

func temperatureFormatter(v interface{}) string {
	return fmt.Sprintf("%d°C", int(v.(float64)))
}

func changeLegendColor(r chart.Renderer, cb chart.Box, chartDefaults chart.Style) {
	r.SetFillColor(chart.DefaultBackgroundColor)
}

var magnitudes = []humanize.RelTimeMagnitude{
	{time.Minute, "now", time.Second},
	{2 * time.Minute, "-1 min %s", 1},
	{time.Hour, "-%d mins", time.Minute},
}
