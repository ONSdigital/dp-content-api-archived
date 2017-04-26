package content

import (
	"database/sql"
	"log"
	"net/http"
)

func GetData(w http.ResponseWriter, r *http.Request, contentQuery, filterQuery *sql.Stmt) {
	uri := r.URL.Query().Get("uri")
	lang := r.URL.Query().Get("lang")
	filter := r.URL.Query().Get("filter")
	if lang == "" {
		lang = "en"
	}
	fullURL := uri + "?lang=" + lang
	var results *sql.Row
	if filter != "" {
		results = filterQuery.QueryRow(filter, filter, fullURL)
	} else {
		results = contentQuery.QueryRow(fullURL)
	}
	var content sql.NullString
	notFound := results.Scan(&content)
	if notFound != nil {
		log.Printf("Data not found. uri : %s, language : %s, %s", fullURL, lang, notFound.Error())
		http.Error(w, "Content not found", http.StatusNotFound)
	}
	w.Write([]byte(content.String))
}
