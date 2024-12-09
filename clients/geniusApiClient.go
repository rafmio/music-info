package clients

import (
	"errors"
	"log"
	"strings"

	"github.com/natecham/genius"

	"musicinfo/models"
)

func GetSongMetadata(queryParams *models.QueryParams, accessToken string) (*models.SongDetail, error) {
	// accessToken := "ycPo8Ic25mjQO6-Kcka0gUeeZDNNDmhCFh1o19YNPLGXY4W95_xXGuPwCDVjBeJA"

	// Create a new client for the Genius API
	client := genius.NewClient(nil, accessToken)
	log.Println("the client for genius.com has been created")

	// Search for the selected song
	results, err := client.Search(queryParams.Song)
	if err != nil {
		return nil, err
	}

	// Extract the song details from the first hit in the search results
	if results.Meta.Status == 200 {
		if len(results.Response.Hits) > 0 {
			// if the found result is single one
			if len(results.Response.Hits) == 1 {
				song := results.Response.Hits[0].Result
				songDetail := FillSongDetail(song)
				log.Println("the requested song has been extracted")
				// filling lyrics (text)
				log.Println("getting lyrics...")
				songDetail.Text, err = GetText(client, songDetail.Link)
				if err != nil {
					log.Println("cannot getting lyrics (text)")
					return nil, err
				}

				log.Println("lyrics (text) has been got successfully")

				return songDetail, nil
			} else {
				// if result of searching contains more than one result
				log.Println("multiple songs found. Trying to choose one of them...")
				song := ChooseSong(queryParams, results.Response.Hits)
				songDetail := FillSongDetail(song)

				log.Println("the requested song has been extracted")

				// filling lyrics (text)
				log.Println("getting lyrics...")
				songDetail.Text, err = GetText(client, songDetail.Link)
				if err != nil {
					log.Println("cannot getting lyrics (text)")
					return nil, err
				}

				log.Println("lyrics (text) has been got successfully")

				return songDetail, nil
			}
		} else {
			log.Printf("error extracting the song. Status: %v", results.Meta.Status)
			return nil, nil
		}
	} else {
		log.Printf("error extracting the song. Status: %v", results.Meta.Status)
		return nil, err
	}
}

// there may be several songs. matches band name (group) and song name
func ChooseSong(queryParams *models.QueryParams, songs []*genius.Hit) *genius.Song {
	log.Printf("ChooseSong() called with group: %s, song: %s", queryParams.Group, queryParams.Song)

	for _, hit := range songs {
		log.Printf("checking song: %s by %s", hit.Result.Title, hit.Result.PrimaryArtist.Name)
		// check if band name matches with "group" name
		if strings.Contains(strings.ToLower(hit.Result.PrimaryArtist.Name), strings.ToLower(queryParams.Group)) {
			log.Printf("match found: %s by %s", hit.Result.Title, hit.Result.PrimaryArtist.Name)
			return hit.Result
		}
	}

	log.Println("ChooseSong(): no matching song found")
	return nil // no matching songs found, return nil
}

func FillSongDetail(song *genius.Song) *models.SongDetail {
	songDetail := &models.SongDetail{
		ID:          song.ID,
		Title:       song.FullTitle,
		ReleaseDate: song.ReleaseDateForDisplay,
		Artist:      song.PrimaryArtist.Name,
		Link:        song.URL,
	}

	log.Println("FillSongDetail(): the models.SongDetail structure has been filled")
	return songDetail
}

func GetText(client *genius.Client, link string) (string, error) {
	// get lyrics
	text, err := client.GetLyrics(link)
	if err != nil {
		log.Println("error getting lyrics (text):", err)
		return "", err
	}

	// check if text (lyrics) is empty
	if text == "" {
		return "", errors.New("lyrics (text) variable is empty")
	}

	return text, nil
}
