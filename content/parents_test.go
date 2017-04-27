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

func TestGetParentContent(t *testing.T) {
	Convey("For a given uri, its bread crumb trail is returned", t, func() {
		uri := "/peoplepopulationandcommunity/populationandmigration/populationestimates/bulletins/annualmidyearpopulationestimates/2015-06-25"
		results := getParentPages(uri)
		So(results, ShouldNotBeNil)
		So(results, ShouldContain, "/peoplepopulationandcommunity?lang=en")
		So(results, ShouldContain, "/peoplepopulationandcommunity/populationandmigration?lang=en")
		So(results, ShouldContain, "/peoplepopulationandcommunity/populationandmigration/populationestimates?lang=en")
	})
}

func TestGetParentReturnsHTTPStatus500Part1(t *testing.T) {
	Convey("From a HTTP GET with a invalid SQL connection, an error code of 500 is returned", t, func() {
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
		GetParent(w, r, stmt)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)

	})
}

func TestGetParentReturnsJsonMessage(t *testing.T) {
	Convey("From a HTTP GET parent json message is returned", t, func() {
		db, mock, err := sqlmock.New()
		So(err, ShouldBeNil)
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"uri", "title", "type"}).
			AddRow("/", "home", "test").
			AddRow("/about", "about", "page")
		mock.ExpectPrepare("SELECT").
			ExpectQuery().
			WillReturnRows(rows)
		db.Begin()
		defer db.Close()
		stmt, err := db.Prepare("SELECT")
		So(err, ShouldBeNil)
		r, err := http.NewRequest("GET", "http://localhost/data?uri=/", nil)
		So(err, ShouldBeNil)
		w := httptest.NewRecorder()
		GetParent(w, r, stmt)
		var nodes []Node
		So(json.Unmarshal(w.Body.Bytes(), &nodes), ShouldBeNil)
		So(nodes, ShouldContain, Node{URI: "/", Description: NodeDescription{Title: "home"}, Type: "test"})
		So(nodes, ShouldContain, Node{URI: "/about", Description: NodeDescription{Title: "about"}, Type: "page"})
	})
}
