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
	return gamedata.GetTranslations().Units.Factions[shortcode]
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
