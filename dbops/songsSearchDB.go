package dbops

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"musicinfo/models"
	"strconv"
	"strings"
)

var (
	ErrSongNotFound = errors.New("no song details found in the DB")
)

func SongsSearchDB(songDetail *models.SongDetail) ([]*models.SongDetail, error) {
	log.Println("the SongsSearchDB() has been called")

	// try to connect to DB
	dbCfg, err := NewDBConfig(dotEnvFile)
	if err != nil {
		log.Println("error creating DBConfig:", err)
		return nil, err
	}

	// setting data source name for DB
	dbCfg.SetDSN()
	log.Println("now DSN set to:", dbCfg.Dsn)

	// connect to DB
	err = dbCfg.EstablishDbConnection()
	if err != nil {
		log.Println("error establishing DB connection:", err)
		return nil, err
	} else {
		log.Println("the connection to DB has been established")
	}
	defer dbCfg.DB.Close()

	// try to make a query
	queryString := buildSongSearchQuery(songDetail)
	log.Println("the query was successfully completed")

	log.Println("try to make SQL-query to DB...")
	songs, err := makeQueryToDB(dbCfg.DB, queryString)
	if err != nil {
		if err != ErrSongNotFound {
			log.Println("error making query to DB:", err)
		}
		if err == ErrSongNotFound {
			log.Println("the song you are looking for has not been found")
		}
		return nil, err
	}
	log.Println("the query to DB was successfully executed")

	return songs, nil
}

func buildSongSearchQuery(song *models.SongDetail) string {
	params := make(map[string]string)
	params["id"] = strconv.Itoa(song.ID)
	params["title"] = song.Title
	params["artist"] = song.Artist
	params["release_date"] = song.ReleaseDate

	selectClause := "SELECT * FROM song_details WHERE %s ;"
	var whereClause []string

	for key, value := range params {
		if value == "" || value == "0" {
			continue
		}
		// Using ILIKE for case-insensitive comparison
		// Escaping single quotes: We continue to escape single quotes in strings
		// using strings.replaceAll(value, "'", """), to avoid SQL\errors.
		whereClause = append(whereClause, fmt.Sprintf("%s ILIKE '%s'", key, strings.ReplaceAll(value, "'", "''")))
	}

	finalWhereClause := strings.Join(whereClause, " AND ")
	finalQueryString := fmt.Sprintf(selectClause, finalWhereClause)

	return finalQueryString
}

// func buildSongSearchQuery(song *models.SongDetail) string {
// 	log.Println("buildSongSearchQuery() has been called")
// 	log.Println("start to building SQL-query string...")

// 	params := make(map[string]string)
// 	params["id"] = strconv.Itoa(song.ID)
// 	params["title"] = song.Title
// 	params["artist"] = song.Artist
// 	params["release_date"] = song.ReleaseDate

// 	selectClause := "SELECT * FROM song_details WHERE %s ;"
// 	var whereClause string

// 	for key, value := range params {
// 		if value == "" || value == "0" {
// 			continue
// 		}
// 		whereClause += key + "=" + "'" + value + "'" + " AND "
// 	}
// 	whereClause = strings.TrimSuffix(whereClause, " AND ")

// 	finalQueryString := fmt.Sprintf(selectClause, whereClause)

// 	log.Println("the SQL-query string has been successfully built:", finalQueryString)

// 	return finalQueryString
// }

func makeQueryToDB(db *sql.DB, query string) ([]*models.SongDetail, error) {
	log.Println("the makeQueryToDB() func has been called")
	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	log.Println("parsing response to []*models.SongDetail...")

	defer rows.Close()
	// Collect the results
	results := make([]*models.SongDetail, 0)

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
		log.Println("error while iterating over rows:", err)
		return nil, err
	}

	if len(results) == 0 {
		log.Println("emtpy 'results', len(results):", len(results), results)
		return nil, ErrSongNotFound
	} else {
		log.Println("non-empty results, len(results):", len(results), results)
	}

	return results, nil
}
