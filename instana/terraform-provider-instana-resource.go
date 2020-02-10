package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

//ResourceHandle resource specific implementation which provides meta data and maps data from/to terraform state. Together with TerraformResource terraform schema resources can be created
type ResourceHandle interface {
	GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource
	Schema() map[string]*schema.Schema
	SchemaVersion() int
	StateUpgraders() []schema.StateUpgrader
	ResourceName() string

	UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error
	ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error)
}

//NewTerraformResource creates a new terraform resource for the given handle
func NewTerraformResource(handle ResourceHandle) TerraformResource {
	return &terraformResourceImpl{
		resourceHandle: handle,
	}
}

//TerraformResource internal simplified representation of a Terraform resource
type TerraformResource interface {
	Create(d *schema.ResourceData, meta interface{}) error
	Read(d *schema.ResourceData, meta interface{}) error
	Update(d *schema.ResourceData, meta interface{}) error
	Delete(d *schema.ResourceData, meta interface{}) error
	ToSchemaResource() *schema.Resource
}

type terraformResourceImpl struct {
	resourceHandle ResourceHandle
}

//Create defines the create operation for the terraform resource
func (r *terraformResourceImpl) Create(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return r.Update(d, meta)
}

//Read defines the read operation for the terraform resource
func (r *terraformResourceImpl) Read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	id := d.Id()
	if len(id) == 0 {
		return fmt.Errorf("ID of %s is missing", r.resourceHandle.ResourceName())
	}
	obj, err := r.resourceHandle.GetResourceFrom(instanaAPI).GetOne(id)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	r.resourceHandle.UpdateState(d, obj)
	return nil
}

//Update defines the update operation for the terraform resource
func (r *terraformResourceImpl) Update(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	obj, err := r.resourceHandle.ConvertStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	updatedObject, err := r.resourceHandle.GetResourceFrom(instanaAPI).Upsert(obj)
	if err != nil {
		return err
	}
	r.resourceHandle.UpdateState(d, updatedObject)
	return nil
}

//Delete defines the delete operation for the terraform resource
func (r *terraformResourceImpl) Delete(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	object, err := r.resourceHandle.ConvertStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	err = r.resourceHandle.GetResourceFrom(instanaAPI).DeleteByID(object.GetID())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func (r *terraformResourceImpl) ToSchemaResource() *schema.Resource {
	return &schema.Resource{
		Create:         r.Create,
		Read:           r.Read,
		Update:         r.Update,
		Delete:         r.Delete,
		Schema:         r.resourceHandle.Schema(),
		SchemaVersion:  r.resourceHandle.SchemaVersion(),
		StateUpgraders: r.resourceHandle.StateUpgraders(),
	}
}
