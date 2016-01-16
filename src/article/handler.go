package article

import (
	"encoding/json"
	"net/http"

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
func buildErrorResponse(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	b, _ := json.Marshal(&BlogError{code, err})
	w.Write(b)
	return
}

func JsonMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		f(w, r)
	}
}

func HandleSaveArticle(w http.ResponseWriter, r *http.Request) {
	article := new(Article)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&article); err != nil {
		info := "Invalid JSON"
		buildErrorResponse(w, http.StatusBadRequest, info)
		return
	}
	err := article.Create()
	if err != nil {
		buildErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	b, _ := json.Marshal(article)
	w.Write(b)
}

func HandleLoadArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := new(Article)
	exist, err := article.Read(vars["articleid"])
	if !exist {
		buildErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		buildErrorResponse(w, http.StatusBadRequest, err.Error())
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
		info := "Invalid JSON"
		buildErrorResponse(w, http.StatusBadRequest, info)
		return
	}
	err := article.Update(vars["articleid"])
	if err != nil {
		buildErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	b, _ := json.Marshal(article)
	w.Write(b)
}

func HandleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := new(Article)
	err := article.Delete(vars["articleid"])
	if err != nil {
		buildErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
