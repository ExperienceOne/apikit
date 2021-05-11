// this file is not autogenerated

package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"

	routing "github.com/go-ozzo/ozzo-routing"
)

var (
	id            = "id"
	password      = "password"
	one           = "1"
	two           = "2"
	three         = "3"
	auth          = "01234567890"
	xSessionID    = "xsession"
	componentType = "WHEELS"
)

var triggerBadRequest = "trigger400"
var fromContext string = "fromContext"

type ReadCloserBuffer struct {
	*bytes.Buffer
}

func (cb *ReadCloserBuffer) Close() (err error) {
	//we don't actually have to do anything here, since the buffer is
	//just some data in memory
	//and the error is initialized to no-error
	return
}

func CreateSession(ctx context.Context, request *CreateSessionRequest) CreateSessionResponse {

	if request.Body.Id == "fromContext" {
		return &CreateSession200Response{
			XAuth: "Hell, yeah!",
		}
	}

	if request.Body.Id != id && request.Body.Password != password {
		return &CreateSession401Response{}
	}
	return &CreateSession200Response{
		XAuth: auth,
	}
}

func PostUpload(ctx context.Context, request *PostUploadRequest) PostUploadResponse {

	defer request.FormData.Upfile.Content.Close()

	buff, err := ioutil.ReadAll(request.FormData.Upfile.Content)
	if err != nil {
		log.Println(fmt.Sprintf("error reading uploaded file (%v)", err))
		return new(PostUpload500Response)
	}

	if string(buff) != "Hello world, file" {
		log.Println("content of file:")
		log.Println(string(buff))
		return new(PostUpload500Response)
	}

	return new(PostUpload200Response)
}

func GetUserInfo(ctx context.Context, request *GetUserInfoRequest) GetUserInfoResponse {

	if request.XAuth == triggerBadRequest {
		return &GetUserInfo400Response{}
	}

	if err := IsContextPopluated(ctx); err != nil {
		log.Println(fmt.Sprintf("error retrieving value from context (%v)", err))
		return &GetUserInfo400Response{}
	}

	if request.XAuth != auth {
		return &GetUserInfo403Response{}
	}

	return &GetUserInfo200Response{
		Body: User{
			GrantedProtocolMappers: map[string]string{"key": "value"},
			Id:                     id,
			Password:               password,
			Permissions:            []string{one, two, three},
		},
	}
}

var ErrPopulatingContext error = errors.New("error populating context")

func IsContextPopluated(ctx context.Context) error {
	raw := ctx.Value("dummy")
	value, ok := raw.(string)
	if !ok {
		return ErrPopulatingContext
	}

	if value != "dummy" {
		return fmt.Errorf("error populating context (actual: '%s', expected: '%s')", value, "dummy")
	}
	return nil
}

func GetBrands(ctx context.Context, request *ListModelsRequest) ListModelsResponse {
	if request.BrandId == "without_first_query_parameter" {
		if request.LanguageId == nil && request.ClassId != nil && request.LineId != nil {
			if len(request.Ids) != 3 {
				return nil
			}
			return new(ListModels200Response)
		}
		return nil
	} else {
		if request.LanguageId != nil {
			return new(ListModels200Response)
		}
	}
	return nil
}

func GetUsers(ctx context.Context, request *GetUsersRequest) GetUsersResponse {
	return &GetUsers200Response{}
}

func GetClient(ctx context.Context, request *GetClientRequest) GetClientResponse {
	if request.ClientId == fromContext {
		return &GetClient200Response{}
	}
	return nil
}

func DownloadImage(ctx context.Context, request *DownloadImageRequest) DownloadImageResponse {

	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))

	// Draw a red dot at (2, 3)
	img.Set(2, 3, color.RGBA{255, 0, 0, 255})

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Println(fmt.Sprintf("error creating temp dir (%v)", err))
		return new(DownloadImage500Response)
	}

	var picture bytes.Buffer
	if err := png.Encode(&picture, img); err != nil {
		log.Println(fmt.Sprintf("error encode raw picture (%v)", err))
		return new(DownloadImage500Response)
	}

	filename := fmt.Sprintf("%s/%s", dir, "test.png")
	if err := ioutil.WriteFile(filename, picture.Bytes(), 0644); err != nil {
		log.Println(fmt.Sprintf("error creating picture file (%v)", err))
		return new(DownloadImage500Response)
	}
	return &DownloadImage200Response{Body: &ReadCloserBuffer{&picture}, ContentType: "image/png"}
}

func DownloadFile(ctx context.Context, request *DownloadFileRequest) DownloadFileResponse {
	return nil
}

const JSONBody string = `
{"widget": {
    "debug": "on",
    "window": {
        "title": "Sample Konfabulator Widget",
        "name": "main_window",
        "width": 500,
        "height": 500
    },
    "image": { 
        "src": "Images/Sun.png",
        "name": "sun1",
        "hOffset": 250,
        "vOffset": 250,
        "alignment": "center"
    },
    "text": {
        "data": "Click Here",
        "size": 36,
        "style": "bold",
        "name": "text1",
        "hOffset": 250,
        "vOffset": 100,
        "alignment": "center",
        "onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
    }
}} 
`

func GenericFileDownload(ctx context.Context, request *GenericFileDownloadRequest) GenericFileDownloadResponse {
	switch request.Ext {
	case ".json":
		buf := new(bytes.Buffer)
		buf.Write([]byte(JSONBody))
		return &GenericFileDownload200Response{Body: &ReadCloserBuffer{buf}, ContentType: "application/json"}
	default:
		log.Println(fmt.Sprintf("error file extension is not supported (%v)", request.Ext))
		return new(GenericFileDownload500Response)
	}
}

func CreateOrUpdateUser(ctx context.Context, request *CreateOrUpdateUserRequest) CreateOrUpdateUserResponse {
	return &CreateOrUpdateUser201Response{}
}

func GetPermissions(ctx context.Context, request *GetPermissionsRequest) GetPermissionsResponse {
	return nil
}

func GetBookings(ctx context.Context, request *GetBookingsRequest) GetBookingsResponse {
	if request.XSessionID == "" {
		return &GetBookings401Response{}
	}
	if request.XSessionID != xSessionID {
		return &GetBookings401Response{}
	}

	bookingID := "3600278183"
	return &GetBookings200Response{
		Body: []Booking{
			{BookingID: &bookingID},
		},
	}
}

func Code(ctx context.Context, request *CodeRequest) CodeResponse {
	if request.Code == "" {
		return &Code400Response{}
	}
	if request.Session == "" {
		return &Code400Response{}
	}
	if len(request.State) == 0 {
		return &Code400Response{}
	}
	return &Code200Response{}
}

func CreateCustomerSession(ctx context.Context, request *CreateCustomerSessionRequest) CreateCustomerSessionResponse {

	if request.Code == "invalid_token" {
		return &CreateCustomerSession401Response{}
	}

	return &CreateCustomerSession201Response{
		Body: Session{
			Registered: true,
			Token:      "1234",
		},
	}
}

func GetClasses(ctx context.Context, request *GetClassesRequest) GetClassesResponse {
	if len(request.ComponentTypes[0]) == 0 || string(request.ComponentTypes[0]) != componentType {
		return &GetClasses400Response{}
	}
	return &GetClasses200Response{}
}

func ListElement(ctx context.Context, request *ListElementsRequest) ListElementsResponse {
	return &ListElements200Response{
		XTotalCount: 5,
	}
}

func GetBooking(ctx context.Context, request *GetBookingRequest) GetBookingResponse {
	return &GetBooking500Response{}
}

func GetRental(ctx context.Context, request *GetRentalRequest) GetRentalResponse {
	return nil
}

var errTriggerPanic error = errors.New("triggerPanic")

func RouterPanicMiddleware() routing.Handler {
	return func(c *routing.Context) error {

		if "/api/user" == c.Request.URL.Path {

			defer func() {
				v := recover()
				if v != nil {
					err, ok := v.(error)

					if !ok {
						log.Print("no trigger panic")
					}

					if err != nil {
						if !errors.Is(err, errTriggerPanic) {
							log.Print("error not a panic for /api/user")
						}
					}
				}
			}()

			panic(errTriggerPanic)
		}
		return c.Next()
	}
}

func RouterPopulateContextMiddleware() routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "dummy", "dummy")
		c.Request = c.Request.WithContext(ctx)
		return c.Next()
	}
}

func HandlerMiddleware() routing.Handler {
	return func(c *routing.Context) error {
		c.SetParam("clientId", fromContext)
		return c.Next()
	}
}

func NestedFileDownload(ctx context.Context, req *DownloadNestedFileRequest) DownloadNestedFileResponse {

	buf := new(bytes.Buffer)
	buf.Write([]byte(JSONBody))
	data := buf.String()

	return &DownloadNestedFile200Response{
		Body: NestedFileStructure{
			Data: &data,
		},
	}
}

func GetShoes(ctx context.Context, req *GetShoesRequest) GetShoesResponse {

	return &GetShoes200Response{Body: Shoes{}}
}

func FileUpload(ctx context.Context, req *FileUploadRequest) FileUploadResponse {

	defer req.FormData.File.Content.Close()

	buff, err := ioutil.ReadAll(req.FormData.File.Content)
	if err != nil {
		log.Println(fmt.Sprintf("error reading uploaded file (%v)", err))
		return new(FileUpload500Response)
	}

	if string(buff) != "Hello world, file" {
		log.Println("content of file:")
		log.Println(string(buff))
		return new(FileUpload500Response)
	}

	return new(FileUpload204Response)
}
