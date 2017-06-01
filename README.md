## Content API

An API used to get data from published meta and static files from different
resource.

### Environment variables
* `S3_URL` defaults to `s3.amazonaws.com`
* `S3_ACCESS_KEY` defaults to "1234"
* `S3_SECRET_ACCESS_KEY` defaults to "1234"
* `S3_BUCKET` defaults to "content"
* `DB_ACCESS` defaults to "user=dp dbname=dp sslmode=disable"
* `GENERATOR_URL` defaults to "http://localhost:8092"
* `PORT` defaults to "8082"
* `TAXONOMY_FILE` defaults to "static/taxonomy.json"
* `PARENT_FILE` defaults to "static/parent.json"
* `HEALTHCHECK_ENDPOINT` defaults to "/healthcheck"

### Interfaces

#### Data
Route : Get /data

Parameters
* uri : location of meta document in the datastore

Errors
* 404 : meta document not found

#### Resource
Route : Get /resource

Parameters
* uri : location of meta document in the S3

Errors
* 404 : S3 file not found

#### Parent
Route : Get /parent

Parameters
* uri : A uri to create bread crumb trail.


#### Taxonomy
Route : Get /taxonomy

Hard coded value. See static-taxonomy.json


#### File Size
Route : Get /filesize

Parameters
* uri : location of static file in the S3

Errors
* 404 : S3 file not found

