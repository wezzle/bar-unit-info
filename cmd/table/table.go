package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wezzle/bar-unit-info/model"
	"github.com/wezzle/bar-unit-info/util"
)

func main() {
	util.LoadTranslations("en")
	util.LoadGridLayouts()

	m := model.NewMainModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}