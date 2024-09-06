package util

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lukegb/dds"
)

var Translations TranslationsT

func LoadTranslations(lang string) TranslationsT {
	f, err := os.Open(fmt.Sprintf("./bar-repo/language/%s/units.json", lang))
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Translations)
	if err != nil {
		panic(err)
	}
	return Translations
}

func NameForRef(ref UnitRef) string {
	return Translations.Units.Names[ref]
}

func DescriptionForRef(ref UnitRef) string {
	return Translations.Units.Descriptions[ref]
}

func FactionForRef(ref UnitRef) string {
	shortcode := ref[0:3]
	return Translations.Units.Factions[shortcode]
}

func LoadImage(ref UnitRef) image.Image {
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
