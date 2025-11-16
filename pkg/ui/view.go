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
	colorPrimary = lipgloss.Color("#00D9FF") // Cyan
	colorAccent  = lipgloss.Color("#C792EA") // Purple

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
	case ViewAbout:
		return m.viewAbout()
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
		"f find",
		"h help",
		"i about",
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
	b.WriteString(styleStationDetail.Render("Tip: Search by name, country (e.g., 'Italy', 'US'), or genre tag"))
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
	shortcuts := "enter search/play • tab switch • ↑/↓ nav • s stop • +/- vol • a bookmark • h help • i about • esc back"
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

		// Render bookmark list with scrolling
		visible := m.VisibleStations()
		end := m.bookmarksScrollOffset + visible
		if end > len(m.bookmarks) {
			end = len(m.bookmarks)
		}

		for i := m.bookmarksScrollOffset; i < end; i++ {
			station := m.bookmarks[i]
			isSelected := i == m.bookmarksCursor
			b.WriteString(m.renderStation(station, isSelected))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	shortcuts := "↑/↓ navigate • enter play • s stop • +/- vol • a/d remove • h help • i about • esc back"
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
		{"f", "Find/Search stations"},
		{"h", "Show this help"},
		{"i", "About Terminal.FM"},
		{"q / Ctrl+C", "Quit application"},
	}

	for _, h := range helps {
		key := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render(h.key)
		desc := lipgloss.NewStyle().Foreground(colorText).Render(h.desc)
		b.WriteString(fmt.Sprintf("  %s  %s\n", key, desc))
	}

	b.WriteString("\n")
	b.WriteString(styleFooter.Render("esc back • i about"))

	return b.String()
}

// viewAbout renders the about screen with credits and info.
func (m Model) viewAbout() string {
	var b strings.Builder

	// Title
	b.WriteString(styleTitle.Render("♫ About Terminal.FM"))
	b.WriteString("\n\n")

	// Version and description
	version := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Render("Version 1.0.0")
	b.WriteString("  " + version)
	b.WriteString("\n\n")

	description := lipgloss.NewStyle().
		Foreground(colorText).
		Render("Internet Radio Player for Your Terminal")
	b.WriteString("  " + description)
	b.WriteString("\n\n")

	// Separator
	separator := lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(strings.Repeat("─", 50))
	b.WriteString("  " + separator)
	b.WriteString("\n\n")

	// Features
	featuresTitle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Render("Features:")
	b.WriteString("  " + featuresTitle)
	b.WriteString("\n\n")

	features := []string{
		"• 30,000+ radio stations worldwide",
		"• Search by name, country, or genre",
		"• Bookmark your favorite stations",
		"• Real-time volume control",
		"• Multi-language support (EN/IT)",
		"• Clean TUI interface with Vim keybindings",
	}

	for _, feature := range features {
		featureText := lipgloss.NewStyle().Foreground(colorText).Render(feature)
		b.WriteString("  " + featureText + "\n")
	}
	b.WriteString("\n")

	// Credits
	creditsTitle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Render("Created by:")
	b.WriteString("  " + creditsTitle)
	b.WriteString("\n\n")

	author := lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true).
		Render("Fulgidus")
	b.WriteString("  " + author)
	b.WriteString("\n\n")

	// GitHub link
	githubLabel := lipgloss.NewStyle().
		Foreground(colorTextDim).
		Render("GitHub: ")
	githubLink := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Underline(true).
		Render("https://github.com/fulgidus/terminal-fm")
	b.WriteString("  " + githubLabel + githubLink)
	b.WriteString("\n\n")

	// Tech stack
	b.WriteString("  " + separator)
	b.WriteString("\n\n")

	techTitle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Render("Built with:")
	b.WriteString("  " + techTitle)
	b.WriteString("\n\n")

	tech := []string{
		"• Go 1.21+ - Programming language",
		"• Charm Wish - SSH server framework",
		"• Bubbletea - Terminal UI framework",
		"• FFplay - Audio streaming",
		"• SQLite - Local storage",
		"• Radio Browser API - Station database",
	}

	for _, t := range tech {
		techText := lipgloss.NewStyle().Foreground(colorTextDim).Render(t)
		b.WriteString("  " + techText + "\n")
	}
	b.WriteString("\n")

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(colorSuccess).
		Italic(true).
		Render("Made with ♥ for the open source community")
	b.WriteString("  " + footer)
	b.WriteString("\n\n")

	b.WriteString(styleFooter.Render("esc back • h help"))

	return b.String()
}
