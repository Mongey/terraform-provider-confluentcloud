package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/terraform"
)

func TestAcc_BasicTopic(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	clusterName := fmt.Sprintf("accTest-%s", u)
	environmentName := "default"

	r.Test(t, r.TestCase{
		Providers:    accProvider(),
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckTopicDestroy,
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(testResourceCluster, clusterName, environmentName),
				Check:  testResourceTopic_noConfigCheck,
			},
		},
	})
}

func testResourceCluster_noConfigCheck(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["confluentcloud_kafka_cluster.test"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	name := instanceState.ID
	bootstrapServers := instanceState["bootstrap_servers"]
	// TODO:
	// 0. Create a api key pair
	// 0. Build a client from the connection
	// 0. refresh metadata?
	if name != instanceState.Attributes["name"] {
		return fmt.Errorf("id doesn't match name")
	}

	client := testProvider.Meta().(*LazyClient)
	topic, err := client.ReadTopic(name)
	if err != nil {
		return err
	}

	if len(topic.Config) != 0 {
		return fmt.Errorf("expected no configs for %s, got %v", name, topic.Config)
	}

	return nil
}

const testResourceCluster = `
resource "confluentcloud_kafka_cluster" "test" {
  name             = "%s"
  service_provider = "aws"
  region           = "eu-west-1"
  availability     = "LOW"
  environment_id   = "%s"
  deployment = {
    sku = "BASIC"
  }
  network_egress  = 1
  network_ingress = 1
  storage         = 5
}
`
