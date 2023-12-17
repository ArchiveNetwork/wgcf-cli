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
	var body []byte
	var err error
	var output []byte
	var response LicenseResponse
	payload := []byte(
		`{
			"license":"` + license + `"
		 }`,
	)

	if body, err = request(payload, token, id, "license"); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	if output, err = json.MarshalIndent(response, "", "    "); err != nil {
		panic(err)
	}
	return string(output), nil
}
