package response

type rawHeaders map[string]string

type Header interface {
	WithHeader(key, value string)
	Header() rawHeaders
}

func (h rawHeaders) WithHeader(key, value string) {
	h[key] = value
}

func (h rawHeaders) Header() rawHeaders {
	return h
}
