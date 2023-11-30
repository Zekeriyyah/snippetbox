package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Zekeriyyah/snippetbox/internal/models"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
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

	//Initializing a new templateCache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//Creating instance of application

	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: DB},
		templateCache: templateCache,
		formDecoder:   formDecoder,
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
