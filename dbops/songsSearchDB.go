package dbops

import (
	"database/sql"
	"encoding/json"
	"log"
	"musicinfo/models"
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
	songs, err := makeSongsSearchQuery(songDetail, dbCfg.DB)
	if err != nil {
		log.Println("error making query:", err)
		return nil, err
	} else {
		log.Println("the query was successfully completed")
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

func makeSongsSearchQuery(songDetail *models.SongDetail, db *sql.DB) ([]*models.SongDetail, error) {
	log.Println("the makeSongsSearchQuery() func has been called")

	var (
		whereClause string
		args        []interface{}
	)

	// If the ID field is set, use it as the primary search criteria
	if songDetail.ID != 0 {
		whereClause = "id = ?"
		args = append(args, songDetail.ID)

		// Construct the SQL query
		query := "SELECT id, title, release_date, artist, text, lyrics, link FROM song_details WHERE " + whereClause

		results, err := queryToDB(db, query, args)
		if err != nil {
			log.Println("error making query to DB:", err)
			return nil, err
		} else {
			log.Println("the queryToDB() function was completed successfully")
		}

		return results, nil
	}

	// Build the WHERE clause based on the provided parameters
	if songDetail.Artist != "" {
		whereClause = "artist = ?"
		args = append(args, songDetail.Artist)
	}
	if songDetail.Title != "" {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "title = ?"
		args = append(args, songDetail.Title)
	}
	if songDetail.ReleaseDate != "" {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "release_date = ?"
		args = append(args, songDetail.ReleaseDate)
	}

	// Construct the SQL query
	query := "SELECT id, title, release_date, artist, text, lyrics, link FROM song_details"
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	results, err := queryToDB(db, query, args)
	if err != nil {
		log.Println("error making query to DB:", err)
		return nil, err
	}

	return results, nil
}

func queryToDB(db *sql.DB, query string, args []interface{}) ([]*models.SongDetail, error) {
	log.Println("the queryToDB() func has been called")
	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
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
	return results, nil
}
