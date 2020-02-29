package ccloud

import (
	"fmt"
	"log"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func kafkaClusterResource() *schema.Resource {
	return &schema.Resource{
		Create: clusterCreate,
		Read:   clusterRead,
		//Update: clusterUpdate,
		Delete: clusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the cluster",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment ID",
			},
			"bootstrap_servers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_provider": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "AWS / GCP",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "where",
			},
			"availability": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "LOW(single-zone) or HIGH(multi-zone)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if val != "LOW" && val != "HIGH" {
						errs = append(errs, fmt.Errorf("%q must be `LOW` or `HIGH`, got: %s", key, v))
					}
					return
				},
			},
		},
	}
}

func clusterCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	name := d.Get("name").(string)
	region := d.Get("region").(string)
	serviceProvider := d.Get("service_provider").(string)
	durability := d.Get("availability").(string)
	accountID := d.Get("environment_id").(string)

	log.Printf("[DEBUG] Creating kafka_cluster")

	req := ccloud.ClusterCreateConfig{
		Name:            name,
		Region:          region,
		ServiceProvider: serviceProvider,
		Storage:         5000, // TODO: paramaterize
		AccountID:       accountID,
		Durability:      durability,
	}
	cluster, err := c.CreateCluster(req)

	if err != nil {
		return err
	}

	d.SetId(cluster.ID)
	err = d.Set("bootstrap_servers", cluster.Endpoint)
	if err != nil {
		return err
	}

	return nil
}

func clusterDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)
	accountID := d.Get("environment_id").(string)

	return c.DeleteCluster(d.Id(), accountID)
}

func clusterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)
	accountID := d.Get("environment_id").(string)

	cluster, err := c.GetCluster(d.Id(), accountID)
	if err != nil {
		return err
	}

	log.Printf("[WARN] hello %s", cluster.APIEndpoint)
	err = d.Set("bootstrap_servers", cluster.Endpoint)
	if err != nil {
		return err
	}

	return nil
}
