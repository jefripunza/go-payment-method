package tripay

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const debug = false

type Tripay struct {
	ApiKey     string
	PrivateKey string
	BaseUrl    string
}

func NewTripay(isProduction bool, apiKey string, privateKey string) *Tripay {
	baseUrl := "https://tripay.co.id/api-sandbox"
	if isProduction {
		baseUrl = "https://tripay.co.id/api"
	}
	return &Tripay{
		ApiKey:     apiKey,
		PrivateKey: privateKey,
		BaseUrl:    baseUrl,
	}
}

// doRequest performs a request to the Tripay API.
func (t *Tripay) doRequest(method, url string, body []byte) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	httpReq, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.ApiKey))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return respBody, resp.StatusCode, fmt.Errorf("tripay API error (status: %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}
