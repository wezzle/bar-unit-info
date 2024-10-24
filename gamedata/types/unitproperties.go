package types

import (
	"fmt"
	"strings"
)

func (p *UnitProperties) SummarizeWeaponTypes() string {
	counts := make(map[string]int)
	for _, wd := range p.WeaponDefs {
		c, ok := counts[wd.WeaponType]
		if !ok {
			c = 0
		}
		c = c + 1
		counts[wd.WeaponType] = c
	}

	summary := []string{}
	for wt, c := range counts {
		desc := wt
		if c > 1 {
			desc = fmt.Sprintf("%dx %s", c, desc)
		}
		summary = append(summary, desc)
	}
	return strings.Join(summary, ", ")
}

func (p *UnitProperties) MaxWeaponRange() float64 {
	weaponRange := 0.0
	for _, wd := range p.WeaponDefs {
		weaponRange = max(weaponRange, wd.Range)
	}
	return weaponRange
}

func (p *UnitProperties) DPS() float64 {
	// Burst = shots per burst
	// burstRate = delay between shots in a burst (seconds)
	// projectiles = projectiles in shot (see sprayAngle)
	// sprayAngle = How inaccurate are individual projectiles in a burst?
	//
	// BeamLasers
	// minIntensity = BeamLaser only. The minimum percentage the weapon's damage can fall-off to over its range. Setting to 1.0 will disable fall off entirely.
	// dynDamageInverted = If true the damage curve is inverted i.e. the weapon does more damage at greater ranges as opposed to less.
	// dynDamageExp = Exponent of the range-dependent damage formula, the default of 0.0 disables dynamic damage, 1.0 means linear scaling, 2.0 quadratic and so on.
	// dynDamageMin = The minimum floor value that range-dependent damage can drop to.
	// dynDamageRange = If set to non-zero values the weapon will use this value in the range-dependant damage formula instead of the actual range.
	// beamtime = The laser maintains it beam for this many seconds, spreading its damage over that time.
	// beamburst = Lets a laser use burst mechanics, but sets `beamtime` to the duration of 1 sim frame.
	//
	// LaserCannon
	//
	//
	// TODO check beamlaser calculation for legbastion, leginc seems to work fine
	dps := 0.0

	for _, weapon := range p.Weapons {
		wd, exists := p.WeaponDefs[strings.ToLower(weapon.Def)]
		if !exists {
			continue
		}

		var damage float64
		if d, exists := wd.Damage[strings.ToLower(weapon.OnlyTargetCategory)]; exists {
			damage = d
		} else if d, exists := wd.Damage["default"]; exists {
			damage = d
		} else {
			continue
		}

		if damage == 0 {
			continue
		}

		damage = damage / wd.ReloadTime
		if wd.Burst != 0 {
			damage = damage * float64(wd.Burst)
		}
		if wd.Projectiles != 0 {
			damage = damage * float64(wd.Projectiles)
		}
		dps = dps + damage
	}

	return dps
}

func (p *UnitProperties) IsBuilding() bool {
	return p.Speed == 0
}
