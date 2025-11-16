# Radio Browser API Integration

This document describes the integration with the [Radio Browser API](https://www.radio-browser.info/), a community-driven database of internet radio stations.

## Table of Contents

- [Overview](#overview)
- [API Endpoints](#api-endpoints)
- [Request Examples](#request-examples)
- [Response Format](#response-format)
- [Caching Strategy](#caching-strategy)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Best Practices](#best-practices)

## Overview

### About Radio Browser

- **Database Size**: 30,000+ radio stations
- **Community-Driven**: Open database maintained by users
- **Free to Use**: No API key required
- **REST API**: Simple HTTP JSON API
- **Update Frequency**: Real-time user submissions

### Base URL

The API uses DNS-based load balancing. Always resolve `all.api.radio-browser.info` to get current servers:

```bash
# Get list of servers
dig all.api.radio-browser.info

# Example servers:
# - de1.api.radio-browser.info
# - nl1.api.radio-browser.info
# - at1.api.radio-browser.info
```

**Base URL**: `https://de1.api.radio-browser.info`

**Official Documentation**: https://api.radio-browser.info/

## API Endpoints

### 1. Search Stations

**Endpoint**: `GET /json/stations/search`

Search stations by various criteria.

**Query Parameters**:
- `name`: Station name (partial match)
- `country`: Country name
- `countrycode`: ISO 3166-1 alpha-2 country code (e.g., "IT", "US")
- `language`: Language name
- `tag`: Genre/tag (e.g., "jazz", "rock")
- `tagList`: Comma-separated tags (AND logic)
- `bitrate`: Minimum bitrate in kbps
- `order`: Sort order (`name`, `votes`, `clickcount`, `bitrate`)
- `reverse`: Reverse sort (true/false)
- `offset`: Pagination offset
- `limit`: Results per page (default: 100, max: 100000)
- `hidebroken`: Hide broken stations (true/false)

**Example**:
```
GET /json/stations/search?country=Italy&tag=jazz&limit=50&order=votes&reverse=true
```

### 2. Get Stations by Country

**Endpoint**: `GET /json/stations/bycountry/{country}`

Get all stations from a specific country.

**Parameters**:
- `{country}`: Country name (e.g., "Italy", "United States")

**Example**:
```
GET /json/stations/bycountry/Italy
```

### 3. Get Stations by Country Code

**Endpoint**: `GET /json/stations/bycountrycode/{code}`

**Parameters**:
- `{code}`: ISO 3166-1 alpha-2 code (e.g., "IT", "US", "UK")

**Example**:
```
GET /json/stations/bycountrycode/IT
```

### 4. Get Stations by Tag

**Endpoint**: `GET /json/stations/bytag/{tag}`

Get stations by genre/tag.

**Parameters**:
- `{tag}`: Tag name (e.g., "jazz", "rock", "classical")

**Example**:
```
GET /json/stations/bytag/jazz
```

### 5. Get Stations by Language

**Endpoint**: `GET /json/stations/bylanguage/{language}`

**Parameters**:
- `{language}`: Language name (e.g., "italian", "english")

**Example**:
```
GET /json/stations/bylanguage/italian
```

### 6. Get Station by UUID

**Endpoint**: `GET /json/stations/byuuid/{uuid}`

Get specific station details.

**Parameters**:
- `{uuid}`: Station UUID

**Example**:
```
GET /json/stations/byuuid/9608b51d-0601-11e8-ae97-52543be04c81
```

### 7. List Countries

**Endpoint**: `GET /json/countries`

Get list of all countries with station counts.

### 8. List Languages

**Endpoint**: `GET /json/languages`

Get list of all languages with station counts.

### 9. List Tags

**Endpoint**: `GET /json/tags`

Get list of all tags/genres with station counts.

### 10. Vote for Station

**Endpoint**: `GET /json/vote/{uuid}`

Increase station's vote count (quality indicator).

**Parameters**:
- `{uuid}`: Station UUID

## Request Examples

### Go HTTP Client

```go
package radiobrowser

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "time"
)

const (
    baseURL   = "https://de1.api.radio-browser.info"
    userAgent = "Terminal.FM/1.0"
)

type Client struct {
    httpClient *http.Client
    baseURL    string
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        baseURL: baseURL,
    }
}

// SearchStations searches for radio stations
func (c *Client) SearchStations(params SearchParams) ([]Station, error) {
    u, err := url.Parse(c.baseURL + "/json/stations/search")
    if err != nil {
        return nil, err
    }

    q := u.Query()
    if params.Name != "" {
        q.Set("name", params.Name)
    }
    if params.Country != "" {
        q.Set("country", params.Country)
    }
    if params.Tag != "" {
        q.Set("tag", params.Tag)
    }
    if params.Language != "" {
        q.Set("language", params.Language)
    }
    if params.Limit > 0 {
        q.Set("limit", fmt.Sprintf("%d", params.Limit))
    }
    q.Set("hidebroken", "true")
    u.RawQuery = q.Encode()

    req, err := http.NewRequest("GET", u.String(), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", userAgent)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }

    var stations []Station
    if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
        return nil, err
    }

    return stations, nil
}
```

### cURL Examples

```bash
# Search by name
curl "https://de1.api.radio-browser.info/json/stations/search?name=jazz&limit=10"

# Get Italian stations
curl "https://de1.api.radio-browser.info/json/stations/bycountrycode/IT?limit=20"

# Get jazz stations with high bitrate
curl "https://de1.api.radio-browser.info/json/stations/bytag/jazz?bitrate=128&order=votes&reverse=true"

# List all countries
curl "https://de1.api.radio-browser.info/json/countries"

# Get specific station
curl "https://de1.api.radio-browser.info/json/stations/byuuid/9608b51d-0601-11e8-ae97-52543be04c81"
```

## Response Format

### Station Object

```json
{
  "stationuuid": "9608b51d-0601-11e8-ae97-52543be04c81",
  "name": "Jazz Radio",
  "url": "https://jazzradio.ice.infomaniak.ch/jazzradio-high.mp3",
  "url_resolved": "https://jazzradio.ice.infomaniak.ch/jazzradio-high.mp3",
  "homepage": "https://www.jazzradio.com",
  "favicon": "https://www.jazzradio.com/favicon.ico",
  "tags": "jazz,smooth jazz,instrumental",
  "country": "Switzerland",
  "countrycode": "CH",
  "state": "",
  "language": "english",
  "languagecodes": "en",
  "votes": 1234,
  "lastchangetime": "2023-01-15 10:30:45",
  "lastchangetime_iso8601": "2023-01-15T10:30:45Z",
  "codec": "MP3",
  "bitrate": 128,
  "hls": 0,
  "lastcheckok": 1,
  "lastchecktime": "2023-11-16 08:15:30",
  "lastchecktime_iso8601": "2023-11-16T08:15:30Z",
  "lastcheckoktime": "2023-11-16 08:15:30",
  "lastcheckoktime_iso8601": "2023-11-16T08:15:30Z",
  "lastlocalchecktime": "2023-11-16 08:15:30",
  "clicktimestamp": "2023-11-16 10:22:15",
  "clicktimestamp_iso8601": "2023-11-16T10:22:15Z",
  "clickcount": 5678,
  "clicktrend": 10,
  "ssl_error": 0,
  "geo_lat": 46.5,
  "geo_long": 6.5
}
```

### Important Fields

- **stationuuid**: Unique identifier (use for bookmarks)
- **name**: Station name
- **url**: Stream URL (use this to play)
- **url_resolved**: Resolved stream URL (may differ if redirects)
- **tags**: Comma-separated genres/tags
- **country** / **countrycode**: Location
- **language** / **languagecodes**: Broadcast language
- **codec**: Audio codec (MP3, AAC, OGG, etc.)
- **bitrate**: Stream bitrate in kbps
- **votes**: Quality indicator (higher = better)
- **lastcheckok**: 1 if station is working, 0 if broken
- **clickcount**: Popularity indicator

### Go Struct

```go
type Station struct {
    StationUUID     string  `json:"stationuuid"`
    Name            string  `json:"name"`
    URL             string  `json:"url"`
    URLResolved     string  `json:"url_resolved"`
    Homepage        string  `json:"homepage"`
    Favicon         string  `json:"favicon"`
    Tags            string  `json:"tags"`
    Country         string  `json:"country"`
    CountryCode     string  `json:"countrycode"`
    State           string  `json:"state"`
    Language        string  `json:"language"`
    LanguageCodes   string  `json:"languagecodes"`
    Votes           int     `json:"votes"`
    Codec           string  `json:"codec"`
    Bitrate         int     `json:"bitrate"`
    HLS             int     `json:"hls"`
    LastCheckOK     int     `json:"lastcheckok"`
    LastCheckTime   string  `json:"lastchecktime_iso8601"`
    ClickCount      int     `json:"clickcount"`
    ClickTrend      int     `json:"clicktrend"`
    GeoLat          float64 `json:"geo_lat"`
    GeoLong         float64 `json:"geo_long"`
}

type SearchParams struct {
    Name       string
    Country    string
    CountryCode string
    Language   string
    Tag        string
    Bitrate    int
    Limit      int
    Offset     int
    Order      string
    Reverse    bool
}
```

## Caching Strategy

### Why Cache?

- Reduce API calls (be a good citizen)
- Improve response times
- Handle API downtime gracefully
- Reduce bandwidth usage

### Cache Implementation

```go
type Cache struct {
    data      map[string]CacheEntry
    mu        sync.RWMutex
    ttl       time.Duration
}

type CacheEntry struct {
    Value     interface{}
    ExpiresAt time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, exists := c.data[key]
    if !exists {
        return nil, false
    }

    if time.Now().After(entry.ExpiresAt) {
        return nil, false
    }

    return entry.Value, true
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = CacheEntry{
        Value:     value,
        ExpiresAt: time.Now().Add(c.ttl),
    }
}
```

### Cache TTL Recommendations

- **Station lists**: 1 hour
- **Station details**: 5 minutes
- **Country/Language/Tag lists**: 24 hours
- **Search results**: 10 minutes

## Error Handling

### Common Errors

1. **Network Errors**: Timeout, connection refused
2. **API Errors**: 4xx, 5xx status codes
3. **Parsing Errors**: Invalid JSON response
4. **Empty Results**: No stations found

### Error Handling Example

```go
func (c *Client) SearchStations(params SearchParams) ([]Station, error) {
    // Try primary server
    stations, err := c.searchStationsFromServer(c.baseURL, params)
    if err == nil {
        return stations, nil
    }

    // Fallback to alternative server
    log.Printf("Primary server failed: %v, trying fallback", err)
    alternativeURL := "https://nl1.api.radio-browser.info"
    stations, err = c.searchStationsFromServer(alternativeURL, params)
    if err != nil {
        return nil, fmt.Errorf("all servers failed: %w", err)
    }

    return stations, nil
}
```

## Rate Limiting

### API Limits

Radio Browser API **does not have strict rate limits**, but recommends:
- Reasonable request frequency
- Use caching when possible
- Set proper User-Agent header

### Implementation

```go
type RateLimiter struct {
    requests chan struct{}
    interval time.Duration
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    rl := &RateLimiter{
        requests: make(chan struct{}, requestsPerSecond),
        interval: time.Second / time.Duration(requestsPerSecond),
    }

    // Fill initial bucket
    for i := 0; i < requestsPerSecond; i++ {
        rl.requests <- struct{}{}
    }

    // Refill bucket
    go func() {
        ticker := time.NewTicker(rl.interval)
        defer ticker.Stop()
        for range ticker.C {
            select {
            case rl.requests <- struct{}{}:
            default:
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Wait() {
    <-rl.requests
}
```

## Best Practices

### 1. Always Use HTTPS

```go
baseURL := "https://de1.api.radio-browser.info" // ✓ Good
baseURL := "http://de1.api.radio-browser.info"  // ✗ Bad
```

### 2. Set User-Agent

```go
req.Header.Set("User-Agent", "Terminal.FM/1.0 (+https://github.com/fulgidus/terminal-fm)")
```

### 3. Handle Broken Stations

```go
// Filter working stations only
params.HideBroken = true

// Or check lastcheckok field
if station.LastCheckOK == 1 {
    // Station is working
}
```

### 4. Sort by Quality

```go
// Get best stations first
params.Order = "votes"
params.Reverse = true
```

### 5. Implement Retry Logic

```go
func (c *Client) searchWithRetry(params SearchParams) ([]Station, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        stations, err := c.SearchStations(params)
        if err == nil {
            return stations, nil
        }
        
        if i < maxRetries-1 {
            time.Sleep(time.Second * time.Duration(i+1))
        }
    }
    return nil, fmt.Errorf("max retries exceeded")
}
```

### 6. Vote for Good Stations

```go
// When user plays a station successfully, vote for it
func (c *Client) VoteForStation(uuid string) error {
    url := fmt.Sprintf("%s/json/vote/%s", c.baseURL, uuid)
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", userAgent)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

---

For architecture details, see [../ARCHITECTURE.md](../ARCHITECTURE.md).

For development setup, see [DEVELOPMENT.md](DEVELOPMENT.md).
