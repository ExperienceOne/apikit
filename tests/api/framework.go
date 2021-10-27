package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type HooksClient struct {
	OnUnknownResponseCode func(response *http.Response, request *http.Request) string
}

func DevHook() HooksClient {
	return HooksClient{
		OnUnknownResponseCode: func(response *http.Response, request *http.Request) string {
			var httpRequestDumpMessage string
			httpRequestDump, err := httputil.DumpRequest(request, true)
			if err != nil {
				httpRequestDumpMessage = fmt.Sprintf("could not dump request (%v)", err.Error())
			} else {
				httpRequestDumpMessage = string(httpRequestDump)
			}

			var httpResponseDumpMessage string
			httpResponseDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				httpResponseDumpMessage = fmt.Sprintf("could not dump response (%v)", err.Error())
			} else {
				httpResponseDumpMessage = string(httpResponseDump)
			}

			message := fmt.Sprintf("unknown response status code %d", response.StatusCode)
			if len(httpRequestDump) != 0 {
				message = message + "\n HTTP Request: \n '" + string(httpResponseDumpMessage) + "' \n"
			}
			if len(httpResponseDump) != 0 {
				message = message + "HTTP Response: \n '" + string(httpRequestDumpMessage) + "'"
			}
			return message
		},
	}
}
func primitiveToString(param reflect.Value) string {

	var value string

	switch param.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = fmt.Sprintf("%d", param.Interface())
	case reflect.Float64:
		value = strconv.FormatFloat(param.Interface().(float64), 'f', -1, 64)
	case reflect.Float32:
		value = strconv.FormatFloat(float64(param.Interface().(float32)), 'f', -1, 32)
	case reflect.String:
		value = fmt.Sprintf("%s", param.Interface())
	case reflect.Bool:
		value = fmt.Sprintf("%t", param.Interface())
	}

	return value
}

func sliceToString(param reflect.Value) string {

	slice := make([]string, param.Len())
	for i := 0; i < param.Len(); i++ {
		slice[i] = toString(param.Index(i).Interface())
	}
	return strings.Join(slice, ",")
}

func toString(param interface{}) string {

	paramReflected := reflect.ValueOf(param)

	for paramReflected.Kind() == reflect.Ptr {
		if paramReflected.IsNil() {
			return ""
		}
		paramReflected = paramReflected.Elem()
	}

	var value string
	if paramReflected.Kind() == reflect.Slice || paramReflected.Kind() == reflect.Array {
		value = sliceToString(paramReflected)
	} else {
		value = primitiveToString(paramReflected)
	}

	return value
}

func stringToPrimitive(s string, param reflect.Value) error {

	if param.Kind() != reflect.Ptr {
		return &ErrValueIsNotPointer{value: param}
	}

	var err error
	elm := param.Elem()

	switch elm.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := int64(0)
		if s != "" {
			val, err = strconv.ParseInt(s, 0, 64)
			if err != nil {
				return err
			}
			if elm.OverflowInt(val) {
				return &ErrTypeValueOverflow{value: s}
			}
		}
		elm.SetInt(val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := uint64(0)
		if s != "" {
			val, err = strconv.ParseUint(s, 0, 64)
			if err != nil {
				return err
			}
			if elm.OverflowUint(val) {
				return &ErrTypeValueOverflow{value: s}
			}
		}
		elm.SetUint(val)

	case reflect.Float32:
		val := float64(0)
		if s != "" {
			val, err = strconv.ParseFloat(s, 32)
			if err != nil {
				return err
			}
		}
		elm.SetFloat(val)

	case reflect.Float64:
		val := float64(0)
		if s != "" {
			val, err = strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
		}
		elm.SetFloat(val)

	case reflect.String:
		elm.SetString(s)

	case reflect.Bool:
		val := false
		if s != "" {
			val, err = strconv.ParseBool(s)
			if err != nil {
				return err
			}
		}
		elm.SetBool(val)

	default:
		return &ErrUnsupportedPrimitiveType{value: param}
	}
	return nil
}

func stringToSlice(s string, param reflect.Value) error {

	values := strings.Split(s, ",")
	elemForInjection := reflect.New(param.Elem().Type().Elem())
	for _, value := range values {
		err := fromString(value, elemForInjection.Interface())
		if err != nil {
			return err
		}
		slice := reflect.Append(param.Elem(), elemForInjection.Elem())
		param.Elem().Set(slice)
	}

	return nil
}

func fromString(s string, param interface{}) (err error) {

	defer func() {
		if v := recover(); v != nil {
			stack := debug.Stack()
			err = &RecoverError{Err: v, Stack: stack}
		}
	}()

	for {
		paramReflected := reflect.ValueOf(param)
		if paramReflected.Kind() != reflect.Ptr {
			return ErrParamIsNotPointer
		}

		if paramReflected.IsNil() {
			return ErrParamIsNil
		}

		kindOfElement := paramReflected.Elem().Kind()
		if kindOfElement == reflect.Ptr {
			ptr := paramReflected.Elem()
			if ptr.IsNil() {
				ptr.Set(reflect.New(ptr.Type().Elem()))
			}
			param = ptr.Interface()
		} else {
			if kindOfElement == reflect.Slice {
				err = stringToSlice(s, paramReflected)
			} else if kindOfElement == reflect.Array || kindOfElement == reflect.Map {
				err = &ErrUnsupportedKind{kind: kindOfElement}
			} else {
				err = stringToPrimitive(s, paramReflected)
			}
			break
		}
	}

	return
}

var ErrParamIsNotPointer error = errors.New("param isn't a pointer")
var ErrParamIsNil error = errors.New("param is nil")

type ErrUnsupportedKind struct {
	kind reflect.Kind
}

func (err *ErrUnsupportedKind) Error() string {
	return fmt.Sprintf("unsupported kind: '%s'", err.kind.String())
}

type ErrValueIsNotPointer struct {
	value reflect.Value
}

func (err *ErrValueIsNotPointer) Error() string {
	return fmt.Sprintf("value is not a pointer: '%s'", err.value.Kind().String())
}

type ErrTypeValueOverflow struct {
	value string
}

func (err *ErrTypeValueOverflow) Error() string {
	return fmt.Sprintf("type overflow: '%s'", err.value)
}

type ErrUnsupportedPrimitiveType struct {
	value reflect.Value
}

func (err *ErrUnsupportedPrimitiveType) Error() string {
	return fmt.Sprintf("unsupported primitive type: '%s'", err.value.Kind().String())
}

type RecoverError struct {
	Err   interface{}
	Stack []byte
}

func (e *RecoverError) Error() string { return fmt.Sprintf("fromString panicked: %v", e.Err) }

var (
	NullError = errors.New("unexpected null value")
	TypeError = errors.New("unexpected type")
)

func JSON(r io.Reader, v interface{}, required bool) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%+v", r))
		}
	}()

	typ := reflect.TypeOf(v)

	if typ.Kind() != reflect.Ptr {
		return errors.New("please pass pointer to target")
	}

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	decoder := json.NewDecoder(r)

	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map || typ.Kind() == reflect.Struct {

		handleNull := func() error {
			if required {
				return NullError
			} else {
				reflect.ValueOf(v).Elem().Set(reflect.Zero(reflect.TypeOf(v).Elem()))
				return nil
			}
		}

		var err error
		var value reflect.Value

		if typ.Kind() == reflect.Slice {

			var abstractSlice []interface{}
			if err := decoder.Decode(&abstractSlice); err != nil {
				return err
			}

			if abstractSlice == nil {
				return handleNull()
			}

			value, err = slice2concrete(typ, abstractSlice)

		} else {

			var abstractMap map[string]interface{}
			if err := decoder.Decode(&abstractMap); err != nil {
				return err
			}

			if abstractMap == nil {
				return handleNull()
			}

			if typ.Kind() == reflect.Map {
				value, err = map2concrete(typ, abstractMap)
			} else {
				value, err = map2object(typ, abstractMap)
			}
		}

		if err != nil {
			return err
		}

		setValue(value, reflect.ValueOf(v).Elem())

	} else {

		if !required {
			return decoder.Decode(v)
		}

		if err := decoder.Decode(&v); err != nil {
			return err
		}

		if v == nil {
			return NullError
		}
	}

	return nil
}

func setValue(value, dest reflect.Value) {

	for dest.Kind() == reflect.Ptr {
		newDest := reflect.New(dest.Type().Elem())
		dest.Set(newDest)
		dest = newDest.Elem()
	}
	dest.Set(value)
}

func setMapValue(key, value, m reflect.Value) {

	mapValue := reflect.New(m.Type().Elem()).Elem()
	setValue(value, mapValue)
	m.SetMapIndex(key, mapValue)
}

func slice2concrete(typ reflect.Type, s []interface{}) (reflect.Value, error) {

	concretSlice := reflect.MakeSlice(typ, len(s), cap(s))

	for i, value := range s {
		concretValue, err := convert(value, typ.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		setValue(concretValue, concretSlice.Index(i))
	}

	return concretSlice, nil
}

func map2concrete(typ reflect.Type, m map[string]interface{}) (reflect.Value, error) {

	concreteMap := reflect.MakeMap(typ)

	for key, value := range m {

		concreteKey, err := convert(key, typ.Key())
		if err != nil {
			return reflect.Value{}, err
		}

		concreteValue, err := convert(value, typ.Elem())
		if err != nil {
			return reflect.Value{}, err
		}

		setMapValue(concreteKey, concreteValue, concreteMap)
	}

	return concreteMap, nil
}

func map2object(typ reflect.Type, m map[string]interface{}) (reflect.Value, error) {

	object := reflect.New(typ).Elem()

	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		if field.Anonymous || !object.Field(i).CanSet() {
			continue
		}

		required := false
		key := field.Name

		jsonTags := strings.Split(field.Tag.Get("json"), ",")
		if len(jsonTags) > 0 {

			if jsonTags[0] == "-" {
				continue
			} else if jsonTags[0] != "" {
				key = jsonTags[0]
			}

			for i := 1; i < len(jsonTags); i++ {
				if jsonTags[i] == "required" {
					required = true
					break
				}
			}
		}

		if value, exists := m[key]; exists && value != nil {

			concreteValue, err := convert(value, field.Type)
			if err != nil {
				return reflect.Value{}, err
			}
			setValue(concreteValue, object.Field(i))

		} else if required {

			return reflect.Value{}, NullError
		}
	}

	return object, nil
}

func convert(data interface{}, typ reflect.Type) (reflect.Value, error) {

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map || typ.Kind() == reflect.Struct {

		var err error
		var value reflect.Value

		if typ.Kind() == reflect.Slice {

			abstractSlice, isSlice := data.([]interface{})
			if !isSlice {
				return reflect.Value{}, TypeError
			}

			value, err = slice2concrete(typ, abstractSlice)

		} else if typ.Kind() == reflect.Map || typ.Kind() == reflect.Struct {

			abstractMap, isMap := data.(map[string]interface{})

			if !isMap {
				return reflect.Value{}, TypeError
			}

			if typ.Kind() == reflect.Map {
				value, err = map2concrete(typ, abstractMap)
			} else {
				value, err = map2object(typ, abstractMap)
			}

		}

		if err != nil {
			return reflect.Value{}, err
		}

		return value, nil

	} else {

		return reflect.ValueOf(data).Convert(typ), nil
	}
}

type ValidationErrorsObject struct {
	Message string                  `json:"message"`
	Errors  []ValidationErrorObject `json:"errors"`
}

type ValidationErrorObject struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Code    string `json:"code"`
}

func NewValidation() *Validator {
	return &Validator{
		validator.New(),
	}
}

type Validator struct {
	*validator.Validate
}

func (v *Validator) ValidateRequest(request interface{}) (*ValidationErrorsObject, error) {

	if err := v.Struct(request); err != nil {

		validationErrors := new(ValidationErrorsObject)
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, err
		}

		validationErrors.Message = "validation failed"
		for _, err := range errors {

			errorCode := fmt.Sprintf("invalid-%s", err.Tag())
			if err.Tag() == "required" {
				errorCode = err.Tag()
			}

			validationError := ValidationErrorObject{
				Message: fmt.Sprint(err),
				Field:   err.Field(),
				Code:    errorCode,
			}
			validationErrors.Errors = append(validationErrors.Errors, validationError)
		}

		return validationErrors, nil
	}

	return nil, nil
}

var (
	GitCommit string = "a8719e7d76227cf28d355fef25f8f275e6d51621"
	GitBranch string = "master"
	GitTag    string = "v1.0.0"
	BuildTime string = "Wed 27 Oct 2021 22:16:50 BST"
)

type VersionInfo struct {
	GoVersion string `json:"go_version"`
	GitTag    string `json:"git_tag"`
	GitCommit string `json:"git_commit"`
	GitBranch string `json:"git_branch"`
	BuildTime string `json:"build_time"`
}

func ApikitVersion() *VersionInfo {
	return &VersionInfo{
		GoVersion: runtime.Version(),
		GitTag:    GitTag,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		BuildTime: BuildTime,
	}
}

func (vi *VersionInfo) PrintTable() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)

	_, err := fmt.Fprintln(w, fmt.Sprintf("Go version: %s", vi.GoVersion))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Git tag: %s", vi.GitTag))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Git commit: %s", vi.GitCommit))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Git branch: %s", vi.GitBranch))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, fmt.Sprintf("Buildtime: %s", vi.BuildTime))
	if err != nil {
		return err
	}

	return w.Flush()
}

type Opts struct {
	Hooks HooksClient
	Ctx   context.Context
}

type httpClientWrapper struct {
	BaseURL string

	*http.Client
}

func newHttpClientWrapper(client *http.Client, baseUrl string) *httpClientWrapper {
	return &httpClientWrapper{
		Client:  client,
		BaseURL: baseUrl,
	}
}

type NewRequest func(string, io.Reader) (*http.Request, error)

func (c *httpClientWrapper) Verb(verb string) NewRequest {
	baseURL := c.BaseURL
	return func(endpoint string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequest(verb, baseURL+endpoint, body)
		if err != nil {
			return nil, err
		}
		return req, err
	}
}

func (c *httpClientWrapper) Get() NewRequest {
	return c.Verb(http.MethodGet)
}

func (c *httpClientWrapper) Into(body io.ReadCloser, r interface{}) error {
	err := json.NewDecoder(body).Decode(&r)
	if err != nil {
		return err
	}
	defer body.Close()
	return nil
}

const (
	contentTypeHeader                    string = "Content-Type"
	contentTypeApplicationJson           string = "application/json"
	contentTypeApplicationHalJson        string = "application/hal+json"
	ContentTypeTextPlain                 string = "text/plain"
	contentTypeMultipartFormData         string = "multipart/form-data"
	contentTypeApplicationFormUrlencoded string = "application/x-www-form-urlencoded"
)

func extractContentType(header string) string {
	if header == "" {
		return ""
	}
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	return strings.TrimSpace(strings.ToLower(header[:i]))
}

func contentTypeInList(types []string, typ string) bool {

	for _, t := range types {
		if t == typ {
			return true
		}
	}
	return false
}
func newNotSupportedContentType(statusCode int, message string) error {
	return &NotSupportedContentType{
		message:    message,
		statusCode: statusCode,
	}
}

type NotSupportedContentType struct {
	message    string
	statusCode int
}

func (e *NotSupportedContentType) Error() string {
	return fmt.Sprintf("error unsupported media type (%s)", e.message)
}

func (e *NotSupportedContentType) StatusCode() int {
	return e.statusCode
}

var newRequestObjectIsNilError = errors.New("request object is nil")

func newErrUnknownResponse(code int) *ErrUnknownResponse {
	return &ErrUnknownResponse{
		code: code,
	}
}

type ErrUnknownResponse struct {
	code int
}

func (err *ErrUnknownResponse) Error() string {
	return fmt.Sprintf("unknown response status code '%d'", err.code)
}

func newErrOnUnknownResponseCode(message string) *ErrOnUnknownResponseCode {
	return &ErrOnUnknownResponseCode{
		Message: message,
	}
}

type ErrOnUnknownResponseCode struct {
	Message string
}

func (err *ErrOnUnknownResponseCode) Error() string {
	return fmt.Sprintf(err.Message)
}
func serveJson(w http.ResponseWriter, status int, v interface{}) error {

	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func serveHalJson(w http.ResponseWriter, status int, v interface{}) error {

	w.Header()["Content-Type"] = []string{"application/hal+json"}
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

type MimeFile struct {
	Header  *multipart.FileHeader
	Content io.ReadCloser
}

func extractUpload(fileID string, r *http.Request) (*MimeFile, error) {

	if err := r.ParseMultipartForm(1024); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile(fileID)
	if err != nil {
		return nil, err
	}

	mimeFile := MimeFile{
		Header:  header,
		Content: file,
	}
	return &mimeFile, nil
}

type contextKey int

const (
	RequestHeaderKey contextKey = 1 + iota
)

type HttpContext interface {
	GetHTTPRequestHeaders() (http.Header, bool)
}

func CreateHttpContext(header http.Header) context.Context {
	ctx := context.Background()
	ctx = hTTPRequestHeaders(ctx, header)
	return ctx
}

func hTTPRequestHeaders(ctx context.Context, header http.Header) context.Context {
	return context.WithValue(ctx, RequestHeaderKey, header)
}

func newHttpContextWrapper(ctx context.Context) HttpContext {
	if ctx == nil {
		ctx = context.Background()
	}
	return &httpContext{
		ctx,
	}
}

type httpContext struct {
	context.Context
}

func (c *httpContext) GetHTTPRequestHeaders() (http.Header, bool) {
	header, ok := c.Value(RequestHeaderKey).(http.Header)
	return header, ok
}

var ErrCollisionMap error = errors.New("header from context overwrites header in request object")

func setRequestHeadersFromContext(httpContext HttpContext, header http.Header) error {

	if httpContext == nil {
		return nil
	}

	headersFromContext, ok := httpContext.GetHTTPRequestHeaders()
	if !ok {
		return nil
	}

	for key, values := range headersFromContext {
		if _, exists := header[key]; exists {
			return ErrCollisionMap
		}

		header[key] = values
	}
	return nil
}

type xHTTPError interface {
	error

	StatusCode() int
}

type httpCodeError struct {
	statusCode int
}

func NewHTTPStatusCodeError(status int) xHTTPError {
	return &httpCodeError{status}
}

func (e *httpCodeError) Error() string {
	return http.StatusText(e.statusCode)
}

func (e *httpCodeError) StatusCode() int {
	return e.statusCode
}

type HttpJsonError struct {
	statusCode int
	Message    interface{}
}

func newJsonHTTPError(status int, message interface{}) xHTTPError {
	return &HttpJsonError{statusCode: status, Message: message}
}

func (e *HttpJsonError) StatusCode() int {
	return e.statusCode
}

func (e *HttpJsonError) Error() string {
	return fmt.Sprintf("%s: %v", http.StatusText(e.statusCode), e.Message)
}

type PrometheusHandler struct {
	counter   *prometheus.CounterVec
	histogram *prometheus.HistogramVec
}

func NewPrometheusHandler(namespace *string) *PrometheusHandler {

	counterName := "api_request_number_total"
	histogramName := "api_request_duration_seconds"

	if namespace != nil {
		counterName = fmt.Sprintf("%s_%s", *namespace, counterName)
		histogramName = fmt.Sprintf("%s_%s", *namespace, histogramName)
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: counterName,
		Help: "Total number API requests sent by the service.",
	}, []string{"handler", "method", "status"})

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: histogramName,
		Help: "Duration of the API requests being handled",
	}, []string{"handler"})

	if err := prometheus.Register(counter); err != nil {
		logrus.WithError(err).Warn("failed to register prometheus counter")
	}

	if err := prometheus.Register(histogram); err != nil {
		logrus.WithError(err).Warn("failed to register prometheus histogram")
	}

	h := &PrometheusHandler{
		counter:   counter,
		histogram: histogram,
	}

	return h
}

func (h *PrometheusHandler) InitMetric(path, method string) {

	h.counter.WithLabelValues(path, method, strconv.Itoa(200))
	h.histogram.WithLabelValues(path)
}

func (h *PrometheusHandler) HandleRequest(path, method string, status int, duration time.Duration) {

	h.counter.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
	h.histogram.WithLabelValues(path).Observe(duration.Seconds())
}

type (
	Timeouts struct {
		ReadTimeout       time.Duration
		ReadHeaderTimeout time.Duration
		WriteTimeout      time.Duration
		IdleTimeout       time.Duration
	}

	ServerOpts struct {
		Timeouts
		ErrorHandler ErrorHandler
		Middleware   []Middleware
		OnStart      func(router *routing.Router)
		Prefix       string
	}

	Middleware struct {
		Handler routing.Handler
		After   bool
	}

	RouteDescription struct {
		Path       string
		Handler    routing.Handler
		Middleware []Middleware
		Method     string
	}

	Server struct {
		Timeouts
		ErrorLogger func(v ...interface{})
		OnStart     func(router *routing.Router)
		server      *http.Server
		Router      *routing.Router
		after       []routing.Handler
		before      []routing.Handler
		SwaggerSpec string
		Prefix      string
	}
)

type ErrorHandler func(v ...interface{})

func newServer(opts *ServerOpts) *Server {

	if opts == nil {
		return &Server{
			ErrorLogger: func(v ...interface{}) {},
		}
	}

	if opts.ErrorHandler == nil {
		opts.ErrorHandler = func(v ...interface{}) {}
	}

	server := &Server{
		ErrorLogger: opts.ErrorHandler,
		Prefix:      "",
	}

	if opts.OnStart != nil {
		server.OnStart = opts.OnStart
	}

	if opts.Prefix != "" {
		server.Prefix = opts.Prefix
	}

	if len(opts.Middleware) != 0 {
		before := make([]routing.Handler, 0)
		after := make([]routing.Handler, 0)

		for _, m := range opts.Middleware {
			if m.After {
				after = append(after, m.Handler)
			} else {
				before = append(before, m.Handler)
			}
		}

		server.after = after
		server.before = before
	}

	server.ReadTimeout = opts.ReadTimeout
	server.ReadHeaderTimeout = opts.ReadHeaderTimeout
	server.WriteTimeout = opts.WriteTimeout
	server.IdleTimeout = opts.IdleTimeout

	return server
}

func (server *Server) makeRouter(routes []RouteDescription) (*routing.Router, error) {

	router := routing.New()
	router.UseEscapedPath = true

	logError := func(format string, a ...interface{}) {
		msg := fmt.Sprintf(format, a...)
		server.ErrorLogger(msg)
	}

	var beforeStack []routing.Handler
	beforeStack = append(beforeStack, errorHandler(logError))
	if server.before != nil {
		beforeStack = append(beforeStack, server.before...)
	}

	prefix := server.Prefix

	if prefix == "/" {
		prefix = ""
	}

	rg := router.Group(prefix)

	rg.Use(beforeStack...)

	var afterStack []routing.Handler
	if server.after != nil {
		afterStack = append(afterStack, server.after...)
	}

	rg.Get("/spec", func(c *routing.Context) error {
		return c.Write(server.SwaggerSpec)
	})

	for _, route := range routes {

		var before, after []routing.Handler

		if route.Middleware != nil {
			for _, m := range route.Middleware {
				if m.After {
					after = append(after, m.Handler)
				} else {
					before = append(before, m.Handler)
				}
			}
		}

		var handler []routing.Handler
		handler = append(handler, before...)
		handler = append(handler, route.Handler)
		handler = append(handler, after...)
		handler = append(handler, afterStack...)

		rg.To(route.Method, route.Path, handler...)
	}

	return router, nil
}

func (server *Server) Start(port int, routes []RouteDescription) error {

	router, err := server.makeRouter(routes)
	if err != nil {
		return err
	}
	if server.OnStart != nil {
		server.OnStart(router)
	}
	server.Router = router

	httpServer := &http.Server{
		ReadTimeout:       server.ReadTimeout,
		ReadHeaderTimeout: server.ReadHeaderTimeout,
		WriteTimeout:      server.WriteTimeout,
		IdleTimeout:       server.IdleTimeout,
		Addr:              ":" + strconv.Itoa(port),
		Handler:           router,
	}
	server.server = httpServer

	return httpServer.ListenAndServe()
}

func (server *Server) Stop() error {

	if server.server != nil {
		deadline, _ := context.WithTimeout(context.TODO(), 30*time.Second)
		return server.server.Shutdown(deadline)
	}

	return nil
}

func errorHandler(logf fault.LogFunc) func(c *routing.Context) (err error) {
	return func(c *routing.Context) (err error) {
		defer func() {

			if e := recover(); e != nil {

				if logf != nil {
					logf("recovered from panic: %v", string(debug.Stack()))
				}
				c.Response.WriteHeader(http.StatusInternalServerError)
				err = nil
				c.Abort()

			} else if err != nil {

				switch errType := err.(type) {
				case *HttpJsonError:
					c.Response.Header()["Content-Type"] = []string{"application/json"}
					c.Response.WriteHeader(errType.StatusCode())
					if e := json.NewEncoder(c.Response).Encode(errType.Message); e != nil && logf != nil {
						logf("failed to write error message: %v", errType.Message)
					}
				case *httpCodeError:
					c.Response.Header()["Content-Type"] = []string{""}
					c.Response.WriteHeader(errType.StatusCode())
				case routing.HTTPError:
					c.Response.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
					c.Response.WriteHeader(errType.StatusCode())
					if _, e := c.Response.Write([]byte(errType.Error())); e != nil {
						logf("failed to write error message: %v", errType.Error())
					}
				}

				err = nil
				c.Abort()
			}
		}()

		return c.Next()
	}
}
