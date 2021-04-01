package ccloud

import (
	"context"
	"fmt"
	"log"
	"time"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func apiKeyResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: apiKeyCreate,
		ReadContext:   apiKeyRead,
		DeleteContext: apiKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"target_resource_type": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default: "kafka_cluster",
				Description: "Type of the resource to which the api keys are created for.",
				ValidateFunc: validation.StringInSlice([]string { "kafka_cluster", "schema_registry" }, false),
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

func apiKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ccloud.Client)

	clusterID := d.Get("cluster_id").(string)
	logicalClusters := d.Get("logical_clusters").([]interface{})
	accountID := d.Get("environment_id").(string)
	userID := d.Get("user_id").(int)
	description := d.Get("description").(string)
	targetResourceType := d.Get("target_resource_type").(string)

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
			return diag.FromErr(err)
		}

		err = d.Set("secret", key.Secret)
		if err != nil {
			return diag.FromErr(err)
		}

		if targetResourceType == "kafka_cluster" {
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

			_, err = stateConf.WaitForStateContext(context.Background())
			if err != nil {
				return diag.FromErr(fmt.Errorf("Error waiting for API Key (%s) to be ready: %s", d.Id(), err))
			}
		} else {
			log.Print("[INFO] Not an api key for a Kafka cluster. Do not wait when it becomes usable.")
		}
	} else {
		log.Printf("[ERROR] Could not create API key: %s", err)
	}

	log.Printf("[INFO] API Key Created successfully: %s", err)
	return diag.FromErr(err)
}

func apiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func apiKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return diag.FromErr(err)
}
