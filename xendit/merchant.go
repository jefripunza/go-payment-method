package xendit

import (
	"encoding/json"
	"fmt"
)

// this is xenPlatform

type MerchantResponse struct {
	Id            string            `json:"id"`
	Created       string            `json:"created"`
	Updated       string            `json:"updated"`
	Email         string            `json:"email"`
	Type          string            `json:"type"`
	PublicProfile map[string]string `json:"public_profile"`
	Country       string            `json:"country"`
	Status        string            `json:"status"`
}

type MerchantCreateRequest struct {
	Email         string            `json:"email"`
	Type          string            `json:"type"`
	PublicProfile map[string]string `json:"public_profile"`
}

type MerchantData struct {
	Id            string            `json:"id"`
	Created       string            `json:"created"`
	Updated       string            `json:"updated"`
	Email         string            `json:"email"`
	Type          string            `json:"type"`
	PublicProfile map[string]string `json:"public_profile"`
	Country       string            `json:"country"`
	Status        string            `json:"status"`
}

type MerchantListResponse struct {
	Data    []MerchantData      `json:"data"`
	HasMore bool                `json:"has_more"`
	Links   []MerchantLinksData `json:"links"`
}

type MerchantLinksData struct {
	Href   string `json:"href"`
	Method string `json:"method"`
	Rel    string `json:"rel"`
}

type MerchantUpdateRequest struct {
	Email         string            `json:"email"`
	PublicProfile map[string]string `json:"public_profile"`
}

// -------------------------------------------------------- //

func (x *Xendit) MerchantCreate(email string, business_name string) (*MerchantResponse, error) {
	url := fmt.Sprintf("%s/v2/accounts", x.BaseUrl)

	payload := MerchantCreateRequest{
		Email: email,
		Type:  "OWNED",
		PublicProfile: map[string]string{
			"business_name": business_name,
		},
	}

	// Convert request to JSON
	payload_json, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, _, err := x.doRequest("POST", url, "", "", payload_json)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result MerchantResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func (x *Xendit) MerchantList() ([]MerchantData, error) {
	url := fmt.Sprintf("%s/v2/accounts", x.BaseUrl)

	resp, _, err := x.doRequest("GET", url, "", "", nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result MerchantListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Data, nil
}

func (x *Xendit) MerchantUpdate(id string, email string, business_name string) (*MerchantResponse, error) {
	url := fmt.Sprintf("%s/v2/accounts/%s", x.BaseUrl, id)

	payload := MerchantUpdateRequest{
		Email: email,
		PublicProfile: map[string]string{
			"business_name": business_name,
		},
	}

	// Convert request to JSON
	payload_json, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, _, err := x.doRequest("PATCH", url, "", "", payload_json)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result MerchantResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
