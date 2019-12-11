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
	accountID, err := getAccountID(c)

	if err != nil {
		return err
	}

	req := ccloud.ClusterCreateConfig{
		Name:            name,
		Region:          region,
		ServiceProvider: serviceProvider,
		Storage:         5000, // TODO: paramaterize
		AccountID:       accountID,
		Durability:      durability,
	}
	cluster, err := c.CreateCluster(req)

	if err == nil {
		d.SetId(cluster.ID)
	}

	return nil
}

func clusterDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func clusterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)
	accountID, err := getAccountID(c)
	if err != nil {
		return err
	}

	cluster, err := c.GetCluster(d.Id(), accountID)
	if err != nil {
		return err
	}

	log.Printf("[WARN] hello %s", cluster.APIEndpoint)
	d.Set("bootstrap_servers", cluster.Endpoint)

	return nil
}

func nameFromRD(d *schema.ResourceData) string {
	return d.Get("name").(string)
}

func configFromRD(d *schema.ResourceData) map[string]string {
	cfg := d.Get("config").(map[string]interface{})
	scfg := d.Get("config_sensitive").(map[string]interface{})
	config := make(map[string]string)
	for k, v := range cfg {
		config[k] = v.(string)
	}
	for k, v := range scfg {
		config[k] = v.(string)
	}

	return config
}

func getAccountID(client *ccloud.Client) (string, error) {
	err := client.Login()
	if err != nil {
		return "", err
	}

	userData, err := client.Me()
	if err != nil {
		return "", err
	}

	return userData.Account.ID, nil
}
