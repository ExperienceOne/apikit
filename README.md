![Experience One](/images/logo.png)

# ExperienceOne Golang APIKit

- [ExperienceOne Golang APIKit](#experienceone-golang-apikit)
  - [Overview](#overview)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Generate standard project structure](#generate-standard-project-structure)
    - [Define the API with OpenAPIv2](#define-the-api-with-openapiv2)
    - [Validate the OpenAPIv2 definition](#validate-the-openapiv2-definition)
    - [Generate API server, client and mock client](#generate-api-server-client-and-mock-client)
    - [Using the API client](#using-the-api-client)
    - [Using the API server](#using-the-api-server)
    - [Generate handler stubs](#generate-handler-stubs)
    - [Generate service stubs](#generate-service-stubs)
  - [Validation of request data](#validation-of-request-data)
    - [Passing information about invalid data to the client](#passing-information-about-invalid-data-to-the-client)
    - [Required and non-required fields](#required-and-non-required-fields)
    - [String validation](#string-validation)
    - [Integer validation](#integer-validation)
    - [Content types](#content-types)
  - [Enum support](#enum-support)
  - [Advanced features](#advanced-features)
    - [Error logging](#error-logging)
    - [Additional routes for the server](#additional-routes-for-the-server)
  - [Middleware components](#middleware-components)
    - [Server-side request and response logging](#server-side-request-and-response-logging)
    - [GDPR compliant request and response logging](#gdpr-compliant-request-and-response-logging)
    - [Client-side request and response logging](#client-side-request-and-response-logging)
    - [Up- and downloading files as streams](#up--and-downloading-files-as-streams)
  - [APIKit development notes](#apikit-development-notes)

## Overview

The APIKit enables the rapid development of Web APIs with Golang by providing the functionality to (re-)generate the communication layer of the API based on an [OpenAPIv2](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md) (Swagger) definition, both for the client and server part. It also helps with one-time generation of stubs for the server-side endpoint handlers.
The generated API code does fully handle (and thus hide) the HTTP layer of the application, so the developer can focus on business logic implementation. Integration between API core and HTTP layer is handled via generated structs and interfaces.
The API HTTP layer can be regenerated from the OpenAPIv2 definition without breaking the integration with the API core code. This enables a "definition first" approach that ensures a 100% match of OpenAPI / Swagger definition and the implemented API.

Specifically the APIKit contains a CLI to generate and update the following API related items:

- standard project directory structures and common files
- validation of OpenAPIv2 (Swagger) definition
- HTTP API client based on an OpenAPIv2 (Swagger) definition
- HTTP API server based on an OpenAPIv2 (Swagger) definition
- stubs for server API endpoint handlers

See the [compliance list](Compliance.md) for a detailed list of supported OpenAPIv2 features.

## Requirements

- Golang 1.13 or higher
- make (if you want to use the Makefile)

## Installation

```bash
git clone git@github.com:ExperienceOne/apikit.git
cd apikit
make install
```

## Usage

### Generate standard project structure

The command `apikit project <dest.dir> <path/of/package>` generates a standard project directory.

```bash
cd go/src
$GOPATH/bin/apikit project myproject myproject
```

The `project` command takes the name of the directory to create and the path of the main package as arguments It creates a CLI stub in the `cmd` subdirectory and a `.gitignore` file.

Note that it is not necessary to use the standard project structure with the APIKit.

### Define the API with OpenAPIv2

The next step is to define your API. It can be useful to use a [Swagger editor](https://editor.swagger.io/) for this task. By convention the Swagger definition is stored in the `doc` subdirectory.

For showcasing the following example definition will be used:

```yaml
---
swagger: "2.0"
info:
  title: "vis-admin"
  version: "1.0.0"
consumes:
- "application/json"
produces:
- "application/json"
paths:
  /api/client/{clientId}:
    put:
      summary: "Create or update client"
      operationId: "CreateOrUpdateClient"
      consumes: []
      parameters:
      - name: "X-Auth"
        in: "header"
        required: true
        type: "string"
        description: "Authentication token"
      - name: "clientId"
        in: "path"
        required: true
        type: "string"
        description: "Client ID"
      - name: "body"
        in: "body"
        required: true
        schema:
          $ref: "#/definitions/Client"
      responses:
        200:
          description: "Updated"
          schema:
            $ref: "#/definitions/Client"
        201:
          description: "Created"
          schema:
            $ref: "#/definitions/Client"
        400:
          description: "Malformed request body"
        403:
          description: "Not authenticated"
        405:
          description: "Not allowed"
definitions:
  Client:
    type: "object"
    required:
    - "id"
    - "name"
    properties:
      id:
        type: "string"
      name:
        type: "string"
      activePresets:
        type: "string"
```

### Validate the OpenAPIv2 definition

The command `apikit validate <api.yaml>` validates against the given OpenAPIv2 / Swagger definition.

```bash
cd myproject
$GOPATH/bin/apikit validate doc/myproject.yaml
```

### Generate API server, client and mock client

The `apikit generate <api.yaml> <dest.dir> <package> <flags>` 
command generates the client, mock client and server API based on an OpenAPIv2 definition.

* (Optional) Use flag `--only-client` to only generate the client component should be generated.
* (Optional) Use flag `--only-server` to only generate the server component should be generated.
* (Optional) Use flag `--mocked` to generate additionally a mocked client which is satisfying the interface of the client (interchangeable). 
  This flag works in combination with `--only-client` or without the flags `--only-client` and `--only-server`. 
  
Note:  that all endpoints need to have an operation ID to generate successfully. The generated code should not be edited manually. Instead, the OpenAPI definition should be updated and the source regenerated with APIKit `generate` command.

#### How to use the mock client?

Please read further about it on this page -> https://github.com/vektra/mockery

#### Example generate server, client and mocked client

```bash
$GOPATH/bin/apikit generate doc/myproject.yaml api api --mocked
```

#### Example generate server and client

```bash
$GOPATH/bin/apikit generate doc/myproject.yaml api api
```

#### Example generate only client and mocked client

```bash
$GOPATH/bin/apikit generate doc/myproject.yaml api api --only-client --mocked
```

#### Example generate only client

```bash
$GOPATH/bin/apikit generate doc/myproject.yaml api api --only-client
```

#### Example generate only server

```bash
$GOPATH/bin/apikit generate doc/myproject.yaml api api --only-server
```

### Using the API client

The `client.go` file contains an `interface` defining the programming interface of the API.

```golang
type VisAdminClient interface {
  CreateOrUpdateClient(request *CreateOrUpdateClientRequest) (CreateOrUpdateClientResponse, error)
}
```

The `request` parameter contains all data that is sent to the server with the API call. Every API call returns a response interface that can be casted to the actual, HTTP return code specific struct, and an `error` code if the communication with the server went wrong. The content of the response is only valid if the `error` return value is `nil`.

The `client.go` file also contains an implementation of the programming interface using HTTP calls to the actual API. There is a constructor for the client implementation which takes an `http.Client` and the API base URL as parameters. Via the `http.Client` network handling and additional `http.RoundTripper` can be configured and integrated with the API client.

The steps for using the client to make an API call are:

- create the client implementation via constructor
- fill request struct with appropriate data
- call the API via the function for the endpoint
- check error return value
- check HTTP status code of the response
- cast response to correct struct

This is illustrated in the following example:

```golang
package myproject_test

import (
  "myproject/api"
  "net/http"
  "testing"
)

const (
  baseUrl      = "http://api.myproject.com/"
  testClientId = "clientId"
  testApiKey   = "apiKey"
  testName     = "name"
)

func TestClient(t *testing.T) {

  // create instance of API client
  client := api.NewVisAdminClient(&http.Client{}, baseUrl, api.Opts{})

  // create request struct
  request := api.CreateOrUpdateClientRequest{
    ClientId: testClientId,
    XAuth:    testApiKey,
    Body: api.Client{
      Id:   testClientId,
      Name: testName,
    },
  }

  // call API with request
  response, err := client.CreateOrUpdateClient(&request)
  if err != nil {
    t.Fatal(err)
  }

  // check HTTP return code
  if response.StatusCode() != http.StatusOK {
    t.Fatalf("API returned error code %d", response.StatusCode())
  }

  // cast response appropriatly and use returned information
  t.Log(response.(*api.CreateOrUpdateClient200Response).Body.Name)
}
```

Note that for testing the API can be mocked by a local implementation of the API client programming interface.

### Using the API server

The file `server.go` contains an a full HTTP server that serves the specified API. Additional to the defined endpoints it contains the `/spec` endpoint that delivers the OpenAPIv2 / Swagger definition the server was generated with. The `/spec` endpoint can be used to visualize the API via UIs like [Swagger UI](https://swagger.io/tools/swagger-ui/).

The implementation of the API endpoints has to be done in handler functions that take the request data as parameter and return an implementation of the endpoint-specific response interface. The handler functions need to be registered before starting the server. If there is no handler registered for an endpoint, the server will respond with status code 404 as default.

The steps for using the API server are:

- create the server via constructor
- register endpoint handlers
- start listening on network socket

The following example shows how the server is initialized and than called by the generated client (see above).

```golang
package myproject_test

import (
  "context"
  "myproject/api"
  "net/http"
  "testing"
  "time"
)

const (
  baseUrl      = "http://localhost:8080"
  testClientId = "clientId"
  testApiKey   = "apiKey"
  testName     = "name"
)

func CreateOrUpdateClient(ctx context.Context, request *api.CreateOrUpdateClientRequest) api.CreateOrUpdateClientResponse {

  return &api.CreateOrUpdateClient200Response{
    Body: request.Body,
  }
}

func TestServer(t *testing.T) {

  // create server instance
  server := api.NewVisAdminServer(nil)

  // register endpoint handler
  server.SetCreateOrUpdateClientHandler(CreateOrUpdateClient)

  // start server to listen for requests
  go func() {
    t.Log(server.Start(8080))
  }()
  time.Sleep(1 * time.Second)

  // create instance of API client
  client := api.NewVisAdminClient(&http.Client{}, baseUrl, api.Opts{})

  // create request struct
  request := api.CreateOrUpdateClientRequest{
    ClientId: testClientId,
    XAuth:    testApiKey,
    Body: api.Client{
      Id:   testClientId,
      Name: testName,
    },
  }

  // call API with request
  response, err := client.CreateOrUpdateClient(&request)
  if err != nil {
    t.Fatal(err)
  }

  // check HTTP return code
  if response.StatusCode() != http.StatusOK {
    t.Fatalf("API returned error code %d", response.StatusCode())
  }

  // cast response appropriate and use returned information
  t.Log(response.(*api.CreateOrUpdateClient200Response).Body.Name)

  // stop server
  server.Stop()
}
```

### Generate handler stubs

For one-time generation of stubs for the endpoint handlers the command `apikit handlers <api.yaml> <dest.go> <package> <api/package/path>` can be used. Especially for large APIs it saves some typing work by writing out the boilerplate code for the handlers.

```bash
$GOPATH/bin/apikit handlers doc/myproject.yaml handlers.go myproject myproject/api
```

### Generate service stubs

For larger projects grouping the endpoint handlers in service classes is convenient. A service stub can be generated with `apikit service <api.yaml> <dest.go> <package> <tag> <api/package/path>`. The command will generate the endpoints that are tagged with `tag` only.

```bash
$GOPATH/bin/apikit service doc/myproject.yaml client_service.go myproject client myproject/api
```

## Validation of request data

The generated server does validate the request against the constraints defined in the OpenAPIv2 specification. If the validation fails, the server will respond with a `400 Bad Request` status code to the client.

### Passing information about invalid data to the client

To pass more information about the invalid data to the client, the `400 Bad Request` response with the body type `ValidationErrors` has to be added to the endpoint. The generated server will automatically use it to send details about malformed data to the client.

```yaml
# ... endpoint definition ...
      responses:
        400:
          description: "Bad request"
          schema:
            $ref: "#/definitions/ValidationErrors"

definitions:
  ValidationErrors:
    type: "object"
    properties:
      Message:
        type: "string"
      Errors:
        type: "array"
        items:
          $ref: "#/definitions/ValidationError"

  ValidationError:
    type: "object"
    properties:
      Message:
        type: "string"
      Field:
        type: "string"
      Code:
        type: "string"

```

### Required and non-required fields

The generated types do reflect if a field is required or not by making not required fields pointers. If a non-required field is not present in the request data or a field is present, but has the JSON value `null`, the pointer will be set to `nil`. Otherwise the pointer will point to the actual data. Required fields are not allowed to be not present or to be `null`. If they are, a `400 Bad Request` response is returned to the client.

The proper handling of required and non-required fields allows to correctly distinguish between zero values ("", 0, false) and non-present values. Note that the default status for fields in OpenAPIv2 is non-required. To avoid nil pointer checks mark as many fields `required` as appropriate.

### String validation

The string validation supports the `minLength`,  `maxLength` and the `pattern` attributes. `pattern` allows to specify a regular expression that the string has to match. If the string has the format `uuid`, `url` or `email`, it is automatically matched against the corresponding regular expression.

### Integer validation

The integer validation supports the `minimum`, `maximum`, `exclusiveMinimum` and `exclusiveMaximum` attributes.

### Array query parameters validation

The array query parameters validation also supports the `minItems` and `maxItems` attributes.

#### Example 
```yaml
   parameters:
     - name: tags 
       in: query
       description: Tags to filter by
       required: true
       type: array
       minItems: 2
       maxItems: 5
       items:
         type: string
```


### Content types

The client needs to send the correct content type header matching the consumes attribute in the OpenAPIv2 definition. If not, the server responds with a `415 Unsupported media type` error.

## Enum support

The APIKit generator has limited enum support. It is required that every enum definition occurs only once in the whole OpenAPIv2 specification. The best way to achieve this is to reference enum types where they are used and have a global definition of the enum type in the `definitions` section of the specification.

```yaml
 Car:
    type: object
    required:
      - driveConcept
    properties:
      driveConcept:
        $ref: '#/definitions/DriveConcept'
```

```yaml
DriveConcept:
    type: string
    enum:
      - COMBUSTOR
      - HYBRID
      - ELECTRIC
      - FUELCELL
      - UNDEFINED
```

This produces the following code in `types.go`:

```go
type DriveConcept string

const (
  DriveConceptCOMBUSTOR DriveConcept = "COMBUSTOR"
  DriveConceptHYBRID    DriveConcept = "HYBRID"
  DriveConceptELECTRIC  DriveConcept = "ELECTRIC"
  DriveConceptFUELCELL  DriveConcept = "FUELCELL"
  DriveConceptUNDEFINED DriveConcept = "UNDEFINED"
)
type (
  Car struct {
    DriveConcept         DriveConcept         `json:"driveConcept" bson:"driveConcept"`
  }
)
```

## URL query parameter defaults

Set URL query parameter defaults for integer and float type based values. 

Note: All other type default values are ignored during code generation. 
We may add support for new types in the future.

```yaml
     parameters:
      - name: "_page"
        in: "query"
        type: "integer"
        default: 1
      - name: "_perPage"
        in: "query"
        type: "integer"
        default: 10
```

## Advanced features

### Error logging

An error logger can be passed to the server via the `ErrorHandler` attribute of the `ServerOpts` struct. The logger will be called when something unexpected goes wrong inside the HTTP layer.

### Additional routes for the server

The `ServerOpts` struct can be used to pass a function to the generated server that is executed during server startup and has access to the internal router object. It is defined as `func(router *routing.Router)`.
With this function routes that are not defined in the OpenAPIv2 specification or are implemented by packages that come with an own HTTP handler can be added to the server.

```golang
import(
    ...
    "github.com/go-ozzo/ozzo-routing"
)

func NewServer() {

  onStart := func(router *routing.Router) {
    router.Any("/metrics", func(ctx *routing.Context) error {
      metrics.ServeHTTP(ctx.Response, ctx.Request)
      return nil
    })
  }

  server := api.NewVisAdminServer(&api.ServerOpts{OnStart: onStart})
```

## Middleware components

Functions that shall be executed every time an endpoint is called can be added to the server and to individual handlers via middleware components. A middleware is defined by a struct that holds the handler function (defined in `ozzo-routing`) and a flag, if the middleware shall be executed before or after the endpoint handler.

```golang
Middleware struct {
  Handler routing.Handler
  After   bool
}
```

Server-wide middleware components can be passed via the `ServerOpts` struct.

```golang
middleware := []api.Middleware{{Handler: middleware.RequestID().Handler}}
server := api.NewVisAdminServer(&api.ServerOpts{Middleware: middleware})
```

Handler-specific middleware can be added via a var-arg list at handler registration.

```golang
sessionMiddleware := xcustomer.Middleware{Handler: middleware.SessionHandler(s.sessions)}
server.SetCreateOrUpdateClientHandler(CreateOrUpdateClient, sessionMiddleware)
```

Middleware can use `c.Request.Context()` object to pass data to the handler function. For example the session handler can check for a valid session token and load session data into the `context` object of the handler if successful.

```golang
func SessionHandler(service session.Service) routing.Handler {

  return func(c *routing.Context) error {

    authToken := c.Request.Header.Get(XAuthTokenHeader)
    if len(authToken) == 0 {
      return xcustomer.NewHTTPStatusCodeError(http.StatusBadRequest)
    }

    sess, err := service.GetSession(authToken)
    if err != nil {
      if err.StatusCode == http.StatusInternalServerError {
        return xcustomer.NewHTTPStatusCodeError(http.StatusInternalServerError)
      }
      return xcustomer.NewHTTPStatusCodeError(http.StatusUnauthorized)
    }

    if sess == nil {
      return xcustomer.NewHTTPStatusCodeError(http.StatusUnauthorized)
    }

        // add session data to the context of the endpoint handler
    ctx := context.WithValue(c.Request.Context(), xcontext.SessionKey, sess)
    c.Request = c.Request.WithContext(ctx)
    return c.Next()
  }
}
```

### Server-side request and response logging

The APIKit `middleware` package provides a component for server-side request and response logging. The middleware records request and response and forwards it to a user-defined log function before sending the data to the client. `LogResponseWriter` is part of the APIKit `middleware` package.

```golang
rrLogger := func(r *http.Request, w *LogResponseWriter, elapsed time.Duration) {
    // ...  print request and response infos as suitable
}
middleware := []api.Middleware{{Handler: middleware.Log(rrLogger)}}
```

There is also a convenience wrapper for the log function.

```golang
timeFormat := time.RFC3339
pathsToIgnore := []string{...} // endpoints not to log
logger := logFunc func(entity middleware.LogEntry, values ...interface{}) {
    // use entity object to print log
}
rrLogger := middleware.NewLogger(timeFormat, pathsToIgnore, logger)
```

### GDPR compliant request and response logging

In order to comply with the GDPR no data that could identify a person must be stored in log files. Therefore, a middleware logger is provided that can mask specified fields in JSON request and response bodies.

```golang
fieldsToMask := []string{...} // list of JSON paths
typesToMask := []reflect.Type{...} // list of types
rrLogger := middleware.NewMaskingLogger(timeFormat, fieldsToMask, typesToMask, pathsToIgnore, logger)
```

### Client-side request and response logging

The APIKit `roundtripper` package contains components for client-side request and response logging. The simplest roundtripper forwards the HTTP request and response to a log function. The `roundtripper.Use()` function wraps the transport of an HTTP client with the given `RoundTripper`.

```golang
logFunc := func(req *http.Request, resp *http.Response) {
    // ... log request and response content ...
}
client := roundtripper.Use(&http.Client{}, roundtripper.Logger(logFunc))
```

Additionally there is a convenience roundtripper that does an `httputil.Dump..()` on request and response and forwards the result a generic logging function.

```golang
logger := func(a ...interface{}) {
  // log a[0]
}

client := roundtripper.Use(&http.Client{}, roundtripper.Dump(logger))
```

The logging function is defined as an interface var-arg list so that it's possible to use  standard functions of logging frameworks like `logrus`. The roundtripper does always call it with one string parameter.

### Up- and downloading files as streams

Up- and downloading binary data in JSON format requires to put whole files in memory while marshalling / unmarshalling the JSON. This can quickly overwhelm the server. A better approach is to handle files as streams. The APIKit supports this via the `type: file` attribute.

When uploading a file the content type `multipart/form-data` (`consumes` attribute) is required.

```yaml
  /upload:
    post:
      summary: Upload a file with others data
      consumes:
      - multipart/form-data
      operationId: PostUpload
      parameters:
      - in: formData
        name: upfile
        type: file
        description: the file to upload
      - in: formData
        name: note
        type: string
        maxLength: 4000
        pattern: "^[0-9a-zA-Z ]*$"
        description: Description of file
      responses:
        '200':
          description: Status 200
        '500':
          description: Status 500
```

The request struct that is passed to the endpoint handler now contains a FormData object that includes an `io.ReadCloser` to access the file content. Note that the `io.ReadCloser` needs to be closed by the handler to prevent memory leaks.

```golang
func PostUpload(ctx context.Context, request *PostUploadRequest) PostUploadResponse {

  defer request.FormData.Upfile.Content.Close()

  data, err := ioutil.ReadAll(request.FormData.Upfile.Content)
  if err != nil {
    log.Println(fmt.Sprintf("error reading uploaded file (%v)", err))
    return new(PostUpload500Response)
  }

    // .. do something with data
  return new(PostUpload200Response)
}
```

To deliver a file as stream to the client, the response body must have the `type: file`. In the example below the `Content-Type` header specified explicitly to overwrite the `Content-Type` generated by the `produces` attribute.

```yaml
  '/download/{image}':
    get:
      produces:
      - image/png
      summary: Retrieve a image
      description: Retrieve a image
      operationId: DownloadImage
      parameters:
      - in: path
        name: image
        description: The image name of the image
        type: string
        required: true
      responses:
        '200':
          description: image to download
          schema:
            type: file
          headers:
            Content-Type:
              type: "string"
        '500':
          description: Internal server error
```

The handle than has to pass an `io.ReadCloser` with the response. The open stream is closed in the HTTP layer of the generated server.

```golang
func DownloadImage(ctx context.Context, request *DownloadImageRequest) DownloadImageResponse {

    file, err := os.Open(request.Name)
    if err != nil {
        return &DownloadImage500Response{}
    }

  return &DownloadImage200Response{Body: file, ContentType: "image/png"}
}
```

## APIKit development notes

If you have found a bug or would like contributed to the project please check out our [Contribution Guidelines](/Contribution.md).

After making changes to the APIKit source do run the testsuite via the `Makefile`.

```bash
make test
```

When there were changes to the `internal/framework` package, update the `framework` files before committing and building the executable.

```bash
make framework
```

When debugging code generation consider enabling the debug mode of the `generate` command. When activated the generator will save the generated code regardless of errors into a temporary file and it print its location to the terminal.
