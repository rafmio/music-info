package dbops

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"musicinfo/models"
	"strings"
)

var (
	ErrSongNotFound = errors.New("no song details found in the DB")
)

func SongsSearchDB(songDetail *models.SongDetail) ([]*models.SongDetail, error) {
	log.Println("The SongsSearchDB() has been called")
	dbCfg, err := NewDBConfig(dotEnvFile)
	if err != nil {
		log.Println("Error creating DBConfig:", err)
		return nil, err
	}
	dbCfg.SetDSN()
	log.Println("Now DSN set to:", dbCfg.Dsn)
	err = dbCfg.EstablishDbConnection()
	if err != nil {
		log.Println("Error establishing DB connection:", err)
		return nil, err
	} else {
		log.Println("The connection to DB has been established")
	}
	defer dbCfg.DB.Close()

	query, params := buildSongSearchQuery(songDetail)
	songs, err := makePreparedQueryToDB(dbCfg.DB, query, params)
	if err != nil {
		if err != ErrSongNotFound {
			log.Println("Error making query to DB:", err)
		}
		if err == ErrSongNotFound {
			log.Println("The song you are looking for has not been found")
		}
		return nil, err
	}
	return songs, nil
}

// func SongsSearchDB(songDetail *models.SongDetail) ([]*models.SongDetail, error) {
// 	log.Println("the SongsSearchDB() has been called")

// 	// try to connect to DB
// 	dbCfg, err := NewDBConfig(dotEnvFile)
// 	if err != nil {
// 		log.Println("error creating DBConfig:", err)
// 		return nil, err
// 	}

// 	// setting data source name for DB
// 	dbCfg.SetDSN()
// 	log.Println("now DSN set to:", dbCfg.Dsn)

// 	// connect to DB
// 	err = dbCfg.EstablishDbConnection()
// 	if err != nil {
// 		log.Println("error establishing DB connection:", err)
// 		return nil, err
// 	} else {
// 		log.Println("the connection to DB has been established")
// 	}
// 	defer dbCfg.DB.Close()

// 	// try to make a query
// 	queryString := buildSongSearchQuery(songDetail)
// 	log.Println("the query was successfully completed")

// 	log.Println("try to make SQL-query to DB...")
// 	songs, err := makeQueryToDB(dbCfg.DB, queryString)
// 	if err != nil {
// 		if err != ErrSongNotFound {
// 			log.Println("error making query to DB:", err)
// 		}
// 		if err == ErrSongNotFound {
// 			log.Println("the song you are looking for has not been found")
// 		}
// 		return nil, err
// 	}
// 	log.Println("the query to DB was successfully executed")

// 	return songs, nil
// }

func buildSongSearchQuery(song *models.SongDetail) (string, []interface{}) {
	var conditions []string
	var params []interface{}
	paramIndex := 1

	if song.ID != 0 {
		conditions = append(conditions, fmt.Sprintf("id = $%d", paramIndex))
		params = append(params, song.ID)
		paramIndex++
	}
	if song.Title != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(title) LIKE LOWER($%d)", paramIndex))
		params = append(params, song.Title)
		paramIndex++
	}
	if song.Artist != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(artist) LIKE LOWER($%d)", paramIndex))
		params = append(params, song.Artist)
		paramIndex++
	}
	if song.ReleaseDate != "" {
		conditions = append(conditions, fmt.Sprintf("release_date = $%d", paramIndex))
		params = append(params, song.ReleaseDate)
		paramIndex++
	}

	queryTemplate := "SELECT * FROM song_details WHERE %s;"
	whereClause := strings.Join(conditions, " AND ")

	return fmt.Sprintf(queryTemplate, whereClause), params
}

// func buildSongSearchQuery(song *models.SongDetail) string {
// 	params := make(map[string]string)
// 	params["id"] = strconv.Itoa(song.ID)
// 	params["title"] = song.Title
// 	params["artist"] = song.Artist
// 	params["release_date"] = song.ReleaseDate

// 	selectClause := "SELECT * FROM song_details WHERE %s ;"
// 	var whereClause []string

// 	for key, value := range params {
// 		if value == "" || value == "0" {
// 			continue
// 		}
// 		// Using 'ILIKE' for case-insensitive comparison
// 		// Escaping single quotes: We continue to escape single quotes in strings
// 		// using strings.replaceAll(value, "'", """), to avoid SQL\errors.
// 		whereClause = append(whereClause, fmt.Sprintf("%s ILIKE '%s'", key, strings.ReplaceAll(value, "'", "''")))
// 	}

// 	finalWhereClause := strings.Join(whereClause, " AND ")
// 	finalQueryString := fmt.Sprintf(selectClause, finalWhereClause)

// 	log.Println("builded SQL query is:", finalQueryString)

// 	return finalQueryString
// }

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

// Функция для выполнения подготовленного запроса
func makePreparedQueryToDB(db *sql.DB, query string, params []interface{}) ([]*models.SongDetail, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		log.Printf("Error executing prepared query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.SongDetail, 0)
	for rows.Next() {
		song := &models.SongDetail{}
		if err := rows.Scan(&song.ID, &song.Title, &song.ReleaseDate, &song.Artist, &song.Text, &song.Lyrics, &song.Link); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		results = append(results, song)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v\n", err)
		return nil, err
	}

	if len(results) == 0 {
		log.Println("Empty results.")
		return nil, ErrSongNotFound
	} else {
		log.Printf("Non-empty results, len(results): %d\n", len(results))
	}

	return results, nil
}

// Вспомогательная функция для преобразования карты параметров в срез значений
// func paramsToSlice(params map[string]interface{}) []interface{} {
// 	values := make([]interface{}, 0, len(params))
// 	for _, v := range params {
// 		values = append(values, v)
// 	}
// 	return values
// }

// func makeQueryToDB(db *sql.DB, query string) ([]*models.SongDetail, error) {
// 	log.Println("the makeQueryToDB() func has been called")
// 	// Execute the query
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Println("parsing response to []*models.SongDetail...")

// 	defer rows.Close()
// 	// Collect the results
// 	results := make([]*models.SongDetail, 0)

// 	for rows.Next() {
// 		song := &models.SongDetail{}
// 		if err := rows.Scan(
// 			&song.ID,
// 			&song.Title,
// 			&song.ReleaseDate,
// 			&song.Artist,
// 			&song.Text,
// 			&song.Lyrics,
// 			&song.Link,
// 		); err != nil {
// 			return nil, err
// 		}
// 		results = append(results, song)
// 	}
// 	if err := rows.Err(); err != nil {
// 		log.Println("error while iterating over rows:", err)
// 		return nil, err
// 	}

// 	if len(results) == 0 {
// 		log.Println("emtpy 'results', len(results):", len(results), results)
// 		return nil, ErrSongNotFound
// 	} else {
// 		log.Println("non-empty results, len(results):", len(results), results)
// 	}

// 	return results, nil
// }
