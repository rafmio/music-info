package dbops

import (
	"log"
	"musicinfo/models"
)

func SongsInsertDB(songDetail *models.SongDetail) error {
	log.Println("the SongsInsertDB() has been called")

	dbConf, err := NewDBConfig(dotEnvFile)
	if err != nil {
		log.Println("error creating DB config from the .env file")
		return err
	} else {
		log.Println("the DBConfig has been created from the.env file")
	}

	dbConf.SetDSN()
	err = dbConf.EstablishDbConnection()
	if err != nil {
		log.Println("error establishing database connection")
		return err
	} else {
		log.Println("the database connection has been established")
	}
	defer dbConf.DB.Close()

	// building the SQL-query string
	query := `
		INSERT INTO song_details (id, title, release_date, artist, lyrics, link)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = dbConf.DB.Exec(query,
		songDetail.ID,
		songDetail.Title,
		songDetail.ReleaseDate,
		songDetail.Artist,
		songDetail.Lyrics,
		songDetail.Link,
	)

	return err
}
