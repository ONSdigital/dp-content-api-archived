package content

import (
	"net/http"
)

func StaticHandler(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Write(data)
}
