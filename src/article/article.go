package article

import (
	"fmt"
	"github.com/jackc/pgx"
)

const (
	ErrorFieldMissing    = "%s is empty"
	ErrorArticleNotFound = "Article not found"
)

type Article struct {
	ID       *int32 `json:"id"`
	AuthorID *int32 `json:"author_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

func (art *Article) Create() error {
	if err := art.Validate(); err != nil {
		return err
	}
	err := conn.QueryRow("INSERT INTO articles(author_id, title, body) VALUES($1,$2,$3) RETURNING id", art.AuthorID, art.Title, art.Body).Scan(&art.ID)
	return err
}

func (art *Article) Read(articleid string) (bool, error) {
	err := conn.QueryRow("SELECT id, author_id, title, body FROM articles WHERE id=$1", articleid).Scan(&art.ID, &art.AuthorID, &art.Title, &art.Body)
	if err == pgx.ErrNoRows {
		return false, fmt.Errorf(ErrorArticleNotFound)
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (art *Article) Update(articleid string) error {
	if err := art.Validate(); err != nil {
		return err
	}
	_, err := conn.Exec("UPDATE articles SET author_id=$1, title=$2, body=$3 WHERE id=$4", art.AuthorID, art.Title, art.Body, articleid)
	return err
}

func (art *Article) Delete(articleid string) error {
	_, err := conn.Exec("DELETE FROM articles WHERE id=$1", articleid)
	return err
}

func (art *Article) Validate() error {
	if art.AuthorID == nil {
		return fmt.Errorf(ErrorFieldMissing, "author_id")
	}
	if art.Title == "" {
		return fmt.Errorf(ErrorFieldMissing, "title")
	}
	if art.Body == "" {
		return fmt.Errorf(ErrorFieldMissing, "body")
	}
	return nil
}
