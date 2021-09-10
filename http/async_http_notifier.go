package async_http

import (
	"io"
	"net/http"
	"sync"
)

type AsyncHttpRequestSender interface {
	ShutDown()
	SendHttpRequestAsync(request AsyncHttpRequest) error
}

type AsyncHttpRequest struct {
	Method        string
	Url           string
	Body          io.Reader
	Headers       map[string]string
	ResultHandler func(resp *http.Response, err error)
	wg            sync.WaitGroup
}

type AsyncHttpNotifier struct {
	httpSender               HttpRequestSender
	pendingHttpRequestsQueue chan AsyncHttpRequest
	numberOfWorkerRoutines   int
	wg                       sync.WaitGroup
}

func (a *AsyncHttpNotifier) ShutDown() {
	close(a.pendingHttpRequestsQueue)
}

func NewAsyncHttpNotifier(sender HttpRequestSender, numberOfWorkerRoutines int) AsyncHttpRequestSender {

	result := &AsyncHttpNotifier{
		httpSender:               sender,
		pendingHttpRequestsQueue: make(chan AsyncHttpRequest, numberOfWorkerRoutines+1),
		numberOfWorkerRoutines:   numberOfWorkerRoutines,
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

func (a *AsyncHttpNotifier) SendHttpRequestAsync(request AsyncHttpRequest) (err error) {

	err = validateHttpRequest(request)

	if err != nil {
		return err
	}

	//queue up the request for sending
	a.pendingHttpRequestsQueue <- request

	return
}

func (a *AsyncHttpNotifier) processQueuedHttpRequestAsync(request AsyncHttpRequest) {
	resp, err := a.httpSender.SendHttpRequest(request.Method, request.Url, request.Headers, request.Body)

	//call the resp handler and pass the results
	request.ResultHandler(resp, err)
	return
}
