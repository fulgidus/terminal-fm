// Package i18n provides internationalization support.
package i18n

import "fmt"

// GetEmbeddedTranslations returns embedded translations for a locale.
func GetEmbeddedTranslations(locale string) map[string]string {
	translations := make(map[string]string)

	// English translations (fallback)
	en := map[string]string{
		"app.title":          "Terminal.FM",
		"station.loading":    "Loading stations...",
		"station.found":      "Found %d stations",
		"station.playing":    "♪ Now Playing: %s",
		"station.volume":     "Vol: %d%%",
		"station.no_results": "No stations found",
		"bookmark.empty":     "No bookmarks yet",
		"bookmark.hint":      "Press 'a' on any station to bookmark it",
		"bookmark.count":     "%d bookmarked stations",
		"bookmark.added":     "Added '%s' to bookmarks",
		"bookmark.removed":   "Removed from bookmarks",
		"bookmark.loading":   "Loading bookmarks...",
		"search.title":       "Search Stations",
		"search.prompt":      "Enter search query:",
		"search.placeholder": "Search stations by name, country, or tag...",
		"search.hint":        "Tip: Enter 2 letters for country code (e.g., 'IT', 'US'), or station name",
		"search.searching":   "Searching...",
		"search.results":     "Found %d stations (Tab to navigate results)",
		"error.play_failed":  "Failed to play: %v",
		"key.up":             "↑/k up",
		"key.down":           "↓/j down",
		"key.play":           "enter play",
		"key.stop":           "s stop",
		"key.volume":         "+/- vol",
		"key.bookmark":       "a bookmark",
		"key.bookmarks":      "b bookmarks",
		"key.search":         "/ search",
		"key.help":           "? help",
		"key.quit":           "q quit",
		"help.title":         "Keyboard Shortcuts",
		"help.up":            "Move cursor up",
		"help.down":          "Move cursor down",
		"help.play":          "Play selected station",
		"help.stop":          "Stop playback",
		"help.volume":        "Volume up/down",
		"help.bookmark":      "Add/Remove bookmark",
		"help.bookmarks":     "Toggle bookmarks view",
		"help.search":        "Search stations",
		"help.help":          "Show this help",
		"help.quit":          "Quit application",
		"help.back":          "Press ESC to go back",
		"footer.back":        "enter play • a remove • esc back",
	}

	// Italian translations
	it := map[string]string{
		"app.title":          "Terminal.FM",
		"station.loading":    "Caricamento stazioni...",
		"station.found":      "Trovate %d stazioni",
		"station.playing":    "♪ In Riproduzione: %s",
		"station.volume":     "Vol: %d%%",
		"station.no_results": "Nessuna stazione trovata",
		"bookmark.empty":     "Nessun preferito",
		"bookmark.hint":      "Premi 'a' su una stazione per aggiungerla ai preferiti",
		"bookmark.count":     "%d stazioni nei preferiti",
		"bookmark.added":     "Aggiunta '%s' ai preferiti",
		"bookmark.removed":   "Rimossa dai preferiti",
		"bookmark.loading":   "Caricamento preferiti...",
		"search.title":       "Cerca Stazioni",
		"search.prompt":      "Inserisci query di ricerca:",
		"search.placeholder": "Cerca stazioni per nome, paese o tag...",
		"search.hint":        "Suggerimento: Inserisci 2 lettere per codice paese (es. 'IT', 'US'), o nome stazione",
		"search.searching":   "Ricerca in corso...",
		"search.results":     "Trovate %d stazioni (Tab per navigare i risultati)",
		"error.play_failed":  "Riproduzione fallita: %v",
		"key.up":             "↑/k su",
		"key.down":           "↓/j giù",
		"key.play":           "invio riproduci",
		"key.stop":           "s ferma",
		"key.volume":         "+/- vol",
		"key.bookmark":       "a preferito",
		"key.bookmarks":      "b preferiti",
		"key.search":         "/ cerca",
		"key.help":           "? aiuto",
		"key.quit":           "q esci",
		"help.title":         "Scorciatoie da Tastiera",
		"help.up":            "Sposta cursore su",
		"help.down":          "Sposta cursore giù",
		"help.play":          "Riproduci stazione selezionata",
		"help.stop":          "Ferma riproduzione",
		"help.volume":        "Volume su/giù",
		"help.bookmark":      "Aggiungi/Rimuovi preferito",
		"help.bookmarks":     "Mostra preferiti",
		"help.search":        "Cerca stazioni",
		"help.help":          "Mostra questo aiuto",
		"help.quit":          "Esci dall'applicazione",
		"help.back":          "Premi ESC per tornare indietro",
		"footer.back":        "invio riproduci • a rimuovi • esc indietro",
	}

	switch locale {
	case "it":
		translations = it
	default:
		translations = en
	}

	return translations
}

// SimpleTranslator provides simple translation without file loading.
type SimpleTranslator struct {
	locale       string
	translations map[string]string
	fallback     map[string]string
}

// NewSimpleTranslator creates a new translator using embedded translations.
func NewSimpleTranslator(locale string) *SimpleTranslator {
	return &SimpleTranslator{
		locale:       locale,
		translations: GetEmbeddedTranslations(locale),
		fallback:     GetEmbeddedTranslations("en"),
	}
}

// T translates a key.
func (t *SimpleTranslator) T(key string) string {
	if val, ok := t.translations[key]; ok {
		return val
	}
	if val, ok := t.fallback[key]; ok {
		return val
	}
	return key
}

// Tf translates with format args.
func (t *SimpleTranslator) Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(t.T(key), args...)
}
