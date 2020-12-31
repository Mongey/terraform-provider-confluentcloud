package ccloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	log.Printf("[DEBUG] Created kafka_cluster %s, Endpoint: %s", cluster.ID, cluster.Endpoint)

	err = d.Set("bootstrap_servers", cluster.Endpoint)
	if err != nil {
		return err
	}

	logicalClusters := []ccloud.LogicalCluster{
		ccloud.LogicalCluster{ID: cluster.ID},
	}

	apiKeyReq := ccloud.ApiKeyCreateRequest{
		AccountID:       accountID,
		LogicalClusters: logicalClusters,
		Description:     "terraform-provider-confluentcloud cluster connection bootstrap",
	}

	log.Printf("[DEBUG] Creating bootstrap keypair")
	key, err := c.CreateAPIKey(&apiKeyReq)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Ready"},
		Refresh:      clusterReady(c, d.Id(), accountID, key.Key, key.Secret),
		Timeout:      300 * time.Second,
		Delay:        3 * time.Second,
		PollInterval: 5 * time.Second,
		MinTimeout:   20 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for cluster to become healthy")
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for cluster (%s) to be ready: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Deleting bootstrap keypair")
	err = c.DeleteAPIKey(fmt.Sprintf("%d", key.ID), accountID, logicalClusters)
	if err != nil {
		log.Printf("[ERROR] Unable to delete bootstrap api key %s", err)
	}

	return nil
}

func clusterReady(client *ccloud.Client, clusterID, accountID, username, password string) resource.StateRefreshFunc {
	return func() (result interface{}, s string, err error) {
		cluster, err := client.GetCluster(clusterID, accountID)
		log.Printf("[DEBUG] Waiting for Cluster to be UP: current status %s %s:%s", cluster.Status, username, password)
		log.Printf("[DEBUG] cluster %v", cluster)

		if err != nil {
			return cluster, "UNKNOWN", err
		}

		log.Printf("[DEBUG] Attempting to connect to %s, created %s", cluster.Endpoint, cluster.Deployment.Created)
		if cluster.Status == "UP" {
			if canConnect(cluster.Endpoint, username, password) {
				return cluster, "Ready", nil
			}
		}

		return cluster, "Pending", nil
	}
}

func canConnect(connection, username, password string) bool {
	bootstrapServers := strings.Replace(connection, "SASL_SSL://", "", 1)
	log.Printf("[INFO] Trying to connect to %s", bootstrapServers)

	cfg := sarama.NewConfig()
	cfg.Net.SASL.Enable = true
	cfg.Net.SASL.User = username
	cfg.Net.SASL.Password = password
	cfg.Net.SASL.Handshake = true
	cfg.Net.TLS.Enable = true

	client, err := sarama.NewClient([]string{bootstrapServers}, cfg)
	if err != nil {
		log.Printf("[ERROR] Could not build client %s", err)
		return false
	}

	err = client.RefreshMetadata()
	if err != nil {
		log.Printf("[ERROR] Could not refresh metadata %s", err)
		return false
	}

	log.Printf("[INFO] Success! Connected to %s", bootstrapServers)
	return true
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
		err = d.Set("deployment", map[string]interface{}{"sku": cluster.Deployment.Sku})
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
