package basket

import (
	"context"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing"
	"net/http"
)

type PostItemHandler func(ctx context.Context, request *PostItemRequest) PostItemResponse

type postItemHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    PostItemHandler
}

func (server *BasketServiceServer) SetPostItemHandler(handler PostItemHandler, middleware ...Middleware) {
	server.postItemHandler = &postItemHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/item", Handler: server.PostItemHandler, Middleware: middleware}}
}

func (server *BasketServiceServer) PostItemHandler(c *routing.Context) error {
	if server.postItemHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: PostItem (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(PostItemRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.Item, false)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PostItem (POST) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			if contentTypeOfResponse != "" {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PostItem (POST) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
				return newNotSupportedContentType(415, contentTypeOfResponse)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PostItem (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.postItemHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: PostItem (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PostItem (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetItemsHandler func(ctx context.Context, request *GetItemsRequest) GetItemsResponse

type getItemsHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetItemsHandler
}

func (server *BasketServiceServer) SetGetItemsHandler(handler GetItemsHandler, middleware ...Middleware) {
	server.getItemsHandler = &getItemsHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/items", Handler: server.GetItemsHandler, Middleware: middleware}}
}

func (server *BasketServiceServer) GetItemsHandler(c *routing.Context) error {
	if server.getItemsHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetItems (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetItemsRequest)
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetItems (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getItemsHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetItems (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetItems (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type BasketServiceServer struct {
	*Server
	Validator       *Validator
	postItemHandler *postItemHandlerRoute
	getItemsHandler *getItemsHandlerRoute
}

func (server *BasketServiceServer) registerValidators() {}

func NewBasketServiceServer(options *ServerOpts) *BasketServiceServer {
	serverWrapper := &BasketServiceServer{Server: newServer(options), Validator: NewValidation()}
	serverWrapper.Server.SwaggerSpec = swagger
	serverWrapper.registerValidators()
	return serverWrapper
}

func (server *BasketServiceServer) Start(port int) error {
	routes := []RouteDescription{}
	if server.postItemHandler != nil {
		routes = append(routes, server.postItemHandler.routeDescription)
	}
	if server.getItemsHandler != nil {
		routes = append(routes, server.getItemsHandler.routeDescription)
	}
	return server.Server.Start(port, routes)
}

const swagger = "{\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"],\"schemes\":[\"http\"],\"swagger\":\"2.0\",\"info\":{\"title\":\"Basket Service\",\"version\":\"1.0.0\"},\"paths\":{\"/item\":{\"post\":{\"operationId\":\"PostItem\",\"parameters\":[{\"name\":\"item\",\"in\":\"body\",\"schema\":{\"$ref\":\"#/definitions/Item\"}}],\"responses\":{\"200\":{\"description\":\"OK\"},\"500\":{\"description\":\"Internal server error (e.g. unexpected condition occurred)\"}}}},\"/items\":{\"get\":{\"operationId\":\"GetItems\",\"responses\":{\"200\":{\"description\":\"Items posted\",\"schema\":{\"$ref\":\"#/definitions/Items\"}}}}}},\"definitions\":{\"Item\":{\"type\":\"object\",\"required\":[\"id\",\"name\",\"price\"],\"properties\":{\"id\":{\"type\":\"integer\"},\"name\":{\"type\":\"string\"},\"price\":{\"type\":\"number\"}}},\"Items\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Item\"}}},\"parameters\":{\"ItemParameter\":{\"name\":\"item\",\"in\":\"body\",\"schema\":{\"$ref\":\"#/definitions/Item\"}}}}"
