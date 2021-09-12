package async_http

import (
	"errors"
	"github.com/nkasozi/refurbed-async-http-lib/http/models"
)

func validateHttpRequest(request models.AsyncHttpRequest) (err error) {
	if request.ResultHandler == nil {
		return errors.New("please supply a result handler function")
	}
	if request.Url == "" {
		return errors.New("please supply a URl to which to send requests to")
	}
	if request.Method == "" {
		return errors.New("please supply an HTTP method to use when sending the request")
	}
	return nil
}
