package content

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ONSdigital/go-ns/log"

	"github.com/ONSdigital/dp-content-api/s3"
)

func GetResource(w http.ResponseWriter, r *http.Request, st *sql.Stmt, s3Client s3.S3Client) {
	uri := r.URL.Query().Get("uri")
	results := st.QueryRow(uri)
	var s3Location sql.NullString
	notFound := results.Scan(&s3Location)
	if notFound != nil {
		log.ErrorC("Resource not found", notFound, log.Data{"uri": uri})
		http.Error(w, "Resource not found", http.StatusNotFound)
	} else {
		s3uri := strings.TrimPrefix(s3Location.String, "s3://"+s3Client.GetBucket()+"/")
		data, err := s3Client.GetObject(s3uri)
		if err != nil {
			log.ErrorC("Resource not found", err, log.Data{"uri": uri})
			http.Error(w, "Content not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(s3uri))
		w.Write(data)
	}
}
