package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fulgidus/terminal-fm/pkg/services/player"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
)

// Color scheme using Lipgloss.
var (
	// Brand colors
	colorPrimary   = lipgloss.Color("#00D9FF") // Cyan
	colorSecondary = lipgloss.Color("#FF6B9D") // Pink
	colorAccent    = lipgloss.Color("#C792EA") // Purple

	// UI colors
	colorText     = lipgloss.Color("#E0E0E0") // Light gray
	colorTextDim  = lipgloss.Color("#808080") // Dim gray
	colorBg       = lipgloss.Color("#1E1E1E") // Dark gray
	colorSelected = lipgloss.Color("#2C2C2C") // Slightly lighter gray
	colorBorder   = lipgloss.Color("#404040") // Border gray

	// Status colors
	colorSuccess = lipgloss.Color("#50FA7B") // Green
	colorError   = lipgloss.Color("#FF5555") // Red
	colorWarning = lipgloss.Color("#FFB86C") // Orange
)

// Styles
var (
	// Title style (top banner)
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Background(colorBg).
			Padding(0, 1)

	// Header info (station count, etc.)
	styleHeader = lipgloss.NewStyle().
			Foreground(colorTextDim).
			Padding(0, 1)

	// Station item (normal)
	styleStation = lipgloss.NewStyle().
			Foreground(colorText).
			Padding(0, 2)

	// Station item (selected)
	styleStationSelected = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(colorSelected).
				Bold(true).
				Padding(0, 2)

	// Station details (country, bitrate, etc.)
	styleStationDetail = lipgloss.NewStyle().
				Foreground(colorTextDim)

	// Footer with shortcuts
	styleFooter = lipgloss.NewStyle().
			Foreground(colorTextDim).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	// Error message
	styleError = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true).
			Padding(0, 1)

	// Loading message
	styleLoading = lipgloss.NewStyle().
			Foreground(colorAccent).
			Padding(0, 1)

	// Status bar styles
	styleStatusBar = lipgloss.NewStyle().
			BorderBottom(true).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1).
			MarginBottom(1)

	styleStatusPlaying = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true)

	styleStatusStopped = lipgloss.NewStyle().
				Foreground(colorTextDim)

	styleStatusBuffering = lipgloss.NewStyle().
				Foreground(colorWarning).
				Bold(true)

	styleStatusError = lipgloss.NewStyle().
				Foreground(colorError).
				Bold(true)
)

// View renders the entire UI (required by Bubbletea).
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	switch m.view {
	case ViewBrowse:
		return m.viewBrowse()
	case ViewSearch:
		return m.viewSearch()
	case ViewBookmarks:
		return m.viewBookmarks()
	case ViewHelp:
		return m.viewHelp()
	default:
		return "Unknown view"
	}
}

// viewBrowse renders the station browsing view.
func (m Model) viewBrowse() string {
	var b strings.Builder

	// Title
	title := styleTitle.Render("♫ " + m.tr.T("app.title"))
	b.WriteString(title)
	b.WriteString("\n")

	// Status bar with player state
	b.WriteString(m.renderStatusBar())
	b.WriteString("\n")

	// Header info
	if m.loading {
		header := styleLoading.Render(m.tr.T("station.loading"))
		b.WriteString(header)
		b.WriteString("\n")
	} else {
		header := styleHeader.Render(m.tr.Tf("station.found", len(m.stations)))
		b.WriteString(header)
		b.WriteString("\n")
	}

	// Error message if any
	if m.errorMsg != "" {
		b.WriteString(styleError.Render(m.errorMsg))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Station list
	if !m.loading && len(m.stations) > 0 {
		b.WriteString(m.renderStationList())
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderStationList renders the scrollable station list.
func (m Model) renderStationList() string {
	var b strings.Builder

	visibleStations := m.VisibleStationList()

	for i, station := range visibleStations {
		absoluteIndex := m.scrollOffset + i
		isSelected := absoluteIndex == m.cursor

		b.WriteString(m.renderStation(station, isSelected))
		b.WriteString("\n")
	}

	return b.String()
}

// renderStation renders a single station item.
func (m Model) renderStation(station radiobrowser.Station, selected bool) string {
	// Format: "► Station Name - Country | Bitrate kbps"
	name := station.Name
	if len(name) > 40 {
		name = name[:37] + "..."
	}

	details := fmt.Sprintf("%s | %d kbps", station.Country, station.Bitrate)

	cursor := " "
	if selected {
		cursor = "►"
	}

	line := fmt.Sprintf("%s %s", cursor, name)
	detailsPart := styleStationDetail.Render(details)

	if selected {
		return styleStationSelected.Render(line) + " " + detailsPart
	}
	return styleStation.Render(line) + " " + detailsPart
}

// renderFooter renders the keyboard shortcuts footer.
func (m Model) renderFooter() string {
	shortcuts := []string{
		"↑/k up",
		"↓/j down",
		"enter play",
		"s stop",
		"+/- vol",
		"b bookmarks",
		"/ search",
		"? help",
		"q quit",
	}

	footer := strings.Join(shortcuts, " • ")
	return styleFooter.Width(m.width).Render(footer)
}

// viewSearch renders the search interface.
func (m Model) viewSearch() string {
	var b strings.Builder

	b.WriteString(styleTitle.Render("♫ Search Stations"))
	b.WriteString("\n")
	
	// Status bar
	b.WriteString(m.renderStatusBar())
	b.WriteString("\n\n")

	// Search input box
	inputLabel := styleHeader.Render("Enter search query:")
	b.WriteString(inputLabel)
	b.WriteString("\n")
	b.WriteString(m.searchInput.View())
	b.WriteString("\n")
	b.WriteString(styleStationDetail.Render("Tip: Enter 2 letters for country code (e.g., 'IT', 'US'), or station name"))
	b.WriteString("\n\n")

	// Show searching status
	if m.searching {
		b.WriteString(styleLoading.Render("Searching..."))
		b.WriteString("\n")
	} else if len(m.searchResults) > 0 {
		// Show results count
		header := styleHeader.Render(fmt.Sprintf("Found %d stations (Tab to navigate results)", len(m.searchResults)))
		b.WriteString(header)
		b.WriteString("\n\n")

		// Render results list
		visible := m.VisibleStations() - 8 // Reserve space for input area
		if visible < 1 {
			visible = 1
		}

		end := m.searchScrollOffset + visible
		if end > len(m.searchResults) {
			end = len(m.searchResults)
		}

		for i := m.searchScrollOffset; i < end; i++ {
			station := m.searchResults[i]
			isSelected := i == m.searchCursor && !m.searchInput.Focused()
			b.WriteString(m.renderStation(station, isSelected))
			b.WriteString("\n")
		}
	} else if m.searchInput.Value() != "" && !m.searching {
		b.WriteString(styleHeader.Render("No results found"))
		b.WriteString("\n")
	}

	// Error message if any
	if m.errorMsg != "" {
		b.WriteString("\n")
		b.WriteString(styleError.Render(m.errorMsg))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	shortcuts := "enter search/play • tab switch focus • ↑/↓ navigate • esc back"
	b.WriteString(styleFooter.Width(m.width).Render(shortcuts))

	return b.String()
}

// viewBookmarks renders the bookmarks view.
func (m Model) viewBookmarks() string {
	var b strings.Builder

	b.WriteString(styleTitle.Render("♫ Bookmarks"))
	b.WriteString("\n")
	
	// Status bar
	b.WriteString(m.renderStatusBar())
	b.WriteString("\n\n")

	if m.bookmarksLoading {
		b.WriteString(styleLoading.Render("Loading bookmarks..."))
		b.WriteString("\n")
	} else if len(m.bookmarks) == 0 {
		b.WriteString(styleHeader.Render("No bookmarks yet"))
		b.WriteString("\n")
		b.WriteString(styleStationDetail.Render("Press 'a' on any station to bookmark it"))
		b.WriteString("\n")
	} else {
		b.WriteString(styleHeader.Render(fmt.Sprintf("%d bookmarked stations", len(m.bookmarks))))
		b.WriteString("\n\n")

		// Render bookmark list (similar to station list)
		for i, station := range m.bookmarks {
			if i >= m.VisibleStations() {
				break
			}
			b.WriteString(m.renderStation(station, false))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	shortcuts := "enter play • a remove • esc back"
	b.WriteString(styleFooter.Width(m.width).Render(shortcuts))

	return b.String()
}

// renderStatusBar renders the player status indicator.
func (m Model) renderStatusBar() string {
	var statusText, statusIcon string
	var statusStyle lipgloss.Style

	currentStation := m.player.GetCurrentStation()
	playerState := m.player.GetState()

	if currentStation != nil && playerState == player.StatePlaying {
		statusIcon = "[STREAMING]"
		statusStyle = styleStatusPlaying
		volume := m.tr.Tf("station.volume", m.player.GetVolume())
		statusText = fmt.Sprintf("%s %s - %s", statusIcon, currentStation.Name, volume)
	} else if m.loading {
		statusIcon = "[BUFFERING]"
		statusStyle = styleStatusBuffering
		statusText = fmt.Sprintf("%s %s", statusIcon, m.tr.T("station.loading"))
	} else if m.errorMsg != "" && currentStation == nil {
		statusIcon = "[ERROR]"
		statusStyle = styleStatusError
		statusText = fmt.Sprintf("%s %s", statusIcon, m.errorMsg)
	} else {
		statusIcon = "[STOPPED]"
		statusStyle = styleStatusStopped
		statusText = fmt.Sprintf("%s %s", statusIcon, m.tr.T("station.stopped"))
	}

	return styleStatusBar.Width(m.width - 2).Render(statusStyle.Render(statusText))
}

// viewHelp renders the help screen with all keyboard shortcuts.
func (m Model) viewHelp() string {
	var b strings.Builder

	b.WriteString(styleTitle.Render("♫ Keyboard Shortcuts"))
	b.WriteString("\n\n")

	helps := []struct {
		key  string
		desc string
	}{
		{"↑ / k", "Move cursor up"},
		{"↓ / j", "Move cursor down"},
		{"Enter / Space", "Play selected station"},
		{"s", "Stop playback"},
		{"+ / -", "Volume up/down"},
		{"a", "Add/Remove bookmark"},
		{"b", "Toggle bookmarks view"},
		{"/", "Search stations"},
		{"?", "Show this help"},
		{"q / Ctrl+C", "Quit application"},
	}

	for _, h := range helps {
		key := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render(h.key)
		desc := lipgloss.NewStyle().Foreground(colorText).Render(h.desc)
		b.WriteString(fmt.Sprintf("  %s  %s\n", key, desc))
	}

	b.WriteString("\n")
	b.WriteString(styleFooter.Render("Press ESC to go back"))

	return b.String()
}
