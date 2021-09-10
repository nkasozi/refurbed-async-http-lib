package async_http

import (
	"github.com/nkasozi/refurbed-async-http-lib/http/models"
	"sync"
)

//go:generate mockgen -destination=mocks/async_http_notifier_mock.go -package=mocks . AsyncHttpRequestSender
type AsyncHttpRequestSender interface {
	ShutDown()
	SendHttpRequestAsync(request models.AsyncHttpRequest) error
}

type AsyncHttpNotifier struct {
	httpSender               HttpRequestSender
	pendingHttpRequestsQueue chan models.AsyncHttpRequest
	numberOfWorkerRoutines   int
	wg                       sync.WaitGroup
	url                      string
}

func (a *AsyncHttpNotifier) ShutDown() {
	close(a.pendingHttpRequestsQueue)
}

func NewAsyncHttpNotifier(url string, sender HttpRequestSender, numberOfWorkerRoutines int) AsyncHttpRequestSender {

	result := &AsyncHttpNotifier{
		httpSender:               sender,
		pendingHttpRequestsQueue: make(chan models.AsyncHttpRequest, numberOfWorkerRoutines+1),
		numberOfWorkerRoutines:   numberOfWorkerRoutines,
		url:                      url,
	}

	result.startChannelProcessingUsingMultipleGoRoutines()

	return result
}

func (a *AsyncHttpNotifier) startChannelProcessingUsingMultipleGoRoutines() {
	go func() {
		// This starts x number of goroutines that wait for something to do
		a.wg.Add(a.numberOfWorkerRoutines)

		for i := 0; i < a.numberOfWorkerRoutines; i++ {
			go func() {
				for {
					pendingRequest, ok := <-a.pendingHttpRequestsQueue

					// if there is nothing to do and the channel has been closed then end the goroutine
					if !ok {
						a.wg.Done()
						return
					}

					// process the pending request
					a.processQueuedHttpRequestAsync(pendingRequest)
				}
			}()
		}

		// Wait for the threads to finish
		a.wg.Wait()
	}()
}

func (a *AsyncHttpNotifier) SendHttpRequestAsync(request models.AsyncHttpRequest) (err error) {

	err = validateHttpRequest(request)

	if err != nil {
		return err
	}

	//queue up the request for sending
	a.pendingHttpRequestsQueue <- request

	return
}

func (a *AsyncHttpNotifier) processQueuedHttpRequestAsync(request models.AsyncHttpRequest) {
	resp, err := a.httpSender.SendHttpRequest(request.Method, a.url, request.Headers, request.Body)

	//call the resp handler and pass the results
	request.ResultHandler(resp, err)
	return
}
