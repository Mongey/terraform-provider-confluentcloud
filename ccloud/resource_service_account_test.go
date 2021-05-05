package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_BasicServiceAccount(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(testServiceAccount_noConfig, u, u),
			},
			{
				ResourceName:      "confluentcloud_service_account.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//lintignore:AT004
const testServiceAccount_noConfig = `
resource "confluentcloud_service_account" "test" {
	name        = "acc-test-%s"
	description = "My cool description - %s"
}
`
