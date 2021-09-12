package async_http

import (
	"github.com/nkasozi/refurbed-async-http-lib/http/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidateHttpRequest(t *testing.T) {
	t.Run("given a valid request, returns no error", func(t *testing.T) {
		request := models.AsyncHttpRequest{
			Url:     "http://test-url",
			Body:    nil,
			Headers: nil,
			ResultHandler: func(resp *http.Response, err error) {

			},
			Method: http.MethodGet,
		}

		validationErr := validateHttpRequest(request)
		assert.NoError(t, validationErr)
	})

	t.Run("given a request with no method, returns an error", func(t *testing.T) {
		request := models.AsyncHttpRequest{
			Body:    nil,
			Headers: nil,
			ResultHandler: func(resp *http.Response, err error) {

			},
			Method: "",
		}

		validationErr := validateHttpRequest(request)
		assert.Error(t, validationErr)
	})

	t.Run("given a request with no handler, returns an error", func(t *testing.T) {
		request := models.AsyncHttpRequest{
			Body:          nil,
			Headers:       nil,
			ResultHandler: nil,
			Method:        http.MethodGet,
		}

		validationErr := validateHttpRequest(request)
		assert.Error(t, validationErr)
	})
}
