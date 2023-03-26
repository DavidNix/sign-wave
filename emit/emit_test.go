package emit

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockHTTPPoster struct {
	GotURL         string
	GotContentType string
	GotBody        string

	StubResponse *http.Response
	StubErr      error
}

func (m *mockHTTPPoster) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	m.GotURL = url
	m.GotContentType = contentType
	b, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	m.GotBody = string(b)

	if m.StubResponse == nil {
		m.StubResponse = &http.Response{StatusCode: 200, Body: io.NopCloser(nil)}
	}

	return m.StubResponse, m.StubErr
}

func TestSendToIngest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var client mockHTTPPoster

		err := SendToIngest(&client, "http://example.com", []int64{1, 2, 3})
		require.NoError(t, err)

		require.Equal(t, "http://example.com/v1/batch", client.GotURL)
		require.Equal(t, "application/json", client.GotContentType)
		require.JSONEq(t, `{"record_ids":[1,2,3]}`, client.GotBody)
	})

	t.Run("post error", func(t *testing.T) {
		var client mockHTTPPoster
		client.StubErr = errors.New("boom")

		err := SendToIngest(&client, "http://example.com", []int64{1, 2, 3})
		require.Error(t, err)
		require.EqualError(t, err, "boom")
	})

	t.Run("bad status code", func(t *testing.T) {
		var client mockHTTPPoster
		client.StubResponse = &http.Response{StatusCode: 500, Body: io.NopCloser(nil)}

		err := SendToIngest(&client, "http://example.com", []int64{1, 2, 3})
		require.Error(t, err)
		require.EqualError(t, err, "bad status code: 500")
	})
}
