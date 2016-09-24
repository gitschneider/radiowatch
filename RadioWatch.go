/*
Package radiowatch provides a small framework to crawl radio stations periodically.
 */
package radiowatch

import (
	"time"
	"sync"

	"strings"
	log "github.com/Sirupsen/logrus"
)

type(
	/*
	A Watcher is the main object.
	It keeps track of the crawlers, starts crawling,
	takes the results and delegates them to the Writer which will persist them.
	 */
	Watcher struct {
		refreshInterval time.Duration
		crawlers        []Crawler
		ticker          *time.Ticker
		writer          Writer
	}

	/*
	A concrete Crawler implementation crawls one specific radio station and returns
	information about the currently played track.
	 */
	Crawler interface {
		/*
		Takes the needed actions to get information about the currently played track.
		Returns the information in a TrackInfo struct or an error.
		 */
		Crawl() (*TrackInfo, error)
		/*
		Returns the name of the radio station it is crawling.
		 */
		Name() string
		/*
		Returns the time at which this crawler should be run next
		 */
		NextCrawlTime() time.Time
	}

	// A Writer takes the TrackInfo and persists it
	Writer interface {
		// Persists the TrackInfo to a concrete medium.
		// How and where it is saved is up to the concrete implementer.
		// Any error is silently dropped. The implementor can write a message to stderr.
		Write(TrackInfo)
	}
)

/*
Returns a new instance of Watcher.
Takes a concrete implementation of Writer which handles the persisting of results
 */
func NewWatcher(resultsWriter Writer) *Watcher {
	w := &Watcher{writer: resultsWriter}
	w.SetInterval("60s")
	return w
}

/*
The watcher will check the crawlers every {interval} seconds if its crawling time
is in the past and should therefore be crawled now.
Uses time.ParseDuration and takes therefore the same input like "60s" or "5m".
Returns only an error if the input was invalid.
 */
func (w *Watcher) SetInterval(interval string) error {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return err
	}

	w.refreshInterval = duration
	return nil
}

/*
Add a concrete Crawler to this watcher
 */
func (w *Watcher) AddCrawler(c Crawler) {
	w.crawlers = append(w.crawlers, c)
}

/*
Add several Crawlers at once
 */
func (w *Watcher) AddCrawlers(crawlers []Crawler) {
	for _, e := range crawlers {
		w.AddCrawler(e)
	}
}

/*
Run all Crawlers.
 */
func (w *Watcher) runCrawlers() {
	start := time.Now()
	tracks := make(chan *TrackInfo)
	var wg sync.WaitGroup
	var counter uint16
	go func() {
		for _, c := range w.crawlers {
			if c.NextCrawlTime().Before(start) {
				wg.Add(1)
				go func(crawler Crawler) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							log.WithFields(log.Fields{
								"crawler": crawler.Name(),
								"message": r,
							}).Error("Crawler panicked")

						}
					}()

					counter++
					track, err := crawler.Crawl()
					if err != nil {
						log.WithFields(log.Fields{
							"error": err.Error(),
							"crawler": crawler.Name(),
						}).Error("Error while crawling")
						return
					}
					log.WithFields(log.Fields{
						"station": track.Station,
						"artist": track.Artist,
						"title": track.Title,
					}).Info("Crawled station")

					track.Artist = strings.TrimSpace(track.Artist)
					track.Title = strings.TrimSpace(track.Title)
					tracks <- track
				}(c)
			}
		}
		wg.Wait()
		close(tracks)
	}()

	for track := range tracks {
		go func() {
			w.writer.Write(*track)
		}()
	}
	if counter > 0 {
		log.WithFields(log.Fields{
			"count": counter,
			"duration": time.Now().Sub(start).Seconds(),
		}).Info("Crawling finished")
	}
}

/*
Starts the crawling. This will check ever {interval} seconds if one of the
crawler should been run and starts it.
 */
func (w *Watcher) StartCrawling() {
	w.runCrawlers()
	w.ticker = time.NewTicker(w.refreshInterval)
	go func() {
		for range w.ticker.C {
			w.runCrawlers()
		}
	}()
}

/*
Stops the crawling
 */
func (w *Watcher) StopCrawling() {
	w.ticker.Stop()
}