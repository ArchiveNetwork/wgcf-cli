package utils

type RegisterOutput struct {
	Endpoint struct {
		V4 string `json:"v4"`
		V6 string `json:"v6"`
	} `json:"endpoint"`
	ReservedStr string `json:"reserved_str"`
	ReservedHex string `json:"reserved_hex"`
	ReservedDec []int  `json:"reserved_dec"`
	PrivateKey  string `json:"private_key"`
	PublicKey   string `json:"public_key"`
	Addresses   struct {
		V4 string `json:"v4"`
		V6 string `json:"v6"`
	} `json:"addresses"`
}

type Response struct {
	ID      string `json:"id"`
	Version string `json:"version,omitempty"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Name    string `json:"name,omitempty"`
	Account struct {
		ID                   string `json:"id"`
		PrivateKey           string `json:"private_key,omitempty"`
		AccountType          string `json:"account_type"`
		Created              string `json:"created,omitempty"`
		Updated              string `json:"updated,omitempty"`
		PremiumData          int    `json:"premium_data,omitempty"`
		Quota                int    `json:"quota,omitempty"`
		Usage                int    `json:"usage,omitempty"`
		WarpPlus             bool   `json:"warp_plus,omitempty"`
		ReferralCount        int    `json:"referral_count,omitempty"`
		ReferralRenewalCount int    `json:"referral_renewal_countdown,omitempty"`
		Role                 string `json:"role,omitempty"`
		License              string `json:"license,omitempty"`
		Managed              string `json:"managed,omitempty"`
		Organization         string `json:"organization,omitempty"`
	} `json:"account"`
	Policy *struct {
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
	} `json:"policy,omitempty"`
	Token     string `json:"token"`
	Warp      bool   `json:"warp_enabled,omitempty"`
	Waitlist  bool   `json:"waitlist_enabled,omitempty"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	TOS       string `json:"tos,omitempty"`
	Place     int    `json:"place,omitempty"`
	Locale    string `json:"locale"`
	Enabled   bool   `json:"enabled,omitempty"`
	InstallID string `json:"install_id"`
	FCMToken  string `json:"fcm_token"`
	SerialNum string `json:"serial_number,omitempty"`
	Config    struct {
		ClientID    string `json:"client_id"`
		ReservedHex string `json:"reserved_hex"`
		ReservedDec []int  `json:"reserved_dec"`
		PrivateKey  string `json:"private_key"`
		Peers       []struct {
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
	Model         string `json:"model,omitempty"`
	OverrideCodes *struct {
		DisableForTime struct {
			Seconds int    `json:"seconds"`
			Secret  string `json:"secret"`
		} `json:"disable_for_time"`
	} `json:"override_codes,omitempty"`
}
