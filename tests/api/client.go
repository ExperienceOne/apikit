package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type VisAdminClient interface {
	GetClientsMethod
	DeleteClientMethod
	GetClientMethod
	CreateOrUpdateClientMethod
	GetViewsSetsMethod
	DeleteViewsSetMethod
	GetViewsSetMethod
	ActivateViewsSetMethod
	CreateOrUpdateViewsSetMethod
	ShowVehicleInViewMethod
	GetPermissionsMethod
	DestroySessionMethod
	GetUserInfoMethod
	CreateSessionMethod
	GetUsersMethod
	DeleteUserMethod
	GetUserMethod
	CreateOrUpdateUserMethod
	GetBookingMethod
	GetBookingsMethod
	ListModelsMethod
	GetClassesMethod
	CodeMethod
	DeleteCustomerSessionMethod
	CreateCustomerSessionMethod
	DownloadNestedFileMethod
	DownloadImageMethod
	ListElementsMethod
	FileUploadMethod
	DownloadFileMethod
	FindByTagsMethod
	GenericFileDownloadMethod
	GetRentalMethod
	GetShoesMethod
	PostUploadMethod
}
type GetClientsMethod interface {
	GetClients(request *GetClientsRequest) (GetClientsResponse, error)
}
type DeleteClientMethod interface {
	DeleteClient(request *DeleteClientRequest) (DeleteClientResponse, error)
}
type GetClientMethod interface {
	GetClient(request *GetClientRequest) (GetClientResponse, error)
}
type CreateOrUpdateClientMethod interface {
	CreateOrUpdateClient(request *CreateOrUpdateClientRequest) (CreateOrUpdateClientResponse, error)
}
type GetViewsSetsMethod interface {
	GetViewsSets(request *GetViewsSetsRequest) (GetViewsSetsResponse, error)
}
type DeleteViewsSetMethod interface {
	DeleteViewsSet(request *DeleteViewsSetRequest) (DeleteViewsSetResponse, error)
}
type GetViewsSetMethod interface {
	GetViewsSet(request *GetViewsSetRequest) (GetViewsSetResponse, error)
}
type ActivateViewsSetMethod interface {
	ActivateViewsSet(request *ActivateViewsSetRequest) (ActivateViewsSetResponse, error)
}
type CreateOrUpdateViewsSetMethod interface {
	CreateOrUpdateViewsSet(request *CreateOrUpdateViewsSetRequest) (CreateOrUpdateViewsSetResponse, error)
}
type ShowVehicleInViewMethod interface {
	ShowVehicleInView(request *ShowVehicleInViewRequest) (ShowVehicleInViewResponse, error)
}
type GetPermissionsMethod interface {
	GetPermissions(request *GetPermissionsRequest) (GetPermissionsResponse, error)
}
type DestroySessionMethod interface {
	DestroySession(request *DestroySessionRequest) (DestroySessionResponse, error)
}
type GetUserInfoMethod interface {
	GetUserInfo(request *GetUserInfoRequest) (GetUserInfoResponse, error)
}
type CreateSessionMethod interface {
	CreateSession(request *CreateSessionRequest) (CreateSessionResponse, error)
}
type GetUsersMethod interface {
	GetUsers(request *GetUsersRequest) (GetUsersResponse, error)
}
type DeleteUserMethod interface {
	DeleteUser(request *DeleteUserRequest) (DeleteUserResponse, error)
}
type GetUserMethod interface {
	GetUser(request *GetUserRequest) (GetUserResponse, error)
}
type CreateOrUpdateUserMethod interface {
	CreateOrUpdateUser(request *CreateOrUpdateUserRequest) (CreateOrUpdateUserResponse, error)
}
type GetBookingMethod interface {
	GetBooking(request *GetBookingRequest) (GetBookingResponse, error)
}
type GetBookingsMethod interface {
	GetBookings(request *GetBookingsRequest) (GetBookingsResponse, error)
}
type ListModelsMethod interface {
	ListModels(request *ListModelsRequest) (ListModelsResponse, error)
}
type GetClassesMethod interface {
	GetClasses(request *GetClassesRequest) (GetClassesResponse, error)
}
type CodeMethod interface {
	Code(request *CodeRequest) (CodeResponse, error)
}
type DeleteCustomerSessionMethod interface {
	DeleteCustomerSession(request *DeleteCustomerSessionRequest) (DeleteCustomerSessionResponse, error)
}
type CreateCustomerSessionMethod interface {
	CreateCustomerSession(request *CreateCustomerSessionRequest) (CreateCustomerSessionResponse, error)
}
type DownloadNestedFileMethod interface {
	DownloadNestedFile(request *DownloadNestedFileRequest) (DownloadNestedFileResponse, error)
}
type DownloadImageMethod interface {
	DownloadImage(request *DownloadImageRequest) (DownloadImageResponse, error)
}
type ListElementsMethod interface {
	ListElements(request *ListElementsRequest) (ListElementsResponse, error)
}
type FileUploadMethod interface {
	FileUpload(request *FileUploadRequest) (FileUploadResponse, error)
}
type DownloadFileMethod interface {
	DownloadFile(request *DownloadFileRequest) (DownloadFileResponse, error)
}
type FindByTagsMethod interface {
	FindByTags(request *FindByTagsRequest) (FindByTagsResponse, error)
}
type GenericFileDownloadMethod interface {
	GenericFileDownload(request *GenericFileDownloadRequest) (GenericFileDownloadResponse, error)
}
type GetRentalMethod interface {
	GetRental(request *GetRentalRequest) (GetRentalResponse, error)
}
type GetShoesMethod interface {
	GetShoes(request *GetShoesRequest) (GetShoesResponse, error)
}
type PostUploadMethod interface {
	PostUpload(request *PostUploadRequest) (PostUploadResponse, error)
}

func NewVisAdminClient(httpClient *http.Client, baseUrl string, options Opts) VisAdminClient {
	return &visAdminClient{httpClient: newHttpClientWrapper(httpClient, baseUrl), baseURL: baseUrl, hooks: options.Hooks, ctx: options.Ctx, xmlMatcher: regexp.MustCompile("^application\\/(.+)xml$")}
}

type visAdminClient struct {
	baseURL    string
	hooks      HooksClient
	ctx        context.Context
	httpClient *httpClientWrapper
	xmlMatcher *regexp.Regexp
}

func (client *visAdminClient) GetClients(request *GetClientsRequest) (GetClientsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetClients200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetClients200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNoContent {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetClients204Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetClients403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) DeleteClient(request *DeleteClientRequest) (DeleteClientResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}"
	method := "DELETE"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteClient200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteClient403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteClient404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetClient(request *GetClientRequest) (GetClientResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetClient200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetClient200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetClient403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetClient404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) CreateOrUpdateClient(request *CreateOrUpdateClientRequest) (CreateOrUpdateClientResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}"
	method := "PUT"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	jsonData := new(bytes.Buffer)
	encodeErr := json.NewEncoder(jsonData).Encode(&request.Body)
	if encodeErr != nil {
		return nil, encodeErr
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, jsonData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentTypeApplicationJson}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateClient200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusCreated {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateClient201Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateClient400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateClient403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusMethodNotAllowed {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateClient405Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetViewsSets(request *GetViewsSetsRequest) (GetViewsSetsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}/views"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetViewsSets200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetViewsSets200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetViewsSets403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) DeleteViewsSet(request *DeleteViewsSetRequest) (DeleteViewsSetResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}/views/{viewsId}"
	method := "DELETE"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	endpoint = strings.Replace(endpoint, "{viewsId}", url.QueryEscape(toString(request.ViewsId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteViewsSet200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteViewsSet403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteViewsSet404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetViewsSet(request *GetViewsSetRequest) (GetViewsSetResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}/views/{viewsId}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	endpoint = strings.Replace(endpoint, "{viewsId}", url.QueryEscape(toString(request.ViewsId)), 1)
	query := make(url.Values)
	query.Add("page", toString(request.Page))
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetViewsSet200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetViewsSet200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetViewsSet403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetViewsSet404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Make this viewset the active one for the client.
func (client *visAdminClient) ActivateViewsSet(request *ActivateViewsSetRequest) (ActivateViewsSetResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}/views/{viewsId}"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	endpoint = strings.Replace(endpoint, "{viewsId}", url.QueryEscape(toString(request.ViewsId)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ActivateViewsSet200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ActivateViewsSet403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ActivateViewsSet404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) CreateOrUpdateViewsSet(request *CreateOrUpdateViewsSetRequest) (CreateOrUpdateViewsSetResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}/views/{viewsId}"
	method := "PUT"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	endpoint = strings.Replace(endpoint, "{viewsId}", url.QueryEscape(toString(request.ViewsId)), 1)
	jsonData := new(bytes.Buffer)
	encodeErr := json.NewEncoder(jsonData).Encode(&request.Body)
	if encodeErr != nil {
		return nil, encodeErr
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, jsonData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentTypeApplicationJson}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateViewsSet200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusCreated {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateViewsSet201Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateViewsSet400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateViewsSet403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusMethodNotAllowed {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateViewsSet405Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) ShowVehicleInView(request *ShowVehicleInViewRequest) (ShowVehicleInViewResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/client/{clientId}/views/{viewsId}/{view}/{breakpoint}/{spec}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{clientId}", url.QueryEscape(toString(request.ClientId)), 1)
	endpoint = strings.Replace(endpoint, "{viewsId}", url.QueryEscape(toString(request.ViewsId)), 1)
	endpoint = strings.Replace(endpoint, "{view}", url.QueryEscape(toString(request.View)), 1)
	endpoint = strings.Replace(endpoint, "{breakpoint}", url.QueryEscape(toString(request.Breakpoint)), 1)
	endpoint = strings.Replace(endpoint, "{spec}", url.QueryEscape(toString(request.Spec)), 1)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ShowVehicleInView200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ShowVehicleInView403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ShowVehicleInView404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

/*
Get the list of permissions
a user can grant to other users.
*/
func (client *visAdminClient) GetPermissions(request *GetPermissionsRequest) (GetPermissionsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/permission"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetPermissions200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetPermissions200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetPermissions403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) DestroySession(request *DestroySessionRequest) (DestroySessionResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/session"
	method := "DELETE"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DestroySession200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DestroySession404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetUserInfo(request *GetUserInfoRequest) (GetUserInfoResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/session"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
	if request.SubID != nil {
		httpRequest.Header["subID"] = []string{toString(request.SubID)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetUserInfo200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetUserInfo200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetUserInfo400Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetUserInfo400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetUserInfo403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) CreateSession(request *CreateSessionRequest) (CreateSessionResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/session"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	jsonData := new(bytes.Buffer)
	encodeErr := json.NewEncoder(jsonData).Encode(&request.Body)
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateSession200Response)
			if err := fromString(httpResponse.Header.Get("X-Auth"), &response.XAuth); err != nil {
				return nil, err
			}
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(CreateSession400Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(CreateSession400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusUnauthorized {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateSession401Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetUsers(request *GetUsersRequest) (GetUsersResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/user"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetUsers200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetUsers200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetUsers403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) DeleteUser(request *DeleteUserRequest) (DeleteUserResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/user/{userId}"
	method := "DELETE"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{userId}", url.QueryEscape(toString(request.UserId)), 1)
	query := make(url.Values)
	if request.AllKeys != nil {
		query.Add("allKeys", toString(request.AllKeys))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteUser200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteUser403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteUser404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetUser(request *GetUserRequest) (GetUserResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/user/{userId}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{userId}", url.QueryEscape(toString(request.UserId)), 1)
	query := make(url.Values)
	if request.AllKeys != nil {
		query.Add("allKeys", toString(request.AllKeys))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetUser200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetUser200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetUser403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetUser404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) CreateOrUpdateUser(request *CreateOrUpdateUserRequest) (CreateOrUpdateUserResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/api/user/{userId}"
	method := "PUT"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{userId}", url.QueryEscape(toString(request.UserId)), 1)
	query := make(url.Values)
	if request.AllKeys != nil {
		query.Add("allKeys", toString(request.AllKeys))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	jsonData := new(bytes.Buffer)
	encodeErr := json.NewEncoder(jsonData).Encode(&request.Body)
	if encodeErr != nil {
		return nil, encodeErr
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, jsonData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentTypeApplicationJson}
	httpRequest.Header["X-Auth"] = []string{toString(request.XAuth)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateUser200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusCreated {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateUser201Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateUser400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateUser403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusMethodNotAllowed {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateOrUpdateUser405Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Get booking of session owner
func (client *visAdminClient) GetBooking(request *GetBookingRequest) (GetBookingResponse, error) {
	return nil, newNotSupportedContentType(415, "no supported content type")
}

// Get bookings of session owner
func (client *visAdminClient) GetBookings(request *GetBookingsRequest) (GetBookingsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/bookings"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	query := make(url.Values)
	if request.Ids != nil {
		query.Add("ids", toString(request.Ids))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Session-ID"] = []string{request.XSessionID}
	if request.Date != nil {
		httpRequest.Header["date"] = []string{toString(request.Date)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetBookings200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetBookings200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetBookings400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusUnauthorized {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetBookings401Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetBookings404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetBookings500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) ListModels(request *ListModelsRequest) (ListModelsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/brands/{brandId}/models"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{brandId}", url.QueryEscape(toString(request.BrandId)), 1)
	query := make(url.Values)
	if request.DriveConcept != nil {
		query.Add("driveConcept", toString(request.DriveConcept))
	}
	if request.LanguageId != nil {
		query.Add("languageId", toString(request.LanguageId))
	}
	if request.ClassId != nil {
		query.Add("classId", toString(request.ClassId))
	}
	if request.LineId != nil {
		query.Add("lineId", toString(request.LineId))
	}
	if request.Ids != nil {
		query.Add("ids", toString(request.Ids))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(ListModels200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(ListModels200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetClasses(request *GetClassesRequest) (GetClassesResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/classes/{productGroup}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{productGroup}", url.QueryEscape(toString(request.ProductGroup)), 1)
	query := make(url.Values)
	if request.ComponentTypes != nil {
		query.Add("componentTypes", toString(request.ComponentTypes))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetClasses200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetClasses200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetClasses400Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetClasses400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) Code(request *CodeRequest) (CodeResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/code"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	query := make(url.Values)
	query.Add("session", toString(request.Session))
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
	queryInBody := make(url.Values)
	if request.State != nil {
		queryInBody.Add("state", toString(request.State))
	}
	if request.ResponseMode != nil {
		queryInBody.Add("response_mode", toString(request.ResponseMode))
	}
	queryInBody.Add("code", toString(request.Code))
	encodedQueryInBody := queryInBody.Encode()
	formData := bytes.NewBufferString(encodedQueryInBody)
	httpRequest, reqErr := http.NewRequest(method, endpoint, formData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentTypeApplicationFormUrlencoded}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(Code200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(Code200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(Code400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusUnauthorized {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(Code401Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusNotFound {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(Code404Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(Code500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

/*
Deletes the user session matching the *X-Auth* header.
*/
func (client *visAdminClient) DeleteCustomerSession(request *DeleteCustomerSessionRequest) (DeleteCustomerSessionResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/customer/session"
	method := "DELETE"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	httpRequest, reqErr := http.NewRequest(method, endpoint, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header["X-Session-ID"] = []string{request.XSessionID}
	if request.XRequestID != nil {
		httpRequest.Header["X-Request-ID"] = []string{toString(request.XRequestID)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusNoContent {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteCustomerSession204Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusUnauthorized {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteCustomerSession401Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DeleteCustomerSession500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

/*
Creates a customer session for a given OpenID authentication token.
*/
func (client *visAdminClient) CreateCustomerSession(request *CreateCustomerSessionRequest) (CreateCustomerSessionResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/customer/session"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	queryInBody := make(url.Values)
	queryInBody.Add("code", toString(request.Code))
	if request.Locale != nil {
		queryInBody.Add("locale", toString(request.Locale))
	}
	encodedQueryInBody := queryInBody.Encode()
	formData := bytes.NewBufferString(encodedQueryInBody)
	httpRequest, reqErr := http.NewRequest(method, endpoint, formData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentTypeApplicationFormUrlencoded}
	if request.XRequestID != nil {
		httpRequest.Header["X-Request-ID"] = []string{toString(request.XRequestID)}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusCreated {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(CreateCustomerSession201Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(CreateCustomerSession201Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusUnauthorized {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateCustomerSession401Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusForbidden {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateCustomerSession403Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusUnprocessableEntity {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(CreateCustomerSession422Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(CreateCustomerSession422Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(CreateCustomerSession500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

/*
Downloads a file that is a property within a nested structure in the response body
*/
func (client *visAdminClient) DownloadNestedFile(request *DownloadNestedFileRequest) (DownloadNestedFileResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/download/nested/file"
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(DownloadNestedFile200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(DownloadNestedFile200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Retrieve a image
func (client *visAdminClient) DownloadImage(request *DownloadImageRequest) (DownloadImageResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/download/{image}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{image}", url.QueryEscape(toString(request.Image)), 1)
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
		httpRequest.Header["Accept"] = []string{"image/png"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeInList(contentTypesForFiles, contentTypeOfResponse) {
			response := new(DownloadImage200Response)
			if err := fromString(httpResponse.Header.Get("Content-Type"), &response.ContentType); err != nil {
				return nil, err
			}
			response.Body = httpResponse.Body
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(DownloadImage500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) ListElements(request *ListElementsRequest) (ListElementsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/elements"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	query := make(url.Values)
	if request.Page != nil {
		query.Add("_page", toString(request.Page))
	}
	if request.PerPage != nil {
		query.Add("_perPage", toString(request.PerPage))
	}
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(ListElements200Response)
			if err := fromString(httpResponse.Header.Get("X-Total-Count"), &response.XTotalCount); err != nil {
				return nil, err
			}
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(ListElements200Response)
			if err := fromString(httpResponse.Header.Get("X-Total-Count"), &response.XTotalCount); err != nil {
				return nil, err
			}
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(ListElements500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) FileUpload(request *FileUploadRequest) (FileUploadResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/file-upload"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	formData := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(formData)
	if request.FormData.File != nil {
		fileWriter0, writerErr0 := bodyWriter.CreateFormFile("file", "file")
		if writerErr0 != nil {
			bodyWriter.Close()
			return nil, writerErr0
		}
		_, copyFileErr0 := io.Copy(fileWriter0, request.FormData.File.Content)
		if copyFileErr0 != nil {
			bodyWriter.Close()
			return nil, copyFileErr0
		}
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	httpRequest, reqErr := http.NewRequest(method, endpoint, formData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentType}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusNoContent {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(FileUpload204Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(FileUpload500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Retrieve a file
func (client *visAdminClient) DownloadFile(request *DownloadFileRequest) (DownloadFileResponse, error) {
	return nil, newNotSupportedContentType(415, "no supported content type")
}

// Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.
func (client *visAdminClient) FindByTags(request *FindByTagsRequest) (FindByTagsResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/findByTags"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	query := make(url.Values)
	query.Add("tags", toString(request.Tags))
	encodedQuery := query.Encode()
	if encodedQuery != "" {
		endpoint += "?" + encodedQuery
	}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(FindByTags200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(FindByTags200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(FindByTags400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// Retrieve a file
func (client *visAdminClient) GenericFileDownload(request *GenericFileDownloadRequest) (GenericFileDownloadResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/generic/download/{ext}"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	endpoint = strings.Replace(endpoint, "{ext}", url.QueryEscape(toString(request.Ext)), 1)
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
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeInList(contentTypesForFiles, contentTypeOfResponse) {
			response := new(GenericFileDownload200Response)
			if err := fromString(httpResponse.Header.Get("Content-Type"), &response.ContentType); err != nil {
				return nil, err
			}
			if err := fromString(httpResponse.Header.Get("Pragma"), &response.Pragma); err != nil {
				return nil, err
			}
			response.Body = httpResponse.Body
			return response, nil
		}
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GenericFileDownload200Response)
			if err := fromString(httpResponse.Header.Get("Content-Type"), &response.ContentType); err != nil {
				return nil, err
			}
			if err := fromString(httpResponse.Header.Get("Pragma"), &response.Pragma); err != nil {
				return nil, err
			}
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GenericFileDownload200Response)
			if err := fromString(httpResponse.Header.Get("Content-Type"), &response.ContentType); err != nil {
				return nil, err
			}
			if err := fromString(httpResponse.Header.Get("Pragma"), &response.Pragma); err != nil {
				return nil, err
			}
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GenericFileDownload500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

// get rental
func (client *visAdminClient) GetRental(request *GetRentalRequest) (GetRentalResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/rental"
	method := "GET"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	jsonData := new(bytes.Buffer)
	encodeErr := json.NewEncoder(jsonData).Encode(&request.Body)
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(GetRental200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetRental400Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetRental400Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) GetShoes(request *GetShoesRequest) (GetShoesResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/shop/shoes"
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
		httpRequest.Header["Accept"] = []string{"application/hal+json"}
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson || contentTypeOfResponse == contentTypeApplicationHalJson {
			response := new(GetShoes200Response)
			decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response.Body)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return response, nil
		} else if contentTypeOfResponse == "" {
			response := new(GetShoes200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}

func (client *visAdminClient) PostUpload(request *PostUploadRequest) (PostUploadResponse, error) {
	if request == nil {
		return nil, newRequestObjectIsNilError
	}
	path := "/upload"
	method := "POST"
	endpoint := client.baseURL + path
	httpContext := newHttpContextWrapper(client.ctx)
	formData := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(formData)
	if request.FormData.Upfile != nil {
		fileWriter0, writerErr0 := bodyWriter.CreateFormFile("upfile", "upfile")
		if writerErr0 != nil {
			bodyWriter.Close()
			return nil, writerErr0
		}
		_, copyFileErr0 := io.Copy(fileWriter0, request.FormData.Upfile.Content)
		if copyFileErr0 != nil {
			bodyWriter.Close()
			return nil, copyFileErr0
		}
	}
	if request.FormData.Note != nil {
		fieldData0 := toString(request.FormData.Note)
		fieldWriter0, fieldErr0 := bodyWriter.CreateFormField("note")
		if fieldErr0 != nil {
			bodyWriter.Close()
			return nil, fieldErr0
		}
		_, writeFieldErr0 := fieldWriter0.Write([]byte(fieldData0))
		if writeFieldErr0 != nil {
			bodyWriter.Close()
			return nil, writeFieldErr0
		}
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	httpRequest, reqErr := http.NewRequest(method, endpoint, formData)
	if reqErr != nil {
		return nil, reqErr
	}
	httpRequest.Header[contentTypeHeader] = []string{contentType}
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
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == http.StatusOK {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(PostUpload200Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if httpResponse.StatusCode == http.StatusInternalServerError {
		contentTypeOfResponse := extractContentType(httpResponse.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == "" {
			response := new(PostUpload500Response)
			return response, nil
		}
		return nil, newNotSupportedContentType(415, contentTypeOfResponse)
	}

	if client.hooks.OnUnknownResponseCode != nil {
		message := client.hooks.OnUnknownResponseCode(httpResponse, httpRequest)
		return nil, newErrOnUnknownResponseCode(message)
	}
	return nil, newErrUnknownResponse(httpResponse.StatusCode)
}
