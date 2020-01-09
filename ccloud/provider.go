package ccloud

import (
	"log"

	confluentcloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
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
			"confluentcloud_kafka_cluster": kafkaClusterResource(),
			"confluentcloud_api_key":       apiKeyResource(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing ConfluentCloud client")
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	c := confluentcloud.NewClient(username, password)

	return c, nil
}
