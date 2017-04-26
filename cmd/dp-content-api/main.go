package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"github.com/ONSdigital/dp-content-api/content"
	"github.com/ONSdigital/dp-content-api/utils"
	"github.com/ONSdigital/go-ns/log"
	_ "github.com/lib/pq"
)

var findMetaDataStatement *sql.Stmt
var findMetaDataWithFilterStatement *sql.Stmt
var findS3DataStatement *sql.Stmt
var parentJSON []byte
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
	taxonomyFile := utils.GetEnvironmentVariable("TAXONOMY_FILE", "static/taxonomy.json")
	parentFile := utils.GetEnvironmentVariable("PARENT_FILE", "static/parent.json")
	log.Namespace = "dp-content-api"

	db, dbErr := sql.Open("postgres", dbSource)
	if dbErr != nil {
		log.ErrorC("Failed to connect to database", dbErr, log.Data{})
		panic(dbErr)
	}
	defer db.Close()
	findMetaDataSQL := "SELECT content FROM metadata WHERE uri = $1"
	findMetaDataWithFilterSQL := "SELECT json_build_object($1::text, content->'description'->>$2) FROM metadata WHERE uri = $3"
	findS3DataSQL := "SELECT s3 FROM s3data WHERE uri = $1"
	findMetaDataStatement = prepareSQLStatement(findMetaDataSQL, db)
	findMetaDataWithFilterStatement = prepareSQLStatement(findMetaDataWithFilterSQL, db)
	findS3DataStatement = prepareSQLStatement(findS3DataSQL, db)
	defer findMetaDataStatement.Close()
	defer findS3DataStatement.Close()

	data, parentErr := ioutil.ReadFile(parentFile)
	if parentErr != nil {
		log.ErrorC("Failed to load static parent data", parentErr, log.Data{})
		panic(parentErr)
	}
	parentJSON = data

	data, taxonomyErr := ioutil.ReadFile(taxonomyFile)
	if taxonomyErr != nil {
		log.ErrorC("Failed to load static parent data", taxonomyErr, log.Data{})
		panic(taxonomyErr)
	}
	taxonomyJSON = data

	log.Debug("Starting content api", log.Data{"port": port, "generator_url": generatorURL})
	// Babbage can use two different url types to call the content-api. One which
	// only contains the endpoint type and another which extends the type and includes
	// collectionID e.g /data/my-collection?param=list. As we don't need the collectionID
	// both endpoints uses the same function handler.
	http.HandleFunc("/data/", getData)
	http.HandleFunc("/data", getData)
	http.HandleFunc("/parent/", getParent)
	http.HandleFunc("/parent", getParent)
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
	content.GetResource(rw, rq, findS3DataStatement)
}

func getData(rw http.ResponseWriter, rq *http.Request) {
	content.GetData(rw, rq, findMetaDataStatement, findMetaDataWithFilterStatement)
}

func getParent(rw http.ResponseWriter, rq *http.Request) {
	content.StaticHandler(rw, rq, parentJSON)
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
