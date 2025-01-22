package models

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var winWidth int

type Pomodoro struct {
	description string
	focus       time.Duration
	rest        time.Duration
}

type State int

const (
	Init State = iota
	Focus
	Rest
)

type model struct {
	options        []Pomodoro
	timer          timer.Model
	selectedOption int
	state          State
}

func InitModel() model {
	return model{
		selectedOption: 0,
		options: []Pomodoro{
			{
				description: "25min focus / 5min rest",
				focus:       time.Duration(25 * time.Minute),
				rest:        time.Duration(5 * time.Minute),
			},
			{
				description: "30min focus / 10min rest",
				focus:       time.Duration(30 * time.Minute),
				rest:        time.Duration(10 * time.Minute),
			},
			{
				description: "30min focus / 5min rest",
				focus:       time.Duration(30 * time.Minute),
				rest:        time.Duration(5 * time.Minute),
			},
			{
				description: "45min focus / 15min rest",
				focus:       time.Duration(45 * time.Minute),
				rest:        time.Duration(15 * time.Minute),
			},
		},
		state: Init,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			return m, tea.Quit

		case "j", "down":
			if m.selectedOption < len(m.options)-1 && m.state == Init {
				m.selectedOption++
			}

		case "k", "up":
			if m.selectedOption > 0 && m.state == Init {
				m.selectedOption--
			}
		case "enter", " ":
			selectedPomodoro := m.options[m.selectedOption]
			m.timer = timer.NewWithInterval(selectedPomodoro.focus, 1*time.Second)
			return InitTimer(Pomodoro{
				focus: selectedPomodoro.focus,
				rest:  selectedPomodoro.rest,
			}).Update(msg)
		}
	case tea.WindowSizeMsg:
		winWidth = msg.Width
	}

	return m, nil
}

func (m model) View() string {
	screenStyle := lipgloss.NewStyle().Width(winWidth).Align(lipgloss.Center)
	options := ""
	for k, pomodoro := range m.options {
		if k == m.selectedOption {
			options += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ea9d34")).
				Render(pomodoro.description) + "\n"
		} else {
			options += pomodoro.description + "\n"
		}
	}
	options = screenStyle.Render(options)
	return options
}
