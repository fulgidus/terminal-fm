// Package ui provides the terminal user interface components.
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fulgidus/terminal-fm/pkg/i18n"
	"github.com/fulgidus/terminal-fm/pkg/services/player"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
	"github.com/fulgidus/terminal-fm/pkg/services/storage"
)

// ViewState represents the current screen being displayed.
type ViewState int

const (
	// ViewBrowse shows the station list.
	ViewBrowse ViewState = iota
	// ViewSearch shows the search interface.
	ViewSearch
	// ViewBookmarks shows bookmarked stations.
	ViewBookmarks
	// ViewHelp shows keyboard shortcuts.
	ViewHelp
)

// Model holds the application state for the TUI.
type Model struct {
	// Core dependencies
	radioClient radiobrowser.Client
	player      player.Player
	store       *storage.Store
	locale      string
	tr          *i18n.SimpleTranslator

	// UI state
	view   ViewState
	width  int
	height int

	// Station browsing
	stations     []radiobrowser.Station
	cursor       int
	scrollOffset int
	loading      bool
	errorMsg     string

	// Search
	searchInput        textinput.Model
	searchResults      []radiobrowser.Station
	searchCursor       int
	searchScrollOffset int
	searching          bool

	// Bookmarks
	bookmarks        []radiobrowser.Station
	bookmarksLoading bool
}

// NewModel creates a new Model with initial state.
func NewModel(radioClient radiobrowser.Client, audioPlayer player.Player, store *storage.Store, locale string) Model {
	// Initialize translator
	tr := i18n.NewSimpleTranslator(locale)

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = tr.T("search.placeholder")
	ti.CharLimit = 100
	ti.Width = 50

	return Model{
		radioClient:   radioClient,
		player:        audioPlayer,
		store:         store,
		locale:        locale,
		tr:            tr,
		view:          ViewBrowse,
		stations:      []radiobrowser.Station{},
		cursor:        0,
		scrollOffset:  0,
		loading:       true,
		bookmarks:     []radiobrowser.Station{},
		searchInput:   ti,
		searchResults: []radiobrowser.Station{},
	}
}

// Init initializes the model (required by Bubbletea).
func (m Model) Init() tea.Cmd {
	// Load stations on startup
	return m.loadStations
}

// loadStations is a command that fetches stations from the API.
func (m Model) loadStations() tea.Msg {
	stations, err := m.radioClient.Search(radiobrowser.SearchParams{
		Limit: 50,
		Order: "votes",
	})
	if err != nil {
		return errMsg{err}
	}
	return stationsLoadedMsg{stations}
}

// loadBookmarks is a command that loads bookmarks from storage.
func (m Model) loadBookmarks() tea.Msg {
	if m.store == nil {
		return errMsg{fmt.Errorf("storage not available")}
	}

	bookmarks, err := m.store.GetBookmarks()
	if err != nil {
		return errMsg{err}
	}
	return bookmarksLoadedMsg{bookmarks}
}

// performSearch executes a search query.
func (m Model) performSearch(query string) tea.Msg {
	if query == "" {
		return searchResultsMsg{[]radiobrowser.Station{}}
	}

	// Try to determine search type
	params := radiobrowser.SearchParams{
		Limit: 50,
		Order: "votes",
	}

	// Simple heuristic: if it's 2 letters, assume country code
	if len(query) == 2 {
		params.Country = query
	} else {
		// Otherwise search by name
		params.Name = query
	}

	stations, err := m.radioClient.Search(params)
	if err != nil {
		return errMsg{err}
	}

	return searchResultsMsg{stations}
}

// Message types for async operations.
type stationsLoadedMsg struct {
	stations []radiobrowser.Station
}

type bookmarksLoadedMsg struct {
	bookmarks []radiobrowser.Station
}

type bookmarkAddedMsg struct {
	station radiobrowser.Station
}

type bookmarkRemovedMsg struct {
	stationUUID string
}

type searchResultsMsg struct {
	results []radiobrowser.Station
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// Helper methods for list navigation.

// VisibleStations returns the number of stations that fit on screen.
func (m Model) VisibleStations() int {
	// Reserve space for:
	// - Title (1 line)
	// - Status bar with top/bottom borders (3 lines)
	// - Status bar margin bottom (1 line)
	// - Header info (1 line)
	// - Spacing before list (1 line)
	// - Spacing after list (1 line)
	// - Footer with border (2 lines)
	reserved := 10
	
	// Add extra line for error message if present
	if m.errorMsg != "" {
		reserved += 1
	}
	
	visible := m.height - reserved
	if visible < 1 {
		visible = 1
	}
	return visible
}

// CanScrollUp returns true if we can scroll up.
func (m Model) CanScrollUp() bool {
	return m.scrollOffset > 0
}

// CanScrollDown returns true if we can scroll down.
func (m Model) CanScrollDown() bool {
	visible := m.VisibleStations()
	return m.scrollOffset+visible < len(m.stations)
}

// UpdateScroll adjusts scrollOffset based on cursor position.
func (m *Model) UpdateScroll() {
	visible := m.VisibleStations()

	// Scroll down if cursor is below visible area
	if m.cursor >= m.scrollOffset+visible {
		m.scrollOffset = m.cursor - visible + 1
	}

	// Scroll up if cursor is above visible area
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
}

// SelectedStation returns the currently selected station, or nil.
func (m Model) SelectedStation() *radiobrowser.Station {
	if len(m.stations) == 0 || m.cursor < 0 || m.cursor >= len(m.stations) {
		return nil
	}
	return &m.stations[m.cursor]
}

// VisibleStationList returns the slice of stations currently visible.
func (m Model) VisibleStationList() []radiobrowser.Station {
	if len(m.stations) == 0 {
		return []radiobrowser.Station{}
	}

	visible := m.VisibleStations()
	end := m.scrollOffset + visible
	if end > len(m.stations) {
		end = len(m.stations)
	}

	return m.stations[m.scrollOffset:end]
}

// Cleanup stops playback and cleans up resources.
func (m *Model) Cleanup() {
	if m.player != nil {
		m.player.Stop()
		// If player implements Cleanup interface, call it
		if cleaner, ok := m.player.(interface{ Cleanup() error }); ok {
			cleaner.Cleanup()
		}
	}
}
