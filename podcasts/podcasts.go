package podcasts

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/SlyMarbo/rss"
)

// Podcast ...
type Podcast struct {
	gorm.Model
	Name     string `json:"name,omitempty"`
	Feed     string `json:"feed,omitempty"`
	Episodes []Episode
}

// Episode ...
type Episode struct {
	gorm.Model
	Name      string    `json:"name,omitempty"`
	Link      string    `json:"link,omitempty"`
	Date      time.Time `json:"date,omitempty"`
	Status    string    `json:"status,omitempty"`
	PodcastID uint
	Podcast   Podcast
}

var db *gorm.DB

// SetDB set database connector
func SetDB(dbArg *gorm.DB) {
	db = dbArg
}

// AddPodcastFromFeed add podcast in database with feed URL
func AddPodcastFromFeed(feedURL string) (podcast Podcast, err error) {

	// Assert podcast don't already exist
	var count int
	db.Model(&podcast).Where("feed = ?", feedURL).Count(&count)

	if count > 0 {
		err = errors.New("podcast already exists")
		return
	}

	var feed *rss.Feed
	feed, err = rss.Fetch(feedURL)

	if err != nil {
		return
	}

	podcast = Podcast{}
	podcast.Feed = feedURL
	podcast.Name = feed.Title

	AddPodcast(&podcast)

	for _, item := range feed.Items {
		var episode = Episode{}
		episode.Name = item.Title
		episode.Link = item.Link
		episode.Date = item.Date
		episode.Status = "notlistened"
		episode.PodcastID = podcast.ID
		AddEpisode(&episode)
	}

	return
}

// AddPodcast store podcast in database
func AddPodcast(podcast *Podcast) {
	db.Create(podcast)
}

// DeletePodcast delete podcast from database
func DeletePodcast(podcastID uint) {
	var podcast Podcast
	epidodes := []Episode{}

	db.First(&podcast, podcastID)
	db.Model(&podcast).Related(&epidodes)

	// Delete episodes
	for _, episode := range epidodes {
		db.Delete(&episode)
	}

	// Delete podcast
	db.Delete(&podcast)
}

// AddEpisode store episode in database
func AddEpisode(episode *Episode) {
	db.Create(episode)
}
