// Package i18n provides internationalization support.
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Translator handles text translations.
type Translator struct {
	locale       string
	translations map[string]string
	fallback     map[string]string
}

// NewTranslator creates a new translator for the given locale.
func NewTranslator(locale, localesPath string) (*Translator, error) {
	t := &Translator{
		locale:       locale,
		translations: make(map[string]string),
		fallback:     make(map[string]string),
	}

	// Load fallback (English)
	if err := t.loadLocale("en", localesPath); err != nil {
		return nil, fmt.Errorf("failed to load fallback locale: %w", err)
	}
	t.fallback = t.translations

	// Load requested locale if different
	if locale != "en" {
		t.translations = make(map[string]string)
		if err := t.loadLocale(locale, localesPath); err != nil {
			// If locale not found, use fallback
			t.translations = t.fallback
		}
	}

	return t, nil
}

// loadLocale loads translations from a JSON file.
func (t *Translator) loadLocale(locale, localesPath string) error {
	filePath := filepath.Join(localesPath, locale+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read locale file: %w", err)
	}

	if err := json.Unmarshal(data, &t.translations); err != nil {
		return fmt.Errorf("failed to parse locale file: %w", err)
	}

	return nil
}

// T translates a key to the current locale.
func (t *Translator) T(key string) string {
	// Try current locale
	if val, ok := t.translations[key]; ok {
		return val
	}

	// Try fallback
	if val, ok := t.fallback[key]; ok {
		return val
	}

	// Return key if not found
	return key
}

// Tf translates a key with format arguments.
func (t *Translator) Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(t.T(key), args...)
}

// GetLocale returns the current locale.
func (t *Translator) GetLocale() string {
	return t.locale
}
