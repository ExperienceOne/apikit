package openapi

type SecurityType string

func (s SecurityType) String() string {
	return string(s)
}

const (
	Basic  SecurityType = "basic"
	ApiKey SecurityType = "apiKey"
)

const BasicAuthHeaderName string = "Authorization"
