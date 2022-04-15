package instana_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"
)

func TestGlobalApplicationAlertConfig(t *testing.T) {
	commonTests := createApplicationAlertConfigTestFor("instana_global_application_alert_config", restapi.GlobalApplicationAlertConfigsResourcePath, NewGlobalApplicationAlertConfigResourceHandle())
	commonTests.run(t)
}
