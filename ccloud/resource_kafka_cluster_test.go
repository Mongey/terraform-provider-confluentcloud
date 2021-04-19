package ccloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func overrideProvider() *schema.Provider {
	return Provider()
}
func accProvider() map[string]*schema.Provider {
	return map[string]*schema.Provider{
		"confluentcloud": overrideProvider(),
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("CONFLUENT_CLOUD_USERNAME") == "" || os.Getenv("CONFLUENT_CLOUD_PASSWORD") == "" {
		t.Skip("CONFLUENT_CLOUD_ environment variables must be set")
	}
}

func TestAcc_BasicCluster(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(testResourceCluster_noConfig, u, u),
			},
			{
				ResourceName: "confluentcloud_kafka_cluster.test",
				ImportState: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.RootModule().Resources
					clusterId := resources["confluentcloud_kafka_cluster.test"].Primary.ID
					envId := resources["confluentcloud_environment.test"].Primary.ID
					return envId + ":" + clusterId, nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

//lintignore:AT004
const testResourceCluster_noConfig = `
resource "confluentcloud_environment" "test" {
  name = "acc_test_environment-%s"
}

resource "confluentcloud_kafka_cluster" "test" {
  name             = "provider-test-%s"
  service_provider = "aws"
  region           = "eu-west-1"
  availability     = "LOW"
  environment_id   = confluentcloud_environment.test.id
  deployment = {
    sku = "BASIC"
  }
  network_egress  = 100
  network_ingress = 100
  storage         = 5000
}
`
