package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"
)

func TestApplicationAlertConfig(t *testing.T) {
	commonTests := createApplicationAlertConfigTestFor("instana_application_alert_config", restapi.ApplicationAlertConfigsResourcePath, NewApplicationAlertConfigResourceHandle())
	commonTests.run(t)
}
