package handlers

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
)

type parseFormHandler struct {
	next http.Handler
	log  *zap.Logger
}

func (self *parseFormHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := self.parseRequest(req)
	if err == nil {
		self.next.ServeHTTP(w, req)
		return
	}

	body, _ := self.readBody(req)
	self.log.Sugar().Errorf("bad request URL: %s, err: %v, body: %s",
		req.RequestURI,
		err,
		body)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

const maxMemory = 10 * 1024 * 1024

func (self *parseFormHandler) parseRequest(req *http.Request) error {
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		contentType, err := getRequestContentType(req)
		if err != nil {
			return err
		}

		switch {
		case contentType == "application/json":
			var body string
			body, err = self.readBody(req)
			if err == nil {
				req.Form["body"] = []string{body}
			}
			return nil
		case contentType == "multipart/form-data":
			err = req.ParseMultipartForm(maxMemory)
			return nil
		}
	}

	return req.ParseForm()
}

func (self *parseFormHandler) readBody(req *http.Request) (body string, err error) {
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

func getRequestContentType(r *http.Request) (string, error) {
	contentType := r.Header.Get("Content-Type")
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var err error
	contentType, _, err = mime.ParseMediaType(contentType)
	return contentType, err
}
