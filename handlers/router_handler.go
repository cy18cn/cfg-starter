package handlers

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"time"
)

type parseFormHandler struct {
	next http.Handler
	log  *zap.Logger
}

func (self *parseFormHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now().UnixNano() / 1e6
	err := self.parseRequest(req)
	if err != nil {
		body, _ := readBody(req)
		self.log.Sugar().Errorf("bad request URL: %s, err: %v, body: %s",
			req.RequestURI,
			err,
			body)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {
		self.next.ServeHTTP(w, req)
		self.log.Sugar().Infof("done to handle request, it takes: %d ms", time.Now().UnixNano()/1e6-startTime)
	}

	return
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
			err = parsePost(req)
			if err != nil {
				return err
			}
			if req.Form == nil {
				if len(req.PostForm) > 0 {
					req.Form = make(url.Values)
					copyValues(req.Form, req.PostForm)
				}
				var newValues url.Values
				if req.URL != nil {
					var e error
					newValues, e = url.ParseQuery(req.URL.RawQuery)
					if err == nil {
						err = e
					}
				}
				if newValues == nil {
					newValues = make(url.Values)
				}
				if req.Form == nil {
					req.Form = newValues
				} else {
					copyValues(req.Form, newValues)
				}
			}
			return nil
		case contentType == "multipart/form-data":
			return req.ParseMultipartForm(maxMemory)
		}
	}

	return req.ParseForm()
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		for _, value := range vs {
			dst.Add(k, value)
		}
	}
}

func parsePost(req *http.Request) (err error) {
	if req.PostForm == nil {
		req.PostForm = make(url.Values)
		var body string
		body, err = readBody(req)
		if err != nil {
			return
		}
		req.PostForm["body"] = []string{body}
	}
	return
}

func readBody(req *http.Request) (body string, err error) {
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
