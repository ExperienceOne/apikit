package basket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

type basketServiceClient struct {
	baseURL    string
	hooks      HooksClient
	ctx        context.Context
	httpClient *httpClientWrapper
	xmlMatcher *regexp.Regexp
}

func (client *basketServiceClient) PostItem(request *PostItemRequest) (PostItemResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/item"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	jsonData := new(bytes.Buffer)
	encodeErr := json.NewEncoder(jsonData).Encode(&request.Item)
	if encodeErr != nil {
		return nil, encodeErr
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, jsonData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentTypeApplicationJson}
	// set all headers from client context
	err := setRequestHeadersFromContext(httpContext, httpRequest.Header)
	if err != nil {
		return nil, err
	}
	if len(httpRequest.Header["accept"]) == 0 && len(httpRequest.Header["Accept"]) == 0 {
		httpRequest.Header["Accept"] = []string{"application/json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode == 200 {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			httpResponse.Body.Close()
			response := new(PostItem200Response)
			return response, nil
		}
		httpResponse.Body.Close()
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == 500 {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			httpResponse.Body.Close()
			response := new(PostItem500Response)
			return response, nil
		}
		httpResponse.Body.Close()
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		httpResponse.Body.Close()
		return nil, errors.New(message)
	}
	httpResponse.Body.Close()
	return nil, newUnknownResponseError(httpResponse.StatusCode)
}

func (client *basketServiceClient) GetItems(request *GetItemsRequest) (GetItemsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/items"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	// set all headers from client context
	err := setRequestHeadersFromContext(httpContext, httpRequest.Header)
	if err != nil {
		return nil, err
	}
	if len(httpRequest.Header["accept"]) == 0 && len(httpRequest.Header["Accept"]) == 0 {
		httpRequest.Header["Accept"] = []string{"application/json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode == 200 {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetItems200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			httpResponse.Body.Close()
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			httpResponse.Body.Close()
			response := new(GetItems200Response)
			return response, nil
		}
		httpResponse.Body.Close()
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		httpResponse.Body.Close()
		return nil, errors.New(message)
	}
	httpResponse.Body.Close()
	return nil, newUnknownResponseError(httpResponse.StatusCode)
}

type BasketServiceClient interface {
	PostItem(request *PostItemRequest) (PostItemResponse, error)
	GetItems(request *GetItemsRequest) (GetItemsResponse, error)
}

func NewBasketServiceClient(httpClient *http.Client, baseUrl string, options Opts) BasketServiceClient {
	return &basketServiceClient{httpClient: newHttpClientWrapper(httpClient, baseUrl), baseURL: baseUrl, hooks: options.Hooks, ctx: options.Ctx, xmlMatcher: regexp.MustCompile("^application\\/(.+)xml$")}
}
