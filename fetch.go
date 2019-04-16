package suggest

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Fetcher
type Fetcher struct {
	client  *http.Client
	headers map[string]string
	timeout time.Duration
}

// NewFetcher
func NewFetcher(headers map[string]string, timeout time.Duration) *Fetcher {
	f := &Fetcher{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		headers: headers,
		timeout: timeout,
	}

	return f
}

// Do send request
func (f *Fetcher) Do(rawurl, query, lang string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range f.headers {
		req.Header.Set(k, v)
	}

	// https://www.google.com/complete/search?hl=en&output=toolbar&q=word
	values := url.Values{}
	values.Set("output", "toolbar")
	values.Set("q", query)
	if len(lang) > 0 {
		values.Set("hl", lang)
	}

	req.URL.RawQuery = values.Encode()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req = req.WithContext(ctx)

	return f.client.Do(req)
}

// Fetch return XML as a structure
func (f *Fetcher) Fetch(rawurl, query, language string) (GoogleSuggestion, error) {
	var gs GoogleSuggestion
	resp, err := f.Do(rawurl, query, language)
	if err != nil {
		return gs, err
	}
	defer func() {
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return gs, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	// xml to struct
	return XMLDecode(resp.Body)
}
