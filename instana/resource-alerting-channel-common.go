package instana

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelFieldName constant value for the schema field name
	AlertingChannelFieldName = "name"
	//AlertingChannelFieldFullName constant value for the schema field full_name
	AlertingChannelFieldFullName = "full_name"
)

var alertingChannelNameSchemaField = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Configures the name of the alerting channel",
}

var alertingChannelFullNameSchemaField = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The the full name field of the alerting channel. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}

func migrateFullNameToName(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[AlertingChannelFieldFullName]; ok {
		state[AlertingChannelFieldName] = state[AlertingChannelFieldFullName]
		delete(state, AlertingChannelFieldFullName)
	}
	return state, nil
}
