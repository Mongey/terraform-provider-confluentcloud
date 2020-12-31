package ccloud

import (
	"context"
	"log"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func environmentResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: environmentCreate,
		ReadContext:   environmentRead,
		UpdateContext: environmentUpdate,
		DeleteContext: environmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The name of the environment",
			},
		},
	}
}

func environmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	name := d.Get("name").(string)

	log.Printf("[INFO] Creating Environment %s", name)
	orgID, err := getOrganizationID(c)
	if err != nil {
		return diag.FromErr(err)
	}

	env, err := c.CreateEnvironment(name, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(env.ID)

	return nil
}

func environmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	newName := d.Get("name").(string)

	log.Printf("[INFO] Updating Environment %s", d.Id())
	orgID, err := getOrganizationID(c)
	if err != nil {
		return diag.FromErr(err)
	}

	env, err := c.UpdateEnvironment(d.Id(), newName, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(env.ID)

	return nil
}

func environmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	log.Printf("[INFO] Reading Environment %s", d.Id())
	env, err := c.GetEnvironment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", env.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func environmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	log.Printf("[INFO] Deleting Environment %s", d.Id())
	err := c.DeleteEnvironment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getOrganizationID(client *ccloud.Client) (int, error) {
	userData, err := client.Me()
	if err != nil {
		return 0, err
	}

	return userData.Account.OrganizationID, nil
}
