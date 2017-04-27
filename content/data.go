package content

import (
	"database/sql"
	"net/http"

	"github.com/ONSdigital/go-ns/log"
)

func GetData(w http.ResponseWriter, r *http.Request, contentQuery, filterQuery *sql.Stmt) {
	uri := r.URL.Query().Get("uri")
	lang := r.URL.Query().Get("lang")
	filter := r.URL.Query().Get("filter")
	if lang == "" {
		lang = "en"
	}
	fullURI := uri + "?lang=" + lang
	var results *sql.Row
	if filter != "" {
		results = filterQuery.QueryRow(filter, "{"+filter+"}", fullURI)
	} else {
		results = contentQuery.QueryRow(fullURI)
	}
	var content sql.NullString
	notFound := results.Scan(&content)
	if notFound != nil {
		log.ErrorC("Data not found.", notFound, log.Data{"uri": fullURI})
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}
	w.Write([]byte(content.String))
}
