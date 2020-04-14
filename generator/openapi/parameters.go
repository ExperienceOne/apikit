package openapi

type ParameterType string

func (p ParameterType) String() string {
	return string(p)
}

const (
	Query    ParameterType = "query"
	Path     ParameterType = "path"
	Header   ParameterType = "header"
	Body     ParameterType = "body"
	FormData ParameterType = "formData"
)
