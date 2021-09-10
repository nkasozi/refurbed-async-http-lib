package async_http

import (
	"github.com/golang/mock/gomock"
	"github.com/nkasozi/refurbed-async-http-lib/http/mocks"
	"github.com/nkasozi/refurbed-async-http-lib/http/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestShutDown(t *testing.T) {
	t.Run("given that we shutdown successfully, channels should be closed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		url := "test-url"
		httpReqSenderMock := mocks.NewMockHttpRequestSender(ctrl)

		notifier := NewAsyncHttpNotifier(url, httpReqSenderMock, 1)

		asyncNotifier, ok := notifier.(*AsyncHttpNotifier)

		if !ok {
			t.Fatalf("notifier is of unexpected type")
			return
		}

		notifier.ShutDown()

		_, isOpen := <-asyncNotifier.pendingHttpRequestsQueue

		if isOpen {
			t.Fatalf("pending requests channel should be in closed state after shutdown")
			return
		}
	})
}

func TestSendHttpRequestAsync(t *testing.T) {
	t.Run("given a valid request, passes validation, returns no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		url := "test-url"
		httpReqSenderMock := mocks.NewMockHttpRequestSender(ctrl)

		notifier := NewAsyncHttpNotifier(url, httpReqSenderMock, 1)

		request := models.AsyncHttpRequest{
			Body:    nil,
			Headers: nil,
			ResultHandler: func(resp *http.Response, err error) {

			},
			Method: http.MethodGet,
		}

		httpReqSenderMock.EXPECT().
			SendHttpRequest(request.Method, url, request.Headers, request.Body).
			Return(&http.Response{}, nil).
			AnyTimes()

		actualResult := notifier.SendHttpRequestAsync(request)
		assert.NoError(t, actualResult)
	})

	t.Run("given a request missing some required fields, returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		url := "test-url"
		httpReqSenderMock := mocks.NewMockHttpRequestSender(ctrl)

		notifier := NewAsyncHttpNotifier(url, httpReqSenderMock, 1)

		request := models.AsyncHttpRequest{}

		actualResult := notifier.SendHttpRequestAsync(request)
		assert.Error(t, actualResult)
	})

	t.Run("given a valid request, passes validation, calls result handler with result, returns no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		url := "test-url"
		httpReqSenderMock := mocks.NewMockHttpRequestSender(ctrl)

		notifier := NewAsyncHttpNotifier(url, httpReqSenderMock, 1)

		isCalled := false

		request := models.AsyncHttpRequest{
			Body:    nil,
			Headers: nil,
			ResultHandler: func(resp *http.Response, err error) {
				isCalled = true
			},
			Method: http.MethodGet,
		}

		httpReqSenderMock.EXPECT().
			SendHttpRequest(request.Method, url, request.Headers, request.Body).
			Return(&http.Response{}, nil).
			AnyTimes()

		actualResult := notifier.SendHttpRequestAsync(request)
		assert.NoError(t, actualResult)
		time.Sleep(2 * time.Second)
		assert.True(t, isCalled)
	})
}
