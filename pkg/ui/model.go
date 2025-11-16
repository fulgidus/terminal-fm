// Package ui provides the terminal user interface components.
package ui

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
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
	
	// UI state
	view          ViewState
	width         int
	height        int
	
	// Station browsing
	stations      []radiobrowser.Station
	cursor        int
	scrollOffset  int
	loading       bool
	errorMsg      string
	
	// Search
	searchInput   string
	searchActive  bool
	
	// Bookmarks
	bookmarks        []radiobrowser.Station
	bookmarksLoading bool
}

// NewModel creates a new Model with initial state.
func NewModel(radioClient radiobrowser.Client, audioPlayer player.Player, store *storage.Store, locale string) Model {
	return Model{
		radioClient: radioClient,
		player:      audioPlayer,
		store:       store,
		locale:      locale,
		view:        ViewBrowse,
		stations:    []radiobrowser.Station{},
		cursor:      0,
		scrollOffset: 0,
		loading:     true,
		bookmarks:   []radiobrowser.Station{},
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

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// Helper methods for list navigation.

// VisibleStations returns the number of stations that fit on screen.
func (m Model) VisibleStations() int {
	// Reserve space for header (3 lines) and footer (2 lines)
	return m.height - 5
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
