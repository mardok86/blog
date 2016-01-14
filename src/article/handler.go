package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

var conn *pgx.Conn

const (
	DBDriver = "postgres"
	DBName   = "marco"
	DBHost   = "127.0.0.1"
)

func Setup() error {
	var err error
	conn, err = pgx.Connect(extractConfig())
	return err
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig
	config.Host = DBHost
	config.Database = DBName

	return config
}

func HandleSaveArticle(w http.ResponseWriter, r *http.Request) {
	article := new(Article)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&article); err != nil {
		fmt.Println("Error decode json")
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := article.Create()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(&BlogError{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(article)
	w.Write(b)
}

func HandleLoadArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := new(Article)
	exist, err := article.Read(vars["articleid"])
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		b, _ := json.Marshal(&BlogError{http.StatusBadRequest, err.Error()})
		w.Write(b)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(&BlogError{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(article)
	w.Write(b)
}

func HandleUpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := new(Article)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&article); err != nil {
		fmt.Println("Database connection error: ")
		os.Exit(1)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := article.Update(vars["articleid"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(&BlogError{http.StatusBadRequest, err.Error()})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(article)
	w.Write(b)
}

func HandleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := new(Article)
	err := article.Delete(vars["articleid"])
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(&BlogError{http.StatusBadRequest, err.Error()})
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
