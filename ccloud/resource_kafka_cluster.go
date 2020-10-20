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
			"storage": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Storage limit(GB)",
			},
			"network_ingress": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Network ingress limit(MBps)",
			},
			"network_egress": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Network egress limit(MBps)",
			},
			"deployment": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Deployment settings.  Currently only `sku` is supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sku": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cku": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "cku",
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
	deployment := d.Get("deployment").(map[string]interface{})
	storage := d.Get("storage").(int)
	networkIngress := d.Get("network_ingress").(int)
	networkEgress := d.Get("network_egress").(int)
	cku := d.Get("cku").(int)

	log.Printf("[DEBUG] Creating kafka_cluster")

	dep := ccloud.ClusterCreateDeploymentConfig{
		AccountID: accountID,
	}

	if val, ok := deployment["sku"]; ok {
		dep.Sku = val.(string)
	} else {
		dep.Sku = "BASIC"
	}

	req := ccloud.ClusterCreateConfig{
		Name:            name,
		Region:          region,
		ServiceProvider: serviceProvider,
		Storage:         storage,
		AccountID:       accountID,
		Durability:      durability,
		Deployment:      dep,
		NetworkIngress:  networkIngress,
		NetworkEgress:   networkEgress,
		Cku:             cku,
	}

	cluster, err := c.CreateCluster(req)
	if err != nil {
		log.Printf("[ERROR] createCluster failed %v, %s", req, err)
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
	if err == nil {
		log.Printf("[WARN] hello %s", cluster.APIEndpoint)
		err = d.Set("bootstrap_servers", cluster.Endpoint)
	}
	if err == nil {
		err = d.Set("name", cluster.Name)
	}
	if err == nil {
		err = d.Set("region", cluster.Region)
	}
	if err == nil {
		err = d.Set("service_provider", cluster.ServiceProvider)
	}
	if err == nil {
		err = d.Set("availability", cluster.Durability)
	}
	if err == nil {
		err = d.Set("deployment", map[string]interface{}{ "sku": cluster.Deployment.Sku })
	}
	if err == nil {
		err = d.Set("storage", cluster.Storage)
	}
	if err == nil {
		err = d.Set("network_ingress", cluster.NetworkIngress)
	}
	if err == nil {
		err = d.Set("network_egress", cluster.NetworkEgress)
	}
	if err == nil {
		err = d.Set("cku", cluster.Cku)
	}
	return err
}
