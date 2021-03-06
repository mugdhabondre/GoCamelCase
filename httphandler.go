package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/mugdhabondre/gocamelcase/gocamelcaseimpl"
)


// Request Handler to take in word and return camelcased output
// TODO: accept input
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phrase := vars["phrase"]
	io.WriteString(w, fmt.Sprint("Input:", phrase, "\n"))
	result, error :=  gocamelcaseimpl.ProcessPhrase(phrase)
	if error != nil {
		io.WriteString(w, fmt.Sprint("Oops, looks like there was an Error:\n" + error.Error() ))
	} else {
		io.WriteString(w, fmt.Sprint("Result:", result, "\n"))
	}
}

// HelloHandler return hello world message`
func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Welcome to the camelCase generator!\n" +
					  "Append a string to the url that needs to be camelcased as :\n" + 
					  "https://<hostname>:<port>/camelcase/<your string>\n" +
					  "example: https://camelcasegenerator.azurewebsites.net/camelcase/ilikeapple")
}

func main() {
	fmt.Println("Server started .....")

	router := mux.NewRouter()
	// HandleFunc registers the handler function for the given pattern
	// in the DefaultServeMux.
	router.HandleFunc("/camelcase/{phrase}", RequestHandler)
	router.HandleFunc("/",  HelloHandler)
	// Get port number from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}

	fmt.Println("Listen serve for port... ", port)
	// in case of error, log.Fatal will exit application
	// Listen and Server in 0.0.0.0:3005
	log.Fatal(http.ListenAndServe(":"+port, router))
}

