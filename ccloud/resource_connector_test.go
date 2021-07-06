package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_CreateConnector(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(testResourceConnector_noConfig, u, u, u, u),
			},
			{
				ResourceName: "confluentcloud_connector.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources["confluentcloud_connector.test"]
					if !ok {
						return "", fmt.Errorf("confluentcloud_connector.test not found")
					}

					return rs.Primary.Attributes["environment_id"] + "/" + rs.Primary.Attributes["cluster_id"] + "/" +
						rs.Primary.Attributes["config.name"], nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

//lintignore:AT004
const testResourceConnector_noConfig = `
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

resource "confluentcloud_api_key" "test" {
  environment_id = confluentcloud_environment.test.id
  cluster_id     = confluentcloud_kafka_cluster.test.id
}

resource "confluentcloud_connector" "test" {
  name           = "acc_test_connector-%s"
  environment_id = confluentcloud_environment.test.id
  cluster_id     = confluentcloud_kafka_cluster.test.id

  config = {
    "topic.creation.enable" = "true",
    "kafka.topic"           = "test-datagen-data",
    "connector.class"       = "DatagenSource",
    "name"                  = "acc_test_connector-%s",
    "kafka.api.key"         = confluentcloud_api_key.test.key,
    "kafka.api.secret"      = confluentcloud_api_key.test.secret,
    "output.data.format"    = "JSON",
    "quickstart"            = "USERS",
    "max.interval"          = "10000",
    "tasks.max"             = "1"
  }
}
`
