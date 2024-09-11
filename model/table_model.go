package model

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wezzle/bar-unit-info/util"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	// Dialog.

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	// status bar

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Background(lipgloss.Color("#6124DF"))
)

type ColumnType int

const (
	CTString ColumnType = iota
	CTInt
	CTFloat
)

func ValueForRowAndColumn(row table.Row, column ColumnWithType, columnIndex int) any {
	ref := row[0]
	properties, _ := util.LoadUnitProperties(ref)
	val := column.ValueByPropertyKey(properties)
	if val == nil {
		return row[columnIndex]
	}
	return val
}

type ColumnWithType struct {
	table.Column
	Type        ColumnType
	PropertyKey string
}

func (c *ColumnWithType) ValueByPropertyKey(p *util.UnitProperties) any {
	switch c.PropertyKey {
	case "metalcost":
		return p.MetalCost
	case "energycost":
		return p.EnergyCost
	case "buildtime":
		return p.Buildtime
	case "techlevel":
		return p.CustomParams.TechLevel
	case "health":
		return p.Health
	case "sightdistance":
		return p.SightDistance
	case "speed":
		return p.Speed
	}
	return nil
}

type TableKeyMap struct {
	table.KeyMap
	Detail        key.Binding
	Left          key.Binding
	Right         key.Binding
	Help          key.Binding
	Quit          key.Binding
	Filter        key.Binding
	FilterConfirm key.Binding
	FilterCancel  key.Binding
	ToggleSort    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LineUp, k.LineDown, k.Left, k.Right, k.ToggleSort, k.Detail, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown, k.Left, k.Right, k.ToggleSort, k.Detail, k.Help, k.Quit},
		{k.GotoTop, k.GotoBottom, k.LineDown, k.PageDown, k.HalfPageUp, k.HalfPageDown},
	}
}

const spacebar = " "

var glyphs = []string{" ▼", " ▲", " •"}

var tableKeys = TableKeyMap{
	KeyMap: table.DefaultKeyMap(),
	Detail: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "show unit detail"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "select column left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "select column right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	FilterConfirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "confirm filter"),
	),
	FilterCancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<esc>", "cancel filter"),
	),
	ToggleSort: key.NewBinding(
		key.WithKeys(spacebar),
		key.WithHelp("<space>", "toggle sort"),
	),
}

func NewTableModel(mainModel *MainModel) Table {
	columns := []ColumnWithType{
		{Column: table.Column{Title: "Ref ▼ •", Width: 20}, Type: CTString},
		{Column: table.Column{Title: "Name", Width: 30}, Type: CTString},
		{Column: table.Column{Title: "Tech level", Width: 15}, Type: CTInt, PropertyKey: "techlevel"},
		{Column: table.Column{Title: "Metal cost", Width: 15}, Type: CTInt, PropertyKey: "metalcost"},
		{Column: table.Column{Title: "Energy cost", Width: 15}, Type: CTInt, PropertyKey: "energycost"},
		{Column: table.Column{Title: "Buildtime", Width: 15}, Type: CTInt, PropertyKey: "buildtime"},
		{Column: table.Column{Title: "Health", Width: 15}, Type: CTInt, PropertyKey: "health"},
		{Column: table.Column{Title: "Sight range", Width: 15}, Type: CTInt, PropertyKey: "sightdistance"},
		{Column: table.Column{Title: "Speed", Width: 15}, Type: CTFloat, PropertyKey: "speed"},
	}

	tableColumns := make([]table.Column, 0)
	for _, c := range columns {
		tableColumns = append(tableColumns, c.Column)
	}

	tableWidth := 0
	defaultCellPadding := 1
	defaultBorderWidth := 1
	for _, c := range columns {
		tableWidth = tableWidth + c.Width + (2 * defaultCellPadding)
	}
	tableWidth = tableWidth + defaultBorderWidth*2

	buildableUnits := make([]util.UnitRef, 0)
	properties := make(map[util.UnitRef]*util.UnitProperties)
	rows := make([]table.Row, 0)

	// Use labs to find buildable units
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

	// Check all buildable units for buildable units of their own
	for _, ref := range buildableUnits {
		up := properties[ref]
		for _, r := range up.BuildOptions {
			unitProperties, err := util.LoadUnitProperties(r)
			if err != nil {
				continue
			}
			buildableUnits = append(buildableUnits, r)
			properties[r] = unitProperties
		}
	}

	sort.Strings(buildableUnits)
	buildableUnits = util.RemoveDuplicate(buildableUnits)

	for _, ref := range buildableUnits {
		up, ok := properties[ref]
		if !ok {
			panic(ref)
		}
		d := time.Second * time.Duration(up.Buildtime/100)
		rows = append(rows, table.Row{
			ref,
			util.NameForRef(ref),
			fmt.Sprintf("T%d", up.CustomParams.TechLevel),
			strconv.Itoa(up.MetalCost),
			strconv.Itoa(up.EnergyCost),
			d.String(),
			strconv.Itoa(up.Health),
			strconv.Itoa(up.SightDistance),
			strconv.FormatFloat(up.Speed, 'f', 1, 64),
		})
	}

	t := table.New(
		table.WithColumns(tableColumns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(40),
		table.WithKeyMap(tableKeys.KeyMap),
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

	ti := textinput.New()
	ti.Placeholder = ""
	ti.CharLimit = 156
	ti.Width = 20
	ti.TextStyle = fishCakeStyle
	ti.Cursor.TextStyle = fishCakeStyle
	ti.Prompt = "/ "

	return Table{
		Table:         t,
		FilterInput:   ti,
		SortCol:       0,
		SelectedCol:   0,
		Reverse:       false,
		DialogShown:   false,
		FilterMode:    false,
		mainModel:     mainModel,
		help:          help.New(),
		tableWidth:    tableWidth,
		columns:       columns,
		rows:          rows,
		columnFilters: make([]string, len(columns)),
	}
}

type Table struct {
	Table       table.Model
	FilterInput textinput.Model
	SortCol     int
	Reverse     bool
	DialogShown bool
	FilterMode  bool
	SelectedCol int

	mainModel  *MainModel
	help       help.Model
	width      int
	height     int
	tableWidth int

	columns       []ColumnWithType
	columnFilters []string
	rows          []table.Row
}

func (m *Table) FilterRows(cf []string) {
	cf[m.SelectedCol] = m.FilterInput.Value()
	filteredRows := make([]table.Row, len(m.rows))
	copy(filteredRows, m.rows)
	for colIndex, f := range cf {
		if f == "" {
			continue
		}
		var filtered []table.Row
		for _, r := range filteredRows {
			found := false
			switch m.columns[colIndex].Type {
			case CTString:
				re := regexp.MustCompile(fmt.Sprintf("(?i)%s", f))
				found = re.Match([]byte(r[colIndex]))
			case CTInt:
				val := ValueForRowAndColumn(r, m.columns[colIndex], colIndex).(int)
				cleanFilterString := string(regexp.MustCompile("[><= ]+").ReplaceAll([]byte(f), []byte("")))
				filterVal, _ := strconv.Atoi(cleanFilterString)
				if strings.Contains(f, ">") {
					found = val > filterVal
				} else if strings.Contains(f, "<") {
					found = val < filterVal
				} else {
					found = val == filterVal
				}
			case CTFloat:
				val := ValueForRowAndColumn(r, m.columns[colIndex], colIndex).(float64)
				cleanFilterString := string(regexp.MustCompile("[><= ]+").ReplaceAll([]byte(f), []byte("")))
				filterVal, _ := strconv.ParseFloat(cleanFilterString, 64)
				if strings.Contains(f, ">") {
					found = val > filterVal
				} else if strings.Contains(f, "<") {
					found = val < filterVal
				} else {
					found = val == filterVal
				}
			}
			if found {
				filtered = append(filtered, r)
			}
		}
		filteredRows = filtered
	}
	m.Table.SetRows(filteredRows)
}

func (m *Table) Init() tea.Cmd {
	return nil
}

func (m *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	reverse := false
	var selectedCol *int
	var sortCol *int

	if m.FilterMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, tableKeys.FilterConfirm):
				m.columnFilters[m.SelectedCol] = m.FilterInput.Value()
				m.FilterMode = false
				return m, cmd
			case key.Matches(msg, tableKeys.FilterCancel):
				m.FilterInput.SetValue("")
				m.FilterMode = false
				m.FilterRows(m.columnFilters)
				return m, cmd
			}
		}

		m.FilterInput, cmd = m.FilterInput.Update(msg)

		cf := make([]string, len(m.columnFilters))
		copy(cf, m.columnFilters)
		m.FilterRows(cf)

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		// case key.Matches(msg, tableKeys.Help):
		// 	if m.Table.Focused() {
		// 		m.Table.Blur()
		// 	} else {
		// 		m.Table.Focus()
		// 	}
		case key.Matches(msg, tableKeys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, tableKeys.Filter):
			// m.DialogShown = true
			m.FilterMode = true
			m.FilterInput.Focus()
			m.FilterInput.SetValue(m.columnFilters[m.SelectedCol])
		case key.Matches(msg, tableKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, tableKeys.Detail):
			return NewUnitModel(m.Table.SelectedRow()[0], m.mainModel), cmd
		case key.Matches(msg, tableKeys.Left):
			s := max(m.SelectedCol-1, 0)
			selectedCol = &s
		case key.Matches(msg, tableKeys.Right):
			s := min(m.SelectedCol+1, len(m.Table.Columns()))
			selectedCol = &s
		case key.Matches(msg, tableKeys.ToggleSort):
			sortCol = &m.SelectedCol
		}

		if sortCol != nil {
			if *sortCol == m.SortCol {
				reverse = !m.Reverse
			}

			cols := m.Table.Columns()
			// Remove sort icons
			for i := range cols {
				cols[i].Title = strings.ReplaceAll(cols[i].Title, " ▲", "")
				cols[i].Title = strings.ReplaceAll(cols[i].Title, " ▼", "")
			}
			// Add sort icon to sortCol
			format := "%s ▼"
			if reverse {
				format = "%s ▲"
			}
			cols[*sortCol].Title = fmt.Sprintf(format, cols[*sortCol].Title)
			m.SortCol = *sortCol
			m.Reverse = reverse

			t := m.columns[*sortCol].Type
			rows := m.Table.Rows()
			sort.Slice(rows, func(i, j int) bool {
				iVal := ValueForRowAndColumn(rows[i], m.columns[*sortCol], *sortCol)
				jVal := ValueForRowAndColumn(rows[j], m.columns[*sortCol], *sortCol)
				var isLess bool
				switch t {
				case CTInt:
					isLess = iVal.(int) < jVal.(int)
				case CTFloat:
					isLess = iVal.(float64) < jVal.(float64)
				default:
					isLess = strings.ToLower(iVal.(string)) < strings.ToLower(jVal.(string))
				}
				if reverse {
					return !isLess
				}
				return isLess
			})
			m.Table.SetRows(rows)
		}

		selectedColUpdate := selectedCol != nil && m.SelectedCol != *selectedCol
		selectedColIsToggled := sortCol != nil && *sortCol == m.SelectedCol
		if selectedColIsToggled {
			selectedCol = &m.SelectedCol
		}
		if selectedColUpdate || selectedColIsToggled {
			cols := m.Table.Columns()
			for i := range cols {
				cols[i].Title = strings.ReplaceAll(cols[i].Title, " •", "")
			}
			cols[*selectedCol].Title = fmt.Sprintf("%s •", cols[*selectedCol].Title)
			m.SelectedCol = *selectedCol
		}
	}
	if sortCol != nil {
		// Exit so we don't run table's keymap
		return m, cmd
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m *Table) View() string {
	// physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	doc := strings.Builder{}
	if m.DialogShown {
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog := lipgloss.Place(m.tableWidth, m.Table.Height(),
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("░"),
			lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}),
		)

		doc.WriteString(dialog + "\n\n")
	} else {
		doc.WriteString(baseStyle.Render(m.Table.View()))

		cleanColTitle := m.columns[m.SelectedCol].Title
		for _, g := range glyphs {
			cleanColTitle = strings.Replace(cleanColTitle, g, "", 1)
		}

		w := lipgloss.Width

		var fishCake string
		if m.FilterMode {
			fishCake = fishCakeStyle.Render(m.FilterInput.View())
		} else if m.columnFilters[m.SelectedCol] != "" {
			fishCake = fishCakeStyle.Render(fmt.Sprintf("Column <%s> is filtered by: %s", cleanColTitle, m.columnFilters[m.SelectedCol]))
		} else {
			fishCake = fishCakeStyle.Render(fmt.Sprintf("Add filter to <%s> by pressing /", cleanColTitle))
		}
		statusVal := statusText.
			Width(m.tableWidth - w(fishCake)).
			Render("")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusVal,
			fishCake,
		)

		doc.WriteString("\n")
		doc.WriteString(statusBarStyle.Width(m.tableWidth).Render(bar))
		doc.WriteString("\n\n")
		doc.WriteString(m.help.View(tableKeys))
	}

	return doc.String()
}
