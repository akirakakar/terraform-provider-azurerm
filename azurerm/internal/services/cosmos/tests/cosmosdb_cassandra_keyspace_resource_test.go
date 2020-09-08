package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbCassandraKeyspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_keyspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMCosmosDbCassandraKeyspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbCassandraKeyspace_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbCassandraKeyspaceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbCassandraKeyspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_keyspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMCosmosDbCassandraKeyspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbCassandraKeyspace_throughput(data, 700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbCassandraKeyspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "700"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbCassandraKeyspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_keyspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.SupportedProviders,
		CheckDestroy:      testCheckAzureRMCosmosDbCassandraKeyspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbCassandraKeyspace_throughput(data, 700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbCassandraKeyspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "700"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbCassandraKeyspace_throughput(data, 1700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbCassandraKeyspaceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "1700"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDbCassandraKeyspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.CassandraClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_cassandra_keyspace" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Cassandra Keyspace %s (account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Cassandra Keyspace %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDbCassandraKeyspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.CassandraClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s') does not exist", name, account)
		}

		return nil
	}
}

func testAccAzureRMCosmosDbCassandraKeyspace_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_keyspace" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra"}), data.RandomInteger)
}

func testAccAzureRMCosmosDbCassandraKeyspace_throughput(data acceptance.TestData, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_keyspace" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name

  throughput = %[3]d
}
`, testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra"}), data.RandomInteger, throughput)
}
