package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	//Implementing flag usage in HTTP Network address
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	//Creating flag dsn for database data source name
	dsn := flag.String("dsn", "web:Awwalweb@db1@/snippetbox?parseTime=true", "Database Data Source Name")

	flag.Parse()

	/**
	//Create a logger for writing information message into a file
	f, errf := os.OpenFile("/tmp/info.go", os.O_RDWR|os.O_CREATE, 0666)
	if errf != nil {
		log.Fatal(errf)
	}
	defer f.Close()
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	**/

	//Create a logger for writing information message into standard output stream
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//Create a logger for writing error message to the terminal
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initializing database
	DB, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	} else {
		infoLog.Println("Database successfully Initialized...")
	}

	defer DB.Close()

	//Creating instance of application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Starting server on port ", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
