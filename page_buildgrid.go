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
	ui.Render(p.tabPane)
	ui.Render(p.grid)
}

func (p *BuildGridPage) HandleEvents(e ui.Event) (Page, error) {
	switch e.ID {
	case "q", "<C-c>":
		return nil, fmt.Errorf("exit")
	case "z":
		p.tabPane.ActiveTabIndex = 0
	case "x":
		p.tabPane.ActiveTabIndex = 1
	case "c":
		p.tabPane.ActiveTabIndex = 2
	case "v":
		p.tabPane.ActiveTabIndex = 3
	case "<Escape>":
		return constructorPage, nil
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		p.grid.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		ui.Render(p.grid)
	}

	if p.tabPane.ActiveTabIndex != p.lastActiveTab {
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

func createBuildGridPage(ref UnitRef) (page *BuildGridPage) {
	page = &BuildGridPage{
		ref: ref,
	}

	var err error
	page.properties, err = loadUnitProperties(ref)
	if err != nil {
		panic(err)
	}

	page.group = unitGrid[ref]

	termWidth, termHeight := ui.TerminalDimensions()

	page.tabPane = widgets.NewTabPane("(Z) Economy", "(X) Combat", "(C) Utility", "(V) Build")
	page.tabPane.SetRect(0, 1, termWidth, 4)
	page.tabPane.Border = true

	page.tabs = make([]*ui.Grid, 4)
	for i := range 4 {
		tab := ui.NewGrid()
		tab.Set(
			ui.NewRow(1.0/3,
				ui.NewCol(1.0/4, createBuildItem(page.group[i][0][0])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][0][1])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][0][2])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][0][3])),
			),
			ui.NewRow(1.0/3,
				ui.NewCol(1.0/4, createBuildItem(page.group[i][1][0])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][1][1])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][1][2])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][1][3])),
			),
			ui.NewRow(1.0/3,
				ui.NewCol(1.0/4, createBuildItem(page.group[i][2][0])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][2][1])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][2][2])),
				ui.NewCol(1.0/4, createBuildItem(page.group[i][2][3])),
			),
		)
		page.tabs[i] = tab
	}

	page.grid = ui.NewGrid()
	page.grid.SetRect(0, 5, termWidth, termHeight-5)
	page.grid.Set(
		ui.NewRow(1.0,
			page.tabs[page.tabPane.ActiveTabIndex],
		),
		// ui.NewRow(0.2, debug),
	)
	page.lastActiveTab = page.tabPane.ActiveTabIndex

	page.Render()
	return
}
