package littlehttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Request struct {
	Method  string
	URL     string
	Headers http.Header
	Body    any
}

func NewRequest(method, url string, headers http.Header, body any) Request {
	return Request{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}
}

func (r *Request) prepare(
	ctx context.Context, marshaller Marshaller, mandatoryHeaders http.Header, prefix string,
) (*http.Request, error) {
	contentType, bodyReader, err := r.prepareBody(marshaller)
	if err != nil {
		return nil, fmt.Errorf("prepare body: %w", err)
	}

	request, err := r.prepareRequestObject(ctx, r.Method, prefix+r.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	request.Header = mergeSourceHeadersIntoDestination(r.Headers, mandatoryHeaders)

	if len(contentType) > 0 {
		request.Header.Set("Content-Type", contentType)
	}

	return request, nil
}

func (r Request) prepareBody(marshaller Marshaller) (string, *bytes.Reader, error) {
	if r.Body == nil {
		return "", nil, nil
	}

	if marshaller == nil {
		return "", nil, errors.New("marshaller not set")
	}

	contentType, bts, err := marshaller(r.Body)
	if err != nil {
		return "", nil, fmt.Errorf("marshaller: %w", err)
	}

	return contentType, bytes.NewReader(bts), nil
}

func (r *Request) prepareRequestObject(ctx context.Context, method, url string, body *bytes.Reader) (*http.Request, error) {
	if body == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	return http.NewRequestWithContext(ctx, method, url, body)
}

func mergeSourceHeadersIntoDestination(destination, source http.Header) http.Header {
	if len(source) == 0 {
		return destination
	}

	if destination == nil {
		destination = make(http.Header, len(source))
	}

	for key, values := range source {
		for _, value := range values {
			source.Add(key, value)
		}
	}

	return destination
}
