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
	defaultBaseValues = map[string]float64{}
)

type BaseValues struct {
	MetalCost      float64
	EnergyCost     float64
	Buildtime      float64
	Health         float64
	Speed          float64
	SightDistance  float64
	RadarDistance  float64
	JammerDistance float64
	SonarDistance  float64
	Buildpower     float64
	DPS            float64
	EPS            float64
	MPS            float64
	ParalyzeTime   float64
	WeaponRange    float64
}

func NewUnitModel(ref types.UnitRef, mainModel *MainModel, baseValues *BaseValues) *Unit {
	m := Unit{}
	m.ref = ref
	m.mainModel = mainModel
	m.name = util.NameForRef(ref)
	m.description = util.DescriptionForRef(ref)
	var ok bool
	m.properties, ok = gamedata.GetUnitPropertiesByRef(ref)
	if !ok {
		panic("unit properties file not generated")
	}

	if baseValues == nil {
		baseValues = &BaseValues{
			MetalCost:      250,
			EnergyCost:     900,
			Buildtime:      1000,
			Health:         150,
			Speed:          1.5,
			SightDistance:  35,
			RadarDistance:  35,
			JammerDistance: 10,
			SonarDistance:  35,
			Buildpower:     3,
			EPS:            250,
			MPS:            250,
			ParalyzeTime:   35,
		}
	}
	m.baseValues = baseValues

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
	m.weaponEps = progress.New(progress.WithSolidFill("#9E6802"), progress.WithoutPercentage())
	m.weaponMps = progress.New(progress.WithSolidFill("#383C3F"), progress.WithoutPercentage())
	m.weaponParalyzeTime = progress.New(progress.WithSolidFill("#1175AE"), progress.WithoutPercentage())

	return &m
}

type Unit struct {
	ref         types.UnitRef
	name        string
	description string
	faction     string
	properties  *types.UnitProperties

	mainModel *MainModel

	// Stats
	metalCost          progress.Model
	energyCost         progress.Model
	buildtime          progress.Model
	health             progress.Model
	sightRange         progress.Model
	speed              progress.Model
	buildpower         progress.Model
	radarRange         progress.Model
	jammerRange        progress.Model
	sonarRange         progress.Model
	weaponDps          progress.Model
	weaponRange        progress.Model
	weaponEps          progress.Model
	weaponMps          progress.Model
	weaponParalyzeTime progress.Model
	// TODO
	// paralyzer
	// TODO
	// stockpile
	// stockpileTime

	baseValues *BaseValues
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

func (m *Unit) PercentageWithBase(value int64, base float64) float64 {
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
		Margin(0, 4, 0, 0).
		Render(m.name))
	titleRow = append(titleRow, lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("246")).
		Padding(0, 1).
		Render(m.ref))
	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Top, titleRow...))

	description := descriptionStyle.Render(m.description)
	sections = append(sections, description)

	d := time.Second * time.Duration(m.properties.Buildtime/100)
	stats := [][]string{
		{"Metal cost", m.metalCost.ViewAs(m.PercentageWithBase(m.properties.MetalCost, m.baseValues.MetalCost)), strconv.FormatInt(m.properties.MetalCost, 10)},
		{"Energy cost", m.energyCost.ViewAs(m.PercentageWithBase(m.properties.EnergyCost, m.baseValues.EnergyCost)), strconv.FormatInt(m.properties.EnergyCost, 10)},
		{"Buildtime", m.buildtime.ViewAs(m.PercentageWithBase(m.properties.Buildtime, m.baseValues.Buildtime)), d.String()},
		{"Health", m.health.ViewAs(m.PercentageWithBase(m.properties.Health, m.baseValues.Health)), strconv.FormatInt(m.properties.Health, 10)},
		{"Speed", m.speed.ViewAs(m.PercentageWithBaseF(m.properties.Speed, m.baseValues.Speed)), strconv.FormatFloat(m.properties.Speed, 'f', 1, 64)},
		{"Sight range", m.sightRange.ViewAs(m.PercentageWithBase(m.properties.SightDistance, m.baseValues.SightDistance)), strconv.FormatInt(m.properties.SightDistance, 10)},
	}

	if m.properties.RadarDistance != 0 {
		stats = append(stats, []string{"Radar range", m.radarRange.ViewAs(m.PercentageWithBase(m.properties.RadarDistance, m.baseValues.RadarDistance)), strconv.FormatInt(m.properties.RadarDistance, 10)})
	}
	if m.properties.JammerDistance != 0 {
		stats = append(stats, []string{"Jammer range", m.jammerRange.ViewAs(m.PercentageWithBase(m.properties.JammerDistance, m.baseValues.JammerDistance)), strconv.FormatInt(m.properties.JammerDistance, 10)})
	}
	if m.properties.SonarDistance != 0 {
		stats = append(stats, []string{"Sonar range", m.sonarRange.ViewAs(m.PercentageWithBase(m.properties.SonarDistance, m.baseValues.SonarDistance)), strconv.FormatInt(m.properties.SonarDistance, 10)})
	}
	if m.properties.Buildpower != 0 {
		stats = append(stats, []string{"Buildpower", m.buildpower.ViewAs(m.PercentageWithBase(m.properties.Buildpower, m.baseValues.Buildpower)), strconv.FormatInt(m.properties.Buildpower, 10)})
	}

	ws := weaponStyle
	if m.properties.ParalyzeTime() != 0.0 {
		ws = ws.Foreground(lipgloss.Color("27"))
	}
	weaponStats := [][]string{
		{"Weapons", ws.Render(m.properties.SummarizeWeaponTypes()), ""},
		{"DPS", m.weaponDps.ViewAs(m.PercentageWithBase(int64(math.Round(m.properties.DPS())), m.baseValues.DPS)), strconv.Itoa(int(math.Round(m.properties.DPS())))},
		{"Weapon range", m.weaponRange.ViewAs(m.PercentageWithBase(int64(m.properties.MaxWeaponRange()), m.baseValues.WeaponRange)), strconv.Itoa(int(m.properties.MaxWeaponRange()))},
	}
	if m.properties.MPS() != 0.0 {
		weaponStats = append(weaponStats, []string{"Metal/s", m.weaponMps.ViewAs(m.PercentageWithBase(int64(math.Round(m.properties.MPS())), m.baseValues.MPS)), strconv.Itoa(int(math.Round(m.properties.MPS())))})
	}
	if m.properties.EPS() != 0.0 {
		weaponStats = append(weaponStats, []string{"Energy/s", m.weaponEps.ViewAs(m.PercentageWithBase(int64(math.Round(m.properties.EPS())), m.baseValues.EPS)), strconv.Itoa(int(math.Round(m.properties.EPS())))})
	}
	if m.properties.ParalyzeTime() != 0.0 {
		weaponStats = append(weaponStats, []string{"Paralyze time", m.weaponParalyzeTime.ViewAs(m.PercentageWithBase(int64(m.properties.ParalyzeTime()), m.baseValues.ParalyzeTime)), strconv.FormatInt(m.properties.ParalyzeTime(), 10)})
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
