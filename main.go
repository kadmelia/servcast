package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"servcast/podcasts"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

// GetPodcasts return all podcast in JSON.
func GetPodcasts(w http.ResponseWriter, r *http.Request) {
	pods := []podcasts.Podcast{}
	db.Find(&pods)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

// GetPodcast return required podcast in JSON.
func GetPodcast(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var podcast podcasts.Podcast

	db.First(&podcast, params["id"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(podcast)
}

// GetEpisodes return required podcast episodes in JSON.
func GetEpisodes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	episodes := []podcasts.Episode{}

	db.Find(&episodes, "podcast_id = ?", params["id"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
}

// CreatePodcast create a podcast with given data.
func CreatePodcast(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	feed := r.Form.Get("feed")

	if feed == "" {
		panic("Can't create podcast without feed")
	}
	var podcast podcasts.Podcast
	podcast, err = podcasts.AddPodcastFromFeed(feed)

	if err != nil {
		http.Error(w, err.Error(), 409)
		return
	}

	json.NewEncoder(w).Encode(podcast)

}

// DeletePodcast delete required podcast
func DeletePodcast(w http.ResponseWriter, r *http.Request) {
	pods := []podcasts.Podcast{}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		panic("invalid id given for podcast delete")
	}

	podcasts.DeletePodcast(uint(id))

	db.Find(&pods)

	json.NewEncoder(w).Encode(pods)
}

// our main function
func main() {
	var err error

	// podcasts.GetPodcastsURL()

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	podcasts.SetDB(db)

	// Migrate the schema
	db.AutoMigrate(&podcasts.Podcast{})
	db.AutoMigrate(&podcasts.Episode{})

	router := mux.NewRouter()
	router.HandleFunc("/podcasts", GetPodcasts).Methods("GET")
	router.HandleFunc("/podcasts/{id}", GetPodcast).Methods("GET")
	router.HandleFunc("/podcasts/{id}/episodes", GetEpisodes).Methods("GET")

	router.HandleFunc("/podcasts", CreatePodcast).Methods("POST")
	router.HandleFunc("/podcasts/{id}", DeletePodcast).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

}
