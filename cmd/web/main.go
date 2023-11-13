package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	//Implementing flag usage in HTTP Network address
	addr := flag.String("addr", ":4000", "HTTP Network Address")

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

	//Creating instance of application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	//Register mux To handle static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//Register the other application routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Println("Starting server on port ", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
