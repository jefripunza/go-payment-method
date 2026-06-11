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

type AccountHolderResponse struct {
	Id         string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Email      string `json:"email,omitempty"`
	WebsiteURL string `json:"website_url,omitempty"`
	Status     string `json:"status,omitempty"`
}

type UpdateAccountHolderRequest struct {
	WebsiteURL string `json:"website_url,omitempty"`
}

type TransferResponse struct {
	Id          string  `json:"id,omitempty"`
	Reference   string  `json:"reference,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Status      string  `json:"status,omitempty"`
	Source      string  `json:"source,omitempty"`
	Destination string  `json:"destination,omitempty"`
	Created     string  `json:"created,omitempty"`
}

// -------------------------------------------------------- //

func (x *Xendit) MerchantCreate(email string, business_name string, forUserId ...string) (*MerchantResponse, error) {
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

	resp, _, err := x.doRequest("POST", url, "", "", payload_json, forUserId...)
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

func (x *Xendit) MerchantList(forUserId ...string) ([]MerchantData, error) {
	url := fmt.Sprintf("%s/v2/accounts", x.BaseUrl)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
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

func (x *Xendit) MerchantUpdate(id string, email string, business_name string, forUserId ...string) (*MerchantResponse, error) {
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

	resp, _, err := x.doRequest("PATCH", url, "", "", payload_json, forUserId...)
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

func (x *Xendit) GetAccountHolder(id string, forUserId ...string) (*AccountHolderResponse, error) {
	url := fmt.Sprintf("%s/account_holders/%s", x.BaseUrl, id)
	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result AccountHolderResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (x *Xendit) AccountHolderUpdate(id string, req *UpdateAccountHolderRequest, forUserId ...string) (*AccountHolderResponse, error) {
	url := fmt.Sprintf("%s/account_holders/%s", x.BaseUrl, id)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, _, err := x.doRequest("PATCH", url, "", "", payload, forUserId...)
	if err != nil {
		return nil, err
	}

	var result AccountHolderResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (x *Xendit) GetTransferByReference(reference string, forUserId ...string) (*TransferResponse, error) {
	url := fmt.Sprintf("%s/transfers/reference=%s", x.BaseUrl, reference)
	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return nil, err
	}

	var result TransferResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
