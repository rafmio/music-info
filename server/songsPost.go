package server

import (
	"encoding/json"
	"errors"
	"log"
	"musicinfo/clients"
	"musicinfo/dbops"
	"musicinfo/models"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	// dotEnvFile = "../config/dbconf.env" // path to the .env file
	dotEnvFile = "config/dbconf.env" // path to the .env file
)

type ExternalApiConfig struct {
	Host        string
	Port        string
	AccessToken string // unexported field
	Path        string
}

// request example:
// curl -X POST \
//   http://localhost:8080/v1/songs \
//   -H 'Content-Type: application/json' \
//   -d '{
//   "group": "Muse",
//   "song": "Supermassive Black Hole"
// }'

func SongsPost(w http.ResponseWriter, r *http.Request) {
	log.Println("the SongsPost() function has been called")

	// check if request's method is POST:
	if r.Method != http.MethodPost {
		log.Println("the SongsSearchGet() receive wrong method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// parse request to the QueryParams struct
	var queryParams models.QueryParams
	err := json.NewDecoder(r.Body).Decode(&queryParams)
	if err != nil {
		log.Println("parsing request to add (POST) new song: invalid request body", err)
		http.Error(w, "parsing request to add (POST) new song: Invalid request body", http.StatusBadRequest)
		return
	} else {
		log.Println("the request's body has been parsed to QueryParams struct")
	}

	// check if the song to be added (POSTed) already exists in the database to avoid duplicates
	songsDetail := new(models.SongDetail)
	songsDetail.Artist = queryParams.Group
	songsDetail.Title = queryParams.Song

	jsonSongs, err := dbops.SongsSearchDB(songsDetail)
	if err != nil {
		log.Println("error searching in DB")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(jsonSongs) > 0 {
		log.Println("the song you are trying to add (POST) already exists in the DB")
		w.WriteHeader(http.StatusConflict) // 409
		return
	}

	// call external API to get song metadata
	songDetail, err := callExternalApi(&queryParams)
	if err != nil {
		log.Println("error getting song metadata from external API:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Println("the song metadata has been retrieved from external API")
	}

	// save the song to the database
	err = dbops.SongsInsertDB(songDetail)
	if err != nil {
		log.Println("error inserting (POST) song into DB:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ???
	// return the newly added song's details
	// jsonSongs, err = dbops.SongsSearchDB(songDetail)
	// if err != nil {
	// 	log.Println("error searching in DB")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// json.NewEncoder(w).Encode(jsonSongs) // marshal and write JSON response to the client
	// ???

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func callExternalApi(queryParams *models.QueryParams) (*models.SongDetail, error) {
	log.Println("the callExternalApi() has been called")

	songDetail := new(models.SongDetail)

	// selects the API to use based on the configuration. By default, it uses the Genius API
	apiConfig, err := selectApi() // apiConfig is of type ExternalApiConfig
	if err != nil {
		log.Println("error selecting API source:", err)
		return nil, err
	}
	if apiConfig.AccessToken == "" {
		songDetail, err = clients.GetSongMetadataExternal(queryParams, apiConfig)
		if err != nil {
			log.Println("error getting song metadata from custom external API:", err)
			return nil, err
		}

	} else {
		songDetail, err = clients.GetSongMetadata(queryParams, apiConfig.AccessToken)
		if err != nil {
			log.Println("error getting song metadata from genius.com API:", err)
			return nil, err
		}
	}

	return songDetail, nil
}

func selectApi() (*models.ExternalApiConfig, error) {
	log.Println("the selectApi() has been called")

	err := godotenv.Load(dotEnvFile)
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	externalApiConfig := new(models.ExternalApiConfig)
	// retrieve API configuration from environment variables
	if os.Getenv("MUSIC_INFO_USE_GENIUS_API") == "true" {
		log.Println("the geinus.com API has been selected as default")

		externalApiConfig.AccessToken = os.Getenv("GENIUS_API_ACCESS_TOKEN")

		if externalApiConfig.AccessToken == "" {
			return nil, errors.New("missing required environment variables for Genius API")
		}
	} else {
		log.Println("the custom external API has been selected")

		externalApiConfig.Host = os.Getenv("EXTERNAL_API_HOST")
		externalApiConfig.Port = os.Getenv("EXTERNAL_API_PORT")
		externalApiConfig.Path = os.Getenv("EXTERNAL_API_PATH")
		externalApiConfig.AccessToken = ""

		if externalApiConfig.Host == "" || externalApiConfig.Port == "" || externalApiConfig.Path == "" {
			return nil, errors.New("missing required environment variables for custom external API")
		}

	}

	return externalApiConfig, nil
}
