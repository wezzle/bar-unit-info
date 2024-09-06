package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/wezzle/bar-unit-info/util"
)

var (
	unitPropertyCache = make(map[util.UnitRef]util.UnitProperties)
	debug             *widgets.List
	constructorPage   *ConstructorPage
)

func debugLine(s string) {
	debug.Rows = append(debug.Rows, s)
	debug.ScrollBottom()
	ui.Render(debug)
}

type Page interface {
	Render()
	HandleEvents(e ui.Event) (Page, error)
}

func main() {
	// Load globals
	util.LoadTranslations("en")
	util.LoadGridLayouts()

	// Start terminal UI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create initial page
	var activePage Page
	constructorPage = createConstructorPage()
	activePage = constructorPage
	// activePage = createUnitPage("corkarg", nil)
	// activePage = createUnitTablePage()
	// activePage.Render()

	// Handle event loop
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		page, err := activePage.HandleEvents(e)
		if err != nil {
			return
		}
		if page != nil {
			activePage = page
			ui.Clear()
			activePage.Render()
		}
	}
}
