package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type BuildGridPage struct {
	grid    *ui.Grid
	tabPane *widgets.TabPane
	tabs    []*ui.Grid

	ref           UnitRef
	properties    *UnitProperties
	group         Group
	lastActiveTab int
}

func (p *BuildGridPage) Render() {
	if p.tabPane != nil {
		ui.Render(p.tabPane)
	}
	ui.Render(p.grid)
}

func (p *BuildGridPage) HandleEvents(e ui.Event) (Page, error) {
	var selectedUnit [2]int
	hasSelectedUnit := false
	switch e.ID {
	case "<C-c>":
		return nil, fmt.Errorf("exit")
	case "z":
		if p.tabPane != nil {
			p.tabPane.ActiveTabIndex = 0
		} else {
			selectedUnit[0] = 2
			selectedUnit[1] = 0
			hasSelectedUnit = true
		}
	case "x":
		if p.tabPane != nil {
			p.tabPane.ActiveTabIndex = 1
		} else {
			selectedUnit[0] = 2
			selectedUnit[1] = 1
			hasSelectedUnit = true
		}
	case "c":
		if p.tabPane != nil {
			p.tabPane.ActiveTabIndex = 2
		} else {
			selectedUnit[0] = 2
			selectedUnit[1] = 2
			hasSelectedUnit = true
		}
	case "v":
		if p.tabPane != nil {
			p.tabPane.ActiveTabIndex = 3
		} else {
			selectedUnit[0] = 2
			selectedUnit[1] = 3
			hasSelectedUnit = true
		}
	case "q":
		selectedUnit[0] = 0
		selectedUnit[1] = 0
		hasSelectedUnit = true
	case "w":
		selectedUnit[0] = 0
		selectedUnit[1] = 1
		hasSelectedUnit = true
	case "e":
		selectedUnit[0] = 0
		selectedUnit[1] = 2
		hasSelectedUnit = true
	case "r":
		selectedUnit[0] = 0
		selectedUnit[1] = 3
		hasSelectedUnit = true
	case "a":
		selectedUnit[0] = 1
		selectedUnit[1] = 0
		hasSelectedUnit = true
	case "s":
		selectedUnit[0] = 1
		selectedUnit[1] = 1
		hasSelectedUnit = true
	case "d":
		selectedUnit[0] = 1
		selectedUnit[1] = 2
		hasSelectedUnit = true
	case "f":
		selectedUnit[0] = 1
		selectedUnit[1] = 3
		hasSelectedUnit = true
	case "<Escape>":
		return constructorPage, nil
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		p.grid.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		p.Render()
	}

	if hasSelectedUnit {
		var activeIndex int
		if p.tabPane != nil {
			activeIndex = p.tabPane.ActiveTabIndex
		} else {
			activeIndex = 0
		}
		units := p.group[activeIndex]
		unit := units[selectedUnit[0]][selectedUnit[1]]
		return createUnitPage(unit, p), nil
	}

	if p.tabPane != nil && p.tabPane.ActiveTabIndex != p.lastActiveTab {
		p.grid.Set(
			ui.NewRow(1.0,
				p.tabs[p.tabPane.ActiveTabIndex],
			),
		)
		p.lastActiveTab = p.tabPane.ActiveTabIndex
	}

	p.Render()
	return nil, nil
}

func createPlaceholder(title string) *widgets.Paragraph {
	placeholder := widgets.NewParagraph()
	placeholder.Title = title
	return placeholder
}

func createBuildItem(ref UnitRef) interface{} {
	if ref == "" {
		return widgets.NewParagraph()
	}
	img := loadImage(ref)
	i := widgets.NewImage(img)
	i.Title = translations.Units.Names[ref]
	return i
}

func createPane(rows GridRow) *ui.Grid {
	pane := ui.NewGrid()
	pane.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/4, createBuildItem(rows[0][0])),
			ui.NewCol(1.0/4, createBuildItem(rows[0][1])),
			ui.NewCol(1.0/4, createBuildItem(rows[0][2])),
			ui.NewCol(1.0/4, createBuildItem(rows[0][3])),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/4, createBuildItem(rows[1][0])),
			ui.NewCol(1.0/4, createBuildItem(rows[1][1])),
			ui.NewCol(1.0/4, createBuildItem(rows[1][2])),
			ui.NewCol(1.0/4, createBuildItem(rows[1][3])),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/4, createBuildItem(rows[2][0])),
			ui.NewCol(1.0/4, createBuildItem(rows[2][1])),
			ui.NewCol(1.0/4, createBuildItem(rows[2][2])),
			ui.NewCol(1.0/4, createBuildItem(rows[2][3])),
		),
	)
	return pane
}

func createBuildGridPage(ref UnitRef) (p *BuildGridPage) {
	p = &BuildGridPage{
		ref: ref,
	}

	var err error
	p.properties, err = loadUnitProperties(ref)
	if err != nil {
		panic(err)
	}

	termWidth, termHeight := ui.TerminalDimensions()

	var ok bool
	if p.group, ok = unitGrid[ref]; ok {
		p.tabPane = widgets.NewTabPane("(Z) Economy", "(X) Combat", "(C) Utility", "(V) Build")
		p.tabPane.SetRect(0, 1, termWidth, 4)
		p.tabPane.Border = true

		p.tabs = make([]*ui.Grid, 4)
		for i := range 4 {
			tab := createPane(p.group[i])
			p.tabs[i] = tab
		}
	} else {
		p.group = make(Group, 1)
		p.group[0] = labGrid[ref]

		p.tabs = make([]*ui.Grid, 1)
		p.tabs[0] = createPane(p.group[0])
	}

	var activeIndex int
	if p.tabPane != nil {
		activeIndex = p.tabPane.ActiveTabIndex
	} else {
		activeIndex = 0
	}

	p.grid = ui.NewGrid()
	p.grid.SetRect(0, 5, termWidth, termHeight-5)
	p.grid.Set(
		ui.NewRow(1.0,
			p.tabs[activeIndex],
		),
		// ui.NewRow(0.2, debug),
	)

	if p.tabPane != nil {
		p.lastActiveTab = p.tabPane.ActiveTabIndex
	}

	p.Render()
	return
}
