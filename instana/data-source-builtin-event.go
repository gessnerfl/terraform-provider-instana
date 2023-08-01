package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewBuiltinEventDataSource creates a new DataSource for Builtin Events
func NewBuiltinEventDataSource() DataSource {
	return &builtInEventDataSource{}
}

const (
	//BuiltinEventSpecificationFieldName constant value for the schema field name
	BuiltinEventSpecificationFieldName = "name"
	//BuiltinEventSpecificationFieldDescription constant value for the schema field description
	BuiltinEventSpecificationFieldDescription = "description"
	//BuiltinEventSpecificationFieldShortPluginID constant value for the schema field short_plugin_id
	BuiltinEventSpecificationFieldShortPluginID = "short_plugin_id"
	//BuiltinEventSpecificationFieldSeverity constant value for the schema field severity
	BuiltinEventSpecificationFieldSeverity = "severity"
	//BuiltinEventSpecificationFieldSeverityCode constant value for the schema field severity_code
	BuiltinEventSpecificationFieldSeverityCode = "severity_code"
	//BuiltinEventSpecificationFieldTriggering constant value for the schema field triggering
	BuiltinEventSpecificationFieldTriggering = "triggering"
	//BuiltinEventSpecificationFieldEnabled constant value for the schema field enabled
	BuiltinEventSpecificationFieldEnabled = "enabled"

	//DataSourceBuiltinEvent the name of the terraform-provider-instana data sourcefor builtin event specifications
	DataSourceBuiltinEvent = "instana_builtin_event_spec"
	//
)

type builtInEventDataSource struct{}

// CreateResource creates the terraform Resource for the data source for Instana builtin events
func (ds *builtInEventDataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		Read: ds.read,
		Schema: map[string]*schema.Schema{
			BuiltinEventSpecificationFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the builtin event",
			},
			BuiltinEventSpecificationFieldDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description text of the builtin event.",
			},
			BuiltinEventSpecificationFieldShortPluginID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plugin id for which the builtin event is created.",
			},
			BuiltinEventSpecificationFieldSeverity: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The severity (WARNING, CRITICAL, etc.) of the builtin event.",
			},
			BuiltinEventSpecificationFieldSeverityCode: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The severity code used by Instana API (5, 10, etc.) of the builtin event.",
			},
			BuiltinEventSpecificationFieldTriggering: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if an incident is triggered the builtin event or not.",
			},
			BuiltinEventSpecificationFieldEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the builtin event is enabled or not",
			},
		},
	}
}

func (ds *builtInEventDataSource) read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	name := d.Get(BuiltinEventSpecificationFieldName).(string)
	shortPluginID := d.Get(BuiltinEventSpecificationFieldShortPluginID).(string)

	data, err := instanaAPI.BuiltinEventSpecifications().GetAll()
	if err != nil {
		return err
	}

	builtInEvent, err := ds.findBuiltInEventByNameAndPluginID(name, shortPluginID, data)

	if err != nil {
		return err
	}

	return ds.updateState(d, builtInEvent)
}

func (ds *builtInEventDataSource) findBuiltInEventByNameAndPluginID(name string, shortPluginID string, data *[]*restapi.BuiltinEventSpecification) (*restapi.BuiltinEventSpecification, error) {
	for _, builtInEvent := range *data {
		if builtInEvent.Name == name && builtInEvent.ShortPluginID == shortPluginID {
			return builtInEvent, nil
		}
	}
	return nil, fmt.Errorf("no built in event found for name '%s' and short plugin ID '%s'", name, shortPluginID)
}

func (ds *builtInEventDataSource) updateState(d *schema.ResourceData, builtInEvent *restapi.BuiltinEventSpecification) error {
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(builtInEvent.Severity)
	if err != nil {
		return err
	}
	d.SetId(builtInEvent.ID)
	d.Set(BuiltinEventSpecificationFieldDescription, builtInEvent.Description)
	d.Set(BuiltinEventSpecificationFieldSeverity, severity)
	d.Set(BuiltinEventSpecificationFieldSeverityCode, builtInEvent.Severity)
	d.Set(BuiltinEventSpecificationFieldTriggering, builtInEvent.Triggering)
	d.Set(BuiltinEventSpecificationFieldEnabled, builtInEvent.Enabled)
	return nil
}
