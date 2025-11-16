package radiobrowser

// Station represents a radio station from Radio Browser API
type Station struct {
	StationUUID   string  `json:"stationuuid"`
	Name          string  `json:"name"`
	URL           string  `json:"url"`
	URLResolved   string  `json:"url_resolved"`
	Homepage      string  `json:"homepage"`
	Favicon       string  `json:"favicon"`
	Tags          string  `json:"tags"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countrycode"`
	State         string  `json:"state"`
	Language      string  `json:"language"`
	LanguageCodes string  `json:"languagecodes"`
	Votes         int     `json:"votes"`
	Codec         string  `json:"codec"`
	Bitrate       int     `json:"bitrate"`
	HLS           int     `json:"hls"`
	LastCheckOK   int     `json:"lastcheckok"`
	ClickCount    int     `json:"clickcount"`
	ClickTrend    int     `json:"clicktrend"`
	GeoLat        float64 `json:"geo_lat"`
	GeoLong       float64 `json:"geo_long"`
}

// SearchParams holds parameters for searching stations
type SearchParams struct {
	Name        string
	Country     string
	CountryCode string
	Language    string
	Tag         string
	Bitrate     int
	Limit       int
	Offset      int
	Order       string
	Reverse     bool
}
