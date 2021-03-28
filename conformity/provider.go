package conformity

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ap-southeast-2", nil),
			},
			"auth_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"conformity_groups": dataSourceGroups(),
		},
	}
}

func marshalData(d *schema.ResourceData, vals map[string]interface{}) {
	for k, v := range vals {
		if k == "id" {
			d.SetId(v.(string))
		} else {
			str, ok := v.(string)
			if ok {
				d.Set(k, str)
			} else {
				d.Set(k, v)
			}
		}
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[DEBUG] Something happened!")
	region := d.Get("region").(string)
	if region == "" {
		log.Println("Defaulting environment in URL config to use default region ap-southeast-1")
	}

	authToken := d.Get("auth_token").(string)

	h := make(http.Header)
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	h.Set("Authorization", authToken)

	headers, exists := d.GetOk("headers")
	if exists {
		for k, v := range headers.(map[string]interface{}) {
			h.Set(k, v.(string))
		}
	}

	return newProviderClient(region, authToken, h)
}

type ProviderClient struct {
	Region    string
	AuthToken string
	Client    *Client
}

func newProviderClient(region, authToken string, headers http.Header) (ProviderClient, error) {
	p := ProviderClient{
		Region:    region,
		AuthToken: authToken,
	}
	p.Client = NewClient(headers, region, authToken)

	return p, nil
}
