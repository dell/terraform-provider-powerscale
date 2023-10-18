package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccSubnetDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + SubnetDatasourceGetAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_subnet.subnet_datasource_test", "subnets.#"),
				),
			},
		},
	})
}

func TestAccSubnetDatasourceGetFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + SubnetDatasourceGetFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_subnet.subnet_datasource_test", "subnets.#"),
				),
			},
		},
	})
}

func TestAccSubnetDatasourceGetFilterError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config:      ProviderConfig + SubnetDatasourceGetFilterErrorConfig,
				ExpectError: regexp.MustCompile("Error reading subnets"),
			},
		},
	})
}

var SubnetDatasourceGetAllConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
}
`

var SubnetDatasourceGetFilterConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
  filter{
    groupnet_name="groupnet0"
  }
}
`

var SubnetDatasourceGetFilterErrorConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
  filter{
    groupnet_name="groupnet-non-existent"
  }
}
`
