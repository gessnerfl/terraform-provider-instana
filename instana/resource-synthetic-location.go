package instana

import (
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// NewSyntheticLocationDataSource creates a new DataSource for Synthetic Locations
func NewSyntheticLocationDataSource() DataSource {
	return &syntheticLocationDataSource{}
}

// ResourceInstanaSyntheticTest the name of the terraform-provider-instana resource to manage synthetic tests
const ResourceInstanaSyntheticLocation = "instana_synthetic_location"

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

// NewSyntheticTestResourceHandle creates the resource handle Synthetic Tests
func (ds *syntheticLocationDataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		Read: ds.read,
		Schema: map[string]*schema.Schema{
			SyntheticLocationFieldLabel: {
				Type:         schema.TypeString,
				Required:     true,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates if the Synthetic test is started or not",
			},
		},
	}
}

func (ds *syntheticLocationDataSource) read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	label := d.Get(SyntheticLocationFieldLabel).(string)
	locationType := d.Get(SyntheticLocationFieldLocationType).(string)

	data, err := instanaAPI.SyntheticLocationConfig().GetAll()
	log.Printf("DEBUG: All locations: %v", data)
	if err != nil {
		return err
	}

	syntheticLocation, err := ds.findSyntheticLocationByLabelAndType(label, locationType, data)

	if err != nil {
		return err
	}

	return ds.updateState(d, syntheticLocation)
}

func (ds *syntheticLocationDataSource) findSyntheticLocationByLabelAndType(label string, locationType string, data *[]restapi.InstanaDataObject) (*restapi.SyntheticLocation, error) {
	for _, e := range *data {
		syntheticLocation, ok := e.(restapi.SyntheticLocation)
		if ok && syntheticLocation.Label == label && syntheticLocation.LocationType == locationType {
			return &syntheticLocation, nil
		}
	}
	return nil, fmt.Errorf("no synthetic location found for label '%s' and location_type '%s'", label, locationType)
}

func (ds *syntheticLocationDataSource) updateState(d *schema.ResourceData, syntheticLocation *restapi.SyntheticLocation) error {
	d.SetId(syntheticLocation.ID)
	d.Set(SyntheticLocationFieldLabel, syntheticLocation.Label)
	d.Set(SyntheticLocationFieldDescription, syntheticLocation.Description)
	d.Set(SyntheticLocationFieldLocationType, syntheticLocation.LocationType)

	return nil
}
