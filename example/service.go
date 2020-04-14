package basket

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service struct {
	server *BasketServiceServer
	items  []Item
}

func NewService() *Service {

	server := NewBasketServiceServer(&ServerOpts{})

	service := &Service{
		server: server,
		items:  make([]Item, 0),
	}

	return service
}

// Start starts the server to listen at port
func (s *Service) Start(port int) error {

	logrus.WithField("port", port).Info("starting basket service")

	s.server.SetGetItemsHandler(s.GetItems)
	s.server.SetPostItemHandler(s.PostItem)

	return s.server.Start(port)
}

// Stop stops the server to listen at port
func (s *Service) Stop() error {

	logrus.Info("stopping basket service")

	return s.server.Stop()
}

func (s *Service) GetItems(ctx context.Context, request *GetItemsRequest) GetItemsResponse {
	return &GetItems200Response{Body: s.items}
}

func (s *Service) PostItem(ctx context.Context, request *PostItemRequest) PostItemResponse {

	if request.Item == nil {
		return &PostItem500Response{}
	}

	s.items = append(s.items, *request.Item)

	return &PostItem200Response{}
}
