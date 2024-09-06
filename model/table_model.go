package model

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wezzle/bar-unit-info/util"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewTableModel(mainModel *MainModel) Table {
	columns := []table.Column{
		{Title: "Ref ▼", Width: 20},
		{Title: "Name", Width: 30},
		{Title: "Metal cost", Width: 20},
		{Title: "Energy cost", Width: 20},
		{Title: "Buildtime", Width: 20},
		{Title: "Health", Width: 20},
		{Title: "Sight range", Width: 20},
		{Title: "Speed", Width: 20},
	}

	buildableUnits := make([]util.UnitRef, 0)
	properties := make(map[util.UnitRef]*util.UnitProperties)
	rows := make([]table.Row, 0)

	// Use labs to find buildable units
	// TODO we might miss units that are buildable by combat engineers and such
	for ref := range util.LabGrid {
		up, err := util.LoadUnitProperties(ref)
		if err != nil {
			continue
		}
		for _, ref := range up.BuildOptions {
			unitProperties, err := util.LoadUnitProperties(ref)
			if err != nil {
				continue
			}
			buildableUnits = append(buildableUnits, ref)
			properties[ref] = unitProperties
		}
	}

	sort.Strings(buildableUnits)
	buildableUnits = util.RemoveDuplicate(buildableUnits)

	for _, ref := range buildableUnits {
		up := properties[ref]
		d := time.Second * time.Duration(up.Buildtime/100)
		rows = append(rows, table.Row{
			ref,
			util.NameForRef(ref),
			strconv.Itoa(up.MetalCost),
			strconv.Itoa(up.EnergyCost),
			d.String(),
			strconv.Itoa(up.Health),
			strconv.Itoa(up.SightDistance),
			strconv.FormatFloat(up.Speed, 'f', 1, 64),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(40),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Table{
		t,
		0,
		false,
		mainModel,
	}
}

type Table struct {
	Table   table.Model
	SortCol int
	Reverse bool

	mainModel *MainModel
}

func (m *Table) Init() tea.Cmd {
	return nil
}

func (m *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	sortStringColumn := -1
	sortIntColumn := -1
	reverse := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return NewUnitModel(m.Table.SelectedRow()[0], m.mainModel), cmd
		case "f1":
			sortStringColumn = 0
		case "f2":
			sortStringColumn = 1
		case "f3":
			sortIntColumn = 2
		case "f4":
			sortIntColumn = 3
		case "f5":
			sortStringColumn = 4
		case "f6":
			sortIntColumn = 5
		case "f7":
			sortIntColumn = 6
		case "f8":
			sortIntColumn = 7
		}

		if sortStringColumn > -1 || sortIntColumn > -1 {
			col := max(sortStringColumn, sortIntColumn)
			if col == m.SortCol && !m.Reverse {
				reverse = true
			}
			cols := m.Table.Columns()
			for i := range cols {
				cols[i].Title = strings.ReplaceAll(cols[i].Title, " ▲", "")
				cols[i].Title = strings.ReplaceAll(cols[i].Title, " ▼", "")
			}
			if reverse {
				cols[col].Title = fmt.Sprintf("%s ▲", cols[col].Title)
			} else {
				cols[col].Title = fmt.Sprintf("%s ▼", cols[col].Title)
			}
			m.SortCol = col
			m.Reverse = reverse
		}
		if sortStringColumn > -1 {
			rows := m.Table.Rows()
			sort.Slice(rows, func(i, j int) bool {
				if reverse {
					return strings.ToLower(rows[i][sortStringColumn]) > strings.ToLower(rows[j][sortStringColumn])
				} else {
					return strings.ToLower(rows[i][sortStringColumn]) < strings.ToLower(rows[j][sortStringColumn])
				}
			})
			m.Table.SetRows(rows)
		}
		if sortIntColumn > -1 {
			rows := m.Table.Rows()
			sort.Slice(rows, func(i, j int) bool {
				intI, _ := strconv.Atoi(rows[i][sortIntColumn])
				intJ, _ := strconv.Atoi(rows[j][sortIntColumn])
				if reverse {
					return intI < intJ
				} else {
					return intI > intJ
				}
			})
			m.Table.SetRows(rows)

		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m *Table) View() string {
	return baseStyle.Render(m.Table.View()) + "\n  " + m.Table.HelpView() + "\n"
}
