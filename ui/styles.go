package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	primaryColor    = lipgloss.Color("#7C3AED") // Purple
	secondaryColor  = lipgloss.Color("#10B981") // Green
	accentColor     = lipgloss.Color("#F59E0B") // Amber
	errorColor      = lipgloss.Color("#EF4444") // Red
	textColor       = lipgloss.Color("#F9FAFB") // Light gray
	mutedColor      = lipgloss.Color("#9CA3AF") // Gray
	backgroundColor = lipgloss.Color("#1F2937") // Dark gray
)

// Base styles
var (
	// Main container style
	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			PaddingBottom(1)

	// ASCII art style
	asciiStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			PaddingBottom(2)

	// Menu item styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2).
			PaddingRight(2).
			PaddingTop(1).
			PaddingBottom(1).
			Margin(0, 2)

	selectedMenuItemStyle = menuItemStyle.Copy().
				Foreground(backgroundColor).
				Background(primaryColor).
				Bold(true)

	// Button styles
	buttonStyle = lipgloss.NewStyle().
			Foreground(backgroundColor).
			Background(secondaryColor).
			Bold(true).
			PaddingLeft(3).
			PaddingRight(3).
			PaddingTop(1).
			PaddingBottom(1).
			Margin(0, 1)

	selectedButtonStyle = buttonStyle.Copy().
				Background(primaryColor)

	// Card styles
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			PaddingLeft(2).
			PaddingRight(2).
			PaddingTop(1).
			PaddingBottom(1).
			Margin(1, 2)

	// Text styles
	textStyle = lipgloss.NewStyle().
			Foreground(textColor)

	mutedTextStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	emphasisStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// Progress bar style
	progressBarStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Background(mutedColor)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Align(lipgloss.Center).
			PaddingTop(1)
)

// ANKI ASCII art
const anktuiASCII = `
 █████╗ ███╗   ██╗██╗  ██╗████████╗██╗   ██╗██╗
██╔══██╗████╗  ██║██║ ██╔╝╚══██╔══╝██║   ██║██║
███████║██╔██╗ ██║█████╔╝    ██║   ██║   ██║██║
██╔══██║██║╚██╗██║██╔═██╗    ██║   ██║   ██║██║
██║  ██║██║ ╚████║██║  ██╗   ██║   ╚██████╔╝██║
╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝`

// Helper functions for centering content
func centerHorizontally(content string, width int) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, content)
}

func centerVertically(content string, height int) string {
	return lipgloss.PlaceVertical(height, lipgloss.Center, content)
}

func centerContent(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
