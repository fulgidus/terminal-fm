package radiobrowser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client interface for Radio Browser API
type Client interface {
	Search(params SearchParams) ([]Station, error)
	GetStationByUUID(uuid string) (*Station, error)
}

// MockClient provides mock data for development
type MockClient struct{}

// NewMockClient creates a new mock client with sample stations
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Search returns mock stations for testing
func (c *MockClient) Search(params SearchParams) ([]Station, error) {
	// Return mock stations for development
	mockStations := []Station{
		{
			StationUUID:   "960b51d-0601-11e8-ae97-52543be04c81",
			Name:          "Jazz Radio",
			URL:           "https://jazzradio.ice.infomaniak.ch/jazzradio-high.mp3",
			URLResolved:   "https://jazzradio.ice.infomaniak.ch/jazzradio-high.mp3",
			Homepage:      "https://www.jazzradio.com",
			Tags:          "jazz,smooth jazz,instrumental",
			Country:       "Switzerland",
			CountryCode:   "CH",
			Language:      "english",
			LanguageCodes: "en",
			Votes:         1234,
			Codec:         "MP3",
			Bitrate:       128,
			LastCheckOK:   1,
			ClickCount:    5678,
		},
		{
			StationUUID:   "a1b2c3d4-1234-5678-9abc-def012345678",
			Name:          "Classic Rock FM",
			URL:           "https://rockfm.example.com/stream",
			URLResolved:   "https://rockfm.example.com/stream",
			Homepage:      "https://www.rockfm.example.com",
			Tags:          "rock,classic rock,70s,80s",
			Country:       "United States",
			CountryCode:   "US",
			Language:      "english",
			LanguageCodes: "en",
			Votes:         2456,
			Codec:         "MP3",
			Bitrate:       192,
			LastCheckOK:   1,
			ClickCount:    8901,
		},
		{
			StationUUID:   "b2c3d4e5-2345-6789-abcd-ef0123456789",
			Name:          "Radio Italia",
			URL:           "https://radioitalia.example.com/live",
			URLResolved:   "https://radioitalia.example.com/live",
			Homepage:      "https://www.radioitalia.it",
			Tags:          "pop,italian,hits",
			Country:       "Italy",
			CountryCode:   "IT",
			Language:      "italian",
			LanguageCodes: "it",
			Votes:         1890,
			Codec:         "AAC",
			Bitrate:       128,
			LastCheckOK:   1,
			ClickCount:    4567,
		},
		{
			StationUUID:   "c3d4e5f6-3456-789a-bcde-f01234567890",
			Name:          "Electronic Beats",
			URL:           "https://electronicbeats.example.com/stream",
			URLResolved:   "https://electronicbeats.example.com/stream",
			Homepage:      "https://www.electronicbeats.de",
			Tags:          "electronic,techno,house,edm",
			Country:       "Germany",
			CountryCode:   "DE",
			Language:      "german",
			LanguageCodes: "de",
			Votes:         3210,
			Codec:         "MP3",
			Bitrate:       320,
			LastCheckOK:   1,
			ClickCount:    9876,
		},
		{
			StationUUID:   "d4e5f6a7-4567-89ab-cdef-012345678901",
			Name:          "Classical Vienna",
			URL:           "https://classicalvienna.example.com/live",
			URLResolved:   "https://classicalvienna.example.com/live",
			Homepage:      "https://www.classicalvienna.at",
			Tags:          "classical,orchestra,mozart,beethoven",
			Country:       "Austria",
			CountryCode:   "AT",
			Language:      "german",
			LanguageCodes: "de",
			Votes:         1567,
			Codec:         "MP3",
			Bitrate:       128,
			LastCheckOK:   1,
			ClickCount:    3456,
		},
	}

	// Simple filtering by name if provided
	if params.Name != "" {
		// Case-insensitive search would go here - for now return all
		filtered := append([]Station{}, mockStations...)
		return filtered, nil
	}

	return mockStations, nil
}

// GetStationByUUID returns a mock station by UUID
func (c *MockClient) GetStationByUUID(uuid string) (*Station, error) {
	stations, _ := c.Search(SearchParams{})
	for _, station := range stations {
		if station.StationUUID == uuid {
			return &station, nil
		}
	}
	return nil, fmt.Errorf("station not found: %s", uuid)
}

// APIClient implements the Client interface using the real Radio Browser API.
type APIClient struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// NewAPIClient creates a new Radio Browser API client.
// It automatically resolves the best API server to use.
func NewAPIClient() (*APIClient, error) {
	client := &APIClient{
		baseURL: "https://de1.api.radio-browser.info", // Default server
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "Terminal.FM/1.0",
	}

	// Try to resolve the best server
	if err := client.resolveBestServer(); err != nil {
		// If resolution fails, continue with default server
		fmt.Printf("Warning: Could not resolve best server, using default: %v\n", err)
	}

	return client, nil
}

// resolveBestServer finds the best Radio Browser API server using DNS.
func (c *APIClient) resolveBestServer() error {
	// Radio Browser uses DNS to distribute load across servers
	// Query all.api.radio-browser.info to get a random server
	resp, err := c.httpClient.Get("https://all.api.radio-browser.info/json/servers")
	if err != nil {
		return fmt.Errorf("failed to resolve servers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server resolution returned status %d", resp.StatusCode)
	}

	var servers []struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return fmt.Errorf("failed to decode servers: %w", err)
	}

	if len(servers) > 0 {
		// Use the first server (they're already randomized by DNS)
		c.baseURL = "https://" + servers[0].Name
	}

	return nil
}

// Search searches for radio stations using the provided parameters.
func (c *APIClient) Search(params SearchParams) ([]Station, error) {
	// Build the search endpoint based on parameters
	endpoint := "/json/stations/search"

	// Build query parameters
	query := url.Values{}

	if params.Name != "" {
		query.Set("name", params.Name)
	}
	if params.Country != "" {
		query.Set("country", params.Country)
	}
	if params.Language != "" {
		query.Set("language", params.Language)
	}
	if params.Tag != "" {
		query.Set("tag", params.Tag)
	}
	if params.Order != "" {
		query.Set("order", params.Order)
		query.Set("reverse", "true") // Most voted/clicked first
	}
	if params.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", params.Limit))
	} else {
		query.Set("limit", "50") // Default limit
	}
	if params.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", params.Offset))
	}

	// Construct full URL
	fullURL := c.baseURL + endpoint
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	// Create request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode response
	var stations []Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Filter out stations with empty URLs or that failed last check
	filtered := make([]Station, 0, len(stations))
	for _, station := range stations {
		if station.URLResolved != "" && station.LastCheckOK == 1 {
			// Clean up tags
			station.Tags = strings.TrimSpace(station.Tags)
			filtered = append(filtered, station)
		}
	}

	return filtered, nil
}

// GetStationByUUID retrieves a specific station by its UUID.
func (c *APIClient) GetStationByUUID(uuid string) (*Station, error) {
	endpoint := fmt.Sprintf("/json/stations/byuuid/%s", uuid)
	fullURL := c.baseURL + endpoint

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var stations []Station
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(stations) == 0 {
		return nil, fmt.Errorf("station not found: %s", uuid)
	}

	return &stations[0], nil
}
