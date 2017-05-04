package content

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-content-api/s3"
	"github.com/ONSdigital/go-ns/log"
)

type FileSizeMessage struct {
	FileSize int64 `json:"fileSize"`
}

func GetFileSize(w http.ResponseWriter, r *http.Request, st *sql.Stmt, s3Client s3.S3Client) {
	uri := r.URL.Query().Get("uri")
	row := st.QueryRow(uri)
	var s3uri sql.NullString
	err := row.Scan(&s3uri)
	if err != nil {
		log.ErrorC("Data not found in datastore", err, log.Data{"uri": uri})
		http.Error(w, "File not found in datastore", http.StatusNotFound)
		return
	}
	s3UriWithOutPrefix := strings.TrimLeft(s3uri.String, "s3://"+s3Client.GetBucket())
	fileSize, err := s3Client.GetFileSize(s3UriWithOutPrefix)
	if err != nil {
		log.ErrorC("Data not found S3", err, log.Data{"uri": s3uri.String})
		http.Error(w, "File not found in S3", http.StatusNotFound)
		return
	}
	data, err := json.Marshal(FileSizeMessage{FileSize: fileSize})
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
