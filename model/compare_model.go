package model

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wezzle/bar-unit-info/gamedata"
	"golang.org/x/term"
)

var paddingStyle = lipgloss.NewStyle().Padding(0, 2)

type CompareKeyMap struct {
	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k CompareKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k CompareKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}

var compareKeys = CompareKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func NewCompareModel(mainModel *MainModel, refs ...string) CompareModel {
	m := CompareModel{
		mainModel:  mainModel,
		UnitModels: make([]*Unit, 0),
	}

	// Calculate base percentages
	bv := &BaseValues{}
	for _, r := range refs {
		up, ok := gamedata.GetUnitPropertiesByRef(r)
		if !ok {
			continue
		}
		bv.MetalCost = max(bv.MetalCost, float64(up.MetalCost)/100)
		bv.EnergyCost = max(bv.EnergyCost, float64(up.EnergyCost)/100)
		bv.Buildtime = max(bv.Buildtime, float64(up.Buildtime)/100)
		bv.Health = max(bv.Health, float64(up.Health)/100)
		bv.Speed = max(bv.Speed, float64(up.Speed)/100)
		bv.SightDistance = max(bv.SightDistance, float64(up.SightDistance)/100)
		bv.RadarDistance = max(bv.RadarDistance, float64(up.RadarDistance)/100)
		bv.JammerDistance = max(bv.JammerDistance, float64(up.JammerDistance)/100)
		bv.SonarDistance = max(bv.SonarDistance, float64(up.SonarDistance)/100)
		bv.Buildpower = max(bv.Buildpower, float64(up.Buildpower)/100)
		bv.DPS = max(bv.DPS, up.DPS()/100)
		bv.WeaponRange = max(bv.WeaponRange, up.MaxWeaponRange()/100)
	}

	components := make([]string, 0)
	for _, r := range refs {
		um := NewUnitModel(r, mainModel, bv)
		m.UnitModels = append(m.UnitModels, um)
		components = append(components, paddingStyle.Render(um.View()))
	}
	m.content = lipgloss.JoinHorizontal(lipgloss.Top, components...)

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	m.viewport = viewport.New(width, height)
	m.viewport.SetContent(m.content)
	m.ready = true

	return m
}

type CompareModel struct {
	UnitModels []*Unit

	viewport  viewport.Model
	mainModel *MainModel
	ready     bool
	content   string
}

func (m CompareModel) Init() tea.Cmd {
	return nil
}

func (m CompareModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, compareKeys.Quit):
			return m.mainModel.TableModel, cmd
		}
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m CompareModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}
