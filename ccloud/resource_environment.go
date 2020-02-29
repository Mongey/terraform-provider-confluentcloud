package ccloud

import (
	"log"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func environmentResource() *schema.Resource {
	return &schema.Resource{
		Create: environmentCreate,
		Read:   environmentRead,
		Update: environmentUpdate,
		Delete: environmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func environmentCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	name := d.Get("name").(string)

	log.Printf("[INFO] Creating Environment %s", name)
	orgID, err := getOrganizationID(c)
	if err != nil {
		return err
	}

	env, err := c.CreateEnvironment(name, orgID)
	if err != nil {
		return err
	}

	d.SetId(env.ID)

	return nil
}

func environmentUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	newName := d.Get("name").(string)

	log.Printf("[INFO] Updating Environment %s", d.Id())
	orgID, err := getOrganizationID(c)
	if err != nil {
		return err
	}

	env, err := c.UpdateEnvironment(d.Id(), newName, orgID)
	if err != nil {
		return err
	}

	d.SetId(env.ID)

	return nil
}

func environmentRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	log.Printf("[INFO] Reading Environment %s", d.Id())
	env, err := c.GetEnvironment(d.Id())
	if err != nil {
		return err
	}

	err = d.Set("name", env.Name)
	if err != nil {
		return err
	}

	return nil
}

func environmentDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	log.Printf("[INFO] Deleting Environment %s", d.Id())
	err := c.DeleteEnvironment(d.Id())
	if err != nil {
		return err
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
