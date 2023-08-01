package instana

import (
	"context"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceMetaData the meta data of a terraform ResourceHandle
type ResourceMetaData struct {
	ResourceName     string
	Schema           map[string]*schema.Schema
	SchemaVersion    int
	SkipIDGeneration bool
	ResourceIDField  *string
}

// ResourceHandle resource specific implementation which provides meta data and maps data from/to terraform state. Together with TerraformResource terraform schema resources can be created
type ResourceHandle[T restapi.InstanaDataObject] interface {
	//MetaData returns the meta data of this ResourceHandle
	MetaData() *ResourceMetaData
	//StateUpgraders returns the slice of state upgraders used to migrate states from one version to another
	StateUpgraders() []schema.StateUpgrader

	//GetRestResource provides the restapi.RestResource used by the ResourceHandle
	GetRestResource(api restapi.InstanaAPI) restapi.RestResource[T]
	//UpdateState updates the state of the resource provided as schema.ResourceData with the input data from the Instana API provided as restapi.InstanaDataObject
	UpdateState(d *schema.ResourceData, obj T, formatter utils.ResourceNameFormatter) error
	//MapStateToDataObject maps the current state of the resource provided as schema.ResourceData to the API model of the Instana API represented as an implementation of restapi.InstanaDataObject
	MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (T, error)
	//SetComputedFields calculate and set the calculated value of computed fields of the given resource
	SetComputedFields(d *schema.ResourceData)
}

// NewTerraformResource creates a new terraform resource for the given handle
func NewTerraformResource[T restapi.InstanaDataObject](handle ResourceHandle[T]) TerraformResource {
	return &terraformResourceImpl[T]{
		resourceHandle: handle,
	}
}

// TerraformResource internal simplified representation of a Terraform resource
type TerraformResource interface {
	Create(d *schema.ResourceData, meta interface{}) error
	Read(d *schema.ResourceData, meta interface{}) error
	Update(d *schema.ResourceData, meta interface{}) error
	Delete(d *schema.ResourceData, meta interface{}) error
	ToSchemaResource() *schema.Resource
}

type terraformResourceImpl[T restapi.InstanaDataObject] struct {
	resourceHandle ResourceHandle[T]
}

// Create defines the create operation for the terraform resource
func (r *terraformResourceImpl[T]) Create(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	if !r.resourceHandle.MetaData().SkipIDGeneration {
		d.SetId(RandomID())
	}
	r.resourceHandle.SetComputedFields(d)

	createRequest, err := r.resourceHandle.MapStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	createdObject, err := r.resourceHandle.GetRestResource(instanaAPI).Create(createRequest)
	if err != nil {
		return err
	}
	r.resourceHandle.UpdateState(d, createdObject, providerMeta.ResourceNameFormatter)
	return nil
}

// Read defines the read operation for the terraform resource
func (r *terraformResourceImpl[T]) Read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	resourceID := r.getResourceID(d)
	if len(resourceID) == 0 {
		return fmt.Errorf("resource ID of %s is missing", r.resourceHandle.MetaData().ResourceName)
	}
	obj, err := r.resourceHandle.GetRestResource(instanaAPI).GetOne(resourceID)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	return r.resourceHandle.UpdateState(d, obj, providerMeta.ResourceNameFormatter)
}

func (r *terraformResourceImpl[T]) getResourceID(d *schema.ResourceData) string {
	if r.resourceHandle.MetaData().ResourceIDField != nil {
		return d.Get(*r.resourceHandle.MetaData().ResourceIDField).(string)
	}
	return d.Id()
}

// Update defines the update operation for the terraform resource
func (r *terraformResourceImpl[T]) Update(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	obj, err := r.resourceHandle.MapStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	updatedObject, err := r.resourceHandle.GetRestResource(instanaAPI).Update(obj)
	if err != nil {
		return err
	}
	return r.resourceHandle.UpdateState(d, updatedObject, providerMeta.ResourceNameFormatter)
}

// Delete defines the delete operation for the terraform resource
func (r *terraformResourceImpl[T]) Delete(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	object, err := r.resourceHandle.MapStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	err = r.resourceHandle.GetRestResource(instanaAPI).DeleteByID(object.GetIDForResourcePath())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func (r *terraformResourceImpl[T]) ToSchemaResource() *schema.Resource {
	metaData := r.resourceHandle.MetaData()
	return &schema.Resource{
		Create: r.Create,
		Read:   r.Read,
		Importer: &schema.ResourceImporter{
			StateContext: r.importState,
		},
		Update:         r.Update,
		Delete:         r.Delete,
		Schema:         metaData.Schema,
		SchemaVersion:  metaData.SchemaVersion,
		StateUpgraders: r.resourceHandle.StateUpgraders(),
	}
}

func (r *terraformResourceImpl[T]) importState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if r.resourceHandle.MetaData().ResourceIDField != nil {
		d.Set(*r.resourceHandle.MetaData().ResourceIDField, d.Id())
	}
	return []*schema.ResourceData{d}, nil
}
