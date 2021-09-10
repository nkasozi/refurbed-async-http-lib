package async_http

import (
	"io"
	"net/http"
)

//go:generate mockgen -destination=mocks/http_request_sender_mock.go -package=mocks . HttpRequestSender
type HttpRequestSender interface {
	SendHttpRequest(method, url string, headers map[string]string, body io.Reader) (resp *http.Response, err error)
}

type DirectHttpRequestSender struct {
	client http.Client
}

func NewHttpRequestSender(httpClient http.Client) HttpRequestSender {
	return &DirectHttpRequestSender{client: httpClient}
}

func (sender DirectHttpRequestSender) SendHttpRequest(method, url string, headers map[string]string, body io.Reader) (resp *http.Response, err error) {
	client := http.Client{}
	request, err := http.NewRequest(method, url, body)

	for headerName, headerValue := range headers {
		request.Header.Set(headerName, headerValue)
	}

	if err != nil {
		return
	}

	resp, err = client.Do(request)

	return
}
