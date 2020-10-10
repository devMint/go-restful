package response

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validResponse = []string{"test", "test"}
	validJSON     = "{\"data\":[\"test\",\"test\"]}"
	validXML      = "<response><data>test</data><data>test</data></response>"

	emptyJSON = "{\"data\":null}"
	emptyXML  = "<response></response>"

	errorsMsg        = errors.New("missing entity")
	errorsJSON       = "{\"type\":\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html\",\"title\":\"Not Found\",\"detail\":\"missing entity\",\"status\":404}"
	errorsJSONCustom = "{\"type\":\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html\",\"title\":\"Not Found\",\"detail\":\"not-found-test\",\"status\":404}"
	errorsXML        = "<response><type>http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html</type><title>Not Found</title><detail>missing entity</detail><status>404</status></response>"
)

func Test_EncodeResponse_ToJSON(t *testing.T) {
	assert.Equal(t, validJSON, Ok(validResponse).GetJSON())
}

func Test_EncodeResponse_ToXML(t *testing.T) {
	assert.Equal(t, validXML, Ok(validResponse).GetXML())
}

func Test_EncodeResponse_EmptyContent_ToJSON(t *testing.T) {
	assert.Equal(t, emptyJSON, Ok().GetJSON())
}

func Test_EncodeResponse_EmptyContent_ToXML(t *testing.T) {
	assert.Equal(t, emptyXML, Ok().GetXML())
}

func Test_EncodeError_ToJSON(t *testing.T) {
	assert.Equal(t, errorsJSON, NotFound(errorsMsg).GetJSON())
}

func Test_EncodeError_ToXML(t *testing.T) {
	assert.Equal(t, errorsXML, NotFound(errorsMsg).GetXML())
}

func Test_EncodeError_ToJSON_CustomMessage(t *testing.T) {
	assert.Equal(t, errorsJSONCustom, NotFound(errorsMsg, "not-found-test").GetJSON())
}

func Test_CustomHeaders(t *testing.T) {
	r := NotFound(errorsMsg)
	r.WithHeader("a", "b")

	value, ok := r.Header()["a"]
	assert.True(t, ok)
	assert.Equal(t, "b", value)
}
