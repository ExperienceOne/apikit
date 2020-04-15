package validation_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/ExperienceOne/apikit/internal/framework/validation"
	"github.com/stretchr/testify/require"

	"github.com/go-playground/validator"
)

type Input1 struct {
	Email *string `json:"email,omitempty" bson:"email,omitempty" validate:"email"`
}

type Input2 struct {
	Name string `json:"Name" bson:"Name" validate:"required"`
}

type Input3 struct {
	Name string `json:"Name" bson:"Name" validate:"min=3"`
}

type Input4 struct {
	Name string `json:"Name" bson:"Name" validate:"max=3"`
}

type Input5 struct {
	State int `json:"State" bson:"State" validate:"min=3"`
}

type Input6 struct {
	State int `json:"State" bson:"State" validate:"max=3"`
}

type Input7 struct {
	State *int `json:"State" bson:"State" validate:"min=3"`
}

func TestValidator(t *testing.T) {

	emailInvalid := "invalidmail.de"
	pint := int(1)

	tests := []struct {
		name           string
		Input          interface{}
		ExpectedOutput validation.ValidationErrorsObject
	}{
		{
			name: "Broken email",
			Input: &Input1{
				Email: &emailInvalid,
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "Email",
						Code:  "invalid-email",
					},
				},
			},
		},
		{
			name: "name is required",
			Input: &Input2{
				Name: "",
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "Name",
						Code:  "required",
					},
				},
			},
		},
		{
			name: "name min is not reached",
			Input: &Input3{
				Name: "te",
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "Name",
						Code:  "invalid-min",
					},
				},
			},
		},
		{
			name: "name is max",
			Input: &Input4{
				Name: "te33",
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "Name",
						Code:  "invalid-max",
					},
				},
			},
		},
		{
			name: "state min is not reached",
			Input: &Input5{
				State: 1,
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "State",
						Code:  "invalid-min",
					},
				},
			},
		},
		{
			name: "state max is not reached",
			Input: &Input6{
				State: 120,
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "State",
						Code:  "invalid-max",
					},
				},
			},
		},
		{
			name: "state min is not reached",
			Input: &Input7{
				State: &pint,
			},
			ExpectedOutput: validation.ValidationErrorsObject{
				Message: "Validation of incoming request faild",
				Errors: []validation.ValidationErrorObject{
					{
						Field: "State",
						Code:  "invalid-min",
					},
				},
			},
		},
	}

	validator := validation.NewValidation()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationErrors, err := validator.ValidateRequest(test.Input)
			if err != nil {
				t.Error(err)
				return
			}

			if validationErrors == nil {
				t.Fatalf("validation errors is nil")
			}

			if len(validationErrors.Errors) > 1 {
				t.Errorf(`unepxected count of errors (%d)`, len(validationErrors.Errors))
			}

			if validationErrors.Errors[0].Field != test.ExpectedOutput.Errors[0].Field {
				t.Errorf(`unexpected field (actual: "%#v", expected: "%#v")`, validationErrors.Errors[0].Field, test.ExpectedOutput.Errors[0].Field)
			}

			if validationErrors.Errors[0].Code != test.ExpectedOutput.Errors[0].Code {
				t.Errorf(`unexpected field (actual: "%#v", expected: "%#v")`, validationErrors.Errors[0].Code, test.ExpectedOutput.Errors[0].Code)
			}
		})
	}
}

type CreateVehicleRequest struct {
	Vin string `json:"id" validate:"createVehicleRequest-Vin"`
}

var (
	validate  = validator.New()
	vinString = "^[A-Z0-9]$" // regex that compiles
	vinRegex  = regexp.MustCompile(vinString)
)

func TestRegexValidation(t *testing.T) {

	err := validate.RegisterValidation("createVehicleRequest-Vin", func(fl validator.FieldLevel) bool {
		return !vinRegex.MatchString(fl.Field().String())
	})
	require.Nil(t, err, "failed to register validation")

	req := CreateVehicleRequest{
		Vin: "!!!!!!!!!!!!!!",
	}

	if errs := validate.Struct(req); errs != nil {
		if err := errs.(validator.ValidationErrors); err != nil {
			if !strings.Contains(err.Error(), "CreateVehicleRequest.Vin") {
				t.Errorf("CreateVehicleRequest.Vin would not executed (%v)", err)
			}
		}
	}

	req = CreateVehicleRequest{
		Vin: "WDDFD51DF1S5D1S5G",
	}

	if errs := validate.Struct(req); errs != nil {
		if err := errs.(validator.ValidationErrors); err != nil {
			if strings.Contains(err.Error(), "CreateVehicleRequest.Vin") {
				t.Errorf("CreateVehicleRequest.Vin would executed (%v)", err)
			}
		}
	}
}

func TestValidatorChain(t *testing.T) {

	v := validation.NewValidation()

	regex1 := regexp.MustCompile(`^[a-zA-Z]*$`)
	regex1Callback := func(fl validator.FieldLevel) bool {
		return regex1.MatchString(fl.Field().String())
	}
	err := v.RegisterValidation("regex1", regex1Callback)
	require.Nil(t, err, "failed to register validation")

	type Rental struct {
		Color  *string `json:"color" bson:"color" validate:"omitempty,min=3,max=20"`
		HomeID *string `json:"homeID" bson:"homeID" validate:"omitempty,regex1"`
	}

	empty := ""
	homeID := "test"
	rental := Rental{Color: &empty, HomeID: &homeID}

	vErrs, err := v.ValidateRequest(rental)
	if err != nil {
		t.Fatal(err)
	}

	if vErrs == nil {
		t.Fatal("expected more errors")
	}

	if len(vErrs.Errors) != 1 || vErrs.Errors[0].Message != "Key: 'Rental.Color' Error:Field validation for 'Color' failed on the 'min' tag" {
		t.Error("unexpected errors")
		for _, err := range vErrs.Errors {
			t.Log(err.Message)
		}
	}
}
