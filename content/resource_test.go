package content

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

type MockS3Client struct {
	throwError bool
}

func TestGetResourceReturns404FromDataStore(t *testing.T) {
	Convey("For a given uri, a resource is not found in a datastore", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnError(fmt.Errorf("Testing HTTP 500 code"))
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/parents?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetResource(w, r, stmt, &MockS3Client{throwError: false})
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestGetResourceReturns404FromS3(t *testing.T) {
	Convey("For a given uri, a resource is not found in a S3", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"uri"}).AddRow("/")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/parents?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetResource(w, r, stmt, &MockS3Client{throwError: true})
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestGetResourceReturnsContentFromS3(t *testing.T) {
	Convey("For a given uri, a resource is not found in a S3", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"uri"}).AddRow("/")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/parents?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetResource(w, r, stmt, &MockS3Client{throwError: false})
		So(string(w.Body.Bytes()), ShouldEqual, "test-data")
	})
}

func (s3 *MockS3Client) GetBucket() string {
	return "mock-bucket"
}

func (s3 *MockS3Client) GetObject(uri string) ([]byte, error) {
	if s3.throwError {
		return nil, fmt.Errorf("Mock S3 error")
	}
	return []byte("test-data"), nil
}
