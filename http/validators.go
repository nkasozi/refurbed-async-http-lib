package async_http

import "errors"

func validateHttpRequest(request AsyncHttpRequest) (err error) {
	if request.ResultHandler == nil {
		return errors.New("please supply a result handler function")
	}
	if request.Url == "" {
		return errors.New("please supply a URL to which to send the request to")
	}
	if request.Method == "" {
		return errors.New("please supply an HTTP method to use when sending the request")
	}
	return nil
}
