package littlehttp

import (
	"encoding/json"
	"net/http"
)

type Unmarshaller func(data []byte, dst any) error
type Marshaller func(src any) ([]byte, error)

type Parameters struct {
	Client              *http.Client
	DefaultUnmarshaller Unmarshaller
	Marshaller          Marshaller

	URLPrefix string
}

func (p *Parameters) SetDefaultsAndValidate() error {
	p.maybeSetDefaultClient()
	p.maybeSetUnmarshaller()

	return nil
}

func (p *Parameters) maybeSetDefaultClient() {
	if p.Client == nil {
		p.Client = http.DefaultClient
	}
}

func (p *Parameters) maybeSetUnmarshaller() {
	if p.DefaultUnmarshaller == nil {
		p.DefaultUnmarshaller = json.Unmarshal
	}
}
