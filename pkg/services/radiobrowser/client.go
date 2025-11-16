package radiobrowser

import (
	"fmt"
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
		filtered := []Station{}
		for _, station := range mockStations {
			// Case-insensitive search would go here
			filtered = append(filtered, station)
		}
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
