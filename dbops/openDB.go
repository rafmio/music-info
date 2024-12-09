package dbops

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dotEnvFile = "../config/dbconf.env"
)

type DBConfig struct {
	DriverName string // e.g. "postgres"
	Host       string // "127.0.0.1", "localhost", etc
	Port       string // port number, e.g. "5432", "8543", etc
	DBName     string // name of DB inside of 'PostgreSQL'
	User       string // username "music_lover", "postgres", etc
	Password   string // password
	SslMode    string // SSL mode, etc "disable", "require", "verify-full", etc"
	Dsn        string // data source name
	DB         *sql.DB
}

func NewDBConfig(dbConfigFilePath string) (*DBConfig, error) {
	log.Println("the NewDBConfig has been called")

	// check if dbConfigFilePath is empty
	if dbConfigFilePath == "" {
		log.Println("database config file path is empty")
		return nil, fmt.Errorf("database config file path is empty")
	}

	// Load environment variables from dbconf.env file
	log.Println("load environment variables from .env file")

	err := godotenv.Load(dotEnvFile)
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	dbCfg := new(DBConfig) // new DBConfig instance

	// Assign environment variables to DBConfig's fields
	dbCfg.DriverName = os.Getenv("POSTGRES_DRIVER_NAME")
	dbCfg.Host = os.Getenv("POSTGRES_HOST")
	dbCfg.Port = os.Getenv("POSTGRES_PORT")
	dbCfg.DBName = os.Getenv("POSTGRES_DB")
	dbCfg.User = os.Getenv("POSTGRES_USER")
	dbCfg.Password = os.Getenv("POSTGRES_PASSWORD")
	dbCfg.SslMode = os.Getenv("POSTGRES_SSL_MODE")

	return dbCfg, nil
}

func (dbc *DBConfig) SetDSN() {
	log.Println("the SetDSN() has been called")

	// Define the format string for the DSN
	formatString := "host=%s port=%s user=%s dbname=%s password=%s sslmode=%s"

	dbc.Dsn = fmt.Sprintf(formatString,
		dbc.Host,
		dbc.Port,
		dbc.User,
		dbc.DBName,
		dbc.Password,
		dbc.SslMode,
	)
}

func (dbc *DBConfig) EstablishDbConnection() error {
	log.Println("the EstablishDbConnection() has been called")

	var err error
	dbc.DB, err = sql.Open(dbc.DriverName, dbc.Dsn)
	if err != nil {
		log.Println("Open database:", err)
		return err
	}

	err = dbc.DB.Ping()
	if err != nil {
		log.Println("Ping database:", err)
	}
	return nil
}
