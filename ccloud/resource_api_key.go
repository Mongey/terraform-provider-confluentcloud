package ccloud

import (
	"fmt"
	"log"
	"time"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func apiKeyResource() *schema.Resource {
	return &schema.Resource{
		Create: apiKeyCreate,
		Read:   apiKeyRead,
		Delete: apiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"logical_clusters": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "Logical Cluster ID List to create API Key",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "User ID",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment ID",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description",
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func apiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	clusterID := d.Get("cluster_id").(string)
	logicalClusters := d.Get("logical_clusters").([]interface{})
	accountID := d.Get("environment_id").(string)
	userID := d.Get("user_id").(int)
	description := d.Get("description").(string)

	logicalClustersReq := []ccloud.LogicalCluster{}

	if len(clusterID) > 0 {
		logicalClustersReq = append(logicalClustersReq, ccloud.LogicalCluster{ID: clusterID})
	}

	for i := range logicalClusters {
		if clusterID != logicalClusters[i].(string) {
			logicalClustersReq = append(logicalClustersReq, ccloud.LogicalCluster{
				ID: logicalClusters[i].(string),
			})
		}
	}

	req := ccloud.ApiKeyCreateRequest{
		AccountID:       accountID,
		UserID:          userID,
		LogicalClusters: logicalClustersReq,
		Description:     description,
	}

	log.Printf("[DEBUG] Creating API key")
	key, err := c.CreateAPIKey(&req)
	if err == nil {
		d.SetId(fmt.Sprintf("%d", key.ID))

		err = d.Set("key", key.Key)
		if err != nil {
			return err
		}

		err = d.Set("secret", key.Secret)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Created API Key, waiting for it become usable")
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Pending"},
			Target:       []string{"Ready"},
			Refresh:      clusterReady(c, clusterID, accountID, key.Key, key.Secret),
			Timeout:      300 * time.Second,
			Delay:        10 * time.Second,
			PollInterval: 5 * time.Second,
			MinTimeout:   20 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for API Key (%s) to be ready: %s", d.Id(), err)
		}
	} else {
		log.Printf("[ERROR] Could not create API key: %s", err)
	}

	log.Printf("[INFO] API Key Created successfully: %s", err)
	return err
}

func apiKeyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func apiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	clusterID := d.Get("cluster_id").(string)
	logicalClusters := d.Get("logical_clusters").([]interface{})
	accountID := d.Get("environment_id").(string)

	logicalClustersReq := []ccloud.LogicalCluster{}

	if len(clusterID) > 0 {
		logicalClustersReq = append(logicalClustersReq, ccloud.LogicalCluster{ID: clusterID})
	}

	for i := range logicalClusters {
		if clusterID != logicalClusters[i].(string) {
			logicalClustersReq = append(logicalClustersReq, ccloud.LogicalCluster{
				ID: logicalClusters[i].(string),
			})
		}
	}

	id := d.Id()
	log.Printf("[INFO] Deleting API key %s in account %s", id, accountID)
	err := c.DeleteAPIKey(id, accountID, logicalClustersReq)

	return err
}
