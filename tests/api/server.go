package api

import (
	"context"
	"encoding/json"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing"
	validator "github.com/go-playground/validator"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

func NewVisAdminServer(options *ServerOpts) *VisAdminServer {
	serverWrapper := &VisAdminServer{Server: newServer(options), Validator: NewValidation()}
	serverWrapper.Server.SwaggerSpec = swagger
	serverWrapper.registerValidators()
	return serverWrapper
}

type VisAdminServer struct {
	*Server
	Validator                     *Validator
	getClientsHandler             *getClientsHandlerRoute
	deleteClientHandler           *deleteClientHandlerRoute
	getClientHandler              *getClientHandlerRoute
	createOrUpdateClientHandler   *createOrUpdateClientHandlerRoute
	getViewsSetsHandler           *getViewsSetsHandlerRoute
	deleteViewsSetHandler         *deleteViewsSetHandlerRoute
	getViewsSetHandler            *getViewsSetHandlerRoute
	activateViewsSetHandler       *activateViewsSetHandlerRoute
	createOrUpdateViewsSetHandler *createOrUpdateViewsSetHandlerRoute
	showVehicleInViewHandler      *showVehicleInViewHandlerRoute
	getPermissionsHandler         *getPermissionsHandlerRoute
	destroySessionHandler         *destroySessionHandlerRoute
	getUserInfoHandler            *getUserInfoHandlerRoute
	createSessionHandler          *createSessionHandlerRoute
	getUsersHandler               *getUsersHandlerRoute
	deleteUserHandler             *deleteUserHandlerRoute
	getUserHandler                *getUserHandlerRoute
	createOrUpdateUserHandler     *createOrUpdateUserHandlerRoute
	getBookingHandler             *getBookingHandlerRoute
	getBookingsHandler            *getBookingsHandlerRoute
	listModelsHandler             *listModelsHandlerRoute
	getClassesHandler             *getClassesHandlerRoute
	codeHandler                   *codeHandlerRoute
	deleteCustomerSessionHandler  *deleteCustomerSessionHandlerRoute
	createCustomerSessionHandler  *createCustomerSessionHandlerRoute
	downloadNestedFileHandler     *downloadNestedFileHandlerRoute
	downloadImageHandler          *downloadImageHandlerRoute
	listElementsHandler           *listElementsHandlerRoute
	fileUploadHandler             *fileUploadHandlerRoute
	downloadFileHandler           *downloadFileHandlerRoute
	findByTagsHandler             *findByTagsHandlerRoute
	genericFileDownloadHandler    *genericFileDownloadHandlerRoute
	getRentalHandler              *getRentalHandlerRoute
	getShoesHandler               *getShoesHandlerRoute
	postUploadHandler             *postUploadHandlerRoute
}
type GetClientsHandler func(ctx context.Context, request *GetClientsRequest) GetClientsResponse

type getClientsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetClientsHandler
}

func (server *VisAdminServer) SetGetClientsHandler(handler GetClientsHandler, middleware ...Middleware) {
	server.getClientsHandler = &getClientsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/client", Handler: server.GetClientsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetClientsHandler(c *routing.Context) error {
	if server.getClientsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetClients (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetClientsRequest)
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetClients (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClients (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getClientsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetClients (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClients (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type DeleteClientHandler func(ctx context.Context, request *DeleteClientRequest) DeleteClientResponse

type deleteClientHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteClientHandler
}

func (server *VisAdminServer) SetDeleteClientHandler(handler DeleteClientHandler, middleware ...Middleware) {
	server.deleteClientHandler = &deleteClientHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/api/client/<clientId>", Handler: server.DeleteClientHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DeleteClientHandler(c *routing.Context) error {
	if server.deleteClientHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteClient (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteClientRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteClient (DELETE) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteClient (DELETE) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteClient (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteClientHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteClient (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteClient (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetClientHandler func(ctx context.Context, request *GetClientRequest) GetClientResponse

type getClientHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetClientHandler
}

func (server *VisAdminServer) SetGetClientHandler(handler GetClientHandler, middleware ...Middleware) {
	server.getClientHandler = &getClientHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/client/<clientId>", Handler: server.GetClientHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetClientHandler(c *routing.Context) error {
	if server.getClientHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetClient (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetClientRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClient (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetClient (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClient (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getClientHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetClient (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClient (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type CreateOrUpdateClientHandler func(ctx context.Context, request *CreateOrUpdateClientRequest) CreateOrUpdateClientResponse

type createOrUpdateClientHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CreateOrUpdateClientHandler
}

func (server *VisAdminServer) SetCreateOrUpdateClientHandler(handler CreateOrUpdateClientHandler, middleware ...Middleware) {
	server.createOrUpdateClientHandler = &createOrUpdateClientHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "PUT", Path: "/api/client/<clientId>", Handler: server.CreateOrUpdateClientHandler, Middleware: middleware}}
}

func (server *VisAdminServer) CreateOrUpdateClientHandler(c *routing.Context) error {
	if server.createOrUpdateClientHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: CreateOrUpdateClient (PUT) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CreateOrUpdateClientRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.Body, true)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateClient (PUT) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateClient (PUT) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateClient (PUT) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateClient (PUT) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateClient (PUT) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.createOrUpdateClientHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: CreateOrUpdateClient (PUT) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateClient (PUT) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetViewsSetsHandler func(ctx context.Context, request *GetViewsSetsRequest) GetViewsSetsResponse

type getViewsSetsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetViewsSetsHandler
}

func (server *VisAdminServer) SetGetViewsSetsHandler(handler GetViewsSetsHandler, middleware ...Middleware) {
	server.getViewsSetsHandler = &getViewsSetsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/client/<clientId>/views", Handler: server.GetViewsSetsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetViewsSetsHandler(c *routing.Context) error {
	if server.getViewsSetsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetViewsSets (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetViewsSetsRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSets (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSets (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSets (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getViewsSetsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetViewsSets (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSets (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type DeleteViewsSetHandler func(ctx context.Context, request *DeleteViewsSetRequest) DeleteViewsSetResponse

type deleteViewsSetHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteViewsSetHandler
}

func (server *VisAdminServer) SetDeleteViewsSetHandler(handler DeleteViewsSetHandler, middleware ...Middleware) {
	server.deleteViewsSetHandler = &deleteViewsSetHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/api/client/<clientId>/views/<viewsId>", Handler: server.DeleteViewsSetHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DeleteViewsSetHandler(c *routing.Context) error {
	if server.deleteViewsSetHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteViewsSet (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteViewsSetRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteViewsSet (DELETE) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("viewsId"), &request.ViewsId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteViewsSet (DELETE) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteViewsSet (DELETE) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteViewsSet (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteViewsSetHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteViewsSet (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteViewsSet (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetViewsSetHandler func(ctx context.Context, request *GetViewsSetRequest) GetViewsSetResponse

type getViewsSetHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetViewsSetHandler
}

func (server *VisAdminServer) SetGetViewsSetHandler(handler GetViewsSetHandler, middleware ...Middleware) {
	server.getViewsSetHandler = &getViewsSetHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/client/<clientId>/views/<viewsId>", Handler: server.GetViewsSetHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetViewsSetHandler(c *routing.Context) error {
	if server.getViewsSetHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetViewsSet (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetViewsSetRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSet (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("viewsId"), &request.ViewsId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSet (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSet (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["page"]) > 0 {
			if err := fromString(c.Request.URL.Query()["page"][0], &request.Page); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSet (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSet (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getViewsSetHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetViewsSet (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetViewsSet (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Make this viewset the active one for the client.
type ActivateViewsSetHandler func(ctx context.Context, request *ActivateViewsSetRequest) ActivateViewsSetResponse

type activateViewsSetHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    ActivateViewsSetHandler
}

func (server *VisAdminServer) SetActivateViewsSetHandler(handler ActivateViewsSetHandler, middleware ...Middleware) {
	server.activateViewsSetHandler = &activateViewsSetHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/api/client/<clientId>/views/<viewsId>", Handler: server.ActivateViewsSetHandler, Middleware: middleware}}
}

func (server *VisAdminServer) ActivateViewsSetHandler(c *routing.Context) error {
	if server.activateViewsSetHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: ActivateViewsSet (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(ActivateViewsSetRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ActivateViewsSet (POST) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("viewsId"), &request.ViewsId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ActivateViewsSet (POST) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ActivateViewsSet (POST) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ActivateViewsSet (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.activateViewsSetHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: ActivateViewsSet (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ActivateViewsSet (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type CreateOrUpdateViewsSetHandler func(ctx context.Context, request *CreateOrUpdateViewsSetRequest) CreateOrUpdateViewsSetResponse

type createOrUpdateViewsSetHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CreateOrUpdateViewsSetHandler
}

func (server *VisAdminServer) SetCreateOrUpdateViewsSetHandler(handler CreateOrUpdateViewsSetHandler, middleware ...Middleware) {
	server.createOrUpdateViewsSetHandler = &createOrUpdateViewsSetHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "PUT", Path: "/api/client/<clientId>/views/<viewsId>", Handler: server.CreateOrUpdateViewsSetHandler, Middleware: middleware}}
}

func (server *VisAdminServer) CreateOrUpdateViewsSetHandler(c *routing.Context) error {
	if server.createOrUpdateViewsSetHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: CreateOrUpdateViewsSet (PUT) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CreateOrUpdateViewsSetRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.Body, true)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("viewsId"), &request.ViewsId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.createOrUpdateViewsSetHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: CreateOrUpdateViewsSet (PUT) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateViewsSet (PUT) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type ShowVehicleInViewHandler func(ctx context.Context, request *ShowVehicleInViewRequest) ShowVehicleInViewResponse

type showVehicleInViewHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    ShowVehicleInViewHandler
}

func (server *VisAdminServer) SetShowVehicleInViewHandler(handler ShowVehicleInViewHandler, middleware ...Middleware) {
	server.showVehicleInViewHandler = &showVehicleInViewHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/client/<clientId>/views/<viewsId>/<view>/<breakpoint>/<spec>", Handler: server.ShowVehicleInViewHandler, Middleware: middleware}}
}

func (server *VisAdminServer) ShowVehicleInViewHandler(c *routing.Context) error {
	if server.showVehicleInViewHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: ShowVehicleInView (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(ShowVehicleInViewRequest)
		if err := fromString(c.Param("clientId"), &request.ClientId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("viewsId"), &request.ViewsId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("view"), &request.View); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("breakpoint"), &request.Breakpoint); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if err := fromString(c.Param("spec"), &request.Spec); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.showVehicleInViewHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: ShowVehicleInView (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ShowVehicleInView (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

/*
Get the list of permissions
a user can grant to other users.
*/
type GetPermissionsHandler func(ctx context.Context, request *GetPermissionsRequest) GetPermissionsResponse

type getPermissionsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetPermissionsHandler
}

func (server *VisAdminServer) SetGetPermissionsHandler(handler GetPermissionsHandler, middleware ...Middleware) {
	server.getPermissionsHandler = &getPermissionsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/permission", Handler: server.GetPermissionsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetPermissionsHandler(c *routing.Context) error {
	if server.getPermissionsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetPermissions (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetPermissionsRequest)
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetPermissions (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetPermissions (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getPermissionsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetPermissions (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetPermissions (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type DestroySessionHandler func(ctx context.Context, request *DestroySessionRequest) DestroySessionResponse

type destroySessionHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DestroySessionHandler
}

func (server *VisAdminServer) SetDestroySessionHandler(handler DestroySessionHandler, middleware ...Middleware) {
	server.destroySessionHandler = &destroySessionHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/api/session", Handler: server.DestroySessionHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DestroySessionHandler(c *routing.Context) error {
	if server.destroySessionHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DestroySession (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DestroySessionRequest)
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: DestroySession (DELETE) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DestroySession (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.destroySessionHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DestroySession (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DestroySession (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetUserInfoHandler func(ctx context.Context, request *GetUserInfoRequest) GetUserInfoResponse

type getUserInfoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetUserInfoHandler
}

func (server *VisAdminServer) SetGetUserInfoHandler(handler GetUserInfoHandler, middleware ...Middleware) {
	server.getUserInfoHandler = &getUserInfoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/session", Handler: server.GetUserInfoHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetUserInfoHandler(c *routing.Context) error {
	if server.getUserInfoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetUserInfo (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetUserInfoRequest)
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetUserInfo (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["subID"]) > 0 {
			if err := fromString(c.Request.Header["subID"][0], &request.SubID); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetUserInfo (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUserInfo (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			c.Response.Header()[contentTypeHeader] = []string{contentTypeApplicationJson}
			c.Response.WriteHeader(http.StatusBadRequest)
			encodeErr := json.NewEncoder(c.Response).Encode(validationErrors)
			if encodeErr != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetUserInfo (GET) could not encode validation response (error: %v)", encodeErr))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			return nil
		}
		response := server.getUserInfoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetUserInfo (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUserInfo (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type CreateSessionHandler func(ctx context.Context, request *CreateSessionRequest) CreateSessionResponse

type createSessionHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CreateSessionHandler
}

func (server *VisAdminServer) SetCreateSessionHandler(handler CreateSessionHandler, middleware ...Middleware) {
	server.createSessionHandler = &createSessionHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/api/session", Handler: server.CreateSessionHandler, Middleware: middleware}}
}

func (server *VisAdminServer) CreateSessionHandler(c *routing.Context) error {
	if server.createSessionHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: CreateSession (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CreateSessionRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.Body, true)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateSession (POST) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateSession (POST) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateSession (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			c.Response.Header()[contentTypeHeader] = []string{contentTypeApplicationJson}
			c.Response.WriteHeader(http.StatusBadRequest)
			encodeErr := json.NewEncoder(c.Response).Encode(validationErrors)
			if encodeErr != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateSession (POST) could not encode validation response (error: %v)", encodeErr))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			return nil
		}
		response := server.createSessionHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: CreateSession (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateSession (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetUsersHandler func(ctx context.Context, request *GetUsersRequest) GetUsersResponse

type getUsersHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetUsersHandler
}

func (server *VisAdminServer) SetGetUsersHandler(handler GetUsersHandler, middleware ...Middleware) {
	server.getUsersHandler = &getUsersHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/user", Handler: server.GetUsersHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetUsersHandler(c *routing.Context) error {
	if server.getUsersHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetUsers (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetUsersRequest)
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetUsers (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUsers (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getUsersHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetUsers (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUsers (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type DeleteUserHandler func(ctx context.Context, request *DeleteUserRequest) DeleteUserResponse

type deleteUserHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteUserHandler
}

func (server *VisAdminServer) SetDeleteUserHandler(handler DeleteUserHandler, middleware ...Middleware) {
	server.deleteUserHandler = &deleteUserHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/api/user/<userId>", Handler: server.DeleteUserHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DeleteUserHandler(c *routing.Context) error {
	if server.deleteUserHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteUser (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteUserRequest)
		if err := fromString(c.Param("userId"), &request.UserId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteUser (DELETE) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteUser (DELETE) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["allKeys"]) > 0 {
			if err := fromString(c.Request.URL.Query()["allKeys"][0], &request.AllKeys); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteUser (DELETE) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteUser (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteUserHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteUser (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteUser (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetUserHandler func(ctx context.Context, request *GetUserRequest) GetUserResponse

type getUserHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetUserHandler
}

func (server *VisAdminServer) SetGetUserHandler(handler GetUserHandler, middleware ...Middleware) {
	server.getUserHandler = &getUserHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/api/user/<userId>", Handler: server.GetUserHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetUserHandler(c *routing.Context) error {
	if server.getUserHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetUser (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetUserRequest)
		if err := fromString(c.Param("userId"), &request.UserId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUser (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetUser (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["allKeys"]) > 0 {
			if err := fromString(c.Request.URL.Query()["allKeys"][0], &request.AllKeys); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetUser (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUser (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getUserHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetUser (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetUser (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type CreateOrUpdateUserHandler func(ctx context.Context, request *CreateOrUpdateUserRequest) CreateOrUpdateUserResponse

type createOrUpdateUserHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CreateOrUpdateUserHandler
}

func (server *VisAdminServer) SetCreateOrUpdateUserHandler(handler CreateOrUpdateUserHandler, middleware ...Middleware) {
	server.createOrUpdateUserHandler = &createOrUpdateUserHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "PUT", Path: "/api/user/<userId>", Handler: server.CreateOrUpdateUserHandler, Middleware: middleware}}
}

func (server *VisAdminServer) CreateOrUpdateUserHandler(c *routing.Context) error {
	if server.createOrUpdateUserHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: CreateOrUpdateUser (PUT) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CreateOrUpdateUserRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.Body, true)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		if err := fromString(c.Param("userId"), &request.UserId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.Header["X-Auth"]) > 0 {
			if err := fromString(c.Request.Header["X-Auth"][0], &request.XAuth); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["allKeys"]) > 0 {
			if err := fromString(c.Request.URL.Query()["allKeys"][0], &request.AllKeys); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.createOrUpdateUserHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: CreateOrUpdateUser (PUT) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateOrUpdateUser (PUT) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Get booking of session owner
type GetBookingHandler func(ctx context.Context, request *GetBookingRequest) GetBookingResponse

type getBookingHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetBookingHandler
}

func (server *VisAdminServer) SetGetBookingHandler(handler GetBookingHandler, middleware ...Middleware) {
	server.getBookingHandler = &getBookingHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/booking", Handler: server.GetBookingHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetBookingHandler(c *routing.Context) error {
	return newNotSupportedContentType(415, "no supported content type")
}

// Get bookings of session owner
type GetBookingsHandler func(ctx context.Context, request *GetBookingsRequest) GetBookingsResponse

type getBookingsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetBookingsHandler
}

func (server *VisAdminServer) SetGetBookingsHandler(handler GetBookingsHandler, middleware ...Middleware) {
	server.getBookingsHandler = &getBookingsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/bookings", Handler: server.GetBookingsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetBookingsHandler(c *routing.Context) error {
	if server.getBookingsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetBookings (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetBookingsRequest)
		request.XSessionID = c.Request.Header.Get("X-Session-ID")
		if len(c.Request.Header["date"]) > 0 {
			if err := fromString(c.Request.Header["date"][0], &request.Date); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetBookings (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		if len(c.Request.URL.Query()["ids"]) > 0 {
			if err := fromString(c.Request.URL.Query()["ids"][0], &request.Ids); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetBookings (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetBookings (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getBookingsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetBookings (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetBookings (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type ListModelsHandler func(ctx context.Context, request *ListModelsRequest) ListModelsResponse

type listModelsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    ListModelsHandler
}

func (server *VisAdminServer) SetListModelsHandler(handler ListModelsHandler, middleware ...Middleware) {
	server.listModelsHandler = &listModelsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/brands/<brandId>/models", Handler: server.ListModelsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) ListModelsHandler(c *routing.Context) error {
	if server.listModelsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: ListModels (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(ListModelsRequest)
		if err := fromString(c.Param("brandId"), &request.BrandId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["driveConcept"]) > 0 {
			if err := fromString(c.Request.URL.Query()["driveConcept"][0], &request.DriveConcept); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		if len(c.Request.URL.Query()["languageId"]) > 0 {
			if err := fromString(c.Request.URL.Query()["languageId"][0], &request.LanguageId); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		if len(c.Request.URL.Query()["classId"]) > 0 {
			if err := fromString(c.Request.URL.Query()["classId"][0], &request.ClassId); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		if len(c.Request.URL.Query()["lineId"]) > 0 {
			if err := fromString(c.Request.URL.Query()["lineId"][0], &request.LineId); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		if len(c.Request.URL.Query()["ids"]) > 0 {
			if err := fromString(c.Request.URL.Query()["ids"][0], &request.Ids); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.listModelsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: ListModels (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListModels (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetClassesHandler func(ctx context.Context, request *GetClassesRequest) GetClassesResponse

type getClassesHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetClassesHandler
}

func (server *VisAdminServer) SetGetClassesHandler(handler GetClassesHandler, middleware ...Middleware) {
	server.getClassesHandler = &getClassesHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/classes/<productGroup>", Handler: server.GetClassesHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetClassesHandler(c *routing.Context) error {
	if server.getClassesHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetClasses (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetClassesRequest)
		if err := fromString(c.Param("productGroup"), &request.ProductGroup); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClasses (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		if len(c.Request.URL.Query()["componentTypes"]) > 0 {
			if err := fromString(c.Request.URL.Query()["componentTypes"][0], &request.ComponentTypes); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetClasses (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClasses (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getClassesHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetClasses (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetClasses (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type CodeHandler func(ctx context.Context, request *CodeRequest) CodeResponse

type codeHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CodeHandler
}

func (server *VisAdminServer) SetCodeHandler(handler CodeHandler, middleware ...Middleware) {
	server.codeHandler = &codeHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/code", Handler: server.CodeHandler, Middleware: middleware}}
}

func (server *VisAdminServer) CodeHandler(c *routing.Context) error {
	if server.codeHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: Code (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CodeRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeMultipartFormData {
		} else if contentTypeOfResponse == contentTypeApplicationFormUrlencoded {
			rawBody, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not read body of incoming request (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			queryInBody, err := url.ParseQuery(string(rawBody))
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not parse raw query string of incoming request (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			if len(queryInBody["state"]) > 0 {
				if err := fromString(queryInBody.Get("state"), &request.State); err != nil {
					server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not convert string to specific type (error: %v)", err))
					return NewHTTPStatusCodeError(http.StatusBadRequest)
				}
			}
			if len(queryInBody["response_mode"]) > 0 {
				if err := fromString(queryInBody.Get("response_mode"), &request.ResponseMode); err != nil {
					server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not convert string to specific type (error: %v)", err))
					return NewHTTPStatusCodeError(http.StatusBadRequest)
				}
			}
			if len(queryInBody["code"]) > 0 {
				if err := fromString(queryInBody.Get("code"), &request.Code); err != nil {
					server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not convert string to specific type (error: %v)", err))
					return NewHTTPStatusCodeError(http.StatusBadRequest)
				}
			} else {
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		if len(c.Request.URL.Query()["session"]) > 0 {
			if err := fromString(c.Request.URL.Query()["session"][0], &request.Session); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.codeHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: Code (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: Code (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

/*
Deletes the user session matching the *X-Auth* header.
*/
type DeleteCustomerSessionHandler func(ctx context.Context, request *DeleteCustomerSessionRequest) DeleteCustomerSessionResponse

type deleteCustomerSessionHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteCustomerSessionHandler
}

func (server *VisAdminServer) SetDeleteCustomerSessionHandler(handler DeleteCustomerSessionHandler, middleware ...Middleware) {
	server.deleteCustomerSessionHandler = &deleteCustomerSessionHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/customer/session", Handler: server.DeleteCustomerSessionHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DeleteCustomerSessionHandler(c *routing.Context) error {
	if server.deleteCustomerSessionHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteCustomerSession (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteCustomerSessionRequest)
		request.XSessionID = c.Request.Header.Get("X-Session-ID")
		if len(c.Request.Header["X-Request-ID"]) > 0 {
			if err := fromString(c.Request.Header["X-Request-ID"][0], &request.XRequestID); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteCustomerSession (DELETE) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteCustomerSession (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteCustomerSessionHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteCustomerSession (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteCustomerSession (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

/*
Creates a customer session for a given OpenID authentication token.
*/
type CreateCustomerSessionHandler func(ctx context.Context, request *CreateCustomerSessionRequest) CreateCustomerSessionResponse

type createCustomerSessionHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    CreateCustomerSessionHandler
}

func (server *VisAdminServer) SetCreateCustomerSessionHandler(handler CreateCustomerSessionHandler, middleware ...Middleware) {
	server.createCustomerSessionHandler = &createCustomerSessionHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/customer/session", Handler: server.CreateCustomerSessionHandler, Middleware: middleware}}
}

func (server *VisAdminServer) CreateCustomerSessionHandler(c *routing.Context) error {
	if server.createCustomerSessionHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: CreateCustomerSession (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(CreateCustomerSessionRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeMultipartFormData {
		} else if contentTypeOfResponse == contentTypeApplicationFormUrlencoded {
			rawBody, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not read body of incoming request (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			queryInBody, err := url.ParseQuery(string(rawBody))
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not parse raw query string of incoming request (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			if len(queryInBody["code"]) > 0 {
				if err := fromString(queryInBody.Get("code"), &request.Code); err != nil {
					server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not convert string to specific type (error: %v)", err))
					return NewHTTPStatusCodeError(http.StatusBadRequest)
				}
			} else {
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
			if len(queryInBody["locale"]) > 0 {
				if err := fromString(queryInBody.Get("locale"), &request.Locale); err != nil {
					server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not convert string to specific type (error: %v)", err))
					return NewHTTPStatusCodeError(http.StatusBadRequest)
				}
			}
		} else {
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		if len(c.Request.Header["X-Request-ID"]) > 0 {
			if err := fromString(c.Request.Header["X-Request-ID"][0], &request.XRequestID); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.createCustomerSessionHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: CreateCustomerSession (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: CreateCustomerSession (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

/*
Downloads a file that is a property within a nested structure in the response body
*/
type DownloadNestedFileHandler func(ctx context.Context, request *DownloadNestedFileRequest) DownloadNestedFileResponse

type downloadNestedFileHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DownloadNestedFileHandler
}

func (server *VisAdminServer) SetDownloadNestedFileHandler(handler DownloadNestedFileHandler, middleware ...Middleware) {
	server.downloadNestedFileHandler = &downloadNestedFileHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/download/nested/file", Handler: server.DownloadNestedFileHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DownloadNestedFileHandler(c *routing.Context) error {
	if server.downloadNestedFileHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DownloadNestedFile (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DownloadNestedFileRequest)
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DownloadNestedFile (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.downloadNestedFileHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DownloadNestedFile (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DownloadNestedFile (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Retrieve a image
type DownloadImageHandler func(ctx context.Context, request *DownloadImageRequest) DownloadImageResponse

type downloadImageHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DownloadImageHandler
}

func (server *VisAdminServer) SetDownloadImageHandler(handler DownloadImageHandler, middleware ...Middleware) {
	server.downloadImageHandler = &downloadImageHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/download/<image>", Handler: server.DownloadImageHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DownloadImageHandler(c *routing.Context) error {
	if server.downloadImageHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DownloadImage (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DownloadImageRequest)
		if err := fromString(c.Param("image"), &request.Image); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DownloadImage (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DownloadImage (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.downloadImageHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DownloadImage (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DownloadImage (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type ListElementsHandler func(ctx context.Context, request *ListElementsRequest) ListElementsResponse

type listElementsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    ListElementsHandler
}

func (server *VisAdminServer) SetListElementsHandler(handler ListElementsHandler, middleware ...Middleware) {
	server.listElementsHandler = &listElementsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/elements", Handler: server.ListElementsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) ListElementsHandler(c *routing.Context) error {
	if server.listElementsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: ListElements (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(ListElementsRequest)
		if len(c.Request.URL.Query()["_page"]) > 0 {
			if err := fromString(c.Request.URL.Query()["_page"][0], &request.Page); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListElements (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			request.Page = new(int64)
			*request.Page = 1
		}
		if len(c.Request.URL.Query()["_perPage"]) > 0 {
			if err := fromString(c.Request.URL.Query()["_perPage"][0], &request.PerPage); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: ListElements (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			request.PerPage = new(int64)
			*request.PerPage = 10
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListElements (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.listElementsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: ListElements (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListElements (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type FileUploadHandler func(ctx context.Context, request *FileUploadRequest) FileUploadResponse

type fileUploadHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    FileUploadHandler
}

func (server *VisAdminServer) SetFileUploadHandler(handler FileUploadHandler, middleware ...Middleware) {
	server.fileUploadHandler = &fileUploadHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/file-upload", Handler: server.FileUploadHandler, Middleware: middleware}}
}

func (server *VisAdminServer) FileUploadHandler(c *routing.Context) error {
	if server.fileUploadHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: FileUpload (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(FileUploadRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeMultipartFormData {
			formData := &request.FormData
			file0, extractErr0 := extractUpload("file", c.Request)
			if extractErr0 != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: FileUpload (POST) could not extract upload from incoming request (error: %v)", extractErr0))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
			formData.File = file0
		} else if contentTypeOfResponse == contentTypeApplicationFormUrlencoded {
		} else {
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: FileUpload (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.fileUploadHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: FileUpload (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: FileUpload (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Retrieve a file
type DownloadFileHandler func(ctx context.Context, request *DownloadFileRequest) DownloadFileResponse

type downloadFileHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DownloadFileHandler
}

func (server *VisAdminServer) SetDownloadFileHandler(handler DownloadFileHandler, middleware ...Middleware) {
	server.downloadFileHandler = &downloadFileHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/filedownload/<file>", Handler: server.DownloadFileHandler, Middleware: middleware}}
}

func (server *VisAdminServer) DownloadFileHandler(c *routing.Context) error {
	return newNotSupportedContentType(415, "no supported content type")
}

// Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.
type FindByTagsHandler func(ctx context.Context, request *FindByTagsRequest) FindByTagsResponse

type findByTagsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    FindByTagsHandler
}

func (server *VisAdminServer) SetFindByTagsHandler(handler FindByTagsHandler, middleware ...Middleware) {
	server.findByTagsHandler = &findByTagsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/findByTags", Handler: server.FindByTagsHandler, Middleware: middleware}}
}

func (server *VisAdminServer) FindByTagsHandler(c *routing.Context) error {
	if server.findByTagsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: FindByTags (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(FindByTagsRequest)
		if len(c.Request.URL.Query()["tags"]) > 0 {
			if err := fromString(c.Request.URL.Query()["tags"][0], &request.Tags); err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: FindByTags (GET) could not convert string to specific type (error: %v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
			// minItems validator constrain at least 2 items
			if len(request.Tags) < 2 {
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
			// maxItems validator constrain at maximum 5 items
			if len(request.Tags) > 5 {
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: FindByTags (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.findByTagsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: FindByTags (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: FindByTags (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// Retrieve a file
type GenericFileDownloadHandler func(ctx context.Context, request *GenericFileDownloadRequest) GenericFileDownloadResponse

type genericFileDownloadHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GenericFileDownloadHandler
}

func (server *VisAdminServer) SetGenericFileDownloadHandler(handler GenericFileDownloadHandler, middleware ...Middleware) {
	server.genericFileDownloadHandler = &genericFileDownloadHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/generic/download/<ext>", Handler: server.GenericFileDownloadHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GenericFileDownloadHandler(c *routing.Context) error {
	if server.genericFileDownloadHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GenericFileDownload (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GenericFileDownloadRequest)
		if err := fromString(c.Param("ext"), &request.Ext); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GenericFileDownload (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GenericFileDownload (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.genericFileDownloadHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GenericFileDownload (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GenericFileDownload (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

// get rental
type GetRentalHandler func(ctx context.Context, request *GetRentalRequest) GetRentalResponse

type getRentalHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetRentalHandler
}

func (server *VisAdminServer) SetGetRentalHandler(handler GetRentalHandler, middleware ...Middleware) {
	server.getRentalHandler = &getRentalHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/rental", Handler: server.GetRentalHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetRentalHandler(c *routing.Context) error {
	if server.getRentalHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetRental (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetRentalRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.Body, true)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetRental (GET) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetRental (GET) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetRental (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			c.Response.Header()[contentTypeHeader] = []string{contentTypeApplicationJson}
			c.Response.WriteHeader(http.StatusBadRequest)
			encodeErr := json.NewEncoder(c.Response).Encode(validationErrors)
			if encodeErr != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: GetRental (GET) could not encode validation response (error: %v)", encodeErr))
				return NewHTTPStatusCodeError(http.StatusInternalServerError)
			}
			return nil
		}
		response := server.getRentalHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetRental (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetRental (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetShoesHandler func(ctx context.Context, request *GetShoesRequest) GetShoesResponse

type getShoesHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetShoesHandler
}

func (server *VisAdminServer) SetGetShoesHandler(handler GetShoesHandler, middleware ...Middleware) {
	server.getShoesHandler = &getShoesHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/shop/shoes", Handler: server.GetShoesHandler, Middleware: middleware}}
}

func (server *VisAdminServer) GetShoesHandler(c *routing.Context) error {
	if server.getShoesHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetShoes (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetShoesRequest)
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetShoes (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getShoesHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetShoes (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetShoes (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type PostUploadHandler func(ctx context.Context, request *PostUploadRequest) PostUploadResponse

type postUploadHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    PostUploadHandler
}

func (server *VisAdminServer) SetPostUploadHandler(handler PostUploadHandler, middleware ...Middleware) {
	server.postUploadHandler = &postUploadHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/upload", Handler: server.PostUploadHandler, Middleware: middleware}}
}

func (server *VisAdminServer) PostUploadHandler(c *routing.Context) error {
	if server.postUploadHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: PostUpload (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(PostUploadRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeMultipartFormData {
			formData := &request.FormData
			file0, extractErr0 := extractUpload("upfile", c.Request)
			if extractErr0 != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PostUpload (POST) could not extract upload from incoming request (error: %v)", extractErr0))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
			formData.Upfile = file0
			if len(c.Request.Form["note"]) > 0 {
				if err := fromString(c.Request.Form["note"][0], &formData.Note); err != nil {
					server.ErrorLogger(fmt.Sprintf("wrap handler: PostUpload (POST) could not convert string to specific type (error: %v)", err))
					return NewHTTPStatusCodeError(http.StatusBadRequest)
				}
			}
		} else if contentTypeOfResponse == contentTypeApplicationFormUrlencoded {
		} else {
			return newNotSupportedContentType(415, contentTypeOfResponse)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PostUpload (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.postUploadHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: PostUpload (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PostUpload (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

func (server *VisAdminServer) registerValidators() {
	regex2 := regexp.MustCompile("^([a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12})?$")
	callbackRegex2 := func(fl validator.FieldLevel) bool {
		return regex2.MatchString(fl.Field().String())
	}
	server.Validator.RegisterValidation("regex2", callbackRegex2)
	server.Validator.RegisterValidation("regex3", callbackRegex2)
	regex7 := regexp.MustCompile("^([a-z]{2})-([A-Z]{2})$")
	callbackRegex7 := func(fl validator.FieldLevel) bool {
		return regex7.MatchString(fl.Field().String())
	}
	server.Validator.RegisterValidation("regex7", callbackRegex7)
	regex5 := regexp.MustCompile("^(http|ftp|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&amp;:/~\\+#]*[\\w\\-\\@?^=%&amp;/~\\+#])?$")
	callbackRegex5 := func(fl validator.FieldLevel) bool {
		return regex5.MatchString(fl.Field().String())
	}
	server.Validator.RegisterValidation("regex5", callbackRegex5)
	server.Validator.RegisterValidation("regex6", callbackRegex5)
	regex1 := regexp.MustCompile("^[a-zA-Z]$")
	callbackRegex1 := func(fl validator.FieldLevel) bool {
		return regex1.MatchString(fl.Field().String())
	}
	server.Validator.RegisterValidation("regex1", callbackRegex1)
	server.Validator.RegisterValidation("regex4", callbackRegex1)
}

func (server *VisAdminServer) Start(port int) error {
	routes := []RouteDescription{}
	if server.getClientsHandler != nil {
		routes = append(routes, server.getClientsHandler.routeDescription)
	}
	if server.deleteClientHandler != nil {
		routes = append(routes, server.deleteClientHandler.routeDescription)
	}
	if server.getClientHandler != nil {
		routes = append(routes, server.getClientHandler.routeDescription)
	}
	if server.createOrUpdateClientHandler != nil {
		routes = append(routes, server.createOrUpdateClientHandler.routeDescription)
	}
	if server.getViewsSetsHandler != nil {
		routes = append(routes, server.getViewsSetsHandler.routeDescription)
	}
	if server.deleteViewsSetHandler != nil {
		routes = append(routes, server.deleteViewsSetHandler.routeDescription)
	}
	if server.getViewsSetHandler != nil {
		routes = append(routes, server.getViewsSetHandler.routeDescription)
	}
	if server.activateViewsSetHandler != nil {
		routes = append(routes, server.activateViewsSetHandler.routeDescription)
	}
	if server.createOrUpdateViewsSetHandler != nil {
		routes = append(routes, server.createOrUpdateViewsSetHandler.routeDescription)
	}
	if server.showVehicleInViewHandler != nil {
		routes = append(routes, server.showVehicleInViewHandler.routeDescription)
	}
	if server.getPermissionsHandler != nil {
		routes = append(routes, server.getPermissionsHandler.routeDescription)
	}
	if server.destroySessionHandler != nil {
		routes = append(routes, server.destroySessionHandler.routeDescription)
	}
	if server.getUserInfoHandler != nil {
		routes = append(routes, server.getUserInfoHandler.routeDescription)
	}
	if server.createSessionHandler != nil {
		routes = append(routes, server.createSessionHandler.routeDescription)
	}
	if server.getUsersHandler != nil {
		routes = append(routes, server.getUsersHandler.routeDescription)
	}
	if server.deleteUserHandler != nil {
		routes = append(routes, server.deleteUserHandler.routeDescription)
	}
	if server.getUserHandler != nil {
		routes = append(routes, server.getUserHandler.routeDescription)
	}
	if server.createOrUpdateUserHandler != nil {
		routes = append(routes, server.createOrUpdateUserHandler.routeDescription)
	}
	if server.getBookingHandler != nil {
		routes = append(routes, server.getBookingHandler.routeDescription)
	}
	if server.getBookingsHandler != nil {
		routes = append(routes, server.getBookingsHandler.routeDescription)
	}
	if server.listModelsHandler != nil {
		routes = append(routes, server.listModelsHandler.routeDescription)
	}
	if server.getClassesHandler != nil {
		routes = append(routes, server.getClassesHandler.routeDescription)
	}
	if server.codeHandler != nil {
		routes = append(routes, server.codeHandler.routeDescription)
	}
	if server.deleteCustomerSessionHandler != nil {
		routes = append(routes, server.deleteCustomerSessionHandler.routeDescription)
	}
	if server.createCustomerSessionHandler != nil {
		routes = append(routes, server.createCustomerSessionHandler.routeDescription)
	}
	if server.downloadNestedFileHandler != nil {
		routes = append(routes, server.downloadNestedFileHandler.routeDescription)
	}
	if server.downloadImageHandler != nil {
		routes = append(routes, server.downloadImageHandler.routeDescription)
	}
	if server.listElementsHandler != nil {
		routes = append(routes, server.listElementsHandler.routeDescription)
	}
	if server.fileUploadHandler != nil {
		routes = append(routes, server.fileUploadHandler.routeDescription)
	}
	if server.downloadFileHandler != nil {
		routes = append(routes, server.downloadFileHandler.routeDescription)
	}
	if server.findByTagsHandler != nil {
		routes = append(routes, server.findByTagsHandler.routeDescription)
	}
	if server.genericFileDownloadHandler != nil {
		routes = append(routes, server.genericFileDownloadHandler.routeDescription)
	}
	if server.getRentalHandler != nil {
		routes = append(routes, server.getRentalHandler.routeDescription)
	}
	if server.getShoesHandler != nil {
		routes = append(routes, server.getShoesHandler.routeDescription)
	}
	if server.postUploadHandler != nil {
		routes = append(routes, server.postUploadHandler.routeDescription)
	}
	return server.Server.Start(port, routes)
}

const swagger = "{\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"],\"swagger\":\"2.0\",\"info\":{\"description\":\"Vehicle Information Service Admin API\",\"title\":\"vis-admin\",\"contact\":{\"name\":\"Max Mustermann\",\"email\":\"max.musterman@fake.de\"},\"version\":\"1.0.0\"},\"paths\":{\"/api/client\":{\"get\":{\"summary\":\"List clients\",\"operationId\":\"GetClients\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Status 200\",\"schema\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Client\"}}},\"204\":{\"description\":\"Status 201\"},\"403\":{\"description\":\"Not authenticated\"}}}},\"/api/client/{clientId}\":{\"get\":{\"summary\":\"Get client\",\"operationId\":\"GetClient\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\",\"schema\":{\"$ref\":\"#/definitions/Client\"}},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"put\":{\"summary\":\"Create or update client\",\"operationId\":\"CreateOrUpdateClient\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true},{\"name\":\"body\",\"in\":\"body\",\"required\":true,\"schema\":{\"$ref\":\"#/definitions/Client\"}}],\"responses\":{\"200\":{\"description\":\"Updated\"},\"201\":{\"description\":\"Created\"},\"400\":{\"description\":\"Malformed request body\"},\"403\":{\"description\":\"Not authenticated\"},\"405\":{\"description\":\"Not allowed\"}}},\"delete\":{\"summary\":\"Delete client\",\"operationId\":\"DeleteClient\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\"},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"parameters\":[{\"type\":\"string\",\"name\":\"clientId\",\"in\":\"path\",\"required\":true}]},\"/api/client/{clientId}/views\":{\"get\":{\"summary\":\"List views sets\",\"operationId\":\"GetViewsSets\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\",\"schema\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/views%20set\"}}},\"403\":{\"description\":\"Not authenticated\"}}},\"parameters\":[{\"type\":\"string\",\"name\":\"clientId\",\"in\":\"path\",\"required\":true}]},\"/api/client/{clientId}/views/{viewsId}\":{\"get\":{\"summary\":\"Get views set\",\"operationId\":\"GetViewsSet\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true},{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"page\",\"in\":\"query\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\",\"schema\":{\"$ref\":\"#/definitions/views%20set\"}},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"put\":{\"summary\":\"Create or update views set\",\"operationId\":\"CreateOrUpdateViewsSet\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true},{\"name\":\"body\",\"in\":\"body\",\"required\":true,\"schema\":{\"$ref\":\"#/definitions/views%20set\"}}],\"responses\":{\"200\":{\"description\":\"Updated\"},\"201\":{\"description\":\"Created\"},\"400\":{\"description\":\"Malformed request body\"},\"403\":{\"description\":\"Not authenticated\"},\"405\":{\"description\":\"Not allowed\"}}},\"post\":{\"description\":\"Make this viewset the active one for the client.\",\"summary\":\"Activate views set\",\"operationId\":\"ActivateViewsSet\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\"},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"delete\":{\"summary\":\"Delete views set\",\"operationId\":\"DeleteViewsSet\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\"},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"parameters\":[{\"type\":\"string\",\"name\":\"clientId\",\"in\":\"path\",\"required\":true},{\"type\":\"string\",\"name\":\"viewsId\",\"in\":\"path\",\"required\":true}]},\"/api/client/{clientId}/views/{viewsId}/{view}/{breakpoint}/{spec}\":{\"get\":{\"summary\":\"Show vehicle in view\",\"operationId\":\"ShowVehicleInView\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\"},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"parameters\":[{\"type\":\"string\",\"name\":\"clientId\",\"in\":\"path\",\"required\":true},{\"type\":\"string\",\"name\":\"viewsId\",\"in\":\"path\",\"required\":true},{\"type\":\"string\",\"name\":\"view\",\"in\":\"path\",\"required\":true},{\"type\":\"string\",\"name\":\"breakpoint\",\"in\":\"path\",\"required\":true},{\"type\":\"string\",\"name\":\"spec\",\"in\":\"path\",\"required\":true}]},\"/api/permission\":{\"get\":{\"description\":\"Get the list of permissions\\na user can grant to other users.\",\"summary\":\"List permissions\",\"operationId\":\"GetPermissions\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Status 200\",\"schema\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}},\"403\":{\"description\":\"Not authenticated\"}}}},\"/api/session\":{\"get\":{\"tags\":[\"SESSION\"],\"summary\":\"Get user info\",\"operationId\":\"GetUserInfo\",\"parameters\":[{\"maxLength\":255,\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true},{\"maximum\":255,\"type\":\"integer\",\"description\":\"session\",\"name\":\"subID\",\"in\":\"header\"}],\"responses\":{\"200\":{\"description\":\"Status 200\",\"schema\":{\"$ref\":\"#/definitions/User\"}},\"400\":{\"description\":\"Malformed request body\",\"schema\":{\"$ref\":\"#/definitions/ValidationErrors\"}},\"403\":{\"description\":\"Not authenticatedq\"}}},\"post\":{\"tags\":[\"SESSION\"],\"summary\":\"Create session\",\"operationId\":\"CreateSession\",\"parameters\":[{\"name\":\"body\",\"in\":\"body\",\"required\":true,\"schema\":{\"type\":\"object\",\"required\":[\"id\",\"password\"],\"properties\":{\"id\":{\"type\":\"string\",\"minLength\":1},\"password\":{\"type\":\"string\",\"minLength\":1}}}}],\"responses\":{\"200\":{\"description\":\"Authentication successful\",\"headers\":{\"X-Auth\":{\"type\":\"string\",\"description\":\"Authentication token\"}}},\"400\":{\"description\":\"Malformed request body\",\"schema\":{\"$ref\":\"#/definitions/ValidationErrors\"}},\"401\":{\"description\":\"Authentication not successful\"}}},\"delete\":{\"summary\":\"Destroy session\",\"operationId\":\"DestroySession\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Session destroyed\"},\"404\":{\"description\":\"Session not found\"}}}},\"/api/user\":{\"get\":{\"summary\":\"List users\",\"operationId\":\"GetUsers\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\",\"schema\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/User\"}}},\"403\":{\"description\":\"Not authenticated\"}}}},\"/api/user/{userId}\":{\"get\":{\"summary\":\"Get user\",\"operationId\":\"GetUser\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\",\"schema\":{\"$ref\":\"#/definitions/User\"}},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"put\":{\"summary\":\"Create or update user\",\"operationId\":\"CreateOrUpdateUser\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true},{\"name\":\"body\",\"in\":\"body\",\"required\":true,\"schema\":{\"$ref\":\"#/definitions/User\"}}],\"responses\":{\"200\":{\"description\":\"Updated\"},\"201\":{\"description\":\"Created\"},\"400\":{\"description\":\"Malformed request body\"},\"403\":{\"description\":\"Not authenticated\"},\"405\":{\"description\":\"Not allowed\"}}},\"delete\":{\"summary\":\"Delete user\",\"operationId\":\"DeleteUser\",\"parameters\":[{\"type\":\"string\",\"description\":\"Authentication token\",\"name\":\"X-Auth\",\"in\":\"header\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Success\"},\"403\":{\"description\":\"Not authenticated\"},\"404\":{\"description\":\"Not found\"}}},\"parameters\":[{\"type\":\"string\",\"name\":\"userId\",\"in\":\"path\",\"required\":true},{\"type\":\"boolean\",\"name\":\"allKeys\",\"in\":\"query\"}]},\"/booking\":{\"get\":{\"security\":[{\"X-Session-ID\":[]}],\"description\":\"Get booking of session owner\",\"consumes\":[\"application/xml\"],\"summary\":\"Get booking\",\"operationId\":\"GetBooking\",\"responses\":{\"200\":{\"description\":\"status 200\",\"schema\":{\"type\":\"string\"}},\"400\":{\"description\":\"status 400\"},\"401\":{\"description\":\"Unauthorized Session Token\"},\"404\":{\"description\":\"Resource Not Found\"},\"500\":{\"description\":\"Malfunction (internal requirements not fulfilled)\"}}}},\"/bookings\":{\"get\":{\"security\":[{\"X-Session-ID\":[]}],\"description\":\"Get bookings of session owner\",\"produces\":[\"application/json\"],\"summary\":\"Get bookings\",\"operationId\":\"GetBookings\",\"parameters\":[{\"type\":\"string\",\"name\":\"date\",\"in\":\"header\"},{\"type\":\"array\",\"items\":{\"type\":\"integer\"},\"name\":\"ids\",\"in\":\"query\"}],\"responses\":{\"200\":{\"description\":\"Success List Booking History\",\"schema\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Booking\"}}},\"400\":{\"description\":\"status 400\"},\"401\":{\"description\":\"Unauthorized Session Token\"},\"404\":{\"description\":\"Resource Not Found\"},\"500\":{\"description\":\"Malfunction (internal requirements not fulfilled)\"}}}},\"/brands/{brandId}/models\":{\"get\":{\"tags\":[\"MODEL\"],\"summary\":\"Get all available models for the given brandId\",\"operationId\":\"ListModels\",\"parameters\":[{\"name\":\"driveConcept\",\"in\":\"query\",\"schema\":{\"$ref\":\"#/definitions/DriveConcept\"}},{\"type\":\"string\",\"x-example\":\"de\",\"name\":\"languageId\",\"in\":\"query\"},{\"type\":\"string\",\"x-example\":\"123\",\"name\":\"classId\",\"in\":\"query\"},{\"type\":\"string\",\"name\":\"lineId\",\"in\":\"query\"},{\"type\":\"array\",\"items\":{\"type\":\"integer\"},\"name\":\"ids\",\"in\":\"query\"}],\"responses\":{\"200\":{\"description\":\"Ok\",\"schema\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Model\"}},\"examples\":{\"application/json\":{\"drive_concept\":\"drive_concept\",\"price\":38,\"technical_information\":null}}}}},\"parameters\":[{\"type\":\"string\",\"name\":\"brandId\",\"in\":\"path\",\"required\":true}]},\"/classes/{productGroup}\":{\"get\":{\"summary\":\"Get all available classes.\",\"operationId\":\"GetClasses\",\"parameters\":[{\"enum\":[\"WHEELS\",\"PAINTS\",\"UPHOLSTERIES\",\"TRIMS\",\"PACKAGES\",\"LINES\",\"SPECIAL_EDITION\",\"SPECIAL_EQUIPMENT\"],\"type\":\"array\",\"items\":{\"type\":\"string\"},\"description\":\"A list of component types separated by a comma case insensitive. If nothing is defined all component types are returned.\",\"name\":\"componentTypes\",\"in\":\"query\"},{\"enum\":[\"PKW\",\"GELAENDEWAGEN\",\"VAN\",\"SPRINTER\",\"CITAN\",\"SMART\"],\"type\":\"string\",\"default\":\"PKW\",\"description\":\"The productGroup of a vehicle case insensitive.\",\"name\":\"productGroup\",\"in\":\"path\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"Successful response\",\"schema\":{\"type\":\"string\"}},\"400\":{\"description\":\"Successful response\",\"schema\":{\"type\":\"string\"}}}}},\"/code\":{\"post\":{\"consumes\":[\"application/x-www-form-urlencoded\"],\"summary\":\"code to token\",\"operationId\":\"Code\",\"parameters\":[{\"type\":\"array\",\"items\":{\"type\":\"integer\"},\"name\":\"state\",\"in\":\"formData\"},{\"type\":\"string\",\"name\":\"response_mode\",\"in\":\"formData\"},{\"type\":\"string\",\"name\":\"code\",\"in\":\"formData\",\"required\":true},{\"type\":\"string\",\"name\":\"session\",\"in\":\"query\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"TBD\",\"schema\":{\"type\":\"string\"}},\"400\":{\"description\":\"status 400\"},\"401\":{\"description\":\"Unauthorized Session code\"},\"404\":{\"description\":\"Resource Not Found\"},\"500\":{\"description\":\"Malfunction (internal requirements not fulfilled)\"}}}},\"/customer/session\":{\"post\":{\"description\":\"Creates a customer session for a given OpenID authentication token.\\n\",\"consumes\":[\"application/x-www-form-urlencoded\"],\"produces\":[\"application/json\"],\"summary\":\"Create session (login)\",\"operationId\":\"CreateCustomerSession\",\"parameters\":[{\"maxLength\":255,\"type\":\"string\",\"description\":\"OpenID authentication token\",\"name\":\"code\",\"in\":\"formData\",\"required\":true},{\"maxLength\":255,\"pattern\":\"^([a-z]{2})-([A-Z]{2})$\",\"type\":\"string\",\"description\":\"default locale\",\"name\":\"locale\",\"in\":\"formData\"},{\"type\":\"string\",\"description\":\"ID of the request in UUIDv4 format\",\"name\":\"X-Request-ID\",\"in\":\"header\"}],\"responses\":{\"201\":{\"description\":\"Session successful created\",\"schema\":{\"$ref\":\"#/definitions/Session\"}},\"401\":{\"description\":\"Invalid OpenID authentication token\"},\"403\":{\"description\":\"Create session with authentication token is forbidden (e.g. Token already used)\\n\"},\"422\":{\"description\":\"Invalid request data\",\"schema\":{\"$ref\":\"#/definitions/ValidationErrors\"}},\"500\":{\"description\":\"Internal server error (e.g. unexpected condition occurred)\"}}},\"delete\":{\"security\":[{\"X-Session-ID\":[]}],\"description\":\"Deletes the user session matching the *X-Auth* header.\\n\",\"summary\":\"Delete session (logout)\",\"operationId\":\"DeleteCustomerSession\",\"parameters\":[{\"type\":\"string\",\"description\":\"ID of the request in UUIDv4 format\",\"name\":\"X-Request-ID\",\"in\":\"header\"}],\"responses\":{\"204\":{\"description\":\"Session successful deleted\"},\"401\":{\"description\":\"Invalid session token\"},\"500\":{\"description\":\"Internal server error (e.g. unexpected condition occurred)\"}}}},\"/download/nested/file\":{\"get\":{\"description\":\"Downloads a file that is a property within a nested structure in the response body\\n\",\"produces\":[\"application/json\"],\"summary\":\"Downloads a nested file\",\"operationId\":\"DownloadNestedFile\",\"responses\":{\"200\":{\"description\":\"Nested file structure\",\"schema\":{\"$ref\":\"#/definitions/NestedFileStructure\"}}}}},\"/download/{image}\":{\"get\":{\"description\":\"Retrieve a image\",\"produces\":[\"image/png\"],\"summary\":\"Retrieve a image\",\"operationId\":\"DownloadImage\",\"parameters\":[{\"type\":\"string\",\"description\":\"The image name of the image\",\"name\":\"image\",\"in\":\"path\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"image to download\",\"schema\":{\"type\":\"file\"},\"headers\":{\"Content-Type\":{\"type\":\"string\"}}},\"500\":{\"description\":\"Malfunction (internal requirements not fulfilled)\"}}}},\"/elements\":{\"get\":{\"summary\":\"ListElements\",\"operationId\":\"ListElements\",\"parameters\":[{\"type\":\"integer\",\"default\":1,\"name\":\"_page\",\"in\":\"query\"},{\"type\":\"integer\",\"default\":10,\"name\":\"_perPage\",\"in\":\"query\"}],\"responses\":{\"200\":{\"description\":\"Status 200\",\"schema\":{\"type\":\"string\"},\"headers\":{\"X-Total-Count\":{\"type\":\"integer\"}}},\"500\":{\"description\":\"Status 500\"}}}},\"/file-upload\":{\"post\":{\"consumes\":[\"multipart/form-data\"],\"summary\":\"File upload\",\"operationId\":\"FileUpload\",\"parameters\":[{\"type\":\"file\",\"description\":\"File to be uploaded in request.\",\"name\":\"file\",\"in\":\"formData\"}],\"responses\":{\"204\":{\"description\":\"File uploaded.\"},\"500\":{\"description\":\"Internal server error\"}}}},\"/filedownload/{file}\":{\"get\":{\"description\":\"Retrieve a file\",\"produces\":[\"text/xml\"],\"summary\":\"Retrieve a file\",\"operationId\":\"DownloadFile\",\"responses\":{\"200\":{\"description\":\"file to download\",\"schema\":{\"type\":\"file\"},\"headers\":{\"Content-Type\":{\"type\":\"string\"}}}}},\"parameters\":[{\"type\":\"string\",\"description\":\"The filename of the file\",\"name\":\"file\",\"in\":\"path\",\"required\":true}]},\"/findByTags\":{\"get\":{\"description\":\"Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.\",\"produces\":[\"application/json\"],\"summary\":\"Finds elements by tags\",\"operationId\":\"FindByTags\",\"parameters\":[{\"maxItems\":5,\"minItems\":2,\"type\":\"array\",\"items\":{\"type\":\"string\"},\"description\":\"Tags to filter by\",\"name\":\"tags\",\"in\":\"query\",\"required\":true}],\"responses\":{\"200\":{\"description\":\"successful operation\",\"schema\":{\"type\":\"string\"}},\"400\":{\"description\":\"Invalid tag value\"}}}},\"/generic/download/{ext}\":{\"get\":{\"description\":\"Retrieve a file\",\"produces\":[\"application/json\"],\"summary\":\"Retrieve a file\",\"operationId\":\"GenericFileDownload\",\"responses\":{\"200\":{\"description\":\"file to download\",\"schema\":{\"type\":\"file\"},\"headers\":{\"Content-Type\":{\"type\":\"string\"},\"Pragma\":{\"type\":\"string\"}}},\"500\":{\"description\":\"Malfunction (internal requirements not fulfilled)\"}}},\"parameters\":[{\"type\":\"string\",\"description\":\"The ext of the file\",\"name\":\"ext\",\"in\":\"path\",\"required\":true}]},\"/rental\":{\"get\":{\"description\":\"get rental\",\"consumes\":[\"application/json\"],\"summary\":\"Get rental\",\"operationId\":\"GetRental\",\"parameters\":[{\"name\":\"body\",\"in\":\"body\",\"required\":true,\"schema\":{\"$ref\":\"#/definitions/Rental\"}}],\"responses\":{\"200\":{\"description\":\"status 200\"},\"400\":{\"description\":\"status 400\",\"schema\":{\"$ref\":\"#/definitions/ValidationErrors\"}}}}},\"/shop/shoes\":{\"get\":{\"produces\":[\"application/hal+json\"],\"summary\":\"Get all shoes\",\"operationId\":\"GetShoes\",\"responses\":{\"200\":{\"description\":\"Successful\",\"schema\":{\"$ref\":\"#/definitions/Shoes\"}}}}},\"/upload\":{\"post\":{\"consumes\":[\"multipart/form-data\"],\"summary\":\"Upload a file with others data\",\"operationId\":\"PostUpload\",\"parameters\":[{\"type\":\"file\",\"description\":\"the file to upload\",\"name\":\"upfile\",\"in\":\"formData\"},{\"maxLength\":4000,\"pattern\":\"^[0-9a-zA-Z ]*$\",\"type\":\"string\",\"description\":\"Description of file\",\"name\":\"note\",\"in\":\"formData\"}],\"responses\":{\"200\":{\"description\":\"Status 200\"},\"500\":{\"description\":\"Status 500\"}}}}},\"definitions\":{\"Address\":{\"type\":\"object\",\"required\":[\"city\",\"country\",\"houseNumber\",\"postalCode\",\"region\",\"street\"],\"properties\":{\"city\":{\"description\":\"City\",\"type\":\"string\"},\"country\":{\"description\":\"Country (ISO 3166)\",\"type\":\"string\"},\"houseNumber\":{\"description\":\"House number\",\"type\":\"string\"},\"postalCode\":{\"description\":\"Postal code\",\"type\":\"string\"},\"region\":{\"description\":\"Region\",\"type\":\"string\"},\"street\":{\"description\":\"Street name\",\"type\":\"string\"}}},\"BasicTypes\":{\"type\":\"object\",\"required\":[\"string\",\"integer\",\"boolean\",\"number\",\"slice\",\"map\"],\"properties\":{\"boolean\":{\"type\":\"boolean\"},\"integer\":{\"type\":\"integer\"},\"map\":{\"type\":\"object\",\"additionalProperties\":{\"type\":\"string\"}},\"number\":{\"type\":\"number\"},\"slice\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}},\"string\":{\"type\":\"string\"}}},\"Booking\":{\"type\":\"object\",\"required\":[\"id\"],\"properties\":{\"bookingID\":{\"type\":\"string\"}}},\"Client\":{\"type\":\"object\",\"required\":[\"id\",\"name\"],\"properties\":{\"activePresets\":{\"type\":\"string\"},\"configuration\":{\"type\":\"object\",\"properties\":{\"bbdCEBaseUrl\":{\"type\":\"string\"},\"bbdCallerIdentifier\":{\"type\":\"string\"},\"bbdDataSupply\":{\"type\":\"string\"},\"bbdImageBackground\":{\"type\":\"string\"},\"bbdImagePerspective\":{\"type\":\"string\"},\"bbdImageType\":{\"type\":\"string\"},\"bbdPassword\":{\"type\":\"string\"},\"bbdProductGroup\":{\"type\":\"string\"},\"bbdSoapMediaProviderUrl\":{\"type\":\"string\"},\"bbdUser\":{\"type\":\"string\"},\"ccoreServiceUrl\":{\"type\":\"string\"},\"cryptKeys\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}},\"healConfigurations\":{\"type\":\"boolean\"}}},\"id\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"}}},\"DriveConcept\":{\"description\":\"The kind of drive concept of a vehicle. Where UNDEFINED is used as the default and/or error case.\",\"type\":\"string\",\"enum\":[\"COMBUSTOR\",\"HYBRID\",\"ELECTRIC\",\"FUELCELL\",\"UNDEFINED\"]},\"EmptySlice\":{\"properties\":{\"EmptySlice\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Price\"}}}},\"Link\":{\"type\":\"object\",\"required\":[\"href\"],\"properties\":{\"href\":{\"type\":\"string\"}}},\"Links\":{\"type\":\"object\",\"required\":[\"self\"],\"properties\":{\"self\":{\"$ref\":\"#/definitions/Link\"}}},\"Model\":{\"type\":\"object\",\"required\":[\"technicalInformation\",\"price\"],\"properties\":{\"driveConcept\":{\"$ref\":\"#/definitions/DriveConcept\"},\"price\":{\"$ref\":\"#/definitions/Price\"},\"technicalInformation\":{\"$ref\":\"#/definitions/TechnicalInformation\"}}},\"NestedFileStructure\":{\"properties\":{\"data\":{\"type\":\"string\"}}},\"Price\":{\"type\":\"object\",\"required\":[\"currency\",\"value\"],\"properties\":{\"currency\":{\"type\":\"string\",\"example\":\"RMB\"},\"value\":{\"type\":\"number\",\"example\":123456.78}}},\"Rental\":{\"type\":\"object\",\"required\":[\"class\",\"lockStatus\",\"status\",\"stationID\",\"maxDoors\",\"minDoors\",\"website\",\"id\"],\"properties\":{\"class\":{\"type\":\"string\",\"maxLength\":20,\"minLength\":3},\"color\":{\"type\":\"string\",\"maxLength\":20,\"minLength\":3},\"homeID\":{\"type\":\"string\",\"pattern\":\"^[a-zA-Z]$\"},\"id\":{\"type\":\"string\",\"format\":\"uuid\"},\"idOptional\":{\"type\":\"string\",\"format\":\"uuid\"},\"lockStatus\":{\"type\":\"integer\",\"format\":\"int32\",\"maximum\":100,\"minimum\":1,\"exclusiveMinimum\":true},\"maxDoors\":{\"type\":\"integer\",\"maximum\":5},\"minDoors\":{\"type\":\"integer\",\"format\":\"int64\",\"minimum\":5},\"optionalInt\":{\"type\":\"integer\"},\"state\":{\"type\":\"integer\",\"format\":\"int64\"},\"stationID\":{\"type\":\"string\",\"pattern\":\"^[a-zA-Z]$\"},\"status\":{\"type\":\"integer\",\"maximum\":49,\"exclusiveMaximum\":true,\"minimum\":46,\"exclusiveMinimum\":true},\"valid\":{\"type\":\"string\",\"maxLength\":255},\"website\":{\"type\":\"string\",\"format\":\"url\"},\"websiteOptional\":{\"type\":\"string\",\"format\":\"url\",\"maxLength\":255}}},\"Session\":{\"type\":\"object\",\"required\":[\"Token\",\"Registered\"],\"properties\":{\"Registered\":{\"description\":\"Indicates if the user is registered at the rental system\",\"type\":\"boolean\"},\"Token\":{\"description\":\"Token used within the X-Session-ID header\",\"type\":\"string\"}}},\"Shoe\":{\"type\":\"object\",\"required\":[\"name\",\"size\",\"color\",\"_links\"],\"properties\":{\"_links\":{\"$ref\":\"#/definitions/Links\"},\"color\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"size\":{\"type\":\"number\"}}},\"Shoes\":{\"type\":\"object\",\"required\":[\"id\",\"_embedded\",\"_links\"],\"properties\":{\"_embedded\":{\"$ref\":\"#/definitions/ShoesEmbedded\"},\"_links\":{\"$ref\":\"#/definitions/Links\"},\"id\":{\"type\":\"string\"}}},\"ShoesEmbedded\":{\"type\":\"object\",\"required\":[\"shop:shoes\"],\"properties\":{\"shop:shoes\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Shoe\"}}}},\"TechnicalInformation\":{\"type\":\"object\",\"required\":[\"transmission\"],\"properties\":{\"transmission\":{\"type\":\"string\",\"example\":\"7G-DCT\"}}},\"User\":{\"type\":\"object\",\"required\":[\"id\",\"password\"],\"properties\":{\"Address\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Address\"}},\"email\":{\"type\":\"string\",\"format\":\"email\",\"maxLength\":255},\"grantedProtocolMappers\":{\"type\":\"object\",\"additionalProperties\":{\"type\":\"string\"}},\"id\":{\"type\":\"string\"},\"password\":{\"type\":\"string\"},\"permissions\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}}},\"ValidationError\":{\"type\":\"object\",\"properties\":{\"Code\":{\"type\":\"string\"},\"Field\":{\"type\":\"string\"},\"Message\":{\"type\":\"string\"}}},\"ValidationErrors\":{\"type\":\"object\",\"properties\":{\"Errors\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/ValidationError\"}},\"Message\":{\"type\":\"string\"}}},\"views set\":{\"type\":\"object\",\"required\":[\"id\"],\"properties\":{\"id\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"},\"views\":{\"description\":\"View definitions in YAML format\",\"type\":\"string\"}}}},\"parameters\":{\"X-Request-ID\":{\"type\":\"string\",\"description\":\"ID of the request in UUIDv4 format\",\"name\":\"X-Request-ID\",\"in\":\"header\"},\"componentType\":{\"enum\":[\"WHEELS\",\"PAINTS\",\"UPHOLSTERIES\",\"TRIMS\",\"PACKAGES\",\"LINES\",\"SPECIAL_EDITION\",\"SPECIAL_EQUIPMENT\"],\"type\":\"array\",\"items\":{\"type\":\"string\"},\"description\":\"A list of component types separated by a comma case insensitive. If nothing is defined all component types are returned.\",\"name\":\"componentTypes\",\"in\":\"query\"},\"fileParam\":{\"type\":\"file\",\"description\":\"File to be uploaded in request.\",\"name\":\"file\",\"in\":\"formData\"},\"productGroup\":{\"enum\":[\"PKW\",\"GELAENDEWAGEN\",\"VAN\",\"SPRINTER\",\"CITAN\",\"SMART\"],\"type\":\"string\",\"default\":\"PKW\",\"description\":\"The productGroup of a vehicle case insensitive.\",\"name\":\"productGroup\",\"in\":\"path\",\"required\":true}},\"securityDefinitions\":{\"X-Session-ID\":{\"type\":\"apiKey\",\"name\":\"X-Session-ID\",\"in\":\"header\"}}}"
