package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_"github.com/jackc/pgx/v4"
	_"github.com/jackc/pgx/v4/stdlib"
	_"github.com/jackc/pgconn"

)

const webPort = ":80"
var counts int64; 


type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {
	// Start the server
	log.Println("Starting the server on port", webPort)
	// todo: connect to db 
	conn := connectToDb() 

	if conn == nil {
		log.Panic("Error connecting to postgres db")
	}

	// set up config 
	app := Config{
		DB: conn,	
		Models: data.New(conn),
	}

	srv := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic("Error starting server", err)
	}


}

func openDB(dsn string) (*sql.DB, error) {

	db, err:= sql.Open("pgx",dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDb() *sql.DB {

	dsn:= os.Getenv("DSN")
	for {
		connection, err:= openDB(dsn)

		if err != nil {
			log.Println("Error: postgres not yet ready", err)
			counts++; 
		} else {
			log.Println("Connected to postgres")
			return connection;
		}

		if counts > 10 {
			log.Panic("Error: postgres not yet ready")
			return nil;
		}

		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue;
	}
}