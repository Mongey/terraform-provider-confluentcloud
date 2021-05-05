package ccloud

import (
	"context"
	"log"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func environmentDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: environmentDataSourceRead,
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

func environmentDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
