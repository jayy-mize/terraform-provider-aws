package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/atest"
)

func TestAccDataSourceAwsOrganizationsDelegatedAdministrators_basic(t *testing.T) {
	var providers []*schema.Provider
	dataSourceName := "data.aws_organizations_delegated_administrators.test"
	servicePrincipal := "config-multiaccountsetup.amazonaws.com"
	dataSourceIdentity := "data.aws_caller_identity.delegated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			atest.PreCheck(t)
			atest.PreCheckAlternateAccount(t)
		},
		ErrorCheck:        atest.ErrorCheck(t, organizations.EndpointsID),
		ProviderFactories: atest.ProviderFactoriesAlternate(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsOrganizationsDelegatedAdministratorsConfig(servicePrincipal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "delegated_administrators.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delegated_administrators.0.id", dataSourceIdentity, "account_id"),
					atest.CheckAttrRfc3339(dataSourceName, "delegated_administrators.0.delegation_enabled_date"),
					atest.CheckAttrRfc3339(dataSourceName, "delegated_administrators.0.joined_timestamp"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsOrganizationsDelegatedAdministrators_multiple(t *testing.T) {
	var providers []*schema.Provider
	dataSourceName := "data.aws_organizations_delegated_administrators.test"
	servicePrincipal := "config-multiaccountsetup.amazonaws.com"
	servicePrincipal2 := "config.amazonaws.com"
	dataSourceIdentity := "data.aws_caller_identity.delegated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			atest.PreCheck(t)
			atest.PreCheckAlternateAccount(t)
		},
		ErrorCheck:        atest.ErrorCheck(t, organizations.EndpointsID),
		ProviderFactories: atest.ProviderFactoriesAlternate(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsOrganizationsDelegatedAdministratorsMultipleConfig(servicePrincipal, servicePrincipal2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "delegated_administrators.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delegated_administrators.0.id", dataSourceIdentity, "account_id"),
					atest.CheckAttrRfc3339(dataSourceName, "delegated_administrators.0.delegation_enabled_date"),
					atest.CheckAttrRfc3339(dataSourceName, "delegated_administrators.0.joined_timestamp"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsOrganizationsDelegatedAdministrators_servicePrincipal(t *testing.T) {
	var providers []*schema.Provider
	dataSourceName := "data.aws_organizations_delegated_administrators.test"
	servicePrincipal := "config-multiaccountsetup.amazonaws.com"
	dataSourceIdentity := "data.aws_caller_identity.delegated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			atest.PreCheck(t)
			atest.PreCheckAlternateAccount(t)
		},
		ErrorCheck:        atest.ErrorCheck(t, organizations.EndpointsID),
		ProviderFactories: atest.ProviderFactoriesAlternate(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsOrganizationsDelegatedAdministratorsServicePrincipalConfig(servicePrincipal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "delegated_administrators.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delegated_administrators.0.id", dataSourceIdentity, "account_id"),
					atest.CheckAttrRfc3339(dataSourceName, "delegated_administrators.0.delegation_enabled_date"),
					atest.CheckAttrRfc3339(dataSourceName, "delegated_administrators.0.joined_timestamp"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsOrganizationsDelegatedAdministrators_empty(t *testing.T) {
	dataSourceName := "data.aws_organizations_delegated_administrators.test"
	servicePrincipal := "config-multiaccountsetup.amazonaws.com"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { atest.PreCheck(t) },
		ErrorCheck:        atest.ErrorCheck(t, organizations.EndpointsID),
		ProviderFactories: atest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsOrganizationsDelegatedAdministratorsEmptyConfig(servicePrincipal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "delegated_administrators.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAwsOrganizationsDelegatedAdministratorsEmptyConfig(servicePrincipal string) string {
	return atest.ConfigProviderAlternateAccount() + fmt.Sprintf(`
data "aws_organizations_delegated_administrators" "test" {
  service_principal = %[1]q
}
`, servicePrincipal)
}

func testAccDataSourceAwsOrganizationsDelegatedAdministratorsConfig(servicePrincipal string) string {
	return atest.ConfigProviderAlternateAccount() + fmt.Sprintf(`
data "aws_caller_identity" "delegated" {
  provider = "awsalternate"
}

resource "aws_organizations_delegated_administrator" "test" {
  account_id        = data.aws_caller_identity.delegated.account_id
  service_principal = %[1]q
}

data "aws_organizations_delegated_administrators" "test" {}
`, servicePrincipal)
}

func testAccDataSourceAwsOrganizationsDelegatedAdministratorsMultipleConfig(servicePrincipal, servicePrincipal2 string) string {
	return atest.ConfigProviderAlternateAccount() + fmt.Sprintf(`
data "aws_caller_identity" "delegated" {
  provider = "awsalternate"
}

resource "aws_organizations_delegated_administrator" "delegated" {
  account_id        = data.aws_caller_identity.delegated.account_id
  service_principal = %[1]q
}

resource "aws_organizations_delegated_administrator" "other_delegated" {
  account_id        = data.aws_caller_identity.delegated.account_id
  service_principal = %[2]q
}

data "aws_organizations_delegated_administrators" "test" {}
`, servicePrincipal, servicePrincipal2)
}

func testAccDataSourceAwsOrganizationsDelegatedAdministratorsServicePrincipalConfig(servicePrincipal string) string {
	return atest.ConfigProviderAlternateAccount() + fmt.Sprintf(`
data "aws_caller_identity" "delegated" {
  provider = "awsalternate"
}

resource "aws_organizations_delegated_administrator" "test" {
  account_id        = data.aws_caller_identity.delegated.account_id
  service_principal = %[1]q
}

data "aws_organizations_delegated_administrators" "test" {
  service_principal = aws_organizations_delegated_administrator.test.service_principal
}
`, servicePrincipal)
}
