package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fulgidus/terminal-fm/pkg/services/radiobrowser"
)

// Color scheme using Lipgloss.
var (
	// Brand colors
	colorPrimary   = lipgloss.Color("#00D9FF") // Cyan
	colorSecondary = lipgloss.Color("#FF6B9D") // Pink
	colorAccent    = lipgloss.Color("#C792EA") // Purple
	
	// UI colors
	colorText      = lipgloss.Color("#E0E0E0") // Light gray
	colorTextDim   = lipgloss.Color("#808080") // Dim gray
	colorBg        = lipgloss.Color("#1E1E1E") // Dark gray
	colorSelected  = lipgloss.Color("#2C2C2C") // Slightly lighter gray
	colorBorder    = lipgloss.Color("#404040") // Border gray
	
	// Status colors
	colorSuccess   = lipgloss.Color("#50FA7B") // Green
	colorError     = lipgloss.Color("#FF5555") // Red
	colorWarning   = lipgloss.Color("#FFB86C") // Orange
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
	title := styleTitle.Render("♫ Terminal.FM")
	b.WriteString(title)
	b.WriteString("\n\n")
	
	// Header info
	if m.loading {
		header := styleLoading.Render("Loading stations...")
		b.WriteString(header)
		b.WriteString("\n")
	} else {
		header := styleHeader.Render(fmt.Sprintf("Found %d stations", len(m.stations)))
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
		"b bookmarks",
		"/ search",
		"? help",
		"q quit",
	}
	
	footer := strings.Join(shortcuts, " • ")
	return styleFooter.Width(m.width).Render(footer)
}

// viewSearch renders the search interface (placeholder for now).
func (m Model) viewSearch() string {
	return styleTitle.Render("Search - Coming Soon") + "\n\n" +
		"Press ESC to go back\n" +
		m.renderFooter()
}

// viewBookmarks renders the bookmarks view (placeholder for now).
func (m Model) viewBookmarks() string {
	var b strings.Builder
	
	b.WriteString(styleTitle.Render("♫ Bookmarks"))
	b.WriteString("\n\n")
	
	if len(m.bookmarks) == 0 {
		b.WriteString(styleHeader.Render("No bookmarks yet"))
		b.WriteString("\n")
	} else {
		b.WriteString(styleHeader.Render(fmt.Sprintf("%d bookmarked stations", len(m.bookmarks))))
		b.WriteString("\n\n")
		// TODO: Render bookmark list
	}
	
	b.WriteString("\n")
	b.WriteString(styleFooter.Render("Press ESC to go back"))
	
	return b.String()
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
		{"Enter", "Play selected station"},
		{"Space", "Pause/Resume playback"},
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
