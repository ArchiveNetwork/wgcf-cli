package main

import (
	"encoding/json"
	"time"
)

func changeLicense(token string, id string, license string) (string, error) {
	type LicenseResponse struct {
		ID                       string    `json:"id"`
		Created                  time.Time `json:"created"`
		Updated                  time.Time `json:"updated"`
		PremiumData              int       `json:"premium_data"`
		Quota                    int       `json:"quota"`
		WarpPlus                 bool      `json:"warp_plus"`
		ReferralCount            int       `json:"referral_count"`
		ReferralRenewalCountdown int       `json:"referral_renewal_countdown"`
		Role                     string    `json:"role"`
	}
	var response LicenseResponse
	payload := []byte(
		`{
			"license":"` + license + `"
		 }`,
	)
	body, err := request(payload, token, id, "license")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	output, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(output), nil
}
