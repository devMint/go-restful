package response

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

var (
	// Ok (HTTP 200)
	// The request has succeeded. The information returned with the response is dependent on the method used in the request
	Ok = createDataResponse(http.StatusOK, true)

	// Created (HTTP 201)
	// The request has been fulfilled and resulted in a new resource being created. The newly created resource can be referenced by the URI(s) returned in the entity of the response, with the most specific URI for the resource given by a Location header field. The response SHOULD include an entity containing a list of resource characteristics and location(s) from which the user or user agent can choose the one most appropriate. The entity format is specified by the media type given in the Content-Type header field. The origin server MUST create the resource before returning the 201 status code. If the action cannot be carried out immediately, the server SHOULD respond with 202 (Accepted) response instead.
	Created = createDataResponse(http.StatusCreated, false)

	// NoContent (HTTP 204)
	// The server has fulfilled the request but does not need to return an entity-body, and might want to return updated metainformation. The response MAY include new or updated metainformation in the form of entity-headers, which if present SHOULD be associated with the requested variant.
	NoContent = createDataResponse(http.StatusNoContent, false)
)

var (
	// BadRequest (HTTP 400)
	// The request could not be understood by the server due to malformed syntax. The client SHOULD NOT repeat the request without modifications.
	BadRequest = createErrorResponse(http.StatusBadRequest)

	// Unauthorized (HTTP 401)
	// The request requires user authentication. The response MUST include a WWW-Authenticate header field containing a challenge applicable to the requested resource. The client MAY repeat the request with a suitable Authorization header field. If the request already included Authorization credentials, then the 401 response indicates that authorization has been refused for those credentials. If the 401 response contains the same challenge as the prior response, and the user agent has already attempted authentication at least once, then the user SHOULD be presented the entity that was given in the response, since that entity might include relevant diagnostic information. HTTP access authentication is explained in "HTTP Authentication: Basic and Digest Access Authentication".
	Unauthorized = createErrorResponse(http.StatusUnauthorized)

	// Forbidden (HTTP 403)
	// The server understood the request, but is refusing to fulfill it. Authorization will not help and the request SHOULD NOT be repeated. If the request method was not HEAD and the server wishes to make public why the request has not been fulfilled, it SHOULD describe the reason for the refusal in the entity. If the server does not wish to make this information available to the client, the status code 404 (Not Found) can be used instead.
	Forbidden = createErrorResponse(http.StatusForbidden)

	// NotFound (HTTP 404)
	// The server has not found anything matching the Request-URI. No indication is given of whether the condition is temporary or permanent. The 410 (Gone) status code SHOULD be used if the server knows, through some internally configurable mechanism, that an old resource is permanently unavailable and has no forwarding address. This status code is commonly used when the server does not wish to reveal exactly why the request has been refused, or when no other response is applicable.
	NotFound = createErrorResponse(http.StatusNotFound)

	// MethodNotAllowed (HTTP 405)
	// The method specified in the Request-Line is not allowed for the resource identified by the Request-URI. The response MUST include an Allow header containing a list of valid methods for the requested resource.
	MethodNotAllowed             = createErrorResponse(http.StatusMethodNotAllowed)
	NotAcceptable                = createErrorResponse(http.StatusNotAcceptable)
	ProxyAuthenticationRequired  = createErrorResponse(http.StatusProxyAuthRequired)
	RequestTimeout               = createErrorResponse(http.StatusRequestTimeout)
	Conflict                     = createErrorResponse(http.StatusConflict)
	Gone                         = createErrorResponse(http.StatusGone)
	LengthRequired               = createErrorResponse(http.StatusLengthRequired)
	PreconditionFailed           = createErrorResponse(http.StatusPreconditionFailed)
	RequestEntityTooLarge        = createErrorResponse(http.StatusRequestEntityTooLarge)
	RequestURITooLong            = createErrorResponse(http.StatusRequestURITooLong)
	UnsupportedMediaType         = createErrorResponse(http.StatusUnsupportedMediaType)
	RequestedRangeNotSatisfiable = createErrorResponse(http.StatusRequestedRangeNotSatisfiable)
	ExpectationFailed            = createErrorResponse(http.StatusExpectationFailed)
)

var (
	InternalServerError     = createErrorResponse(http.StatusInternalServerError)
	NotImplemented          = createErrorResponse(http.StatusNotImplemented)
	BadGateway              = createErrorResponse(http.StatusBadGateway)
	ServiceUnavailable      = createErrorResponse(http.StatusServiceUnavailable)
	GatewayTimeout          = createErrorResponse(http.StatusGatewayTimeout)
	HTTPVersionNotSupported = createErrorResponse(http.StatusHTTPVersionNotSupported)
)

type Response interface {
	StatusCode() int
	GetJSON() string
	GetXML() string

	Header
}

type errorMessage struct {
	rawHeaders `json:"-" xml:"-"`
	XMLName    xml.Name `json:"-" xml:"response"`
	Type       string   `json:"type,omitempty" xml:"type"`
	Title      string   `json:"title" xml:"title"`
	Detail     string   `json:"detail" xml:"detail"`
	Status     int      `json:"status" xml:"status"`
}

func (e errorMessage) StatusCode() int { return e.Status }

func (e errorMessage) GetJSON() string { return toJSON(e) }

func (e errorMessage) GetXML() string { return toXML(e) }

func (e errorMessage) Error() string { return e.Detail }

type dataResponse struct {
	rawHeaders `json:"-" xml:"-"`
	XMLName    xml.Name    `json:"-" xml:"response"`
	Data       interface{} `json:"data" xml:"data"`
	Status     int         `json:"-" xml:"-"`
}

func (o dataResponse) StatusCode() int { return o.Status }

func (o dataResponse) GetJSON() string { return toJSON(o) }

func (o dataResponse) GetXML() string { return toXML(o) }

func createDataResponse(statusCode int, renderBody bool) func(data ...interface{}) Response {
	return func(data ...interface{}) Response {
		var body interface{}
		if renderBody && len(data) == 1 {
			body = data[0]
		}
		if renderBody && len(data) > 1 {
			body = data
		}

		return dataResponse{
			Data:       body,
			Status:     statusCode,
			rawHeaders: rawHeaders{},
		}
	}
}

func createErrorResponse(statusCode int) func(err error) Response {
	return func(err error) Response {
		errResponse, ok := err.(errorMessage)
		if ok {
			return errResponse
		}

		return errorMessage{
			Type:       "http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
			Title:      http.StatusText(statusCode),
			Status:     statusCode,
			Detail:     err.Error(),
			rawHeaders: rawHeaders{},
		}
	}
}

func toJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}

	return string(b)
}

func toXML(v interface{}) string {
	c, err := xml.Marshal(v)
	if err != nil {
		return ""
	}

	return string(c)
}
