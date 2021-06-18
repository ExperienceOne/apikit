package todo

import "net/http"

var contentTypesForFiles = []string{"application/json", "image/png", "image/jpeg", "image/tiff", "image/webp", "image/svg+xml", "image/gif", "image/tiff", "image/x-icon", "application/pdf", "application/octet-stream"}

type Todo struct {
	Completed bool   `bson:"completed,required" json:"completed,required" xml:"completed,required"`
	Id        int64  `bson:"id,required" json:"id,required" xml:"id,required"`
	Order     int64  `bson:"order,required" json:"order,required" xml:"order,required"`
	Title     string `bson:"title,required" json:"title,required" xml:"title,required"`
	Url       string `bson:"url,required" json:"url,required" xml:"url,required"`
}

type TodoList []Todo
type Object1 struct {
	Completed *bool   `bson:"completed,omitempty" json:"completed,omitempty" xml:"completed,omitempty"`
	Order     *int64  `bson:"order,omitempty" json:"order,omitempty" xml:"order,omitempty"`
	Title     *string `bson:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
}

type Object2 struct {
	Title string `bson:"title,required" json:"title,required" xml:"title,required"`
}

type DeleteTodosRequest struct{}

type DeleteTodosResponse interface {
	isDeleteTodosResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Ok
type DeleteTodos204Response struct{}

func (r *DeleteTodos204Response) isDeleteTodosResponse() {}

func (r *DeleteTodos204Response) StatusCode() int {
	return 204
}

func (r *DeleteTodos204Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(204)
	return nil
}

type ListTodosRequest struct{}

type ListTodosResponse interface {
	isListTodosResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// List of todos
type ListTodos200Response struct {
	Body TodoList
}

func (r *ListTodos200Response) isListTodosResponse() {}

func (r *ListTodos200Response) StatusCode() int {
	return 200
}

func (r *ListTodos200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type Object3 struct {
	Title string `bson:"title,required" json:"title,required" xml:"title,required"`
}

type PostTodoRequest struct {
	TodoPost Object3
}

type PostTodoResponse interface {
	isPostTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Created
type PostTodo201Response struct {
	Body Todo
}

func (r *PostTodo201Response) isPostTodoResponse() {}

func (r *PostTodo201Response) StatusCode() int {
	return 201
}

func (r *PostTodo201Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 201, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type DeleteTodoRequest struct {
	TodoId int64
}

type DeleteTodoResponse interface {
	isDeleteTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Ok
type DeleteTodo204Response struct{}

func (r *DeleteTodo204Response) isDeleteTodoResponse() {}

func (r *DeleteTodo204Response) StatusCode() int {
	return 204
}

func (r *DeleteTodo204Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(204)
	return nil
}

// Not found
type DeleteTodo404Response struct{}

func (r *DeleteTodo404Response) isDeleteTodoResponse() {}

func (r *DeleteTodo404Response) StatusCode() int {
	return 404
}

func (r *DeleteTodo404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type GetTodoRequest struct {
	TodoId int64
}

type GetTodoResponse interface {
	isGetTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Successful
type GetTodo200Response struct {
	Body Todo
}

func (r *GetTodo200Response) isGetTodoResponse() {}

func (r *GetTodo200Response) StatusCode() int {
	return 200
}

func (r *GetTodo200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not found
type GetTodo404Response struct{}

func (r *GetTodo404Response) isGetTodoResponse() {}

func (r *GetTodo404Response) StatusCode() int {
	return 404
}

func (r *GetTodo404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type Object4 struct {
	Completed *bool   `bson:"completed,omitempty" json:"completed,omitempty" xml:"completed,omitempty"`
	Order     *int64  `bson:"order,omitempty" json:"order,omitempty" xml:"order,omitempty"`
	Title     *string `bson:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
}

type PatchTodoRequest struct {
	TodoId    int64
	TodoPatch Object4
}

type PatchTodoResponse interface {
	isPatchTodoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Successful
type PatchTodo200Response struct {
	Body Todo
}

func (r *PatchTodo200Response) isPatchTodoResponse() {}

func (r *PatchTodo200Response) StatusCode() int {
	return 200
}

func (r *PatchTodo200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not found
type PatchTodo404Response struct{}

func (r *PatchTodo404Response) isPatchTodoResponse() {}

func (r *PatchTodo404Response) StatusCode() int {
	return 404
}

func (r *PatchTodo404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}
