package main

import (
	"embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wezzle/bar-unit-info/model"
	"github.com/wezzle/bar-unit-info/util"
)

//go:embed bar-repo/luaui bar-repo/units bar-repo/language
var repoFiles embed.FS

func main() {
	util.InitFS(repoFiles)
	util.LoadTranslations("en")
	util.LoadGridLayouts()

	// p, err := util.LoadUnitProperties("legmineb")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%+v\n", p.CustomParams)
	// return

	m := model.NewMainModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
