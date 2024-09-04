package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ConstructorPage struct {
	list    *widgets.List
	image   *widgets.Image
	content *widgets.Paragraph
	grid    *ui.Grid

	constructorList []Constructor

	previousKey         string
	previousSelectedRow int
}

func (p *ConstructorPage) Render() {
	ui.Render(p.grid)
}

func (p *ConstructorPage) HandleEvents(e ui.Event) (Page, error) {
	switch e.ID {
	case "q", "<C-c>":
		return nil, fmt.Errorf("exit")
	case "j", "<Down>":
		p.list.ScrollDown()
	case "k", "<Up>":
		p.list.ScrollUp()
	case "<C-d>":
		p.list.ScrollHalfPageDown()
	case "<C-u>":
		p.list.ScrollHalfPageUp()
	case "<C-f>":
		p.list.ScrollPageDown()
	case "<C-b>":
		p.list.ScrollPageUp()
	case "g":
		if p.previousKey == "g" {
			p.list.ScrollTop()
		}
	case "<Home>":
		p.list.ScrollTop()
	case "G", "<End>":
		p.list.ScrollBottom()
	case "<Enter>":
		selectedConstructor := UnitRef(p.constructorList[p.list.SelectedRow])
		return createBuildGridPage(selectedConstructor), nil
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		p.grid.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		ui.Render(p.grid)
	}

	if p.previousSelectedRow != p.list.SelectedRow {
		selectedConstructor := UnitRef(p.constructorList[p.list.SelectedRow])
		img := loadImage(selectedConstructor)
		p.image.Image = img
		p.content.Text = fmt.Sprintf("%s\n\n%s", translations.Units.Names[selectedConstructor], translations.Units.Descriptions[selectedConstructor])
		p.previousSelectedRow = p.list.SelectedRow
	}

	if p.previousKey == "g" {
		p.previousKey = ""
	} else {
		p.previousKey = e.ID
	}

	p.Render()
	return nil, nil
}

func createConstructorPage() (page *ConstructorPage) {
	debug = widgets.NewList()
	debug.Rows = make([]string, 0)
	debug.Title = "Debug"
	debug.WrapText = true

	page = &ConstructorPage{}
	page.constructorList = make([]Constructor, 0)
	for constructor := range unitGrid {
		if strings.Contains(constructor, "lvl") {
			continue
		}
		properties, err := loadUnitProperties(constructor)
		if err != nil || properties.buildOptions == nil {
			log.Printf("Skipping %s, no properties\n", constructor)
			continue
		}
		page.constructorList = append(page.constructorList, constructor)
	}

	sort.Strings(page.constructorList)

	page.list = widgets.NewList()
	page.list.Title = "Constructors"
	page.list.Rows = make([]string, 0)
	for index, constructor := range page.constructorList {
		page.list.Rows = append(page.list.Rows, fmt.Sprintf("[%02d] %s", index, constructor))
	}
	page.list.TextStyle = ui.NewStyle(ui.ColorYellow)
	page.list.WrapText = false

	firstConstructor := UnitRef(page.constructorList[0])

	img := loadImage(firstConstructor)
	page.image = widgets.NewImage(img)
	page.image.Title = "Preview"

	page.grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	page.grid.SetRect(0, 0, termWidth, termHeight)

	page.content = widgets.NewParagraph()
	page.content.Text = fmt.Sprintf("%s\n\n%s", translations.Units.Names[firstConstructor], translations.Units.Descriptions[firstConstructor])

	page.grid.Set(
		ui.NewRow(0.8,
			ui.NewCol(1.0/2, page.list),
			ui.NewCol(1.0/2,
				ui.NewRow(1.0/2, page.image),
				ui.NewRow(1.0/2, page.content),
			),
		),
		ui.NewRow(0.2, debug),
	)

	page.Render()
	return
}
