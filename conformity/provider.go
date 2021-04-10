package conformity

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"terraform-provider-conformity/conformity/groups"
	"terraform-provider-conformity/conformity/models"
	"terraform-provider-conformity/conformity/provider"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
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
		ResourcesMap: map[string]*schema.Resource{
			"conformity_group": groups.ResourceGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"conformity_groups": groups.DataSourceGroups(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	region := d.Get("region").(string)
	token := d.Get("auth_token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	h := make(http.Header)
	h.Set("Content-Type", "application/vnd.api+json")
	h.Set("Accept", "application/json")
	h.Set("Authorization", "ApiKey "+token)

	headers, exists := d.GetOk("headers")
	if exists {
		for k, v := range headers.(map[string]interface{}) {
			h.Set(k, v.(string))
		}
	}

	p := models.ProviderClient{
		Region:    region,
		AuthToken: token,
	}
	p.Client = provider.NewClient(h, &region, &token)

	return p, diags
}
