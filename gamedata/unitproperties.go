package gamedata

import "github.com/wezzle/bar-unit-info/gamedata/types"

func GetUnitPropertiesByRef() types.UnitPropertiesByRef {
	return unitPropertiesData
}

func GetUnitProperties(ref string) (types.UnitProperties, bool) {
	up, ok := unitPropertiesData[ref]
    return up, ok
}

