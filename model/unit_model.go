package model

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/wezzle/bar-unit-info/gamedata"
	"github.com/wezzle/bar-unit-info/gamedata/types"
	"github.com/wezzle/bar-unit-info/util"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	labelStyle       = lipgloss.NewStyle().Margin(0, 1, 0, 0).Foreground(lipgloss.Color("241"))
	descriptionStyle = lipgloss.NewStyle().Margin(1, 0, 0).Foreground(lipgloss.Color("245"))
	padding          = lipgloss.NewStyle().Margin(1, 0, 0)
	weaponStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#cc0000"))
	factionColors    = map[string]string{
		"Armada": "27",
		"Cortex": "124",
		"Legion": "34",
	}
)

func NewUnitModel(ref types.UnitRef, mainModel *MainModel) *Unit {
	m := Unit{}
	m.ref = ref
	m.mainModel = mainModel
	m.name = util.NameForRef(ref)
	m.description = util.DescriptionForRef(ref)
	var ok bool
	m.properties, ok = gamedata.GetUnitProperties(ref)
	if !ok {
		panic("unit properties file not generated")
	}

	m.faction = util.FactionForRef(ref)
	m.metalCost = progress.New(progress.WithSolidFill("#383C3F"), progress.WithoutPercentage())
	m.energyCost = progress.New(progress.WithSolidFill("#9E6802"), progress.WithoutPercentage())
	m.buildtime = progress.New(progress.WithSolidFill("#FEED53"), progress.WithoutPercentage())
	m.health = progress.New(progress.WithSolidFill("#49AE11"), progress.WithoutPercentage())
	m.sightRange = progress.New(progress.WithSolidFill("#C6C8C9"), progress.WithoutPercentage())
	m.speed = progress.New(progress.WithSolidFill("#1175AE"), progress.WithoutPercentage())
	m.buildpower = progress.New(progress.WithSolidFill("#6e17a3"), progress.WithoutPercentage())
	m.radarRange = progress.New(progress.WithSolidFill("#43e029"), progress.WithoutPercentage())
	m.jammerRange = progress.New(progress.WithSolidFill("#ea9896"), progress.WithoutPercentage())
	m.sonarRange = progress.New(progress.WithSolidFill("#29a3e8"), progress.WithoutPercentage())

	m.weaponDps = progress.New(progress.WithSolidFill("#cc0000"), progress.WithoutPercentage())
	m.weaponRange = progress.New(progress.WithSolidFill("#c3807f"), progress.WithoutPercentage())

	return &m
}

type Unit struct {
	ref         types.UnitRef
	name        string
	description string
	faction     string
	properties  types.UnitProperties

	mainModel *MainModel

	// Stats
	metalCost   progress.Model
	energyCost  progress.Model
	buildtime   progress.Model
	health      progress.Model
	sightRange  progress.Model
	speed       progress.Model
	buildpower  progress.Model
	radarRange  progress.Model
	jammerRange progress.Model
	sonarRange  progress.Model
	weaponDps   progress.Model
	weaponRange progress.Model
}

func (m *Unit) Init() tea.Cmd {
	return nil
}

func (m *Unit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m.mainModel.TableModel, cmd
		}
	}

	return m, cmd
}

func (m *Unit) RenderBar(labelWidth int, label string, progress string, maxValueWidth int, value string) string {
	v := value
	for range maxValueWidth - len(value) {
		v = " " + v
	}
	bar := []string{
		labelStyle.Width(labelWidth + 1).Render(fmt.Sprintf("%s:", label)),
		progress,
	}
	if value != "" {
		bar = append(bar, lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Margin(0, 0, 0, 1).Render(v))
	}
	return padding.Render(lipgloss.JoinHorizontal(lipgloss.Top, bar...))
}

func (m *Unit) PercentageWithBase(value int, base float64) float64 {
	return m.PercentageWithBaseF(float64(value), base)
}

func (m *Unit) PercentageWithBaseF(value float64, base float64) float64 {
	return min(value/base, 100.0) / 100
}

func (m *Unit) View() string {
	var sections []string

	var titleRow []string
	titleRow = append(titleRow, lipgloss.NewStyle().
		Background(lipgloss.Color(factionColors[m.faction])).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		Margin(0, 4, 0, 0).
		Render(m.faction))
	titleRow = append(titleRow, lipgloss.NewStyle().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		Render(m.name))
	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, titleRow...))

	description := descriptionStyle.Render(m.description)
	sections = append(sections, description)

	d := time.Second * time.Duration(m.properties.Buildtime/100)
	stats := [][]string{
		{"Metal cost", m.metalCost.ViewAs(m.PercentageWithBase(m.properties.MetalCost, 250)), strconv.Itoa(m.properties.MetalCost)},
		{"Energy cost", m.energyCost.ViewAs(m.PercentageWithBase(m.properties.EnergyCost, 900)), strconv.Itoa(m.properties.EnergyCost)},
		{"Buildtime", m.buildtime.ViewAs(m.PercentageWithBase(m.properties.Buildtime, 1000)), d.String()},
		{"Health", m.health.ViewAs(m.PercentageWithBase(m.properties.Health, 150)), strconv.Itoa(m.properties.Health)},
		{"Speed", m.speed.ViewAs(m.PercentageWithBaseF(m.properties.Speed, 1.5)), strconv.Itoa(m.properties.SightDistance)},
		{"Sight range", m.sightRange.ViewAs(m.PercentageWithBase(m.properties.SightDistance, 35)), strconv.Itoa(m.properties.SightDistance)},
	}

	if m.properties.RadarDistance != 0 {
		stats = append(stats, []string{"Radar range", m.radarRange.ViewAs(m.PercentageWithBase(m.properties.RadarDistance, 35)), strconv.Itoa(m.properties.RadarDistance)})
	}
	if m.properties.JammerDistance != 0 {
		stats = append(stats, []string{"Jammer range", m.jammerRange.ViewAs(m.PercentageWithBase(m.properties.JammerDistance, 10)), strconv.Itoa(m.properties.JammerDistance)})
	}
	if m.properties.SonarDistance != 0 {
		stats = append(stats, []string{"Sonar range", m.sonarRange.ViewAs(m.PercentageWithBase(m.properties.SonarDistance, 35)), strconv.Itoa(m.properties.SonarDistance)})
	}
	if m.properties.Buildpower != 0 {
		stats = append(stats, []string{"Buildpower", m.buildpower.ViewAs(m.PercentageWithBase(m.properties.Buildpower, 3)), strconv.Itoa(m.properties.Buildpower)})
	}

	weaponStats := [][]string{
		{"Weapons", weaponStyle.Render(m.properties.SummarizeWeaponTypes()), ""},
		{"DPS", m.weaponDps.ViewAs(m.PercentageWithBase(int(math.Round(m.properties.DPS())), 15)), strconv.Itoa(int(math.Round(m.properties.DPS())))},
		{"Weapon range", m.weaponRange.ViewAs(m.PercentageWithBase(int(m.properties.MaxWeaponRange()), 20)), strconv.Itoa(int(m.properties.MaxWeaponRange()))},
	}

	allStats := append(stats, weaponStats...)

	maxLabelWidth := 0
	maxValueWidth := 0
	for _, stat := range allStats {
		maxLabelWidth = max(ansi.StringWidth(stat[0]), maxLabelWidth)
		maxValueWidth = max(ansi.StringWidth(stat[2]), maxValueWidth)
	}

	for _, stat := range stats {
		sections = append(sections, m.RenderBar(maxLabelWidth, stat[0], stat[1], maxValueWidth, stat[2]))
	}

	weaponSections := make([]string, 0)
	for _, stat := range weaponStats {
		weaponSections = append(weaponSections, m.RenderBar(maxLabelWidth, stat[0], stat[1], maxValueWidth, stat[2]))
	}

	sections = append(sections, padding.Render(lipgloss.JoinVertical(lipgloss.Left, weaponSections...)))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
