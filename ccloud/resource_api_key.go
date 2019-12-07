package ccloud

import (
	"fmt"
	"log"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
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
				Required:    true,
				ForceNew:    true,
				Description: "",
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
	accountID, err := getAccountID(c)
	if err != nil {
		return err
	}

	req := ccloud.ApiKeyCreateRequest{
		AccountID: accountID,
		LogicalClusters: []ccloud.LogicalCluster{
			ccloud.LogicalCluster{
				ID: clusterID,
			},
		},
	}

	key, err := c.CreateAPIKey(&req)
	if err == nil {
		d.SetId(fmt.Sprintf("%d", key.ID))
		d.Set("key", key.Key)
		d.Set("secret", key.Secret)
	} else {
		log.Printf("[WARN] err creating: %s", err)
	}

	return err
}

func apiKeyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func apiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
