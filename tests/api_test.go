package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ExperienceOne/apikit/generator/stringutil"
	"github.com/ExperienceOne/apikit/tests/api"

	routing "github.com/go-ozzo/ozzo-routing"
)

var (
	id       = "id"
	password = "password"
	auth     = "01234567890"
)

var VisAdminClient api.VisAdminClient
var fakeClient api.VisAdminClient
var BadRequestVisAdminClient api.VisAdminClient
var testServerWrapper *api.VisAdminServer

func TestMain(m *testing.M) {

	middlewares := []api.Middleware{
		{Handler: api.RouterPanicMiddleware()},
		{Handler: api.RouterPopulateContextMiddleware()},
	}

	testServerWrapper = api.NewVisAdminServer(&api.ServerOpts{
		Middleware:   middlewares,
		ErrorHandler: log.Println,
		OnStart: func(router *routing.Router) {
			router.Get("/no-spec-route", func(context *routing.Context) error {
				context.Write(1)
				return nil
			})
		},
	})

	testServerWrapper.SetGetUsersHandler(api.GetUsers)
	testServerWrapper.SetGetClientHandler(api.GetClient, api.Middleware{Handler: api.HandlerMiddleware()})
	testServerWrapper.SetCreateSessionHandler(api.CreateSession)
	testServerWrapper.SetCreateOrUpdateUserHandler(api.CreateOrUpdateUser)
	testServerWrapper.SetGetUserInfoHandler(api.GetUserInfo)
	testServerWrapper.SetDownloadImageHandler(api.DownloadImage)
	testServerWrapper.SetGetPermissionsHandler(api.GetPermissions)
	testServerWrapper.SetPostUploadHandler(api.PostUpload)
	testServerWrapper.SetDownloadFileHandler(api.DownloadFile)
	testServerWrapper.SetGetBookingsHandler(api.GetBookings)
	testServerWrapper.SetGetBookingHandler(api.GetBooking)
	testServerWrapper.SetGetClassesHandler(api.GetClasses)
	testServerWrapper.SetListModelsHandler(api.GetBrands)
	testServerWrapper.SetGenericFileDownloadHandler(api.GenericFileDownload)
	testServerWrapper.SetGetRentalHandler(api.GetRental)
	testServerWrapper.SetListElementsHandler(api.ListElement)
	testServerWrapper.SetCodeHandler(api.Code)
	testServerWrapper.SetCreateCustomerSessionHandler(api.CreateCustomerSession)
	testServerWrapper.SetDownloadNestedFileHandler(api.NestedFileDownload)
	testServerWrapper.SetGetShoesHandler(api.GetShoes)
	testServerWrapper.SetFileUploadHandler(api.FileUpload)
	testServerWrapper.SetFindByTagsHandler(api.FindByTags)

	go testServerWrapper.Start(4567)
	time.Sleep(1 * time.Second)

	opts := api.Opts{
		Hooks: api.DevHook(),
		Ctx:   nil,
	}

	VisAdminClient = api.NewVisAdminClient(new(http.Client), "http://localhost:4567", opts)

	// use round tripper to figure out if bad request code is set
	lt := new(api.BadRequestHTTPTransport)
	c := &http.Client{
		Transport: lt,
		Timeout:   time.Second * 3,
	}
	BadRequestVisAdminClient = api.NewVisAdminClient(c, "http://localhost:4567", opts)

	// setup a fake client and server to check the raw request and response object
	contextApiKey := "x-context-api-key"
	contextApiKeyValue := "test"
	testHandler := func(resp http.ResponseWriter, req *http.Request) {
		apiKey := req.Header.Get(contextApiKey)
		if apiKey == contextApiKeyValue {
			resp.WriteHeader(http.StatusNoContent)
			return
		}
		resp.WriteHeader(http.StatusForbidden)
	}

	fakeMux := http.NewServeMux()
	fakeMux.HandleFunc("/api/client", testHandler)
	fakeServer := httptest.NewServer(fakeMux)
	headers := make(http.Header)
	headers.Set(contextApiKey, contextApiKeyValue)

	httpContext := api.CreateHttpContext(headers)
	opts.Ctx = httpContext

	fakeClient = api.NewVisAdminClient(new(http.Client), fakeServer.URL, opts)

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestNoSpec(t *testing.T) {
	resp, err := http.Get("http://localhost:4567/no-spec-route")
	if err != nil {
		t.Fatalf("http get has failed (%v)", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("response is bad (http status: %v)", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read body of response")
	}

	if string(respBody) != "1" {
		t.Fatalf("anwers is not 1, got=%s", string(respBody))
	}
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	request := &api.GetUsersRequest{
		XAuth: "sessionID",
	}

	response, err := VisAdminClient.GetUsers(request)
	if err != nil {
		t.Fatalf("error sending GetUsers GET request: %v", err)
	}

	if _, ok := response.(*api.GetUsers200Response); !ok {
		t.Fatalf("error GetUsers response is bad: %#v", err)
	}
}

func TestGetClient(t *testing.T) {
	t.Parallel()

	request := &api.GetClientRequest{
		XAuth:    "sessionID",
		ClientId: "test",
	}

	response, err := VisAdminClient.GetClient(request)
	if err != nil {
		t.Fatalf("error sending GetClient GET request: %v", err)
	}

	if _, ok := response.(*api.GetClient200Response); !ok {
		t.Fatalf("error GetClient response is bad: %#v", err)
	}
}

func TestGetBookings200(t *testing.T) {
	t.Parallel()

	response, err := VisAdminClient.GetBookings(&api.GetBookingsRequest{XSessionID: "xsession"})
	if err != nil {
		t.Fatalf("error sending GetBookings200 GET request: %v", err)
		return
	}

	getBookings200Response, ok := response.(*api.GetBookings200Response)
	if !ok {
		t.Fatalf("error GetBookings200 response is bad: %#v", response)
		return
	}

	// We had an deps issue in the past that didn't handle a property with suffix ID properly
	// For instance lib converts modelID to modelId, that isn't ok
	// check if BookingID doesn't trigger a panic
	if len(getBookings200Response.Body) == 0 || *getBookings200Response.Body[0].BookingID != "3600278183" {
		t.Fatalf("error GetBookings200 response is bad: %#v", response)
	}
}

func TestGetBookings401(t *testing.T) {
	t.Parallel()

	response, err := VisAdminClient.GetBookings(&api.GetBookingsRequest{XSessionID: "xsession1"})
	if err != nil {
		t.Fatalf("error sending GetBookings401 GET request: %v", err)
		return
	}

	if _, ok := response.(*api.GetBookings401Response); !ok {
		t.Fatalf("error GetBookings401 GET reponse is bad: %#v", response)
	}
}

func TestGetBookingErrNotSupportedContentType(t *testing.T) {
	t.Parallel()

	response, err := VisAdminClient.GetBooking(&api.GetBookingRequest{})
	if err == nil {
		t.Fatalf("error sending GetBooking GET request: %#v", response)
		return
	}

	if !strings.Contains(err.Error(), "error unsupported media type") {
		t.Fatalf("error sending GetBooking GET request: %v", err)
	}
}

func TestGetUserInfo400(t *testing.T) {
	t.Parallel()

	response, err := BadRequestVisAdminClient.GetUserInfo(&api.GetUserInfoRequest{XAuth: "trigger400"})
	if err != nil {
		t.Fatalf("error sending GetUserInfo GET request: %#v", err)
		return
	}

	if _, ok := response.(*api.GetUserInfo400Response); !ok {
		t.Fatalf("response to get user info request without auth header is not 400: %#v", response)
	}
}

// Test case for post upload
func TestPostUpload(t *testing.T) {
	t.Parallel()

	note := "Hello world"
	buffer := bytes.NewBufferString("Hello world, file")
	request := &api.PostUploadRequest{
		FormData: api.PostUploadRequestFormData{
			Upfile: &api.MimeFile{
				Content: ioutil.NopCloser(buffer),
			},
			Note: &note,
		},
	}

	resp, err := VisAdminClient.PostUpload(request)
	if err != nil {
		t.Fatalf("error sending upload error: %v", err)
	}

	_, ok := resp.(*api.PostUpload200Response)
	if !ok {
		t.Fatalf("error PostUploadResponse is bad: %#v", resp)
	}
}

// Test case for download a image
func TestDownloadImage(t *testing.T) {
	t.Parallel()

	request := &api.DownloadImageRequest{
		Image: "dummy.png",
	}

	resp, err := VisAdminClient.DownloadImage(request)
	if err != nil {
		t.Fatalf("image download error: %v", err)
	}

	resp200, ok := resp.(*api.DownloadImage200Response)
	if !ok {
		t.Fatalf("image download faild: %#v", resp)
	}

	if resp200.Body == nil {
		t.Fatalf("image is nil")
	}

	data, err := ioutil.ReadAll(resp200.Body)
	if err != nil {
		t.Fatalf("image download error: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("got empty image")
	}

	resp200.Body.Close()
}

func TestGetShoes(t *testing.T) {
	t.Parallel()

	request := &api.GetShoesRequest{}

	resp, err := VisAdminClient.GetShoes(request)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := resp.(*api.GetShoes200Response)
	if !ok {
		t.Fatalf("get hal+json faild: %#v", resp)
	}
}

// Test case for not supported content type
func TestDownloadFileNotSupportedContentType(t *testing.T) {
	t.Parallel()

	request := &api.DownloadFileRequest{
		File: "dummy.png",
	}

	_, err := VisAdminClient.DownloadFile(request)
	if err == nil {
		t.Fatal("error is nil")
	}

	if !strings.Contains(err.Error(), "media type") {
		t.Fatalf("error message is bad (%v)", err)
	}
}

// Test case for missing session data
func TestCreateSessionFail(t *testing.T) {
	t.Parallel()

	_, err := VisAdminClient.CreateSession(nil)
	if err == nil {
		t.Fatal("sending nil request did not cause error")
	}
}

// test case to trigger bad request
func TestCreateSessionBadRequest(t *testing.T) {
	t.Parallel()

	response, err := VisAdminClient.CreateSession(&api.CreateSessionRequest{})
	if err != nil {
		t.Fatalf("error sending CreateSession empty POST body caused error: %v", err)
		return
	}

	if _, ok := response.(*api.CreateSession400Response); !ok {
		t.Fatalf("error CreateSession response is bad: %#v", response)
	}
}

// test case to simulate a user workflow
func TestUserWorkflow(t *testing.T) {
	t.Parallel()

	response, err := VisAdminClient.CreateSession(&api.CreateSessionRequest{Body: api.Object2{Id: id, Password: password}})
	if err != nil {
		t.Fatalf("error sending CreateSession POST request: %v", err)
		return
	}

	if response200, ok := response.(*api.CreateSession200Response); !ok {
		t.Fatalf("response to valid create session request is not 200: %#v", response)
	} else {
		if response200.XAuth != auth {
			t.Fatal("create session response is missing X-Auth header")
		} else {
			auth := response200.XAuth

			response, err := VisAdminClient.GetUserInfo(&api.GetUserInfoRequest{XAuth: "skfsdfj"})
			if err != nil {
				t.Fatalf("error sending GetUserInfo GET request: %v", err)
				return
			}

			if _, ok := response.(*api.GetUserInfo403Response); !ok {
				t.Fatalf("response to get user info request without auth header is not 403: %#v", response)
			}

			response, err = VisAdminClient.GetUserInfo(&api.GetUserInfoRequest{XAuth: auth})
			if err != nil {
				t.Fatalf("error sending GetUserInfo GET request: %v", err)
				return
			}

			if response, ok := response.(*api.GetUserInfo200Response); !ok {
				t.Fatalf("response to valid get user info request is not 200: %#v", response)
			} else {
				if response.Body.Id != id || response.Body.Password != password || len(response.Body.Permissions) != 3 {
					t.Fatalf("unexpected get user info response body: %v", response.Body)
				}

				mapValue, ok := response.Body.GrantedProtocolMappers["key"]
				if !ok {
					t.Fatalf("unexpected get user info response body: %v", response.Body)
				}
				if mapValue != "value" {
					t.Fatalf("unexpected get user info response body: %v", response.Body)
				}
			}

		}
	}
}

func TestGetClasses(t *testing.T) {
	t.Parallel()

	request := new(api.GetClassesRequest)
	request.ComponentTypes = append(request.ComponentTypes, api.ComponentTypesWHEELS)
	request.ProductGroup = api.ProductGroupPKW

	response, err := VisAdminClient.GetClasses(request)
	if err != nil {
		t.Fatalf("error sending GetClasses GET request: %v", err)
		return
	}

	if _, ok := response.(*api.GetClasses200Response); !ok {
		t.Fatalf("error GetClasses response is bad: %#v", response)
		return
	}
}

// test case to check if auto generated validation is working
func TestCheckEmail(t *testing.T) {
	t.Parallel()

	// check optional validation
	invalidEmail := "sdifisifjs"
	invalidEmailResponse, _ := BadRequestVisAdminClient.CreateOrUpdateUser(&api.CreateOrUpdateUserRequest{
		XAuth:  "999",
		UserId: "9",
		Body: api.User{
			Id:       "9",
			Password: "aosjdfosjdf",
			Email:    &invalidEmail,
		},
	})
	if _, ok := invalidEmailResponse.(*api.CreateOrUpdateUser400Response); !ok {
		t.Fatalf("invalid email: unexpected CreateOrUpdateUserResponse : %#v", invalidEmailResponse)
	}
}

func TestCheckEmailOmitIfEmpty(t *testing.T) {
	t.Parallel()

	emailResponse, _ := VisAdminClient.CreateOrUpdateUser(&api.CreateOrUpdateUserRequest{
		XAuth:  "999",
		UserId: "9",
		Body: api.User{
			Id:       "9",
			Password: "aosjdfosjdf",
		},
	})

	if _, ok := emailResponse.(*api.CreateOrUpdateUser201Response); !ok {
		t.Fatalf("without email: unexpected CreateOrUpdateUserResponse : %#v", emailResponse)
	}
}

// test case to if an empty response is handled safe
func TestHandleSafePanic(t *testing.T) {
	t.Parallel()

	_, err := VisAdminClient.GetPermissions(&api.GetPermissionsRequest{XAuth: "99999999"})
	if err == nil || strings.Contains(err.Error(), "unknown response status code '500'") {
		t.Fatalf("unexpected error: %#v", err.Error())
	}

	t.Log(err.Error())
}

// test case to if an empty response is handled safe
func TestSetHeaderFromContex(t *testing.T) {
	t.Parallel()

	req := &api.GetClientsRequest{
		XAuth: "test",
	}

	resp, err := fakeClient.GetClients(req)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := resp.(*api.GetClients204Response); !ok {
		t.Fatalf("unexpected error: %#v", resp)
	}
}

func TestGetListModelsQueryChaining(t *testing.T) {
	t.Parallel()

	testString := "test"
	classId := ""
	req := &api.ListModelsRequest{
		BrandId:    "without_first_query_parameter",
		LanguageId: nil,
		ClassId:    &testString,
		LineId:     &classId,
		Ids:        []int64{1, 2, 3},
	}

	resp, err := VisAdminClient.ListModels(req)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := resp.(*api.ListModels200Response); !ok {
		t.Fatalf("unexpected error: %#v", resp)
	}

	req = &api.ListModelsRequest{
		BrandId:    "with_first_query_parameter",
		LanguageId: &testString,
	}

	resp, err = VisAdminClient.ListModels(req)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := resp.(*api.ListModels200Response); !ok {
		t.Fatalf("unexpected error: %#v", resp)
	}
}

func TestGetJSONFile(t *testing.T) {
	t.Parallel()

	request := &api.GenericFileDownloadRequest{
		Ext: ".json",
	}

	resp, err := VisAdminClient.GenericFileDownload(request)
	if err != nil {
		t.Fatalf("json download error: %v", err)
	}

	resp200, ok := resp.(*api.GenericFileDownload200Response)
	if !ok {
		t.Fatalf("json download faild: %#v", resp)
	}

	if resp200.Body == nil {
		t.Fatalf("json is empty")
	}
}

func TestGetRental(t *testing.T) {
	t.Parallel()

	fields := []string{"Class", "LockStatus", "MaxDoors", "MinDoors", "Status", "Color", "Valid", "StationID", "HomeID", "Id", "Website"}

	resp, err := VisAdminClient.GetRental(&api.GetRentalRequest{
		Body: api.Rental{},
	})
	if err != nil {
		t.Fatal(err)
	}

	if resp400, ok := resp.(*api.GetRental400Response); ok {
		if len(resp400.Body.Errors) == 0 {
			t.Fatalf("errors is empty")
		}
		for _, err := range resp400.Body.Errors {
			if !stringutil.InStringSlice(fields, *err.Field) {
				t.Error(*err.Field)
			}
		}
	} else {
		t.Fatal(fmt.Sprintf("resp is bad (%v)", resp))
	}

	color := "re"
	uuid := "sfdsf"
	website := "http:example.de"
	website2 := ""

	fields = append(fields, "WebsiteOptional")

	resp, err = VisAdminClient.GetRental(&api.GetRentalRequest{Body: api.Rental{
		Color:           &color,
		Id:              uuid,
		Website:         website,
		WebsiteOptional: &website2,
	}})
	if err != nil {
		t.Fatal(err)
	}

	if resp400, ok := resp.(*api.GetRental400Response); ok {
		if len(resp400.Body.Errors) == 0 {
			t.Fatalf("errors is empty")
		}
		for _, err := range resp400.Body.Errors {
			if !stringutil.InStringSlice(fields, *err.Field) {
				t.Error(*err.Field)
			}
		}
	} else {
		t.Fatal(fmt.Sprintf("resp is bad (%v)", resp))
	}

	homeID := "r088888"

	resp, err = VisAdminClient.GetRental(&api.GetRentalRequest{Body: api.Rental{HomeID: &homeID}})
	if err != nil {
		t.Fatal(err)
	}

	if resp400, ok := resp.(*api.GetRental400Response); ok {
		if len(resp400.Body.Errors) == 0 {
			t.Fatalf("errors is empty")
		}
		for _, err := range resp400.Body.Errors {
			if !stringutil.InStringSlice(fields, *err.Field) {
				t.Error(err.Field)
			}
		}
	} else {
		t.Fatal(fmt.Sprintf("resp is bad (%v)", resp))
	}
}

func TestListElements(t *testing.T) {
	t.Parallel()

	resp, err := VisAdminClient.ListElements(&api.ListElementsRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if resp200, ok := resp.(*api.ListElements200Response); ok {
		if resp200.XTotalCount != 5 {
			t.Error(fmt.Sprintf("total count is bad (%d)", resp200.XTotalCount))
		}
	} else {
		t.Fatal(fmt.Sprintf("resp is bad (%#v)", resp))
	}
}

func TestCode(t *testing.T) {
	t.Parallel()

	resp, err := VisAdminClient.Code(&api.CodeRequest{
		State:   []int64{1, 2, 3},
		Session: "test",
		Code:    "rueeerjfid-ifjsfjfs-osdjfj",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, ok := resp.(*api.Code200Response)
	if !ok {
		t.Error("response is bad")
	}
}

func TestCreateCustomerSession_InvalidToken(t *testing.T) {

	request := &api.CreateCustomerSessionRequest{Code: "invalid_token"}
	resp, err := VisAdminClient.CreateCustomerSession(request)
	if err != nil {
		t.Fatalf("unexpected error = %+v", err)
	}

	got := resp.StatusCode()
	if got != http.StatusUnauthorized {
		t.Fatalf("CreateCustomerSession() got = %+v, want = %+v", got, http.StatusCreated)
	}
}

func TestCreateCustomerSession_Success(t *testing.T) {

	request := &api.CreateCustomerSessionRequest{Code: "abc"}
	resp, err := VisAdminClient.CreateCustomerSession(request)
	if err != nil {
		t.Fatalf("unexpected error = %+v", err)
	}

	got := resp.StatusCode()
	if got != http.StatusCreated {
		t.Fatalf("CreateCustomerSession() got = %+v, want = %+v", got, http.StatusCreated)
	}
}

func TestGetClientVersion(t *testing.T) {
	t.Parallel()
	version := api.ApikitVersion()
	isValidVersion(t, version)
}
func TestGetServerVersion(t *testing.T) {
	version := api.ApikitVersion()
	isValidVersion(t, version)
}

func isValidVersion(t *testing.T, version *api.VersionInfo) {
	if version == nil {
		t.Errorf("version is nil")
	}
	if version.GitCommit == "" {
		t.Error("git commit is empty")
	}
	if version.GoVersion == "" {
		t.Error("go version is empty")
	}
	if version.GitBranch == "" {
		t.Error("git branch is empty")
	}
	if version.GitTag == "" {
		t.Error("go tag is empty")
	}
	if version.BuildTime == "" {
		t.Error("build time is empty")
	}
}

func TestDownloadNestedFile(t *testing.T) {
	t.Parallel()

	resp, err := VisAdminClient.DownloadNestedFile(&api.DownloadNestedFileRequest{})
	if err != nil {
		t.Fatalf("nested file download failed: %v", err)
	}

	resp200, ok := resp.(*api.DownloadNestedFile200Response)
	if !ok {
		t.Fatalf("nested file download failed with status code: %#v", resp)
	}

	if resp200.Body.Data == nil {
		t.Fatal("file is nil")
	}
}

func TestFileUpload(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBufferString("Hello world, file")
	req := api.FileUploadRequest{
		FormData: api.FileUploadRequestFormData{
			File: &api.MimeFile{
				Content: ioutil.NopCloser(buffer),
			},
		},
	}

	resp, err := VisAdminClient.FileUpload(&req)
	if err != nil {
		t.Fatalf("error during file upload: %v", err)
	}

	_, ok := resp.(*api.FileUpload204Response)
	if !ok {
		t.Fatalf("unexpected FileUploadResponse: %#v", resp)
	}
}

func TestMinMaxItemsValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tags       []string
		statusCode int
	}{
		{
			name:       "valid",
			tags:       []string{"tag1", "tag2", "tag3"},
			statusCode: http.StatusOK,
		},
		{
			name:       "valid",
			tags:       []string{"tag1", "tag2"},
			statusCode: http.StatusOK,
		},
		{
			name:       "valid",
			tags:       []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
			statusCode: http.StatusOK,
		},
		{
			name:       "min_invalid",
			tags:       []string{"tag1"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "max_invalid",
			tags:       []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req := api.FindByTagsRequest{
				Tags: test.tags,
			}

			resp, err := VisAdminClient.FindByTags(&req)
			if err != nil {
				t.Fatalf("error find by tags failed: %v", err)
			}

			if test.statusCode != resp.StatusCode() {
				t.Fatalf("response status is bad, want:'%d', got:'%d'", test.statusCode, resp.StatusCode())
			}
		})
	}
}
