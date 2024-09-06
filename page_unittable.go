package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/wezzle/bar-unit-info/util"
)

type UnitTablePage struct {
	// UI
	grid  *ui.Grid
	table *widgets.Table

	// State
	buildableUnits []util.UnitRef
	properties     map[util.UnitRef]*util.UnitProperties

	// Stats
	metalCost  *widgets.Gauge
	energyCost *widgets.Gauge
	buildtime  *widgets.Gauge
	health     *widgets.Gauge
	sightRange *widgets.Gauge
	speed      *widgets.Gauge
}

func (p *UnitTablePage) Render() {
	ui.Render(p.grid)
}

func (p *UnitTablePage) HandleEvents(e ui.Event) (Page, error) {
	switch e.ID {
	case "<C-c>":
		return nil, fmt.Errorf("exit")
	case "<Escape>":
		return constructorPage, nil
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		p.grid.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		p.Render()
	}

	p.Render()
	return nil, nil
}

func createUnitTablePage() (p *UnitTablePage) {
	p = &UnitTablePage{
		buildableUnits: make([]util.UnitRef, 0),
		properties:     make(map[util.UnitRef]*util.UnitProperties),
	}

	termWidth, termHeight := ui.TerminalDimensions()

	// Use labs to find buildable units
	// TODO we might miss units that are buildable by combat engineers and such
	for ref := range util.LabGrid {
		properties, err := util.LoadUnitProperties(ref)
		if err != nil {
			log.Printf("Skipping lab %s, no properties\n", ref)
			continue
		}
		for _, ref := range properties.BuildOptions {
			unitProperties, err := util.LoadUnitProperties(ref)
			if err != nil {
				log.Printf("Skipping unit %s, no properties\n", ref)
				continue

			}
			p.buildableUnits = append(p.buildableUnits, ref)
			p.properties[ref] = unitProperties
		}
	}

	p.table = widgets.NewTable()
	p.table.Rows = [][]string{
		{"Ref", "Name", "Metal cost", "Energy cost", "Buildtime", "Health", "Sight range", "Speed"},
	}
	for _, ref := range p.buildableUnits {
		up := p.properties[ref]
		d := time.Second * time.Duration(up.Buildtime/100)
		p.table.Rows = append(p.table.Rows, []string{
			ref,
			util.NameForRef(ref),
			strconv.Itoa(up.MetalCost),
			strconv.Itoa(up.EnergyCost),
			d.String(),
			strconv.Itoa(up.Health),
			strconv.Itoa(up.SightDistance),
			strconv.FormatFloat(up.Speed, 'f', -1, 64),
		})
	}
	p.table.TextStyle = ui.NewStyle(ui.ColorWhite)

	p.grid = ui.NewGrid()
	p.grid.SetRect(0, 0, termWidth, termHeight)
	p.grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0, p.table),
		),
	)
	return p
}
