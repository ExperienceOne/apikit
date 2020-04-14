package basket

import "net/http"

var contentTypesForFiles = []string{"application/json", "image/png", "image/jpeg", "image/tiff", "image/webp", "image/svg+xml", "image/gif", "image/tiff", "image/x-icon", "application/pdf"}

type Item struct {
	Id    int64   `bson:"id,required" json:"id,required" xml:"id,required"`
	Name  string  `bson:"name,required" json:"name,required" xml:"name,required"`
	Price float64 `bson:"price,required" json:"price,required" xml:"price,required"`
}

type Items []Item
type PostItemRequest struct {
	Item *Item
}

type PostItemResponse interface {
	isPostItemResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// OK
type PostItem200Response struct{}

func (r *PostItem200Response) isPostItemResponse() {}

func (r *PostItem200Response) StatusCode() int {
	return 200
}

func (r *PostItem200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Internal server error (e.g. unexpected condition occurred)
type PostItem500Response struct{}

func (r *PostItem500Response) isPostItemResponse() {}

func (r *PostItem500Response) StatusCode() int {
	return 500
}

func (r *PostItem500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type GetItemsRequest struct{}

type GetItemsResponse interface {
	isGetItemsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Items posted
type GetItems200Response struct {
	Body Items
}

func (r *GetItems200Response) isGetItemsResponse() {}

func (r *GetItems200Response) StatusCode() int {
	return 200
}

func (r *GetItems200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}
