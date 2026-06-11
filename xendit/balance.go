package xendit

import (
	"encoding/json"
	"fmt"
)

func (x *Xendit) GetBalance(forUserId ...string) (float64, error) {
	url := fmt.Sprintf("%s/balance", x.BaseUrl)

	resp, _, err := x.doRequest("GET", url, "", "", nil, forUserId...)
	if err != nil {
		return 0, err
	}

	var result struct {
		Balance int `json:"balance"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, err
	}

	return float64(result.Balance), nil
}

