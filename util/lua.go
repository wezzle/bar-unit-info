package util

import (
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

	properties := UnitProperties{}

	p := LuaTableParser{data}

	// Simple stats assignments

	properties.MetalCost = IgnoreError("metalcost", p.Int)
	if properties.MetalCost == 0 {
		properties.MetalCost = IgnoreError("buildcostmetal", p.Int)
	}

	properties.EnergyCost = IgnoreError("energycost", p.Int)
	if properties.EnergyCost == 0 {
		properties.EnergyCost = IgnoreError("buildcostenergy", p.Int)
	}

	properties.Buildtime = IgnoreError("buildtime", p.Int)
	properties.Health = IgnoreError("health", p.Int)
	properties.SightDistance = int(IgnoreError("sightdistance", p.Float64))
	properties.Speed = IgnoreError("speed", p.Float64)
	properties.Buildpower = IgnoreError("workertime", p.OptionalInt)
	properties.RadarDistance = IgnoreError("radardistance", p.OptionalInt)
	properties.JammerDistance = IgnoreError("radardistancejam", p.OptionalInt)
	properties.SonarDistance = IgnoreError("sonardistance", p.OptionalInt)

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

	properties.BuildOptions = buildOptions
	properties.CustomParams = &customParams
	properties.WeaponDefs = ParseWeaponDefs(data)

	unitPropertyCache[ref] = properties
	return &properties, nil
}

func ParseWeaponDefs(data *lua.LTable) []WeaponDef {
	defs := make([]WeaponDef, 0)
	wd := data.RawGetString("weapondefs")
	if wd.Type() != lua.LTTable {
		return nil
	}

	wd.(*lua.LTable).ForEach(func(k lua.LValue, v lua.LValue) {
		vT := v.(*lua.LTable)
		p := LuaTableParser{vT}
		damage := IgnoreError("damage", p.Table)
		def := WeaponDef{
			Name:             IgnoreError("name", p.String),
			WeaponType:       IgnoreError("weapontype", p.String),
			Id:               IgnoreError("id", p.Int),
			CustomParams:     map[string]string{},
			AvoidFriendly:    IgnoreError("avoidfriendly", p.Bool),
			AvoidFeature:     IgnoreError("avoidfeature", p.Bool),
			AvoidNeutral:     IgnoreError("avoidneutral", p.Bool),
			AvoidGround:      IgnoreError("avoidground", p.Bool),
			AvoidCloaked:     IgnoreError("avoidcloaked", p.Bool),
			CollideEnemy:     IgnoreError("collideenemy", p.Bool),
			CollideFriendly:  IgnoreError("collidefriendly", p.Bool),
			CollideFeature:   IgnoreError("collidefeature", p.Bool),
			CollideNeutral:   IgnoreError("collideneutral", p.Bool),
			CollideFireBase:  IgnoreError("collidefirebase", p.Bool),
			CollideNonTarget: IgnoreError("collidenontarget", p.Bool),
			CollideGround:    IgnoreError("collideground", p.Bool),
			CollideCloaked:   IgnoreError("collidecloaked", p.Bool),
			Damage: Damage{
				Default: IgnoreError("default", damage.Float64),
			},
			ExplosionSpeed:           IgnoreError("explosionspeed", p.Float64),
			ImpactOnly:               IgnoreError("impactonly", p.Bool),
			NoSelfDamage:             IgnoreError("noselfdamage", p.Bool),
			NoExplode:                IgnoreError("noexplode", p.Bool),
			Burnblow:                 IgnoreError("burnblow", p.Bool),
			DamageAreaOfEffect:       IgnoreError("damageareaofeffect", p.Float64),
			EdgeEffectiveness:        IgnoreError("edgeeffectiveness", p.Float64),
			CollisionSize:            IgnoreError("collisionsize", p.Float64),
			WeaponVelocity:           IgnoreError("weaponvelocity", p.Float64),
			StartVelocity:            IgnoreError("startvelocity", p.Float64),
			Weaponacceleration:       IgnoreError("weaponacceleration", p.Float64),
			ReloadTime:               IgnoreError("reloadtime", p.Float64),
			BurstRate:                IgnoreError("burstrate", p.Float64),
			Burst:                    IgnoreError("burst", p.Int),
			Projectiles:              IgnoreError("projectiles", p.Int),
			WaterBounce:              IgnoreError("waterbounce", p.Bool),
			GroundBounce:             IgnoreError("groundbounce", p.Bool),
			BounceSlip:               IgnoreError("bounceslip", p.Float64),
			BounceRebound:            IgnoreError("bouncerebound", p.Float64),
			NumBounce:                IgnoreError("numbounce", p.Int),
			ImpulseFactor:            IgnoreError("impulsefactor", p.Float64),
			ImpulseBoost:             IgnoreError("impulseboost", p.Float64),
			CraterMult:               IgnoreError("cratermult", p.Float64),
			CraterBoost:              IgnoreError("craterboost", p.Float64),
			CraterAreaOfEffect:       IgnoreError("craterareaofeffect", p.Float64),
			Waterweapon:              IgnoreError("waterweapon", p.Bool),
			Submissile:               IgnoreError("submissile", p.Bool),
			FireSubmersed:            IgnoreError("firesubmersed", p.Bool),
			Commandfire:              IgnoreError("commandfire", p.Bool),
			Range:                    IgnoreError("range", p.Float64),
			Heightmod:                IgnoreError("heightmod", p.Float64),
			TargetBorder:             IgnoreError("targetborder", p.Float64),
			CylinderTargeting:        IgnoreError("cylindertargeting", p.Float64),
			Turret:                   IgnoreError("turret", p.Bool),
			FixedLauncher:            IgnoreError("fixedlauncher", p.Bool),
			Tolerance:                IgnoreError("tolerance", p.Float64),
			Firetolerance:            IgnoreError("firetolerance", p.Float64),
			HighTrajectory:           IgnoreError("hightrajectory", p.Int),
			TrajectoryHeight:         IgnoreError("trajectoryheight", p.Float64),
			Tracks:                   IgnoreError("tracks", p.Bool),
			Wobble:                   IgnoreError("wobble", p.Float64),
			Dance:                    IgnoreError("dance", p.Float64),
			GravityAffected:          IgnoreError("gravityaffected", p.Bool),
			MyGravity:                IgnoreError("mygravity", p.Float64),
			CanAttackGround:          IgnoreError("canattackground", p.Bool),
			WeaponTimer:              IgnoreError("weapontimer", p.Float64),
			Flighttime:               IgnoreError("flighttime", p.Float64),
			Turnrate:                 IgnoreError("turnrate", p.Float64),
			HeightBoostFactor:        IgnoreError("heightboostfactor", p.Float64),
			ProximityPriority:        IgnoreError("proximitypriority", p.Float64),
			AllowNonBlockingAim:      IgnoreError("allownonblockingaim", p.Bool),
			Accuracy:                 IgnoreError("accuracy", p.Float64),
			SprayAngle:               IgnoreError("sprayangle", p.Float64),
			MovingAccuracy:           IgnoreError("movingaccuracy", p.Float64),
			TargetMoveError:          IgnoreError("targetmoveerror", p.Float64),
			LeadLimit:                IgnoreError("leadlimit", p.Float64),
			LeadBonus:                IgnoreError("leadbonus", p.Float64),
			PredictBoost:             IgnoreError("predictboost", p.Float64),
			OwnerExpAccWeight:        IgnoreError("ownerexpaccweight", p.Float64),
			MinIntensity:             IgnoreError("minintensity", p.Float64),
			Duration:                 IgnoreError("duration", p.Float64),
			Beamtime:                 IgnoreError("beamtime", p.Float64),
			Beamburst:                IgnoreError("beamburst", p.Bool),
			BeamTTL:                  IgnoreError("beamttl", p.Int),
			SweepFire:                IgnoreError("sweepfire", p.Bool),
			LargeBeamLaser:           IgnoreError("largebeamlaser", p.Bool),
			SizeGrowth:               IgnoreError("sizegrowth", p.Float64),
			FlameGfxTime:             IgnoreError("flamegfxtime", p.Float64),
			MetalPerShot:             IgnoreError("metalpershot", p.Float64),
			EnergyPerShot:            IgnoreError("energypershot", p.Float64),
			FireStarter:              IgnoreError("firestarter", p.Float64),
			Paralyzer:                IgnoreError("paralyzer", p.Bool),
			ParalyzeTime:             IgnoreError("paralyzetime", p.Int),
			Stockpile:                IgnoreError("stockpile", p.Bool),
			StockpileTime:            IgnoreError("stockpiletime", p.Float64),
			Targetable:               IgnoreError("targetable", p.Int),
			Interceptor:              IgnoreError("interceptor", p.Int),
			InterceptedByShieldType:  IgnoreError("interceptedbyshieldtype", p.Int64),
			Coverage:                 IgnoreError("coverage", p.Float64),
			InterceptSolo:            IgnoreError("interceptsolo", p.Bool),
			DynDamageInverted:        IgnoreError("dyndamageinverted", p.Bool),
			DynDamageExp:             IgnoreError("dyndamageexp", p.Float64),
			DynDamageMin:             IgnoreError("dyndamagemin", p.Float64),
			DynDamageRange:           IgnoreError("dyndamagerange", p.Float64),
			Shield:                   Shield{},
			RechargeDelay:            IgnoreError("rechargedelay", p.Float64),
			Model:                    IgnoreError("model", p.String),
			Size:                     IgnoreError("size", p.Float64),
			ScarGlowColorMap:         IgnoreError("scarglowcolormap", p.String),
			ScarIndices:              ScarIndices{},
			ExplosionScar:            IgnoreError("explosionscar", p.Bool),
			ScarDiameter:             IgnoreError("scardiameter", p.Float64),
			ScarAlpha:                IgnoreError("scaralpha", p.Float64),
			ScarGlow:                 IgnoreError("scarglow", p.Float64),
			ScarTtl:                  IgnoreError("scarttl", p.Float64),
			ScarGlowTtl:              IgnoreError("scarglowttl", p.Float64),
			ScarDotElimination:       IgnoreError("scardotelimination", p.Float64),
			ScarProjVector:           [4]float64{},
			ScarColorTint:            [4]float64{},
			AlwaysVisible:            IgnoreError("alwaysvisible", p.Bool),
			CameraShake:              IgnoreError("camerashake", p.Float64),
			SmokeTrail:               IgnoreError("smoketrail", p.Bool),
			SmokeTrailCastShadow:     IgnoreError("smoketrailcastshadow", p.Bool),
			SmokePeriod:              IgnoreError("smokeperiod", p.Int),
			SmokeTime:                IgnoreError("smoketime", p.Int),
			SmokeSize:                IgnoreError("smokesize", p.Float64),
			SmokeColor:               IgnoreError("smokecolor", p.Float64),
			CastShadow:               IgnoreError("castshadow", p.Bool),
			SizeDecay:                IgnoreError("sizedecay", p.Float64),
			AlphaDecay:               IgnoreError("alphadecay", p.Float64),
			Separation:               IgnoreError("separation", p.Float64),
			NoGap:                    IgnoreError("nogap", p.Bool),
			Stages:                   IgnoreError("stages", p.Int),
			LodDistance:              IgnoreError("loddistance", p.Int),
			Thickness:                IgnoreError("thickness", p.Float64),
			CoreThickness:            IgnoreError("corethickness", p.Float64),
			LaserFlareSize:           IgnoreError("laserflaresize", p.Float64),
			TileLength:               IgnoreError("tilelength", p.Float64),
			ScrollSpeed:              IgnoreError("scrollspeed", p.Float64),
			PulseSpeed:               IgnoreError("pulsespeed", p.Float64),
			BeamDecay:                IgnoreError("beamdecay", p.Float64),
			FalloffRate:              IgnoreError("falloffrate", p.Float64),
			Hardstop:                 IgnoreError("hardstop", p.Bool),
			RgbColor:                 [3]float64{},
			RgbColor2:                [3]float64{},
			Intensity:                IgnoreError("intensity", p.Float64),
			Colormap:                 IgnoreError("colormap", p.String),
			CegTag:                   IgnoreError("cegtag", p.String),
			ExplosionGenerator:       IgnoreError("explosiongenerator", p.String),
			BounceExplosionGenerator: IgnoreError("bounceexplosiongenerator", p.String),
			SoundTrigger:             IgnoreError("soundtrigger", p.Bool),
			SoundStart:               IgnoreError("soundstart", p.String),
			SoundHitDry:              IgnoreError("soundhitdry", p.String),
			SoundHitWet:              IgnoreError("soundhitwet", p.String),
			SoundStartVolume:         IgnoreError("soundstartvolume", p.Float64),
			SoundHitDryVolume:        IgnoreError("soundhitdryvolume", p.Float64),
			SoundHitWetVolume:        IgnoreError("soundhitwetvolume", p.Float64),
		}
		defs = append(defs, def)
	})
	return defs
}
