package ccloud

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

	confluentcloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
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
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"confluentcloud_environment": environmentDataSource(),
			"confluentcloud_service_account": serviceAccountDataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"confluentcloud_kafka_cluster":   kafkaClusterResource(),
			"confluentcloud_api_key":         apiKeyResource(),
			"confluentcloud_environment":     environmentResource(),
			"confluentcloud_schema_registry": schemaRegistryResource(),
			"confluentcloud_service_account": serviceAccountResource(),
			"confluentcloud_connector":       connectorResource(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing ConfluentCloud client")
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	wait := 2

	var diags diag.Diagnostics
	c := confluentcloud.NewClient(username, password, false)

	loginErr := c.Login()
	if loginErr == nil {
		return c, diags
	}

	err := resource.RetryContext(ctx, 30*time.Minute, func() *resource.RetryError {
		err := c.Login()

		if err != nil && strings.Contains(err.Error(), "Exceeded rate limit") {
			rand.Seed(time.Now().UnixNano())
			randomNumberMilliseconds := rand.Intn(1001)
			timeSleep := time.Duration(wait)*time.Second + time.Duration(randomNumberMilliseconds)

			log.Printf("[INFO] ConfluentCloud API rate limit exceeded, retrying in %s.", timeSleep)
			time.Sleep(timeSleep)
			wait = wait * 2
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return c, diag.FromErr(err)
}
