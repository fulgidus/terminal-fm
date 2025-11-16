package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the model (required by Bubbletea).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	// Window size changed
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	
	// Stations loaded successfully
	case stationsLoadedMsg:
		m.stations = msg.stations
		m.loading = false
		m.errorMsg = ""
		if len(m.stations) > 0 {
			m.cursor = 0
		}
		return m, nil
	
	// Error occurred
	case errMsg:
		m.loading = false
		m.errorMsg = msg.Error()
		return m, nil
	
	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	
	return m, nil
}

// handleKeyPress processes keyboard input based on current view.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts (work in all views)
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	
	// View-specific shortcuts
	switch m.view {
	case ViewBrowse:
		return m.handleBrowseKeys(msg)
	case ViewSearch:
		return m.handleSearchKeys(msg)
	case ViewBookmarks:
		return m.handleBookmarksKeys(msg)
	case ViewHelp:
		return m.handleHelpKeys(msg)
	}
	
	return m, nil
}

// handleBrowseKeys handles keyboard input in the browse view.
func (m Model) handleBrowseKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	
	// Navigation
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.UpdateScroll()
		}
		return m, nil
	
	case "down", "j":
		if m.cursor < len(m.stations)-1 {
			m.cursor++
			m.UpdateScroll()
		}
		return m, nil
	
	case "pgup":
		visible := m.VisibleStations()
		m.cursor -= visible
		if m.cursor < 0 {
			m.cursor = 0
		}
		m.UpdateScroll()
		return m, nil
	
	case "pgdown":
		visible := m.VisibleStations()
		m.cursor += visible
		if m.cursor >= len(m.stations) {
			m.cursor = len(m.stations) - 1
		}
		m.UpdateScroll()
		return m, nil
	
	case "home", "g":
		m.cursor = 0
		m.UpdateScroll()
		return m, nil
	
	case "end", "G":
		if len(m.stations) > 0 {
			m.cursor = len(m.stations) - 1
		}
		m.UpdateScroll()
		return m, nil
	
	// Actions (placeholders for now)
	case "enter", " ":
		// TODO: Play selected station in Week 2
		return m, nil
	
	case "s":
		// TODO: Stop playback in Week 2
		return m, nil
	
	case "a":
		// TODO: Add to bookmarks in Week 2
		return m, nil
	
	// View switching
	case "b":
		m.view = ViewBookmarks
		return m, nil
	
	case "/":
		m.view = ViewSearch
		return m, nil
	
	case "?":
		m.view = ViewHelp
		return m, nil
	}
	
	return m, nil
}

// handleSearchKeys handles keyboard input in the search view.
func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.view = ViewBrowse
		return m, nil
	
	// TODO: Implement search input in Week 2
	}
	
	return m, nil
}

// handleBookmarksKeys handles keyboard input in the bookmarks view.
func (m Model) handleBookmarksKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.view = ViewBrowse
		return m, nil
	
	// TODO: Implement bookmark navigation in Week 2
	}
	
	return m, nil
}

// handleHelpKeys handles keyboard input in the help view.
func (m Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "?":
		m.view = ViewBrowse
		return m, nil
	}
	
	return m, nil
}
