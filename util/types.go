package util

type (
	UnitRef     = string
	GridCol     []UnitRef
	GridRow     []GridCol
	Group       []GridRow
	Constructor = UnitRef
	TUnitGrid   map[Constructor]Group
	Lab         = UnitRef
	TLabGrid    map[Lab]GridRow
	WeaponDef   struct {
		Range int
	}
	CustomParams struct {
		TechLevel int
		UnitGroup string
	}
	UnitProperties struct {
		MetalCost     int
		EnergyCost    int
		Buildtime     int
		BuildOptions  []UnitRef
		Health        int
		SightDistance int
		Speed         float64
		WeaponDefs    []WeaponDef
		CustomParams  *CustomParams
	}
	TranslationsT struct {
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
