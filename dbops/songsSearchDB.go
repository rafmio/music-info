package dbops

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"musicinfo/models"
	"strings"
)

func SongsSearchDB(songDetail *models.SongDetail) ([]byte, error) {
	log.Println("the SongsSearchDB() has been called")

	// try to connect to DB
	dbCfg, err := NewDBConfig(dotEnvFile)
	if err != nil {
		log.Println("Error creating DBConfig:", err)
		return nil, err
	}

	// setting data source name for DB
	dbCfg.SetDSN()
	log.Println("now DSN set to:", dbCfg.Dsn)

	err = dbCfg.EstablishDbConnection()
	if err != nil {
		log.Println("error establishing DB connection:", err)
		return nil, err
	} else {
		log.Println("the connection to DB has been established")
	}
	defer dbCfg.DB.Close()

	// try to make a query
	log.Println("trying to make a query...")
	query, params := buildSongSearchQuery(songDetail)
	if err != nil {
		log.Println("error making query:", err)
		return nil, err
	} else {
		log.Println("the query was successfully completed")
	}

	songs, err := makeQueryToDB(dbCfg.DB, query, params)
	if err != nil {
		log.Println("error making query to DB")
	}

	log.Println("encoding data to JSON")
	jsonSongs, err := json.Marshal(songs)
	if err != nil {
		log.Println("error encoding data to JSON")
	} else {
		log.Println("the data is encoded in JSON successfully")
	}

	return jsonSongs, nil
}

func buildSongSearchQuery(song *models.SongDetail) (string, []interface{}) {
	log.Println("start to building SQL-query string...")

	var queryParts []string
	var params []interface{}

	if song.ID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("id = $%d", len(params)+1))
		params = append(params, song.ID)
	}

	if song.Title != "" {
		queryParts = append(queryParts, fmt.Sprintf("title ILIKE $%d", len(params)+1))
		params = append(params, "%"+song.Title+"%") // Добавляем wildcards для частичного совпадения
	}

	if song.ReleaseDate != "" { // Проверяем наличие непустой даты
		queryParts = append(queryParts, fmt.Sprintf("release_date = $%d", len(params)+1))
		params = append(params, song.ReleaseDate)
	}

	if song.Artist != "" {
		queryParts = append(queryParts, fmt.Sprintf("artist ILIKE $%d", len(params)+1))
		params = append(params, "%"+song.Artist+"%") // Добавляем wildcards для частичного совпадения
	}

	baseQuery := `
        SELECT seq_num, id, title, release_date, artist, lyrics, link
        FROM song_details
        WHERE %s`

	finalQuery := fmt.Sprintf(baseQuery, strings.Join(queryParts, " AND "))

	log.Println("the query string is built")

	return finalQuery, params
}

func makeQueryToDB(db *sql.DB, query string, args []interface{}) ([]*models.SongDetail, error) {
	log.Println("the makeQueryToDB() func has been called")
	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	log.Println("parsing response to []*models.SongDetail...")

	defer rows.Close()
	// Collect the results
	var results []*models.SongDetail
	for rows.Next() {
		song := &models.SongDetail{}
		if err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.ReleaseDate,
			&song.Artist,
			&song.Text,
			&song.Lyrics,
			&song.Link,
		); err != nil {
			return nil, err
		}
		results = append(results, song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	log.Println("results:", results)
	log.Println("len(results):", len(results))
	return results, nil
}
