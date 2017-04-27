package content

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ONSdigital/go-ns/log"
	"github.com/lib/pq"
)

type NodeDescription struct {
	Title string `json:"title"`
}

type Node struct {
	URI         string          `json:"uri"`
	Description NodeDescription `json:"description"`
	Type        string          `json:"type"`
}

func GetParent(w http.ResponseWriter, r *http.Request, parentQuery *sql.Stmt) {
	uri := r.URL.Query().Get("uri")
	parents := getParentPages(uri)
	nodes := []Node{}
	rows, err := parentQuery.Query(pq.Array(parents))
	if err != nil {
		log.Error(err, log.Data{"stmt": parentQuery, "uri": uri})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var uri, title, pageType sql.NullString
		if err = rows.Scan(&uri, &title, &pageType); rows.Err() != nil {
			log.Error(err, log.Data{"rows": rows})
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		description := NodeDescription{Title: title.String}
		node := Node{URI: uri.String, Type: pageType.String, Description: description}
		nodes = append(nodes, node)
	}
	data, err := json.Marshal(nodes)
	if err != nil {
		log.Error(err, log.Data{"nodes": nodes})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func getParentPages(uri string) []string {
	uriParts := strings.Split(uri[1:], "/")
	parents := []string{"/?lang=en"}
	var baseURI string
	for i := 0; i < len(uriParts); i++ {
		if uriParts[i] == "timeseries" || uriParts[i] == "bulletins" || uriParts[i] == "datasets" {
			break
		}
		baseURI = baseURI + "/" + uriParts[i]
		if baseURI == uri {
			break
		}
		parents = append(parents, baseURI+"?lang=en")
	}
	return parents
}
