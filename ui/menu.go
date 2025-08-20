package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuOption represents a menu option
type MenuOption struct {
	Label       string
	Description string
	Action      func() tea.Msg
}

// MenuModel represents the main menu
type MenuModel struct {
	options  []MenuOption
	selected int
	width    int
	height   int
}

// NewMenuModel creates a new menu model
func NewMenuModel() *MenuModel {
	return &MenuModel{
		options: []MenuOption{
			{
				Label:       "Study Decks",
				Description: "Review and study your flashcards",
				Action: func() tea.Msg {
					return NavigateMsg{Screen: DeckListScreen}
				},
			},
			{
				Label:       "Manage Decks",
				Description: "Create, edit, and organize your decks",
				Action: func() tea.Msg {
					return NavigateMsg{Screen: DeckManagerScreen}
				},
			},
			{
				Label:       "Statistics",
				Description: "View your learning progress",
				Action: func() tea.Msg {
					// TODO: Implement statistics screen
					return NavigateMsg{Screen: MenuScreen}
				},
			},
			{
				Label:       "Quit",
				Description: "Exit AnkTUI",
				Action: func() tea.Msg {
					return tea.Quit()
				},
			},
		},
		selected: 0,
	}
}

// SetSize sets the terminal size
func (m *MenuModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.options)-1 {
				m.selected++
			}
		case "enter", " ":
			return m, m.options[m.selected].Action
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *MenuModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Create the ASCII art title
	asciiArt := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(anktuiASCII)

	// Create menu items
	var menuItems []string
	for i, option := range m.options {
		style := menuItemStyle
		if i == m.selected {
			style = lipgloss.NewStyle().
				Foreground(backgroundColor).
				Background(primaryColor).
				Bold(true).
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				PaddingBottom(1).
				Margin(0, 2)
		}

		menuItems = append(menuItems, style.Render(option.Label))
	}

	// Join menu items
	menu := lipgloss.JoinVertical(lipgloss.Center, menuItems...)

	// Add description for selected item
	description := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(2).
		Render(m.options[m.selected].Description)

	// Add help text
	help := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(4).
		Render("Use ↑/↓ arrows or j/k to navigate • Enter to select • q to quit")

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		asciiArt,
		menu,
		description,
		help,
	)

	// Center everything in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
