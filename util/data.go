package util

import (
	"fmt"
	"image"
	"io/fs"
	"os"
	"regexp"

	"github.com/lukegb/dds"
	"github.com/wezzle/bar-unit-info/gamedata"
	"github.com/wezzle/bar-unit-info/gamedata/types"
)

func NameForRef(ref types.UnitRef) string {
	return gamedata.GetTranslations().Units.Names[ref]
}

func DescriptionForRef(ref types.UnitRef) string {
	return gamedata.GetTranslations().Units.Descriptions[ref]
}

func FactionForRef(ref types.UnitRef) string {
	shortcode := ref[0:3]
	if shortcode == "lee" {
		shortcode = "leg"
	}
	return gamedata.GetTranslations().Units.Factions[shortcode]
}

func OtherFactions(faction string, includeRandom bool) []string {
	factions := []string{}
	for _, f := range gamedata.GetTranslations().Units.Factions {
		if f == "Random" && !includeRandom {
			continue
		}
		if f != faction {
			factions = append(factions, f)
		}
	}
	return factions
}

func PrefixForFaction(faction string) string {
	for prefix, data := range gamedata.GetTranslations().Units.Factions {
		if data == faction {
			return prefix
		}
	}
	return "random"
}

func LoadImage(ref types.UnitRef) image.Image {
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

func findUnitPropertiesFile(ref types.UnitRef) (string, error) {
	r, err := regexp.Compile(fmt.Sprintf("%s.lua$", ref))
	if err != nil {
		return "", err
	}

	file := ""
	err = fs.WalkDir(repoFiles, "bar-repo/units", func(path string, d os.DirEntry, err error) error {
		if err == nil && r.MatchString(path) {
			file = path
			return fs.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

func CounterpartForBuilding(ref types.UnitRef) []types.UnitRef {
	refs := make([]types.UnitRef, 0)
	constructors := gamedata.IsBuiltByUnits(ref)
	faction := FactionForRef(ref)
	for _, otherFaction := range OtherFactions(faction, false) {
		for _, constructorSuffix := range []string{"ca", "aca", "cv", "acv", "ck", "ack", "acsub"} {
			constructorRef := fmt.Sprintf("%s%s", PrefixForFaction(faction), constructorSuffix)
			path, exists := constructors[constructorRef]
			if !exists {
				continue
			}

			fConstructorRef := fmt.Sprintf("%s%s", PrefixForFaction(otherFaction), constructorSuffix)
			fConstructorGroup, exists := gamedata.GetUnitGrid()[fConstructorRef]
			if !exists {
				continue
			}

			refs = append(refs, fConstructorGroup[path[0]][path[1]][path[2]])
			break
		}
	}
	return refs
}
