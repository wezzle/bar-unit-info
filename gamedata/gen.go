//go:build ignore

package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wezzle/bar-unit-info/gamedata/parser"
)

func main() {
	templates, err := filepath.Glob("templates/*.go.tmpl")
	if err != nil {
		panic(err)
	}

	for _, t := range templates {
		tpl, err := template.ParseFiles(t)
		if err != nil {
			panic(err)
		}

		base := filepath.Base(t)
		filename := strings.TrimSuffix(base, filepath.Ext(base))
		slog.Info("creating template", "file", filename)

		f, err := os.Create(filename)
		if err != nil {
			slog.Error("failed to create file", "error", err)
		}

		data := struct {
			Var string
		}{}
		switch filename {
		case "labgrid.go":
			_, labGrid := parser.LoadGridLayouts()
			data.Var = fmt.Sprintf("%#v\n", labGrid)
		case "unitgrid.go":
			unitGrid, _ := parser.LoadGridLayouts()
			data.Var = fmt.Sprintf("%#v\n", unitGrid)
		case "unitproperties.go":
			upByRef := parser.LoadAllUnitProperties()
			data.Var = fmt.Sprintf("%#v\n", upByRef)
		case "translations.go":
			t := parser.LoadTranslations("en")
			data.Var = fmt.Sprintf("%#v\n", t)
		}

		tpl.Execute(f, data)
		f.Close()
	}
}