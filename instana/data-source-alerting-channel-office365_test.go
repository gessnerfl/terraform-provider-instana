package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const dataSourceAlertingChannelOffice365Definition = `
data "instana_alerting_channel_office_365" "example" {
  name = "my-office-356-channel"
}
`
const dataSourceAlertingChannelOffice365DefinitionPath = "data.instana_alerting_channel_office_365.example"

const alertingChannelOffice365ServerResponseTemplate = `
[{
	"id"     	 : "12345",
	"name"   	 : "other1",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook url"
},{
	"id"     	 : "23456",
	"name"   	 : "my-office-356-channel",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook url"
},{
	"id"     	 : "34567",
	"name"   	 : "other2",
	"kind"   	 : "OFFICE_365",
	"webhookUrl" : "webhook url"
}]
`

func TestReadOfAlertingChannelOffice365DataSourceWithMockServer(t *testing.T) {
	httpServer := createMockHttpServerForDataSource(restapi.AlertingChannelsResourcePath, newStringContentResponseProvider(alertingChannelOffice365ServerResponseTemplate))
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAlertingChannelOffice356DateSourceTestStep(httpServer.GetPort()),
		},
	})
}

func createAlertingChannelOffice356DateSourceTestStep(httpPort int64) resource.TestStep {
	config := appendProviderConfig(dataSourceAlertingChannelOffice365Definition, httpPort)
	return resource.TestStep{
		Config: config,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(dataSourceAlertingChannelOffice365DefinitionPath, "id", "23456"),
			resource.TestCheckResourceAttr(dataSourceAlertingChannelOffice365DefinitionPath, AlertingChannelFieldName, "my-office-356-channel"),
		),
	}
}
