package hooks

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type HooksClient struct {
	OnUnknownResponseCode func(response *http.Response, request *http.Request) string
}

func DevHook() HooksClient {
	return HooksClient{
		OnUnknownResponseCode: func(response *http.Response, request *http.Request) string {
			var httpRequestDumpMessage string
			httpRequestDump, err := httputil.DumpRequest(request, true)
			if err != nil {
				httpRequestDumpMessage = fmt.Sprintf("could not dump request (%v)", err.Error())
			} else {
				httpRequestDumpMessage = string(httpRequestDump)
			}

			var httpResponseDumpMessage string
			httpResponseDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				httpResponseDumpMessage = fmt.Sprintf("could not dump response (%v)", err.Error())
			} else {
				httpResponseDumpMessage = string(httpResponseDump)
			}

			message := fmt.Sprintf("unknown response status code %d", response.StatusCode)
			if len(httpRequestDump) != 0 {
				message = message + "\n HTTP Request: \n '" + string(httpResponseDumpMessage) + "' \n"
			}
			if len(httpResponseDump) != 0 {
				message = message + "HTTP Response: \n '" + string(httpRequestDumpMessage) + "'"
			}
			return message
		},
	}
}
