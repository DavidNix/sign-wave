package emit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	RecordIDs []int64 `json:"record_ids"`
}

type HTTPPoster interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// SendToIngest sends a batch of records to the ingest endpoint.
func SendToIngest(client HTTPPoster, baseURL string, recordIDs []int64) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("bad base url: %w", err)
	}
	u.Path = "/v1/batch"

	body := Request{RecordIDs: recordIDs}
	bz, err := json.Marshal(body)
	if err != nil {
		panic(err) // Should not be possible.
	}
	resp, err := client.Post(u.String(), "application/json", bytes.NewReader(bz))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}
