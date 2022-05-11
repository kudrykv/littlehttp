package littlehttp_test

import (
	"context"
	"littlehttp"
	"net/http"
	"testing"
)

type roundTrip struct {
	statusCode int
}

func (r roundTrip) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: r.statusCode,
	}, nil
}

func TestLittleHTTP_Get(t *testing.T) {
	client, err := littlehttp.New(littlehttp.Parameters{
		Client:    &http.Client{Transport: &roundTrip{statusCode: 200}},
		URLPrefix: "http://localhost:8080",
	})
	expectNoError(t, err)

	response, err := client.Do(context.Background(), littlehttp.Request{})
	expectNoError(t, err)

	expectEquals(t, response.IsSuccessful(), true)
}

func expectNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func expectEquals(t *testing.T, expected, actual any) {
	t.Helper()

	if expected != actual {
		t.Fatalf("expected != actual: %v != %v", expected, actual)
	}
}
