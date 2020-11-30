package ccloud

import (
	"log"
	"strings"
	"time"

	confluentcloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	log.Printf("[INFO] Creating Provider")
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONFLUENT_CLOUD_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CONFLUENT_CLOUD_PASSWORD", ""),
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"confluentcloud_kafka_cluster":   kafkaClusterResource(),
			"confluentcloud_api_key":         apiKeyResource(),
			"confluentcloud_environment":     environmentResource(),
			"confluentcloud_schema_registry": schemaRegistryResource(),
			"confluentcloud_service_account": serviceAccountResource(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ConfluentCloud client")
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	c := confluentcloud.NewClient(username, password)

	loginE := c.Login()

	if loginE == nil {
		return c, loginE
	}

	return c, resource.Retry(30*time.Minute, func() *resource.RetryError {
		err := c.Login()

		if strings.Contains(err.Error(), "Exceeded rate limit") {
			log.Printf("[INFO] ConfluentCloud API rate limit exceeded, retrying.")
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}
