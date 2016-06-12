package radiowatch

import (
	"time"
	"regexp"
	"strings"
)

/*
Contains information to a track
 */
type TrackInfo struct {
	// The title of the track
	Title      string `json:"title"`
	// The name of the artist
	Artist     string `json:"artist"`
	// The time at which this information was crawled
	CrawlTime  time.Time `json:"crawl_time"`
	// The name of the station at which this track was played
	Station    string `json:"station"`
	normalized string
}

/*
Returns the stations name in a normalized way so it can be used as a file name.
Converts each non-word character (\W) to a lodash (_).
 */
func (t *TrackInfo) NormalizedStationName() string {
	if len(t.normalized) > 0 {
		return t.normalized
	}

	regex, _ := regexp.Compile("\\W+")
	return strings.Trim(regex.ReplaceAllString(strings.ToLower(t.Station), "_"), "_")
}
