package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_SchemaRegistry(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(testResourceSchemaRegistry_noConfig, u, u),
			},
			{
				ResourceName: "confluentcloud_schema_registry.test",
				ImportState:  true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					resources := state.RootModule().Resources
					schemaRegistryID := resources["confluentcloud_schema_registry.test"].Primary.ID
					envID := resources["confluentcloud_environment.test"].Primary.ID
					return (envID + "/" + schemaRegistryID), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
					"service_provider",
				},
			},
		},
	})
}

//lintignore:AT004
const testResourceSchemaRegistry_noConfig = `
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

resource "confluentcloud_schema_registry" "test" {
  environment_id   = confluentcloud_environment.test.id
  service_provider = "aws"
  region           = "EU"

  # Requires at least one kafka cluster to enable the schema registry in the environment.
  depends_on = [confluentcloud_kafka_cluster.test]
}
`
