package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	//Implementing flag usage in HTTP Network address
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()

	//Create a logger for writing information message
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//Create a logger for writing error message to the terminal
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	//Register mux To handle static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//Register the other application routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	infoLog.Println("Starting server on port ", *addr)
	errorLog.Fatal(http.ListenAndServe(*addr, mux))
}
