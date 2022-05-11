package littlehttp

import (
	"context"
	"fmt"
	"net/http"
)

type LittleHTTP struct {
	client              *http.Client
	defaultUnmarshaller Unmarshaller
	marshaller          Marshaller

	mandatoryHeaders http.Header
	urlPrefix        string
}

func New(params Parameters) (*LittleHTTP, error) {
	if err := params.SetDefaultsAndValidate(); err != nil {
		return nil, fmt.Errorf("set defaults and validate: %w", err)
	}

	return &LittleHTTP{
		client:              params.Client,
		defaultUnmarshaller: params.DefaultUnmarshaller,
		marshaller:          params.Marshaller,
		urlPrefix:           params.URLPrefix,
		mandatoryHeaders:    params.MandatoryHeaders,
	}, nil
}

func (r *LittleHTTP) SetMandatoryHeaders(headers http.Header) *LittleHTTP {
	r.mandatoryHeaders = make(http.Header, len(headers))

	for key, values := range headers {
		cp := make([]string, len(values))
		copy(cp, values)

		r.mandatoryHeaders[key] = values
	}

	return r
}

func (r *LittleHTTP) Get(ctx context.Context, url string, headers http.Header) (*Response, error) {
	response, err := r.Do(ctx, Request{Method: http.MethodGet, URL: url, Headers: headers})
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	return response, nil
}

func (r *LittleHTTP) Do(ctx context.Context, req Request) (*Response, error) {
	request, err := req.prepare(ctx, r.marshaller, r.mandatoryHeaders, r.urlPrefix)
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	response := newResponse(r.defaultUnmarshaller)
	if response.raw, err = r.client.Do(request); err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	return response, nil
}
