package ccloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func serviceAccountResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: serviceAccountCreate,
		ReadContext:   serviceAccountRead,
		DeleteContext: serviceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Account Name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Account Description",
			},
		},
	}
}

func serviceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	req := ccloud.ServiceAccountCreateRequest{
		Name:        name,
		Description: description,
	}

	serviceAccount, err := c.CreateServiceAccount(&req)
	if err == nil {
		d.SetId(fmt.Sprintf("%d", serviceAccount.ID))

		err = d.Set("name", serviceAccount.Name)
		if err != nil {
			return diag.FromErr(err)
		}

		err = d.Set("description", serviceAccount.Description)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		log.Printf("[ERROR] Could not create Service Account: %s", err)
	}

	return diag.FromErr(err)
}

func serviceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	ID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not parse Service Account ID %s to int", d.Id())
		return diag.FromErr(err)
	}
	account, err := getServiceAccount(c, ID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", account.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", account.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getServiceAccount(client *ccloud.Client, id int) (*ccloud.ServiceAccount, error) {
	accounts, err := client.ListServiceAccounts()
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.ID == id {
			return &account, nil
		}
	}

	return nil, errors.New("Unable to find service account")
}

func serviceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	ID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not parse Service Account ID %s to int", d.Id())
		return diag.FromErr(err)
	}

	err = c.DeleteServiceAccount(ID)
	if err != nil {
		log.Printf("[ERROR] Service Account can not be deleted: %d", ID)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Service Account deleted: %d", ID)

	return nil
}
