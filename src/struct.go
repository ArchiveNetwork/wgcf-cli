package main

type Response interface {
	isResponse()
}

type NormalResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Model   string `json:"model"`
	Name    string `json:"name"`
	Key     string `json:"key"`
	Account struct {
		ID                   string `json:"id"`
		PrivateKey           string `json:"private_key"`
		ReservedHex          string `json:"reserved_hex"`
		ReservedDec          []int  `json:"reserved_dec"`
		AccountType          string `json:"account_type"`
		Created              string `json:"created"`
		Updated              string `json:"updated"`
		PremiumData          int    `json:"premium_data"`
		Quota                int    `json:"quota"`
		Usage                int    `json:"usage"`
		WarpPlus             bool   `json:"warp_plus"`
		ReferralCount        int    `json:"referral_count"`
		ReferralRenewalCount int    `json:"referral_renewal_countdown"`
		Role                 string `json:"role"`
		License              string `json:"license"`
	} `json:"account"`
	Config struct {
		ClientID string `json:"client_id"`
		Peers    []struct {
			PublicKey string `json:"public_key"`
			Endpoint  struct {
				V4   string `json:"v4"`
				V6   string `json:"v6"`
				Host string `json:"host"`
			} `json:"endpoint"`
		} `json:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `json:"v4"`
				V6 string `json:"v6"`
			} `json:"addresses"`
		} `json:"interface"`
		Services struct {
			HTTPProxy string `json:"http_proxy"`
		} `json:"services"`
	} `json:"config"`
	Token     string `json:"token"`
	Warp      bool   `json:"warp_enabled"`
	Waitlist  bool   `json:"waitlist_enabled"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	TOS       string `json:"tos"`
	Place     int    `json:"place"`
	Locale    string `json:"locale"`
	Enabled   bool   `json:"enabled"`
	InstallID string `json:"install_id"`
	FCMToken  string `json:"fcm_token"`
	SerialNum string `json:"serial_number"`
}

func (r NormalResponse) isResponse() {}

type TeamResponse struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Account struct {
		ID           string `json:"id"`
		PrivateKey   string `json:"private_key"`
		ReservedHex  string `json:"reserved_hex"`
		ReservedDec  []int  `json:"reserved_dec"`
		AccountType  string `json:"account_type"`
		Managed      string `json:"managed"`
		Organization string `json:"organization"`
	} `json:"account"`
	Policy struct {
		ServiceModeV2 struct {
			Mode string `json:"mode"`
		} `json:"service_mode_v2"`
		DisableAutoFallback bool `json:"disable_auto_fallback"`
		FallbackDomains     []struct {
			Suffix string `json:"suffix"`
		} `json:"fallback_domains"`
		Exclude []struct {
			Address     string `json:"address"`
			Description string `json:"description,omitempty"`
		} `json:"exclude"`
		GatewayUniqueID  string `json:"gateway_unique_id"`
		AppURL           string `json:"app_url"`
		Organization     string `json:"organization"`
		CaptivePortal    int    `json:"captive_portal"`
		AllowModeSwitch  bool   `json:"allow_mode_switch"`
		AllowedToLeave   bool   `json:"allowed_to_leave"`
		ExcludeOfficeIPs bool   `json:"exclude_office_ips"`
	} `json:"policy"`
	Token    string `json:"token"`
	Locale   string `json:"locale"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	FcmToken string `json:"fcm_token"`
	Config   struct {
		ClientID string `json:"client_id"`
		Peers    []struct {
			PublicKey string `json:"public_key"`
			Endpoint  struct {
				V4   string `json:"v4"`
				V6   string `json:"v6"`
				Host string `json:"host"`
			} `json:"endpoint"`
		} `json:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `json:"v4"`
				V6 string `json:"v6"`
			} `json:"addresses"`
		} `json:"interface"`
		Services struct {
			HTTPProxy string `json:"http_proxy"`
		} `json:"services"`
	} `json:"config"`
	InstallID     string `json:"install_id"`
	Model         string `json:"model"`
	OverrideCodes struct {
		DisableForTime struct {
			Seconds int    `json:"seconds"`
			Secret  string `json:"secret"`
		} `json:"disable_for_time"`
	} `json:"override_codes"`
}

func (r TeamResponse) isResponse() {}

type ResponseHolder *NormalResponse
