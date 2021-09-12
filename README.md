# refurbed-async-http-lib

Its a go library that helps you make async calls to an end point and also allows you to handle the responses asynchronously later

## Installation
```
go get -u github.com/nkasozi/refurbed-async-http-lib
```

## Usage

```
import(
  async_http "github.com/nkasozi/refurbed-async-http-lib/http"
	"github.com/nkasozi/refurbed-async-http-lib/http/models"
  )

    //build and send the http message
		request := models.AsyncHttpRequest{
			Url:    notifyUrl,
			Body:   strings.NewReader("Hello Worl"),
			Method: http.MethodPost,
			ResultHandler: func(resp *http.Response, err error) {
        
        //handle send error
				if err != nil {
					fmt.Printf("Request Failed because of error: [%v]\n", err)
					return
				}

        //handle an error response
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Request Failed With Status Code: [%v]\n", resp.StatusCode)
					return
				}

        //handle success response
				fmt.Printf("Recieved Success Response for Request\n")
			},
		}

		asyncHttpRequestSender.SendHttpRequestAsync(request)
```

