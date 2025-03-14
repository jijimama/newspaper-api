// Package presenter provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package presenter

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// NewspaperCreateRequest defines model for NewspaperCreateRequest.
type NewspaperCreateRequest struct {
	ColumnName string `json:"columnName"`
	Title      string `json:"title"`
}

// NewspaperUpdateRequest defines model for NewspaperUpdateRequest.
type NewspaperUpdateRequest struct {
	ColumnName *string `json:"columnName,omitempty"`
	Title      *string `json:"title,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewspaperResponse defines model for NewspaperResponse.
type NewspaperResponse struct {
	ColumnName *string `json:"columnName,omitempty"`
	Id         *int    `json:"id,omitempty"`
	Title      *string `json:"title,omitempty"`
}

// NewspaperCreateRequestBody defines model for NewspaperCreateRequestBody.
type NewspaperCreateRequestBody = NewspaperCreateRequest

// NewspaperUpdateRequestBody defines model for NewspaperUpdateRequestBody.
type NewspaperUpdateRequestBody = NewspaperUpdateRequest

// CreateNewspaperJSONRequestBody defines body for CreateNewspaper for application/json ContentType.
type CreateNewspaperJSONRequestBody = NewspaperCreateRequest

// UpdateNewspaperByIdJSONRequestBody defines body for UpdateNewspaperById for application/json ContentType.
type UpdateNewspaperByIdJSONRequestBody = NewspaperUpdateRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// CreateNewspaperWithBody request with any body
	CreateNewspaperWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateNewspaper(ctx context.Context, body CreateNewspaperJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteNewspaperById request
	DeleteNewspaperById(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetNewspaperById request
	GetNewspaperById(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateNewspaperByIdWithBody request with any body
	UpdateNewspaperByIdWithBody(ctx context.Context, id int, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateNewspaperById(ctx context.Context, id int, body UpdateNewspaperByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) CreateNewspaperWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateNewspaperRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateNewspaper(ctx context.Context, body CreateNewspaperJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateNewspaperRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteNewspaperById(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteNewspaperByIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetNewspaperById(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetNewspaperByIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateNewspaperByIdWithBody(ctx context.Context, id int, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateNewspaperByIdRequestWithBody(c.Server, id, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateNewspaperById(ctx context.Context, id int, body UpdateNewspaperByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateNewspaperByIdRequest(c.Server, id, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewCreateNewspaperRequest calls the generic CreateNewspaper builder with application/json body
func NewCreateNewspaperRequest(server string, body CreateNewspaperJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateNewspaperRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateNewspaperRequestWithBody generates requests for CreateNewspaper with any type of body
func NewCreateNewspaperRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/newspapers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteNewspaperByIdRequest generates requests for DeleteNewspaperById
func NewDeleteNewspaperByIdRequest(server string, id int) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/newspapers/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetNewspaperByIdRequest generates requests for GetNewspaperById
func NewGetNewspaperByIdRequest(server string, id int) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/newspapers/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateNewspaperByIdRequest calls the generic UpdateNewspaperById builder with application/json body
func NewUpdateNewspaperByIdRequest(server string, id int, body UpdateNewspaperByIdJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateNewspaperByIdRequestWithBody(server, id, "application/json", bodyReader)
}

// NewUpdateNewspaperByIdRequestWithBody generates requests for UpdateNewspaperById with any type of body
func NewUpdateNewspaperByIdRequestWithBody(server string, id int, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/newspapers/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// CreateNewspaperWithBodyWithResponse request with any body
	CreateNewspaperWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateNewspaperResponse, error)

	CreateNewspaperWithResponse(ctx context.Context, body CreateNewspaperJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateNewspaperResponse, error)

	// DeleteNewspaperByIdWithResponse request
	DeleteNewspaperByIdWithResponse(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*DeleteNewspaperByIdResponse, error)

	// GetNewspaperByIdWithResponse request
	GetNewspaperByIdWithResponse(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*GetNewspaperByIdResponse, error)

	// UpdateNewspaperByIdWithBodyWithResponse request with any body
	UpdateNewspaperByIdWithBodyWithResponse(ctx context.Context, id int, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateNewspaperByIdResponse, error)

	UpdateNewspaperByIdWithResponse(ctx context.Context, id int, body UpdateNewspaperByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateNewspaperByIdResponse, error)
}

type CreateNewspaperResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *NewspaperResponse
	JSON400      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r CreateNewspaperResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateNewspaperResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteNewspaperByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *ErrorResponse
	JSON404      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r DeleteNewspaperByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteNewspaperByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetNewspaperByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *NewspaperResponse
	JSON400      *ErrorResponse
	JSON404      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r GetNewspaperByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetNewspaperByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateNewspaperByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *NewspaperResponse
	JSON400      *ErrorResponse
	JSON404      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r UpdateNewspaperByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateNewspaperByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// CreateNewspaperWithBodyWithResponse request with arbitrary body returning *CreateNewspaperResponse
func (c *ClientWithResponses) CreateNewspaperWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateNewspaperResponse, error) {
	rsp, err := c.CreateNewspaperWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateNewspaperResponse(rsp)
}

func (c *ClientWithResponses) CreateNewspaperWithResponse(ctx context.Context, body CreateNewspaperJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateNewspaperResponse, error) {
	rsp, err := c.CreateNewspaper(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateNewspaperResponse(rsp)
}

// DeleteNewspaperByIdWithResponse request returning *DeleteNewspaperByIdResponse
func (c *ClientWithResponses) DeleteNewspaperByIdWithResponse(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*DeleteNewspaperByIdResponse, error) {
	rsp, err := c.DeleteNewspaperById(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteNewspaperByIdResponse(rsp)
}

// GetNewspaperByIdWithResponse request returning *GetNewspaperByIdResponse
func (c *ClientWithResponses) GetNewspaperByIdWithResponse(ctx context.Context, id int, reqEditors ...RequestEditorFn) (*GetNewspaperByIdResponse, error) {
	rsp, err := c.GetNewspaperById(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetNewspaperByIdResponse(rsp)
}

// UpdateNewspaperByIdWithBodyWithResponse request with arbitrary body returning *UpdateNewspaperByIdResponse
func (c *ClientWithResponses) UpdateNewspaperByIdWithBodyWithResponse(ctx context.Context, id int, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateNewspaperByIdResponse, error) {
	rsp, err := c.UpdateNewspaperByIdWithBody(ctx, id, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateNewspaperByIdResponse(rsp)
}

func (c *ClientWithResponses) UpdateNewspaperByIdWithResponse(ctx context.Context, id int, body UpdateNewspaperByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateNewspaperByIdResponse, error) {
	rsp, err := c.UpdateNewspaperById(ctx, id, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateNewspaperByIdResponse(rsp)
}

// ParseCreateNewspaperResponse parses an HTTP response from a CreateNewspaperWithResponse call
func ParseCreateNewspaperResponse(rsp *http.Response) (*CreateNewspaperResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateNewspaperResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest NewspaperResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	}

	return response, nil
}

// ParseDeleteNewspaperByIdResponse parses an HTTP response from a DeleteNewspaperByIdWithResponse call
func ParseDeleteNewspaperByIdResponse(rsp *http.Response) (*DeleteNewspaperByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteNewspaperByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseGetNewspaperByIdResponse parses an HTTP response from a GetNewspaperByIdWithResponse call
func ParseGetNewspaperByIdResponse(rsp *http.Response) (*GetNewspaperByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetNewspaperByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest NewspaperResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseUpdateNewspaperByIdResponse parses an HTTP response from a UpdateNewspaperByIdWithResponse call
func ParseUpdateNewspaperByIdResponse(rsp *http.Response) (*UpdateNewspaperByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateNewspaperByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest NewspaperResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new newspaper
	// (POST /newspapers)
	CreateNewspaper(ctx echo.Context) error
	// Delete a newspaper by ID
	// (DELETE /newspapers/{id})
	DeleteNewspaperById(ctx echo.Context, id int) error
	// Find newspaper by ID
	// (GET /newspapers/{id})
	GetNewspaperById(ctx echo.Context, id int) error
	// Update a newspaper by ID
	// (PATCH /newspapers/{id})
	UpdateNewspaperById(ctx echo.Context, id int) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateNewspaper converts echo context to params.
func (w *ServerInterfaceWrapper) CreateNewspaper(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateNewspaper(ctx)
	return err
}

// DeleteNewspaperById converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteNewspaperById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteNewspaperById(ctx, id)
	return err
}

// GetNewspaperById converts echo context to params.
func (w *ServerInterfaceWrapper) GetNewspaperById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetNewspaperById(ctx, id)
	return err
}

// UpdateNewspaperById converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateNewspaperById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateNewspaperById(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/newspapers", wrapper.CreateNewspaper)
	router.DELETE(baseURL+"/newspapers/:id", wrapper.DeleteNewspaperById)
	router.GET(baseURL+"/newspapers/:id", wrapper.GetNewspaperById)
	router.PATCH(baseURL+"/newspapers/:id", wrapper.UpdateNewspaperById)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9SW30/bMBDH/5XotkerTgFpKG8DNtQXNKHtCfXBJEdrlNjGvoKiKv/7ZDsNKWn4tY5p",
	"4iXYd77vfe6HuoZcV0YrVOQgW4PFuxU6OtGFxHBwgQ/OCIP21KIgvOzua3+ba0WoyH8KY0qZC5Ja8Vun",
	"lT9z+RIr4b8+W7yBDD7xx3A83jq+OwQ0TdOwx/i/TPG342+FCPEbBhad0cpFGt+s1fayPXmTAGO1QUst",
	"1QqdE4vwAtUGIQNHVqoFxIh3K2mxgOyqM5yzjaG+vsU8qGNQoMutND4kZIBeHPSZ7UFprstVpS5EtUss",
	"A1n0jqUiXGCQQJLKkfRezENt5Ccb9sGvLdd4U75Z/DMi+zWIZqz/2rAcY526T1E7yEl1o4Nx9IJznXz9",
	"MUt+YmVKQV71PVoXuU4n6ST1IbRBJYyEDA4n6eQQGBhByyCOd/DDv0bHDLz+0C+zAjKI1Lt8gfXWRj02",
	"bVubhT+zVp7O3EE6HX+ztePDlm8YHKXpy57bIx0abVVVwtZdoolIFD4kqpcviYXzrdGjNfeuPXx8LYvG",
	"hy+wRMIhxbNw3gk/qWdFqIQVFVLgf7UG6QvnqwMMVGgZP3L9/iS7QtYb36ez2MwHPI+irP7MXejktF0Q",
	"7wTnvY7+DHckEnG3G+C6TmZnY8AZLHBHe54j/Quq6Qd26R5gf5eqeD1oIyhfDlHHVfdBtN+7Yoa/HJr/",
	"vXgxpddPSvBGe78pycqWkMGSyGScp5Pwlx2nxykXRvL7KTTsiVGpc1EutaPnzaYHX8Jr022zefM7AAD/",
	"/4QLOnRpCgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
