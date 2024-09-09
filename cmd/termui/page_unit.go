package main

import (
	"fmt"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/wezzle/bar-unit-info/util"
)

type UnitPage struct {
	// UI
	grid  *ui.Grid
	hero  *widgets.Paragraph
	image *widgets.Image

	// Stats
	metalCost  *widgets.Gauge
	energyCost *widgets.Gauge
	buildtime  *widgets.Gauge
	health     *widgets.Gauge
	sightRange *widgets.Gauge
	speed      *widgets.Gauge

	ref        util.UnitRef
	properties *util.UnitProperties

	parentPage *BuildGridPage
}

func (p *UnitPage) Render() {
	ui.Render(p.grid)
}

func (p *UnitPage) HandleEvents(e ui.Event) (Page, error) {
	switch e.ID {
	case "<C-c>":
		return nil, fmt.Errorf("exit")
	case "<Escape>":
		return p.parentPage, nil
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		p.grid.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		p.Render()
	}

	p.Render()
	return nil, nil
}

func createUnitPage(ref util.UnitRef, parent *BuildGridPage) (p *UnitPage) {
	p = &UnitPage{
		parentPage: parent,
		ref:        ref,
	}

	var err error
	p.properties, err = util.LoadUnitProperties(ref)
	if err != nil {
		panic(err)
	}

	p.hero = widgets.NewParagraph()
	p.hero.Text = fmt.Sprintf("%s\n\n%s", util.NameForRef(ref), util.DescriptionForRef(ref))

	img := util.LoadImage(ref)
	p.image = widgets.NewImage(img)
	p.image.Title = "Preview"

	// Stats
	p.metalCost = widgets.NewGauge()
	p.metalCost.Title = "Metal cost"
	p.metalCost.Percent = min(int(float64(p.properties.MetalCost)/250.0), 100)
	p.metalCost.BarColor = ui.ColorWhite
	p.metalCost.BorderStyle.Fg = ui.ColorWhite
	p.metalCost.Label = strconv.Itoa(p.properties.MetalCost)

	p.energyCost = widgets.NewGauge()
	p.energyCost.Title = "Energy cost"
	p.energyCost.Percent = min(int(float64(p.properties.EnergyCost)/900), 100)
	p.energyCost.BarColor = ui.ColorYellow
	p.energyCost.BorderStyle.Fg = ui.ColorWhite
	p.energyCost.Label = strconv.Itoa(p.properties.EnergyCost)

	d := time.Second * time.Duration(p.properties.Buildtime/100)
	p.buildtime = widgets.NewGauge()
	p.buildtime.Title = "Buildtime"
	p.buildtime.Percent = min(int(float64(p.properties.Buildtime)/1000), 100)
	p.buildtime.BarColor = ui.ColorYellow
	p.buildtime.BorderStyle.Fg = ui.ColorWhite
	p.buildtime.Label = d.String()

	p.health = widgets.NewGauge()
	p.health.Title = "Health"
	p.health.Percent = min(int(float64(p.properties.Health)/150), 100)
	p.health.BarColor = ui.ColorGreen
	p.health.BorderStyle.Fg = ui.ColorWhite
	p.health.Label = strconv.Itoa(p.properties.Health)

	p.sightRange = widgets.NewGauge()
	p.sightRange.Title = "Sight range"
	p.sightRange.Percent = min(int(float64(p.properties.SightDistance)/35), 100)
	p.sightRange.BarColor = ui.ColorWhite
	p.sightRange.BorderStyle.Fg = ui.ColorWhite
	p.sightRange.Label = strconv.Itoa(p.properties.SightDistance)

	p.speed = widgets.NewGauge()
	p.speed.Title = "Speed"
	p.speed.Percent = min(int(float64(p.properties.Speed)/1.5), 100)
	p.speed.BarColor = ui.ColorBlue
	p.speed.BorderStyle.Fg = ui.ColorWhite
	p.speed.Label = strconv.Itoa(int(p.properties.Speed))

	gaugeCount := 6.0

	termWidth, termHeight := ui.TerminalDimensions()

	p.grid = ui.NewGrid()
	p.grid.SetRect(0, 0, termWidth, termHeight)
	p.grid.Set(
		ui.NewRow(0.2,
			ui.NewCol(0.8, p.hero),
			ui.NewCol(0.2, p.image),
		),
		ui.NewRow(0.6,
			ui.NewCol(1.0,
				ui.NewRow(1.0/gaugeCount, p.metalCost),
				ui.NewRow(1.0/gaugeCount, p.energyCost),
				ui.NewRow(1.0/gaugeCount, p.buildtime),
				ui.NewRow(1.0/gaugeCount, p.health),
				ui.NewRow(1.0/gaugeCount, p.sightRange),
				ui.NewRow(1.0/gaugeCount, p.speed),
			),
		),
		// ui.NewRow(0.2, debug),
	)
	return p
}
