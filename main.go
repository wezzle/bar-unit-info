package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lukegb/dds"
	lua "github.com/yuin/gopher-lua"
)

func isScavengers(L *lua.LState) int {
	L.Push(lua.LFalse)
	return 1
}

func getModOptions(L *lua.LState) int {
	t := lua.LTable{}
	t.RawSetString("forceallunits", lua.LFalse)
	L.Push(&t)
	return 1
}

type (
	UnitRef        = string
	GridCol        []UnitRef
	GridRow        []GridCol
	Group          []GridRow
	Constructor    = UnitRef
	UnitGrid       map[Constructor]Group
	UnitProperties struct {
		metalcost    int
		energycost   int
		buildtime    int
		buildOptions []UnitRef
	}
	Translations struct {
		Units struct {
			Factions                  map[string]string  `json:"factions"`
			Dead                      string             `json:"dead"`
			Heap                      string             `json:"heap"`
			DecoyCommanderNameTag     string             `json:"decoyCommanderNameTag"`
			Scavenger                 string             `json:"scavenger"`
			ScavCommanderNameTag      string             `json:"scavCommanderNameTag"`
			ScavDecoyCommanderNameTag string             `json:"scavDecoyCommanderNameTag"`
			Names                     map[UnitRef]string `json:"names"`
			Descriptions              map[UnitRef]string `json:"descriptions"`
		} `json:"units"`
	}
)

var (
	unitPropertyCache = make(map[UnitRef]UnitProperties)
	debug             *widgets.List
	translations      Translations
	unitGrid          UnitGrid
	constructorPage   *ConstructorPage
)

func indexFromLValue(v lua.LValue) int {
	index, err := strconv.Atoi(v.String())
	if err != nil {
		panic(err)
	}
	return index - 1
}

func loadImage(ref UnitRef) image.Image {
	r, err := os.Open(fmt.Sprintf("./bar-repo/unitpics/%s.dds", ref))
	if err != nil {
		return nil
	}
	img, err := dds.Decode(r)
	if err != nil {
		return nil
	}
	return img
}

func loadTranslations(lang string) (translations Translations) {
	f, err := os.Open(fmt.Sprintf("./bar-repo/language/%s/units.json", lang))
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&translations)
	if err != nil {
		panic(err)
	}
	return
}

func loadUnitGrid() UnitGrid {
	L := lua.NewState()

	springTable := lua.LTable{}
	utilitiesTable := lua.LTable{}
	gametypeTable := lua.LTable{}
	gametypeTable.RawSetString("IsScavengers", L.NewFunction(isScavengers))
	utilitiesTable.RawSetString("Gametype", &gametypeTable)
	springTable.RawSetString("Utilities", &utilitiesTable)
	springTable.RawSetString("GetModOptions", L.NewFunction(getModOptions))
	L.SetGlobal("Spring", &springTable)
	defer L.Close()
	if err := L.DoFile("./bar-repo/luaui/configs/gridmenu_layouts.lua"); err != nil {
		panic(err)
	}

	lv := L.Get(-1)

	unitGrid := make(UnitGrid)

	v := lv.(*lua.LTable).RawGetString("UnitGrids")
	v.(*lua.LTable).ForEach(func(k lua.LValue, v lua.LValue) {
		constructor := Constructor(k.String())
		unitGrid[constructor] = make(Group, 4)

		v.(*lua.LTable).ForEach(func(k lua.LValue, group lua.LValue) {
			groupIndex := indexFromLValue(k)
			unitGrid[constructor][groupIndex] = make(GridRow, 3)
			group.(*lua.LTable).ForEach(func(k lua.LValue, units lua.LValue) {
				rowIndex := indexFromLValue(k)
				unitGrid[constructor][groupIndex][rowIndex] = make(GridCol, 4)
				units.(*lua.LTable).ForEach(func(k lua.LValue, unitName lua.LValue) {
					colIndex := indexFromLValue(k)
					unitGrid[constructor][groupIndex][rowIndex][colIndex] = UnitRef(unitName.String())
				})
			})
		})
	})

	return unitGrid
}

func findUnitPropertiesFile(ref UnitRef) (string, error) {
	r, err := regexp.Compile(fmt.Sprintf("%s.lua$", ref))
	if err != nil {
		return "", err
	}

	file := ""
	err = filepath.WalkDir("./bar-repo/units", func(path string, d os.DirEntry, err error) error {
		if err == nil && r.MatchString(path) {
			file = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

func loadUnitProperties(ref UnitRef) (*UnitProperties, error) {
	if properties, ok := unitPropertyCache[ref]; ok {
		return &properties, nil
	}

	L := lua.NewState()
	defer L.Close()

	unitFilepath, err := findUnitPropertiesFile(ref)
	if err != nil {
		return nil, err
	}
	if err := L.DoFile(unitFilepath); err != nil {
		return nil, err
	}

	lv := L.Get(-1)
	data := lv.(*lua.LTable).RawGetString(ref).(*lua.LTable)

	metalcost, _ := strconv.Atoi(data.RawGetString("metalcost").String())
	energycost, _ := strconv.Atoi(data.RawGetString("energycost").String())
	buildtime, _ := strconv.Atoi(data.RawGetString("buildtime").String())

	bo := data.RawGetString("buildoptions")
	var buildOptions []UnitRef
	if bo.Type() == lua.LTTable {
		buildOptions = make([]UnitRef, 0)
		bo.(*lua.LTable).ForEach(func(index lua.LValue, v lua.LValue) {
			buildOptions = append(buildOptions, v.String())
		})
	}

	properties := UnitProperties{
		metalcost:    metalcost,
		energycost:   energycost,
		buildtime:    buildtime,
		buildOptions: buildOptions,
	}
	unitPropertyCache[ref] = properties
	return &properties, nil
}

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
	translations = loadTranslations("en")
	unitGrid = loadUnitGrid()

	// Start terminal UI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create initial page
	var activePage Page
	constructorPage = createConstructorPage()
	activePage = constructorPage

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
