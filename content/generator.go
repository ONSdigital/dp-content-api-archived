package content

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-publish-pipeline/utils"
)

func ExportHandler(w http.ResponseWriter, r *http.Request, generatorURL string) {
	r.ParseForm()
	res, _ := http.PostForm(generatorURL+r.URL.String(), r.PostForm)
	body, _ := ioutil.ReadAll(res.Body)
	if strings.Contains(r.PostFormValue("format"), "csv") {
		utils.SetCSVContentHeader(w)
	} else {
		utils.SetXLSContentHeader(w)
	}
	w.Write(body)
}

func GeneratorHandler(w http.ResponseWriter, r *http.Request, generatorURL string) {
	res, _ := http.Get(generatorURL + r.URL.String())
	body, _ := ioutil.ReadAll(res.Body)
	if strings.Contains(r.URL.Query().Get("format"), "csv") {
		utils.SetCSVContentHeader(w)
	} else {
		utils.SetXLSContentHeader(w)
	}
	w.Write(body)
}
