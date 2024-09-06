package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

func NewMainModel() MainModel {
	m := MainModel{}
	t := NewTableModel(&m)
	m.TableModel = &t
	m.activeModel = m.TableModel
	return m
}

type MainModel struct {
	TableModel *Table
	UnitModel  *Unit

	activeModel tea.Model
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	return m.activeModel.View()
}
