package main

import (
	"fmt"
	"os"
	"strings"
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
	TaskManager
	TaskSelection
)

type model struct {
	options        []pomodoro
	timer          timer.Model
	selectedOption int
	state          State
	tasks          []string
	selectedTasks  []string
	newTask        string
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
		selectedTasks:  make([]string, 0),
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
		state:   Init,
		tasks:   make([]string, 0),
		newTask: "",
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
		if m.state == TaskManager {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl":
				return m, tea.Quit
			case "esc":
				m.state = Init
				return m, nil
			case "enter":
				if m.newTask != "" {
					m.tasks = append(m.tasks, m.newTask)
					m.newTask = ""
				}
			case "backspace":
				if len(m.newTask) > 0 {
					m.newTask = m.newTask[:len(m.newTask)-1]
				}
			default:
				m.newTask += msg.String()
			}
			return m, nil
		}

		if m.state == TaskSelection {
			switch msg.String() {
			case "j", "down":
				if m.selectedOption < len(m.selectedTasks)-1 {
					m.selectedOption++
				}

			case "k", "up":
				if m.selectedOption > 0 {
					m.selectedOption--
				}

			case "space":
				if len(m.tasks) > 0 {
					currentTask := m.tasks[m.selectedOption]
					if contains(m.selectedTasks, currentTask) {
						m.selectedTasks = remove(m.selectedTasks, currentTask)
					} else {
						m.selectedTasks = append(m.selectedTasks, currentTask)
					}
				}
				return m, nil
			case "enter":
				m.state = Focus
				return m, nil
			case "q":
				m.state = Init
				return m, nil
			}
		}

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
				m.state = TaskSelection
				return m, m.timer.Init()
			case Focus, Rest:
				return m, m.timer.Toggle()
			default:
				return m, nil
			}

		case "a":
			if m.state == Init {
				m.state = TaskManager
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

		legend := "\nPress 'a' to manage tasks"

		return screenStyle.Render(options + legend)

	case TaskSelection:
		var sb strings.Builder
		sb.WriteString("Which task will you be working on?\n\n")

		if len(m.tasks) == 0 {
			sb.WriteString("No tasks yet\n")
		} else {
			for i, task := range m.tasks {
				prefix := "  "
				if contains(m.selectedTasks, task) {
					prefix = "O "
				}
				sb.WriteString(fmt.Sprintf("%s%d. %s\n", prefix, i+1, task))
			}
		}

		sb.WriteString("\n\n\nPress space to select task, Enter to continue, q to go back")
		return screenStyle.Render(sb.String())

	case TaskManager:
		var sb strings.Builder
		sb.WriteString("Tasks\n")

		if len(m.tasks) == 0 {
			sb.WriteString("No tasks yet\n")
		} else {
			sb.WriteString(renderTasks(m))
		}

		sb.WriteString("\nNew task: " + m.newTask + "█\n")
		sb.WriteString("\n\n\nPress <Enter> to add task, <esc> to go back")

		return screenStyle.Render(sb.String())

	case Focus, Rest:
		var message string
		if m.state == Focus {
			message = "Focus!"
		} else {
			message = "Rest"
		}

		tasksView := screenStyle.
			MarginTop(1).
			Render(renderTasks(m))

		focusView := screenStyle.
			MarginTop(1).
			Render(message)

		timerView := screenStyle.
			MarginTop(1).
			Bold(true).
			Render(m.timer.View())

		pauseLine := screenStyle.
			MarginTop(1).
			Bold(true).
			Render("-PAUSE-")

		if !m.timer.Running() {
			return lipgloss.JoinVertical(lipgloss.Top, focusView, tasksView, timerView, pauseLine)
		}

		return lipgloss.JoinVertical(lipgloss.Top, focusView, tasksView, timerView)
	default:
		panic(fmt.Sprintf("unexpected State: %#v", m.state))
	}
}

func renderTasks(m model) string {
	var sb strings.Builder
	for i, task := range m.tasks {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, task))
	}
	return sb.String()
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func remove(slice []string, str string) []string {
	result := make([]string, 0)
	for _, v := range slice {
		if v != str {
			result = append(result, v)
		}
	}
	return result
}
