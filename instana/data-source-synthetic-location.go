package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// NewSyntheticLocationDataSource creates a new DataSource for Synthetic Locations
func NewSyntheticLocationDataSource() DataSource {
	return &syntheticLocationDataSource{}
}

const (
	//SyntheticLocationFieldLabel constant value for the schema field label
	SyntheticLocationFieldLabel = "label"
	//SyntheticLocationFieldDescription constant value for the computed schema field description
	SyntheticLocationFieldDescription = "description"
	//SyntheticLocationFieldLocationType constant value for the schema field location_type
	SyntheticLocationFieldLocationType = "location_type"
	//DataSourceSyntheticLocation the name of the terraform-provider-instana data sourcefor synthetic location specifications
	DataSourceSyntheticLocation = "instana_synthetic_location"
)

type syntheticLocationDataSource struct{}

// CreateResource creates the resource handle Synthetic Locations
func (ds *syntheticLocationDataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		Read: ds.read,
		Schema: map[string]*schema.Schema{
			SyntheticLocationFieldLabel: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Friendly name of the Synthetic Location",
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			SyntheticLocationFieldDescription: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The description of the Synthetic location",
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			SyntheticLocationFieldLocationType: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Indicates if the location is public or private",
				ValidateFunc: validation.StringInSlice([]string{"Public", "Private"}, true),
			},
		},
	}
}

func (ds *syntheticLocationDataSource) read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	label := d.Get(SyntheticLocationFieldLabel).(string)
	locationType := d.Get(SyntheticLocationFieldLocationType).(string)

	data, err := instanaAPI.SyntheticLocation().GetAll()
	if err != nil {
		return err
	}

	syntheticLocation, err := ds.findSyntheticLocations(label, locationType, data)

	if err != nil {
		return err
	}

	return ds.updateState(d, syntheticLocation)
}

func (ds *syntheticLocationDataSource) findSyntheticLocations(label string, locationType string, data *[]*restapi.SyntheticLocation) (*restapi.SyntheticLocation, error) {
	for _, e := range *data {
		if e.Label == label && e.LocationType == locationType {
			return e, nil
		}
	}
	return nil, fmt.Errorf("no synthetic location found")
}

func (ds *syntheticLocationDataSource) updateState(d *schema.ResourceData, syntheticLocation *restapi.SyntheticLocation) error {
	d.SetId(syntheticLocation.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		SyntheticLocationFieldLabel:        syntheticLocation.Label,
		SyntheticLocationFieldDescription:  syntheticLocation.Description,
		SyntheticLocationFieldLocationType: syntheticLocation.LocationType,
	})
}
