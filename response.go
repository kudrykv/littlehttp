package littlehttp

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Response struct {
	defaultUnmarshaller Unmarshaller

	raw *http.Response

	once  *sync.Once
	mutex *sync.Mutex

	bodyBytes []byte
}

func newResponse(defaultUnmarshaller Unmarshaller) *Response {
	return &Response{
		once:  &sync.Once{},
		mutex: &sync.Mutex{},

		defaultUnmarshaller: defaultUnmarshaller,
	}
}

func (r *Response) Unmarshal(dst any) error {
	bodyBytes, err := r.readBodyOnce()
	if err != nil {
		return fmt.Errorf("read body once: %w", err)
	}

	if err = r.unmarshal(r.raw.Header.Get("Content-Type"), bodyBytes, dst); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

func (r *Response) readBodyOnce() ([]byte, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var err error

	r.once.Do(func() { r.bodyBytes, err = r.readBody() })

	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return r.bodyBytes, nil
}

func (r *Response) readBody() ([]byte, error) {
	bts, err := io.ReadAll(r.raw.Body)
	if err != nil {
		return nil, fmt.Errorf("read all body: %w", err)
	}

	if err = r.raw.Body.Close(); err != nil {
		return nil, fmt.Errorf("close body: %w", err)
	}

	return bts, nil
}

func (r *Response) unmarshal(contentType string, bts []byte, dst any) error {
	unmarshal, err := r.pickUnmarshaller(contentType)
	if err != nil {
		return fmt.Errorf("pick unmarshaller for %s: %w", contentType, err)
	}

	if err = unmarshal(bts, dst); err != nil {
		return fmt.Errorf("unmarshal content-type %s: %w", contentType, err)
	}

	return nil
}

func (r *Response) pickUnmarshaller(contentType string) (Unmarshaller, error) {
	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return json.Unmarshal, nil
	case strings.HasPrefix(contentType, "application/xml"):
		return xml.Unmarshal, nil
	default:
		if r.defaultUnmarshaller == nil {
			return nil, errors.New("default unmarshaller not set")
		}

		return r.defaultUnmarshaller, nil
	}
}

func (r *Response) IsSuccessful() bool {
	return r.raw.StatusCode >= 200 && r.raw.StatusCode < 300
}
