package models

import "terraform-provider-conformity/conformity/provider"

type ProviderClient struct {
	Region    string
	AuthToken string
	Client    *provider.Client
}
