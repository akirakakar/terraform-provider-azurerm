package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbSqlContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMCosmosDbSqlContainerDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMCosmosDbSqlContainer_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbSqlContainer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMCosmosDbSqlContainerDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMCosmosDbSqlContainer_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlContainerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbSqlContainer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMCosmosDbSqlContainerDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMCosmosDbSqlContainer_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlContainerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_ttl", "500"),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "600"),
				),
			},
			data.ImportStep(),
			{

				Config: testAccAzureRMCosmosDbSqlContainer_update(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlContainerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_ttl", "1000"),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "400"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDbSqlContainerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.SqlClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_sql_container" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		database := rs.Primary.Attributes["database_name"]

		resp, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos SQL Container %s (account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos SQL Container %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDbSqlContainerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.SqlClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		database := rs.Primary.Attributes["database_name"]

		resp, err := client.GetSQLContainer(ctx, resourceGroup, account, database, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos Container '%s' (account: '%s') does not exist", name, account)
		}

		return nil
	}
}

func testAccAzureRMCosmosDbSqlContainer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
}
`, testAccAzureRMCosmosDbSqlDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbSqlContainer_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_path  = "/definition/id"
  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
  default_ttl = 500
  throughput  = 600
}
`, testAccAzureRMCosmosDbSqlDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbSqlContainer_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_path  = "/definition/id"
  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
  default_ttl = 1000
  throughput  = 400
}
`, testAccAzureRMCosmosDbSqlDatabase_basic(data), data.RandomInteger)
}
