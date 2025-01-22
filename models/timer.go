package models

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type timerModel struct {
	timer    timer.Model
	state    State
	pomodoro Pomodoro
}

func InitTimer(pomodoro Pomodoro) timerModel {

	pomodoroModel := timerModel{
		state:    Focus,
		pomodoro: pomodoro,
	}
	pomodoroModel.Init()
	return pomodoroModel
}

func (m timerModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		switch m.state {
		case Focus:
			// Switch to rest state
			m.timer = timer.NewWithInterval(m.pomodoro.rest, 1*time.Second)
			go playRestNotification()
			m.state = Rest
			return m, m.timer.Init()
		case Rest:
			// Switch to focus state
			m.timer = timer.NewWithInterval(m.pomodoro.focus, 1*time.Second)
			go playFocusNotification()
			m.state = Focus
			return m, m.timer.Init()
		}

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		winWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.timer.Toggle()
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return InitModel(), nil
		}

	}
	return m, nil

}

func (m timerModel) View() string {
	var message string
	if m.state == Focus {
		message = "Focus!"
	} else {
		message = "Rest"
	}

	screenStyle := lipgloss.NewStyle().Width(winWidth).Align(lipgloss.Center)
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
}
