package ccloud

import (
	"log"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaRegistryResource() *schema.Resource {
	return &schema.Resource{
		Create: schemaRegistryCreate,
		Read:   schemaRegistryRead,
		Delete: schemaRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment ID",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "where",
			},
			"service_provider": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud provider",
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	environment := d.Get("environment_id").(string)
	region := d.Get("region").(string)
	service_provider := d.Get("service_provider").(string)

	log.Printf("[INFO] Creating Schema Registry %s", environment)

	reg, err := c.CreateSchemaRegistry(environment, region, service_provider)
	if err != nil {
		return err
	}

	d.SetId(reg.ID)
	err = d.Set("endpoint", reg.Endpoint)
	if err != nil {
		return err
	}

	return nil
}

func schemaRegistryRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	environment := d.Get("environment_id").(string)
	log.Printf("[INFO] Reading Schema Registry %s", environment)

	env, err := c.GetSchemaRegistry(environment)
	if err != nil {
		return err
	}

	err = d.Set("environment_id", environment)
	if err != nil {
		err = d.Set("endpoint", env.Endpoint)
	}

	return err
}

func schemaRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Schema registry cannot be deleted: %s", d.Id())
	return nil
}
