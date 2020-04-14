package api

import (
	"io"
	"net/http"
)

var contentTypesForFiles = []string{"application/json", "image/png", "image/jpeg", "image/tiff", "image/webp", "image/svg+xml", "image/gif", "image/tiff", "image/x-icon", "application/pdf"}

type Address struct {
	City        string `bson:"city,required" json:"city,required" xml:"city,required"`
	Country     string `bson:"country,required" json:"country,required" xml:"country,required"`
	HouseNumber string `bson:"houseNumber,required" json:"houseNumber,required" xml:"houseNumber,required"`
	PostalCode  string `bson:"postalCode,required" json:"postalCode,required" xml:"postalCode,required"`
	Region      string `bson:"region,required" json:"region,required" xml:"region,required"`
	Street      string `bson:"street,required" json:"street,required" xml:"street,required"`
}

type BasicTypes struct {
	Boolean bool              `bson:"boolean,required" json:"boolean,required" xml:"boolean,required"`
	Integer int64             `bson:"integer,required" json:"integer,required" xml:"integer,required"`
	Map     map[string]string `bson:"map,required" json:"map,required" xml:"map,required"`
	Number  float64           `bson:"number,required" json:"number,required" xml:"number,required"`
	Slice   []string          `bson:"slice,required" json:"slice,required" xml:"slice,required"`
	String  string            `bson:"string,required" json:"string,required" xml:"string,required"`
}

type Booking struct {
	BookingID *string `bson:"bookingID,omitempty" json:"bookingID,omitempty" xml:"bookingID,omitempty"`
}

type Object1 struct {
	BbdCEBaseUrl            *string  `bson:"bbdCEBaseUrl,omitempty" json:"bbdCEBaseUrl,omitempty" xml:"bbdCEBaseUrl,omitempty"`
	BbdCallerIdentifier     *string  `bson:"bbdCallerIdentifier,omitempty" json:"bbdCallerIdentifier,omitempty" xml:"bbdCallerIdentifier,omitempty"`
	BbdDataSupply           *string  `bson:"bbdDataSupply,omitempty" json:"bbdDataSupply,omitempty" xml:"bbdDataSupply,omitempty"`
	BbdImageBackground      *string  `bson:"bbdImageBackground,omitempty" json:"bbdImageBackground,omitempty" xml:"bbdImageBackground,omitempty"`
	BbdImagePerspective     *string  `bson:"bbdImagePerspective,omitempty" json:"bbdImagePerspective,omitempty" xml:"bbdImagePerspective,omitempty"`
	BbdImageType            *string  `bson:"bbdImageType,omitempty" json:"bbdImageType,omitempty" xml:"bbdImageType,omitempty"`
	BbdPassword             *string  `bson:"bbdPassword,omitempty" json:"bbdPassword,omitempty" xml:"bbdPassword,omitempty"`
	BbdProductGroup         *string  `bson:"bbdProductGroup,omitempty" json:"bbdProductGroup,omitempty" xml:"bbdProductGroup,omitempty"`
	BbdSoapMediaProviderUrl *string  `bson:"bbdSoapMediaProviderUrl,omitempty" json:"bbdSoapMediaProviderUrl,omitempty" xml:"bbdSoapMediaProviderUrl,omitempty"`
	BbdUser                 *string  `bson:"bbdUser,omitempty" json:"bbdUser,omitempty" xml:"bbdUser,omitempty"`
	CcoreServiceUrl         *string  `bson:"ccoreServiceUrl,omitempty" json:"ccoreServiceUrl,omitempty" xml:"ccoreServiceUrl,omitempty"`
	CryptKeys               []string `bson:"cryptKeys,omitempty" json:"cryptKeys,omitempty" xml:"cryptKeys,omitempty"`
	HealConfigurations      *bool    `bson:"healConfigurations,omitempty" json:"healConfigurations,omitempty" xml:"healConfigurations,omitempty"`
}

type Client struct {
	ActivePresets *string `bson:"activePresets,omitempty" json:"activePresets,omitempty" xml:"activePresets,omitempty"`
	Configuration Object1 `bson:"configuration,omitempty" json:"configuration,omitempty" xml:"configuration,omitempty"`
	Id            string  `bson:"id,required" json:"id,required" xml:"id,required"`
	Name          string  `bson:"name,required" json:"name,required" xml:"name,required"`
}

// The kind of drive concept of a vehicle. Where UNDEFINED is used as the default and/or error case.
type DriveConcept string

const (
	DriveConceptCOMBUSTOR DriveConcept = "COMBUSTOR"
	DriveConceptHYBRID    DriveConcept = "HYBRID"
	DriveConceptELECTRIC  DriveConcept = "ELECTRIC"
	DriveConceptFUELCELL  DriveConcept = "FUELCELL"
	DriveConceptUNDEFINED DriveConcept = "UNDEFINED"
)

type EmptySlice struct {
	EmptySlice []Price `bson:"EmptySlice,omitempty" json:"EmptySlice,omitempty" validate:"omitempty,gt=0,dive" xml:"EmptySlice,omitempty"`
}

type Link struct {
	Href string `bson:"href,required" json:"href,required" xml:"href,required"`
}

type Links struct {
	Self Link `bson:"self,required" json:"self,required" xml:"self,required"`
}

type Model struct {
	DriveConcept         *DriveConcept        `bson:"driveConcept,omitempty" json:"driveConcept,omitempty" xml:"driveConcept,omitempty"`
	Price                Price                `bson:"price,required" json:"price,required" xml:"price,required"`
	TechnicalInformation TechnicalInformation `bson:"technicalInformation,required" json:"technicalInformation,required" xml:"technicalInformation,required"`
}

type NestedFileStructure struct {
	Data *string `bson:"data,omitempty" json:"data,omitempty" xml:"data,omitempty"`
}

type Price struct {
	Currency string  `bson:"currency,required" json:"currency,required" xml:"currency,required"`
	Value    float64 `bson:"value,required" json:"value,required" xml:"value,required"`
}

type Rental struct {
	Class           string  `bson:"class,required" json:"class,required" validate:"min=3,max=20" xml:"class,required"`
	Color           *string `bson:"color,omitempty" json:"color,omitempty" validate:"omitempty,min=3,max=20" xml:"color,omitempty"`
	HomeID          *string `bson:"homeID,omitempty" json:"homeID,omitempty" validate:"omitempty,regex1" xml:"homeID,omitempty"`
	Id              string  `bson:"id,required" json:"id,required" validate:"regex2" xml:"id,required"`
	IdOptional      *string `bson:"idOptional,omitempty" json:"idOptional,omitempty" validate:"omitempty,regex3" xml:"idOptional,omitempty"`
	LockStatus      int32   `bson:"lockStatus,required" json:"lockStatus,required" validate:"min=1,max=100" xml:"lockStatus,required"`
	MaxDoors        int64   `bson:"maxDoors,required" json:"maxDoors,required" validate:"max=5" xml:"maxDoors,required"`
	MinDoors        int64   `bson:"minDoors,required" json:"minDoors,required" validate:"min=5" xml:"minDoors,required"`
	OptionalInt     *int64  `bson:"optionalInt,omitempty" json:"optionalInt,omitempty" xml:"optionalInt,omitempty"`
	State           *int64  `bson:"state,omitempty" json:"state,omitempty" xml:"state,omitempty"`
	StationID       string  `bson:"stationID,required" json:"stationID,required" validate:"regex4" xml:"stationID,required"`
	Status          int64   `bson:"status,required" json:"status,required" validate:"min=46,max=49" xml:"status,required"`
	Valid           *string `bson:"valid,omitempty" json:"valid,omitempty" validate:"omitempty,max=255" xml:"valid,omitempty"`
	Website         string  `bson:"website,required" json:"website,required" validate:"regex5" xml:"website,required"`
	WebsiteOptional *string `bson:"websiteOptional,omitempty" json:"websiteOptional,omitempty" validate:"omitempty,regex6,max=255" xml:"websiteOptional,omitempty"`
}

type Session struct {
	Registered bool   `bson:"Registered,required" json:"Registered,required" xml:"Registered,required"`
	Token      string `bson:"Token,required" json:"Token,required" xml:"Token,required"`
}

type Shoe struct {
	Links Links   `bson:"_links,required" json:"_links,required" xml:"_links,required"`
	Color string  `bson:"color,required" json:"color,required" xml:"color,required"`
	Name  string  `bson:"name,required" json:"name,required" xml:"name,required"`
	Size  float64 `bson:"size,required" json:"size,required" xml:"size,required"`
}

type Shoes struct {
	Embedded ShoesEmbedded `bson:"_embedded,required" json:"_embedded,required" xml:"_embedded,required"`
	Links    Links         `bson:"_links,required" json:"_links,required" xml:"_links,required"`
	Id       string        `bson:"id,required" json:"id,required" xml:"id,required"`
}

type ShoesEmbedded struct {
	ShopShoes []Shoe `bson:"shop:shoes,required" json:"shop:shoes,required" validate:"dive" xml:"shop:shoes,required"`
}

type TechnicalInformation struct {
	Transmission string `bson:"transmission,required" json:"transmission,required" xml:"transmission,required"`
}

type User struct {
	Address                []Address         `bson:"Address,omitempty" json:"Address,omitempty" validate:"omitempty,gt=0,dive" xml:"Address,omitempty"`
	Email                  *string           `bson:"email,omitempty" json:"email,omitempty" validate:"omitempty,email,max=255" xml:"email,omitempty"`
	GrantedProtocolMappers map[string]string `bson:"grantedProtocolMappers,omitempty" json:"grantedProtocolMappers,omitempty" xml:"grantedProtocolMappers,omitempty"`
	Id                     string            `bson:"id,required" json:"id,required" xml:"id,required"`
	Password               string            `bson:"password,required" json:"password,required" xml:"password,required"`
	Permissions            []string          `bson:"permissions,omitempty" json:"permissions,omitempty" xml:"permissions,omitempty"`
}

type ValidationError struct {
	Code    *string `bson:"Code,omitempty" json:"Code,omitempty" xml:"Code,omitempty"`
	Field   *string `bson:"Field,omitempty" json:"Field,omitempty" xml:"Field,omitempty"`
	Message *string `bson:"Message,omitempty" json:"Message,omitempty" xml:"Message,omitempty"`
}

type ValidationErrors struct {
	Errors  []ValidationError `bson:"Errors,omitempty" json:"Errors,omitempty" validate:"omitempty,gt=0,dive" xml:"Errors,omitempty"`
	Message *string           `bson:"Message,omitempty" json:"Message,omitempty" xml:"Message,omitempty"`
}

type ViewsSet struct {
	Id    string  `bson:"id,required" json:"id,required" xml:"id,required"`
	Name  *string `bson:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	Views *string `bson:"views,omitempty" json:"views,omitempty" xml:"views,omitempty"`
}

// ID of the request in UUIDv4 format
// A list of component types separated by a comma case insensitive. If nothing is defined all component types are returned.
type ComponentTypes string

const (
	ComponentTypesWHEELS           ComponentTypes = "WHEELS"
	ComponentTypesPAINTS           ComponentTypes = "PAINTS"
	ComponentTypesUPHOLSTERIES     ComponentTypes = "UPHOLSTERIES"
	ComponentTypesTRIMS            ComponentTypes = "TRIMS"
	ComponentTypesPACKAGES         ComponentTypes = "PACKAGES"
	ComponentTypesLINES            ComponentTypes = "LINES"
	ComponentTypesSPECIALEDITION   ComponentTypes = "SPECIAL_EDITION"
	ComponentTypesSPECIALEQUIPMENT ComponentTypes = "SPECIAL_EQUIPMENT"
)

// The productGroup of a vehicle case insensitive.
type ProductGroup string

const (
	ProductGroupPKW           ProductGroup = "PKW"
	ProductGroupGELAENDEWAGEN ProductGroup = "GELAENDEWAGEN"
	ProductGroupVAN           ProductGroup = "VAN"
	ProductGroupSPRINTER      ProductGroup = "SPRINTER"
	ProductGroupCITAN         ProductGroup = "CITAN"
	ProductGroupSMART         ProductGroup = "SMART"
)

type GetClientsRequest struct {
	XAuth string
}

type GetClientsResponse interface {
	isGetClientsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Status 200
type GetClients200Response struct {
	Body []Client `validate:"dive"`
}

func (r *GetClients200Response) isGetClientsResponse() {}

func (r *GetClients200Response) StatusCode() int {
	return 200
}

func (r *GetClients200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Status 201
type GetClients204Response struct{}

func (r *GetClients204Response) isGetClientsResponse() {}

func (r *GetClients204Response) StatusCode() int {
	return 204
}

func (r *GetClients204Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(204)
	return nil
}

// Not authenticated
type GetClients403Response struct{}

func (r *GetClients403Response) isGetClientsResponse() {}

func (r *GetClients403Response) StatusCode() int {
	return 403
}

func (r *GetClients403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

type DeleteClientRequest struct {
	ClientId string
	XAuth    string
}

type DeleteClientResponse interface {
	isDeleteClientResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type DeleteClient200Response struct{}

func (r *DeleteClient200Response) isDeleteClientResponse() {}

func (r *DeleteClient200Response) StatusCode() int {
	return 200
}

func (r *DeleteClient200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Not authenticated
type DeleteClient403Response struct{}

func (r *DeleteClient403Response) isDeleteClientResponse() {}

func (r *DeleteClient403Response) StatusCode() int {
	return 403
}

func (r *DeleteClient403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type DeleteClient404Response struct{}

func (r *DeleteClient404Response) isDeleteClientResponse() {}

func (r *DeleteClient404Response) StatusCode() int {
	return 404
}

func (r *DeleteClient404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type GetClientRequest struct {
	ClientId string
	XAuth    string
}

type GetClientResponse interface {
	isGetClientResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type GetClient200Response struct {
	Body Client
}

func (r *GetClient200Response) isGetClientResponse() {}

func (r *GetClient200Response) StatusCode() int {
	return 200
}

func (r *GetClient200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticated
type GetClient403Response struct{}

func (r *GetClient403Response) isGetClientResponse() {}

func (r *GetClient403Response) StatusCode() int {
	return 403
}

func (r *GetClient403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type GetClient404Response struct{}

func (r *GetClient404Response) isGetClientResponse() {}

func (r *GetClient404Response) StatusCode() int {
	return 404
}

func (r *GetClient404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type CreateOrUpdateClientRequest struct {
	ClientId string
	XAuth    string
	Body     Client
}

type CreateOrUpdateClientResponse interface {
	isCreateOrUpdateClientResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Updated
type CreateOrUpdateClient200Response struct{}

func (r *CreateOrUpdateClient200Response) isCreateOrUpdateClientResponse() {}

func (r *CreateOrUpdateClient200Response) StatusCode() int {
	return 200
}

func (r *CreateOrUpdateClient200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Created
type CreateOrUpdateClient201Response struct{}

func (r *CreateOrUpdateClient201Response) isCreateOrUpdateClientResponse() {}

func (r *CreateOrUpdateClient201Response) StatusCode() int {
	return 201
}

func (r *CreateOrUpdateClient201Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(201)
	return nil
}

// Malformed request body
type CreateOrUpdateClient400Response struct{}

func (r *CreateOrUpdateClient400Response) isCreateOrUpdateClientResponse() {}

func (r *CreateOrUpdateClient400Response) StatusCode() int {
	return 400
}

func (r *CreateOrUpdateClient400Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(400)
	return nil
}

// Not authenticated
type CreateOrUpdateClient403Response struct{}

func (r *CreateOrUpdateClient403Response) isCreateOrUpdateClientResponse() {}

func (r *CreateOrUpdateClient403Response) StatusCode() int {
	return 403
}

func (r *CreateOrUpdateClient403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not allowed
type CreateOrUpdateClient405Response struct{}

func (r *CreateOrUpdateClient405Response) isCreateOrUpdateClientResponse() {}

func (r *CreateOrUpdateClient405Response) StatusCode() int {
	return 405
}

func (r *CreateOrUpdateClient405Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(405)
	return nil
}

type GetViewsSetsRequest struct {
	ClientId string
	XAuth    string
}

type GetViewsSetsResponse interface {
	isGetViewsSetsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type GetViewsSets200Response struct {
	Body []ViewsSet `validate:"dive"`
}

func (r *GetViewsSets200Response) isGetViewsSetsResponse() {}

func (r *GetViewsSets200Response) StatusCode() int {
	return 200
}

func (r *GetViewsSets200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticated
type GetViewsSets403Response struct{}

func (r *GetViewsSets403Response) isGetViewsSetsResponse() {}

func (r *GetViewsSets403Response) StatusCode() int {
	return 403
}

func (r *GetViewsSets403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

type DeleteViewsSetRequest struct {
	ClientId string
	ViewsId  string
	XAuth    string
}

type DeleteViewsSetResponse interface {
	isDeleteViewsSetResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type DeleteViewsSet200Response struct{}

func (r *DeleteViewsSet200Response) isDeleteViewsSetResponse() {}

func (r *DeleteViewsSet200Response) StatusCode() int {
	return 200
}

func (r *DeleteViewsSet200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Not authenticated
type DeleteViewsSet403Response struct{}

func (r *DeleteViewsSet403Response) isDeleteViewsSetResponse() {}

func (r *DeleteViewsSet403Response) StatusCode() int {
	return 403
}

func (r *DeleteViewsSet403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type DeleteViewsSet404Response struct{}

func (r *DeleteViewsSet404Response) isDeleteViewsSetResponse() {}

func (r *DeleteViewsSet404Response) StatusCode() int {
	return 404
}

func (r *DeleteViewsSet404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type GetViewsSetRequest struct {
	ClientId string
	ViewsId  string
	Page     string
	XAuth    string
}

type GetViewsSetResponse interface {
	isGetViewsSetResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type GetViewsSet200Response struct {
	Body ViewsSet
}

func (r *GetViewsSet200Response) isGetViewsSetResponse() {}

func (r *GetViewsSet200Response) StatusCode() int {
	return 200
}

func (r *GetViewsSet200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticated
type GetViewsSet403Response struct{}

func (r *GetViewsSet403Response) isGetViewsSetResponse() {}

func (r *GetViewsSet403Response) StatusCode() int {
	return 403
}

func (r *GetViewsSet403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type GetViewsSet404Response struct{}

func (r *GetViewsSet404Response) isGetViewsSetResponse() {}

func (r *GetViewsSet404Response) StatusCode() int {
	return 404
}

func (r *GetViewsSet404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type ActivateViewsSetRequest struct {
	ClientId string
	ViewsId  string
	XAuth    string
}

type ActivateViewsSetResponse interface {
	isActivateViewsSetResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type ActivateViewsSet200Response struct{}

func (r *ActivateViewsSet200Response) isActivateViewsSetResponse() {}

func (r *ActivateViewsSet200Response) StatusCode() int {
	return 200
}

func (r *ActivateViewsSet200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Not authenticated
type ActivateViewsSet403Response struct{}

func (r *ActivateViewsSet403Response) isActivateViewsSetResponse() {}

func (r *ActivateViewsSet403Response) StatusCode() int {
	return 403
}

func (r *ActivateViewsSet403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type ActivateViewsSet404Response struct{}

func (r *ActivateViewsSet404Response) isActivateViewsSetResponse() {}

func (r *ActivateViewsSet404Response) StatusCode() int {
	return 404
}

func (r *ActivateViewsSet404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type CreateOrUpdateViewsSetRequest struct {
	ClientId string
	ViewsId  string
	XAuth    string
	Body     ViewsSet
}

type CreateOrUpdateViewsSetResponse interface {
	isCreateOrUpdateViewsSetResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Updated
type CreateOrUpdateViewsSet200Response struct{}

func (r *CreateOrUpdateViewsSet200Response) isCreateOrUpdateViewsSetResponse() {}

func (r *CreateOrUpdateViewsSet200Response) StatusCode() int {
	return 200
}

func (r *CreateOrUpdateViewsSet200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Created
type CreateOrUpdateViewsSet201Response struct{}

func (r *CreateOrUpdateViewsSet201Response) isCreateOrUpdateViewsSetResponse() {}

func (r *CreateOrUpdateViewsSet201Response) StatusCode() int {
	return 201
}

func (r *CreateOrUpdateViewsSet201Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(201)
	return nil
}

// Malformed request body
type CreateOrUpdateViewsSet400Response struct{}

func (r *CreateOrUpdateViewsSet400Response) isCreateOrUpdateViewsSetResponse() {}

func (r *CreateOrUpdateViewsSet400Response) StatusCode() int {
	return 400
}

func (r *CreateOrUpdateViewsSet400Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(400)
	return nil
}

// Not authenticated
type CreateOrUpdateViewsSet403Response struct{}

func (r *CreateOrUpdateViewsSet403Response) isCreateOrUpdateViewsSetResponse() {}

func (r *CreateOrUpdateViewsSet403Response) StatusCode() int {
	return 403
}

func (r *CreateOrUpdateViewsSet403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not allowed
type CreateOrUpdateViewsSet405Response struct{}

func (r *CreateOrUpdateViewsSet405Response) isCreateOrUpdateViewsSetResponse() {}

func (r *CreateOrUpdateViewsSet405Response) StatusCode() int {
	return 405
}

func (r *CreateOrUpdateViewsSet405Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(405)
	return nil
}

type ShowVehicleInViewRequest struct {
	ClientId   string
	ViewsId    string
	View       string
	Breakpoint string
	Spec       string
	XAuth      string
}

type ShowVehicleInViewResponse interface {
	isShowVehicleInViewResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type ShowVehicleInView200Response struct{}

func (r *ShowVehicleInView200Response) isShowVehicleInViewResponse() {}

func (r *ShowVehicleInView200Response) StatusCode() int {
	return 200
}

func (r *ShowVehicleInView200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Not authenticated
type ShowVehicleInView403Response struct{}

func (r *ShowVehicleInView403Response) isShowVehicleInViewResponse() {}

func (r *ShowVehicleInView403Response) StatusCode() int {
	return 403
}

func (r *ShowVehicleInView403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type ShowVehicleInView404Response struct{}

func (r *ShowVehicleInView404Response) isShowVehicleInViewResponse() {}

func (r *ShowVehicleInView404Response) StatusCode() int {
	return 404
}

func (r *ShowVehicleInView404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type GetPermissionsRequest struct {
	XAuth string
}

type GetPermissionsResponse interface {
	isGetPermissionsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Status 200
type GetPermissions200Response struct {
	Body []string
}

func (r *GetPermissions200Response) isGetPermissionsResponse() {}

func (r *GetPermissions200Response) StatusCode() int {
	return 200
}

func (r *GetPermissions200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticated
type GetPermissions403Response struct{}

func (r *GetPermissions403Response) isGetPermissionsResponse() {}

func (r *GetPermissions403Response) StatusCode() int {
	return 403
}

func (r *GetPermissions403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

type DestroySessionRequest struct {
	XAuth string
}

type DestroySessionResponse interface {
	isDestroySessionResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Session destroyed
type DestroySession200Response struct{}

func (r *DestroySession200Response) isDestroySessionResponse() {}

func (r *DestroySession200Response) StatusCode() int {
	return 200
}

func (r *DestroySession200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Session not found
type DestroySession404Response struct{}

func (r *DestroySession404Response) isDestroySessionResponse() {}

func (r *DestroySession404Response) StatusCode() int {
	return 404
}

func (r *DestroySession404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type GetUserInfoRequest struct {
	XAuth string `validate:"max=255"`
	SubID *int64 `validate:"omitempty,max=255"`
}

type GetUserInfoResponse interface {
	isGetUserInfoResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Status 200
type GetUserInfo200Response struct {
	Body User
}

func (r *GetUserInfo200Response) isGetUserInfoResponse() {}

func (r *GetUserInfo200Response) StatusCode() int {
	return 200
}

func (r *GetUserInfo200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Malformed request body
type GetUserInfo400Response struct {
	Body ValidationErrors
}

func (r *GetUserInfo400Response) isGetUserInfoResponse() {}

func (r *GetUserInfo400Response) StatusCode() int {
	return 400
}

func (r *GetUserInfo400Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 400, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticatedq
type GetUserInfo403Response struct{}

func (r *GetUserInfo403Response) isGetUserInfoResponse() {}

func (r *GetUserInfo403Response) StatusCode() int {
	return 403
}

func (r *GetUserInfo403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

type Object2 struct {
	Id       string `bson:"id,required" json:"id,required" validate:"min=1" xml:"id,required"`
	Password string `bson:"password,required" json:"password,required" validate:"min=1" xml:"password,required"`
}

type CreateSessionRequest struct {
	Body Object2
}

type CreateSessionResponse interface {
	isCreateSessionResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Authentication successful
type CreateSession200Response struct {
	XAuth string
}

func (r *CreateSession200Response) isCreateSessionResponse() {}

func (r *CreateSession200Response) StatusCode() int {
	return 200
}

func (r *CreateSession200Response) write(response http.ResponseWriter) error {
	response.Header()["X-Auth"] = []string{toString(r.XAuth)}
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Malformed request body
type CreateSession400Response struct {
	Body ValidationErrors
}

func (r *CreateSession400Response) isCreateSessionResponse() {}

func (r *CreateSession400Response) StatusCode() int {
	return 400
}

func (r *CreateSession400Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 400, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Authentication not successful
type CreateSession401Response struct{}

func (r *CreateSession401Response) isCreateSessionResponse() {}

func (r *CreateSession401Response) StatusCode() int {
	return 401
}

func (r *CreateSession401Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(401)
	return nil
}

type GetUsersRequest struct {
	XAuth string
}

type GetUsersResponse interface {
	isGetUsersResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type GetUsers200Response struct {
	Body []User `validate:"dive"`
}

func (r *GetUsers200Response) isGetUsersResponse() {}

func (r *GetUsers200Response) StatusCode() int {
	return 200
}

func (r *GetUsers200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticated
type GetUsers403Response struct{}

func (r *GetUsers403Response) isGetUsersResponse() {}

func (r *GetUsers403Response) StatusCode() int {
	return 403
}

func (r *GetUsers403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

type DeleteUserRequest struct {
	UserId  string
	AllKeys *bool
	XAuth   string
}

type DeleteUserResponse interface {
	isDeleteUserResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type DeleteUser200Response struct{}

func (r *DeleteUser200Response) isDeleteUserResponse() {}

func (r *DeleteUser200Response) StatusCode() int {
	return 200
}

func (r *DeleteUser200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Not authenticated
type DeleteUser403Response struct{}

func (r *DeleteUser403Response) isDeleteUserResponse() {}

func (r *DeleteUser403Response) StatusCode() int {
	return 403
}

func (r *DeleteUser403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type DeleteUser404Response struct{}

func (r *DeleteUser404Response) isDeleteUserResponse() {}

func (r *DeleteUser404Response) StatusCode() int {
	return 404
}

func (r *DeleteUser404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type GetUserRequest struct {
	UserId  string
	AllKeys *bool
	XAuth   string
}

type GetUserResponse interface {
	isGetUserResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success
type GetUser200Response struct {
	Body User
}

func (r *GetUser200Response) isGetUserResponse() {}

func (r *GetUser200Response) StatusCode() int {
	return 200
}

func (r *GetUser200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Not authenticated
type GetUser403Response struct{}

func (r *GetUser403Response) isGetUserResponse() {}

func (r *GetUser403Response) StatusCode() int {
	return 403
}

func (r *GetUser403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not found
type GetUser404Response struct{}

func (r *GetUser404Response) isGetUserResponse() {}

func (r *GetUser404Response) StatusCode() int {
	return 404
}

func (r *GetUser404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

type CreateOrUpdateUserRequest struct {
	UserId  string
	AllKeys *bool
	XAuth   string
	Body    User
}

type CreateOrUpdateUserResponse interface {
	isCreateOrUpdateUserResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Updated
type CreateOrUpdateUser200Response struct{}

func (r *CreateOrUpdateUser200Response) isCreateOrUpdateUserResponse() {}

func (r *CreateOrUpdateUser200Response) StatusCode() int {
	return 200
}

func (r *CreateOrUpdateUser200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Created
type CreateOrUpdateUser201Response struct{}

func (r *CreateOrUpdateUser201Response) isCreateOrUpdateUserResponse() {}

func (r *CreateOrUpdateUser201Response) StatusCode() int {
	return 201
}

func (r *CreateOrUpdateUser201Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(201)
	return nil
}

// Malformed request body
type CreateOrUpdateUser400Response struct{}

func (r *CreateOrUpdateUser400Response) isCreateOrUpdateUserResponse() {}

func (r *CreateOrUpdateUser400Response) StatusCode() int {
	return 400
}

func (r *CreateOrUpdateUser400Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(400)
	return nil
}

// Not authenticated
type CreateOrUpdateUser403Response struct{}

func (r *CreateOrUpdateUser403Response) isCreateOrUpdateUserResponse() {}

func (r *CreateOrUpdateUser403Response) StatusCode() int {
	return 403
}

func (r *CreateOrUpdateUser403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Not allowed
type CreateOrUpdateUser405Response struct{}

func (r *CreateOrUpdateUser405Response) isCreateOrUpdateUserResponse() {}

func (r *CreateOrUpdateUser405Response) StatusCode() int {
	return 405
}

func (r *CreateOrUpdateUser405Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(405)
	return nil
}

type GetBookingRequest struct {
	XSessionID string
}

type GetBookingResponse interface {
	isGetBookingResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// status 200
type GetBooking200Response struct {
	Body string
}

func (r *GetBooking200Response) isGetBookingResponse() {}

func (r *GetBooking200Response) StatusCode() int {
	return 200
}

func (r *GetBooking200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// status 400
type GetBooking400Response struct{}

func (r *GetBooking400Response) isGetBookingResponse() {}

func (r *GetBooking400Response) StatusCode() int {
	return 400
}

func (r *GetBooking400Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(400)
	return nil
}

// Unauthorized Session Token
type GetBooking401Response struct{}

func (r *GetBooking401Response) isGetBookingResponse() {}

func (r *GetBooking401Response) StatusCode() int {
	return 401
}

func (r *GetBooking401Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(401)
	return nil
}

// Resource Not Found
type GetBooking404Response struct{}

func (r *GetBooking404Response) isGetBookingResponse() {}

func (r *GetBooking404Response) StatusCode() int {
	return 404
}

func (r *GetBooking404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

// Malfunction (internal requirements not fulfilled)
type GetBooking500Response struct{}

func (r *GetBooking500Response) isGetBookingResponse() {}

func (r *GetBooking500Response) StatusCode() int {
	return 500
}

func (r *GetBooking500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type GetBookingsRequest struct {
	Ids        []int64
	Date       *string
	XSessionID string
}

type GetBookingsResponse interface {
	isGetBookingsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Success List Booking History
type GetBookings200Response struct {
	Body []Booking `validate:"dive"`
}

func (r *GetBookings200Response) isGetBookingsResponse() {}

func (r *GetBookings200Response) StatusCode() int {
	return 200
}

func (r *GetBookings200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// status 400
type GetBookings400Response struct{}

func (r *GetBookings400Response) isGetBookingsResponse() {}

func (r *GetBookings400Response) StatusCode() int {
	return 400
}

func (r *GetBookings400Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(400)
	return nil
}

// Unauthorized Session Token
type GetBookings401Response struct{}

func (r *GetBookings401Response) isGetBookingsResponse() {}

func (r *GetBookings401Response) StatusCode() int {
	return 401
}

func (r *GetBookings401Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(401)
	return nil
}

// Resource Not Found
type GetBookings404Response struct{}

func (r *GetBookings404Response) isGetBookingsResponse() {}

func (r *GetBookings404Response) StatusCode() int {
	return 404
}

func (r *GetBookings404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

// Malfunction (internal requirements not fulfilled)
type GetBookings500Response struct{}

func (r *GetBookings500Response) isGetBookingsResponse() {}

func (r *GetBookings500Response) StatusCode() int {
	return 500
}

func (r *GetBookings500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type ListModelsRequest struct {
	BrandId      string
	DriveConcept *DriveConcept
	LanguageId   *string
	ClassId      *string
	LineId       *string
	Ids          []int64
}

type ListModelsResponse interface {
	isListModelsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Ok
type ListModels200Response struct {
	Body []Model `validate:"dive"`
}

func (r *ListModels200Response) isListModelsResponse() {}

func (r *ListModels200Response) StatusCode() int {
	return 200
}

func (r *ListModels200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type GetClassesRequest struct {
	ProductGroup   ProductGroup
	ComponentTypes []ComponentTypes
}

type GetClassesResponse interface {
	isGetClassesResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Successful response
type GetClasses200Response struct {
	Body string
}

func (r *GetClasses200Response) isGetClassesResponse() {}

func (r *GetClasses200Response) StatusCode() int {
	return 200
}

func (r *GetClasses200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Successful response
type GetClasses400Response struct {
	Body string
}

func (r *GetClasses400Response) isGetClassesResponse() {}

func (r *GetClasses400Response) StatusCode() int {
	return 400
}

func (r *GetClasses400Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 400, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type CodeRequest struct {
	Session      string
	State        []int64
	ResponseMode *string
	Code         string
}

type CodeResponse interface {
	isCodeResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// TBD
type Code200Response struct {
	Body string
}

func (r *Code200Response) isCodeResponse() {}

func (r *Code200Response) StatusCode() int {
	return 200
}

func (r *Code200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// status 400
type Code400Response struct{}

func (r *Code400Response) isCodeResponse() {}

func (r *Code400Response) StatusCode() int {
	return 400
}

func (r *Code400Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(400)
	return nil
}

// Unauthorized Session code
type Code401Response struct{}

func (r *Code401Response) isCodeResponse() {}

func (r *Code401Response) StatusCode() int {
	return 401
}

func (r *Code401Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(401)
	return nil
}

// Resource Not Found
type Code404Response struct{}

func (r *Code404Response) isCodeResponse() {}

func (r *Code404Response) StatusCode() int {
	return 404
}

func (r *Code404Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(404)
	return nil
}

// Malfunction (internal requirements not fulfilled)
type Code500Response struct{}

func (r *Code500Response) isCodeResponse() {}

func (r *Code500Response) StatusCode() int {
	return 500
}

func (r *Code500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type DeleteCustomerSessionRequest struct {
	XRequestID *string
	XSessionID string
}

type DeleteCustomerSessionResponse interface {
	isDeleteCustomerSessionResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Session successful deleted
type DeleteCustomerSession204Response struct{}

func (r *DeleteCustomerSession204Response) isDeleteCustomerSessionResponse() {}

func (r *DeleteCustomerSession204Response) StatusCode() int {
	return 204
}

func (r *DeleteCustomerSession204Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(204)
	return nil
}

// Invalid session token
type DeleteCustomerSession401Response struct{}

func (r *DeleteCustomerSession401Response) isDeleteCustomerSessionResponse() {}

func (r *DeleteCustomerSession401Response) StatusCode() int {
	return 401
}

func (r *DeleteCustomerSession401Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(401)
	return nil
}

// Internal server error (e.g. unexpected condition occurred)
type DeleteCustomerSession500Response struct{}

func (r *DeleteCustomerSession500Response) isDeleteCustomerSessionResponse() {}

func (r *DeleteCustomerSession500Response) StatusCode() int {
	return 500
}

func (r *DeleteCustomerSession500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type CreateCustomerSessionRequest struct {
	XRequestID *string
	Code       string  `validate:"max=255"`
	Locale     *string `validate:"omitempty,regex7,max=255"`
}

type CreateCustomerSessionResponse interface {
	isCreateCustomerSessionResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Session successful created
type CreateCustomerSession201Response struct {
	Body Session
}

func (r *CreateCustomerSession201Response) isCreateCustomerSessionResponse() {}

func (r *CreateCustomerSession201Response) StatusCode() int {
	return 201
}

func (r *CreateCustomerSession201Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 201, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Invalid OpenID authentication token
type CreateCustomerSession401Response struct{}

func (r *CreateCustomerSession401Response) isCreateCustomerSessionResponse() {}

func (r *CreateCustomerSession401Response) StatusCode() int {
	return 401
}

func (r *CreateCustomerSession401Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(401)
	return nil
}

/*
Create session with authentication token is forbidden (e.g. Token already used)
*/
type CreateCustomerSession403Response struct{}

func (r *CreateCustomerSession403Response) isCreateCustomerSessionResponse() {}

func (r *CreateCustomerSession403Response) StatusCode() int {
	return 403
}

func (r *CreateCustomerSession403Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(403)
	return nil
}

// Invalid request data
type CreateCustomerSession422Response struct {
	Body ValidationErrors
}

func (r *CreateCustomerSession422Response) isCreateCustomerSessionResponse() {}

func (r *CreateCustomerSession422Response) StatusCode() int {
	return 422
}

func (r *CreateCustomerSession422Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 422, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Internal server error (e.g. unexpected condition occurred)
type CreateCustomerSession500Response struct{}

func (r *CreateCustomerSession500Response) isCreateCustomerSessionResponse() {}

func (r *CreateCustomerSession500Response) StatusCode() int {
	return 500
}

func (r *CreateCustomerSession500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type DownloadNestedFileRequest struct{}

type DownloadNestedFileResponse interface {
	isDownloadNestedFileResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Nested file structure
type DownloadNestedFile200Response struct {
	Body NestedFileStructure
}

func (r *DownloadNestedFile200Response) isDownloadNestedFileResponse() {}

func (r *DownloadNestedFile200Response) StatusCode() int {
	return 200
}

func (r *DownloadNestedFile200Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type DownloadImageRequest struct {
	Image string
}

type DownloadImageResponse interface {
	isDownloadImageResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// image to download
type DownloadImage200Response struct {
	Body        io.ReadCloser
	ContentType string
}

func (r *DownloadImage200Response) isDownloadImageResponse() {}

func (r *DownloadImage200Response) StatusCode() int {
	return 200
}

func (r *DownloadImage200Response) write(response http.ResponseWriter) error {
	response.Header()["Content-Type"] = []string{toString(r.ContentType)}
	if _, err := io.Copy(response, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	if err := r.Body.Close(); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Malfunction (internal requirements not fulfilled)
type DownloadImage500Response struct{}

func (r *DownloadImage500Response) isDownloadImageResponse() {}

func (r *DownloadImage500Response) StatusCode() int {
	return 500
}

func (r *DownloadImage500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type ListElementsRequest struct {
	Page    *int64
	PerPage *int64
}

type ListElementsResponse interface {
	isListElementsResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Status 200
type ListElements200Response struct {
	Body        string
	XTotalCount int64
}

func (r *ListElements200Response) isListElementsResponse() {}

func (r *ListElements200Response) StatusCode() int {
	return 200
}

func (r *ListElements200Response) write(response http.ResponseWriter) error {
	response.Header()["X-Total-Count"] = []string{toString(r.XTotalCount)}
	if err := serveJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Status 500
type ListElements500Response struct{}

func (r *ListElements500Response) isListElementsResponse() {}

func (r *ListElements500Response) StatusCode() int {
	return 500
}

func (r *ListElements500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type DownloadFileRequest struct {
	File string
}

type DownloadFileResponse interface {
	isDownloadFileResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// file to download
type DownloadFile200Response struct {
	Body        io.ReadCloser
	ContentType string
}

func (r *DownloadFile200Response) isDownloadFileResponse() {}

func (r *DownloadFile200Response) StatusCode() int {
	return 200
}

func (r *DownloadFile200Response) write(response http.ResponseWriter) error {
	response.Header()["Content-Type"] = []string{toString(r.ContentType)}
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

type GenericFileDownloadRequest struct {
	Ext string
}

type GenericFileDownloadResponse interface {
	isGenericFileDownloadResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// file to download
type GenericFileDownload200Response struct {
	Body        io.ReadCloser
	ContentType string
	Pragma      string
}

func (r *GenericFileDownload200Response) isGenericFileDownloadResponse() {}

func (r *GenericFileDownload200Response) StatusCode() int {
	return 200
}

func (r *GenericFileDownload200Response) write(response http.ResponseWriter) error {
	response.Header()["Content-Type"] = []string{toString(r.ContentType)}
	response.Header()["Pragma"] = []string{toString(r.Pragma)}
	if _, err := io.Copy(response, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	if err := r.Body.Close(); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

// Malfunction (internal requirements not fulfilled)
type GenericFileDownload500Response struct{}

func (r *GenericFileDownload500Response) isGenericFileDownloadResponse() {}

func (r *GenericFileDownload500Response) StatusCode() int {
	return 500
}

func (r *GenericFileDownload500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}

type GetRentalRequest struct {
	Body Rental
}

type GetRentalResponse interface {
	isGetRentalResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// status 200
type GetRental200Response struct{}

func (r *GetRental200Response) isGetRentalResponse() {}

func (r *GetRental200Response) StatusCode() int {
	return 200
}

func (r *GetRental200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// status 400
type GetRental400Response struct {
	Body ValidationErrors
}

func (r *GetRental400Response) isGetRentalResponse() {}

func (r *GetRental400Response) StatusCode() int {
	return 400
}

func (r *GetRental400Response) write(response http.ResponseWriter) error {
	if err := serveJson(response, 400, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type GetShoesRequest struct{}

type GetShoesResponse interface {
	isGetShoesResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Successful
type GetShoes200Response struct {
	Body Shoes
}

func (r *GetShoes200Response) isGetShoesResponse() {}

func (r *GetShoes200Response) StatusCode() int {
	return 200
}

func (r *GetShoes200Response) write(response http.ResponseWriter) error {
	if err := serveHalJson(response, 200, r.Body); err != nil {
		return NewHTTPStatusCodeError(http.StatusInternalServerError)
	}
	return nil
}

type PostUploadRequestFormData struct {
	Upfile *MimeFile
	Note   *string
}

type PostUploadRequest struct {
	FormData PostUploadRequestFormData
}

type PostUploadResponse interface {
	isPostUploadResponse()
	StatusCode() int
	write(response http.ResponseWriter) error
}

// Status 200
type PostUpload200Response struct{}

func (r *PostUpload200Response) isPostUploadResponse() {}

func (r *PostUpload200Response) StatusCode() int {
	return 200
}

func (r *PostUpload200Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(200)
	return nil
}

// Status 500
type PostUpload500Response struct{}

func (r *PostUpload500Response) isPostUploadResponse() {}

func (r *PostUpload500Response) StatusCode() int {
	return 500
}

func (r *PostUpload500Response) write(response http.ResponseWriter) error {
	response.Header()[contentTypeHeader] = []string{}
	response.WriteHeader(500)
	return nil
}
