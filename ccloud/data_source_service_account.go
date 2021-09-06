package ccloud

import (
	"context"
	"log"
	"strconv"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func serviceAccountDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: serviceAccountDataSourceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The name of the service account",
			},
		},
	}
}

func serviceAccountDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Service Account %s", name)
	serviceAccounts, err := c.ListServiceAccounts()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, serviceAccount := range serviceAccounts {
		if serviceAccount.Name == name {
			d.SetId(strconv.Itoa(serviceAccount.ID))
			err := d.Set("name", serviceAccount.Name)

			if err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	return nil
}
