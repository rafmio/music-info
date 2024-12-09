package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"musicinfo/models"

	"net/http"
)

func GetSongMetadataExternal(queryParams *models.QueryParams, config *models.ExternalApiConfig) (*models.SongDetail, error) {
	log.Println("try to get song's metadata from external API")

	// converting the request body to JSON
	requestBodyJSON, err := json.Marshal(queryParams)
	if err != nil {
		log.Println("error marshaling request body:", err)
		return nil, err
	} else {
		log.Println("the request body for external (custom) API has been marshalled to JSON")
	}

	// building URL for the HTTP-request
	url := fmt.Sprintf("http://%s:%s%s", config.Host, config.Port, config.Path)
	log.Printf("the URL for request has been build: %s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		log.Println("error creating request to external (custom) API:", err)
		return nil, err
	} else {
		log.Println("the HTTP-request for external (custom) API has been created")
	}

	// adding the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// make HTTP-request
	// client := &http.Client{}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error making request to external (custom) API:", err)
		return nil, err
	} else {
		log.Println("the HTTP-request to external (custom) API has been made")
	}
	defer resp.Body.Close()

	// reading the request
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body from the external (custom) API:", err)
		return nil, err
	}

	// unmarshalling the response to SongDetail struct
	var songDetail models.SongDetail
	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		log.Println("error unmarshalling response body from the external (custom) API:", err)
		return nil, err
	} else {
		log.Println("the response body from external (custom) API has been unmarshaled to SongDetail struct")
	}

	return &songDetail, nil
}
