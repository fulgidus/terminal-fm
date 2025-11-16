package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
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

	// Bookmarks loaded successfully
	case bookmarksLoadedMsg:
		m.bookmarks = msg.bookmarks
		m.bookmarksLoading = false
		// Reset cursor if it's out of bounds
		if m.bookmarksCursor >= len(m.bookmarks) {
			m.bookmarksCursor = 0
			m.bookmarksScrollOffset = 0
		}
		return m, nil

	// Bookmark added
	case bookmarkAddedMsg:
		m.errorMsg = fmt.Sprintf("Added '%s' to bookmarks", msg.station.Name)
		// Reload bookmarks
		return m, m.loadBookmarks

	// Bookmark removed
	case bookmarkRemovedMsg:
		m.errorMsg = "Removed from bookmarks"
		// Reload bookmarks
		return m, m.loadBookmarks

	// Search results received
	case searchResultsMsg:
		m.searchResults = msg.results
		m.searching = false
		m.searchCursor = 0
		m.searchScrollOffset = 0
		if len(msg.results) == 0 {
			m.errorMsg = "No stations found"
		} else {
			m.errorMsg = ""
		}
		return m, nil

	// Error occurred
	case errMsg:
		m.loading = false
		m.bookmarksLoading = false
		m.searching = false
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
		// Cleanup before quitting
		m.Cleanup()
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
	case ViewAbout:
		return m.handleAboutKeys(msg)
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

	// Actions
	case "enter", " ":
		// Play/pause selected station
		station := m.SelectedStation()
		if station != nil {
			currentStation := m.player.GetCurrentStation()
			if currentStation != nil && currentStation.StationUUID == station.StationUUID {
				// Stop if already playing this station
				_ = m.player.Stop()
			} else {
				// Play the selected station
				if err := m.player.Play(station); err != nil {
					m.errorMsg = fmt.Sprintf("Failed to play: %v", err)
				} else {
					m.errorMsg = ""
				}
			}
		}
		return m, nil

	case "s":
		// Stop playback
		_ = m.player.Stop()
		m.errorMsg = ""
		return m, nil

	case "=", "+":
		// Increase volume
		currentVol := m.player.GetVolume()
		if currentVol < 100 {
			_ = m.player.SetVolume(currentVol + 10)
		}
		return m, nil

	case "-", "_":
		// Decrease volume
		currentVol := m.player.GetVolume()
		if currentVol > 0 {
			_ = m.player.SetVolume(currentVol - 10)
		}
		return m, nil

	case "a":
		// Add/remove bookmark
		station := m.SelectedStation()
		if station != nil && m.store != nil {
			// Check if already bookmarked
			isBookmarked, err := m.store.IsBookmarked(station.StationUUID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Error checking bookmark: %v", err)
				return m, nil
			}

			if isBookmarked {
				// Remove bookmark
				return m, func() tea.Msg {
					if err := m.store.RemoveBookmark(station.StationUUID); err != nil {
						return errMsg{err}
					}
					return bookmarkRemovedMsg{station.StationUUID}
				}
			} else {
				// Add bookmark
				return m, func() tea.Msg {
					if err := m.store.AddBookmark(station); err != nil {
						return errMsg{err}
					}
					return bookmarkAddedMsg{*station}
				}
			}
		}
		return m, nil

	// View switching
	case "b":
		m.view = ViewBookmarks
		m.bookmarksLoading = true
		// Load bookmarks when switching to bookmarks view
		return m, m.loadBookmarks

	case "f", "/":
		// 'f' for find (international keyboard friendly), '/' still works
		m.view = ViewSearch
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		m.searchResults = []radiobrowser.Station{}
		m.errorMsg = ""
		return m, textinput.Blink

	case "h", "?":
		// 'h' for help (international keyboard friendly), '?' still works
		m.view = ViewHelp
		return m, nil

	case "i":
		// Info/About view
		m.view = ViewAbout
		return m, nil
	}

	return m, nil
}

// handleSearchKeys handles keyboard input in the search view.
func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle text input first if focused (except for special keys)
	if m.searchInput.Focused() {
		switch msg.String() {
		case "esc":
			m.searchInput.Blur()
			m.view = ViewBrowse
			return m, nil
		case "enter":
			// Execute search
			query := m.searchInput.Value()
			if query != "" {
				m.searching = true
				m.errorMsg = ""
				return m, func() tea.Msg {
					return m.performSearch(query)
				}
			}
			return m, nil
		case "tab":
			// Switch focus to results
			m.searchInput.Blur()
			return m, nil
		default:
			// Pass all other keys to the text input
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}
	}

	// Handle navigation and commands when input is NOT focused
	switch msg.String() {
	case "esc":
		m.view = ViewBrowse
		return m, nil

	case "enter":
		// Play selected station from results
		if len(m.searchResults) > 0 {
			station := &m.searchResults[m.searchCursor]
			if err := m.player.Play(station); err != nil {
				m.errorMsg = fmt.Sprintf("Failed to play: %v", err)
			} else {
				m.errorMsg = ""
			}
		}
		return m, nil

	case "tab":
		// Switch focus back to input
		m.searchInput.Focus()
		cmd = textinput.Blink
		return m, cmd

	case "up", "k":
		if m.searchCursor > 0 {
			m.searchCursor--
			m.updateSearchScroll()
		}
		return m, nil

	case "down", "j":
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
			m.updateSearchScroll()
		}
		return m, nil

	case "s":
		// Stop playback
		_ = m.player.Stop()
		m.errorMsg = ""
		return m, nil

	case "=", "+":
		// Increase volume
		currentVol := m.player.GetVolume()
		if currentVol < 100 {
			_ = m.player.SetVolume(currentVol + 10)
		}
		return m, nil

	case "-", "_":
		// Decrease volume
		currentVol := m.player.GetVolume()
		if currentVol > 0 {
			_ = m.player.SetVolume(currentVol - 10)
		}
		return m, nil

	case "a":
		// Add/remove bookmark for selected result
		if len(m.searchResults) > 0 && m.store != nil {
			station := &m.searchResults[m.searchCursor]

			// Check if already bookmarked
			isBookmarked, err := m.store.IsBookmarked(station.StationUUID)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Error checking bookmark: %v", err)
				return m, nil
			}

			if isBookmarked {
				// Remove bookmark
				return m, func() tea.Msg {
					if err := m.store.RemoveBookmark(station.StationUUID); err != nil {
						return errMsg{err}
					}
					return bookmarkRemovedMsg{station.StationUUID}
				}
			} else {
				// Add bookmark
				return m, func() tea.Msg {
					if err := m.store.AddBookmark(station); err != nil {
						return errMsg{err}
					}
					return bookmarkAddedMsg{*station}
				}
			}
		}
		return m, nil
	}

	return m, cmd
}

// updateSearchScroll adjusts search scroll offset based on cursor position.
func (m *Model) updateSearchScroll() {
	visible := m.VisibleStations()

	// Scroll down if cursor is below visible area
	if m.searchCursor >= m.searchScrollOffset+visible {
		m.searchScrollOffset = m.searchCursor - visible + 1
	}

	// Scroll up if cursor is above visible area
	if m.searchCursor < m.searchScrollOffset {
		m.searchScrollOffset = m.searchCursor
	}
}

// handleBookmarksKeys handles keyboard input in the bookmarks view.
func (m Model) handleBookmarksKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.view = ViewBrowse
		return m, nil

	// Navigation
	case "up", "k":
		if m.bookmarksCursor > 0 {
			m.bookmarksCursor--
			m.updateBookmarksScroll()
		}
		return m, nil

	case "down", "j":
		if m.bookmarksCursor < len(m.bookmarks)-1 {
			m.bookmarksCursor++
			m.updateBookmarksScroll()
		}
		return m, nil

	case "pgup":
		visible := m.VisibleStations()
		m.bookmarksCursor -= visible
		if m.bookmarksCursor < 0 {
			m.bookmarksCursor = 0
		}
		m.updateBookmarksScroll()
		return m, nil

	case "pgdown":
		visible := m.VisibleStations()
		m.bookmarksCursor += visible
		if m.bookmarksCursor >= len(m.bookmarks) {
			m.bookmarksCursor = len(m.bookmarks) - 1
		}
		m.updateBookmarksScroll()
		return m, nil

	case "home", "g":
		m.bookmarksCursor = 0
		m.updateBookmarksScroll()
		return m, nil

	case "end", "G":
		if len(m.bookmarks) > 0 {
			m.bookmarksCursor = len(m.bookmarks) - 1
		}
		m.updateBookmarksScroll()
		return m, nil

	// Actions
	case "enter", " ":
		// Play selected bookmark
		if len(m.bookmarks) > 0 && m.bookmarksCursor < len(m.bookmarks) {
			station := &m.bookmarks[m.bookmarksCursor]
			currentStation := m.player.GetCurrentStation()
			if currentStation != nil && currentStation.StationUUID == station.StationUUID {
				// Stop if already playing this station
				_ = m.player.Stop()
			} else {
				// Play the selected station
				if err := m.player.Play(station); err != nil {
					m.errorMsg = fmt.Sprintf("Failed to play: %v", err)
				} else {
					m.errorMsg = ""
				}
			}
		}
		return m, nil

	case "s":
		// Stop playback
		_ = m.player.Stop()
		m.errorMsg = ""
		return m, nil

	case "=", "+":
		// Increase volume
		currentVol := m.player.GetVolume()
		if currentVol < 100 {
			_ = m.player.SetVolume(currentVol + 10)
		}
		return m, nil

	case "-", "_":
		// Decrease volume
		currentVol := m.player.GetVolume()
		if currentVol > 0 {
			_ = m.player.SetVolume(currentVol - 10)
		}
		return m, nil

	case "a", "d":
		// Remove bookmark (a for add/remove toggle, d for delete)
		if len(m.bookmarks) > 0 && m.bookmarksCursor < len(m.bookmarks) && m.store != nil {
			station := &m.bookmarks[m.bookmarksCursor]

			return m, func() tea.Msg {
				if err := m.store.RemoveBookmark(station.StationUUID); err != nil {
					return errMsg{err}
				}
				return bookmarkRemovedMsg{station.StationUUID}
			}
		}
		return m, nil
	}

	return m, nil
}

// updateBookmarksScroll adjusts bookmarks scroll offset based on cursor position.
func (m *Model) updateBookmarksScroll() {
	visible := m.VisibleStations()

	// Scroll down if cursor is below visible area
	if m.bookmarksCursor >= m.bookmarksScrollOffset+visible {
		m.bookmarksScrollOffset = m.bookmarksCursor - visible + 1
	}

	// Scroll up if cursor is above visible area
	if m.bookmarksCursor < m.bookmarksScrollOffset {
		m.bookmarksScrollOffset = m.bookmarksCursor
	}
}

// handleHelpKeys handles keyboard input in the help view.
func (m Model) handleHelpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "h", "?":
		m.view = ViewBrowse
		return m, nil
	case "i":
		m.view = ViewAbout
		return m, nil
	}

	return m, nil
}

// handleAboutKeys handles keyboard input in the about view.
func (m Model) handleAboutKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "i":
		m.view = ViewBrowse
		return m, nil
	case "h", "?":
		m.view = ViewHelp
		return m, nil
	}

	return m, nil
}
