package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetFileSize404FromDataStore(t *testing.T) {
	Convey("For a given uri, a resource is not found in a datastore", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnError(fmt.Errorf("Testing HTTP 404 code"))
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/filesize?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetFileSize(w, r, stmt, &MockS3Client{throwError: false})
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestGetFileSize404FromS3(t *testing.T) {
	Convey("For a given uri, a resource is not found in a datastore", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"uri"}).AddRow("s3:://content/my/data.xls")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnRows(row)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/filesize?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetFileSize(w, r, stmt, &MockS3Client{throwError: true})
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestGetFileSizeContent(t *testing.T) {
	Convey("For a given uri, a resource is not found in a datastore", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"uri"}).AddRow("s3:://content/my/data.xls")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnRows(row)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/filesize?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetFileSize(w, r, stmt, &MockS3Client{throwError: false})
		So(w.Code, ShouldEqual, http.StatusOK)
		var message FileSizeMessage
		json.Unmarshal(w.Body.Bytes(), &message)
		So(message.FileSize, ShouldEqual, 0)
	})
}
