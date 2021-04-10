package groups

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-conformity/conformity/models"
	"terraform-provider-conformity/conformity/utils"
	"time"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Println("[DEBUG] Starting group creation")
	provider := m.(models.ProviderClient)
	client := provider.Client
	var diags diag.Diagnostics

	group := models.Group{}
	group.Attributes.Name = d.Get("name").(string)
	group.Attributes.Tags = utils.ExpandStringList(d.Get("tags").([]interface{}))

	gs := GetGroupService(client)
	ng, err := gs.DoCreateGroup(client.BaseUrl.String(), &group)
	log.Printf("Group Post output2: %v\n", ng)
	if err != nil {
		log.Printf("[ERROR] Group Post error: %v\n", err)
		return diag.FromErr(err)
	}

	d.SetId(ng.ID)

	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	groupID := d.Id()
	log.Printf("[DEBUG] Starting group read for group id %s\n", groupID)
	provider := m.(models.ProviderClient)
	client := provider.Client
	var diags diag.Diagnostics

	gs := GetGroupService(client)
	group, err := gs.DoGetGroup(client.BaseUrl.String(), groupID)
	if err != nil {
		log.Printf("Group Read error: %v\n", err)
		return diag.FromErr(err)
	}
	if err := d.Set("name", group.Attributes.Name); err != nil {
		log.Printf("Group Read Name not found: %v\n", err)
		return diag.FromErr(err)
	}
	if err := d.Set("tags", group.Attributes.Tags); err != nil {
		log.Printf("Group Read tags not found: %v\n", err)
		return diag.FromErr(err)
	}
	log.Printf("Group Read output: %v\n", group)

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Starting group update")
	provider := m.(models.ProviderClient)
	client := provider.Client
	groupID := d.Id()
	if d.HasChanges("name", "tags") {
		group := models.Group{}
		group.ID = groupID
		group.Attributes.Name = d.Get("name").(string)
		group.Attributes.Tags = utils.ExpandStringList(d.Get("tags").([]interface{}))

		gs := GetGroupService(client)
		ug, err := gs.DoUpdateGroup(client.BaseUrl.String(), &group)
		log.Printf("Group Patch output2: %v\n", ug)
		if err != nil {
			log.Printf("Group Post error: %v\n", err)
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	groupID := d.Id()
	log.Printf("[DEBUG] Starting group delete for group id %s\n", groupID)
	provider := m.(models.ProviderClient)
	client := provider.Client
	var diags diag.Diagnostics

	gs := GetGroupService(client)
	err := gs.DoDeleteGroup(client.BaseUrl.String(), groupID)
	if err != nil {
		log.Printf("[DEBUG] Group delete error: %v\n", err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
