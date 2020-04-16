package todo

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Service struct {
	server *TodoServiceServer
	todos  TodoList
	nextId int
}

func NewService() *Service {

	server := NewTodoServiceServer(&ServerOpts{})

	service := &Service{
		server: server,
		todos:  make(TodoList, 0),
		nextId: 0,
	}

	return service
}

// Start starts the server to listen at port
func (s *Service) Start(port int) error {

	logrus.WithField("port", port).Info("starting todo service")

	s.server.SetListTodosHandler(s.ListTodos)
	s.server.SetPostTodoHandler(s.PostTodo)
	s.server.SetDeleteTodosHandler(s.DeleteTodos)
	s.server.SetGetTodoHandler(s.GetTodo)
	s.server.SetPatchTodoHandler(s.PatchTodo)
	s.server.SetDeleteTodoHandler(s.DeleteTodo)

	return s.server.Start(port)
}

// Stop stops the server to listen at port
func (s *Service) Stop() error {

	logrus.Info("stopping todo service")

	return s.server.Stop()
}

func (s *Service) ListTodos(ctx context.Context, request *ListTodosRequest) ListTodosResponse {
	return &ListTodos200Response{Body: s.todos}
}

func (s *Service) PostTodo(ctx context.Context, request *PostTodoRequest) PostTodoResponse {

	todo := Todo{
		Completed: false,
		Id:        int64(s.nextId),
		Order:     0,
		Title:     request.TodoPost.Title,
		Url:       fmt.Sprintf("todos/%d", s.nextId),
	}
	s.nextId++

	s.todos = append(s.todos, todo)

	return &PostTodo201Response{Body: todo}
}

func (s *Service) DeleteTodos(ctx context.Context, request *DeleteTodosRequest) DeleteTodosResponse {

	s.todos = s.todos[:0]

	return &DeleteTodos204Response{}
}

func (s *Service) GetTodo(ctx context.Context, request *GetTodoRequest) GetTodoResponse {

	for _, todo := range s.todos {
		if todo.Id == request.TodoId {
			return &GetTodo200Response{Body: todo}
		}
	}

	return &GetTodo404Response{}
}

func (s *Service) PatchTodo(ctx context.Context, request *PatchTodoRequest) PatchTodoResponse {

	for i, todo := range s.todos {
		if todo.Id == request.TodoId {

			if request.TodoPatch.Title != nil {
				todo.Title = *request.TodoPatch.Title
			}

			if request.TodoPatch.Order != nil {
				todo.Order = *request.TodoPatch.Order
			}

			if request.TodoPatch.Completed != nil {
				todo.Completed = *request.TodoPatch.Completed
			}

			s.todos[i] = todo

			return &PatchTodo200Response{Body: todo}
		}
	}

	return &PatchTodo404Response{}
}

func (s *Service) DeleteTodo(ctx context.Context, request *DeleteTodoRequest) DeleteTodoResponse {

	index := -1
	for i, todo := range s.todos {
		if todo.Id == request.TodoId {
			index = i
		}
	}

	if index == -1 {
		return &DeleteTodo404Response{}
	}

	s.todos[index] = s.todos[len(s.todos)-1]
	s.todos = s.todos[:len(s.todos)-1]

	return &DeleteTodo204Response{}
}
