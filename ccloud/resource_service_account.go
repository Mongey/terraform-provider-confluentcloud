package ccloud

import (
	"fmt"
	"log"
	"strconv"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func serviceAccountResource() *schema.Resource {
	return &schema.Resource{
		Create: serviceAccountCreate,
		Read:   serviceAccountRead,
		Delete: serviceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
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

func serviceAccountCreate(d *schema.ResourceData, meta interface{}) error {
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
			return err
		}

		err = d.Set("description", serviceAccount.Description)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[ERROR] Could not create Service Account: %s", err)
	}

	return err
}

func serviceAccountRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func serviceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	ID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not parse Service Account ID %s to int", d.Id())
		return err
	}

	err = c.DeleteServiceAccount(ID)
	if err != nil {
		log.Printf("[ERROR] Service Account can not be deleted: %d", ID)
		return err
	}

	log.Printf("[INFO] Service Account deleted: %d", ID)

	return nil
}
