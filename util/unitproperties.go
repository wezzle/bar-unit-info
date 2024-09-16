package util

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
	// TODO check beamlaser calculations, the following units are off: corkorg corjugg
	dps := 0.0
	for _, wd := range p.WeaponDefs {
		if wd.Damage.Default == 0 {
			continue
		}
		damage := wd.Damage.Default / wd.ReloadTime
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
