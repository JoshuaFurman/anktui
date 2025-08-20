package ui

import (
	"anktui/config"
	"anktui/models"
	"anktui/storage"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents different screens in the application
type Screen int

const (
	MenuScreen Screen = iota
	DeckListScreen
	StudyScreen
	DeckManagerScreen
	CardEditorScreen
)

// App represents the main application model
type App struct {
	config  *config.Config
	storage storage.Storage

	// UI state
	currentScreen Screen
	width         int
	height        int

	// Models for different screens
	menu        *MenuModel
	deckList    *DeckListModel
	study       *StudyModel
	deckManager *DeckManagerModel
	cardEditor  *CardEditorModel

	// Data
	decks          []*models.Deck
	currentDeck    *models.Deck
	currentSession *models.StudySession

	// Error state
	errorMessage string
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, store storage.Storage) *App {
	return &App{
		config:        cfg,
		storage:       store,
		currentScreen: MenuScreen,
		menu:          NewMenuModel(),
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	// Load all decks on startup
	return tea.Cmd(func() tea.Msg {
		decks, err := a.storage.LoadAllDecks()
		if err != nil {
			return ErrorMsg{err}
		}
		return DecksLoadedMsg{decks}
	})
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// Update all models with new size
		if a.menu != nil {
			a.menu.SetSize(msg.Width, msg.Height)
		}
		if a.deckList != nil {
			a.deckList.SetSize(msg.Width, msg.Height)
		}
		if a.study != nil {
			a.study.SetSize(msg.Width, msg.Height)
		}
		// TODO: Update other models

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if a.currentScreen == MenuScreen {
				return a, tea.Quit
			}
			// For other screens, go back to menu
			a.currentScreen = MenuScreen
			a.errorMessage = ""
			return a, nil
		}

	case DecksLoadedMsg:
		a.decks = msg.Decks
		a.errorMessage = ""

	case ErrorMsg:
		a.errorMessage = msg.Error.Error()

	case NavigateMsg:
		return a.handleNavigation(msg)

	case SaveDeckMsg:
		// Save the deck to storage
		return a, tea.Cmd(func() tea.Msg {
			if err := a.storage.SaveDeck(msg.Deck); err != nil {
				return ErrorMsg{err}
			}
			return nil
		})
	}

	// Route update to current screen
	switch a.currentScreen {
	case MenuScreen:
		if a.menu != nil {
			newModel, newCmd := a.menu.Update(msg)
			a.menu = newModel.(*MenuModel)
			cmd = newCmd
		}

	case DeckListScreen:
		if a.deckList == nil {
			a.deckList = NewDeckListModel(a.decks)
			a.deckList.SetSize(a.width, a.height)
		}
		newModel, newCmd := a.deckList.Update(msg)
		a.deckList = newModel.(*DeckListModel)
		cmd = newCmd

	case StudyScreen:
		if a.study != nil {
			newModel, newCmd := a.study.Update(msg)
			a.study = newModel.(*StudyModel)
			cmd = newCmd
		}

		// TODO: Handle other screens
	}

	return a, cmd
}

// View implements tea.Model
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	var content string

	// Show error message if there is one
	if a.errorMessage != "" {
		errorView := errorStyle.Render(fmt.Sprintf("Error: %s", a.errorMessage))
		content = centerContent(errorView, a.width, a.height)
		return content
	}

	// Route view to current screen
	switch a.currentScreen {
	case MenuScreen:
		if a.menu != nil {
			content = a.menu.View()
		}

	case DeckListScreen:
		if a.deckList != nil {
			content = a.deckList.View()
		}

	case StudyScreen:
		if a.study != nil {
			content = a.study.View()
		}

	// TODO: Handle other screens
	default:
		content = "Screen not implemented yet"
	}

	return content
}

// handleNavigation handles navigation messages between screens
func (a *App) handleNavigation(msg NavigateMsg) (tea.Model, tea.Cmd) {
	switch msg.Screen {
	case DeckListScreen:
		a.currentScreen = DeckListScreen
		if a.deckList == nil {
			a.deckList = NewDeckListModel(a.decks)
			a.deckList.SetSize(a.width, a.height)
		}

	case MenuScreen:
		a.currentScreen = MenuScreen

	case StudyScreen:
		a.currentScreen = StudyScreen
		if deck, ok := msg.Data.(*models.Deck); ok {
			a.currentDeck = deck
			a.study = NewStudyModel(deck, a.config.StudySession.CardsPerSession)
			a.study.SetSize(a.width, a.height)
		}

		// TODO: Handle other screen navigations
	}

	return a, nil
}

// Message types
type DecksLoadedMsg struct {
	Decks []*models.Deck
}

type ErrorMsg struct {
	Error error
}

type NavigateMsg struct {
	Screen Screen
	Data   interface{}
}
