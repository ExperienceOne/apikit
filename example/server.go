package todo

import (
	"context"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing"
	"net/http"
)

type DeleteTodosHandler func(ctx context.Context, request *DeleteTodosRequest) DeleteTodosResponse

type deleteTodosHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteTodosHandler
}

func (server *TodoServiceServer) SetDeleteTodosHandler(handler DeleteTodosHandler, middleware ...Middleware) {
	server.deleteTodosHandler = &deleteTodosHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/todos", Handler: server.DeleteTodosHandler, Middleware: middleware}}
}

func (server *TodoServiceServer) DeleteTodosHandler(c *routing.Context) error {
	if server.deleteTodosHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteTodos (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteTodosRequest)
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodos (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteTodosHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteTodos (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodos (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type ListTodosHandler func(ctx context.Context, request *ListTodosRequest) ListTodosResponse

type listTodosHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    ListTodosHandler
}

func (server *TodoServiceServer) SetListTodosHandler(handler ListTodosHandler, middleware ...Middleware) {
	server.listTodosHandler = &listTodosHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/todos", Handler: server.ListTodosHandler, Middleware: middleware}}
}

func (server *TodoServiceServer) ListTodosHandler(c *routing.Context) error {
	if server.listTodosHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: ListTodos (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(ListTodosRequest)
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListTodos (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.listTodosHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: ListTodos (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: ListTodos (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type PostTodoHandler func(ctx context.Context, request *PostTodoRequest) PostTodoResponse

type postTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    PostTodoHandler
}

func (server *TodoServiceServer) SetPostTodoHandler(handler PostTodoHandler, middleware ...Middleware) {
	server.postTodoHandler = &postTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "POST", Path: "/todos", Handler: server.PostTodoHandler, Middleware: middleware}}
}

func (server *TodoServiceServer) PostTodoHandler(c *routing.Context) error {
	if server.postTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: PostTodo (POST) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(PostTodoRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.TodoPost, false)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PostTodo (POST) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			if contentTypeOfResponse != "" {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PostTodo (POST) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
				return newNotSupportedContentType(415, contentTypeOfResponse)
			}
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PostTodo (POST) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.postTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: PostTodo (POST) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PostTodo (POST) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type DeleteTodoHandler func(ctx context.Context, request *DeleteTodoRequest) DeleteTodoResponse

type deleteTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    DeleteTodoHandler
}

func (server *TodoServiceServer) SetDeleteTodoHandler(handler DeleteTodoHandler, middleware ...Middleware) {
	server.deleteTodoHandler = &deleteTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "DELETE", Path: "/todos/<todoId>", Handler: server.DeleteTodoHandler, Middleware: middleware}}
}

func (server *TodoServiceServer) DeleteTodoHandler(c *routing.Context) error {
	if server.deleteTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: DeleteTodo (DELETE) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(DeleteTodoRequest)
		if err := fromString(c.Param("todoId"), &request.TodoId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodo (DELETE) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodo (DELETE) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.deleteTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: DeleteTodo (DELETE) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: DeleteTodo (DELETE) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type GetTodoHandler func(ctx context.Context, request *GetTodoRequest) GetTodoResponse

type getTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    GetTodoHandler
}

func (server *TodoServiceServer) SetGetTodoHandler(handler GetTodoHandler, middleware ...Middleware) {
	server.getTodoHandler = &getTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "GET", Path: "/todos/<todoId>", Handler: server.GetTodoHandler, Middleware: middleware}}
}

func (server *TodoServiceServer) GetTodoHandler(c *routing.Context) error {
	if server.getTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: GetTodo (GET) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(GetTodoRequest)
		if err := fromString(c.Param("todoId"), &request.TodoId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodo (GET) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodo (GET) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.getTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: GetTodo (GET) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: GetTodo (GET) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type PatchTodoHandler func(ctx context.Context, request *PatchTodoRequest) PatchTodoResponse

type patchTodoHandlerRoute struct {
	routeDescription RouteDescription
	customHandler    PatchTodoHandler
}

func (server *TodoServiceServer) SetPatchTodoHandler(handler PatchTodoHandler, middleware ...Middleware) {
	server.patchTodoHandler = &patchTodoHandlerRoute{customHandler: handler, routeDescription: RouteDescription{Method: "PATCH", Path: "/todos/<todoId>", Handler: server.PatchTodoHandler, Middleware: middleware}}
}

func (server *TodoServiceServer) PatchTodoHandler(c *routing.Context) error {
	if server.patchTodoHandler.customHandler == nil {
		server.ErrorLogger("wrap handler: PatchTodo (PATCH) endpoint is not registered")
		return NewHTTPStatusCodeError(http.StatusNotFound)
	} else {
		request := new(PatchTodoRequest)
		contentTypeOfResponse := extractContentType(c.Request.Header.Get(contentTypeHeader))
		if contentTypeOfResponse == contentTypeApplicationJson {
			err := JSON(c.Request.Body, &request.TodoPatch, false)
			if err != nil {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PatchTodo (PATCH) could not decode request body of incoming request (%v)", err))
				return NewHTTPStatusCodeError(http.StatusBadRequest)
			}
		} else {
			if contentTypeOfResponse != "" {
				server.ErrorLogger(fmt.Sprintf("wrap handler: PatchTodo (PATCH) content type of incoming request is bad (want: application/json, got: %s)", contentTypeOfResponse))
				return newNotSupportedContentType(415, contentTypeOfResponse)
			}
		}
		if err := fromString(c.Param("todoId"), &request.TodoId); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PatchTodo (PATCH) could not convert string to specific type (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		validationErrors, err := server.Validator.ValidateRequest(request)
		if err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PatchTodo (PATCH) could not validate incoming request (error: %v)", err))
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if validationErrors != nil {
			return NewHTTPStatusCodeError(http.StatusBadRequest)
		}
		response := server.patchTodoHandler.customHandler(c.Request.Context(), request)
		if response == nil {
			server.ErrorLogger("wrap handler: PatchTodo (PATCH) received a nil response object")
			return NewHTTPStatusCodeError(http.StatusInternalServerError)
		}
		if err := response.write(c.Response); err != nil {
			server.ErrorLogger(fmt.Sprintf("wrap handler: PatchTodo (PATCH) could not send response (error: %v)", err))
			return err
		}
	}
	return nil
}

type TodoServiceServer struct {
	*Server
	Validator          *Validator
	deleteTodosHandler *deleteTodosHandlerRoute
	listTodosHandler   *listTodosHandlerRoute
	postTodoHandler    *postTodoHandlerRoute
	deleteTodoHandler  *deleteTodoHandlerRoute
	getTodoHandler     *getTodoHandlerRoute
	patchTodoHandler   *patchTodoHandlerRoute
}

func (server *TodoServiceServer) registerValidators() {}

func NewTodoServiceServer(options *ServerOpts) *TodoServiceServer {
	serverWrapper := &TodoServiceServer{Server: newServer(options), Validator: NewValidation()}
	serverWrapper.Server.SwaggerSpec = swagger
	serverWrapper.registerValidators()
	return serverWrapper
}

func (server *TodoServiceServer) Start(port int) error {
	routes := []RouteDescription{}
	if server.deleteTodosHandler != nil {
		routes = append(routes, server.deleteTodosHandler.routeDescription)
	}
	if server.listTodosHandler != nil {
		routes = append(routes, server.listTodosHandler.routeDescription)
	}
	if server.postTodoHandler != nil {
		routes = append(routes, server.postTodoHandler.routeDescription)
	}
	if server.deleteTodoHandler != nil {
		routes = append(routes, server.deleteTodoHandler.routeDescription)
	}
	if server.getTodoHandler != nil {
		routes = append(routes, server.getTodoHandler.routeDescription)
	}
	if server.patchTodoHandler != nil {
		routes = append(routes, server.patchTodoHandler.routeDescription)
	}
	return server.Server.Start(port, routes)
}

const swagger = "{\"consumes\":[\"application/json\"],\"produces\":[\"application/json\"],\"schemes\":[\"http\"],\"swagger\":\"2.0\",\"info\":{\"title\":\"Todo Service\",\"version\":\"1.0.0\"},\"host\":\"localhost:9001\",\"paths\":{\"/todos\":{\"get\":{\"operationId\":\"ListTodos\",\"responses\":{\"200\":{\"description\":\"List of todos\",\"schema\":{\"$ref\":\"#/definitions/TodoList\"}}}},\"post\":{\"operationId\":\"PostTodo\",\"parameters\":[{\"name\":\"todoPost\",\"in\":\"body\",\"schema\":{\"type\":\"object\",\"required\":[\"title\"],\"properties\":{\"title\":{\"type\":\"string\"}}}}],\"responses\":{\"201\":{\"description\":\"Created\",\"schema\":{\"$ref\":\"#/definitions/Todo\"}}}},\"delete\":{\"operationId\":\"DeleteTodos\",\"responses\":{\"204\":{\"description\":\"Ok\"}}}},\"/todos/{todoId}\":{\"get\":{\"operationId\":\"GetTodo\",\"responses\":{\"200\":{\"description\":\"Successful\",\"schema\":{\"$ref\":\"#/definitions/Todo\"}},\"404\":{\"description\":\"Not found\"}}},\"delete\":{\"operationId\":\"DeleteTodo\",\"responses\":{\"204\":{\"description\":\"Ok\"},\"404\":{\"description\":\"Not found\"}}},\"patch\":{\"operationId\":\"PatchTodo\",\"parameters\":[{\"name\":\"TodoPatch\",\"in\":\"body\",\"schema\":{\"type\":\"object\",\"properties\":{\"completed\":{\"type\":\"boolean\"},\"order\":{\"type\":\"integer\"},\"title\":{\"type\":\"string\"}}}}],\"responses\":{\"200\":{\"description\":\"Successful\",\"schema\":{\"$ref\":\"#/definitions/Todo\"}},\"404\":{\"description\":\"Not found\"}}},\"parameters\":[{\"type\":\"integer\",\"name\":\"todoId\",\"in\":\"path\",\"required\":true}]}},\"definitions\":{\"Todo\":{\"type\":\"object\",\"required\":[\"id\",\"title\",\"order\",\"completed\",\"url\"],\"properties\":{\"completed\":{\"type\":\"boolean\"},\"id\":{\"type\":\"integer\"},\"order\":{\"type\":\"integer\"},\"title\":{\"type\":\"string\"},\"url\":{\"type\":\"string\"}}},\"TodoList\":{\"type\":\"array\",\"items\":{\"$ref\":\"#/definitions/Todo\"}}},\"parameters\":{\"TodoId\":{\"type\":\"integer\",\"name\":\"todoId\",\"in\":\"path\",\"required\":true},\"TodoPatch\":{\"name\":\"TodoPatch\",\"in\":\"body\",\"schema\":{\"type\":\"object\",\"properties\":{\"completed\":{\"type\":\"boolean\"},\"order\":{\"type\":\"integer\"},\"title\":{\"type\":\"string\"}}}},\"TodoPost\":{\"name\":\"todoPost\",\"in\":\"body\",\"schema\":{\"type\":\"object\",\"required\":[\"title\"],\"properties\":{\"title\":{\"type\":\"string\"}}}}}}"
