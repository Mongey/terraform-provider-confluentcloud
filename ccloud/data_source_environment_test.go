package ccloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnvironmentDataSourceTest(t *testing.T) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("acc_test_environment-%s", u)

	r.ParallelTest(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: accProvider(),
		Steps: []r.TestStep{
			{
				Config: testAccEnvironmentFilterConfig(name),
				Check:  checkDatasourceDashboardListAttrs(name),
			},
		},
	})
}

func checkDatasourceDashboardListAttrs(uniq string) r.TestCheckFunc {
	return r.ComposeTestCheckFunc(
		r.TestCheckResourceAttr(
			"data.confluentcloud_environment.test", "name", uniq),
		r.TestCheckResourceAttrSet(
			"data.confluentcloud_environment.test", "id"),
	)
}

func testAccEnvironmentFilterConfig(uniq string) string {
	return fmt.Sprintf(`
resource "confluentcloud_environment" "test" {
  name = "%s"
}

data "confluentcloud_environment" "test" {
  name = confluentcloud_environment.test.name
}`, uniq)
}
