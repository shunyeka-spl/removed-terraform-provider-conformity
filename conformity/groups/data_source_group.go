package groups

import (
	"context"
	"log"
	"strconv"
	"terraform-provider-conformity/conformity/models"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attributes_tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"attributes_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attributes_created_date": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"attributes_last_modified_date": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"relationships_organisation_data_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"relationships_organisation_data_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"relationships_accounts_data": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Something happened twice!")
	provider := m.(models.ProviderClient)
	client := provider.Client
	var diags diag.Diagnostics
	gs := GetGroupService(client)
	groups, err := gs.DoGetGroups(client.BaseUrl.String())
	log.Printf("Raw output2: %v\n", groups)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
