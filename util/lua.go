package util

import (
	"math"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

var (
	UnitGrid          TUnitGrid
	LabGrid           TLabGrid
	unitPropertyCache = make(map[UnitRef]UnitProperties)
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

func indexFromLValue(v lua.LValue) int {
	index, err := strconv.Atoi(v.String())
	if err != nil {
		panic(err)
	}
	return index - 1
}

func loadUnitGrid(v *lua.LTable) TUnitGrid {
	grid := make(TUnitGrid)

	v.ForEach(func(k lua.LValue, v lua.LValue) {
		constructor := Constructor(k.String())
		grid[constructor] = make(Group, 4)
		v.(*lua.LTable).ForEach(func(k lua.LValue, group lua.LValue) {
			groupIndex := indexFromLValue(k)
			if groupIndex >= 4 {
				return
			}
			grid[constructor][groupIndex] = make(GridRow, 3)
			group.(*lua.LTable).ForEach(func(k lua.LValue, units lua.LValue) {
				rowIndex := indexFromLValue(k)
				if rowIndex >= 3 {
					return
				}
				grid[constructor][groupIndex][rowIndex] = make(GridCol, 4)
				units.(*lua.LTable).ForEach(func(k lua.LValue, unitName lua.LValue) {
					colIndex := indexFromLValue(k)
					if colIndex >= 4 {
						return
					}
					grid[constructor][groupIndex][rowIndex][colIndex] = UnitRef(unitName.String())
				})
			})
		})
	})

	return grid
}

func loadLabGrid(v *lua.LTable) TLabGrid {
	grid := make(TLabGrid)

	v.ForEach(func(k lua.LValue, v lua.LValue) {
		lab := Constructor(k.String())
		grid[lab] = make(GridRow, 3)
		for i := range grid[lab] {
			grid[lab][i] = make(GridCol, 4)
		}

		rowIndex := 0
		v.(*lua.LTable).ForEach(func(k lua.LValue, unitName lua.LValue) {
			colIndex := indexFromLValue(k)
			if colIndex%4 == 0 && colIndex != 0 {
				rowIndex = rowIndex + 1
			}
			grid[lab][rowIndex][colIndex%4] = UnitRef(unitName.String())
		})
	})

	return grid
}

func LoadGridLayouts() {
	fileContents, err := repoFiles.ReadFile("bar-repo/luaui/configs/gridmenu_layouts.lua")
	if err != nil {
		panic(err)
	}

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
	if err := L.DoString(string(fileContents)); err != nil {
		panic(err)
	}

	lv := L.Get(-1).(*lua.LTable)

	UnitGrid = loadUnitGrid(lv.RawGetString("UnitGrids").(*lua.LTable))
	LabGrid = loadLabGrid(lv.RawGetString("LabGrids").(*lua.LTable))
}

func LoadUnitProperties(ref UnitRef) (*UnitProperties, error) {
	if properties, ok := unitPropertyCache[ref]; ok {
		return &properties, nil
	}

	unitFilepath, err := findUnitPropertiesFile(ref)
	if err != nil {
		return nil, err
	}
	fileContents, err := repoFiles.ReadFile(unitFilepath)
	if err != nil {
		panic(err)
	}

	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(string(fileContents)); err != nil {
		return nil, err
	}

	lv := L.Get(-1)
	data := lv.(*lua.LTable).RawGetString(ref).(*lua.LTable)

	// Simple stats
	metalcost, _ := strconv.Atoi(data.RawGetString("metalcost").String())
	if metalcost == 0 {
		metalcost, _ = strconv.Atoi(data.RawGetString("buildcostmetal").String())
	}
	energycost, _ := strconv.Atoi(data.RawGetString("energycost").String())
	if energycost == 0 {
		energycost, _ = strconv.Atoi(data.RawGetString("buildcostenergy").String())
	}
	buildtime, _ := strconv.Atoi(data.RawGetString("buildtime").String())
	health, _ := strconv.Atoi(data.RawGetString("health").String())
	sightdistance, err := strconv.Atoi(data.RawGetString("sightdistance").String())
	if err != nil {
		floatSightDistance, _ := strconv.ParseFloat(data.RawGetString("sightdistance").String(), 64)
		sightdistance = int(math.Round(floatSightDistance))
	}
	speed, _ := strconv.ParseFloat(data.RawGetString("speed").String(), 64)

	var workertime *int
	if data.RawGetString("workertime").Type() != lua.LTNil {
		workertimeVal, _ := strconv.Atoi(data.RawGetString("workertime").String())
		workertime = &workertimeVal
	}

	// Build option slice
	bo := data.RawGetString("buildoptions")
	var buildOptions []UnitRef
	if bo.Type() == lua.LTTable {
		buildOptions = make([]UnitRef, 0)
		bo.(*lua.LTable).ForEach(func(index lua.LValue, v lua.LValue) {
			buildOptions = append(buildOptions, v.String())
		})
	}

	// Custom params
	cp := data.RawGetString("customparams")
	customParams := CustomParams{}
	if cp.Type() == lua.LTTable {
		customParams.TechLevel, _ = strconv.Atoi(cp.(*lua.LTable).RawGetString("techlevel").String())
		customParams.UnitGroup = cp.(*lua.LTable).RawGetString("unitgroup").String()
	}
	// Find lab that produces this one and get techlevel from that unit
	if customParams.TechLevel == 0 && !strings.Contains(customParams.UnitGroup, "builder") {
		found := false
		for labRef := range LabGrid {
			lp, _ := LoadUnitProperties(labRef)
			for _, bo := range lp.BuildOptions {
				if bo == ref {
					found = true
					customParams.TechLevel = lp.CustomParams.TechLevel
					break
				}
			}
			if found {
				break
			}
		}
	}
	// Find a unit that produces this one and get techlevel from that unit
	if customParams.TechLevel == 0 {
		found := false
		for ref, up := range unitPropertyCache {
			for _, bo := range up.BuildOptions {
				if bo == ref {
					found = true
					customParams.TechLevel = up.CustomParams.TechLevel
					break
				}
			}
			if found {
				break
			}
		}
	}
	// Default tech level to 1 if not found
	if customParams.TechLevel == 0 {
		customParams.TechLevel = 1
	}

	properties := UnitProperties{
		MetalCost:     metalcost,
		EnergyCost:    energycost,
		Buildtime:     buildtime,
		BuildOptions:  buildOptions,
		Health:        health,
		SightDistance: sightdistance,
		Speed:         speed,
		Buildpower:    workertime,
		CustomParams:  &customParams,
	}
	unitPropertyCache[ref] = properties
	return &properties, nil
}
