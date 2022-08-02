package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServiceAccountDataSourceTest(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("acc_test-%s", u)

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: testAccServiceAccountFilterConfig(name),
				Check:  checkServiceAccountDatasourceDashboardListAttrs(name),
			},
		},
	})
}

func checkServiceAccountDatasourceDashboardListAttrs(uniq string) r.TestCheckFunc {
	return r.ComposeTestCheckFunc(
		r.TestCheckResourceAttr(
			"data.confluentcloud_service_account.test", "name", uniq),
		r.TestCheckResourceAttrSet(
			"data.confluentcloud_service_account.test", "id"),
	)
}

func testAccServiceAccountFilterConfig(uniq string) string {
	return fmt.Sprintf(`
resource "confluentcloud_service_account" "test" {
  name = "%s"
  description    = "service account acc test"
}

data "confluentcloud_service_account" "test" {
  name = confluentcloud_service_account.test.name
}`, uniq)
}
