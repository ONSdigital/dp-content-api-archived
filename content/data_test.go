package content

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDataReturns404(t *testing.T) {
	Convey("From a URI parameter, a ONS page is return", t, func() {
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
		r, err := http.NewRequest("GET", "http://localhost/data?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetData(w, r, stmt, stmt)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

func TestGetDataReturnsPageContent(t *testing.T) {
	Convey("From a URI parameter, a ONS page is return", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"content"}).AddRow("page data")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/data?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetData(w, r, stmt, stmt)
		So(string(w.Body.Bytes()), ShouldEqual, "page data")
	})
}

func TestGetDataReturnsFilteredData(t *testing.T) {
	Convey("From a URI parameter, a ONS page is return", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"content"}).AddRow("page data")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		defer stmt.Close()
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/data?uri=/&filter=description", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetData(w, r, nil, stmt)
		So(string(w.Body.Bytes()), ShouldEqual, "page data")
	})
}
