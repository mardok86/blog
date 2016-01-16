package main

import (
	"article"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	err := article.Setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/articles", article.JsonMiddleware(article.HandleSaveArticle)).Methods("POST")
	router.HandleFunc("/v1/articles/{articleid:[0-9]+}", article.JsonMiddleware(article.HandleLoadArticle)).Methods("GET")
	router.HandleFunc("/v1/articles/{articleid:[0-9]+}", article.JsonMiddleware(article.HandleUpdateArticle)).Methods("PUT")
	router.HandleFunc("/v1/articles/{articleid:[0-9]+}", article.JsonMiddleware(article.HandleDeleteArticle)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":4444", router))
}
