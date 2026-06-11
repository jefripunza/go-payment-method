package xendit

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type Xendit struct {
	Authorization string
	CallbackToken string
	ForUserId     string
	BaseUrl       string
}

func NewXendit(apiKey string, callbackToken string, defaultForUserId ...string) *Xendit {
	auth := apiKey + ":"
	authorization := base64.StdEncoding.EncodeToString([]byte(auth))
	var fUserId string
	if len(defaultForUserId) > 0 {
		fUserId = defaultForUserId[0]
	}
	return &Xendit{
		Authorization: authorization,
		CallbackToken: callbackToken,
		ForUserId:     fUserId,
		BaseUrl:       "https://api.xendit.co",
	}
}

func (x *Xendit) doRequest(method, url string, apiVersion string, businessId string, body []byte, forUserId ...string) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	httpReq, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Basic %s", x.Authorization))
	fUserId := x.ForUserId
	if len(forUserId) > 0 && forUserId[0] != "" {
		fUserId = forUserId[0]
	}
	if fUserId != "" {
		httpReq.Header.Set("for-user-id", fUserId)
	}
	if apiVersion != "" {
		httpReq.Header.Set("api-version", apiVersion)
	}
	if businessId != "" {
		httpReq.Header.Set("business-id", businessId)
	}

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
		return respBody, resp.StatusCode, fmt.Errorf("xendit API error (status: %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}
