package res

import (
	"io"
	"mime"
	"net/http"
	"path"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	b, err := Open(p)
	if err != nil {
		http.NotFound(rw, req)
		return
	}
	t := mime.TypeByExtension(path.Ext(p))
	rw.Header().Set("Content-Type", t)
	rw.WriteHeader(http.StatusOK)
	_, _ = io.Copy(rw, b)
}
