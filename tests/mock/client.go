package api

import (
	"context"
	"encoding/json"
	"net/http"
)

func NewVisAdminClientMock(httpClient *http.Client, baseUrl string, ctx ...context.Context) *VisAdminClientMock {
	return &VisAdminClientMock{}
}

func (client *VisAdminClientMock) GetClients(request *GetClientsRequest) (GetClientsResponse, error) {
	if client.GetClientsStatusCode == 200 {
		response := new(GetClients200Response)
		return response, nil
	}
	if client.GetClientsStatusCode == 204 {
		response := new(GetClients204Response)
		return response, nil
	}
	if client.GetClientsStatusCode == 403 {
		response := new(GetClients403Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) DeleteClient(request *DeleteClientRequest) (DeleteClientResponse, error) {
	if client.DeleteClientStatusCode == 200 {
		response := new(DeleteClient200Response)
		return response, nil
	}
	if client.DeleteClientStatusCode == 403 {
		response := new(DeleteClient403Response)
		return response, nil
	}
	if client.DeleteClientStatusCode == 404 {
		response := new(DeleteClient404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetClient(request *GetClientRequest) (GetClientResponse, error) {
	if client.GetClientStatusCode == 200 {
		response := new(GetClient200Response)
		return response, nil
	}
	if client.GetClientStatusCode == 403 {
		response := new(GetClient403Response)
		return response, nil
	}
	if client.GetClientStatusCode == 404 {
		response := new(GetClient404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) CreateOrUpdateClient(request *CreateOrUpdateClientRequest) (CreateOrUpdateClientResponse, error) {
	if client.CreateOrUpdateClientStatusCode == 200 {
		response := new(CreateOrUpdateClient200Response)
		return response, nil
	}
	if client.CreateOrUpdateClientStatusCode == 201 {
		response := new(CreateOrUpdateClient201Response)
		return response, nil
	}
	if client.CreateOrUpdateClientStatusCode == 400 {
		response := new(CreateOrUpdateClient400Response)
		return response, nil
	}
	if client.CreateOrUpdateClientStatusCode == 403 {
		response := new(CreateOrUpdateClient403Response)
		return response, nil
	}
	if client.CreateOrUpdateClientStatusCode == 405 {
		response := new(CreateOrUpdateClient405Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetViewsSets(request *GetViewsSetsRequest) (GetViewsSetsResponse, error) {
	if client.GetViewsSetsStatusCode == 200 {
		response := new(GetViewsSets200Response)
		return response, nil
	}
	if client.GetViewsSetsStatusCode == 403 {
		response := new(GetViewsSets403Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) DeleteViewsSet(request *DeleteViewsSetRequest) (DeleteViewsSetResponse, error) {
	if client.DeleteViewsSetStatusCode == 200 {
		response := new(DeleteViewsSet200Response)
		return response, nil
	}
	if client.DeleteViewsSetStatusCode == 403 {
		response := new(DeleteViewsSet403Response)
		return response, nil
	}
	if client.DeleteViewsSetStatusCode == 404 {
		response := new(DeleteViewsSet404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetViewsSet(request *GetViewsSetRequest) (GetViewsSetResponse, error) {
	if client.GetViewsSetStatusCode == 200 {
		response := new(GetViewsSet200Response)
		return response, nil
	}
	if client.GetViewsSetStatusCode == 403 {
		response := new(GetViewsSet403Response)
		return response, nil
	}
	if client.GetViewsSetStatusCode == 404 {
		response := new(GetViewsSet404Response)
		return response, nil
	}
	return nil, nil
}

// Make this viewset the active one for the client.
func (client *VisAdminClientMock) ActivateViewsSet(request *ActivateViewsSetRequest) (ActivateViewsSetResponse, error) {
	if client.ActivateViewsSetStatusCode == 200 {
		response := new(ActivateViewsSet200Response)
		return response, nil
	}
	if client.ActivateViewsSetStatusCode == 403 {
		response := new(ActivateViewsSet403Response)
		return response, nil
	}
	if client.ActivateViewsSetStatusCode == 404 {
		response := new(ActivateViewsSet404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) CreateOrUpdateViewsSet(request *CreateOrUpdateViewsSetRequest) (CreateOrUpdateViewsSetResponse, error) {
	if client.CreateOrUpdateViewsSetStatusCode == 200 {
		response := new(CreateOrUpdateViewsSet200Response)
		return response, nil
	}
	if client.CreateOrUpdateViewsSetStatusCode == 201 {
		response := new(CreateOrUpdateViewsSet201Response)
		return response, nil
	}
	if client.CreateOrUpdateViewsSetStatusCode == 400 {
		response := new(CreateOrUpdateViewsSet400Response)
		return response, nil
	}
	if client.CreateOrUpdateViewsSetStatusCode == 403 {
		response := new(CreateOrUpdateViewsSet403Response)
		return response, nil
	}
	if client.CreateOrUpdateViewsSetStatusCode == 405 {
		response := new(CreateOrUpdateViewsSet405Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) ShowVehicleInView(request *ShowVehicleInViewRequest) (ShowVehicleInViewResponse, error) {
	if client.ShowVehicleInViewStatusCode == 200 {
		response := new(ShowVehicleInView200Response)
		return response, nil
	}
	if client.ShowVehicleInViewStatusCode == 403 {
		response := new(ShowVehicleInView403Response)
		return response, nil
	}
	if client.ShowVehicleInViewStatusCode == 404 {
		response := new(ShowVehicleInView404Response)
		return response, nil
	}
	return nil, nil
}

/*
Get the list of permissions
a user can grant to other users.
*/
func (client *VisAdminClientMock) GetPermissions(request *GetPermissionsRequest) (GetPermissionsResponse, error) {
	if client.GetPermissionsStatusCode == 200 {
		response := new(GetPermissions200Response)
		return response, nil
	}
	if client.GetPermissionsStatusCode == 403 {
		response := new(GetPermissions403Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) DestroySession(request *DestroySessionRequest) (DestroySessionResponse, error) {
	if client.DestroySessionStatusCode == 200 {
		response := new(DestroySession200Response)
		return response, nil
	}
	if client.DestroySessionStatusCode == 404 {
		response := new(DestroySession404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetUserInfo(request *GetUserInfoRequest) (GetUserInfoResponse, error) {
	if client.GetUserInfoStatusCode == 200 {
		response := new(GetUserInfo200Response)
		return response, nil
	}
	if client.GetUserInfoStatusCode == 400 {
		response := new(GetUserInfo400Response)
		return response, nil
	}
	if client.GetUserInfoStatusCode == 403 {
		response := new(GetUserInfo403Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) CreateSession(request *CreateSessionRequest) (CreateSessionResponse, error) {
	if client.CreateSessionStatusCode == 200 {
		response := new(CreateSession200Response)
		return response, nil
	}
	if client.CreateSessionStatusCode == 400 {
		response := new(CreateSession400Response)
		return response, nil
	}
	if client.CreateSessionStatusCode == 401 {
		response := new(CreateSession401Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetUsers(request *GetUsersRequest) (GetUsersResponse, error) {
	if client.GetUsersStatusCode == 200 {
		response := new(GetUsers200Response)
		return response, nil
	}
	if client.GetUsersStatusCode == 403 {
		response := new(GetUsers403Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) DeleteUser(request *DeleteUserRequest) (DeleteUserResponse, error) {
	if client.DeleteUserStatusCode == 200 {
		response := new(DeleteUser200Response)
		return response, nil
	}
	if client.DeleteUserStatusCode == 403 {
		response := new(DeleteUser403Response)
		return response, nil
	}
	if client.DeleteUserStatusCode == 404 {
		response := new(DeleteUser404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetUser(request *GetUserRequest) (GetUserResponse, error) {
	if client.GetUserStatusCode == 200 {
		response := new(GetUser200Response)
		return response, nil
	}
	if client.GetUserStatusCode == 403 {
		response := new(GetUser403Response)
		return response, nil
	}
	if client.GetUserStatusCode == 404 {
		response := new(GetUser404Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) CreateOrUpdateUser(request *CreateOrUpdateUserRequest) (CreateOrUpdateUserResponse, error) {
	if client.CreateOrUpdateUserStatusCode == 200 {
		response := new(CreateOrUpdateUser200Response)
		return response, nil
	}
	if client.CreateOrUpdateUserStatusCode == 201 {
		response := new(CreateOrUpdateUser201Response)
		return response, nil
	}
	if client.CreateOrUpdateUserStatusCode == 400 {
		response := new(CreateOrUpdateUser400Response)
		return response, nil
	}
	if client.CreateOrUpdateUserStatusCode == 403 {
		response := new(CreateOrUpdateUser403Response)
		return response, nil
	}
	if client.CreateOrUpdateUserStatusCode == 405 {
		response := new(CreateOrUpdateUser405Response)
		return response, nil
	}
	return nil, nil
}

// Get booking of session owner
func (client *VisAdminClientMock) GetBooking(request *GetBookingRequest) (GetBookingResponse, error) {
	if client.GetBookingStatusCode == 200 {
		response := new(GetBooking200Response)
		return response, nil
	}
	if client.GetBookingStatusCode == 400 {
		response := new(GetBooking400Response)
		return response, nil
	}
	if client.GetBookingStatusCode == 401 {
		response := new(GetBooking401Response)
		return response, nil
	}
	if client.GetBookingStatusCode == 404 {
		response := new(GetBooking404Response)
		return response, nil
	}
	if client.GetBookingStatusCode == 500 {
		response := new(GetBooking500Response)
		return response, nil
	}
	return nil, nil
}

// Get bookings of session owner
func (client *VisAdminClientMock) GetBookings(request *GetBookingsRequest) (GetBookingsResponse, error) {
	if client.GetBookingsStatusCode == 200 {
		response := new(GetBookings200Response)
		return response, nil
	}
	if client.GetBookingsStatusCode == 400 {
		response := new(GetBookings400Response)
		return response, nil
	}
	if client.GetBookingsStatusCode == 401 {
		response := new(GetBookings401Response)
		return response, nil
	}
	if client.GetBookingsStatusCode == 404 {
		response := new(GetBookings404Response)
		return response, nil
	}
	if client.GetBookingsStatusCode == 500 {
		response := new(GetBookings500Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) ListModels(request *ListModelsRequest) (ListModelsResponse, error) {
	if client.ListModelsStatusCode == 200 {
		data := "{\"drive_concept\":\"drive_concept\",\"price\":38,\"technical_information\":null}"
		response := new(ListModels200Response)
		responseBody := &response.Body
		err := json.Unmarshal([]byte(data), responseBody)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) GetClasses(request *GetClassesRequest) (GetClassesResponse, error) {
	if client.GetClassesStatusCode == 200 {
		response := new(GetClasses200Response)
		return response, nil
	}
	if client.GetClassesStatusCode == 400 {
		response := new(GetClasses400Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) Code(request *CodeRequest) (CodeResponse, error) {
	if client.CodeStatusCode == 200 {
		response := new(Code200Response)
		return response, nil
	}
	if client.CodeStatusCode == 400 {
		response := new(Code400Response)
		return response, nil
	}
	if client.CodeStatusCode == 401 {
		response := new(Code401Response)
		return response, nil
	}
	if client.CodeStatusCode == 404 {
		response := new(Code404Response)
		return response, nil
	}
	if client.CodeStatusCode == 500 {
		response := new(Code500Response)
		return response, nil
	}
	return nil, nil
}

/*
Deletes the user session matching the *X-Auth* header.
*/
func (client *VisAdminClientMock) DeleteCustomerSession(request *DeleteCustomerSessionRequest) (DeleteCustomerSessionResponse, error) {
	if client.DeleteCustomerSessionStatusCode == 204 {
		response := new(DeleteCustomerSession204Response)
		return response, nil
	}
	if client.DeleteCustomerSessionStatusCode == 401 {
		response := new(DeleteCustomerSession401Response)
		return response, nil
	}
	if client.DeleteCustomerSessionStatusCode == 500 {
		response := new(DeleteCustomerSession500Response)
		return response, nil
	}
	return nil, nil
}

/*
Creates a customer session for a given OpenID authentication token.
*/
func (client *VisAdminClientMock) CreateCustomerSession(request *CreateCustomerSessionRequest) (CreateCustomerSessionResponse, error) {
	if client.CreateCustomerSessionStatusCode == 201 {
		response := new(CreateCustomerSession201Response)
		return response, nil
	}
	if client.CreateCustomerSessionStatusCode == 401 {
		response := new(CreateCustomerSession401Response)
		return response, nil
	}
	if client.CreateCustomerSessionStatusCode == 403 {
		response := new(CreateCustomerSession403Response)
		return response, nil
	}
	if client.CreateCustomerSessionStatusCode == 422 {
		response := new(CreateCustomerSession422Response)
		return response, nil
	}
	if client.CreateCustomerSessionStatusCode == 500 {
		response := new(CreateCustomerSession500Response)
		return response, nil
	}
	return nil, nil
}

/*
Downloads a file that is a property within a nested structure in the response body
*/
func (client *VisAdminClientMock) DownloadNestedFile(request *DownloadNestedFileRequest) (DownloadNestedFileResponse, error) {
	if client.DownloadNestedFileStatusCode == 200 {
		response := new(DownloadNestedFile200Response)
		return response, nil
	}
	return nil, nil
}

// Retrieve a image
func (client *VisAdminClientMock) DownloadImage(request *DownloadImageRequest) (DownloadImageResponse, error) {
	if client.DownloadImageStatusCode == 200 {
		response := new(DownloadImage200Response)
		return response, nil
	}
	if client.DownloadImageStatusCode == 500 {
		response := new(DownloadImage500Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) ListElements(request *ListElementsRequest) (ListElementsResponse, error) {
	if client.ListElementsStatusCode == 200 {
		response := new(ListElements200Response)
		return response, nil
	}
	if client.ListElementsStatusCode == 500 {
		response := new(ListElements500Response)
		return response, nil
	}
	return nil, nil
}

// Retrieve a file
func (client *VisAdminClientMock) DownloadFile(request *DownloadFileRequest) (DownloadFileResponse, error) {
	if client.DownloadFileStatusCode == 200 {
		response := new(DownloadFile200Response)
		return response, nil
	}
	return nil, nil
}

// Retrieve a file
func (client *VisAdminClientMock) GenericFileDownload(request *GenericFileDownloadRequest) (GenericFileDownloadResponse, error) {
	if client.GenericFileDownloadStatusCode == 200 {
		response := new(GenericFileDownload200Response)
		return response, nil
	}
	if client.GenericFileDownloadStatusCode == 500 {
		response := new(GenericFileDownload500Response)
		return response, nil
	}
	return nil, nil
}

// get rental
func (client *VisAdminClientMock) GetRental(request *GetRentalRequest) (GetRentalResponse, error) {
	if client.GetRentalStatusCode == 200 {
		response := new(GetRental200Response)
		return response, nil
	}
	if client.GetRentalStatusCode == 400 {
		response := new(GetRental400Response)
		return response, nil
	}
	return nil, nil
}

func (client *VisAdminClientMock) PostUpload(request *PostUploadRequest) (PostUploadResponse, error) {
	if client.PostUploadStatusCode == 200 {
		response := new(PostUpload200Response)
		return response, nil
	}
	if client.PostUploadStatusCode == 500 {
		response := new(PostUpload500Response)
		return response, nil
	}
	return nil, nil
}

type VisAdminClientMock struct {
	GetClientsStatusCode             int
	DeleteClientStatusCode           int
	GetClientStatusCode              int
	CreateOrUpdateClientStatusCode   int
	GetViewsSetsStatusCode           int
	DeleteViewsSetStatusCode         int
	GetViewsSetStatusCode            int
	ActivateViewsSetStatusCode       int
	CreateOrUpdateViewsSetStatusCode int
	ShowVehicleInViewStatusCode      int
	GetPermissionsStatusCode         int
	DestroySessionStatusCode         int
	GetUserInfoStatusCode            int
	CreateSessionStatusCode          int
	GetUsersStatusCode               int
	DeleteUserStatusCode             int
	GetUserStatusCode                int
	CreateOrUpdateUserStatusCode     int
	GetBookingStatusCode             int
	GetBookingsStatusCode            int
	ListModelsStatusCode             int
	GetClassesStatusCode             int
	CodeStatusCode                   int
	DeleteCustomerSessionStatusCode  int
	CreateCustomerSessionStatusCode  int
	DownloadNestedFileStatusCode     int
	DownloadImageStatusCode          int
	ListElementsStatusCode           int
	DownloadFileStatusCode           int
	GenericFileDownloadStatusCode    int
	GetRentalStatusCode              int
	PostUploadStatusCode             int
}
