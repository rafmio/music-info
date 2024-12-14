package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"musicinfo/dbops"
	"musicinfo/models"
)

// http://localhost:8080/v1/songs/search?id=1232&title=Supermassive%20Black%20Hole&artist=Muse&release_date=10-05-2006
func SongsSearchGet(w http.ResponseWriter, r *http.Request) {
	log.Println("inside the SongsSearchGet() func")

	// check if request's method is GET:
	if r.Method != http.MethodGet {
		log.Println("the SongsSearchGet() receive wrong method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	songDetail := new(models.SongDetail)
	// parse query parameters:
	songDetail.ID, _ = strconv.Atoi(r.URL.Query().Get("id"))
	songDetail.Title = r.URL.Query().Get("title")
	songDetail.Artist = r.URL.Query().Get("artist")
	songDetail.ReleaseDate = r.URL.Query().Get("release_date")

	// TODO: Implement pagination

	// search song
	log.Println("trying to find song")
	songs, err := dbops.SongsSearchDB(songDetail)
	if err != nil {
		if err == dbops.ErrSongNotFound {
			log.Println("the song you are looking for has not been found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Println("error searching for song:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		jsonSongs, err := json.Marshal(songs)
		if err != nil {
			log.Println("error marshaling songs:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonSongs)
		if err != nil {
			log.Println("error writing response:", err)
		}
	}
}
