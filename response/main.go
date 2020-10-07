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

	// Accepted (HTTP 202)
	// The request has been accepted for processing, but the processing has not been completed. The request might or might not eventually be acted upon, as it might be disallowed when processing actually takes place. There is no facility for re-sending a status code from an asynchronous operation such as this.
	Accepted = createDataResponse(http.StatusAccepted, true)

	// NonAuthoritaviteInformation (HTTP 203)
	// The returned metainformation in the entity-header is not the definitive set as available from the origin server, but is gathered from a local or a third-party copy. The set presented MAY be a subset or superset of the original version. For example, including local annotation information about the resource might result in a superset of the metainformation known by the origin server. Use of this response code is not required and is only appropriate when the response would otherwise be 200 (OK).
	NonAuthoritaviteInformation = createDataResponse(http.StatusNonAuthoritativeInfo, true)

	// NoContent (HTTP 204)
	// The server has fulfilled the request but does not need to return an entity-body, and might want to return updated metainformation. The response MAY include new or updated metainformation in the form of entity-headers, which if present SHOULD be associated with the requested variant.
	NoContent = createDataResponse(http.StatusNoContent, false)

	// ResetContent (HTTP 205)
	// The server has fulfilled the request and the user agent SHOULD reset the document view which caused the request to be sent. This response is primarily intended to allow input for actions to take place via user input, followed by a clearing of the form in which the input is given so that the user can easily initiate another input action. The response MUST NOT include an entity.
	ResetContent = createDataResponse(http.StatusResetContent, false)

	// PartialContent (HTTP 206)
	// The server has fulfilled the partial GET request for the resource. The request MUST have included a Range header field indicating the desired range, and MAY have included an If-Range header field to make the request conditional.
	PartialContent = createDataResponse(http.StatusPartialContent, true)
)

var (
	// MultipleChoices (HTTP 300)
	// The requested resource corresponds to any one of a set of representations, each with its own specific location, and agent- driven negotiation information is being provided so that the user (or user agent) can select a preferred representation and redirect its request to that location.
	MultipleChoices = createRedirectResponse(http.StatusMultipleChoices)

	// MovedPermanently (HTTP 301)
	// The requested resource has been assigned a new permanent URI and any future references to this resource SHOULD use one of the returned URIs. Clients with link editing capabilities ought to automatically re-link references to the Request-URI to one or more of the new references returned by the server, where possible. This response is cacheable unless indicated otherwise.
	MovedPermanently = createRedirectResponse(http.StatusMovedPermanently)

	// Found (HTTP 302)
	// The requested resource resides temporarily under a different URI. Since the redirection might be altered on occasion, the client SHOULD continue to use the Request-URI for future requests. This response is only cacheable if indicated by a Cache-Control or Expires header field.
	Found = createRedirectResponse(http.StatusFound)

	// SeeOther (HTTP 303)
	// The response to the request can be found under a different URI and SHOULD be retrieved using a GET method on that resource. This method exists primarily to allow the output of a POST-activated script to redirect the user agent to a selected resource. The new URI is not a substitute reference for the originally requested resource. The 303 response MUST NOT be cached, but the response to the second (redirected) request might be cacheable.
	SeeOther = createRedirectResponse(http.StatusSeeOther)

	// NotModified (HTTP 304)
	// If the client has performed a conditional GET request and access is allowed, but the document has not been modified, the server SHOULD respond with this status code. The 304 response MUST NOT contain a message-body, and thus is always terminated by the first empty line after the header fields.
	NotModified = createRedirectResponse(http.StatusNotModified)

	// UseProxy (HTTP 305)
	// The requested resource MUST be accessed through the proxy given by the Location field. The Location field gives the URI of the proxy. The recipient is expected to repeat this single request via the proxy. 305 responses MUST only be generated by origin servers.
	UseProxy = createRedirectResponse(http.StatusUseProxy)

	// TemporaryRedirect (HTTP 307)
	// The requested resource resides temporarily under a different URI. Since the redirection MAY be altered on occasion, the client SHOULD continue to use the Request-URI for future requests. This response is only cacheable if indicated by a Cache-Control or Expires header field.
	TemporaryRedirect = createRedirectResponse(http.StatusTemporaryRedirect)
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
	MethodNotAllowed = createErrorResponse(http.StatusMethodNotAllowed)

	// NotAcceptable (HTTP 406)
	// The resource identified by the request is only capable of generating response entities which have content characteristics not acceptable according to the accept headers sent in the request.
	NotAcceptable = createErrorResponse(http.StatusNotAcceptable)

	// ProxyAuthenticationRequired (HTTP 407)
	// This code is similar to 401 (Unauthorized), but indicates that the client must first authenticate itself with the proxy. The proxy MUST return a Proxy-Authenticate header fields containing a challenge applicable to the proxy for the requested resource. The client MAY repeat the request with a suitable Proxy-Authorization header field. HTTP access authentication is explained in "HTTP Authentication: Basic and Digest Access Authentication".
	ProxyAuthenticationRequired = createErrorResponse(http.StatusProxyAuthRequired)

	// RequestTimeout (HTTP 408)
	// The client did not produce a request within the time that the server was prepared to wait. The client MAY repeat the request without modifications at any later time.
	RequestTimeout = createErrorResponse(http.StatusRequestTimeout)

	// Conflict (HTTP 409)
	// The request could not be completed due to a conflict with the current state of the resource. This code is only allowed in situations where it is expected that the user might be able to resolve the conflict and resubmit the request. The response body SHOULD include enough
	Conflict = createErrorResponse(http.StatusConflict)

	// Gone (HTTP 410)
	// The requested resource is no longer available at the server and no forwarding address is known. This condition is expected to be considered permanent. Clients with link editing capabilities SHOULD delete references to the Request-URI after user approval. If the server does not know, or has no facility to determine, whether or not the condition is permanent, the status code 404 (Not Found) SHOULD be used instead. This response is cacheable unless indicated otherwise.
	Gone = createErrorResponse(http.StatusGone)

	// LengthRequired (HTTP 411)
	// The server refuses to accept the request without a defined Content- Length. The client MAY repeat the request if it adds a valid Content-Length header field containing the length of the message-body in the request message.
	LengthRequired = createErrorResponse(http.StatusLengthRequired)

	// PreconditionFailed (HTTP 412)
	// The precondition given in one or more of the request-header fields evaluated to false when it was tested on the server. This response code allows the client to place preconditions on the current resource metainformation (header field data) and thus prevent the requested method from being applied to a resource other than the one intended.
	PreconditionFailed = createErrorResponse(http.StatusPreconditionFailed)

	// RequestEntityTooLarge (HTTP 413)
	// The server is refusing to process a request because the request entity is larger than the server is willing or able to process. The server MAY close the connection to prevent the client from continuing the request.
	RequestEntityTooLarge = createErrorResponse(http.StatusRequestEntityTooLarge)

	// RequestURITooLong (HTTP 414)
	// The server is refusing to service the request because the Request-URI is longer than the server is willing to interpret. This rare condition is only likely to occur when a client has improperly converted a POST request to a GET request with long query information, when the client has descended into a URI "black hole" of redirection (e.g., a redirected URI prefix that points to a suffix of itself), or when the server is under attack by a client attempting to exploit security holes present in some servers using fixed-length buffers for reading or manipulating the Request-URI.
	RequestURITooLong = createErrorResponse(http.StatusRequestURITooLong)

	// UnsupportedMediaType (HTTP 415)
	// The server is refusing to service the request because the entity of the request is in a format not supported by the requested resource for the requested method.
	UnsupportedMediaType = createErrorResponse(http.StatusUnsupportedMediaType)

	// RequestedRangeNotSatisfiable (HTTP 416)
	// A server SHOULD return a response with this status code if a request included a Range request-header field, and none of the range-specifier values in this field overlap the current extent of the selected resource, and the request did not include an If-Range request-header field. (For byte-ranges, this means that the first- byte-pos of all of the byte-range-spec values were greater than the current length of the selected resource.)
	RequestedRangeNotSatisfiable = createErrorResponse(http.StatusRequestedRangeNotSatisfiable)

	// ExpectationFailed (HTTP 417)
	// The expectation given in an Expect request-header field could not be met by this server, or, if the server is a proxy, the server has unambiguous evidence that the request could not be met by the next-hop server.
	ExpectationFailed = createErrorResponse(http.StatusExpectationFailed)
)

var (
	// InternalServerError (HTTP 500)
	// The server encountered an unexpected condition which prevented it from fulfilling the request.
	InternalServerError = createErrorResponse(http.StatusInternalServerError)

	// NotImplemented (HTTP 501)
	// The server does not support the functionality required to fulfill the request. This is the appropriate response when the server does not recognize the request method and is not capable of supporting it for any resource.
	NotImplemented = createErrorResponse(http.StatusNotImplemented)

	// BadGateway (HTTP 502)
	// The server, while acting as a gateway or proxy, received an invalid response from the upstream server it accessed in attempting to fulfill the request.
	BadGateway = createErrorResponse(http.StatusBadGateway)

	// ServiceUnavailable (HTTP 503)
	// The server is currently unable to handle the request due to a temporary overloading or maintenance of the server. The implication is that this is a temporary condition which will be alleviated after some delay. If known, the length of the delay MAY be indicated in a Retry-After header. If no Retry-After is given, the client SHOULD handle the response as it would for a 500 response.
	ServiceUnavailable = createErrorResponse(http.StatusServiceUnavailable)

	// GatewayTimeout (HTTP 504)
	// The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server specified by the URI (e.g. HTTP, FTP, LDAP) or some other auxiliary server (e.g. DNS) it needed to access in attempting to complete the request.
	GatewayTimeout = createErrorResponse(http.StatusGatewayTimeout)

	// HTTPVersionNotSupported (HTTP 505)
	// The server does not support, or refuses to support, the HTTP protocol version that was used in the request message. The server is indicating that it is unable or unwilling to complete the request using the same major version as the client other than with this error message. The response SHOULD contain an entity describing why that version is not supported and what other protocols are supported by that server.
	HTTPVersionNotSupported = createErrorResponse(http.StatusHTTPVersionNotSupported)
)

// Response Generic representation of data returned by HTTP resource
type Response interface {
	StatusCode() int
	GetJSON() string
	GetXML() string

	header
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

type redirectResponse struct {
	rawHeaders `json:"-" xml:"-"`
	Status     int `json:"-" xml:"-"`
}

func (o redirectResponse) StatusCode() int { return o.Status }

func (o redirectResponse) GetJSON() string { return "" }

func (o redirectResponse) GetXML() string { return "" }

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

func createErrorResponse(statusCode int) func(err error, msg ...string) Response {
	return func(err error, msg ...string) Response {
		errResponse, ok := err.(errorMessage)
		if ok {
			return errResponse
		}

		errMessage := err.Error()
		if len(msg) == 1 {
			errMessage = msg[0]
		}

		return errorMessage{
			Type:       "http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
			Title:      http.StatusText(statusCode),
			Status:     statusCode,
			Detail:     errMessage,
			rawHeaders: rawHeaders{},
		}
	}
}

func createRedirectResponse(statusCode int) func(url string) Response {
	return func(url string) Response {
		return redirectResponse{
			Status: statusCode,
			rawHeaders: rawHeaders{
				"Location": url,
			},
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
