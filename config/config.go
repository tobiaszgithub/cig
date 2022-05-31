package config

type ConfigurationFile struct {
	ActiveTenantKey string
	Tenants         []struct {
		Key           string
		ApiURL        string
		Authorization struct {
			Type         string
			Username     string
			Password     string
			ClientID     string
			ClientSecret string
			TokenURL     string
		}
	}
}

type Configuration struct {
	Key           string
	ApiURL        string
	Authorization struct {
		Type         string
		Username     string
		Password     string
		ClientID     string
		ClientSecret string
		TokenURL     string
	}
}
