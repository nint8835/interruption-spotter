package spotdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const dataURL = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"

type Fetcher struct {
	lastEtag string
}

func (f *Fetcher) ShouldFetch(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, dataURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	etag := resp.Header.Get("ETag")
	return etag != f.lastEtag, nil
}

func (f *Fetcher) Fetch(ctx context.Context) (SpotDataFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dataURL, nil)
	if err != nil {
		return SpotDataFile{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SpotDataFile{}, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SpotDataFile{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data SpotDataFile
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return SpotDataFile{}, fmt.Errorf("failed to decode response: %w", err)
	}

	f.lastEtag = resp.Header.Get("ETag")
	return data, nil
}
