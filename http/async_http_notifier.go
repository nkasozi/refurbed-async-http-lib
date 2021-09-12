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

type asyncHttpRequestsSender struct {
	httpSender               HttpRequestSender
	pendingHttpRequestsQueue chan models.AsyncHttpRequest
	numberOfWorkerRoutines   int
	wg                       sync.WaitGroup
}

func (a *asyncHttpRequestsSender) ShutDown() {
	close(a.pendingHttpRequestsQueue)
}

func NewAsyncHttpRequestsSender(sender HttpRequestSender, numberOfWorkerRoutines int) AsyncHttpRequestSender {

	result := &asyncHttpRequestsSender{
		httpSender:               sender,
		pendingHttpRequestsQueue: make(chan models.AsyncHttpRequest, numberOfWorkerRoutines+1),
		numberOfWorkerRoutines:   numberOfWorkerRoutines,
	}

	result.startChannelProcessingUsingMultipleGoRoutines()

	return result
}

func (a *asyncHttpRequestsSender) startChannelProcessingUsingMultipleGoRoutines() {
	go func() {
		// This starts x number of goroutines that wait for something to do
		a.wg.Add(a.numberOfWorkerRoutines)

		for i := 0; i < a.numberOfWorkerRoutines; i++ {
			go func() {
				for {
					pendingRequest, isOpen := <-a.pendingHttpRequestsQueue

					// if there is nothing to do and the channel has been closed then end the goroutine
					if !isOpen {
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

func (a *asyncHttpRequestsSender) SendHttpRequestAsync(request models.AsyncHttpRequest) (err error) {

	err = validateHttpRequest(request)

	if err != nil {
		return err
	}

	//queue up the request for sending
	a.pendingHttpRequestsQueue <- request

	return
}

func (a *asyncHttpRequestsSender) processQueuedHttpRequestAsync(request models.AsyncHttpRequest) {
	resp, err := a.httpSender.SendHttpRequest(request.Method, request.Url, request.Headers, request.Body)

	//call the resp handler and pass the results
	request.ResultHandler(resp, err)
	return
}
