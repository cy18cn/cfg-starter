package handlers

import (
	"errors"
	"github.com/cy18cn/zlog"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

type ParseFormHandler struct {
	next http.Handler
}

func (self *ParseFormHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := parseRequest(req)
	if err == nil {
		self.next.ServeHTTP(w, req)
	}

	var body string
	body, err = readBody(req)
	zlog.Errorf("bad request URL: %s, err: %v, body: %s",
		req.RequestURI,
		err,
		body)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

const maxMemory = 10 * 1024 * 1024

func parseRequest(req *http.Request) (err error) {
	contentType := req.Header.Get("Content-Type")
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	contentType, _, err = mime.ParseMediaType(contentType)
	switch {
	case contentType == "application/json":
		var body string
		body, err = readBody(req)
		if err == nil {
			req.Form["body"] = []string{body}
		}
	case contentType == "application/x-www-form-urlencoded":
		err = req.ParseForm()
	case contentType == "multipart/form-data":
		err = req.ParseMultipartForm(maxMemory)
	default:
		err = errors.New("unsupported content type")
	}

	return
}

func readBody(req *http.Request) (body string, err error) {
	if req.Form == nil {
		var reader io.Reader = req.Body
		maxFormSize := int64(10 << 20) // 10 MB is a lot of text.
		b, e := ioutil.ReadAll(reader)
		if e != nil {
			if err == nil {
				err = e
			}
			return
		}

		if int64(len(b)) > maxFormSize {
			err = errors.New("http: POST too large")
			return
		}

		body = string(b)
	}

	return
}
