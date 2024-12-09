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
}
