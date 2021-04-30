package unit

import (
	"testing"

	"github.com/ExperienceOne/apikit/tests/api"
	"github.com/icrowley/fake"
)

func TestInputValidation(t *testing.T) {
	string256 := fake.CharactersN(256)
	int256 := int64(256)

	website := "http://website"
	uuid := "03402834-23804"
	websiteOptional := ""

	tests := []struct {
		name         string
		input        interface{}
		countOfError int
	}{
		{
			name: "GetUserInfoRequest",
			input: api.GetUserInfoRequest{
				XAuth: string256,
				SubID: &int256,
			},
			countOfError: 2,
		},
		{
			name: "GetRentalRequest",
			input: api.GetRentalRequest{
				Body: api.Rental{
					Website:         website,
					Id:              uuid,
					WebsiteOptional: &websiteOptional,
				},
			},
			countOfError: 8,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !IsRequestInvalid(t, test.input, test.countOfError) {
				t.Fatalf("validation errors is empty")
			}
		})
	}
}

func TestInputIsEmptyValidation(t *testing.T) {

	uuidOptional := ""
	websiteOptional := ""
	emailOptional := ""

	tests := []struct {
		name         string
		input        interface{}
		countOfError int
	}{
		{
			name: "GetRentalRequest",
			input: api.GetRentalRequest{
				Body: api.Rental{
					IdOptional:      &uuidOptional,
					WebsiteOptional: &websiteOptional,
				},
			},
			countOfError: 7,
		},
		{
			name: "user",
			input: api.User{
				Email: &emailOptional,
			},
			countOfError: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !IsRequestInvalid(t, test.input, test.countOfError) {
				t.Fatalf("validation errors is empty")
			}
		})
	}
}

func TestInputIsEmptyBasicList(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		countOfError int
	}{
		{
			name:         "IsEmptySlice",
			input:        api.EmptySlice{},
			countOfError: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IsRequestInvalid(t, test.input, test.countOfError)
		})
	}
}

func TestInputIsNotEmptyBasicList(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		countOfError int
	}{
		{
			name: "IsEmptySlice",
			input: api.EmptySlice{EmptySlice: []api.Price{{
				Currency: "DE",
				Value:    1.53,
			}}},
			countOfError: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IsRequestInvalid(t, test.input, test.countOfError)
		})
	}
}

func TestInputRequiredButEmpty(t *testing.T) {

	// though all fields are marked as required we expect no errors
	// as nil values in the JSON have to be checked by the handler-layer
	// here at struct validation level, zero values are totally okay
	IsRequestInvalid(t, api.BasicTypes{}, 0)
}

func IsRequestInvalid(t *testing.T, request interface{}, countOfError int) bool {
	s := api.NewVisAdminServer(nil)

	validationErrors, err := s.Validator.ValidateRequest(request)
	if err != nil {
		t.Fatalf("could not validate request object")
	}

	if validationErrors == nil {
		return false
	}

	if len(validationErrors.Errors) != countOfError {
		t.Errorf("count of error is bad (got: %v, want: %v)", len(validationErrors.Errors), countOfError)
		for _, err := range validationErrors.Errors {
			t.Log("field:" + err.Message)
		}
	}

	return true
}
