package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//ResourceMetaData the meta data of a terraform ResourceHandle
type ResourceMetaData struct {
	ResourceName  string
	Schema        map[string]*schema.Schema
	SchemaVersion int
}

//ResourceHandle resource specific implementation which provides meta data and maps data from/to terraform state. Together with TerraformResource terraform schema resources can be created
type ResourceHandle interface {
	//MetaData returns the meta data of this ResourceHandle
	MetaData() *ResourceMetaData
	//StateUpgraders returns the slice of state upgraders used to migrate states from one version to another
	StateUpgraders() []schema.StateUpgrader

	//GetRestResource provides the restapi.RestResource used by the ResourceHandle
	GetRestResource(api restapi.InstanaAPI) restapi.RestResource
	//UpdateState updates the state of the resource provided as schema.ResourceData with the actual data from the Instana API provided as restapi.InstanaDataObject
	UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error
	//MapStateToDataObject maps the current state of the resource provided as schema.ResourceData to the API model of the Instana API represented as an implementation of restapi.InstanaDataObject
	MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error)
	//SetComputedFields calculate and set the calculated value of computed fields of the given resource
	SetComputedFields(d *schema.ResourceData)
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
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	d.SetId(RandomID())
	r.resourceHandle.SetComputedFields(d)

	createRequest, err := r.resourceHandle.MapStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	createdObject, err := r.resourceHandle.GetRestResource(instanaAPI).Create(createRequest)
	if err != nil {
		return err
	}
	r.resourceHandle.UpdateState(d, createdObject)
	return nil
}

//Read defines the read operation for the terraform resource
func (r *terraformResourceImpl) Read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	id := d.Id()
	if len(id) == 0 {
		return fmt.Errorf("ID of %s is missing", r.resourceHandle.MetaData().ResourceName)
	}
	obj, err := r.resourceHandle.GetRestResource(instanaAPI).GetOne(id)
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

	obj, err := r.resourceHandle.MapStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	updatedObject, err := r.resourceHandle.GetRestResource(instanaAPI).Update(obj)
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

	object, err := r.resourceHandle.MapStateToDataObject(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	err = r.resourceHandle.GetRestResource(instanaAPI).DeleteByID(object.GetID())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func (r *terraformResourceImpl) ToSchemaResource() *schema.Resource {
	metaData := r.resourceHandle.MetaData()
	return &schema.Resource{
		Create:         r.Create,
		Read:           r.Read,
		Update:         r.Update,
		Delete:         r.Delete,
		Schema:         metaData.Schema,
		SchemaVersion:  metaData.SchemaVersion,
		StateUpgraders: r.resourceHandle.StateUpgraders(),
	}
}
