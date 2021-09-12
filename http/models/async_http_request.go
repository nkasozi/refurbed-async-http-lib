package models

import (
	"io"
	"net/http"
	"sync"
)

type AsyncHttpRequest struct {
	Url           string
	Body          io.Reader
	Headers       map[string]string
	ResultHandler func(resp *http.Response, err error)
	wg            sync.WaitGroup
	Method        string
}
