package ccloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaRegistryResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: schemaRegistryCreate,
		ReadContext:   schemaRegistryRead,
		DeleteContext: schemaRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schemaRegistryImport,
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

func schemaRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	environmentID := d.Get("environment_id").(string)
	region := d.Get("region").(string)
	serviceProvider := d.Get("service_provider").(string)

	log.Printf("[INFO] Creating Schema Registry %s", environmentID)

	reg, err := c.CreateSchemaRegistry(environmentID, region, serviceProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created Schema Registry %s", reg.ID)

	d.SetId(reg.ID)

	err = d.Set("endpoint", reg.Endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func schemaRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	log.Printf("[INFO] Reading Schema Registry %s", d.Id())
	environmentID := d.Get("environment_id").(string)

	err := d.Set("environment_id", environmentID)
	if err != nil {
		return diag.FromErr(err)
	}

	registry, err := c.GetSchemaRegistry(environmentID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if registry == nil {
		return diag.Errorf("Unable to read schema registry")
	}

	err = d.Set("endpoint", registry.Endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

func schemaRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[WARN] Schema registry cannot be deleted: %s", d.Id())
	return nil
}

func schemaRegistryImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	envIDAndClusterID := d.Id()
	parts := strings.Split(envIDAndClusterID, "/")

	var err error
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for schema registry cluster import: expected '<env ID>/<cluster ID>'")
	}

	d.SetId(parts[1])
	err = d.Set("environment_id", parts[0])
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
