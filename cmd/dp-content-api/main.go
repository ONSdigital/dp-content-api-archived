package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-content-api/content"
	"github.com/ONSdigital/dp-content-api/s3"
	"github.com/ONSdigital/dp-content-api/utils"
	"github.com/ONSdigital/go-ns/log"
	_ "github.com/lib/pq"
)

var findMetaDataStatement *sql.Stmt
var findMetaDataWithFilterStatement *sql.Stmt
var findS3DataStatement *sql.Stmt
var parentStatement *sql.Stmt
var s3Client s3.S3Client
var taxonomyJSON []byte
var generatorURL string

func prepareSQLStatement(sql string, db *sql.DB) *sql.Stmt {
	statement, err := db.Prepare(sql)
	if err != nil {
		log.ErrorC("Error: Could not prepare statement on database", err, log.Data{"sql": sql})
		panic(err)
	}
	return statement
}

func main() {
	dbSource := utils.GetEnvironmentVariable("DB_ACCESS", "user=dp dbname=dp sslmode=disable")
	port := utils.GetEnvironmentVariable("PORT", "8082")
	generatorURL = utils.GetEnvironmentVariable("GENERATOR_URL", "localhost:8092")
	s3BucketName := utils.GetEnvironmentVariable("S3_BUCKET", "content")
	s3Endpoint := utils.GetEnvironmentVariable("S3_URL", "localhost:4000")
	accessKeyID := utils.GetEnvironmentVariable("S3_ACCESS_KEY", "1234")
	secretAccessKey := utils.GetEnvironmentVariable("S3_SECRET_ACCESS_KEY", "1234")
	taxonomyFile := utils.GetEnvironmentVariable("TAXONOMY_FILE", "static/taxonomy.json")
	log.Namespace = "dp-content-api"

	db, dbErr := sql.Open("postgres", dbSource)
	if dbErr != nil {
		log.ErrorC("Failed to connect to database", dbErr, log.Data{})
		panic(dbErr)
	}
	defer db.Close()
	findMetaDataSQL := "SELECT content FROM metadata WHERE uri = $1"
	findMetaDataWithFilterSQL := "SELECT json_build_object($1::text, content#>$2, 'uri', content->>'uri') FROM metadata WHERE uri = $3"
	findS3DataSQL := "SELECT s3 FROM s3data WHERE uri = $1"
	parentDataSQL := "SELECT uri, content->'description'->>'title', content->'type' FROM metadata WHERE uri = ANY($1) ORDER BY length(uri) ASC;"
	findMetaDataStatement = prepareSQLStatement(findMetaDataSQL, db)
	findMetaDataWithFilterStatement = prepareSQLStatement(findMetaDataWithFilterSQL, db)
	findS3DataStatement = prepareSQLStatement(findS3DataSQL, db)
	parentStatement = prepareSQLStatement(parentDataSQL, db)
	defer findMetaDataStatement.Close()
	defer findMetaDataWithFilterStatement.Close()
	defer findS3DataStatement.Close()
	defer parentStatement.Close()

	s3Client = s3.CreateClient(s3BucketName, s3Endpoint, accessKeyID, secretAccessKey, false)

	data, taxonomyErr := ioutil.ReadFile(taxonomyFile)
	if taxonomyErr != nil {
		log.ErrorC("Failed to load static parent data", taxonomyErr, log.Data{})
		panic(taxonomyErr)
	}
	taxonomyJSON = data

	log.Debug("Starting content api", log.Data{"port": port, "generator_url": generatorURL,
		"s3_bucket": s3BucketName, "s3_endpoint": s3Endpoint})
	// Babbage can use two different url types to call the content-api. One which
	// only contains the endpoint type and another which extends the type and includes
	// collectionID e.g /data/my-collection?param=list. As we don't need the collectionID
	// both endpoints uses the same function handler.
	http.HandleFunc("/data/", getData)
	http.HandleFunc("/data", getData)
	http.HandleFunc("/parents/", getParent)
	http.HandleFunc("/parents", getParent)
	http.HandleFunc("/resource/", getResource)
	http.HandleFunc("/resource", getResource)
	http.HandleFunc("/taxonomy/", getTaxonomy)
	http.HandleFunc("/taxonomy", getTaxonomy)
	http.HandleFunc("/generator/", generatorHandler)
	http.HandleFunc("/generator", generatorHandler)
	http.HandleFunc("/export/", exportHandler)
	http.HandleFunc("/export", exportHandler)
	serverErr := http.ListenAndServe(":"+port, nil)
	if serverErr != nil {
		log.ErrorC("Failed to start http server", serverErr, log.Data{})
		panic(serverErr)
	}
}

func getResource(rw http.ResponseWriter, rq *http.Request) {
	content.GetResource(rw, rq, findS3DataStatement, s3Client)
}

func getData(rw http.ResponseWriter, rq *http.Request) {
	content.GetData(rw, rq, findMetaDataStatement, findMetaDataWithFilterStatement)
}

func getParent(rw http.ResponseWriter, rq *http.Request) {
	log.Debug("Data", log.Data{"parms": rq.URL.Query()})
	content.GetParent(rw, rq, parentStatement)
}

func getTaxonomy(rw http.ResponseWriter, rq *http.Request) {
	content.StaticHandler(rw, rq, taxonomyJSON)
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	content.ExportHandler(w, r, generatorURL)
}

func generatorHandler(w http.ResponseWriter, r *http.Request) {
	content.GeneratorHandler(w, r, generatorURL)
}
