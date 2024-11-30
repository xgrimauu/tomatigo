package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var winWidth int

type pomodoro struct {
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
	options        []pomodoro
	timer          timer.Model
	selectedOption int
	state          State
}

func main() {
	program := tea.NewProgram(initModel(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		os.Exit(1)
	}
}

func initModel() model {
	return model{
		selectedOption: 0,
		options: []pomodoro{
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
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		switch m.state {
		case Focus:
			// Switch to rest state
			selectedPomodoro := m.options[m.selectedOption]
			m.timer = timer.NewWithInterval(selectedPomodoro.rest, 1*time.Second)
			go playRestNotification()
			m.state = Rest
			return m, m.timer.Init()
		case Rest:
			// Switch to focus state
			selectedPomodoro := m.options[m.selectedOption]
			m.timer = timer.NewWithInterval(selectedPomodoro.focus, 1*time.Second)
			go playFocusNotification()
			m.state = Focus
			return m, m.timer.Init()
		}

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.state == Init {
				return m, tea.Quit
			}
			m.state = Init
			return m, nil

		case "j", "down":
			if m.selectedOption < len(m.options)-1 && m.state == Init {
				m.selectedOption++
			}

		case "k", "up":
			if m.selectedOption > 0 && m.state == Init {
				m.selectedOption--
			}
		case "enter", " ":
			switch m.state {
			case Init:
				selectedPomodoro := m.options[m.selectedOption]
				m.timer = timer.NewWithInterval(selectedPomodoro.focus, 1*time.Second)
				m.state = Focus
				return m, m.timer.Init()
			case Focus, Rest:
				return m, m.timer.Toggle()
			default:
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		winWidth = msg.Width
	}

	return m, nil
}

func (m model) View() string {
	screenStyle := lipgloss.NewStyle().Width(winWidth).Align(lipgloss.Center)
	switch m.state {
	case Init:
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

	case Focus, Rest:
		var message string
		if m.state == Focus {
			message = "Focus!"
		} else {
			message = "Rest"
		}

		focusView := screenStyle.
			MarginTop(4).
			Render(message)

		timerView := screenStyle.
			MarginTop(2).
			Bold(true).
			Render(m.timer.View())

		if !m.timer.Running() {
			pauseView := screenStyle.
				MarginTop(1).
				Bold(true).
				Render("-PAUSE-")
			return lipgloss.JoinVertical(lipgloss.Center, focusView, timerView, pauseView)
		}

		return lipgloss.JoinVertical(lipgloss.Center, focusView, timerView)
	default:
		panic(fmt.Sprintf("unexpected State: %#v", m.state))
	}
}
